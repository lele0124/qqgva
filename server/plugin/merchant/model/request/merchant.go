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
	IsEnabled      *bool       `json:"isEnabled" form:"isEnabled"` // 改为*bool类型，与前端保持一致
	MerchantType   *uint       `json:"merchantType" form:"merchantType"` // 改为*uint类型，与前端保持一致
	MerchantLevel  *uint       `json:"merchantLevel" form:"merchantLevel"` // 改为*uint类型，与前端保持一致
	request.PageInfo
	Sort  string `json:"sort" form:"sort"`
	Order string `json:"order" form:"order"`
}
