package notification

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const killSwitchIDPrefix = "kill-switch:"

func (s *Service) ApplyKillSwitch(feed RemoteFeed, currentVersion string) (Item, bool, error) {
	id := killSwitchIDPrefix + NormalizeVersion(currentVersion)
	kind := normalizeKind(feed.KillSwitch.Kind)
	title := fallbackString(feed.KillSwitch.Title, "Notifications are disabled for this version")
	body := fallbackString(feed.KillSwitch.Body, "Update Virga Player to continue receiving notifications.")
	if strings.TrimSpace(feed.KillSwitch.GitHubURL) != "" {
		body = fmt.Sprintf("%s\n%s", body, strings.TrimSpace(feed.KillSwitch.GitHubURL))
	}

	item := Item{
		ID:        id,
		Kind:      kind,
		Title:     title,
		Body:      body,
		Version:   strings.TrimSpace(currentVersion),
		CreatedAt: s.now().UTC(),
	}
	added, err := s.upsert(item)
	if err != nil {
		return Item{}, false, err
	}
	return item, added, nil
}

func (s *Service) RemoveKillSwitchNotices() (bool, error) {
	filtered := make([]Item, 0, len(s.state.Items))
	changed := false
	for _, item := range s.state.Items {
		if strings.HasPrefix(item.ID, killSwitchIDPrefix) {
			changed = true
			continue
		}
		filtered = append(filtered, item)
	}
	if !changed {
		return false, nil
	}
	s.state.Items = filtered
	return true, saveState(s.path, s.state)
}

func (s *Service) ApplyRemoteNotifications(feed RemoteFeed, currentVersion string) (int, error) {
	added := 0
	for _, entry := range feed.Notifications {
		item := remoteNotificationToItem(entry, currentVersion, s.now().UTC())
		if item.ID == "" {
			continue
		}
		created, err := s.upsert(item)
		if err != nil {
			return added, err
		}
		if created {
			added++
		}
	}
	return added, nil
}

func remoteNotificationToItem(entry RemoteNotificationItem, currentVersion string, now time.Time) Item {
	id := remoteNotificationID(entry)
	if id == "" {
		return Item{}
	}

	createdAt := now
	if parsed, err := time.Parse(time.RFC3339Nano, strings.TrimSpace(entry.CreatedAt)); err == nil {
		createdAt = parsed.UTC()
	}

	version := strings.TrimSpace(entry.Version)
	if version == "" {
		version = strings.TrimSpace(currentVersion)
	}

	return Item{
		ID:        id,
		Kind:      normalizeKind(entry.Kind),
		Title:     fallbackString(entry.Title, unsupportedTitle),
		Body:      fallbackString(entry.Body, unsupportedBodyText),
		Version:   version,
		CreatedAt: createdAt,
	}
}

func remoteNotificationID(entry RemoteNotificationItem) string {
	raw := strings.TrimSpace(string(entry.ID))
	if raw == "" {
		return ""
	}

	var value string
	if err := json.Unmarshal(entry.ID, &value); err == nil {
		value = strings.TrimSpace(value)
		if value == "" {
			return ""
		}
		return "remote:" + value
	}

	var intValue int
	if err := json.Unmarshal(entry.ID, &intValue); err == nil {
		return "remote:" + strconv.Itoa(intValue)
	}

	trimmed := strings.Trim(raw, "\"")
	if trimmed == "" {
		return ""
	}
	return "remote:" + trimmed
}

func normalizeKind(raw string) Kind {
	switch Kind(strings.ToLower(strings.TrimSpace(raw))) {
	case KindWelcome:
		return KindWelcome
	case KindUpdate:
		return KindUpdate
	default:
		return KindInfo
	}
}
