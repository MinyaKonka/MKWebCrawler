package downloader

import (
	base "MKWebCrawler/base"
	middleware "MKWebCrawler/middleware"
	"logging"
	"net/http"
)

// 日志记录器
var logger logging.Logger = base.NewLogger()

// ID生成器
var downloaderIDGenerator middleware.IDGenerator32 = middleware.NewIDGenerator32()

// 生成并返回ID
func generateDownloaderID() uint32 {
	return downloaderIDGenerator.GetUint32()
}

// 网页下载器接口
type MKPageDownloader interface {
	ID() uint32                                                // 获得ID
	Download(request base.MKRequest) (*base.MKResponse, error) // 根据请求下载网页并返回响应
}

// 创建网页下载器
func NewPageDownloader(client *http.Client) MKPageDownload {

	id := generateDownloaderID()
	if client == nil {
		client = &http.Client()
	}

	return &mk_PageDownloader{
		id:         id,
		httpClient: *client,
	}
}

// 网页下载器实现类型
type mk_PageDownloader struct {
	id         uint32      // ID
	httpClient http.Client // HTTP客户端
}

func (downloader *mk_PageDownloader) ID() uint32 {
	return downloader.id
}

func (downloader *mk_PageDownloader) Download(request base.MKRequest) (*base.MKResponse, error) {
	httpRequest := request.Request()
	logger.Infof("请求【url = %s】\n", httpRequest.URL)

	httpResponse, err := downloader.httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	return base.NewResponse(httpResponse, request.Depth()), nil
}
