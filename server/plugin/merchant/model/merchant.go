package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"time"
)

// MerchantType 商户类型枚举
type MerchantType uint

// MerchantType枚举常量定义
const (
	MerchantTypeEnterprise MerchantType = 1 // 企业
	MerchantTypeIndividual MerchantType = 2 // 个体
)

// MerchantLevel 商户等级枚举
type MerchantLevel uint

// MerchantLevel枚举常量定义
const (
	MerchantLevelNormal MerchantLevel = 1 // 普通商户
	MerchantLevelPremium MerchantLevel = 2 // 高级商户
	MerchantLevelVIP MerchantLevel = 3 // VIP商户
)

// Merchant 商户信息 结构体
type Merchant struct {
	global.GVA_MODEL
	MerchantName      string      `json:"merchantName" form:"merchantName" gorm:"type:varchar(100);not null;comment:商户名称;column:merchant_name;size:100;index" binding:"required"` // 规则：必填，商户名称，长度1-100字符，创建普通索引以提升查询效率
	MerchantIcon      *string     `json:"merchantIcon" form:"merchantIcon" gorm:"type:varchar(255);comment:商户图标URL;column:merchant_icon;size:255;"`                           // 规则：可选，商户图标URL，长度不超过255字符
	ParentID          uint        `json:"parentID" form:"parentID" gorm:"comment:父商户ID;column:parent_id;index"`                                                       // 规则：必选，父商户ID，用于构建商户层级结构，建立索引优化查询
	MerchantType      MerchantType `json:"merchantType" form:"merchantType" gorm:"not null;comment:商户类型;column:merchant_type;" binding:"required"`                      // 规则：必填，商户类型枚举值（1-企业 2-个体）
	BusinessLicense   *string     `json:"businessLicense" form:"businessLicense" gorm:"type:varchar(100);comment:营业执照号;column:business_license;size:100;"`                // 规则：可选，营业执照号，长度不超过100字符
	LegalPerson       *string     `json:"legalPerson" form:"legalPerson" gorm:"type:varchar(50);comment:法人代表;column:legal_person;size:50;"`                              // 规则：可选，法人代表姓名，长度不超过50字符
	RegisteredAddress *string     `json:"registeredAddress" form:"registeredAddress" gorm:"type:varchar(255);comment:注册地址;column:registered_address;size:255;"`         // 规则：可选，注册地址，长度不超过255字符
	BusinessScope     *string     `json:"businessScope" form:"businessScope" gorm:"type:varchar(255);comment:经营范围;column:business_scope;size:255;"`                       // 规则：可选，经营范围，长度不超过255字符
	IsEnabled         bool        `json:"isEnabled" form:"isEnabled" gorm:"default:true;not null;comment:商户开关状态;column:is_enabled;index"`                              // 规则：默认为true（正常），商户开关状态，建立索引优化查询
	ValidStartTime    *time.Time  `json:"validStartTime" form:"validStartTime" gorm:"comment:有效开始时间;column:valid_start_time;"`                                        // 规则：可选，有效开始时间
	ValidEndTime      *time.Time  `json:"validEndTime" form:"validEndTime" gorm:"comment:有效结束时间;column:valid_end_time;"`                                              // 规则：可选，有效结束时间
	MerchantLevel     MerchantLevel `json:"merchantLevel" form:"merchantLevel" gorm:"not null;comment:商户等级;column:merchant_level;" binding:"required"`                 // 规则：必填，商户等级枚举值（1-普通商户 2-高级商户 3-VIP商户）
	Address           *string     `json:"address" form:"address" gorm:"type:varchar(255);comment:地址;column:address;size:255;"`                                             // 规则：可选，地址，长度不超过255字符
}

// TableName 商户信息 Merchant自定义表名 merchants
func (Merchant) TableName() string {
	return "merchants"
}
