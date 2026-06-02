package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func FetchRemoteFeed(ctx context.Context, url string) (RemoteFeed, error) {
	trimmedURL := strings.TrimSpace(url)
	if trimmedURL == "" {
		return RemoteFeed{}, fmt.Errorf("remote notifications URL is empty")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, trimmedURL, nil)
	if err != nil {
		return RemoteFeed{}, err
	}

	client := &http.Client{Timeout: remoteRequestTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return RemoteFeed{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return RemoteFeed{}, fmt.Errorf("remote notifications request failed: status=%d", resp.StatusCode)
	}

	var feed RemoteFeed
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&feed); err != nil {
		return RemoteFeed{}, err
	}
	return feed, nil
}

func RemoteCheckTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), remoteRequestTimeout+5*time.Second)
}
