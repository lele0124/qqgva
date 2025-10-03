# Merchant 模块数据结构文档

## 文档信息
- **最后更新时间**: 2025-10-04 15:30:00
- **版本**: v2.2
- **维护者**: 智能体系统

## 任务状态概览
- 文件检查与定位: **已完成** (2025-10-04 15:30:00)
- 内容读取与更新: **已完成** (2025-10-04 15:30:15)
- 代码详细检查与数据库验证: **已完成** (2025-10-04 15:31:00)
- 优化建议生成: **已完成** (2025-10-04 15:32:00)
- 结构体重构示例: **已完成** (2025-10-04 15:32:30)

## 结构体定义
```go
package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"time"
)

// Merchant 商户信息 结构体
type Merchant struct {
	global.GVA_MODEL
	MerchantName      *string    `json:"merchantName" form:"merchantName" gorm:"comment:商户名称;column:merchant_name;size:100;index" binding:"required"` // 规则：必填，商户名称，长度1-100字符，创建普通索引以提升查询效率
	MerchantIcon      *string    `json:"merchantIcon" form:"merchantIcon" gorm:"comment:商户图标URL;column:merchant_icon;size:255;"`                      // 规则：可选，商户图标URL，长度不超过255字符
	ParentID          *uint      `json:"parentID" form:"parentID" gorm:"comment:父商户ID;column:parent_id;index"`                                        // 规则：可选，父商户ID，用于构建商户层级结构，建立索引优化查询
	MerchantType      *uint      `json:"merchantType" form:"merchantType" gorm:"comment:商户类型;column:merchant_type;" binding:"required"`               // 规则：必填，商户类型枚举值（1-企业 2-个体）
	BusinessLicense   *string    `json:"businessLicense" form:"businessLicense" gorm:"comment:营业执照号;column:business_license;size:100;"`               // 规则：可选，营业执照号，长度不超过100字符
	LegalPerson       *string    `json:"legalPerson" form:"legalPerson" gorm:"comment:法人代表;column:legal_person;size:50;"`                             // 规则：可选，法人代表姓名，长度不超过50字符
	RegisteredAddress *string    `json:"registeredAddress" form:"registeredAddress" gorm:"comment:注册地址;column:registered_address;size:255;"`          // 规则：可选，注册地址，长度不超过255字符
	BusinessScope     *string    `json:"businessScope" form:"businessScope" gorm:"comment:经营范围;column:business_scope;size:255;"`                      // 规则：可选，经营范围，长度不超过255字符
	IsEnabled         bool       `json:"isEnabled" form:"isEnabled" gorm:"default:true;comment:商户开关状态;column:is_enabled;index"`                       // 规则：默认为true（正常），商户开关状态，建立索引优化查询
	ValidStartTime    *time.Time `json:"validStartTime" form:"validStartTime" gorm:"comment:有效开始时间;column:valid_start_time;"`                         // 规则：可选，有效开始时间
	ValidEndTime      *time.Time `json:"validEndTime" form:"validEndTime" gorm:"comment:有效结束时间;column:valid_end_time;"`                               // 规则：可选，有效结束时间
	MerchantLevel     *uint      `json:"merchantLevel" form:"merchantLevel" gorm:"comment:商户等级;column:merchant_level;" binding:"required"`            // 规则：必填，商户等级枚举值（1-普通商户 2-高级商户 3-VIP商户）
}

// TableName 商户信息 Merchant自定义表名 merchants
func (Merchant) TableName() string {
	return "merchants"
}

// Validate 验证商户字段是否符合业务规则
func (m *Merchant) Validate() error {
	// 验证逻辑将在优化建议中实现
	return nil
}
```

## 字段业务逻辑详解
- **GVA_MODEL**: 全局模型，包含ID、CreatedAt、UpdatedAt、DeletedAt等基础字段。
- **MerchantName**: 必填字段，存储商户名称。业务规则要求长度在1-100字符之间，并在数据库中创建了索引以优化按名称查询的性能。
- **MerchantIcon**: 可选字段，存储商户图标URL地址，最大长度为255字符。
- **ParentID**: 可选字段，表示父商户ID，用于构建商户的层级结构，建立索引优化层级查询效率。
- **MerchantType**: 必填字段，表示商户类型，使用枚举值（1-企业 2-个体）。
- **BusinessLicense**: 可选字段，存储营业执照号码，最大长度为100字符。
- **LegalPerson**: 可选字段，存储法人代表姓名，最大长度为50字符。
- **RegisteredAddress**: 可选字段，存储注册地址，最大长度为255字符。
- **BusinessScope**: 可选字段，存储经营范围描述，最大长度为255字符。
- **IsEnabled**: 必填字段，默认为true，表示商户的启用状态，建立索引优化状态筛选查询。
- **ValidStartTime**: 可选字段，表示商户有效期的开始时间。
- **ValidEndTime**: 可选字段，表示商户有效期的结束时间。
- **MerchantLevel**: 必填字段，表示商户等级，使用枚举值（1-普通商户 2-高级商户 3-VIP商户）。

## 实现状况
### 检查时间: 2025-10-04 15:31:00
- ✅ 后端Go定义完整
- ⚠️ 数据库字段映射部分不一致
- ❌ 前端类型定义待检查
- ❌ 数据库迁移脚本待检查
- ❌ API接口定义待检查

### 数据库验证结果
通过MCP工具查询，`merchants`表结构信息如下（2025-10-04 17:28:45更新）：

**关于global.GVA_MODEL的说明和系统提示词优化建议：**

`global.GVA_MODEL`是Go语言结构体中的嵌入类型(embedding)，是一个基础模型，在GORM框架中会被展开为具体字段映射到数据库表中，而不是作为一个单独的字段存在。

### global.GVA_MODEL的具体定义
```go
// GVA_MODEL 基础模型结构体定义（位于server/global/model.go）
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
```

### 系统提示词优化建议
为了让智能体能够通过MCP正确查询到global.GVA_MODEL在数据库对应的字段，建议在系统提示词中加入以下内容：

1. **明确说明嵌入结构体特性**：
   > `global.GVA_MODEL`是Go语言的结构体嵌入类型(embedding)，在GORM框架中，嵌入的结构体会被自动展开，其字段会直接作为主结构体的字段映射到数据库表中，而不是以嵌套结构体的形式存在。

2. **详细列出字段映射关系**：
   > `global.GVA_MODEL`包含以下字段，这些字段会被映射到数据库表中：
   > - ID (uint) -> 数据库字段: `id` (bigint, 主键)
   > - UUID (uuid.UUID) -> 数据库字段: `uuid` (text)
   > - OperatorId (uint) -> 数据库字段: `operator_id` (bigint)
   > - OperatorName (string) -> 数据库字段: `operator_name` (text)
   > - OperatorMerchantId (uint) -> 数据库字段: `operator_merchant_id` (bigint)
   > - OperatorMerchantName (string) -> 数据库字段: `operator_merchant_name` (text)
   > - UpdatedAt (time.Time) -> 数据库字段: `updated_at` (timestamp with time zone)
   > - DeletedAt (gorm.DeletedAt) -> 数据库字段: `deleted_at` (timestamp with time zone，用于软删除)

3. **提供MCP查询示例**：
   > 当需要查询包含global.GVA_MODEL的表结构时，应当使用以下SQL查询，它会返回所有字段（包括来自GVA_MODEL的字段）：
   > ```sql
   > SELECT column_name, data_type, character_maximum_length, is_nullable, column_default 
   > FROM information_schema.columns 
   > WHERE table_name = '目标表名' 
   > ORDER BY ordinal_position;
   > ```

4. **添加验证方法**：
   > 在代码生成或验证阶段，确保生成的结构体正确嵌入了global.GVA_MODEL，并验证其对应的数据库字段是否完整存在于表结构中。

通过这些优化，智能体将能够更准确地理解和查询global.GVA_MODEL在数据库中对应的字段信息。

### 数据库字段映射关系总结
这些字段会被自动映射到数据库表中，但在MCP数据库查询结果中，是以独立字段形式显示的（如`id`、`uuid`、`updated_at`等），而不是以`global.GVA_MODEL`整体形式出现。

### merchants表完整字段信息
- `id`: bigint, 主键，自增，非空
- `uuid`: text, 可选
- `operator_id`: bigint, 可选
- `operator_name`: text, 可选
- `operator_merchant_id`: bigint, 可选
- `operator_merchant_name`: text, 可选
- `updated_at`: timestamp with time zone, 可选
- `deleted_at`: timestamp with time zone, 可选
- `merchant_name`: character varying(100), 可选
- `contact_person`: character varying(50), 可选
- `contact_phone`: character varying(20), 可选
- `address`: character varying(255), 可选
- `business_scope`: character varying(255), 可选
- `is_enabled`: boolean, 可选，默认true
- `merchant_icon`: character varying(255), 可选
- `parent_id`: bigint, 可选
- `merchant_type`: bigint, 可选
- `business_license`: character varying(100), 可选
- `legal_person`: character varying(50), 可选
- `registered_address`: character varying(255), 可选
- `valid_start_time`: timestamp with time zone, 可选
- `valid_end_time`: timestamp with time zone, 可选
- `merchant_level`: bigint, 可选

**不一致点**:
1. 数据库中存在结构体未定义的字段：`contact_person`, `contact_phone`, `address`
2. 结构体中标记为必填的字段在数据库中均为可选
3. 结构体中的`*uint`类型与数据库中的`bigint`类型需要统一

### 前端与API检查结果
待执行

## 优化建议
### 生成时间: 2025-10-04 15:32:00
1. **高优先级**: 同步数据库和结构体字段，添加缺失的`ContactPerson`, `ContactPhone`, `Address`字段
2. **高优先级**: 在数据库层面为必填字段添加`NOT NULL`约束，确保数据完整性
3. **中优先级**: 完善Validate方法，添加完整的字段业务规则验证逻辑
4. **中优先级**: 统一类型映射，将结构体中的`*uint`类型与数据库中的`bigint`类型保持一致
5. **低优先级**: 为枚举类型字段添加常量定义，提高代码可读性和可维护性
6. **低优先级**: 考虑添加版本控制字段，便于数据变更追踪

## 重构建议
```go
package model

import (
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"regexp"
	"time"
)

// 商户类型枚举常量
type MerchantType uint

const (
	MerchantTypeEnterprise MerchantType = 1 // 企业商户
	MerchantTypeIndividual MerchantType = 2 // 个体商户
)

// 商户等级枚举常量
type MerchantLevel uint

const (
	MerchantLevelNormal MerchantLevel = 1 // 普通商户
	MerchantLevelAdvanced MerchantLevel = 2 // 高级商户
	MerchantLevelVIP MerchantLevel = 3 // VIP商户
)

// Merchant 商户信息 结构体
type Merchant struct {
	global.GVA_MODEL
	MerchantName      string     `json:"merchantName" form:"merchantName" gorm:"type:varchar(100);not null;comment:商户名称;column:merchant_name;size:100;index" binding:"required"` // 规则：必填，商户名称，长度1-100字符，创建普通索引以提升查询效率
	MerchantIcon      *string    `json:"merchantIcon" form:"merchantIcon" gorm:"type:varchar(255);comment:商户图标URL;column:merchant_icon;size:255;"`                           // 规则：可选，商户图标URL，长度不超过255字符
	ParentID          uint      `json:"parentID" form:"parentID" gorm:"comment:父商户ID;column:parent_id;index"`                                                       // 规则：必选，父商户ID，用于构建商户层级结构，建立索引优化查询
	MerchantType      MerchantType `json:"merchantType" form:"merchantType" gorm:"not null;comment:商户类型;column:merchant_type;" binding:"required"`                      // 规则：必填，商户类型枚举值（1-企业 2-个体）
	BusinessLicense   *string    `json:"businessLicense" form:"businessLicense" gorm:"type:varchar(100);comment:营业执照号;column:business_license;size:100;"`                // 规则：可选，营业执照号，长度不超过100字符
	LegalPerson       *string    `json:"legalPerson" form:"legalPerson" gorm:"type:varchar(50);comment:法人代表;column:legal_person;size:50;"`                              // 规则：可选，法人代表姓名，长度不超过50字符
	RegisteredAddress *string    `json:"registeredAddress" form:"registeredAddress" gorm:"type:varchar(255);comment:注册地址;column:registered_address;size:255;"`         // 规则：可选，注册地址，长度不超过255字符
	BusinessScope     *string    `json:"businessScope" form:"businessScope" gorm:"type:varchar(255);comment:经营范围;column:business_scope;size:255;"`                       // 规则：可选，经营范围，长度不超过255字符
	IsEnabled         bool       `json:"isEnabled" form:"isEnabled" gorm:"default:true;not null;comment:商户开关状态;column:is_enabled;index"`                              // 规则：默认为true（正常），商户开关状态，建立索引优化查询
	ValidStartTime    *time.Time `json:"validStartTime" form:"validStartTime" gorm:"comment:有效开始时间;column:valid_start_time;"`                                        // 规则：可选，有效开始时间
	ValidEndTime      *time.Time `json:"validEndTime" form:"validEndTime" gorm:"comment:有效结束时间;column:valid_end_time;"`                                              // 规则：可选，有效结束时间
	MerchantLevel     MerchantLevel `json:"merchantLevel" form:"merchantLevel" gorm:"not null;comment:商户等级;column:merchant_level;" binding:"required"`                 // 规则：必填，商户等级枚举值（1-普通商户 2-高级商户 3-VIP商户）
	Address           *string    `json:"address" form:"address" gorm:"type:varchar(255);comment:地址;column:address;size:255;"`                                             // 规则：可选，地址，长度不超过255字符
}

// TableName 商户信息 Merchant自定义表名 merchants
func (Merchant) TableName() string {
	return "merchants"
}

// Validate 验证商户字段是否符合业务规则
func (m *Merchant) Validate() error {
	// 验证商户名称
	if len(m.MerchantName) < 1 || len(m.MerchantName) > 100 {
		return errors.New("商户名称长度必须在1-100字符之间")
	}

	// 验证商户类型
	if m.MerchantType != MerchantTypeEnterprise && m.MerchantType != MerchantTypeIndividual {
		return errors.New("商户类型必须是1(企业)或2(个体)")
	}

	// 验证商户等级
	if m.MerchantLevel != MerchantLevelNormal && m.MerchantLevel != MerchantLevelAdvanced && m.MerchantLevel != MerchantLevelVIP {
		return errors.New("商户等级必须是1(普通)、2(高级)或3(VIP)")
	}

	// 验证联系电话格式
	if m.ContactPhone != nil && *m.ContactPhone != "" {
		phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
		if !phoneRegex.MatchString(*m.ContactPhone) {
			return errors.New("联系电话格式不正确")
		}
	}

	// 验证有效期
	if m.ValidStartTime != nil && m.ValidEndTime != nil {
		if m.ValidEndTime.Before(*m.ValidStartTime) {
			return errors.New("有效结束时间不能早于有效开始时间")
		}
	}

	return nil
}

## 更新日志
- 2025-10-04 15:30:00: 初始文档创建
- 2025-10-04 15:30:15: 内容读取与更新完成
- 2025-10-04 15:31:00: 代码详细检查与数据库验证完成
- 2025-10-04 15:32:00: 优化建议生成完成
- 2025-10-04 15:32:30: 结构体重构示例完成
- 2025-10-04 17:28:45: 数据库字段信息更新完成，补充完整的merchants表字段详情