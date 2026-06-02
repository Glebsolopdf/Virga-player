package music

import playerctlcmd "virga-player/music/playerctlcmd"

func getPlayerctlTrack() *TrackInfo {
	if !playerctlcmd.Available() {
		return nil
	}

	metadata, ok := playerctlcmd.ReadTrackMetadata()
	if !ok {
		return getPlayerctlTrackFallback()
	}
	artworkURL := metadata.ArtworkURL
	if artworkURL == "" {
		artworkURL = getArtworkURLFromTrackURL(metadata.TrackURL)
	}

	return &TrackInfo{
		Title:      metadata.Title,
		Artist:     metadata.Artist,
		Album:      metadata.Album,
		Duration:   metadata.Duration,
		Elapsed:    playerctlcmd.PositionSeconds(),
		Paused:     playerctlcmd.IsPaused(),
		ArtworkURL: artworkURL,
		Source:     "playerctl",
	}
}

func getPlayerctlTrackFallback() *TrackInfo {
	title := playerctlcmd.FirstMetadataValue("xesam:title", "title")
	if title == "" {
		return nil
	}

	artist := playerctlcmd.FirstMetadataValue("xesam:artist", "artist")
	album := playerctlcmd.FirstMetadataValue("xesam:album", "album")

	artworkURL := getArtworkURL()

	duration := playerctlcmd.DurationFromMicros(playerctlcmd.FirstMetadataValue("mpris:length"))
	if duration == 0 {
		duration = playerctlcmd.DurationFromMicros(playerctlcmd.FirstMetadataValue("xesam:length"))
	}

	return &TrackInfo{
		Title:      title,
		Artist:     artist,
		Album:      album,
		Duration:   duration,
		Elapsed:    playerctlcmd.PositionSeconds(),
		Paused:     playerctlcmd.IsPaused(),
		ArtworkURL: artworkURL,
		Source:     "playerctl",
	}
}
