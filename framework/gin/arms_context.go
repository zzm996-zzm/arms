package gin

import (
	"context"
)

var _ IRequest = new(Context)
var _ IResponse = new(Context)
var _ context.Context = new(Context)

//TODO:支持配置文件
// const defaultMultipartMemory = 32 << 20 // 32 MB

func (ctx *Context) BaseContext() context.Context {
	return ctx.Request.Context()
}
