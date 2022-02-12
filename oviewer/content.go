package oviewer

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/uniseg"
)

// content represents one character on the terminal.
// content is a value that can be set in SetContent() of tcell.
type content struct {
	width int
	mainc rune
	combc []rune
	style tcell.Style
}

// lineContents represents one line of contents.
type lineContents []content

// The states of the ANSI escape code parser.
const (
	ansiText = iota
	ansiEscape
	ansiSubstring
	ansiControlSequence
)

// DefaultContent is a blank Content.
var DefaultContent = content{
	mainc: 0,
	combc: nil,
	width: 0,
	style: tcell.StyleDefault,
}

// EOFContent is "~" only.
var EOFContent = content{
	mainc: '~',
	combc: nil,
	width: 1,
	style: tcell.StyleDefault.Foreground(tcell.ColorGray),
}

var csiCache sync.Map

// parseString converts a string to lineContents.
// parseString includes escape sequences and tabs.
func parseString(str string, tabWidth int) lineContents {
	lc := make(lineContents, 0, len(str))
	state := ansiText
	csiParameter := new(bytes.Buffer)
	style := tcell.StyleDefault
	tabX := 0
	b := 0
	bsFlag := false // backspace(^H) flag
	var bsContent content

	gr := uniseg.NewGraphemes(str)
	for gr.Next() {
		runeValue := gr.Runes()[0]
		c := DefaultContent
		switch state {
		case ansiEscape:
			switch runeValue {
			case '[': // Control Sequence Introducer.
				csiParameter.Reset()
				state = ansiControlSequence
				continue
			case 'c': // Reset.
				style = tcell.StyleDefault
				state = ansiText
				continue
			case 'P', ']', 'X', '^', '_': // Substrings and commands.
				state = ansiSubstring
				continue
			default: // Ignore.
				state = ansiText
			}
		case ansiSubstring:
			if runeValue == 0x1b {
				state = ansiControlSequence
			}
			continue
		case ansiControlSequence:
			if runeValue == 'm' {
				style = csToStyle(style, csiParameter.String())
			} else if runeValue >= 'A' && runeValue <= 'T' {
				// Ignore.
			} else {
				if runeValue >= 0x30 && runeValue <= 0x3f {
					csiParameter.WriteRune(runeValue)
					continue
				}
			}
			state = ansiText
			continue
		}

		switch runeValue {
		case 0x1b:
			state = ansiEscape
			continue
		case '\n':
			continue
		}

		switch runewidth.RuneWidth(runeValue) {
		case 0:
			switch runeValue {
			case '\t': // TAB
				switch {
				case tabWidth > 0:
					tabStop := tabWidth - (tabX % tabWidth)
					c.width = 1
					c.style = style
					c.mainc = rune('\t')
					lc = append(lc, c)
					tabX++
					c.mainc = 0
					for i := 0; i < tabStop-1; i++ {
						lc = append(lc, c)
						tabX++
					}
				case tabWidth < 0:
					c.width = 1
					c.style = style.Reverse(true)
					c.mainc = rune('\\')
					lc = append(lc, c)
					c.mainc = rune('t')
					lc = append(lc, c)
					tabX += 2
				default:
				}
				continue
			case '\b': // BackSpace
				if len(lc) == 0 {
					continue
				}
				bsFlag = true
				bsContent = lastContent(lc)
				b -= (1 + len(string(bsContent.mainc)))
				if bsContent.width > 1 {
					lc = lc[:len(lc)-2]
				} else {
					lc = lc[:len(lc)-1]
				}
				continue
			}
			lastC := lastContent(lc)
			lastC.combc = append(lastC.combc, runeValue)
			n := len(lc) - lastC.width
			if n >= 0 && len(lc) > 0 {
				lc[n] = lastC
			}
		case 1:
			c.mainc = runeValue
			if len(gr.Runes()) > 1 {
				c.combc = gr.Runes()[1:]
			}
			c.width = 1
			c.style = style
			if bsFlag {
				c.style = overstrike(bsContent.mainc, runeValue, style)
				bsFlag = false
				bsContent = DefaultContent
			}
			lc = append(lc, c)
			tabX++
		case 2:
			c.mainc = runeValue
			if len(gr.Runes()) > 1 {
				c.combc = gr.Runes()[1:]
			}
			c.width = 2
			c.style = style
			if bsFlag {
				c.style = overstrike(bsContent.mainc, runeValue, style)
				bsFlag = false
				bsContent = DefaultContent
			}
			lc = append(lc, c, DefaultContent)
			tabX += 2
		}
	}
	return lc
}

// overstrike returns an overstrike tcell.Style.
func overstrike(p, m rune, style tcell.Style) tcell.Style {
	if p == m {
		style = OverStrikeStyle
	} else if p == '_' {
		style = OverLineStyle
	}
	return style
}

// lastContent returns the last character of Contents.
func lastContent(lc lineContents) content {
	n := len(lc)
	if n == 0 {
		return content{}
	}
	if (n > 1) && (lc[n-2].width > 1) {
		return lc[n-2]
	}
	return lc[n-1]
}

// csToStyle returns tcell.Style from the control sequence.
func csToStyle(style tcell.Style, params string) tcell.Style {
	if params == "0" || params == "" || params == ";" {
		return tcell.StyleDefault.Normal()
	}
	if s, ok := csiCache.Load(params); ok {
		style = applyStyle(style, s.(ovStyle))
		return style
	}

	s := parseCSI(params)
	csiCache.Store(params, s)
	style = applyStyle(style, s)
	return style
}

func parseCSI(params string) ovStyle {
	style := ovStyle{}
	fields := strings.Split(params, ";")
	fl := len(fields)
	for index := 0; index < fl; index++ {
		field := fields[index]
		switch field {
		case "1", "01":
			style.Bold = true
		case "2", "02":
			style.Dim = true
		case "3", "03":
			style.Italic = true
		case "4", "04":
			style.Underline = true
		case "5", "05":
			style.Blink = true
		case "6", "06":
			style.Blink = true
		case "7", "07":
			style.Reverse = true
		case "8", "08":
			style.Reverse = true
		case "9", "09":
			style.StrikeThrough = true
		case "22", "24", "25", "27":
			style = ovStyle{}
		case "30", "31", "32", "33", "34", "35", "36", "37":
			colorNumber, _ := strconv.Atoi(field)
			style.Foreground = colorName(int(tcell.Color(colorNumber - 30)))
		case "39":
			style.Foreground = "default"
		case "40", "41", "42", "43", "44", "45", "46", "47":
			colorNumber, _ := strconv.Atoi(field)
			style.Background = colorName(int(tcell.Color(colorNumber - 40)))
		case "49":
			style.Background = "default"
		case "90", "91", "92", "93", "94", "95", "96", "97":
			colorNumber, _ := strconv.Atoi(field)
			style.Foreground = colorName(int(tcell.Color(colorNumber - 82)))
		case "100", "101", "102", "103", "104", "105", "106", "107":
			colorNumber, _ := strconv.Atoi(field)
			style.Background = colorName(int(tcell.Color(colorNumber - 92)))
		case "38", "48":
			var color string
			var color2 string
			var field2 string
			if len(fields) > index+1 {
				index++
				if fields[index] == "5" && len(fields) > index+1 { // 8-bit colors.
					index++
					c, _ := strconv.Atoi(fields[index])
					color = colorName(c)
					// both
					if len(fields) > index+3 {
						index++
						field2 = fields[index]
						index++
						// fields[index] == 5
						index++
						c, _ := strconv.Atoi(fields[index])
						color2 = colorName(c)
					}
				} else if fields[index] == "2" && len(fields) > index+3 { // 24-bit colors.
					index++
					red, _ := strconv.Atoi(fields[index])
					index++
					green, _ := strconv.Atoi(fields[index])
					index++
					blue, _ := strconv.Atoi(fields[index])
					color = fmt.Sprintf("#%02x%02x%02x", red, green, blue)
					// both
					if len(fields) > index+5 {
						index++
						field2 = fields[index]
						index++
						// fields[index] == 2
						index++
						red, _ := strconv.Atoi(fields[index])
						index++
						green, _ := strconv.Atoi(fields[index])
						index++
						blue, _ := strconv.Atoi(fields[index])
						color2 = fmt.Sprintf("#%02x%02x%02x", red, green, blue)
					}
				}
			}
			if len(color) > 0 {
				if field == "38" {
					style.Foreground = color
				} else {
					style.Background = color
				}
				// both
				if len(color2) > 0 {
					if field2 == "38" {
						style.Foreground = color2
					} else {
						style.Background = color2
					}
				}
			}
		}
	}
	return style
}

func colorName(colorNumber int) string {
	if colorNumber <= 7 {
		return lookupColor(colorNumber)
	} else if colorNumber <= 15 {
		return lookupColor(colorNumber)
	} else if colorNumber <= 231 {
		red := (colorNumber - 16) / 36
		green := ((colorNumber - 16) / 6) % 6
		blue := (colorNumber - 16) % 6
		return fmt.Sprintf("#%02x%02x%02x", 255*red/5, 255*green/5, 255*blue/5)
	} else if colorNumber <= 255 {
		grey := 255 * (colorNumber - 232) / 23
		return fmt.Sprintf("#%02x%02x%02x", grey, grey, grey)
	}
	return ""
}

// lookupColor returns the color name from the color number.
func lookupColor(colorNumber int) string {
	if colorNumber < 0 || colorNumber > 15 {
		return "black"
	}
	return [...]string{
		"black",
		"maroon",
		"green",
		"olive",
		"navy",
		"purple",
		"teal",
		"silver",
		"gray",
		"red",
		"lime",
		"yellow",
		"blue",
		"fuchsia",
		"aqua",
		"white",
	}[colorNumber]
}

// StrToContents converts a single-line string into a one line of contents.
// Parse escape sequences, etc.
// 1 Content matches the characters displayed on the screen.
func StrToContents(str string, tabWidth int) lineContents {
	return parseString(str, tabWidth)
}

// ContentsToStr returns a converted string
// and byte position, as well as the content position conversion table.
func ContentsToStr(lc lineContents) (string, map[int]int) {
	var buff bytes.Buffer
	posCV := make(map[int]int)

	bn := 0
	for n, c := range lc {
		if c.mainc == 0 {
			continue
		}
		posCV[bn] = n
		_, err := buff.WriteRune(c.mainc)
		if err != nil {
			log.Println(err)
		}
		bn += len(string(c.mainc))
		for _, r := range c.combc {
			_, err := buff.WriteRune(r)
			if err != nil {
				log.Println(err)
			}
			bn += len(string(r))
		}
	}
	str := buff.String()
	posCV[bn] = len(lc)
	return str, posCV
}
