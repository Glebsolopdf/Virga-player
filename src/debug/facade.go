package debug

import mgr "virga-player/debug/manager"

type Manager = mgr.Manager

func NewManager(enabled, forced bool) *Manager {
	return mgr.NewManager(enabled, forced)
}
