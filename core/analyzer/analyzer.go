package analyzer

import (
	base "core/base"
	middleware "core/middleware"
	"errors"
	"fmt"
	"logging"
	"net/url"
)

// 日志记录器
var logger logging.Logger = base.NewLogger()

// ID生成器
var analyzerIDGenerator middle.MKIDGenerator = middle.NewIDGenerator32()

// 生成并返回ID
func generatorAnalyzerId() uint32 {
	return analyzerIDGenerator.GetUint32()
}

// 分析器接口
type MKAnalyzer interface {
	ID() uint32                                                                           // 获得ID
	Analyze(parsers []MKParseResponse, response base.MKResponse) ([]base.MKData, []error) // 根据规则分析响应并返回请求和条目
}

// 创建分析器
func NewAnalyzer() MKAnalyzer {
	return &mk_analyzer{id: generatorAnalyzerId()}
}

// 分析器的实现类型
type mk_analyzer struct {
	id uint32 // ID
}

func (analyzer *mk_analyzer) ID() uint32 {
	return analyzer.id
}

func (analyzer *mk_analyzer) Analyze(
	parsers []MKParseResponse,
	response base.MKResponse) (dataList []base.MKData, errorList []error) {

	if parsers == nil {
		err := errors.New("响应解析器列表无效！")
		return nil, err
	}

	httpResponse := response.Response()
	if httpResponse == nil {
		err := errors.New("响应无效")
		return nil, err
	}

	var requestURL *url.URL = httpResponse.Request.URL
	logger.Infof("解析响应(request url = %s)... \n", requestURL)
	depth := response.Depth()

	// 解析HTTP响应
	dataList = make([]base.MKData, 0)
	errorList = make([]error, 0)

	for i, parse := range parsers {

		if parse == nil {
			err := errors.New(fmt.Sprintf("文档解析器[%d]无效！", i))
			errorList = append(errorList, err)
			continue
		}

		pDataList, pErrorList := parse(httpResponse, depth)
		if pDataList != nil {
			for _, data := range pDataList {
				dataList = appendDataList(dataList, data, depth)
			}
		}

		if errorList != nil {
			for _, err := range pErrorList {
				errorList = append(errorList, err)
			}
		}
	}

	return dataList, errorList
}

func appendDataList(dataList []base.MKData, data base.MKData, depth uint32) []base.MKData {

	if data == nil {
		return dataList
	}

	request, ok := data.(*base.MKRequest)
	if !ok {
		return append(dataList, data)
	}

	newDepth := depth + 1
	if request.Depth() != newDepth {
		request = base.NewRequest(request.Request(), newDepth)
	}

	return append(dataList, request)
}

func appendErrorList(errorList []error, err error) []error {
	if err == nil {
		return errorList
	}

	return append(errorList, err)
}
