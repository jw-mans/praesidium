package monitor

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func GetExternalIP(ipCheckURL string) (string, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(ipCheckURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch external IP: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	return string(body), nil
}
