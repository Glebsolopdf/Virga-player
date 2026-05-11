![Virga logo](<for readme/virga.png>)

![Go](https://img.shields.io/badge/Go-1.25+-blue.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)
![Platform](https://img.shields.io/badge/Platform-Linux-yellow.svg)

Virga Player is a terminal application written in Go for visualizing music playback with rain particle effects and track metadata.

![Player preview 1](<for readme/prew.png>)
![Player preview 2](<for readme/prew1.png>)

*Russian documentation: [README.ru.md](README.ru.md)*

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                      Main application                           │
│                          (app/app.go)                          │
└───────────────────────────┬─────────────────────────────────────┘
                            │
        ┌───────────────────┼───────────────────┐
        │                   │                   │
        ▼                   ▼                   ▼
┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│ Terminal     │    │ Event        │    │ Settings     │
│ rendering    │    │ handler      │    │ manager      │
│ (renderer/)  │    │(app/interact)│    │ (settings/)  │
└──────────────┘    └──────────────┘    └──────────────┘
        │                   │                   │
        └───────────────────┼───────────────────┘
                            │
        ┌───────────────────┼───────────────────────────┐
        │                   │                           │
        ▼                   ▼                           ▼
┌──────────────┐    ┌────────────────┐    ┌─────────────────────┐
│ Animation    │    │ Particle       │    │ Audio analysis      │
│ engine       │    │ system         │    │                     │
│ (animation/) │    │ (rain/)        │    │ • frequency bands   │
│              │    │                │    │ • envelope tracking │
│• FPS control │    │• physics       │    │ • audio capture     │
│• timing      │    │• music-reactive│    │ (audio/analyzer.go) │
│• main loop   │    │• rendering     │    │                     │
└──────────────┘    └────────────────┘    └─────────────────────┘
                            │
        ┌───────────────────┼───────────────────┐
        │                   │                   │
        ▼                   ▼                   ▼
┌──────────────┐    ┌────────────────┐    ┌──────────────┐
│ Music data   │    │ Artwork        │    │ Scene        │
│              │    │ display        │    │ rendering    │
│ • Playerctl  │    │ • Sixel (PNG)  │    │ • background │
│ • JSON file  │    │ • text mode    │    │ • buildings  │
│ • fallback   │    │ • animations   │    │ • UI elements│
│ (music/)     │    │ (artwork/)     │    │ (scene/)     │
└──────────────┘    └────────────────┘    └──────────────┘
```

### Main components

#### 1. App (`app/`)
Central coordinator of application lifecycle and subsystem integration.

- `app.go` - main App structure with subsystems
- `init.go` - initialization and setup
- `lifecycle.go` - start, stop, cleanup
- `interaction.go` - keyboard event handling
- `settings_flow.go` - settings menu navigation logic
- `tick.go` - main loop timing
- `install/` - installation and environment helpers
- `bootstrap/` - bootstrap procedures
- `state/` - application state
- `events/` - internal event handlers

#### 2. Animation engine (`animation/`)
Handles frame timing and FPS control.

- `engine.go` - FPS limiter and timing calculation
- Provides a stable animation loop, default 60 FPS

#### 3. Particle system (`rain/`)
Physics-based particle simulation for the rain effect.

- `types.go` - `Particle` and `ParticleSystem` structures
- `system.go` - particle lifecycle management
- `spawn.go` - spawn logic for static and music-reactive particles
- `update.go` - physics update logic (velocity, position, acceleration)
- `draw.go` - particle rendering into the screen buffer
- `particle.go` - individual particle behavior

**Physics model:**
- Particles spawn at the top with direction-dependent velocity
- Gravity increases `VelY` each frame
- Horizontal velocity `VelX` simulates wind
- Particle length grows over time
- Collision detection removes particles at the bottom

#### 4. Audio analyzer (`audio/`)
Captures audio in real time and analyzes frequency content.

- `analyzer.go` - audio monitor integration and FFT analysis
- Captures PCM audio at 11025 Hz using `parec`
- Splits audio into three frequency bands:
  - low (60-180 Hz)
  - mid (500-2000 Hz)
  - high (2800-5000 Hz)
- Computes RMS envelope for overall dynamics
- Returns normalized 0-1 values for use in effects

#### 5. Music data (`music/`)
Fetches track metadata from multiple sources.

- `track.go` - main `TrackInfo` structure and caching
- `playerctl.go` - MPRIS player integration
- `json_default.go` - JSON fallback support
- `artwork_path.go` - artwork path resolution
- `artwork_lookup.go` - artwork discovery
- `format.go` - track metadata formatting
- `mpd.go` - MPD support placeholder

**Source priority:**
1. Playerctl (MPRIS)
2. JSON file `/tmp/virga-player/track.json`
3. fallback empty data

#### 6. Artwork display (`app/artwork/`)
Renders album artwork in terminal or text fallback.

- `artwork.go` - artwork state and lifecycle
- `draw.go` - rendering coordination
- `image_io.go` - image loading and file handling
- `image_render.go` - image conversion for terminal display
- `render_sixel.go` - Sixel rendering
- `render_text.go` - text/Unicode fallback rendering
- `sixel_support.go` - terminal capability detection

**Rendering strategy:**
- Detects Sixel support
- Converts PNG to 256×256 Sixel if supported
- Falls back to colored text layout otherwise
- Applies audio-driven effects such as fade and pulse

#### 7. Terminal renderer (`renderer/`)
Low-level tcell screen abstraction and buffer management.

- `renderer.go` - frame buffer and drawing pipeline

#### 8. Settings (`settings/`)
Configuration and settings UI management.

- `config.go` - config load/save (JSON)
- `theme.go` - theme loading and current theme access
- `page.go` - settings page abstraction
- `page/handler.go` - menu item handlers
- `page/menu.go` - menu rendering and navigation
- `page/page.go` - page layout
- `page/render.go` - UI rendering logic
- `theme/defaults.go` - default CSS theme values
- `theme/loader.go` - load/create theme file
- `theme/parser.go` - parse CSS theme variables
- `theme/theme.go` - theme structure and current theme

**Config path:** `~/.config/virga-player/config.json`

**Available settings:**
- FPS (default 60)
- max_particles (default 220)
- rain_speed (default 100)
- rain_enabled
- music_reactive
- music_reactive_intensity (default 100)
- direction (`right-to-left`, `left-to-right`, `straight`, `random`)
- cover_animation
- player

#### 9. Scene (`scene/`)
Background rendering and scene composition.

- `scene.go` - scene orchestration
- renders static city background
- overlays rain particle effect
- positions artwork display
- handles multi-layer composition

## Getting started

### Requirements

- Go 1.25 or later
- Linux
- local audio backend (for audio analysis)
- ImageMagick `convert` (for Sixel artwork)
- Terminal with 24-bit color

### Installation

```bash
git clone https://github.com/Glebsolopdf/Virga-playerl
cd Virga-player
cd src
go mod download
go build -o ../virga-player main.go
./Virga-player
```

## Configuration

Config is stored in `~/.config/virga-player/config.json`:

```json
{
  "fps": 60,
  "max_particles": 220,
  "rain_speed": 100,
  "rain_enabled": true,
  "music_reactive": false,
  "music_reactive_intensity": 100,
  "cover_animation": false,
  "direction": "random",
  "player": false
}
```

### Config parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `fps` | int | 60 | frames per second (15-120) |
| `max_particles` | int | 220 | maximum particles (20-500) |
| `rain_speed` | int | 100 | base rain speed (25-300) |
| `rain_enabled` | bool | true | enable rain animation |
| `music_reactive` | bool | false | enable music-reactive behavior |
| `music_reactive_intensity` | int | 100 | intensity scaling (20-200) |
| `direction` | string | `random` | rain direction |
| `cover_animation` | bool | false | enable artwork animation |
| `player` | bool | false | show player information widget |

### Theme

Theme is stored in CSS at `~/.config/virga-player/style.css`.
The file is created automatically on first run.

Example theme variables:

```css
:root {
  --bg: transparent;
  --message-text: white;
  --track-title: white;
  --track-artist: green;
  --track-album: yellow;
  --track-time: gray;
  --rain-head: white;
  --rain-tail: gray;
}
```

## Usage

### Keyboard controls

| Key | Action |
|-----|--------|
| S | open settings menu |
| Enter | toggle/select option |
| ESC | exit application |

## Cover rendering

- Sixel: used when supported by terminal
- Text mode: fallback when Sixel is unavailable
- Requires ImageMagick `convert`

## Music reactivity

- Audio analysis via local audio backend
- Three frequency bands: 60-180 Hz, 500-2000 Hz, 2800-5000 Hz
- RMS envelope calculation
- `musicReactive` scales particle parameters based on analysis

## Performance

Configuration limits:
- `fps`: 15-120
- `max_particles`: 20-500
- `rain_speed`: 25-300
- `music_reactive_intensity`: 20-200

## Troubleshooting

### No audio reaction

- Verify your audio subsystem with `pactl info`
- Ensure audio source is available
- Restart the application

### Cover not rendering

- Verify `artwork_url`
- Check `convert` is installed
- Check file permissions

### Terminal rendering issues

- Verify `COLORTERM`
- Try a different terminal emulator

### Audio backend issues

- Check audio subsystem status
- `groups $USER`
- `pactl list short sinks`

## Project structure

```
├── build.sh
├── for readme
│   ├── prew1.png
│   ├── prew.png
│   └── virga.png
├── install.sh
├── LICENSE
├── README.md
├── README.ru.md
├── src
│   ├── animation
│   │   └── engine.go
│   ├── app
│   │   ├── app.go
│   │   ├── artwork
│   │   │   ├── artwork.go
│   │   │   ├── draw.go
│   │   │   ├── image_io.go
│   │   │   ├── image_render.go
│   │   │   ├── render_sixel.go
│   │   │   ├── render_text.go
│   │   │   └── sixel_support.go
│   │   ├── bootstrap
│   │   │   └── bootstrap.go
│   │   ├── events
│   │   │   └── events.go
│   │   ├── frame
│   │   │   ├── frame.go
│   │   │   ├── hit.go
│   │   │   └── render.go
│   │   ├── init.go
│   │   ├── install
│   │   │   ├── install.go
│   │   │   ├── shell.go
│   │   │   ├── system.go
│   │   │   ├── user.go
│   │   │   └── utils.go
│   │   ├── interaction.go
│   │   ├── lifecycle.go
│   │   ├── message
│   │   │   └── message.go
│   │   ├── player
│   │   │   └── player.go
│   │   ├── settings_flow.go
│   │   ├── state
│   │   │   └── state.go
│   │   └── tick.go
│   ├── audio
│   │   ├── analysis.go
│   │   ├── analyzer.go
│   │   ├── dsp.go
│   │   ├── monitor_source.go
│   │   └── types.go
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   ├── music
│   │   ├── artwork_lookup.go
│   │   ├── artwork_path.go
│   │   ├── format.go
│   │   ├── json_default.go
│   │   ├── mpd.go
│   │   ├── playerctl.go
│   │   └── track.go
│   ├── rain
│   │   ├── draw.go
│   │   ├── particle.go
│   │   ├── spawn.go
│   │   ├── system.go
│   │   ├── types.go
│   │   └── update.go
│   ├── renderer
│   │   └── renderer.go
│   ├── scene
│   │   ├── draw.go
│   │   ├── generate.go
│   │   ├── scene.go
│   │   └── types.go
│   └── settings
│       ├── config.go
│       ├── page
│       │   ├── handler.go
│       │   ├── menu.go
│       │   ├── page.go
│       │   └── render.go
│       ├── page.go
│       ├── theme
│       │   ├── defaults.go
│       │   ├── loader.go
│       │   ├── parser
│       │   │   ├── color.go
│       │   │   ├── component.go
│       │   │   ├── rgb.go
│       │   │   └── rune.go
│       │   ├── parser.go
│       │   └── theme.go
│       └── theme.go
└── virga-player

22 directories, 77 files
```

## Development

### Dependencies

- `github.com/gdamore/tcell/v2`
- `golang.org/x/image`

Dependencies are Go-compatible and do not require C/C++.

## Environment variables

- `COLORTERM` - set to `truecolor` for 24-bit colors
- `TERM` - terminal type, e.g. `xterm-256color`, `xterm-kitty`
- `HOME` - used for `$HOME/.config/virga-player/`

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

## Links

- [Sixel Graphics Format](https://en.wikipedia.org/wiki/Sixel)
- [ANSI Escape Codes](https://en.wikipedia.org/wiki/ANSI_escape_code)
- [Audio backend documentation](https://www.freedesktop.org/wiki/Software/)
- [MPRIS Specification](https://specifications.freedesktop.org/mpris-spec/)
- [tcell Documentation](https://github.com/gdamore/tcell)
- [Go Image Package](https://pkg.go.dev/image)
