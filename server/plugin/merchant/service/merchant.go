package service

import (
	"context"
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/merchant/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/merchant/model/request"
	"github.com/google/uuid"
)

var Merchant = new(merchant)

type merchant struct{}

// CreateMerchant 创建商户信息记录
// Author [yourname](https://github.com/yourname)
func (s *merchant) CreateMerchant(ctx context.Context, merchant *model.Merchant) (err error) {
	// 验证商户名称不能为空
	if merchant.MerchantName == nil || *merchant.MerchantName == "" {
		return errors.New("商户名称不能为空")
	}

	// 验证商户类型不能为空
	if merchant.MerchantType == nil || *merchant.MerchantType == 0 {
		return errors.New("商户类型不能为空")
	}

	// IsEnabled已经是普通bool类型，不需要再检查nil值

	// 验证商户等级不能为空
	if merchant.MerchantLevel == nil || *merchant.MerchantLevel == 0 || *merchant.MerchantLevel > 3 {
		return errors.New("商户等级必须为1(普通)、2(高级)或3(VIP)")
	}

	// 生成UUID
	merchant.UUID = uuid.New()
	err = global.GVA_DB.Create(merchant).Error
	return err
}

// DeleteMerchant 删除商户信息记录
// Author [yourname](https://github.com/yourname)
func (s *merchant) DeleteMerchant(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&model.Merchant{}, "id = ?", ID).Error
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
	// 检查UUID是否为空或全零UUID，如果是则生成新的UUID
	if merchant.UUID == uuid.Nil || merchant.UUID.String() == "00000000-0000-0000-0000-000000000000" {
		merchant.UUID = uuid.New()
	}

	// 验证商户名称不能为空
	if merchant.MerchantName == nil || *merchant.MerchantName == "" {
		return errors.New("商户名称不能为空")
	}

	// 验证商户类型不能为空
	if merchant.MerchantType == nil || *merchant.MerchantType == 0 {
		return errors.New("商户类型不能为空")
	}

	// IsEnabled已经是普通bool类型，不需要再检查nil值

	// 验证商户等级不能为空
	if merchant.MerchantLevel == nil || *merchant.MerchantLevel == 0 || *merchant.MerchantLevel > 3 {
		return errors.New("商户等级必须为1(普通)、2(高级)或3(VIP)")
	}

	err = global.GVA_DB.Model(&model.Merchant{}).Where("id = ?", merchant.ID).Updates(&merchant).Error
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
		db = db.Where("merchant_name LIKE ?", "%"+*info.MerchantName+"%")
	}
	if info.Address != nil && *info.Address != "" {
		db = db.Where("address LIKE ?", "%"+*info.Address+"%")
	}
	if info.BusinessScope != nil && *info.BusinessScope != "" {
		db = db.Where("business_scope LIKE ?", "%"+*info.BusinessScope+"%")
	}
	// 只有当IsEnabled不为nil且不为空字符串时才应用过滤条件
	// 直接使用布尔类型进行过滤
	if info.IsEnabled != nil {
		db = db.Where("is_enabled = ?", *info.IsEnabled)
	}

	// 直接使用uint类型进行过滤，不需要类型转换
	if info.MerchantType != nil {
		db = db.Where("merchant_type = ?", *info.MerchantType)
	}
	if info.MerchantLevel != nil {
		db = db.Where("merchant_level = ?", *info.MerchantLevel)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	var OrderStr string
	orderMap := make(map[string]bool)
	orderMap["id"] = true
	orderMap["created_at"] = true
	orderMap["updated_at"] = true
	orderMap["merchant_name"] = true
	// 默认按更新时间倒序排序
	if info.Sort == "" {
		OrderStr = "updated_at desc"
	} else if orderMap[info.Sort] {
		OrderStr = info.Sort
		if info.Order == "descending" || info.Order == "desc" {
			OrderStr = OrderStr + " desc"
		}
	}

	if OrderStr != "" {
		db = db.Order(OrderStr)
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Find(&merchants).Error
	return merchants, total, err
}

func (s *merchant) GetMerchantPublic(ctx context.Context) {

}
