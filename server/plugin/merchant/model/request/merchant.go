package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/merchant/model"
)

// MerchantSearch 商户搜索请求模型
type MerchantSearch struct {
	CreatedAtRange []time.Time          `json:"createdAtRange" form:"createdAtRange[]"`     // 创建时间范围
	MerchantName   *string              `json:"merchantName" form:"merchantName"`           // 商户名称（模糊搜索）
	Address        *string              `json:"address" form:"address"`                     // 地址（模糊搜索）
	BusinessScope  *string              `json:"businessScope" form:"businessScope"`         // 经营范围（模糊搜索）
	IsEnabled      *bool                `json:"isEnabled" form:"isEnabled"`                 // 商户开关状态
	MerchantType   *model.MerchantType  `json:"merchantType" form:"merchantType"`           // 商户类型（枚举值）
	MerchantLevel  *model.MerchantLevel `json:"merchantLevel" form:"merchantLevel"`         // 商户等级（枚举值）
	request.PageInfo                                                             // 分页信息
	Sort  string `json:"sort" form:"sort"`                                     // 排序字段
	Order string `json:"order" form:"order"`                                   // 排序方式（asc/desc）
}