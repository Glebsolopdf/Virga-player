![Virga prew](<for readme/prew.png>) 

![Go](https://img.shields.io/badge/Go-1.25+-blue.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)
![Platform](https://img.shields.io/badge/Platform-Linux-yellow.svg)

Virga Player is a terminal application written in Go for visualizing music playback with rain particle effects and track metadata.

## Player Compatibility

| Player | Status | Notes |
|:---|:---:|:---|
| ![Spotify](https://img.shields.io/badge/Spotify_(Client)-ffffff?style=flat&logo=spotify&logoColor=000000&labelColor=ffffff) | ![Full](https://img.shields.io/badge/Full-2ea44f?style=flat&logo=check&logoColor=ffffff) | Perfect metadata sync, cover art works flawlessly. |
| ![SoundCloud](https://img.shields.io/badge/SoundCloud_(Web)-ffffff?style=flat&logo=soundcloud&logoColor=000000&labelColor=ffffff) | ![Mostly Stable](https://img.shields.io/badge/Mostly_Stable-94d3a2?style=flat&logo=check&logoColor=ffffff) | Works great, tracks metadata and artwork correctly. |
| ![Yandex Music](https://img.shields.io/badge/Yandex_Music_(Web)-ffffff?style=flat&logo=yandex&logoColor=000000&labelColor=ffffff) | ![Mostly Stable](https://img.shields.io/badge/Mostly_Stable-94d3a2?style=flat&logo=check&logoColor=ffffff) | Highly functional, but web-extension might briefly drop cover art on track change. |
| ![Spotify](https://img.shields.io/badge/Spotify_(Web)-ffffff?style=flat&logo=spotify&logoColor=000000&labelColor=ffffff) | ![Bugs](https://img.shields.io/badge/Minor_Bugs-dfb317?style=flat&logo=alert&logoColor=ffffff) | Syncs metadata and controls fine, but fails to extract cover art. |
| ![YouTube Music](https://img.shields.io/badge/YouTube_Music-ffffff?style=flat&logo=youtubemusic&logoColor=000000&labelColor=ffffff) | ![Bugs](https://img.shields.io/badge/Minor_Bugs-dfb317?style=flat&logo=alert&logoColor=ffffff) | Metadata and artwork sync properly, but fails to track duration. |
| ![Apple Music](https://img.shields.io/badge/Apple_Music_(Web)-ffffff?style=flat&logo=apple-music&logoColor=000000&labelColor=ffffff) | ![Bugs](https://img.shields.io/badge/Minor_Bugs-dfb317?style=flat&logo=alert&logoColor=ffffff) | Web version via extension fails to track duration; artwork updates only on the first track of an album. |

### Status Legend
* ![Full](https://img.shields.io/badge/Full-2ea44f?style=flat&logo=check&logoColor=ffffff) — Works flawlessly out of the box.
* ![Mostly Stable](https://img.shields.io/badge/Mostly_Stable-94d3a2?style=flat&logo=check&logoColor=ffffff) — Highly functional, but rare minor glitches may occur due to browser MPRIS extensions.
* ![Bugs](https://img.shields.io/badge/Minor_Bugs-dfb317?style=flat&logo=alert&logoColor=ffffff) — Functional, but missing some features (like duration or artwork) due to web-player limitations.
* ![Not Supported](https://img.shields.io/badge/Not_Supported-954a4a?style=flat&logo=octicons&logoColor=ffffff) — Does not work with Virga Player at the moment.

### Requirements

- Go 1.25 or later
- local audio backend (PulseAudio/PipeWire) with `pactl` available
- ImageMagick `convert` (for Sixel artwork)
- Terminal with 24-bit color

### Installation

>**The installer for Arch, Debian, Fedora, and Void will automatically install dependencies and build the project.**

```bash
curl -sSL https://raw.githubusercontent.com/Glebsolopdf/Virga-player/main/install.sh | bash
```

## Arch 
```bash
sudo pacman -S go git imagemagick
```

## Void 
```bash
sudo xbps-install -S go git ImageMagick
```

## Debian
```bash
sudo apt install golang-go git imagemagick
```

## Fedora
```bash
sudo dnf install golang git ImageMagick
```

>**After you have installed all the dependencies, you can use these commands to quickly compile Virga and get started.**

```bash
git clone https://github.com/Glebsolopdf/Virga-player
cd Virga-player
cd src
go mod download
go build -o ../virga-player main.go
./virga-player
```

## Usage

> **Note:** The binary file will add itself to the `PATH`, so you can run Virga by typing `virga` or `virgaplayer` in the terminal.
```bash
virga
```
Use the `--debug` flag to enable debug logging and the in-app debug overlay:

```bash
virga --debug
```

When enabled, Virga shows a debug overlay with runtime diagnostics and log messages.

## Configuration

Config is stored in `~/.config/virga-player/config.json`:

```json
{
  "fps": 60,
  "max_particles": 220,
  "rain_speed": 100,
  "pulse_speed": 100,
  "pulse_mode": "rain",
  "rain_enabled": true,
  "music_reactive": false,
  "music_reactive_intensity": 100,
  "rain_in_front_of_player": true,
  "direction": "random",
  "player": false
}
```

`rain_in_front_of_player` controls the layer order when player mode is enabled: `true` draws rain over the player, `false` keeps rain behind the player.

`pulse_speed` controls the base rise/fade speed for pulse effects on both rain and cover artwork. Virga also adapts that speed to recent beat/transient timing in the current track, so fast songs tend to get shorter, quicker pulses while slower songs keep a longer pulse tail.

`pulse_mode` controls where pulse is allowed to appear: `off`, `rain`, `cover`, or `all`.

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
  --rain-layer-very-near: white;
  --rain-layer-near: lightcyan;
  --rain-layer-mid: white;
  --rain-layer-far: lightgray;
  --rain-layer-very-far: darkgray;

  --settings-title: white;
  --settings-hint: gray;
  --settings-text: white;
  --settings-selected-fg: black;
  --settings-selected-bg: white;

  --timeline-char-left: '[';
  --timeline-char-right: ']';
  --timeline-char-played: '█';
  --timeline-char-current: '▌';
  --timeline-char-empty: '░';

  --rain-char-body: '│';
  --rain-char-head: '•';
  --rain-char-left: '/';
  --rain-char-right: '\\';
  --artwork-char-block: '▀';
}
```

### Keyboard controls

| Key | Action |
|-----|--------|
| S | open settings menu |
| Enter | toggle/select option |
| ESC | exit application |

## Troubleshooting

### No audio reaction

- Verify your audio subsystem with `pactl info`
- Ensure `pactl` is installed and available in PATH
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

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

![Virga logo](<for readme/virga.png>)
