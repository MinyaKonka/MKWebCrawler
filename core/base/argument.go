package base

import (
	"errors"
	"fmt"
)

// 参数容器的接口
type MKArguments interface {
	// 自检参数的有效性，并在必要时返回可以说明问题的错误值
	// 若结果值为nil，则说明未发现问题，否则就意味着自检未通过。
	Check() error

	// 获得参数容器的字符串表现形式
	String() string
}

// 通道参数的容器的描述模板
var channelArgumentsTemplate string = "{ request channel length: %d, " +
	"response channel length: %d," +
	"item channel length: %d," +
	"error channel length: %d }"

// 通道参数的容器
type ChannelArguments struct {
	requestChannelLength  uint   // 请求通道的长度
	responseChannelLength uint   // 响应通道的长度
	itemChannelLength     uint   // 条目通道的长度
	errorChannelLength    uint   // 错误通道的长度
	description           string // 描述
}

// 创建通道参数的容器
func NewChannelArguments(
	requestChannelLength uint,
	responseChannelLength uint,
	itemChannelLength uint,
	errorChannelLength uint) ChannelArguments {

	return &ChannelArguments{
		requestChannelLength:  requestChannelLength,
		responseChannelLength: responseChannelLength,
		itemChannelLength:     itemChannelLength,
		errorChannelLength:    errorChannelLength,
	}
}

func (arguments *ChannelArguments) Check() error {
	if arguments.requestChannelLength == 0 {
		return errors.New("请求通道的容量不能为0！\n")
	}

	if arguments.responseChannelLength == 0 {
		return errors.New("响应通道的容量不能为0！\n")
	}

	if arguments.itemChannelLength == 0 {
		return errors.New("条目通道容量不能为0！\n")
	}

	if arguments.errorChannelLength == 0 {
		return errors.New("错误通道容量不能为0！\n")
	}

	return nil
}

func (arguments *ChannelArguments) String() string {
	if arguments.description == "" {
		arguments.description =
			fmt.Sprintf(channelArgumentsTemplate,
				arguments.requestChannelLength,
				arguments.responseChannelLength,
				arguments.itemChannelLength,
				arguments.errorChannelLength)
	}

	return arguments.description
}

// 获取请求通道的长度
func (arguments *ChannelArguments) RequestChannelLength() uint {
	return arguments.requestChannelLength
}

// 获取响应通道的长度
func (arguments *ChannelArguments) ResponseChannelLength() uint {
	return arguments.responseChannelLength
}

// 获得条目通道的长度
func (argument *ChannelArguments) ItemChannelLength() uint {
	return argument.itemChannelLength
}

// 获取错误通道的长度
func (argument *ChannelArguments) ErrorChannelLength() uint {
	return argument.errorChannelLength
}

// 池基本参数描述模板
var poolArgumentsTemplate string = "{ page downloader pool size: %d, analyzer pool size: %d }"

// 池基本参数容器
type PoolArguments struct {
	pageDownloaderPoolSize uint32 // 网页下载池的尺寸
	analyzerPoolSize       uint32 // 分析器池的尺寸
	description            string // 描述
}

// 创建池基本参数的容器
func NewPoolArguments(
	pageDownloaderPoolSize uint32,
	analyzerPoolSize uint32) PoolArguments {

	return PoolArguments{
		pageDownloaderPoolSize: pageDownloaderPoolSize,
		analyzerPoolSize:       analyzerPoolSize,
	}
}

func (arguments *PoolArguments) Check() error {

	if arguments.pageDownloaderPoolSize == 0 {
		return errors.New("网页下载池大小不能为0！\n")
	}

	if arguments.analyzerPoolSize == 0 {
		return errors.New("分析器池大小不能为0！\n")
	}

	return nil
}

func (arguments *PoolArguments) String() string {

	if arguments.description == "" {
		arguments.description =
			fmt.Sprintf(poolArgumentsTemplate,
				arguments.pageDownloaderPoolSize,
				arguments.analyzerPoolSize)
	}

	return arguments.description
}

// 获得网页下载器池的尺寸
func (arguments *PoolArguments) PageDownloaderPoolSize() uint32 {
	return arguments.pageDownloaderPoolSize
}

// 获得分析器池的尺寸
func (arguments *PoolArguments) AnalyzerPoolSize() uint32 {
	return arguments.analyzerPoolSize
}
