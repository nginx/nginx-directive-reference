package atom

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

type entry struct {
	Link struct {
		Href string `xml:"href,attr"`
	} `xml:"link"`
}

type feed struct {
	Entry []entry `xml:"entry"`
}

type config struct {
	Client http.Client
}
type Option = func(*config)

// WithHttpClient uses the provided client instead of the default for
// downloading tarballs.
func WithHttpClient(c http.Client) Option {
	return func(o *config) { o.Client = c }
}

// GetVersion gets the first link of the first entry of the XML file
func GetVersion(ctx context.Context, url string, opts ...Option) (string, error) {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}
	body, _ := openURL(ctx, url, cfg.Client)
	return parseXML(body)
}
func openURL(ctx context.Context, url string, client http.Client) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to download %s: %w", url, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to download %s: %s", url, res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %w", err)
	}
	return body, nil
}
func parseXML(XMLContent []byte) (string, error) {
	var f feed
	err := xml.Unmarshal(XMLContent, &f)
	if err != nil {
		return "", fmt.Errorf("unable to parse XML: %w", err)
	}
	if len(f.Entry) > 0 {
		return f.Entry[0].Link.Href, nil
	}
	return "", fmt.Errorf("no entry was found in the feed")
}
