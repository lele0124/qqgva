package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type MerchantSearch struct {
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	MerchantName   *string     `json:"merchantName" form:"merchantName"`
	Address        *string     `json:"address" form:"address"`
	BusinessScope  *string     `json:"businessScope" form:"businessScope"`
	IsEnabled      *string     `json:"isEnabled" form:"isEnabled"` // 修改为string类型接收前端数据
	MerchantType   *string     `json:"merchantType" form:"merchantType"` // 修改为string类型接收前端数据
	MerchantLevel  *string     `json:"merchantLevel" form:"merchantLevel"` // 修改为string类型接收前端数据
	request.PageInfo
	Sort  string `json:"sort" form:"sort"`
	Order string `json:"order" form:"order"`
}
