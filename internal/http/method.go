package http

// Method 表示 HTTP 请求方法（对应 hutool-http 的 Method 枚举）。
type Method string

const (
	MethodGet     Method = "GET"
	MethodPost    Method = "POST"
	MethodHead    Method = "HEAD"
	MethodOptions Method = "OPTIONS"
	MethodPut     Method = "PUT"
	MethodDelete  Method = "DELETE"
	MethodTrace   Method = "TRACE"
	MethodConnect Method = "CONNECT"
	MethodPatch   Method = "PATCH"
)

// String 返回方法字符串。
func (m Method) String() string { return string(m) }
