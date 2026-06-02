package notification

import "time"

func (s *Service) ShouldRunRemoteCheck(now time.Time) (bool, error) {
	last, ok, err := loadMetaTime(s.path, metaKeyLastNotificationsCheck)
	if err != nil {
		return false, err
	}
	if !ok {
		return true, nil
	}
	return now.UTC().Sub(last) >= remoteCheckInterval, nil
}

func (s *Service) SetLastRemoteCheck(at time.Time) error {
	return saveMetaTime(s.path, metaKeyLastNotificationsCheck, at)
}

func (s *Service) RemoteSupportState() (bool, bool, error) {
	return loadMetaBool(s.path, metaKeyRemoteSupported)
}

func (s *Service) SetRemoteSupportState(supported bool) error {
	return saveMetaBool(s.path, metaKeyRemoteSupported, supported)
}
