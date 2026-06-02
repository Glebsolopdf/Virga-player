package app

import (
	"time"

	"virga-player/notification"
	"virga-player/settings"
	"virga-player/version"
)

func (a *App) processRemoteNotificationSyncResults() {
	for {
		select {
		case result := <-a.remoteNotificationResults:
			if a.uninstallInProgress.Load() {
				continue
			}
			a.applyRemoteNotificationSyncResult(result)
		default:
			return
		}
	}
}

func (a *App) applyRemoteNotificationSyncResult(result remoteNotificationSyncResult) {
	if a.notifications != nil && !result.checkedAt.IsZero() {
		if err := a.notifications.SetLastRemoteCheck(result.checkedAt); err != nil && a.debug != nil {
			a.debug.Warnf("notifications last-check save failed: %v", err)
		}
	}

	if result.err != nil {
		if a.debug != nil {
			a.debug.Warnf("remote notification sync failed: %v", result.err)
		}
		return
	}
	if result.skipped || result.feed == nil || a.notifications == nil {
		return
	}

	supported := notificationVersionSupported(*result.feed)
	a.notificationsSupported = supported
	if err := a.notifications.SetRemoteSupportState(supported); err != nil && a.debug != nil {
		a.debug.Warnf("notifications support-state save failed: %v", err)
	}

	if supported {
		if _, err := a.notifications.RemoveKillSwitchNotices(); err != nil && a.debug != nil {
			a.debug.Warnf("kill-switch cleanup failed: %v", err)
		}
		if added, err := a.notifications.ApplyRemoteNotifications(*result.feed, version.AppVersion); err != nil {
			if a.debug != nil {
				a.debug.Warnf("remote notifications apply failed: %v", err)
			}
		} else if added > 0 && a.debug != nil {
			a.debug.Infof("remote notifications added: %d", added)
		}
	} else {
		if a.cfg != nil {
			a.cfg.NotificationsEnabled = false
			_ = settings.SaveConfig(a.cfg)
		}
		if a.settingsPage != nil && a.settingsPage.Config != nil {
			a.settingsPage.Config.NotificationsEnabled = false
		}
		if _, added, err := a.notifications.ApplyKillSwitch(*result.feed, version.AppVersion); err != nil {
			if a.debug != nil {
				a.debug.Warnf("kill-switch apply failed: %v", err)
			}
		} else if added && a.debug != nil {
			a.debug.Infof("kill-switch notification added")
		}
	}

	if a.debug != nil && !result.checkedAt.IsZero() {
		a.debug.Infof("remote notification sync finished at %s", result.checkedAt.Format(time.RFC3339Nano))
	}
	a.refreshNotificationsPageBindings()
	a.armNotificationToast()
}

func notificationVersionSupported(feed notification.RemoteFeed) bool {
	allowed := feed.AllowedVersions()
	if len(allowed) == 0 {
		return true
	}
	return notification.IsVersionAllowed(version.AppVersion, allowed)
}

func (a *App) refreshNotificationsPageBindings() {
	if a.settingsPage == nil {
		return
	}
	a.settingsPage.SetNotificationsSupported(a.notificationsSupported)
	a.settingsPage.SetNotifications(a.notificationItems(), a.openNotificationsPage)
}
