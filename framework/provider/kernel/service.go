package kernel

import (
	"net/http"

	"github.com/zzm996-zzm/arms/framework/contract"
	"github.com/zzm996-zzm/arms/framework/gin"
)

var _ contract.Kernel = new(ArmsKernelService)

type ArmsKernelService struct {
	engine *gin.Engine
}

func NewArmsKernelService(params ...interface{}) (interface{}, error) {
	httpEngine := params[0].(*gin.Engine)
	return &ArmsKernelService{engine: httpEngine}, nil
}

func (k *ArmsKernelService) HttpEngine() http.Handler {
	return k.engine
}
