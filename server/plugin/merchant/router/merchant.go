package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

var Merchant = new(merchant)

type merchant struct {}

// Init 初始化 商户信息 路由信息
func (r *merchant) Init(public *gin.RouterGroup, private *gin.RouterGroup) {
	{
	    group := private.Group("merchant").Use(middleware.OperationRecord())
		group.POST("createMerchant", apiMerchant.CreateMerchant)   // 新建商户信息
		group.DELETE("deleteMerchant", apiMerchant.DeleteMerchant) // 删除商户信息
		group.DELETE("deleteMerchantByIds", apiMerchant.DeleteMerchantByIds) // 批量删除商户信息
		group.PUT("updateMerchant", apiMerchant.UpdateMerchant)    // 更新商户信息
	}
	{
	    group := private.Group("merchant")
		group.GET("findMerchant", apiMerchant.FindMerchant)        // 根据ID获取商户信息
		group.GET("getMerchantList", apiMerchant.GetMerchantList)  // 获取商户信息列表
	}
	{
	    group := public.Group("merchant")
	    group.GET("getMerchantPublic", apiMerchant.GetMerchantPublic)  // 商户信息开放接口
	}
}
