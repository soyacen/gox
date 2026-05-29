package httpx

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestToCurl_GET(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com/api/users?page=1&size=10", nil)
	curl := ToCurl(req)
	if !strings.HasPrefix(curl, "curl ") {
		t.Errorf("expected curl prefix, got: %s", curl)
	}
	if !strings.Contains(curl, "http://example.com/api/users?page=1&size=10") {
		t.Errorf("expected URL in curl, got: %s", curl)
	}
	// GET 请求不应包含 -X
	if strings.Contains(curl, "-X") {
		t.Errorf("GET request should not have -X flag, got: %s", curl)
	}
	t.Log(curl)
}

func TestToCurl_POST_WithBody(t *testing.T) {
	body := `{"name":"jack","age":18}`
	req, _ := http.NewRequest(http.MethodPost, "http://example.com/api/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	curl := ToCurl(req)
	if !strings.Contains(curl, "-X POST") {
		t.Errorf("expected -X POST, got: %s", curl)
	}
	// 应使用 --data-raw 而非 -d（避免 @ 开头的值被当作文件路径）
	if !strings.Contains(curl, "--data-raw") {
		t.Errorf("expected --data-raw flag, got: %s", curl)
	}
	if !strings.Contains(curl, `{"name":"jack","age":18}`) {
		t.Errorf("expected body in curl, got: %s", curl)
	}
	if !strings.Contains(curl, "-H") {
		t.Errorf("expected -H flag, got: %s", curl)
	}
	t.Log(curl)
}

func TestToCurl_PUT(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPut, "http://example.com/api/users/1", strings.NewReader(`{"name":"tom"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token123")

	curl := ToCurl(req)
	if !strings.Contains(curl, "-X PUT") {
		t.Errorf("expected -X PUT, got: %s", curl)
	}
	t.Log(curl)
}

func TestToCurl_DELETE(t *testing.T) {
	req, _ := http.NewRequest(http.MethodDelete, "http://example.com/api/users/1", nil)
	curl := ToCurl(req)
	if !strings.Contains(curl, "-X DELETE") {
		t.Errorf("expected -X DELETE, got: %s", curl)
	}
	t.Log(curl)
}

func TestToCurl_PATCH(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPatch, "http://example.com/api/users/1", strings.NewReader(`{"status":"active"}`))
	req.Header.Set("Content-Type", "application/json")

	curl := ToCurl(req)
	if !strings.Contains(curl, "-X PATCH") {
		t.Errorf("expected -X PATCH, got: %s", curl)
	}
	if !strings.Contains(curl, "--data-raw") {
		t.Errorf("expected --data-raw for PATCH body, got: %s", curl)
	}
	t.Log(curl)
}

func TestToCurl_WithCookies(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "abc123"})
	req.AddCookie(&http.Cookie{Name: "lang", Value: "zh"})

	curl := ToCurl(req)
	if !strings.Contains(curl, "-b") {
		t.Errorf("expected -b flag for cookies, got: %s", curl)
	}
	if !strings.Contains(curl, "session=abc123") {
		t.Errorf("expected cookie value, got: %s", curl)
	}
	// Cookie 不应在 -H 中重复出现
	if strings.Count(curl, "session=abc123") > 1 {
		t.Errorf("cookie should not appear in both -H and -b, got: %s", curl)
	}
	t.Log(curl)
}

func TestToCurl_WithHeaders(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Custom", "hello world")

	curl := ToCurl(req)
	if !strings.Contains(curl, "-H 'Accept: application/json'") {
		t.Errorf("expected Accept header, got: %s", curl)
	}
	if !strings.Contains(curl, "-H 'X-Custom: hello world'") {
		t.Errorf("expected X-Custom header, got: %s", curl)
	}
	t.Log(curl)
}

func TestToCurl_HostHeader(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com/path", nil)
	// 当 req.Host 与 URL.Host 不同时，应输出 Host 头
	req.Host = "api.example.com"

	curl := ToCurl(req)
	if !strings.Contains(curl, "-H 'Host: api.example.com'") {
		t.Errorf("expected Host header when req.Host differs from URL.Host, got: %s", curl)
	}
	t.Log(curl)
}

func TestToCurl_HostHeader_SameHost(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com/path", nil)
	// req.Host 与 URL.Host 相同时，不应重复输出 Host 头

	curl := ToCurl(req)
	if strings.Contains(curl, "-H 'Host:") {
		t.Errorf("should not have Host header when req.Host == URL.Host, got: %s", curl)
	}
	t.Log(curl)
}

func TestToCurl_BodyRestore(t *testing.T) {
	bodyContent := `{"key":"value"}`
	req, _ := http.NewRequest(http.MethodPost, "http://example.com", strings.NewReader(bodyContent))

	_ = ToCurl(req)

	// 验证请求体已恢复
	restored, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("failed to read restored body: %v", err)
	}
	if string(restored) != bodyContent {
		t.Errorf("body not restored, expected %q, got %q", bodyContent, string(restored))
	}
}

func TestToCurl_BodyRestore_WithGetBody(t *testing.T) {
	bodyContent := `{"key":"value"}`
	req, _ := http.NewRequest(http.MethodPost, "http://example.com", strings.NewReader(bodyContent))
	// 模拟 http.Request 设置了 GetBody 的情况
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(strings.NewReader(bodyContent)), nil
	}

	curl := ToCurl(req)
	if !strings.Contains(curl, bodyContent) {
		t.Errorf("expected body from GetBody, got: %s", curl)
	}
	t.Log(curl)
}

func TestToCurlWithOptions_Compressed(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	curl := ToCurlWithOptions(req, CurlOptions{Compressed: true})
	if !strings.Contains(curl, "--compressed") {
		t.Errorf("expected --compressed flag, got: %s", curl)
	}
	t.Log(curl)
}

func TestToCurlWithOptions_Insecure(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://example.com", nil)
	curl := ToCurlWithOptions(req, CurlOptions{Insecure: true})
	if !strings.Contains(curl, "-k") {
		t.Errorf("expected -k flag, got: %s", curl)
	}
	t.Log(curl)
}

func TestToCurlWithOptions_FollowRedirects(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	curl := ToCurlWithOptions(req, CurlOptions{FollowRedirects: true})
	if !strings.Contains(curl, "-L") {
		t.Errorf("expected -L flag, got: %s", curl)
	}
	t.Log(curl)
}

func TestToCurlWithOptions_Pretty(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "http://example.com/api", strings.NewReader(`{"a":1}`))
	req.Header.Set("Content-Type", "application/json")
	curl := ToCurlWithOptions(req, CurlOptions{Pretty: true})
	if !strings.Contains(curl, "\\\n") {
		t.Errorf("expected multi-line format, got: %s", curl)
	}
	// flag 和值应在同一行
	if !strings.Contains(curl, "-X POST") {
		t.Errorf("expected -X POST on same line, got: %s", curl)
	}
	if !strings.Contains(curl, "--data-raw '{") {
		t.Errorf("expected --data-raw with body on same line, got: %s", curl)
	}
	t.Log(curl)
}

func TestToCurl_BinaryBody(t *testing.T) {
	binaryData := []byte{0x00, 0x01, 0x02, 0x03, 0xFF}
	req, _ := http.NewRequest(http.MethodPost, "http://example.com/upload", bytes.NewReader(binaryData))
	curl := ToCurl(req)
	if !strings.Contains(curl, "--data-binary") {
		t.Errorf("expected --data-binary for binary body, got: %s", curl)
	}
	// 应有 base64 编码注释
	if !strings.Contains(curl, "# body is base64 encoded") {
		t.Errorf("expected base64 comment for binary data, got: %s", curl)
	}
	t.Log(curl)
}

func TestToCurl_DataRaw_AtSign(t *testing.T) {
	// body 以 @ 开头，必须用 --data-raw 而非 -d，否则 curl 会当作文件路径
	req, _ := http.NewRequest(http.MethodPost, "http://example.com", strings.NewReader("@file.txt"))
	curl := ToCurl(req)
	if !strings.Contains(curl, "--data-raw") {
		t.Errorf("expected --data-raw for @ prefixed body, got: %s", curl)
	}
	if strings.Contains(curl, " -d ") {
		t.Errorf("should not use -d for @ prefixed body, got: %s", curl)
	}
	t.Log(curl)
}

func TestToCurl_ShellEscape(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple", "simple"},
		{"hello world", "'hello world'"},
		{"it's", "'it'\\''s'"},
		{"", "''"},
		{"http://example.com?a=1&b=2", "http://example.com?a=1&b=2"},
		{"value with $dollar", "'value with $dollar'"},
		{"value with \"quotes\"", "'value with \"quotes\"'"},
	}
	for _, tt := range tests {
		result := shellEscape(tt.input)
		if result != tt.expected {
			t.Errorf("shellEscape(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestToCurl_NilBody(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	curl := ToCurl(req)
	if strings.Contains(curl, "--data-raw") {
		t.Errorf("nil body should not have --data-raw flag, got: %s", curl)
	}
	t.Log(curl)
}

func TestToCurl_EmptyMethod(t *testing.T) {
	req, _ := http.NewRequest("", "http://example.com", nil)
	curl := ToCurl(req)
	// 空方法应默认为 GET，不应出现 -X
	if strings.Contains(curl, "-X") {
		t.Errorf("empty method should default to GET without -X, got: %s", curl)
	}
	t.Log(curl)
}

func TestToCurl_AllOptions(t *testing.T) {
	body := `{"query":"test"}`
	req, _ := http.NewRequest(http.MethodPost, "https://api.example.com/search", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer my-token-123")
	req.AddCookie(&http.Cookie{Name: "session", Value: "xyz"})
	req.Host = "api.example.com"

	curl := ToCurlWithOptions(req, CurlOptions{
		Compressed:      true,
		Insecure:        true,
		FollowRedirects: true,
		Pretty:          true,
	})

	// 验证所有关键元素都存在
	checks := []string{
		"-X POST",
		"https://api.example.com/search",
		"-H 'Content-Type: application/json'",
		"-H 'Authorization: Bearer my-token-123'",
		"-b session=xyz",
		"--data-raw",
		"--compressed",
		"-k",
		"-L",
	}
	for _, check := range checks {
		if !strings.Contains(curl, check) {
			t.Errorf("expected %q in curl, got:\n%s", check, curl)
		}
	}
	t.Log(curl)
}
