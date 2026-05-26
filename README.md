# Awesome Autostart Manager (AAM)

A lightweight Windows GUI tool for managing which applications launch at startup. View, add, rename, and delete autostart entries — across both the **Windows Registry** and the **Startup folder** — from a single interface.

![Demo](https://media.giphy.com/media/v1.Y2lkPTc5MGI3NjExMmE1YjMxMzI5ZGZiNzUwZDQ0MGZmZDkyZGM5N2FkNjY4NTM5YjM5MCZlcD12MV9pbnRlcm5hbF9naWZzX2dpZklkJmN0PWc/U8iMFtS32920OmMshz/giphy.gif)

## Features

- **View all autostart apps** — reads both `HKCU\Software\Microsoft\Windows\CurrentVersion\Run` and the user's Startup folder, deduplicated into a single list.
- **Add apps** — pick any `.exe` via a file dialog; a shortcut is created in the Startup folder automatically.
- **Rename entries** — rename an autostart entry in both the registry and Startup folder at once.
- **Delete entries** — remove an app from autostart (attempts both registry and Startup folder regardless of where it was added).
- **Standalone binary** — all icons are embedded into the executable; no external files needed.

## Project Structure

```
├── main.go          # App init, window creation
├── ui.go            # GUI layout, table with cached app list
├── autostart.go     # Windows autostart CRUD (registry + Startup folder)
├── bundled.go       # Embedded PNG assets via //go:embed
├── saves/           # Source icon files (not needed at runtime)
├── go.mod
└── go.sum
```

## Getting Started

### Prerequisites

- **Windows 10/11**
- **[Go 1.23+](https://go.dev/dl/)** (for building from source)
- A C compiler for CGo (required by Fyne) — [MSYS2 MinGW-w64](https://www.msys2.org/) or [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) work well

### Build & Run

```powershell
git clone https://github.com/paxamans/awesomeProject
cd awesomeProject
go build -o aam.exe
./aam.exe
```

The resulting `aam.exe` is fully portable — copy it anywhere and it will work without the `saves/` folder.

### Quick Check (no run)

```powershell
go vet ./...     # static analysis
go build ./...   # compile check
```

## How It Works

1. On launch, the app reads autostart entries from **two sources**:
   - `.lnk` shortcuts in `%APPDATA%\Microsoft\Windows\Start Menu\Programs\Startup`
   - Named values in `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`
2. Results are deduplicated and cached in memory — the registry/filesystem is only hit on startup or after you make a change.
3. Adding an app creates a `.lnk` shortcut in the Startup folder via COM (`WScript.Shell`).
4. Deleting or renaming attempts both locations, so it works regardless of where the entry was originally created.

## Contributing

Contributions welcome! See the [Pull Requests](https://github.com/paxamans/awesomeProject/pulls) page.

## Bug Reports

Found a problem? [Open an issue](https://github.com/paxamans/awesomeProject/issues).

## Credits

Built by [paxamans](https://github.com/paxamans).
