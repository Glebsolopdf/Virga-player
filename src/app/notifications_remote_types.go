package app

import (
	"time"
	"virga-player/notification"
)

type remoteNotificationSyncResult struct {
	feed      *notification.RemoteFeed
	checkedAt time.Time
	err       error
	skipped   bool
}
