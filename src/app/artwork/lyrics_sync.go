package artwork

func (a *Artwork) SetSyncedLyrics(lyrics string) {
	a.mu.Lock()
	if lyrics == a.lyricsRaw {
		a.mu.Unlock()
		return
	}

	a.lyricsRaw = lyrics
	a.lyricsParseToken++
	parseToken := a.lyricsParseToken
	if lyrics == "" {
		a.Lyrics = nil
		a.mu.Unlock()
		return
	}

	a.Lyrics = nil
	a.mu.Unlock()

	go func(raw string, token uint64) {
		parsed := parseSyncedLyrics(raw)

		a.mu.Lock()
		defer a.mu.Unlock()
		if token != a.lyricsParseToken || raw != a.lyricsRaw {
			return
		}
		a.Lyrics = parsed
	}(lyrics, parseToken)
}
