
package service

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/merchant/model"
    "github.com/flipped-aurora/gin-vue-admin/server/plugin/merchant/model/request"
    "gorm.io/gorm"
)

var MerchantBusinessRecord = new(merchant_record)

type merchant_record struct {}
// CreateMerchantBusinessRecord 创建商户业务记录记录
// Author [yourname](https://github.com/yourname)
func (s *merchant_record) CreateMerchantBusinessRecord(ctx context.Context, merchant_record *model.MerchantBusinessRecord) (err error) {
	err = global.GVA_DB.Create(merchant_record).Error
	return err
}

// DeleteMerchantBusinessRecord 删除商户业务记录记录
// Author [yourname](https://github.com/yourname)
func (s *merchant_record) DeleteMerchantBusinessRecord(ctx context.Context, ID string,userID uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
	    if err := tx.Model(&model.MerchantBusinessRecord{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
              return err
        }
        if err = tx.Delete(&model.MerchantBusinessRecord{},"id = ?",ID).Error; err != nil {
              return err
        }
        return nil
	})
	return err
}

// DeleteMerchantBusinessRecordByIds 批量删除商户业务记录记录
// Author [yourname](https://github.com/yourname)
func (s *merchant_record) DeleteMerchantBusinessRecordByIds(ctx context.Context, IDs []string,deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
	    if err := tx.Model(&model.MerchantBusinessRecord{}).Where("id in ?", IDs).Update("deleted_by", deleted_by).Error; err != nil {
            return err
        }
        if err := tx.Where("id in ?", IDs).Delete(&model.MerchantBusinessRecord{}).Error; err != nil {
            return err
        }
        return nil
    })
	return err
}

// UpdateMerchantBusinessRecord 更新商户业务记录记录
// Author [yourname](https://github.com/yourname)
func (s *merchant_record) UpdateMerchantBusinessRecord(ctx context.Context, merchant_record model.MerchantBusinessRecord) (err error) {
	err = global.GVA_DB.Model(&model.MerchantBusinessRecord{}).Where("id = ?",merchant_record.ID).Updates(&merchant_record).Error
	return err
}

// GetMerchantBusinessRecord 根据ID获取商户业务记录记录
// Author [yourname](https://github.com/yourname)
func (s *merchant_record) GetMerchantBusinessRecord(ctx context.Context, ID string) (merchant_record model.MerchantBusinessRecord, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&merchant_record).Error
	return
}
// GetMerchantBusinessRecordInfoList 分页获取商户业务记录记录
// Author [yourname](https://github.com/yourname)
func (s *merchant_record) GetMerchantBusinessRecordInfoList(ctx context.Context, info request.MerchantBusinessRecordSearch) (list []model.MerchantBusinessRecord, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&model.MerchantBusinessRecord{})
    var merchant_records []model.MerchantBusinessRecord
    // 如果有条件搜索 下方会自动创建搜索语句
    if len(info.CreatedAtRange) == 2 {
     db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
    }
  
    if info.MerchantID != nil && *info.MerchantID != "" {
        db = db.Where("merchant_id = ?", *info.MerchantID)
    }
    if info.RecordType != nil && *info.RecordType != "" {
        db = db.Where("record_type LIKE ?", "%"+ *info.RecordType+"%")
    }
    if info.Amount != nil {
        db = db.Where("amount >= ?", *info.Amount)
    }
    if info.Description != nil && *info.Description != "" {
        db = db.Where("description LIKE ?", "%"+ *info.Description+"%")
    }
			if len(info.RecordTimeRange) == 2 {
				db = db.Where("record_time BETWEEN ? AND ? ", info.RecordTimeRange[0], info.RecordTimeRange[1])
			}
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }

	if limit != 0 {
       db = db.Limit(limit).Offset(offset)
    }
	err = db.Find(&merchant_records).Error
	return  merchant_records, total, err
}
func (s *merchant_record)GetMerchantBusinessRecordDataSource(ctx context.Context) (res map[string][]map[string]any, err error) {
	res = make(map[string][]map[string]any)
	
	   merchantId := make([]map[string]any, 0)
	   global.GVA_DB.Table("merchants").Where("deleted_at IS NULL").Select("merchant_name as label,id as value").Scan(&merchantId)
	   res["merchantId"] = merchantId
	return
}

func (s *merchant_record)GetMerchantBusinessRecordPublic(ctx context.Context) {

}
