package music

func getPlayerctlTrack() *TrackInfo {
	if !playerctlAvailable() {
		return nil
	}

	meta, ok := readTrackMeta()
	if !ok {
		return getPlayerctlTrackFallback()
	}
	artworkURL := meta.ArtworkURL
	if artworkURL == "" {
		artworkURL = getArtworkURLFromTrackURL(meta.TrackURL)
	}

	return &TrackInfo{
		Title:      meta.Title,
		Artist:     meta.Artist,
		Album:      meta.Album,
		Duration:   meta.Duration,
		Elapsed:    playerctlPosition(),
		Paused:     playerctlIsPaused(),
		ArtworkURL: artworkURL,
		Source:     "playerctl",
	}
}

func getPlayerctlTrackFallback() *TrackInfo {
	title := firstMetaValue("xesam:title", "title")
	if title == "" {
		return nil
	}

	artist := firstMetaValue("xesam:artist", "artist")
	album := firstMetaValue("xesam:album", "album")
	artworkURL := getArtworkURL()

	duration := durationFromMicros(firstMetaValue("mpris:length"))
	if duration == 0 {
		duration = durationFromMicros(firstMetaValue("xesam:length"))
	}

	return &TrackInfo{
		Title:      title,
		Artist:     artist,
		Album:      album,
		Duration:   duration,
		Elapsed:    playerctlPosition(),
		Paused:     playerctlIsPaused(),
		ArtworkURL: artworkURL,
		Source:     "playerctl",
	}
}
