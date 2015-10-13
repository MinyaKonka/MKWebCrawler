package base

import (
	"bytes"
	"fmt"
)

// 错误域
type ErrorDomain string

// 错误编码
type ErrorCode int

// 错误域常量
const (
	ERR_DOMAIN_DOWNLOADER     ErrorDomain = "Downloader Error"
	ERR_DOMAIN_ANALYZER       ErrorDomain = "Analyzer Error"
	ERR_DOMAIN_ITEM_PROCESSOR ErrorDomain = "Item Processor Error"
)

// 错误编码常量
const (
	ERR_CODE_NONE ErrorCode = 0 // 无错误
)

// 错误接口
type MKError interface {
	Domain() ErrorDomain // 获取错误类型
	Code() ErrorCode     // 获取错误编码
	Error() string       // 获取错误提示信息
}

// 错误实现
type mk_error struct {
	domain      ErrorDomain // 错误域
	code        ErrorCode   // 错误编码
	message     string      // 错误提示信息
	fullMessage string      // 完整错误提示信息
}

// 创建一个新的爬虫错误
func NewError(domain ErrorDomain, code ErrorCode, message string) MKError {
	return &mk_error{
		domain:  domain,
		code:    code,
		message: message,
	}
}

// 获取错误域
func (err *mk_error) Domain() ErrorDomain {
	return err.domain
}

// 获取错误编码
func (err *mk_error) Code() ErrorCode {
	return err.code
}

// 获取错误提示信息
func (err *mk_error) Error() string {
	if err.fullMessage == "" {
		err.generateFullErrorMessage()
	}

	return err.fullMessage
}

// 生成错误提示信息
func (err *mk_error) generateFullErrorMessage() {
	var buffer bytes.Buffer

	buffer.WriteString("Error: ")

	if err.domain != "" && err.code != ERR_CODE_NONE {
		buffer.WriteString(string(err.domain))
		buffer.WriteString("[%d]", err.code)
		buffer.WriteString(": ")
	}

	buffer.WriteString(err.message)

	err.fullMessage = fmt.Sprintf("%s\n", buffer.String())
}
