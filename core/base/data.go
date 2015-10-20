package base

import (
	"net/http"
)

// 数据接口
type MKData interface {
	Valid() bool // 数据是否有效
}

/*
 * 请求
 */
type MKRequest struct {
	request *http.Request // HTTP请求指针
	depth   uint32        // 请求深度
}

// 创建新的请求
func NewRequest(request *http.Request, depth uint32) *MKRequest {
	return &MKRequest{
		request: request,
		depth:   depth,
	}
}

// 获取HTTP请求
func (request *MKRequest) Request() *http.Request {
	return request.request
}

// 获取深度值
func (request *MKRequest) Depth() uint32 {
	return request.depth
}

// 数据是否有效
func (request *MKRequest) Valid() bool {
	return request.request != nil && request.request.URL != nil
}

/*
 *	响应
 */
type MKResponse struct {
	response *http.Response
	depth    uint32
}

// 创建新的响应
func NewResponse(response *http.Response, depth uint32) *MKResponse {
	return &MKResponse{
		response: response,
		depth:    depth,
	}
}

// 获取HTTP响应
func (response *MKResponse) Response() *http.Response {
	return response.response
}

// 获取深度值
func (response *MKResponse) Depth() uint32 {
	return response.depth
}

// 数据是否有效
func (response *MKResponse) Valid() bool {
	return response.response != nil && response.response.Body != nil
}

/*
 * 条目
 */
type MKItem map[string]interface{}

// 数据是否有效
func (item MKItem) Valid() bool {
	return item != nil
}
