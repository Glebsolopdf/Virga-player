package notification

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

const toastMessage = "You have unread notifications. Open Settings > Notifications."

const notificationsIntroID = "info:notifications-intro"

type Service struct {
	path  string
	state State
	now   func() time.Time
}

func Load() (*Service, error) {
	path := StatePath()
	state, err := loadState(path)
	if err != nil {
		return nil, err
	}
	return &Service{
		path:  path,
		state: state,
		now:   time.Now,
	}, nil
}

func NewForPath(path string) *Service {
	return &Service{
		path: path,
		state: State{
			Items: []Item{},
		},
		now: time.Now,
	}
}

func (s *Service) ClearAll() error {
	s.state.Items = []Item{}
	return clearState(s.path)
}

func (s *Service) Items() []Item {
	items := append([]Item(nil), s.state.Items...)
	slices.SortFunc(items, func(left, right Item) int {
		if left.CreatedAt.Equal(right.CreatedAt) {
			return strings.Compare(right.ID, left.ID)
		}
		if left.CreatedAt.After(right.CreatedAt) {
			return -1
		}
		return 1
	})
	return items
}

func (s *Service) UnreadCount() int {
	count := 0
	for _, item := range s.state.Items {
		if !item.IsRead() {
			count++
		}
	}
	return count
}

func (s *Service) HasUnread() bool {
	return s.UnreadCount() > 0
}

func (s *Service) ToastMessage() string {
	return toastMessage
}

func (s *Service) EnsureWelcome(firstRun bool, currentVersion string) (Item, bool, error) {
	if !firstRun {
		return Item{}, false, nil
	}

	item := Item{
		ID:        welcomeID(currentVersion),
		Kind:      KindWelcome,
		Title:     "Welcome to Virga Player",
		Body:      welcomeBody(currentVersion),
		Version:   currentVersion,
		CreatedAt: s.now().UTC(),
	}

	added, err := s.upsert(item)
	return item, added, err
}

func (s *Service) EnsureNotificationsIntro() (Item, bool, error) {
	item := Item{
		ID:        notificationsIntroID,
		Kind:      KindInfo,
		Title:     "Welcome to the notifications section!",
		Body:      "Here you will receive important information regarding updates, lyrics files, and a few logs. This feature is not fully implemented yet!",
		CreatedAt: s.now().UTC(),
	}

	added, err := s.upsert(item)
	return item, added, err
}

func (s *Service) RemoveUpdateNotices() (bool, error) {
	filtered := make([]Item, 0, len(s.state.Items))
	changed := false
	for _, item := range s.state.Items {
		if item.Kind == KindUpdate {
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

func (s *Service) AddUpdateNotice(latestVersion string) (Item, bool, error) {
	version := strings.TrimSpace(latestVersion)
	if version == "" {
		return Item{}, false, nil
	}

	item := Item{
		ID:        updateID(version),
		Kind:      KindUpdate,
		Title:     "New Virga Player version available",
		Body:      fmt.Sprintf("A new Virga Player version (%s) is available. Please update.", version),
		Version:   version,
		CreatedAt: s.now().UTC(),
	}

	added, err := s.upsert(item)
	return item, added, err
}

func (s *Service) MarkAllRead() (bool, error) {
	now := s.now().UTC()
	changed := false
	for index := range s.state.Items {
		if s.state.Items[index].ReadAt != nil {
			continue
		}
		stamp := now
		s.state.Items[index].ReadAt = &stamp
		changed = true
	}
	if !changed {
		return false, nil
	}
	return true, saveState(s.path, s.state)
}

func (s *Service) upsert(item Item) (bool, error) {
	for _, existing := range s.state.Items {
		if existing.ID == item.ID {
			return false, nil
		}
	}
	s.state.Items = append(s.state.Items, item)
	return true, saveState(s.path, s.state)
}

func welcomeID(currentVersion string) string {
	version := strings.TrimSpace(currentVersion)
	if version == "" {
		version = "unknown"
	}
	return "welcome:" + version
}

func updateID(version string) string {
	return "update:" + strings.TrimSpace(version)
}

func welcomeBody(currentVersion string) string {
	return fmt.Sprintf(
		"Welcome to Virga Player! This terminal utility turns your currently playing music into a rain-driven visual scene with optional player and synced lyrics. How to use it: start Virga, press S to open settings, use the arrow keys to move through categories, and press Enter to save and exit settings. In Settings you can change rain behavior, audio reactivity, player visuals, lyrics options, and debug mode. Current version: %s.",
		strings.TrimSpace(currentVersion),
	)
}
