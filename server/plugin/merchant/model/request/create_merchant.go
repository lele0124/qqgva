package request

import (
	"time"
)

// CreateMerchantRequest 创建商户的请求模型
// 该模型使用string类型的时间字段，以避免类型转换错误

type CreateMerchantRequest struct {
	MerchantName      string `json:"merchantName" form:"merchantName" binding:"required"`
	MerchantIcon      string `json:"merchantIcon" form:"merchantIcon"`
	ParentID          uint   `json:"parentID" form:"parentID" binding:"required"`
	MerchantType      uint   `json:"merchantType" form:"merchantType" binding:"required"` // 与数据模型保持一致
	BusinessLicense   string `json:"businessLicense" form:"businessLicense"`
	LegalPerson       string `json:"legalPerson" form:"legalPerson"`
	RegisteredAddress string `json:"registeredAddress" form:"registeredAddress"`
	BusinessScope     string `json:"businessScope" form:"businessScope"`
	IsEnabled         bool   `json:"isEnabled" form:"isEnabled"` // 与数据模型保持一致
	ValidStartTime    string `json:"validStartTime" form:"validStartTime"` // 使用string类型接收时间
	ValidEndTime      string `json:"validEndTime" form:"validEndTime"`     // 使用string类型接收时间
	MerchantLevel     uint   `json:"merchantLevel" form:"merchantLevel" binding:"required"` // 与数据模型保持一致
	Address           string `json:"address" form:"address"`
}

// ToMerchantModel 将请求模型转换为数据模型
func (req *CreateMerchantRequest) ToMerchantModel() (model interface{}, err error) {
	// 直接使用请求模型中的类型，不需要转换
	result := map[string]interface{}{
		"MerchantName":      req.MerchantName,
		"MerchantIcon":      req.MerchantIcon,
		"ParentID":          req.ParentID,
		"MerchantType":      req.MerchantType,
		"BusinessLicense":   req.BusinessLicense,
		"LegalPerson":       req.LegalPerson,
		"RegisteredAddress": req.RegisteredAddress,
		"BusinessScope":     req.BusinessScope,
		"IsEnabled":         req.IsEnabled,
		"MerchantLevel":     req.MerchantLevel,
		"Address":           req.Address,
	}

	// 处理时间字段，只有非空时才尝试解析
	if req.ValidStartTime != "" {
		startTime, timeErr := time.Parse(time.RFC3339, req.ValidStartTime)
		if timeErr != nil {
			// 尝试其他常见的时间格式
			startTime, timeErr = time.Parse("2006-01-02 15:04:05", req.ValidStartTime)
			if timeErr != nil {
				return nil, timeErr
			}
		}
		result["ValidStartTime"] = &startTime
	}

	if req.ValidEndTime != "" {
		endTime, timeErr := time.Parse(time.RFC3339, req.ValidEndTime)
		if timeErr != nil {
			// 尝试其他常见的时间格式
			endTime, timeErr = time.Parse("2006-01-02 15:04:05", req.ValidEndTime)
			if timeErr != nil {
				return nil, timeErr
			}
		}
		result["ValidEndTime"] = &endTime
	}

	return result, nil
}