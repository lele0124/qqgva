package request

import (
	"time"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/merchant/model"
)

// UpdateMerchantRequest 更新商户的请求模型
// 该模型使用与数据模型一致的类型定义

type UpdateMerchantRequest struct {
	ID                uint                 `json:"id" form:"id" binding:"required"`
	MerchantName      *string              `json:"merchantName" form:"merchantName"`
	MerchantIcon      *string              `json:"merchantIcon" form:"merchantIcon"`
	ParentID          *uint                `json:"parentID" form:"parentID"`
	MerchantType      *model.MerchantType  `json:"merchantType" form:"merchantType"`          // 与数据模型保持一致
	BusinessLicense   *string              `json:"businessLicense" form:"businessLicense"`
	LegalPerson       *string              `json:"legalPerson" form:"legalPerson"`
	RegisteredAddress *string              `json:"registeredAddress" form:"registeredAddress"`
	BusinessScope     *string              `json:"businessScope" form:"businessScope"`
	IsEnabled         *bool                `json:"isEnabled" form:"isEnabled"`               // 与数据模型保持一致
	ValidStartTime    *time.Time           `json:"validStartTime" form:"validStartTime"`      // 使用*time.Time类型接收时间
	ValidEndTime      *time.Time           `json:"validEndTime" form:"validEndTime"`          // 使用*time.Time类型接收时间
	MerchantLevel     *model.MerchantLevel `json:"merchantLevel" form:"merchantLevel"`        // 与数据模型保持一致
	Address           *string              `json:"address" form:"address"`
}

// ToMerchantModel 将请求模型转换为数据模型
func (req *UpdateMerchantRequest) ToMerchantModel() (model.Merchant, error) {
	// 构造Merchant模型
	merchant := model.Merchant{
		// ID将在GORM更新时使用
	}

	// 处理可选字段
	if req.MerchantName != nil {
		merchant.MerchantName = *req.MerchantName
	}
	if req.MerchantIcon != nil {
		merchant.MerchantIcon = req.MerchantIcon
	}
	if req.ParentID != nil {
		merchant.ParentID = *req.ParentID
	}
	if req.MerchantType != nil {
		merchant.MerchantType = *req.MerchantType
	}
	if req.BusinessLicense != nil {
		merchant.BusinessLicense = req.BusinessLicense
	}
	if req.LegalPerson != nil {
		merchant.LegalPerson = req.LegalPerson
	}
	if req.RegisteredAddress != nil {
		merchant.RegisteredAddress = req.RegisteredAddress
	}
	if req.BusinessScope != nil {
		merchant.BusinessScope = req.BusinessScope
	}
	if req.IsEnabled != nil {
		merchant.IsEnabled = *req.IsEnabled
	}
	if req.ValidStartTime != nil {
		merchant.ValidStartTime = req.ValidStartTime
	}
	if req.ValidEndTime != nil {
		merchant.ValidEndTime = req.ValidEndTime
	}
	if req.MerchantLevel != nil {
		merchant.MerchantLevel = *req.MerchantLevel
	}
	if req.Address != nil {
		merchant.Address = req.Address
	}

	return merchant, nil
}