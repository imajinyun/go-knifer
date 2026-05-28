// Package gkhttp 对应 hutool-http，提供 HTTP 客户端、下载、Cookie、UserAgent、SimpleServer 等工具。
//
// 与 hutool-http 不同，本包基于 Go 标准库 net/http 进行二次封装，提供链式 API：
//
//	body := gkhttp.Get("https://example.com").Execute().Body()
//	resp := gkhttp.NewRequest(gkhttp.MethodPost, url).
//	            Form(map[string]any{"a": 1}).
//	            Timeout(5 * time.Second).
//	            Execute()
package http
