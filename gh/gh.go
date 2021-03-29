package gh

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type releaseResponse struct {
	TagName string `json:"tag_name"`
}

// GetLatestReleaseTagName find the latest release tag name.
func GetLatestReleaseTagName(githubURL, githubToken, owner, repositoryName string) (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	latestReleaseURL := fmt.Sprintf("https://%s/repos/%s/%s/releases/latest", githubURL, owner, repositoryName)

	req, err := http.NewRequest("GET", latestReleaseURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to construct new request for latest release: %w", err)
	}

	if len(githubToken) != 0 {
		req.Header.Add("Authorization", "Bearer "+githubToken)
	}

	resp, err := client.Do(req)
	if err != nil {
		logResponseBody(resp)
		return "", fmt.Errorf("failed to get latest release tag name: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		logResponseBody(resp)
		return "", fmt.Errorf("failed to get latest release tag name on GitHub (%q), status: %s", latestReleaseURL, resp.Status)
	}

	var release releaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("failed to decode github response for latest release: %w", err)
	}

	if err := resp.Body.Close(); err != nil {
		return "", fmt.Errorf("closing github response body stream: %w", err)
	}

	return release.TagName, nil
}

func logResponseBody(resp *http.Response) {
	if resp.Body == nil {
		log.Println("The response body is empty")
		return
	}

	defer safeClose(resp.Body.Close)

	body, errBody := ioutil.ReadAll(resp.Body)
	if errBody != nil {
		log.Println(errBody)
		return
	}

	log.Println("Body:", string(body))
}

func safeClose(fn func() error) {
	if err := fn(); err != nil {
		log.Println(err)
	}
}
