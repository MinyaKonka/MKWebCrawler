package analyzer

import (
	base "core/base"
	"net/http"
)

// 被用于解析HTTP响应的函数类型
type MKParseResponse func(response *http.Response, depth uint32) ([]base.MKData, []error)
