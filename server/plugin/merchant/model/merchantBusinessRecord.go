
package model
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"time"
)

// MerchantBusinessRecord 商户业务记录 结构体
type MerchantBusinessRecord struct {
    global.GVA_MODEL
  MerchantID  *string `json:"merchantId" form:"merchantId" gorm:"comment:商户ID;column:merchant_id;size:36;" binding:"required"`  //商户ID
  RecordType  *string `json:"recordType" form:"recordType" gorm:"comment:记录类型;column:record_type;size:50;" binding:"required"`  //记录类型
  Amount  *float64 `json:"amount" form:"amount" gorm:"default:0;comment:金额;column:amount;" binding:"required"`  //金额
  Description  *string `json:"description" form:"description" gorm:"comment:描述;column:description;size:255;"`  //描述
  RecordTime  *time.Time `json:"recordTime" form:"recordTime" gorm:"comment:记录时间;column:record_time;" binding:"required"`  //记录时间
    CreatedBy  uint   `gorm:"column:created_by;comment:创建者"`
    UpdatedBy  uint   `gorm:"column:updated_by;comment:更新者"`
    DeletedBy  uint   `gorm:"column:deleted_by;comment:删除者"`
}


// TableName 商户业务记录 MerchantBusinessRecord自定义表名 merchant_business_records
func (MerchantBusinessRecord) TableName() string {
    return "merchant_business_records"
}







