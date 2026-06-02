package app

import (
	"time"

	"virga-player/notification"
)

const notificationToastDuration = 5 * time.Second

func (a *App) notificationsVisible() bool {
	if !a.notificationsSupported {
		return true
	}
	return a.cfg == nil || a.cfg.NotificationsEnabled
}

func (a *App) armNotificationToast() {
	if a.cfg == nil || !a.cfg.NotifyUnreadToast {
		if a.notificationToast != nil {
			a.notificationToast.Hide()
		}
		return
	}
	if !a.notificationsVisible() {
		if a.notificationToast != nil {
			a.notificationToast.Hide()
		}
		return
	}
	if a.notificationToast == nil || a.notifications == nil || !a.notifications.HasUnread() {
		return
	}
	a.notificationToast.Show(a.notifications.ToastMessage(), notificationToastDuration)
}

func (a *App) notificationItems() []notification.Item {
	if a.notifications == nil || !a.notificationsVisible() {
		return nil
	}
	return a.notifications.Items()
}

func (a *App) openNotificationsPage() []notification.Item {
	if a.notifications == nil || !a.notificationsVisible() {
		return nil
	}
	changed, err := a.notifications.MarkAllRead()
	if err != nil {
		if a.debug != nil {
			a.debug.Warnf("notification mark-read failed: %v", err)
		}
		return a.notifications.Items()
	}
	if changed && a.notificationToast != nil {
		a.notificationToast.Hide()
	}
	return a.notifications.Items()
}
