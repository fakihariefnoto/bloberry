package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client talks to a Bloberry server over the standard envelope
// {data?, messages?: [{code, content}]}, mirroring web/src/lib/api.ts.
type Client struct {
	BaseURL string
	Token   string
	HTTP    *http.Client
}

func NewClient(baseURL, token string) *Client {
	return &Client{
		BaseURL: strings.TrimSuffix(baseURL, "/"),
		Token:   token,
		HTTP:    &http.Client{Timeout: 60 * time.Second},
	}
}

type Message struct {
	Code    string `json:"code"`
	Content string `json:"content"`
}

type Envelope struct {
	Data     json.RawMessage `json:"data,omitempty"`
	Messages []Message       `json:"messages,omitempty"`
}

type HTTPError struct {
	Status  int
	Code    string
	Content string
}

func (e *HTTPError) Error() string { return e.Content }

// Do performs a request, decodes the envelope, and maps HTTP errors to exit
// codes (cli/README.md §Exit codes).
func (c *Client) Do(method, path string, body io.Reader, contentType string) (*Envelope, int, error) {
	req, err := http.NewRequest(method, c.BaseURL+path, body)
	if err != nil {
		return nil, ExitFailed, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, ExitBackendDown, exit(ExitBackendDown, "storage backend unreachable: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return &Envelope{}, resp.StatusCode, nil
	}

	env := &Envelope{}
	raw, _ := io.ReadAll(resp.Body)
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, env)
	}

	if resp.StatusCode >= 400 {
		code := "error"
		content := http.StatusText(resp.StatusCode)
		if len(env.Messages) > 0 {
			code = env.Messages[0].Code
			content = env.Messages[0].Content
		}
		return env, exitCodeFor(resp.StatusCode, code), &HTTPError{Status: resp.StatusCode, Code: code, Content: content}
	}
	return env, resp.StatusCode, nil
}

func exitCodeFor(status int, code string) int {
	switch {
	case status == http.StatusUnauthorized || code == "unauthorized" || code == "refresh_invalid" ||
		code == "key_revoked" || code == "key_expired" || code == "totp_required":
		return ExitNotAuthed
	case status == http.StatusForbidden || code == "forbidden" || code == "no_invitation":
		return ExitForbidden
	case status == http.StatusNotFound || code == "not_found":
		return ExitNotFound
	case status == http.StatusUnprocessableEntity && code == "quota_exceeded":
		return ExitQuota
	case status == http.StatusConflict:
		return ExitConflict
	case status == http.StatusBadGateway || status == http.StatusBadRequest && code == "backend_unreachable":
		return ExitBackendDown
	default:
		return ExitFailed
	}
}

// JSON decodes the envelope's data field.
func (e *Envelope) JSON(v interface{}) error {
	if len(e.Data) == 0 {
		return nil
	}
	return json.Unmarshal(e.Data, v)
}

var _ = fmt.Sprintf
var _ = errors.New
