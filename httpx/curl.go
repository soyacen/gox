package httpx

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
)

// CurlOptions 控制 ToCurl 生成的 curl 命令格式。
type CurlOptions struct {
	// Compressed 是否添加 --compressed 标志
	Compressed bool
	// Insecure 是否添加 -k/--insecure 标志（跳过 TLS 证书验证）
	Insecure bool
	// FollowRedirects 是否添加 -L/--location 标志（跟随重定向）
	FollowRedirects bool
	// Pretty 是否使用多行格式（使用反斜杠换行）
	Pretty bool
}

// ToCurl 将 *http.Request 转换为等价的、可直接在 shell 中运行的 curl 命令字符串。
// 请求体（如果有）会被读取，然后通过 io.NopCloser 恢复，不影响后续使用。
// 对于无法读取的请求体，将忽略 body 部分。
//
// 示例输出：
//
//	curl -X POST 'http://example.com/api' -H 'Content-Type: application/json' --data-raw '{"name":"jack"}'
func ToCurl(req *http.Request) string {
	return ToCurlWithOptions(req, CurlOptions{})
}

// ToCurlWithOptions 将 *http.Request 转换为等价的 curl 命令字符串，支持自定义选项。
func ToCurlWithOptions(req *http.Request, opts CurlOptions) string {
	var parts []string
	parts = append(parts, "curl")

	// 方法（GET 可省略 -X，其他方法显式指定）
	method := req.Method
	if method == "" {
		method = http.MethodGet
	}
	if method != http.MethodGet {
		parts = append(parts, "-X", method)
	}

	// URL（使用 String() 保留完整的 scheme://host/path?query）
	urlStr := req.URL.String()
	parts = append(parts, shellEscape(urlStr))

	// Host 头：Go 的 http.Request 将 Host 存在 req.Host，不在 req.Header 中
	// curl 默认用 URL 中的 Host，仅当 req.Host 与 URL.Host 不同时需要显式设置
	if req.Host != "" && req.Host != req.URL.Host {
		parts = append(parts, "-H", shellEscape("Host: "+req.Host))
	}

	// Headers（按 key 排序以保证输出确定性）
	// 跳过 Cookie 头（已通过 -b 单独处理）和 Host 头（已通过 -H Host: 处理）
	hasCookies := len(req.Cookies()) > 0
	if len(req.Header) > 0 {
		keys := make([]string, 0, len(req.Header))
		for k := range req.Header {
			ck := http.CanonicalHeaderKey(k)
			if hasCookies && ck == "Cookie" {
				continue
			}
			if ck == "Host" {
				continue
			}
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			for _, v := range req.Header[k] {
				parts = append(parts, "-H", shellEscape(k+": "+v))
			}
		}
	}

	// Cookies
	if len(req.Cookies()) > 0 {
		var cookieParts []string
		for _, c := range req.Cookies() {
			cookieParts = append(cookieParts, c.Name+"="+c.Value)
		}
		parts = append(parts, "-b", shellEscape(strings.Join(cookieParts, "; ")))
	}

	// Body：使用 --data-raw（避免 -d 将 @ 开头的值当作文件路径）
	// 二进制数据使用 base64 编码并通过 shell 解码管道传递
	if req.Body != nil {
		bodyBytes, err := readAndRestoreBody(req)
		if err == nil && len(bodyBytes) > 0 {
			if isBinaryData(bodyBytes) {
				// 二进制数据：base64 编码后通过管道传给 curl --data-binary @-
				encoded := base64.StdEncoding.EncodeToString(bodyBytes)
				parts = append(parts, "--data-binary", shellEscape(encoded))
				// 添加注释提示这是 base64 编码的数据
				parts = append([]string{fmt.Sprintf("# body is base64 encoded (%d bytes binary data)", len(bodyBytes))}, parts...)
			} else {
				parts = append(parts, "--data-raw", shellEscape(string(bodyBytes)))
			}
		}
	}

	// 选项标志
	if opts.Compressed {
		parts = append(parts, "--compressed")
	}
	if opts.Insecure {
		parts = append(parts, "-k")
	}
	if opts.FollowRedirects {
		parts = append(parts, "-L")
	}

	// 组装
	if opts.Pretty {
		return prettyJoin(parts)
	}
	return strings.Join(parts, " ")
}

// readAndRestoreBody 读取请求体并恢复，保证后续还能再次读取。
func readAndRestoreBody(req *http.Request) ([]byte, error) {
	// 优先使用 GetBody（*http.Request 在设置 body 时可能已设置 GetBody）
	if req.GetBody != nil {
		rc, err := req.GetBody()
		if err != nil {
			return nil, err
		}
		defer rc.Close()
		return io.ReadAll(rc)
	}
	// 回退：直接读取并恢复
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	return bodyBytes, nil
}

// shellEscape 对字符串进行 shell 转义，确保可直接粘贴到 shell 运行。
// 规则：空字符串输出两个单引号；只含安全字符原样输出；含特殊字符用单引号包裹，内部单引号替换。
func shellEscape(s string) string {
	if s == "" {
		return "''"
	}
	safe := true
	for _, c := range s {
		if !isShellSafe(c) {
			safe = false
			break
		}
	}
	if safe {
		return s
	}
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

// isShellSafe 判断字符是否为 shell 安全字符（不需要引号包裹）。
func isShellSafe(c rune) bool {
	if c >= 'a' && c <= 'z' {
		return true
	}
	if c >= 'A' && c <= 'Z' {
		return true
	}
	if c >= '0' && c <= '9' {
		return true
	}
	switch c {
	case '-', '_', '.', '/', ':', '@', '=', ',', '?', '&', '#', '%', '+', '~':
		return true
	}
	return false
}

// prettyJoin 将 curl 命令各部分按 flag+value 分组，用反斜杠换行连接。
// 例如：
//
//	curl \
//	  -X POST \
//	  'http://example.com' \
//	  -H 'Content-Type: application/json'
func prettyJoin(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	// flags 需要和下一个参数绑定在一起（不出现在行尾单独一行）
	flagSet := map[string]bool{
		"-X": true, "-H": true, "-d": true, "-b": true,
		"--data-binary": true, "--data-raw": true, "--data": true,
		"-e": true, "-u": true, "--user": true, "-o": true, "-w": true,
	}
	var lines []string
	i := 0
	for i < len(parts) {
		if flagSet[parts[i]] && i+1 < len(parts) {
			lines = append(lines, parts[i]+" "+parts[i+1])
			i += 2
		} else {
			lines = append(lines, parts[i])
			i++
		}
	}
	return strings.Join(lines, " \\\n  ")
}

// isBinaryData 判断数据是否为二进制数据（包含不可打印的控制字符）。
func isBinaryData(data []byte) bool {
	for _, b := range data {
		if b < 0x09 || (b > 0x0d && b < 0x20 && b != 0x1b) {
			return true
		}
	}
	return false
}
