package main

import (
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"
)

// baseURL is where the running stack is expected. The app must already be
// up, e.g. via `docker compose up -d`, before running this test.
func baseURL() string {
	if v := os.Getenv("APP_URL"); v != "" {
		return v
	}
	return "http://localhost:5000"
}

// waitForApp polls the endpoint until it responds or the timeout elapses,
// giving `docker compose up -d` time to finish starting the containers.
func waitForApp(t *testing.T, url string) {
	t.Helper()
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	t.Fatalf("app at %s did not become ready; is it running via `docker compose up -d`?", url)
}

func get(t *testing.T, url string) string {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("GET %s failed: %v", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET %s returned status %d", url, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading body failed: %v", err)
	}
	return string(body)
}

var hitCountRe = regexp.MustCompile(`Hello World! I have been seen (\d+) times\.`)

func parseCount(t *testing.T, body string) int {
	t.Helper()
	m := hitCountRe.FindStringSubmatch(body)
	if m == nil {
		t.Fatalf("unexpected response body: %q", body)
	}
	n, err := strconv.Atoi(m[1])
	if err != nil {
		t.Fatalf("could not parse count from %q: %v", body, err)
	}
	return n
}

// TestHelloEndpoint verifies the running stack responds with the greeting.
func TestHelloEndpoint(t *testing.T) {
	url := baseURL() + "/"
	waitForApp(t, url)

	body := get(t, url)
	if !hitCountRe.MatchString(body) {
		t.Fatalf("response %q did not match expected greeting", body)
	}
}

// TestHitCounterIncrements verifies the redis-backed counter goes up
// across requests, exercising both the web and redis containers.
func TestHitCounterIncrements(t *testing.T) {
	url := baseURL() + "/"
	waitForApp(t, url)

	first := parseCount(t, get(t, url))
	second := parseCount(t, get(t, url))

	if second <= first {
		t.Fatalf("expected hit count to increase, got %d then %d", first, second)
	}
}
