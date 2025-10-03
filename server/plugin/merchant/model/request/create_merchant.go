package request

import (
	"time"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/merchant/model"
)

// CreateMerchantRequest 创建商户的请求模型
// 该模型使用与数据模型一致的类型定义

type CreateMerchantRequest struct {
	MerchantName      string               `json:"merchantName" form:"merchantName" binding:"required"`
	MerchantIcon      string               `json:"merchantIcon" form:"merchantIcon"`
	ParentID          uint                 `json:"parentID" form:"parentID" binding:"required"`
	MerchantType      model.MerchantType   `json:"merchantType" form:"merchantType" binding:"required"`      // 与数据模型保持一致
	BusinessLicense   string               `json:"businessLicense" form:"businessLicense"`
	LegalPerson       string               `json:"legalPerson" form:"legalPerson"`
	RegisteredAddress string               `json:"registeredAddress" form:"registeredAddress"`
	BusinessScope     string               `json:"businessScope" form:"businessScope"`
	IsEnabled         bool                 `json:"isEnabled" form:"isEnabled"`                               // 与数据模型保持一致
	ValidStartTime    *time.Time           `json:"validStartTime" form:"validStartTime"`                     // 使用*time.Time类型接收时间
	ValidEndTime      *time.Time           `json:"validEndTime" form:"validEndTime"`                         // 使用*time.Time类型接收时间
	MerchantLevel     model.MerchantLevel  `json:"merchantLevel" form:"merchantLevel" binding:"required"`    // 与数据模型保持一致
	Address           string               `json:"address" form:"address"`
}

// ToMerchantModel 将请求模型转换为数据模型
func (req *CreateMerchantRequest) ToMerchantModel() (model.Merchant, error) {
	// 直接构造Merchant模型
	merchant := model.Merchant{
		MerchantName:      req.MerchantName,
		ParentID:          req.ParentID,
		MerchantType:      req.MerchantType,
		BusinessLicense:   &req.BusinessLicense,
		LegalPerson:       &req.LegalPerson,
		RegisteredAddress: &req.RegisteredAddress,
		BusinessScope:     &req.BusinessScope,
		IsEnabled:         req.IsEnabled,
		MerchantLevel:     req.MerchantLevel,
		Address:           &req.Address,
	}

	// 处理可选字段
	if req.MerchantIcon != "" {
		merchant.MerchantIcon = &req.MerchantIcon
	}

	// 处理时间字段
	merchant.ValidStartTime = req.ValidStartTime
	merchant.ValidEndTime = req.ValidEndTime

	return merchant, nil
}
