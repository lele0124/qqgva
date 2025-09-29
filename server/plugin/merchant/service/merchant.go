
package service

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/merchant/model"
    "github.com/flipped-aurora/gin-vue-admin/server/plugin/merchant/model/request"
)

var Merchant = new(merchant)

type merchant struct {}
// CreateMerchant 创建商户信息记录
// Author [yourname](https://github.com/yourname)
func (s *merchant) CreateMerchant(ctx context.Context, merchant *model.Merchant) (err error) {
	err = global.GVA_DB.Create(merchant).Error
	return err
}

// DeleteMerchant 删除商户信息记录
// Author [yourname](https://github.com/yourname)
func (s *merchant) DeleteMerchant(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&model.Merchant{},"id = ?",ID).Error
	return err
}

// DeleteMerchantByIds 批量删除商户信息记录
// Author [yourname](https://github.com/yourname)
func (s *merchant) DeleteMerchantByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Where("id in ?", IDs).Delete(&model.Merchant{}).Error
	return err
}

// UpdateMerchant 更新商户信息记录
// Author [yourname](https://github.com/yourname)
func (s *merchant) UpdateMerchant(ctx context.Context, merchant model.Merchant) (err error) {
	err = global.GVA_DB.Model(&model.Merchant{}).Where("id = ?",merchant.ID).Updates(&merchant).Error
	return err
}

// GetMerchant 根据ID获取商户信息记录
// Author [yourname](https://github.com/yourname)
func (s *merchant) GetMerchant(ctx context.Context, ID string) (merchant model.Merchant, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&merchant).Error
	return
}
// GetMerchantInfoList 分页获取商户信息记录
// Author [yourname](https://github.com/yourname)
func (s *merchant) GetMerchantInfoList(ctx context.Context, info request.MerchantSearch) (list []model.Merchant, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&model.Merchant{})
    var merchants []model.Merchant
    // 如果有条件搜索 下方会自动创建搜索语句
    if len(info.CreatedAtRange) == 2 {
     db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
    }
  
    if info.MerchantName != nil && *info.MerchantName != "" {
        db = db.Where("merchant_name LIKE ?", "%"+ *info.MerchantName+"%")
    }
    if info.ContactPerson != nil && *info.ContactPerson != "" {
        db = db.Where("contact_person LIKE ?", "%"+ *info.ContactPerson+"%")
    }
    if info.ContactPhone != nil && *info.ContactPhone != "" {
        db = db.Where("contact_phone LIKE ?", "%"+ *info.ContactPhone+"%")
    }
    if info.Address != nil && *info.Address != "" {
        db = db.Where("address LIKE ?", "%"+ *info.Address+"%")
    }
    if info.BusinessScope != nil && *info.BusinessScope != "" {
        db = db.Where("business_scope LIKE ?", "%"+ *info.BusinessScope+"%")
    }
    if info.IsEnabled != nil {
        db = db.Where("is_enabled = ?", *info.IsEnabled)
    }
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }
        var OrderStr string
        orderMap := make(map[string]bool)
        orderMap["id"] = true
        orderMap["created_at"] = true
        orderMap["merchant_name"] = true
       if orderMap[info.Sort] {
          OrderStr = info.Sort
          if info.Order == "descending" {
             OrderStr = OrderStr + " desc"
          }
          db = db.Order(OrderStr)
       }

	if limit != 0 {
       db = db.Limit(limit).Offset(offset)
    }
	err = db.Find(&merchants).Error
	return  merchants, total, err
}

func (s *merchant)GetMerchantPublic(ctx context.Context) {

}
