package itempipeline

import (
	base "core/base"
)

// 被用来处理条目的函数类型
type MKProcessItem func(item base.MKItem) (result base.Item, err error)
