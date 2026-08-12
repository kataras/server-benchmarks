package runner

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"
)

// downloadTimeout bounds a BodyFile download as a whole.
const downloadTimeout = 30 * time.Second

// fetchBodyFile downloads rawURL into dir and returns the local file path.
// The URL's file extension is preserved (a later Content-Type inference
// depends on it).
func fetchBodyFile(ctx context.Context, rawURL, dir string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid BodyFile URL %q: %w", rawURL, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return "", fmt.Errorf("request BodyFile %q: %w", rawURL, err)
	}

	client := &http.Client{Timeout: downloadTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("download BodyFile %q: %w", rawURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download BodyFile %q: unexpected status %s", rawURL, resp.Status)
	}

	ext := path.Ext(u.Path)
	if ext == "" {
		ext = ".bin"
	}

	name := filepath.Join(dir, "body"+ext)

	f, err := os.Create(name)
	if err != nil {
		return "", fmt.Errorf("create BodyFile %q: %w", name, err)
	}

	if _, err := io.Copy(f, resp.Body); err != nil {
		f.Close()
		return "", fmt.Errorf("save BodyFile %q: %w", name, err)
	}

	if err := f.Close(); err != nil {
		return "", fmt.Errorf("save BodyFile %q: %w", name, err)
	}

	return name, nil
}
