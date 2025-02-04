# ov - feature rich terminal pager

[![PkgGoDev](https://pkg.go.dev/badge/github.com/noborus/ov)](https://pkg.go.dev/github.com/noborus/ov)
[![Actions Status](https://github.com/noborus/ov/workflows/Go/badge.svg)](https://github.com/noborus/ov/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/noborus/ov)](https://goreportcard.com/report/github.com/noborus/ov)

ov is a terminal pager.

* Can be used instead of `less` or `more`.  It can also be used instead of `tail -f`.
* `ov` also has an effective function for tabular text.

![ov1.png](https://raw.githubusercontent.com/noborus/ov/master/docs/ov1.png)

<!-- vscode-markdown-toc -->
* 1. [Feature](#Feature)
* 2. [Install](#Install)
	* 2.1. [deb package](#debpackage)
	* 2.2. [rpm package](#rpmpackage)
	* 2.3. [MacPorts (macOS)](#MacPortsmacOS)
	* 2.4. [Homebrew(macOS or Linux)](#HomebrewmacOSorLinux)
	* 2.5. [pkg (FreeBSD)](#pkgFreeBSD)
	* 2.6. [Binary](#Binary)
	* 2.7. [go install](#goinstall)
	* 2.8. [go get(details or developer version)](#gogetdetailsordeveloperversion)
* 3. [Usage](#Usage)
	* 3.1. [Basic usage](#Basicusage)
	* 3.2. [Config](#Config)
	* 3.3. [Section](#Section)
	* 3.4. [Follow mode](#Followmode)
	* 3.5. [Follow all mode](#Followallmode)
	* 3.6. [Follow section mode](#Followsectionmode)
	* 3.7. [Exec mode](#Execmode)
	* 3.8. [Search](#Search)
	* 3.9. [Mark](#Mark)
	* 3.10. [Watch](#Watch)
	* 3.11. [Mouse support](#Mousesupport)
* 4. [Command option](#Commandoption)
* 5. [Key bindings](#Keybindings)
* 6. [Customize](#Customize)
	* 6.1. [Style customization](#Stylecustomization)
	* 6.2. [Key binding customization](#Keybindingcustomization)

<!-- vscode-markdown-toc-config
	numbering=true
	autoSave=true
	/vscode-markdown-toc-config -->
<!-- /vscode-markdown-toc -->

##  1. <a name='Feature'></a>Feature

* Better support for Unicode and East Asian Width.
* Support for compressed files (gzip, bzip2, zstd, lz4, xz).
* Columns support column mode that can be selected by delimiter.
* The header row can always be displayed.
* Dynamic wrap/nowrap switchable.
* Supports alternating row style changes.
* Shortcut keys are customizable.
* The style of the effect is customizable.
* Supports follow-mode (like tail -f).
* Supports following multiple files and switching when updated.
* Supports the execution of commands that toggle both stdout and stder for display.
* Supports incremental search and regular expression search.

##  2. <a name='Install'></a>Install

###  2.1. <a name='debpackage'></a>deb package

You can download the package from [releases](https://github.com/noborus/ov/releases).

```console
curl -L -O https://github.com/noborus/ov/releases/download/vx.x.x/ov_x.x.x-1_amd64.deb
sudo dpkg -i ov_x.x.x-1_amd64.deb
```

###  2.2. <a name='rpmpackage'></a>rpm package

You can download the package from [releases](https://github.com/noborus/ov/releases).

```console
sudo rpm -ivh https://github.com/noborus/ov/releases/download/vx.x.x/ov_x.x.x-1_amd64.rpm
```

###  2.3. <a name='MacPortsmacOS'></a>MacPorts (macOS)

```console
sudo port install ov
```

###  2.4. <a name='HomebrewmacOSorLinux'></a>Homebrew(macOS or Linux)

```console
brew install noborus/tap/ov
```

###  2.5. <a name='pkgFreeBSD'></a>pkg (FreeBSD)

```console
pkg install ov
```

###  2.6. <a name='Binary'></a>Binary

You can download the binary from [releases](https://github.com/noborus/ov/releases).

```console
curl -L -O https://github.com/noborus/ov/releases/download/vx.x.x/ov_x.x.x_linux_amd64.zip
unzip ov_x.x.x_linux_amd64.zip
sudo install ov /usr/local/bin
```

###  2.7. <a name='goinstall'></a>go install

It will be installed in $GOPATH/bin by the following command.

```console
go install github.com/noborus/ov@latest
```

###  2.8. <a name='gogetdetailsordeveloperversion'></a>go get(details or developer version)

First of all, download only with the following command without installing it.

```console
go get -d github.com/noborus/ov
cd $GOPATH/src/github.com/noborus/ov
```

Next, to install to $GOPATH/bin, run the make install command.

```console
make install
```

Or, install it in a PATH location for other users to use
(For example, in /usr/local/bin).

```console
make
sudo install ov /usr/local/bin
```

##  3. <a name='Usage'></a>Usage

###  3.1. <a name='Basicusage'></a>Basic usage

ov supports open file name or standard input.

```console
ov filename
```

```console
cat filename|ov
```

Used by other commands by setting the environment variable **PAGER**.

```console
export PAGER=ov
```

See the [ov site](https://noborus.github.io/ov/) for more use cases.

###  3.2. <a name='Config'></a>Config

You can set style and key bindings in the setting file.

Create a `.ov.yaml` file in your user's home directory.

for example.

```filepath
$HOME/.ov.yaml
```

Windows.

```filepath
%USERPROFILE%/.ov.yaml
```

Please refer to the sample [ov.yaml](https://raw.githubusercontent.com/noborus/ov/master/ov.yaml) configuration file.

If you like `less` key bindings, copy  [ov-less.yaml](https://raw.githubusercontent.com/noborus/ov/master/ov-less.yaml) and use it.

###  3.3. <a name='Section'></a>Section

You specify `--section-delimiter`, you can move up and down in section units.
The start of the section can be adjusted with `--section-start`.

![section.png](docs/section.png)

The section-delimiter is written in a regular expression (for example: "^#").
(Line breaks are not included in matching lines).

For example, if you specify "^diff" for a diff that contains multiple files,
you can move the diff for each file.

###  3.4. <a name='Followmode'></a>Follow mode

Output appended data and move it to the bottom line (like `tail -f`).

```console
ov --follow-mode /var/log/syslog
```

```console
(while :; do echo random-$RANDOM; sleep 0.1; done;)|./ov  --follow-mode
```

###  3.5. <a name='Followallmode'></a>Follow all mode

Same as follow-mode, and switches to the last updated file when there are multiple files.

```console
ov --follow-all /var/log/nginx/access.log /var/log/nginx/error.log
```

###  3.6. <a name='Followsectionmode'></a>Follow section mode

Follow mode is line-by-line, while follow section mode is section-by-section.
Follow section mode displays the bottom section.
The following example is displayed from the header (#) at the bottom.

```console
ov --section-delimiter "^#" --follow-section README.md
```

 [Watch](#Watch) mode is a mode in which `--follow-section` and
 `--section-delimiter "^\f"` are automatically set.

###  3.7. <a name='Execmode'></a>Exec mode

Execute the command to display stdout/stderr separately.
Arguments after (`--`) are interpreted as command arguments.

Shows the stderr screen as soon as an error occurs, when used with `--follow-all`.

```console
ov --follow-all --exec -- make
```

###  3.8. <a name='Search'></a>Search

Search by forward search `/` key(default) or the backward search `?` key(defualt).
Search can be toggled between incremental search, regular expression search, and case sensitivity.
Displayed when the following are enabled in the search input prompt:

| Function | display | (Default)key |command option |
|:---------|:--------|:----|:--------------|
| Incremental search | (I) | alt+i | --incremental |
| Regular expression search | (R) | alt+r | --regexp-search  |
| Case sensitive | (Aa) | alt+c |  -i, --case-sensitive |

###  3.9. <a name='Mark'></a>Mark

Mark the display position with the `m` key(default).
The mark is decorated with `StyleMarkLine` and `MarkStyleWidth`.

Marks can be erased individually with the `M` key(default).
It is also possible to delete all marks with the `ctrl + delete` key(default).

Use the `>`next and `<`previous (default) key to move to the marked position.

###  3.10. <a name='Watch'></a>Watch

`ov` has a watch mode that reads the file every N seconds and adds it to the end.
When you reach EOF, add '\f' instead.
Go further to the last section.
The default is'section-delimiter', so the last loaded content is displayed.

for example.

```console
ov --watch 1 /proc/meminfo
```

###  3.11. <a name='Mousesupport'></a>Mouse support

The ov makes the mouse support its control.
This can be disabled with the option `--disable-mouse`.

If mouse support is enabled, tabs and line breaks will be interpreted correctly when copying.

Copying to the clipboard uses [atotto/clipboard](https://github.com/atotto/clipboard).
For this reason, the 'xclip' or 'xsel' command is required in Linux/Unix environments.

Selecting the range with the mouse and then left-clicking will copy it to the clipboard.

Pasting in ov is done with the middle button.
In other applications, it is pasted from the clipboard (often by pressing the right-click).

##  4. <a name='Commandoption'></a>Command option

```console
$ ov --help
ov is a feature rich pager(such as more/less).
It supports various compressed files(gzip, bzip2, zstd, lz4, and xz).

Usage:
  ov [flags]

Flags:
  -C, --alternate-rows             alternately change the line color
  -i, --case-sensitive             case-sensitive in search
  -d, --column-delimiter string    column delimiter (default ",")
  -c, --column-mode                column mode
      --completion string          generate completion script [bash|zsh|fish|powershell]
      --config string              config file (default is $HOME/.ov.yaml)
      --debug                      debug mode
      --disable-mouse              disable mouse support
  -e, --exec                       exec command
  -X, --exit-write                 output the current screen when exiting
  -a, --exit-write-after int       NUM after the current lines when exiting
  -b, --exit-write-before int      NUM before the current lines when exiting
  -A, --follow-all                 follow all
  -f, --follow-mode                follow mode
      --follow-section             follow section
  -H, --header int                 number of header rows to fix
  -h, --help                       help for ov
      --help-key                   display key bind information
      --incsearch                  incremental search (default true)
  -n, --line-number                line number mode
  -F, --quit-if-one-screen         quit if the output fits on one screen
      --regexp-search              regular expression search
      --section-delimiter string   section delimiter
      --section-start int          section start position
      --skip-lines int             skip the number of lines
  -x, --tab-width int              tab stop width (default 8)
  -v, --version                    display version information
  -T, --watch int                  watch mode interval
  -w, --wrap                       wrap mode (default true)
```

It can also be changed after startup.

##  5. <a name='Keybindings'></a>Key bindings

```console
 [Escape], [q]                * quit
 [ctrl+c]                     * cancel
 [Q]                          * output screen and quit
 [ctrl+q]                     * set output screen and quit
 [ctrl+z]                     * suspend
 [h], [ctrl+F1], [ctrl+alt+c] * display help screen
 [ctrl+F2], [ctrl+alt+e]      * display log screen
 [ctrl+l]                     * screen sync
 [ctrl+f]                     * follow mode toggle
 [ctrl+a]                     * follow all mode toggle
 [ctrl+alt+r]                 * enable/disable mouse

	Moving

 [Enter], [Down], [ctrl+N]    * forward by one line
 [Up], [ctrl+p]               * backward by one line
 [Home]                       * go to begin of line
 [End]                        * go to end of line
 [PageDown], [ctrl+v]         * forward by page
 [PageUp], [ctrl+b]           * backward by page
 [ctrl+d]                     * forward a half page
 [ctrl+u]                     * backward a half page
 [left]                       * scroll to left
 [right]                      * scroll to right
 [ctrl+left]                  * scroll left half screen
 [ctrl+right]                 * scroll right half screen
 [g]                          * number of go to line

	Move document

 []]                          * next document
 [[]                          * previous document
 [ctrl+k]                     * close current document

	Mark position

 [m]                          * mark current position
 [M]                          * remove mark current position
 [ctrl+delete]                * remove all mark
 [>]                          * move to next marked position
 [<]                          * move to previous marked position

	Search

 [/]                          * forward search mode
 [?]                          * backward search mode
 [n]                          * repeat forward search
 [N]                          * repeat backward search

	Change display

 [w], [W]                     * wrap/nowrap toggle
 [c]                          * column mode toggle
 [C]                          * color to alternate rows toggle
 [G]                          * line number toggle

	Change Display with Input

 [p], [P]                     * view mode selection
 [d]                          * delimiter string
 [H]                          * number of header lines
 [ctrl+s]                     * number of skip lines
 [t]                          * TAB width

	Section

 [alt+d]                      * section delimiter regular expression
 [ctrl+F3], [alt+s]           * section start position
 [space]                      * next section
 [^]                          * previous section
 [9]                          * last section
 [F2]                         * follow section mode toggle

	Close and reload

 [ctrl+F9], [ctrl+alt+s]      * close file
 [F5], [ctrl+alt+l]           * reload file
 [F4], [ctrl+alt+w]           * watch mode
 [ctrl+w]                     * set watch interval

	Key binding when typing

 [alt+c]                      * case-sensitive toggle
 [alt+r]                      * regular expression search toggle
 [alt+i]                      * incremental search toggle
```

##  6. <a name='Customize'></a>Customize

###  6.1. <a name='Stylecustomization'></a>Style customization

You can customize the following items.

* StyleAlternate
* StyleHeader
* StyleOverStrike
* StyleOverLine
* StyleLineNumber
* StyleSearchHighlight
* StyleColumnHighlight
* StyleMarkLine
* StyleSectionLine

Specifies the color name for the foreground and background [colors](https://pkg.go.dev/github.com/gdamore/tcell/v2#pkg-constants).
Specify bool values for Reverse, Bold, Blink, Dim, Italic, and Underline.

[Example]

```yaml
StyleAlternate:
  Background: "gray"
  Bold: true
  Underline: true
```

| item name | value | example |
|:----------|:------|:--------|
| Foreground | "color name" or "rgb" | "red" |
| Background | "color name" or "rgb" | "#2a2a2a" |
| Reverse | true/false | true |
| Bold | true/false | true |
| Blink | true/false | true |
| Dim | true/false | false |
| Italic | true/false | false |
| Underline | true/false | false |

###  6.2. <a name='Keybindingcustomization'></a>Key binding customization

You can customize key bindings.

[Example]

```yaml
    down:
        - "Enter"
        - "Down"
        - "ctrl+N"
    up:
        - "Up"
        - "ctrl+p"
```

See [ov.yaml](https://github.com/noborus/ov/blob/master/ov.yaml) for more information..
