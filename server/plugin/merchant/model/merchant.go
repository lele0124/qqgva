
package model
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Merchant 商户信息 结构体
type Merchant struct {
    global.GVA_MODEL
  MerchantName  *string `json:"merchantName" form:"merchantName" gorm:"comment:商户名称;column:merchant_name;size:100;" binding:"required"`  //商户名称
  ContactPerson  *string `json:"contactPerson" form:"contactPerson" gorm:"comment:联系人;column:contact_person;size:50;" binding:"required"`  //联系人
  ContactPhone  *string `json:"contactPhone" form:"contactPhone" gorm:"comment:联系电话;column:contact_phone;size:20;" binding:"required"`  //联系电话
  Address  *string `json:"address" form:"address" gorm:"comment:商户地址;column:address;size:255;"`  //商户地址
  BusinessScope  *string `json:"businessScope" form:"businessScope" gorm:"comment:经营范围;column:business_scope;size:255;"`  //经营范围
  IsEnabled  *bool `json:"isEnabled" form:"isEnabled" gorm:"default:true;comment:是否启用;column:is_enabled;"`  //是否启用
}


// TableName 商户信息 Merchant自定义表名 merchants
func (Merchant) TableName() string {
    return "merchants"
}







