package global

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GVA_MODEL struct {
	ID                   uint           `gorm:"primarykey" json:"ID"`                              // 主键ID
	UUID                 uuid.UUID      `json:"uuid" gorm:"index;comment:全局UUID"`                  // 全局UUID
	OperatorId           uint           `json:"operatorId" gorm:"index;comment:操作人ID"`             // 操作人ID
	OperatorName         string         `json:"operatorName" gorm:"index;comment:操作人姓名"`           // 操作人姓名
	OperatorMerchantId   uint           `json:"operatorMerchantId" gorm:"index;comment:操作人商户ID"`   // 操作人商户ID
	OperatorMerchantName string         `json:"operatorMerchantName" gorm:"index;comment:操作人商户名称"` // 操作人商户名称
	UpdatedAt            time.Time      `json:"updatedAt" gorm:"index;comment:更新时间"`               // 更新时间
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"deletedAt"`                            // 删除时间
}
