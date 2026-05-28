package artworklookup

import (
	"html"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var ogImageRegexp = regexp.MustCompile(`(?i)property=["']og:image["'][^>]*content=["']([^"']+)["']|content=["']([^"']+)["'][^>]*property=["']og:image["']`)
var twitterImageRegexp = regexp.MustCompile(`(?i)name=["']twitter:image["'][^>]*content=["']([^"']+)["']|content=["']([^"']+)["'][^>]*name=["']twitter:image["']`)

func fetchTrackPageArtwork(trackURL string) string {
	client := &http.Client{Timeout: 4 * time.Second}
	req, err := http.NewRequest(http.MethodGet, trackURL, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ""
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil {
		return ""
	}

	return extractMetaImageURL(string(body))
}

func extractMetaImageURL(htmlBody string) string {
	for _, re := range []*regexp.Regexp{ogImageRegexp, twitterImageRegexp} {
		matches := re.FindStringSubmatch(htmlBody)
		if len(matches) < 2 {
			continue
		}

		for i := 1; i < len(matches); i++ {
			candidate := strings.TrimSpace(html.UnescapeString(matches[i]))
			if strings.HasPrefix(candidate, "http://") || strings.HasPrefix(candidate, "https://") {
				return candidate
			}
		}
	}

	return ""
}
