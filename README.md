# reamd

**A read-only Markdown reader for your terminal.** `reamd` opens one document in the current terminal, lets you scroll through it, and returns you to the shell when you quit.

> [!NOTE]
> `reamd` is in early development. Core CommonMark and selected GFM rendering, file loading, scrolling TUI, plain pipe output, link navigation, text selection/copy, and theme behavior are implemented. Final visual polish is tracked separately. Manual smoke testing has been recorded for Windows Terminal and GNOME Terminal on Ubuntu WSL 2; Terminal.app has not been tested.

## Usage

```sh
# Open one Markdown file in the current terminal
reamd README.md

# Choose a theme
reamd --theme transparent README.md

# Show usage or version information
reamd --help
reamd --version
```

With no arguments, `reamd` shows help and usage examples and exits with code 0. It accepts one file path per run and does not read from stdin. In an interactive terminal it opens a read-only TUI; when output is redirected or piped, it prints plain text without ANSI styling.

Current keys: `Up`/`Down` and `j`/`k` scroll one line, `PageUp`/`PageDown` scroll one screen, `Home`/`End` jump to the start or end, `Tab`/`Shift+Tab` focus links, `Enter` activates a focused link, `Ctrl+C` copies selected text, `?` opens help, and `q` or `Esc` quits. On mouse-capable terminals, drag selects visible text, the wheel scrolls, and the scrollbar can be clicked or dragged; clicking a link activates it. HTTP(S) links open in the default browser only after activation; anchors jump to headings. Terminals smaller than 40 columns by 10 rows show a short resize message.

## Current status

The reader displays core CommonMark content plus GFM tables, task lists, strikethrough, and visible URL/email autolinks in a scrollable TUI or as plain rendered text when output is redirected. Links can be focused and activated with the keyboard or mouse. Wide tables switch to label/value rows when they do not fit.

```sh
# Run from source
go run ./cmd/reamd README.md

# Build a local binary
go build -o reamd ./cmd/reamd
```

## Planned MVP

- Read a Markdown file and show it in a scrollable, read-only TUI.
- Support CommonMark and GitHub Flavored Markdown (GFM), including headings, lists, quotes, tables, links, and code blocks. Bare URLs and email addresses are recognized as autolinks.
- Adapt layout to terminal width. Switch wide tables to a stacked layout when needed.
- Provide in-app zoom by adjusting layout and wrapping. This will not change the terminal emulator's actual font size.
- Support `auto`, `dark`, `light`, and `transparent` themes.
- Support keyboard and mouse navigation, text selection, and copying rendered text to the system clipboard.
- Show image alt text and URLs without downloading images. Open HTTP(S) links in the default browser only after explicit user activation; never open links automatically.
- Show Mermaid source as a fallback; diagram rendering is a later possibility.
- Offer plain text output when writing to a pipe, without ANSI styling by default.

## Design principles

- Keep Markdown parsing, terminal layout, and the TUI separate.
- Do not execute embedded HTML or pass terminal control sequences from document content through to the terminal.
- Do not make network requests while rendering a document. Explicit activation may open an HTTP(S) link in the default browser.

## Project status

Go 1.25 is the approved minimum project toolchain. Goldmark is the Markdown parser and tcell is the TUI library. The application uses `uniseg` for display width, `x/term` for terminal detection, and `github.com/atotto/clipboard` behind an internal adapter for clipboard writes. Installation method and release process have not been selected. The distribution goal is to provide ready-to-run binaries so users do not need to install Go themselves.

See the [technical analysis and MVP proposal](docs/markdown-terminal-reader-analysis.md) for the REA observations, design decisions, open questions, and acceptance criteria.

See the [command specification](docs/command-spec.md) for the planned arguments, themes, and keyboard controls.

See the [acceptance criteria and Markdown fixtures](docs/acceptance-criteria.md) for how MVP behavior will be verified.

See [DESIGN.md](docs/DESIGN.md) for the Ocean visual direction and TUI design constraints.

See the [development workflow](docs/development-workflow.md) for the recommended implementation phases and verification gates.

See the [development ticket checklists](docs/tickets/README.md) for agent-sized tasks, prerequisites, and acceptance checks.
