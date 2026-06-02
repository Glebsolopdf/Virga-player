package notification

import "encoding/json"

type RemoteFeed struct {
	AllowedVersionsAlias []string                 `json:"allowed-version"`
	AllowedVersionsList  []string                 `json:"allowed_versions"`
	KillSwitch           RemoteKillSwitch         `json:"kill_switch_notification"`
	Notifications        []RemoteNotificationItem `json:"notifications"`
}

func (f RemoteFeed) AllowedVersions() []string {
	if len(f.AllowedVersionsList) > 0 {
		return f.AllowedVersionsList
	}
	return f.AllowedVersionsAlias
}

type RemoteKillSwitch struct {
	Kind      string `json:"kind"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	GitHubURL string `json:"github_url"`
}

type RemoteNotificationItem struct {
	ID         json.RawMessage `json:"id"`
	Kind       string          `json:"kind"`
	Title      string          `json:"title"`
	Body       string          `json:"body"`
	Version    string          `json:"version"`
	CreatedAt  string          `json:"created_at"`
	MinVersion string          `json:"min_version"`
	MaxVersion string          `json:"max_version"`
}
