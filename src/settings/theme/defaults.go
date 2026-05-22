package theme

const defaultThemeCSS = `/* Virga Player theme */
:root {
  --bg: transparent;
  --message-text: white;

  --track-title: white;
  --track-artist: green;
  --track-album: yellow;
  --track-time: gray;

  --timeline-bracket: silver;
  --timeline-played: green;
  --timeline-current: green;
  --timeline-remaining: gray;

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
  --settings-danger: red;
  --settings-danger-bg: maroon;

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
`
