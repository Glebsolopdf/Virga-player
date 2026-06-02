package app

import (
	"time"

	"virga-player/notification"
	"virga-player/version"
)

func (a *App) startRemoteNotificationSync() {
	if a.uninstallInProgress.Load() {
		return
	}
	if a.notifications == nil || a.remoteNotificationResults == nil {
		return
	}
	if !a.notifyRemoteInFlight.CompareAndSwap(false, true) {
		return
	}

	service := a.notifications
	go func() {
		defer a.notifyRemoteInFlight.Store(false)

		if a.uninstallInProgress.Load() {
			return
		}

		now := time.Now().UTC()
		shouldCheck, err := service.ShouldRunRemoteCheck(now)
		if err != nil {
			a.pushRemoteNotificationSyncResult(remoteNotificationSyncResult{err: err})
			return
		}
		if !shouldCheck {
			a.pushRemoteNotificationSyncResult(remoteNotificationSyncResult{skipped: true})
			return
		}

		ctx, cancel := notification.RemoteCheckTimeoutContext()
		feed, fetchErr := notification.FetchRemoteFeed(ctx, version.Notifications)
		cancel()

		completedAt := time.Now().UTC()
		if fetchErr != nil {
			a.pushRemoteNotificationSyncResult(remoteNotificationSyncResult{err: fetchErr, checkedAt: completedAt})
			return
		}

		a.pushRemoteNotificationSyncResult(remoteNotificationSyncResult{feed: &feed, checkedAt: completedAt})
	}()
}

func (a *App) pushRemoteNotificationSyncResult(result remoteNotificationSyncResult) {
	select {
	case a.remoteNotificationResults <- result:
	default:
	}
}
