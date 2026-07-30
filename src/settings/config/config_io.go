package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func LoadOrCreateConfig() (*Config, bool, error) {
	path := ConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := DefaultConfig()
			return cfg, true, SaveConfig(cfg)
		}
		return DefaultConfig(), false, err
	}

	var raw map[string]json.RawMessage
	migrate := false
	if json.Unmarshal(data, &raw) == nil {
		for _, k := range []string{"lyrics_mode", "lyrics_visible", "lyrics_auto_save_after_sec", "lyrics_double_confirm", "player_rain_layer", "lyrics_rain_layer"} {
			if _, ok := raw[k]; !ok {
				migrate = true
				break
			}
		}
	}

	cfg := DefaultConfig()
	if err := json.Unmarshal(data, cfg); err != nil {
		return DefaultConfig(), false, err
	}

	// migrate old rain_in_front_of_player -> player_rain_layer
	if raw != nil {
		if rawVal, ok := raw["rain_in_front_of_player"]; ok {
			var front *bool
			if json.Unmarshal(rawVal, &front) == nil && front != nil {
				if *front {
					cfg.PlayerRainLayer = RainLayerBehind
				} else {
					cfg.PlayerRainLayer = RainLayerFront
				}
			}
		}
	}

	cfg.normalize()
	if migrate {
		if err := SaveConfig(cfg); err != nil {
			return cfg, false, err
		}
	}
	return cfg, false, nil
}

func SaveConfig(cfg *Config) error {
	path := ConfigPath()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func ConfigPath() string {
	return filepath.Join(ConfigDirPath(), "config.json")
}

func ConfigDirPath() string {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return filepath.Join("/tmp", "virga-player")
		}
		configHome = filepath.Join(home, ".config")
	}
	return filepath.Join(configHome, "virga-player")
}
