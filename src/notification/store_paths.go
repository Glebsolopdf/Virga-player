package notification

import (
	"path/filepath"
	"time"

	"virga-player/settings"
)

const (
	fileName                      = "notifications.db"
	legacyJSONFileName            = "notifications.json"
	notificationMaxAge            = 5 * 24 * time.Hour
	remoteCheckInterval           = 15 * time.Minute
	remoteRequestTimeout          = 10 * time.Second
	metaKeyLastNotificationsCheck = "lastnotificationscheck"
	metaKeyRemoteSupported        = "notifications_version_supported"
	unsupportedTitle              = "Unsupported notification format"
	unsupportedBodyText           = "Some fields from this notification cannot be displayed in this version."
)

func StatePath() string {
	return filepath.Join(settings.ConfigDirPath(), fileName)
}

func legacyStatePathFor(path string) string {
	return filepath.Join(filepath.Dir(path), legacyJSONFileName)
}
