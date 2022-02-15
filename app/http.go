package app

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// FourOFour returns a standard 404 meesage
func FourOFour(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, "404 page not found")
}

// HTTPError returns a standard 404 meesage
func HTTPError(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, msg)
}

// DownloadToString returns a []byte slice from a URL
func DownloadToString(url string) (string, error) {
	// timeout in case URL is unreachable, default 30s
	timeout := time.Duration(15 * time.Second)

	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// req.Header.Set("User-Agent", UserAgent)

	// Get the data
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err := fmt.Errorf("Error downloading %s", url)
		return "", err
	}

	if !strings.Contains(resp.Header.Get("Content-Type"), "text/plain") {
		err := fmt.Errorf("Content type %s not supported (%s)", resp.Header.Get("Content-Type"), url)
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
