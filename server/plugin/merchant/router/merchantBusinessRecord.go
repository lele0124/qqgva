package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

var MerchantBusinessRecord = new(merchant_record)

type merchant_record struct {}

// Init 初始化 商户业务记录 路由信息
func (r *merchant_record) Init(public *gin.RouterGroup, private *gin.RouterGroup) {
	{
	    group := private.Group("merchant_record").Use(middleware.OperationRecord())
		group.POST("createMerchantBusinessRecord", apiMerchantBusinessRecord.CreateMerchantBusinessRecord)   // 新建商户业务记录
		group.DELETE("deleteMerchantBusinessRecord", apiMerchantBusinessRecord.DeleteMerchantBusinessRecord) // 删除商户业务记录
		group.DELETE("deleteMerchantBusinessRecordByIds", apiMerchantBusinessRecord.DeleteMerchantBusinessRecordByIds) // 批量删除商户业务记录
		group.PUT("updateMerchantBusinessRecord", apiMerchantBusinessRecord.UpdateMerchantBusinessRecord)    // 更新商户业务记录
	}
	{
	    group := private.Group("merchant_record")
		group.GET("findMerchantBusinessRecord", apiMerchantBusinessRecord.FindMerchantBusinessRecord)        // 根据ID获取商户业务记录
		group.GET("getMerchantBusinessRecordList", apiMerchantBusinessRecord.GetMerchantBusinessRecordList)  // 获取商户业务记录列表
	}
	{
	    group := public.Group("merchant_record")
	    group.GET("getMerchantBusinessRecordDataSource", apiMerchantBusinessRecord.GetMerchantBusinessRecordDataSource)  // 获取商户业务记录数据源
	    group.GET("getMerchantBusinessRecordPublic", apiMerchantBusinessRecord.GetMerchantBusinessRecordPublic)  // 商户业务记录开放接口
	}
}
