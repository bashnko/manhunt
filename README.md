# manhunt

manhunt is a lightweight desktop search launcher for Linux.

It lets you:

- search the web from one prompt
- use short engine keywords like `gg`, `yt`, `rd`, `so`
- open saved bookmarks quickly
- add bookmarks from the launcher UI

When you run manhunt for the first time, it creates your config file automatically at:

`~/.config/manhunt/config.json`

## Installation

### 1. Go install

Use this if you have Go installed.

```bash
go install github.com/bashnko/manhunt/cmd/manhunt@latest
```

Make sure your Go bin directory is in `PATH` (usually `~/go/bin`).

### 2. npm

If you prefer npm, install the package named `manhunt`:

```bash
npm install -g manhunt
```

### 3. Manual build

Clone and build from source:

```bash
git clone https://github.com/bashnko/manhunt.git
cd manhunt
go build -o manhunt ./cmd/manhunt
sudo install -m 0755 manhunt /usr/local/bin/manhunt
```

Or run directly without installing:

```bash
go run ./cmd/manhunt
```

## Setup

manhunt currently requires one external dependency:

- [rofi](https://github.com/davatorium/rofi)

Install it with your distro package manager.

Examples:

```bash
# Debian/Ubuntu
sudo apt install rofi

# Arch
sudo pacman -S rofi

# Fedora
sudo dnf install rofi
```

Once both are installed (`rofi` + `manhunt`), proceed with these steps:

1. Run `manhunt` from your terminal once type `manhunt` (or `manhunt init`) to create `~/.config/manhunt/config.json`.
2. Bind `manhunt` to a hotkey in your compositor/window manager.
3. Press the hotkey and start searching.

## Keymap examples

Add one of these bindings to your desktop environment config.

> [!NOTE]
> Keymap configuration can differ based on your system setup, compositor version, and how you manage your keybindings.
> Treat the examples below as templates and adjust paths/syntax to match your own configuration.

<details>
<summary>Hyprland</summary>

File: `~/.config/hypr/hyprland.conf`

```ini
bind = SUPER, SPACE, exec, manhunt
```

</details>

<details>
<summary>i3wm</summary>

File: `~/.config/i3/config`

```i3
bindsym $mod+space exec --no-startup-id manhunt
```

</details>

<details>
<summary>niri</summary>

File: `~/.config/niri/config.kdl`

```kdl
binds {
	Mod+Space { spawn "manhunt"; }
}
```

</details>

<details>
<summary>Sway</summary>

File: `~/.config/sway/config`

```sway
bindsym $mod+space exec manhunt
```

</details>

<details>
<summary>bspwm (sxhkd)</summary>

File: `~/.config/sxhkd/sxhkdrc`

```text
super + space
	manhunt
```

</details>

After changing keybind configs, reload your compositor/window manager config.

## Usage

- type normal text to search with the default engine
- type `engine query` to use a specific engine, for example: `yt lofi mix`
- type `:links` to browse saved bookmarks
- type `:add_url` to add a bookmark interactively

Default engines:

- `gg` Google
- `yt` YouTube
- `rd` Reddit
- `so` Stack Overflow

## Example config

manhunt is fully extensible and configurable. You can customize engines, commands, prefixes, and bookmark shortcuts to match your workflow.

```json
{
  "DefaultEngine": "gg",
  "CommandPrefix": ":",
  "LinksCommand": ":links",
  "AddURLCommand": ":add_url",
  "SearchEngines": {
    "gg": "https://www.google.com/search?q=%s",
    "ms": "https://music.youtube.com/search?q=%s",
    "rd": "https://www.reddit.com/search/?q=%s",
    "so": "https://stackoverflow.com/search?q=%s",
    "yt": "https://www.youtube.com/results?search_query=%s"
  },
  "Bookmarks": [
    {
      "keyword": "vercel",
      "name": "vercel",
      "url": "https://vercel.app/"
    },
    {
      "keyword": "gh",
      "name": "GitHub",
      "url": "https://github.com"
    }
  ]
}
```
