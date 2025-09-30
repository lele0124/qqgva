package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/google/uuid"
	"time"
)

// Merchant 商户信息 结构体
type Merchant struct {
	global.GVA_MODEL
	UUID              uuid.UUID  `json:"uuid" form:"uuid" gorm:"comment:唯一标识;column:uuid;type:uuid;"`                                                 //唯一标识
	MerchantName      *string    `json:"merchantName" form:"merchantName" gorm:"comment:商户名称;column:merchant_name;size:100;index" binding:"required"` //商户名称
	MerchantIcon      *string    `json:"merchantIcon" form:"merchantIcon" gorm:"comment:商户图标URL;column:merchant_icon;size:255;"`                      //商户图标URL
	ParentID          *uint      `json:"parentID" form:"parentID" gorm:"comment:父商户ID;column:parent_id;index"`                                        //父商户ID
	MerchantType      *uint      `json:"merchantType" form:"merchantType" gorm:"comment:商户类型;column:merchant_type;" binding:"required"`               //商户类型：1-企业 2-个体
	BusinessLicense   *string    `json:"businessLicense" form:"businessLicense" gorm:"comment:营业执照号;column:business_license;size:100;"`               //营业执照号
	LegalPerson       *string    `json:"legalPerson" form:"legalPerson" gorm:"comment:法人代表;column:legal_person;size:50;"`                             //法人代表
	RegisteredAddress *string    `json:"registeredAddress" form:"registeredAddress" gorm:"comment:注册地址;column:registered_address;size:255;"`          //注册地址
	BusinessScope     *string    `json:"businessScope" form:"businessScope" gorm:"comment:经营范围;column:business_scope;size:255;"`                      //经营范围
	IsEnabled         bool       `json:"isEnabled" form:"isEnabled" gorm:"default:true;comment:商户开关状态;column:is_enabled;index"`                       //商户开关状态：true-正常 false-关闭
	ValidStartTime    *time.Time `json:"validStartTime" form:"validStartTime" gorm:"comment:有效开始时间;column:valid_start_time;"`                         //有效开始时间
	ValidEndTime      *time.Time `json:"validEndTime" form:"validEndTime" gorm:"comment:有效结束时间;column:valid_end_time;"`                               //有效结束时间
	MerchantLevel     *uint      `json:"merchantLevel" form:"merchantLevel" gorm:"comment:商户等级;column:merchant_level;" binding:"required"`            //商户等级：1-普通商户 2-高级商户 3-VIP商户
}

// TableName 商户信息 Merchant自定义表名 merchants
func (Merchant) TableName() string {
	return "merchants"
}
