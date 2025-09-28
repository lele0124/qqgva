# 商户管理功能设计文档

## 概述

商户管理功能是基于 gin-vue-admin 系统的多租户核心业务模块，实现商户维度的数据隔离和权限管理。该功能深度集成现有的用户管理（sys_user）、角色管理（sys_authority）和权限控制体系，以及未来新开发的新功能模块中，通过在核心数据表中增加商户ID字段，实现多租户架构下的精细化权限控制。

### 核心业务特性
- **多租户数据隔离**：通过商户ID实现数据维度的完全隔离
- **员工独立账户管理**：员工可以在多个商户中拥有不同的独立账户，每个账户有独立的用户名、手机号和密码
- **跨商户身份支持**：同一手机号或用户名可在不同商户中创建不同的员工账户，账户间完全独立
- **角色权限继承**：基于现有Casbin RBAC体系，增加商户维度的权限控制
- **统一认证体系**：保持JWT认证机制，扩展支持商户上下文和多账户登录选择

## 技术架构

### 多租户架构设计策略
基于现有gin-vue-admin架构，采用共享数据库、独立Schema的多租户模式：

- **数据隔离方式**：在核心业务表中增加merchant_id字段实现行级数据隔离
- **权限扩展策略**：扩展Casbin RBAC模型，增加商户维度的权限控制
- **认证体系升级**：保持JWT机制，在Token中增加当前商户上下文
- **前端状态管理**：扩展Pinia Store，管理当前商户状态

```mermaid
graph TB
    subgraph "前端单租户层"
        A[商户管理页面]
        B[员工管理组件]
        C[权限控制指令]
    end
    
    subgraph "API接口层"
        D[商户管理API]
        E[员工管理API]
        F[角色管理API]
    end
    
    subgraph "业务服务层"
        G[商户信息服务]
        H[员工管理服务]
        I[角色权限服务]
    end
    
    subgraph "数据模型层"
        J[商户信息模型]
        K[扩展用户模型]
        L[扩展角色模型]
    end
    
    subgraph "权限控制层"
        M[扩展Casbin引擎]
        N[商户权限策略]
        O[数据隔离中间件]
    end
    
    A --> D
    B --> E
    C --> F
    
    D --> G
    E --> H
    F --> I
    
    G --> J
    H --> K
    I --> L
    
    D --> M
    E --> M
    F --> M
    M --> N
    M --> O
```

## 多租户数据模型设计

### 现有表结构扩展

#### 现有数据兼容性处理

**重要说明**：为保证与现有系统的兼容性，在进行多租户改造时需要特别处理现有数据：

- **现有员工数据处理**：所有现有的sys_user表中的员工记录，在扩展MerchantID字段后，默认设置为商户ID = 1
- **现有角色数据处理**：所有现有的sys_authority表中的角色记录，在扩展MerchantID字段后，默认设置为商户ID = 1  
- **默认商户创建**：系统需要预先创建一个ID为1的默认商户，作为现有数据的归属商户
- **超级管理员处理**：根据新的设计，超级管理员也归属于默认商户（ID=1），但通过RoleType=1标识其特殊权限
- **系统角色标识**：现有的系统级角色（如超级管理员角色）通过RoleType字段标识为1

**数据迁移策略**：
```sql
-- 创建默认商户（ID=1）
INSERT INTO sys_merchant (id, merchant_code, merchant_name, merchant_type, level, path, contact_name, contact_phone, contact_email, status, is_enabled, valid_start_time, valid_end_time, merchant_level, operator_id, operator_name, operator_merchant_id, operator_merchant_name) 
VALUES (1, 'DEFAULT_MERCHANT', '默认商户', 'ENTERPRISE', 1, '1', '系统管理员', '13800000000', 'admin@system.com', 'ACTIVE', 1, '2024-01-01 00:00:00', '2099-12-31 23:59:59', 'VIP', 1, '系统', 1, '默认商户');

-- 更新所有现有员工数据，设置merchant_id为1
UPDATE sys_user SET merchant_id = 1 WHERE merchant_id IS NULL;

-- 更新所有现有角色数据，设置merchant_id和role_type
UPDATE sys_authority SET merchant_id = 1, role_type = 3 WHERE merchant_id IS NULL AND authority_name != '超级管理员';
UPDATE sys_authority SET merchant_id = 1, role_type = 1 WHERE authority_name = '超级管理员';


```

#### sys_user 表扩展（员工表）
在现有sys_user表基础上增加商户相关字段：

| 新增字段名 | 类型 | 必填 | 索引 | 说明 | 示例值 |
|-----------|------|------|------|------|--------|
| MerchantID | uint | 是 | 复合索引 | 所属商户ID | 1 |

**字段关系说明**：
- **id**：sys_user表的主键，每条记录的唯一标识
- **phone**：手机号，允许同一手机号在不同商户中创建不同的员工账户
- **merchant_id**：商户ID，标识该记录属于哪个商户
- **username**：登录用户名，在同一商户内必须唯一，不同商户间可以重复
- **authority_id**：角色ID，可以在不同商户中拥有不同角色
- **name**：真实姓名，同一手机号的员工在不同商户中可以使用不同姓名
- **password**：登录密码，每个商户的员工账户有独立的密码

**索引设计优化**：


**优化索引策略**：

移除手机号全局唯一限制，采用商户内字段唯一约束：

**索引设计逻辑**：
- **商户内手机号唯一**：通过`(phone, merchant_id, deleted_at)`确保同一商户内手机号唯一
- **商户内用户名唯一**：通过`(username, merchant_id, deleted_at)`确保同一商户内用户名唯一
- **支持跨商户账户**：同一手机号或用户名可以在不同商户中创建不同的员工账户
- **用户名灵活性**：同一手机号在不同商户中可以使用不同的用户名
- **软删除支持**：通过`deleted_at`字段支持软删除逻辑

**索引创建语句**：
```sql
-- 删除原有的手机号相关索引
DROP INDEX IF EXISTS idx_phone_unique ON sys_user;
DROP INDEX IF EXISTS idx_phone_merchant_unique ON sys_user;

-- 创建商户内手机号唯一索引
CREATE UNIQUE INDEX idx_phone_merchant_unique ON sys_user (phone, merchant_id, deleted_at);

-- 创建商户内用户名唯一索引
CREATE UNIQUE INDEX idx_username_merchant_unique ON sys_user (username, merchant_id, deleted_at);

-- 创建手机号查询索引（用于多账户登录检测）
CREATE INDEX idx_phone_lookup ON sys_user (phone, deleted_at);

-- 创建用户名查询索引（用于多账户登录检测）
CREATE INDEX idx_username_lookup ON sys_user (username, deleted_at);
```

**业务逻辑设计**：

1. **员工创建流程**：
   - 管理员在创建员工时填写：手机号、用户名、商户ID、密码等信息
   - 系统校验：手机号和用户名在该商户内不重复
   - 创建成功后生成唯一的sys_user记录ID
   - 该员工可以使用手机号或用户名+密码登录对应商户

2. **多账户登录检测**：
   - 用户输入手机号或用户名后，系统查询sys_user表
   - 如果找不到任何记录，提示"用户不存在或无登录权限"
   - 如果找到一条记录，直接进入密码验证流程
   - 如果找到多条记录，显示所有对应的商户列表供用户选择

3. **商户选择和密码验证**：
   - 用户选择具体商户后，系统获取该商户下的具体账户信息
   - 验证用户输入的密码是否与该账户的密码匹配
   - 检查商户状态和用户状态是否正常
   - 验证通过后生成JWT Token，包含用户信息和商户上下文

**数据模型示例**：
```sql
-- 员工张三在商户A中的账户
INSERT INTO sys_user (id, username, phone, merchant_id, name, authority_id, password) 
VALUES (1, 'zhangsan_sales', '13800138000', 1, '张三', 5, 'encrypted_password_1');

-- 员工张三在商户B中的账户（不同的用户名和密码）
INSERT INTO sys_user (id, username, phone, merchant_id, name, authority_id, password) 
VALUES (2, 'zhangsan_tech', '13800138000', 2, '张三', 8, 'encrypted_password_2');

-- 员工李四在商户A中的账户（用户名可以与张三在商户B中的相同）
INSERT INTO sys_user (id, username, phone, merchant_id, name, authority_id, password) 
VALUES (3, 'zhangsan_tech', '13900139000', 1, '李四', 6, 'encrypted_password_3');
```

**方案优势**：
2. **灵活性强**：支持同一手机号在多个商户中创建不同账户
3. **数据隔离**：每个账户独立管理，包括独立的密码和权限
4. **登录便捷**：支持手机号或用户名登录，自动检测多账户情况
5. **扩展性好**：在现有表结构基础上只需调整索引，迁移成本低

**数据迁移策略**：
```sql
-- 1. 删除原有的手机号全局唯一索引
DROP INDEX IF EXISTS idx_phone_unique ON sys_user;

-- 2. 为现有数据设置默认商户
UPDATE sys_user SET merchant_id = 1 WHERE merchant_id IS NULL;

-- 3. 创建新的索引
CREATE UNIQUE INDEX idx_phone_merchant_unique ON sys_user (phone, merchant_id, deleted_at);
CREATE UNIQUE INDEX idx_username_merchant_unique ON sys_user (username, merchant_id, deleted_at);
CREATE INDEX idx_phone_lookup ON sys_user (phone, deleted_at);
CREATE INDEX idx_username_lookup ON sys_user (username, deleted_at);

-- 4. 验证数据完整性
SELECT phone, merchant_id, COUNT(*) as count 
FROM sys_user 
WHERE deleted_at IS NULL 
GROUP BY phone, merchant_id 
HAVING COUNT(*) > 1;




#### sys_authority 表扩展（角色表）
在现有sys_authority表基础上增加商户维度：

| 新增字段名 | 类型 | 必填 | 索引 | 说明 | 示例值 |
|-----------|------|------|------|------|--------|
| MerchantID | uint | 是 | 普通索引 | 所属商户ID | 1 |
| RoleType | int | 是 | 普通索引 | 角色类型：1-超级管理员 2-商户管理员 3-商户自定义角色 | 1 |

**角色类型说明**：

1. **超级管理员**（MerchantID = 1, RoleType = 1）
   - 属于默认商户（ID=1），但拥有跨商户权限
   - 通过权限控制实现跨商户数据访问
   - 可以管理所有商户和系统配置

2. **商户管理员**（MerchantID = 具体值, RoleType = 2）
   - 属于特定商户的系统预设管理员角色
   - 只能管理所属商户及其子商户的数据
   - 可以创建和管理商户内的员工和自定义角色

3. **商户自定义角色**（MerchantID = 具体值, RoleType = 3）
   - 由商户管理员创建的自定义角色
   - 权限范围受商户管理员分配限制
   - 如：销售专员、财务人员、客服等

**数据示例**：
```sql
-- 超级管理员（属于默认商户，但有跨商户权限）
INSERT INTO sys_authority (authority_name, merchant_id, role_type) 
VALUES ('超级管理员', 1, 1);

-- 商户1的管理员
INSERT INTO sys_authority (authority_name, merchant_id, role_type) 
VALUES ('商户管理员', 1, 2);

-- 商户2的管理员  
INSERT INTO sys_authority (authority_name, merchant_id, role_type) 
VALUES ('商户管理员', 2, 2);

-- 商户1的自定义角色
INSERT INTO sys_authority (authority_name, merchant_id, role_type) 
VALUES ('销售专员', 1, 3);

-- 商户2的自定义角色
INSERT INTO sys_authority (authority_name, merchant_id, role_type) 
VALUES ('财务人员', 2, 3);
```

#### 商户基础信息模型（sys_merchant）
新增商户主表：

| 字段名 | 类型 | 必填 | 索引 | 说明 | 示例值 |
|--------|------|------|------|------|--------|
| ID | uint | 是 | 主键 | 主键ID | 1 |
| MerchantCode | string | 是 | 唯一索引 | 商户编码 | MERCH20240001 |
| MerchantName | string | 是 | 普通索引 | 商户名称 | XX科技有限公司 |
| MerchantIcon | string | 否 | 无 | 商户图标URL | /uploads/icons/merchant_1.png |
| ParentID | uint | 否 | 普通索引 | 父商户ID（NULL表示顶级商户） | 1 |
| MerchantType | string | 是 | 无 | 商户类型 | ENTERPRISE, INDIVIDUAL |
| Level | int | 是 | 无 | 商户层级（1为顶级） | 1 |
| Path | string | 是 | 无 | 层级路径（如：1/2/3） | 1/2 |
| ContactName | string | 是 | 无 | 联系人姓名 | 张三 |
| ContactPhone | string | 是 | 无 | 联系电话 | 13800138000 |
| ContactEmail | string | 是 | 无 | 联系邮箱 | contact@example.com |
| BusinessLicense | string | 否 | 无 | 营业执照号 | 91110000000000000X |
| LegalPerson | string | 否 | 无 | 法人代表 | 李四 |
| RegisteredAddress | string | 否 | 无 | 注册地址 | 北京市朝阳区XX路XX号 |
| BusinessScope | string | 否 | 无 | 经营范围 | 技术开发、技术服务 |

| IsEnabled | int | 是 | 普通索引 | 商户开关状态：1-正常 0-关闭 | 1 |
| ValidStartTime | time.Time | 否 | 无 | 有效开始时间 | 2024-01-01 00:00:00 |
| ValidEndTime | time.Time | 否 | 无 | 有效结束时间 | 2024-12-31 23:59:59 |
| MerchantLevel | string | 是 | 无 | 商户等级 | BASIC, PREMIUM, VIP |

| OperatorID | uint | 是 | 无 | 操作者用户ID | 1 |
| OperatorName | string | 是 | 无 | 操作者姓名 | 张三 |
| OperatorMerchantID | uint | 否 | 无 | 操作者所属商户ID | 1 |
| OperatorMerchantName | string | 否 | 无 | 操作者所属商户名称 | XX科技有限公司 |
| CreatedAt | time.Time | 是 | 无 | 创建时间 | 2024-01-01 10:00:00 |
| UpdatedAt | time.Time | 是 | 无 | 更新时间 | 2024-01-02 15:30:00 |
| DeletedAt | gorm.DeletedAt | 否 | 无 | 删除时间 | NULL |

**层级关系说明**：
- **ParentID**：指向父商户的ID，NULL表示顶级商户
- **Level**：商户在层级中的深度，顶级商户为1，子商户为2，依此类推
- **Path**：记录从根节点到当前节点的完整路径，便于层级查询



#### 商户状态变更记录模型（sys_merchant_status_log）

| 字段名 | 类型 | 必填 | 索引 | 说明 | 示例值 |
|--------|------|------|------|------|--------|
| ID | uint | 是 | 主键 | 主键ID | 1 |
| MerchantID | uint | 是 | 外键索引 | 商户ID | 1 |
| PreviousStatus | string | 是 | 无 | 变更前状态 | ACTIVE |
| NewStatus | string | 是 | 无 | 变更后状态 | SUSPENDED |
| ChangeReason | string | 是 | 无 | 变更原因 | 系统管理员操作 |
| OperatorID | uint | 是 | 无 | 操作者用户ID | 2 |
| OperatorName | string | 是 | 无 | 操作者姓名 | 张三 |
| OperatorMerchantID | uint | 否 | 无 | 操作者所属商户ID | 1 |
| CreatedAt | time.Time | 是 | 无 | 变更时间 | 2024-01-02 15:00:00 |

### 商户状态控制设计

#### 商户开关状态说明

**IsEnabled 字段作用**：
- **1（正常）**：商户处于正常运营状态，商户内所有用户可以正常登录和使用系统
- **0（关闭）**：商户处于关闭状态，商户内所有用户无法登录系统平台

**有效时间设计**：
- **ValidStartTime**：商户有效开始时间，用于记录商户的合同或授权起始时间
- **ValidEndTime**：商户有效结束时间，用于记录商户的合同或授权结束时间
- **注意**：有效时间段目前仅做记录用途，不对登录逻辑进行处理

#### 登录限制逻辑

**影响范围**：
当商户状态设置为关闭（IsEnabled = 0）时，以下用户将无法登录系统：
- 该商户的所有员工（RoleType = 3）
- 该商户的管理员（RoleType = 2）
- **例外**：超级管理员（RoleType = 1）不受商户状态影响，可以正常登录

**登录校验流程**：

```mermaid
sequenceDiagram
    participant U as 用户
    participant API as 登录API
    participant AUTH as 认证服务
    participant DB as 数据库
    
    U->>API: 提交登录请求
    API->>AUTH: 验证用户名密码
    AUTH->>DB: 查询用户信息
    DB-->>AUTH: 返回用户数据
    
    alt 超级管理员
        AUTH->>AUTH: 跳过商户状态检查
        AUTH-->>API: 登录成功
    else 商户用户
        AUTH->>DB: 查询所属商户状态
        DB-->>AUTH: 返回商户IsEnabled字段
        
        alt IsEnabled = 1
            AUTH-->>API: 登录成功
        else IsEnabled = 0
            AUTH-->>API: 登录失败：商户状态关闭
        end
    end
    
    API-->>U: 返回登录结果
```

**错误提示信息**：
```json
{
    "code": 40303,
    "message": "商户状态关闭，请联系管理员",
    "data": {
        "merchantId": 2,
        "merchantName": "XX科技有限公司",
        "isEnabled": 0
    }
}
```

**实现要点**：
1. 在用户登录时，需要同时查询用户所属商户的IsEnabled状态
2. 超级管理员不受商户状态限制，可以管理关闭的商户
3. 已登录的用户在商户被关闭后，需要在下次请求时进行状态检查
4. 提供明确的错误提示，帮助用户理解限制原因

### 商户层级关系设计

```mermaid
erDiagram
    SysMerchant ||--o{ SysMerchant : "父子关系"

    
    SysMerchant ||--o{ SysMerchant : "父子关系"
    SysMerchant ||--o{ SysUser : "拥有员工"
    SysMerchant ||--o{ SysAuthority : "拥有角色"
    SysMerchant ||--o{ SysMerchantStatusLog : "状态记录"
    
    SysUser }|--|| SysAuthority : "主角色"
    SysUser ||--o{ SysUserAuthority : "多角色"
    
    SysAuthority ||--o{ SysUserAuthority : "用户关联"
    
    SysMerchant {
        uint ID PK
        string MerchantCode UK
        string MerchantName
        string MerchantIcon
        uint ParentID FK
        string MerchantType
        int Level
        string Path
        string Status
        int IsEnabled
        time ValidStartTime
        time ValidEndTime
        string MerchantLevel
        uint AdminUserID FK
        uint OperatorID FK
        string OperatorName
        uint OperatorMerchantID
        string OperatorMerchantName
    }
    
    SysUser {
        uint ID PK
        string Username UK
        string Phone UK
        uint MerchantID FK
        uint AuthorityId FK
        bool Enable
    }
    
    SysAuthority {
        uint AuthorityId PK
        string AuthorityName
        uint MerchantID FK
        int RoleType
        uint ParentId FK
    }
    
    SysMerchantStatusLog {
        uint ID PK
        uint MerchantID FK
        string PreviousStatus
        string NewStatus
        string ChangeReason
        uint OperatorID FK
        string OperatorName
    }
```

### 商户层级结构示例

```mermaid
graph TD
    A[总公司 - Level 1]
    B[分公司A - Level 2]
    C[分公司B - Level 2]
    D[部门A1 - Level 3]
    E[部门A2 - Level 3]
    F[部门B1 - Level 3]
    G[小组A1-1 - Level 4]
    H[小组A1-2 - Level 4]
    
    A --> B
    A --> C
    B --> D
    B --> E
    C --> F
    D --> G
    D --> H
    
    style A fill:#e1f5fe
    style B fill:#f3e5f5
    style C fill:#f3e5f5
    style D fill:#e8f5e8
    style E fill:#e8f5e8
    style F fill:#e8f5e8
    style G fill:#fff3e0
    style H fill:#fff3e0
```

## 多租户API接口设计

### 接口规范说明
- 遵循 RESTful API 设计规范，增加商户上下文
- 统一使用 JSON 格式进行数据交换
- 扩展JWT认证，Token中包含当前用户身份和商户ID信息
- 集成商户维度的Casbin权限控制
- 所有业务接口自动进行商户数据隔离
- 支持多商户身份切换和上下文管理

### 多商户登录功能设计

支持同一用户在多个商户中拥有不同身份的智能登录功能。

#### 登录流程设计

```mermaid
sequenceDiagram
    participant U as 用户
    participant LOGIN as 登录界面
    participant API as 登录API
    participant AUTH as 认证服务
    participant DB as 数据库
    
    U->>LOGIN: 输入手机号或用户名
    LOGIN->>API: 提交登录请求
    API->>AUTH: 验证手机号用户名
    AUTH->>DB: 根据手机号或用户名查询所有身份记录
    DB-->>AUTH: 返回用户身份列表
    
    alt 手机号或用户名不存在
        AUTH-->>API: 返回登录失败
        API-->>LOGIN: 显示错误信息
    else 只有一个身份
        AUTH->>AUTH: 验证密码和商户状态
        alt 验证成功且商户正常
            AUTH->>AUTH: 生成JWT Token
            AUTH-->>API: 返回登录成功+Token
            API-->>LOGIN: 跳转到主界面
        else 密码错误或商户关闭
            AUTH-->>API: 返回对应错误信息
            API-->>LOGIN: 显示错误信息
        end
    else 有多个身份

            AUTH-->>API: 返回商户选择列表
            API-->>LOGIN: 显示商户选择界面
            U->>LOGIN: 选择商户身份和密码
        AUTH->>AUTH: 验证密码（对应商户身份和密码）和商户状态
        alt 密码验证成功
            LOGIN->>API: 提交商户选择
            API->>AUTH: 验证选择的商户身份
            AUTH->>AUTH: 检查选择商户状态
            alt 商户正常
                AUTH->>AUTH: 生成JWT Token
                AUTH-->>API: 返回登录成功+Token
                API-->>LOGIN: 跳转到主界面
            else 商户关闭
                AUTH-->>API: 返回商户状态错误
                API-->>LOGIN: 显示商户关闭信息
            end
        else 密码验证失败
            AUTH-->>API: 返回密码错误
            API-->>LOGIN: 显示密码错误信息
        end
    end
```

#### 登录接口设计

**第一步：统一登录接口**
- **接口路径**：`POST /api/v1/auth/login`
- **请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| phone | string | 是 | 手机号 |
| password | string | 是 | 密码 |

**响应格式** - 只有一个身份：
```json
{
    "code": 0,
    "message": "登录成功",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "user": {
            "id": 1,
            "username": "zhangsan_sales",
            "name": "张三",
            "phone": "13800138000",
            "merchantInfo": {
                "merchantId": 1,
                "merchantName": "XX科技有限公司",
                "merchantIcon": "/uploads/icons/merchant_1.png",
                "authorityId": 5,
                "authorityName": "销售经理",
                "roleType": 3
            }
        },
        "needMerchantSelect": false
    }
}
```

**响应格式** - 多个账户：
```json
{
    "code": 10001,
    "message": "请选择登录商户",
    "data": {
        "name": "张三",
        "phone": "13800138000",
        "accounts": [
            {
                "id": 1,
                "merchantId": 1,
                "merchantName": "XX科技有限公司",
                "merchantIcon": "/uploads/icons/merchant_1.png",
                "username": "zhangsan_sales",
                "authorityName": "销售经理",
                "roleType": 3,
                "merchantStatus": "ACTIVE",
                "merchantEnabled": 1,
                "isDefault": true
            },
            {
                "id": 2,
                "merchantId": 2,
                "merchantName": "YY贸易有限公司",
                "merchantIcon": "/uploads/icons/merchant_2.png",
                "username": "zhangsan_tech",
                "authorityName": "技术顾问",
                "roleType": 3,
                "merchantStatus": "ACTIVE",
                "merchantEnabled": 1,
                "isDefault": false
            }
        ],
        "needMerchantSelect": true,
        "tempToken": "temp_token_for_merchant_selection"
    }
}
```

**第二步：商户身份选择接口**
- **接口路径**：`POST /api/v1/auth/select-identity`
- **请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| identityId | uint | 是 | 选择的账户记录ID（sys_user.id） |
| tempToken | string | 是 | 临时令牌（用于验证） |

**响应格式**：
```json
{
    "code": 0,
    "message": "登录成功",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "user": {
            "id": 2,
            "username": "zhangsan_tech",
            "name": "张三",
            "phone": "13800138000",
            "merchantInfo": {
                "merchantId": 2,
                "merchantName": "YY贸易有限公司",
                "merchantIcon": "/uploads/icons/merchant_2.png",
                "authorityId": 8,
                "authorityName": "技术顾问",
                "roleType": 3
            }
        }
    }
}
```

#### JWT Token 扩展设计

在JWT Payload中增加用户账户和商户上下文信息：

```json
{
    "userId": 2,
    "username": "zhangsan_tech",
    "phone": "13800138000",
    "name": "张三",
    "merchantId": 2,
    "merchantName": "YY贸易有限公司",
    "authorityId": 8,
    "authorityName": "技术顾问",
    "roleType": 3,
    "iat": 1640995200,
    "exp": 1641081600
}
```

#### 账户切换功能

**账户切换接口**
- **接口路径**：`POST /api/v1/auth/switch-account`
- **权限要求**：需要已登录用户的有效Token
- **请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| targetAccountId | uint | 是 | 目标账户记录ID |

**功能说明**：
- 用户在不退出的情况下切换到另一个商户账户
- 系统验证目标账户是否属于当前用户（相同手机号）
- 检查目标商户的状态和账户在该商户中的状态
- 生成新的JWT Token包含目标账户信息

**响应格式**：
```json
{
    "code": 0,
    "message": "账户切换成功",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "user": {
            "id": 3,
            "username": "zhangsan_finance",
            "name": "张三",
            "phone": "13800138000",
            "merchantInfo": {
                "merchantId": 3,
                "merchantName": "ZZ金融服务有限公司",
                "merchantIcon": "/uploads/icons/merchant_3.png",
                "authorityId": 12,
                "authorityName": "财务专员",
                "roleType": 3
            }
        }
    }
}
```

### 商户信息管理接口

#### 创建商户
- **接口路径**：`POST /api/v1/merchant`
- **权限要求**：系统管理员权限（不受商户限制）
- **请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| merchantName | string | 是 | 商户名称 |
| merchantIcon | string | 否 | 商户图标URL |
| parentId | uint | 否 | 父商户ID（不填为顶级商户） |
| merchantType | string | 是 | 商户类型（ENTERPRISE/INDIVIDUAL） |
| contactName | string | 是 | 联系人姓名 |
| contactPhone | string | 是 | 联系电话 |
| contactEmail | string | 是 | 联系邮箱 |
| businessLicense | string | 否 | 营业执照号 |
| legalPerson | string | 否 | 法人代表 |
| registeredAddress | string | 否 | 注册地址 |
| businessScope | string | 否 | 经营范围 |
| merchantLevel | string | 是 | 商户等级（BASIC/PREMIUM/VIP） |
| isEnabled | int | 否 | 商户开关状态（默认1-正常） |
| validStartTime | string | 否 | 有效开始时间 |
| validEndTime | string | 否 | 有效结束时间 |


**层级关系处理**：
- 如果提供parentId，系统会自动计算level和path
- 顶级商户：level=1, path=merchantId
- 子商户：level=父商户level+1, path=父商户path/merchantId

**请求示例**：
```json
{
  "merchantName": "XX科技有限公司",
  "merchantIcon": "/uploads/icons/merchant_logo.png",
  "parentId": 1,
  "merchantType": "ENTERPRISE",
  "contactName": "张三",
  "contactPhone": "13800138000",
  "contactEmail": "contact@xxtech.com",
  "businessLicense": "91110000000000000X",
  "legalPerson": "张三",
  "registeredAddress": "北京市朝阳区XX路XX号",
  "businessScope": "技术开发、技术服务",
  "merchantLevel": "PREMIUM",
  "isEnabled": 1,
  "validStartTime": "2024-01-01 00:00:00",
  "validEndTime": "2024-12-31 23:59:59",

}
```

**响应格式**：
```json
{
  "code": 0,
  "message": "创建成功",
  "data": {
    "merchantId": 5,
    "merchantCode": "MERCH20240005",
    "merchantName": "XX科技有限公司",
    "level": 2,
    "path": "1/5",

    "createdAt": "2024-01-01T10:00:00Z"
  }
}
```

#### 查询商户层级结构
- **接口路径**：`GET /api/v1/merchant/tree`
- **权限要求**：根据角色自动过滤数据
- **查询参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| rootId | uint | 否 | 指定根节点ID（不填则显示所有顶级商户） |
| maxLevel | int | 否 | 最大层级深度（限制展示层数） |
| includeInactive | bool | 否 | 是否包含非活跃商户 |

- **响应格式**：
```json
{
    "code": 0,
    "message": "获取成功",
    "data": {
        "merchantTree": [
            {
                "merchantId": 1,
                "merchantName": "总公司",
                "merchantIcon": "/uploads/icons/merchant_1.png",
                "parentId": null,
                "level": 1,
                "path": "1",
                "status": "ACTIVE",
                "children": [
                    {
                        "merchantId": 2,
                        "merchantName": "分公司A",
                        "merchantIcon": "/uploads/icons/merchant_2.png",
                        "parentId": 1,
                        "level": 2,
                        "path": "1/2",
                        "status": "ACTIVE",
                        "children": []
                    }
                ]
            }
        ]
    }
}
```

#### 查询商户子级列表
- **接口路径**：`GET /api/v1/merchant/:id/children`
- **权限要求**：商户管理员或系统管理员
- **功能说明**：获取指定商户的直接子级商户列表

#### 查询商户祖先路径
- **接口路径**：`GET /api/v1/merchant/:id/ancestors`
- **权限要求**：商户管理员或系统管理员
- **功能说明**：获取从根节点到当前商户的完整路径

#### 移动商户位置
- **接口路径**：`PUT /api/v1/merchant/:id/move`
- **权限要求**：仅系统管理员可操作
- **请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| newParentId | uint | 否 | 新的父商户ID（NULL表示移动到顶级） |

**校验规则**：
- 不能将商户移动到其自身或其子孙节点下
- 移动后自动更新所有子孙节点的level和path

#### 查询商户列表
- **接口路径**：`GET /api/v1/merchant/list`
- **权限要求**：根据角色自动过滤数据
  - 系统管理员：可查看所有商户
  - 商户管理员：只能查看所属商户及其子级商户
- **查询参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码（默认1） |
| pageSize | int | 否 | 每页数量（默认10） |
| merchantName | string | 否 | 商户名称（模糊查询） |
| merchantType | string | 否 | 商户类型筛选 |
| status | string | 否 | 状态筛选 |
| parentId | uint | 否 | 父商户ID筛选 |
| level | int | 否 | 层级筛选 |
| treeView | bool | 否 | 是否返回树形结构 |

#### 更新商户信息
- **接口路径**：`PUT /api/v1/merchant/:id`
- **权限要求**：商户管理员或系统管理员
- **数据隔离**：自动校验商户权限，只能修改所属商户

#### 商户状态管理
- **接口路径**：`PUT /api/v1/merchant/:id/status`
- **权限要求**：仅系统管理员可操作

#### 创建商户员工
- **接口路径**：`POST /api/v1/merchant/user`
- **权限要求**：商户管理员权限或超级管理员
- **功能说明**：
- 在当前商户中创建新员工或为一个手机号在当前商户中创建账户
- 支持同一手机号在不同商户中创建不同账户
- 每个账户拥有独立的用户名、密码、角色
- **请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| phone | string | 是 | 手机号 |　（在当前商户内唯一）
| name | string | 是 | 真实姓名 |
| username | string | 是 | 登录用户名（在当前商户内唯一） |
| password | string | 是 | 登录密码 |
| email | string | 否 | 邮箱地址 |
| nickName | string | 否 | 显示昵称 |
| headerImg | string | 否 | 头像地址 |
| authorityId | uint | 是 | 角色ID（必须是当前商户的角色） |
| enable | int | 否 | 启用状态（默认1-启用） |

**业务逻辑**：
1. **创建新员工**：
   - 检查手机号和用户名在当前商户内是否已存在
   - 在当前商户中创建新的员工记录
   - 设置用户在该商户中的特定信息（用户名、角色、密码等）

2. **支持同一手机号或用户名在不同商户中创建账户**：
   - 同一手机号或用户名可以在不同商户中创建不同的员工账户
   - 每个账户拥有独立的用户名、密码、角色
   - 同一商户内手机号和用户名必须唯一

**响应格式**：
```json
{
    "code": 0,
    "message": "创建成功",
    "data": {
        "id": 5,
        "username": "lisi_sales",
        "name": "李四",
        "phone": "13800138001",
        "merchantId": 1,
        "authorityId": 6,
        "createdAt": "2024-01-01T10:00:00Z"
    }
}
```

#### 查询商户员工列表
- **接口路径**：`GET /api/v1/merchant/user/list`
- **权限要求**：商户管理员或超级管理员
- **数据隔离**：自动过滤，只返回当前商户的员工
- **查询参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码（默认1） |
| pageSize | int | 否 | 每页数量（默认10） |
| name | string | 否 | 姓名模糊查询 |
| phone | string | 否 | 手机号模糊查询 |
| username | string | 否 | 用户名模糊查询 |
| authorityId | uint | 否 | 角色ID筛选 |
| enable | int | 否 | 启用状态筛选 |

**响应格式**：
```json
{
    "code": 0,
    "message": "获取成功",
    "data": {
        "list": [
            {
                "id": 1,
                "username": "zhangsan_sales",
                "name": "张三",
                "phone": "13800138000",
                "email": "zhangsan@example.com",
                "nickName": "张三",
                "headerImg": "/uploads/avatar/zhangsan.jpg",
                "enable": 1,
                "merchantId": 1,
                "authorityInfo": {
                    "authorityId": 5,
                    "authorityName": "销售经理",
                    "roleType": 3
                },
                "createdAt": "2024-01-01T10:00:00Z",
                "hasOtherAccounts": true,
                "otherMerchantCount": 2
            }
        ],
        "total": 25,
        "page": 1,
        "pageSize": 10
    }
}
```

#### 员工账户状态管理
- **接口路径**：`PUT /api/v1/merchant/user/:id/status`
- **权限要求**：商户管理员或超级管理员
- **请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| enable | bool | 是 | 在当前商户中的启用状态 |
| reason | string | 否 | 状态变更原因 |

**注意**：此操作仅影响该员工在当前商户中的状态，不影响其在其他商户中的账户。

#### 移除商户员工
- **接口路径**：`DELETE /api/v1/merchant/user/:id`
- **权限要求**：商户管理员或超级管理员
- **功能说明**：仅删除该员工在当前商户中的账户记录，不影响其在其他商户中的账户

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| reason | string | 否 | 删除原因 |
| transferAuthorityId | uint | 否 | 如果该员工有待处理事务，转移给的角色ID |

**响应格式**：
```json
{
    "code": 0,
    "message": "移除成功",
    "data": {
        "removedAccountId": 5,
        "merchantId": 1,
        "hasOtherAccounts": false,
        "message": "已移除该用户在当前商户中的账户"
    }
}
```

### 员工角色权限管理接口

#### 获取商户内可分配角色权限
- **接口路径**：`GET /api/v1/merchant/role/permissions/assignable`
- **权限要求**：商户管理员或超级管理员
- **功能说明**：获取当前商户内可以分配给员工角色的权限列表

#### 分配角色权限
- **接口路径**：`PUT /api/v1/merchant/role/:roleId/permissions`
- **权限要求**：商户管理员或超级管理员
- **请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| permissions | array | 是 | 权限代码数组 |

#### 员工角色分配
- **接口路径**：`PUT /api/v1/merchant/user/:userId/roles`
- **权限要求**：商户管理员或超级管理员
- **请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| roleIds | array | 是 | 角色ID数组 |

### 商户权限分配管理接口

#### 获取可分配权限列表
- **接口路径**：`GET /api/v1/system/permissions/assignable`
- **权限要求**：仅超级管理员可操作
- **功能说明**：获取所有可以分配给商户管理员的权限列表
- **响应格式**：
```json
{
    "code": 0,
    "message": "获取成功",
    "data": {
        "permissions": [
            {
                "permissionCode": "merchant:info:view",
                "permissionName": "商户信息查看",
                "category": "merchant_info",
                "description": "查看商户基本信息",
                "isRequired": true
            },
            {
                "permissionCode": "merchant:info:update",
                "permissionName": "商户信息编辑",
                "category": "merchant_info",
                "description": "编辑商户基本信息",
                "isRequired": false
            },
            {
                "permissionCode": "merchant:user:list",
                "permissionName": "商户员工查看",
                "category": "employee_management",
                "description": "查看商户员工列表",
                "isRequired": false
            },
            {
                "permissionCode": "merchant:user:create",
                "permissionName": "商户员工创建",
                "category": "employee_management",
                "description": "创建新的商户员工",
                "isRequired": false
            },
            {
                "permissionCode": "merchant:user:update",
                "permissionName": "商户员工编辑",
                "category": "employee_management",
                "description": "编辑商户员工信息",
                "isRequired": false
            },
            {
                "permissionCode": "merchant:user:delete",
                "permissionName": "商户员工删除",
                "category": "employee_management",
                "description": "删除商户员工",
                "isRequired": false
            },
            {
                "permissionCode": "merchant:role:list",
                "permissionName": "商户角色查看",
                "category": "role_management",
                "description": "查看商户角色列表",
                "isRequired": false
            },
            {
                "permissionCode": "merchant:role:create",
                "permissionName": "商户角色创建",
                "category": "role_management",
                "description": "创建新的商户角色",
                "isRequired": false
            },
            {
                "permissionCode": "merchant:role:update",
                "permissionName": "商户角色编辑",
                "category": "role_management",
                "description": "编辑商户角色信息",
                "isRequired": false
            },
            {
                "permissionCode": "merchant:role:delete",
                "permissionName": "商户角色删除",
                "category": "role_management",
                "description": "删除商户角色",
                "isRequired": false
            }
        ],
        "categories": [
            {
                "categoryCode": "merchant_info",
                "categoryName": "商户信息管理",
                "description": "商户基本信息管理相关权限"
            },
            {
                "categoryCode": "employee_management",
                "categoryName": "员工管理",
                "description": "商户员工管理相关权限"
            },
            {
                "categoryCode": "role_management",
                "categoryName": "角色管理",
                "description": "商户角色管理相关权限"
            }
        ]
    }
}
```

#### 分配商户权限
- **接口路径**：`PUT /api/v1/merchant/:id/permissions`
- **权限要求**：仅超级管理员可操作
- **请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| permissions | array | 是 | 权限代码数组 |
| assignToAdmin | bool | 否 | 是否同时分配给商户管理员 |

- **请求示例**：
```json
{
    "permissions": [
        "merchant:info:view",
        "merchant:info:update",
        "merchant:user:list",
        "merchant:user:create",
        "merchant:user:update",
        "merchant:role:list",
        "merchant:role:create"
    ],
    "assignToAdmin": true
}
```

- **响应格式**：
```json
{
    "code": 0,
    "message": "权限分配成功",
    "data": {
        "merchantId": 1,
        "assignedPermissions": [
            "merchant:info:view",
            "merchant:info:update",
            "merchant:user:list",
            "merchant:user:create",
            "merchant:user:update",
            "merchant:role:list",
            "merchant:role:create"
        ],
        "adminUserId": 2,
        "adminUpdated": true,
        "assignTime": "2024-01-01T10:00:00Z"
    }
}
```

#### 获取商户权限
- **接口路径**：`GET /api/v1/merchant/:id/permissions`
- **权限要求**：超级管理员或商户管理员
- **功能说明**：获取指定商户已分配的权限列表
- **响应格式**：
```json
{
    "code": 0,
    "message": "获取成功",
    "data": {
        "merchantId": 1,
        "merchantName": "XX科技有限公司",
        "permissions": [
            {
                "permissionCode": "merchant:info:view",
                "permissionName": "商户信息查看",
                "category": "merchant_info",
                "assignedAt": "2024-01-01T10:00:00Z",
                "assignedBy": "super_admin"
            }
        ],
        "lastUpdated": "2024-01-01T10:00:00Z"
    }
}
```

## 多租户业务流程设计

### 三级权限体系流程设计

#### 超级管理员分配商户权限流程

```mermaid
sequenceDiagram
    participant SA as 超级管理员
    participant API as 权限管理API
    participant PS as 权限服务
    participant CS as Casbin服务
    participant DB as 数据库
    
    SA->>API: 查看可分配权限列表
    API->>PS: 获取所有非系统级权限
    PS-->>API: 返回权限列表
    API-->>SA: 显示可分配权限
    
    SA->>API: 为商户分配权限
    API->>PS: 验证权限合法性
    PS->>CS: 为商户管理员角色分配权限
    CS->>DB: 更新Casbin策略表
    
    alt 同时分配给管理员
        PS->>CS: 为管理员用户添加权限
        CS->>DB: 更新用户权限表
    end
    
    PS-->>API: 返回分配结果
    API-->>SA: 显示分配成功
```

#### 商户管理员分配员工权限流程

```mermaid
sequenceDiagram
    participant MA as 商户管理员
    participant API as 角色管理API
    participant RS as 角色服务
    participant CS as Casbin服务
    participant DB as 数据库
    
    MA->>API: 查看商户内可分配权限
    API->>RS: 获取当前商户已有权限
    RS-->>API: 返回权限列表
    API-->>MA: 显示可分配权限
    
    MA->>API: 为员工角色分配权限
    API->>RS: 验证权限合法性
    note right of RS: 检查是否在商户允许范围内
    RS->>CS: 为员工角色分配权限
    CS->>DB: 更新角色权限表
    
    MA->>API: 为员工分配角色
    API->>RS: 验证角色合法性
    RS->>CS: 为员工添加角色
    CS->>DB: 更新用户角色关联
    
    RS-->>API: 返回分配结果
    API-->>MA: 显示分配成功
```

### 商户移动操作流程

```mermaid
sequenceDiagram
    participant Admin as 系统管理员
    participant API as API接口
    participant MS as 商户服务
    participant DB as 数据库
    
    Admin->>API: PUT /api/v1/merchant/:id/move
    API->>MS: 调用商户移动服务
    
    MS->>DB: 查询当前商户信息
    MS->>DB: 查询新父商户信息
    MS->>MS: 校验是否形成循环引用
    
    alt 校验通过
        MS->>DB: 更新当前商户的ParentID
        MS->>MS: 计算新的level和path
        MS->>DB: 更新当前商户level和path
        
        loop 更新所有子孙节点
            MS->>DB: 查询子孙节点
            MS->>MS: 重新计算level和path
            MS->>DB: 批量更新子孙节点
        end
        
        MS-->>API: 返回成功结果
    else 校验失败
        MS-->>API: 返回错误信息
    end
    
    API-->>Admin: 返回操作结果
```

### 单一商户员工登录流程

```mermaid
sequenceDiagram
    participant U as 员工
    participant F as 前端页面
    participant API as 登录API
    participant AUTH as 认证服务
    participant DB as 数据库
    
    U->>F: 输入用户名密码
    F->>API: POST /api/v1/login
    API->>AUTH: 验证用户凭据
    AUTH->>DB: 查询用户信息
    AUTH->>AUTH: 生成包含商户ID的JWT
    AUTH-->>F: 返回Token和商户信息
    F->>F: 进入商户后台
```

### 商户内数据隔离流程

```mermaid
sequenceDiagram
    participant F as 前端请求
    participant MW as 数据隔离中间件
    participant API as API接口
    participant SVC as 业务服务
    participant DB as 数据库
    
    F->>MW: 发起业务请求
    MW->>MW: 从jWT提取商户ID
    MW->>MW: 设置商户上下文
    MW->>API: 转发请求
    
    API->>SVC: 调用业务服务
    SVC->>SVC: 获取当前商户ID
    
    alt 查询操作
        SVC->>DB: WHERE merchant_id = 当前商户ID
    else 创建操作
        SVC->>SVC: 自动设置 merchant_id
        SVC->>DB: INSERT 包含商户ID的数据
    else 更新操作
        SVC->>DB: UPDATE WHERE id=? AND merchant_id=当前商户ID
    end
    
    DB-->>SVC: 返回结果
    SVC-->>API: 返回业务数据
    API-->>MW: 返回响应
    MW-->>F: 返回最终结果
```



### 原有系统界面兼容性改造

#### 员工列表页面改造

**前端页面修改（user.vue）**

1. **表格列增加**
   - 在现有表格中增加"商户信息"列，显示商户ID和商户名称
   - 位置：在"用户角色"列之后，"状态"列之前

2. **搜索条件扩展**
   - 增加商户筛选下拉框，支持按商户过滤员工列表
   - 超级管理员可看到所有商户选项
   - 商户管理员只能看到自己的商户及子商户

3. **数据获取逻辑**
   - 调用API时自动根据当前用户权限过滤数据
   - 界面显示时关联查询商户信息

**后端接口修改（GetUserList）**

1. **请求参数扩展**
   - GetUserList结构体增加MerchantId字段，支持按商户筛选
   - 增加MerchantName字段，支持按商户名称模糊搜索

2. **查询逻辑增强**
   - 在GetUserInfoList方法中增加商户关联查询
   - 使用Preload加载商户信息
   - 根据用户权限自动过滤可见数据范围

3. **返回数据结构**
   - SysUser结构体关联Merchant信息
   - 返回数据包含完整的商户ID和商户名称

#### 角色列表页面改造

**前端页面修改（authority.vue）**

1. **表格列增加**
   - 在现有表格中增加"所属商户"列，显示商户ID和商户名称
   - 位置：在"角色名称"列之后，"更新时间"列之前

2. **商户筛选功能**
   - 增加商户下拉筛选器
   - 根据用户权限控制可见商户范围

3. **角色创建/编辑**
   - 新增角色时需要选择所属商户
   - 商户管理员只能在自己的商户下创建角色

**后端接口修改（GetAuthorityList）**

1. **查询逻辑修改**
   - GetAuthorityInfoList方法增加商户过滤逻辑
   - 根据当前用户的商户权限过滤角色列表
   - 支持商户层级查询（父子商户关系）

2. **数据关联**
   - SysAuthority结构体关联Merchant信息
   - 查询时预加载商户数据

#### 接口兼容性设计

**请求参数结构扩展**

| 接口 | 原有参数 | 新增参数 | 说明 |
|------|----------|----------|---------|
| GetUserList | Username, NickName, Phone, Email | MerchantId, MerchantName | 商户筛选支持 |
| GetAuthorityList | PageInfo | MerchantId, MerchantName | 商户角色筛选 |

**响应数据结构扩展**

| 实体 | 原有字段 | 新增字段 | 说明 |
|------|----------|----------|---------|
| SysUser | 基础用户信息 | MerchantId, MerchantName, Merchants | 商户关联信息 |
| SysAuthority | 基础角色信息 | MerchantId, MerchantName | 商户归属信息 |

**权限控制策略**

1. **数据可见性**
   - 超级管理员：查看所有商户的员工和角色
   - 商户管理员：仅查看本商户及子商户的员工和角色
   - 普通员工：仅查看本人信息

2. **操作权限**
   - 员工管理：只能管理所属商户范围内的员工
   - 角色管理：只能管理所属商户的角色
   - 跨商户操作需要相应权限验证

#### 界面交互优化

**商户信息展示**
- 采用两行显示格式：第一行显示商户名称，第二行显示商户ID
- 使用颜色区分不同商户，提升用户体验
- 支持点击商户名称快速筛选该商户下的数据

**筛选器设计**
- 商户下拉框支持搜索功能
- 记住用户的筛选偏好
- 提供"全部商户"选项（仅超级管理员可见）

**兼容性保证**
- 现有功能保持不变
- 新增字段采用非必填设计
- 向前兼容旧版本API调用

### 现有员工登录界面改造

#### 登录流程优化设计

**多商户员工登录场景分析**：
- 员工在登录时，系统检查该手机号在多个商户中是否存在
- 如果员工只属于一个商户：直接进入该商户后台
- 如果员工属于多个商户：显示商户选择界面，让用户选择要进入的商户
- 保持现有单商户用户体验不变，多商户功能对原用户透明

#### 多商户检测接口

**新增接口：检测用户多商户账户**
- **接口路径**：`POST /api/v1/auth/check-merchants`
- **调用时机**：用户输入用户名/手机号后，密码验证前
- **请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| identifier | string | 是 | 用户标识（手机号或用户名） |
| identifierType | string | 否 | 标识类型：phone/username（默认自动检测） |

- **后端检测逻辑**：
```sql
-- 自动检测标识类型（正则匹配手机号格式）
SELECT u.id, u.username, u.phone, u.merchant_id, u.name,
       m.merchant_name, m.merchant_icon, m.status,
       CASE 
         WHEN a.role_type = 1 THEN '超级管理员'
         WHEN a.role_type = 2 THEN '商户管理员'
         ELSE '员工'
       END as user_role
FROM sys_user u
LEFT JOIN sys_merchant m ON u.merchant_id = m.id
LEFT JOIN sys_authority a ON u.authority_id = a.authority_id
WHERE (u.phone = ? OR u.username = ?) 
  AND u.deleted_at IS NULL 
  AND u.enable = 1
  AND m.is_enabled = 1;
```

- **响应格式**：
```json
{
  "code": 0,
  "message": "检测成功",
  "data": {
    "isMultiMerchant": true,
    "merchantCount": 2,
    "identifierType": "phone",
    "merchants": [
      {
        "userId": 10,
        "username": "zhangsan_a",
        "merchantId": 1,
        "merchantName": "XX科技有限公司",
        "merchantIcon": "/uploads/icons/merchant_1.png",
        "userRole": "商户管理员",
        "status": "ACTIVE"
      },
      {
        "userId": 25,
        "username": "zhangsan_b", 
        "merchantId": 3,
        "merchantName": "YY贸易公司",
        "merchantIcon": "/uploads/icons/merchant_3.png",
        "userRole": "员工",
        "status": "ACTIVE"
      }
    ]
  }
}
```

#### 多商户登录流程

```mermaid
sequenceDiagram
    participant U as 员工
    participant F as 前端页面
    participant API as 登录API
    participant AUTH as 认证服务
    participant DB as 数据库
    
    U->>F: 输入手机号
    F->>API: POST /api/v1/auth/check-merchants
    API->>DB: 查询手机号对应的商户列表
    
    alt 单一商户员工
        DB-->>API: 返回单个商户信息
        API-->>F: isMultiMerchant=false
        F->>F: 显示密码输入框
        U->>F: 输入密码
        F->>API: POST /api/v1/login (包含merchantId)
        API->>AUTH: 验证密码并生成Token
        AUTH-->>F: 返回Token和商户信息
        F->>F: 直接进入商户后台
    else 多商户员工
        DB-->>API: 返回多个商户信息
        API-->>F: isMultiMerchant=true + 商户列表
        F->>F: 显示商户选择界面
        U->>F: 选择目标商户
        F->>F: 显示密码输入框
        U->>F: 输入密码
        F->>API: POST /api/v1/login (包含选中merchantId)
        API->>AUTH: 验证密码并生成Token
        AUTH-->>F: 返回Token和商户信息
        F->>F: 进入选定商户后台
    end
```

#### 登录接口扩展

**现有登录接口修改（/api/v1/base/login）**

**请求参数扩展**：
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| username | string | 是 | 用户名或手机号 |
| password | string | 是 | 密码 |
| captcha | string | 否 | 验证码 |
| captchaId | string | 否 | 验证码ID |
| merchantId | uint | 否 | 选择的商户ID（多商户时必填） |

**响应数据结构扩展**：
```json
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "user": {
      "userName": "admin",
      "nickName": "管理员",
      "uuid": "xxx-xxx-xxx",
      "userId": 1,
      "authorityId": "super_admin",
      "merchantId": 1,
      "merchantName": "XX科技有限公司",
      "roleType": 2,
      "phone": "13800138000"
    },
    "token": "jwt_token_here",
    "expiresAt": "2024-01-01T10:00:00Z"
  }
}
```

#### 商户选择页面设计

**页面路由设计**
- 路由地址：在登录页面内部状态切换，不单独设置路由
- 页面状态：`loginStep: 'phone' | 'merchant-select' | 'password'`

**商户选择界面组件**
```vue
<template>
  <div class="merchant-select-container">
    <div class="select-header">
      <h3>选择工作商户</h3>
      <p>检测到您在多个商户中有账户，请选择要登录的商户</p>
    </div>
    
    <div class="merchant-list">
      <div 
        v-for="merchant in merchants" 
        :key="merchant.merchantId"
        class="merchant-item"
        :class="{ active: selectedMerchantId === merchant.merchantId }"
        @click="selectMerchant(merchant)"
      >
        <div class="merchant-icon">
          <img v-if="merchant.merchantIcon" :src="merchant.merchantIcon" :alt="merchant.merchantName" />
          <div v-else class="default-icon">{{ merchant.merchantName.charAt(0) }}</div>
        </div>
        <div class="merchant-info">
          <h4>{{ merchant.merchantName }}</h4>
          <span class="role-tag">{{ merchant.userRole }}</span>
        </div>
        <div class="select-indicator">
          <el-icon v-if="selectedMerchantId === merchant.merchantId"><Check /></el-icon>
        </div>
      </div>
    </div>
    
    <div class="action-buttons">
      <el-button @click="goBack">返回</el-button>
      <el-button type="primary" :disabled="!selectedMerchantId" @click="continueLogin">
        继续登录
      </el-button>
    </div>
  </div>
</template>
```

**交互逻辑**：
1. 用户输入手机号后，自动检测多商户身份
2. 如果是多商户，显示商户选择界面
3. 用户选择商户后，显示密码输入框
4. 提交登录时携带选中的商户ID
5. 支持返回重新选择商户

#### 前端状态管理

**登录状态管理扩展**
```javascript
// stores/auth.js
export const useAuthStore = defineStore('auth', {
  state: () => ({
    loginStep: 'phone', // 'phone' | 'merchant-select' | 'password'
    phoneNumber: '',
    isMultiMerchant: false,
    availableMerchants: [],
    selectedMerchantId: null,
    selectedMerchant: null
  }),
  
  actions: {
    // 检测多商户身份
    async checkMerchants(phone) {
      const response = await checkMerchants({ phone })
      this.phoneNumber = phone
      this.isMultiMerchant = response.data.isMultiMerchant
      this.availableMerchants = response.data.merchants || []
      
      if (this.isMultiMerchant) {
        this.loginStep = 'merchant-select'
      } else {
        this.loginStep = 'password'
        // 单商户时自动设置商户ID
        if (this.availableMerchants.length > 0) {
          this.selectedMerchantId = this.availableMerchants[0].merchantId
          this.selectedMerchant = this.availableMerchants[0]
        }
      }
    },
    
    // 选择商户
    selectMerchant(merchant) {
      this.selectedMerchantId = merchant.merchantId
      this.selectedMerchant = merchant
      this.loginStep = 'password'
    },
    
    // 登录
    async login(password) {
      const loginData = {
        username: this.phoneNumber,
        password,
        merchantId: this.selectedMerchantId
      }
      
      const response = await login(loginData)
      // 处理登录成功逻辑
      return response
    },
    
    // 重置状态
    resetLoginState() {
      this.loginStep = 'phone'
      this.phoneNumber = ''
      this.isMultiMerchant = false
      this.availableMerchants = []
      this.selectedMerchantId = null
      this.selectedMerchant = null
    }
  }
})
```

#### 兼容性保障措施

**1. 现有用户体验保持不变**
- 单商户员工登录流程保持不变，直接进入后台系统
- 登录界面UI保持原有设计，用户无感知
- 现有的自动登录、图片验证码、记住密码等功能正常工作

**2. 多商户功能透明升级**
- 对现有单商户用户无任何影响
- 多商户功能仅在需要时才会显示，不会干扰现有用户
- 支持渐进式升级，新增多商户员工时自动生效

**3. 错误处理和降级**
- 商户检测接口失败时，默认使用单商户模式
- 商户状态关闭时提供明确的错误提示
- 网络异常时，保留基础的登录功能




    merchantCount: 0,
    isMultiMerchant: false,
    userMerchants: []
  }),
  
  actions: {
    // 登录方法增强
    async login(loginForm) {
      const response = await login(loginForm)
      
      // 保存用户基础信息
      this.userInfo = response.data.user
      this.token = response.data.token
      this.merchantCount = response.data.user.merchantCount
      this.isMultiMerchant = response.data.user.isMultiMerchant
      
      // 设置Token
      setToken(response.data.token)
      
      // 如果是多商户用户，获取商户列表
      if (this.isMultiMerchant) {
        await this.fetchUserMerchants()
      }
      
      return response.data
    },
    
    // 获取用户商户列表
    async fetchUserMerchants() {
      const response = await getUserMerchants()
      this.userMerchants = response.data.merchants
      
      
})
```

#### 兼容性保障措施

**1. 现有用户体验保持不变**
- 现有单商户用户登录后直接进入后台，无需额外选择
- 登录界面UI保持原有设计，用户无感知
- 现有的自动登录、记住密码等功能正常工作

**2. 渐进式升级支持**
- 支持用户逐步被添加到多个商户
- 首次成为多商户用户时，引导用户了解商户选择功能
#### 兼容性保障措施

**1. 现有用户体验保持不变**
- 员工登录后直接进入所属商户的后台系统
- 登录界面UI保持原有设计，用户无感知
- 现有的自动登录、记住密码等功能正常工作

**2. 错误处理和降级**
- 商户状态关闭时提供明确的错误提示
- 网络异常时，保留基础的登录功能



## 多租户前端组件架构

### 商户层级管理组件结构

```mermaid
graph TB
    subgraph "商户管理主页面"
        A[MerchantManagement.vue - 商户管理主页]
        B[MerchantList.vue - 商户列表]
        C[MerchantTree.vue - 商户树形结构]
        D[MerchantForm.vue - 商户表单]
        E[MerchantDetail.vue - 商户详情]
    end
    
    subgraph "商户层级管理"
        F[MerchantHierarchy.vue - 层级结构管理]
        G[MerchantMove.vue - 商户移动组件]
        H[HierarchyBreadcrumb.vue - 层级面包屑]
    end
    

    
    subgraph "商户角色管理"
        K[merchantRoleList.vue - 商户角色列表]
        L[merchantRoleForm.vue - 角色表单]
        M[RolePermission.vue - 角色权限配置]
    end
    
    subgraph "通用组件"
        Q[merchantContext.vue - 商户上下文组件]
        R[DataIsolation.vue - 数据隔离组件]
        S[MerchantBadge.vue - 商户标记组件]
        T[TreeSelect.vue - 树形选择器]
    end
    
    A --> C
    B --> D
    C --> D
    C --> E
    C --> F
    C --> G
    C --> H
    D --> K
    D --> N
    E --> I
    E --> J
    F --> T
    H --> I
    H --> J
    K --> L
    K --> M
    N --> O
    N --> P
    
    Q --> C
    Q --> K
    Q --> N
    R --> D
    R --> K
    R --> N
    T --> F
    T --> I
```

### 商户层级组件设计说明

#### 商户树形结构组件（MerchantTree.vue）
- **功能职责**：以树形结构展示商户层级关系
- **交互功能**：支持节点展开/折叠、拖拽排序、右键菜单
- **操作支持**：新增子商户、编辑商户、移动商户、删除商户
- **数据加载**：支持懒加载和全量加载两种模式

#### 层级结构管理组件（MerchantHierarchy.vue）
- **功能职责**：提供商户层级结构的全面管理界面
- **视图切换**：支持树形视图、列表视图、组织架构图等多种展示方式
- **搜索过滤**：支持按商户名称、层级、状态等条件进行过滤
- **批量操作**：支持批量更改状态、批量移动等操作

#### 商户移动组件（MerchantMove.vue）
- **功能职责**：处理商户在层级结构中的移动操作
- **选择器**：提供树形选择器选择新的父商户
- **校验机制**：防止循环引用、防止移动到子节点
- **影响预览**：显示移动后将影响的所有子孙节点

#### 层级面包屑组件（HierarchyBreadcrumb.vue）
- **功能职责**：展示当前商户在层级结构中的位置
- **导航功能**：点击面包屑中的任意层级可快速导航
- **层级信息**：显示每个层级的商户名称和状态
- **响应式设计**：在移动设备上自动折叠显示

#### 树形选择器组件（TreeSelect.vue）
- **功能职责**：为表单提供层级商户选择功能
- **选择模式**：支持单选、多选、级联选择等模式
- **搜索功能**：支持按商户名称进行实时搜索过滤
- **数据限制**：可配置只显示特定层级或状态的商户### 多租户状态管理设计

#### Pinia Store 扩展结构
``javascript
// stores/merchant.js - 多租户状态管理
export const usemerchantStore = defineStore('merchant', {
  state: () => ({
    // 当前商户信息
    currentMerchant: null,
    // 用户所属商户列表
    userMerchants: [],
    // 商户树结构
    merchantTree: [],
    // 商户员工列表
    merchantUsers: [],
    // 商户角色列表
    merchantRoles: [],
    // 搜索条件
    searchConditions: {},
    // 分页信息
    pagination: {
      page: 1,
      pageSize: 10,
      total: 0
    }
  }),
  
  getters: {
    // 当前商户ID
    currentMerchantId: (state) => state.currentMerchant?.merchantId,
    // 当前商户名称
    currentMerchantName: (state) => state.currentMerchant?.merchantName,
    // 是否多商户用户
    isMultimerchant: (state) => state.userMerchants.length > 1,
    // 当前用户在当前商户的角色
    currentUserRole: (state) => state.currentMerchant?.userRole
  },
  
  actions: {
    // 获取用户所属商户列表
    async fetchUserMerchants() {
      try {
        const response = await api.get('/api/v1/user/merchants')
        this.userMerchants = response.data.merchants
        return response.data
      } catch (error) {
        console.error('获取商户列表失败:', error)
        throw error
      }
    },
    
    // 切换商户
    async switchMerchant(merchantId) {
      try {
        const response = await api.post('/api/v1/user/switch-merchant', {
          merchantId
        })
        
        // 更新Token
        const newToken = response.data.token
        localStorage.setItem('token', newToken)
        
        // 更新当前商户信息
        this.currentMerchant = response.data.merchantInfo
        
        // 重新初始化权限信息
        await useUserStore().getUserInfo()
        
        return response.data
      } catch (error) {
        console.error('切换商户失败:', error)
        throw error
      }
    },
    
    // 获取商户列表（管理员）
    async fetchMerchantList(params) {
      try {
        const response = await api.get('/api/v1/merchant/list', { params })
        this.merchantList = response.data.list
        this.pagination = response.data.pagination
        return response.data
      } catch (error) {
        console.error('获取商户列表失败:', error)
        throw error
      }
    },
    
    // 获取商户树形结构
    async fetchMerchantTree(params) {
      try {
        const response = await api.get('/api/v1/merchant/tree', { params })
        this.merchantTree = response.data.merchantTree
        return response.data
      } catch (error) {
        console.error('获取商户树结构失败:', error)
        throw error
      }
    },
    
    // 创建商户（支持层级）
    async createMerchant(merchantData) {
      try {
        const response = await api.post('/api/v1/merchant', merchantData)
        // 刷新列表和树结构
        await Promise.all([
          this.fetchMerchantList(this.searchConditions),
          this.fetchMerchantTree()
        ])
        return response.data
      } catch (error) {
        console.error('创建商户失败:', error)
        throw error
      }
    },
    
    // 移动商户位置
    async moveMerchant(merchantId, newParentId) {
      try {
        const response = await api.put(`/api/v1/merchant/${merchantId}/move`, {
          newParentId
        })
        // 刷新列表和树结构
        await Promise.all([
          this.fetchMerchantList(this.searchConditions),
          this.fetchMerchantTree()
        ])
        return response.data
      } catch (error) {
        console.error('移动商户失败:', error)
        throw error
      }
    },
    
    // 获取商户员工列表
    async fetchmerchantUsers(params) {
      try {
        const response = await api.get('/api/v1/merchant/user/list', { params })
        this.merchantUsers = response.data.list
        return response.data
      } catch (error) {
        console.error('获取商户员工列表失败:', error)
        throw error
      }
    },
    
    // 创建商户员工
    async createmerchantUser(userData) {
      try {
        const response = await api.post('/api/v1/merchant/user', userData)
        await this.fetchmerchantUsers()
        return response.data
      } catch (error) {
        console.error('创建商户员工失败:', error)
        throw error
      }
    },
    
    // 获取商户角色列表
    async fetchmerchantRoles(params) {
      try {
        const response = await api.get('/api/v1/merchant/authority/list', { params })
        this.merchantRoles = response.data.list
        return response.data
      } catch (error) {
        console.error('获取商户角色列表失败:', error)
        throw error
      }
    },
    
    // 创建商户角色
    async createmerchantRole(roleData) {
      try {
        const response = await api.post('/api/v1/merchant/authority', roleData)
        await this.fetchmerchantRoles()
        return response.data
      } catch (error) {
        console.error('创建商户角色失败:', error)
        throw error
      }
    }
  }
})
```

## 多租户权限控制设计

### Casbin 多租户权限策略扩展

#### 扩展权限模型
基于现有的 Casbin RBAC 体系，设计三级权限体系：

**原有模型**：`p = sub, obj, act`
**扩展模型**：`p = sub, obj, act, merchant`

其中：
- `sub`：主体（用户角色）
- `obj`：对象（API路径或菜单资源）
- `act`：操作（HTTP方法或菜单操作）
- `merchant`：租户（商户ID，超级管理员为*）

### 三级权限体系详细设计

#### 1. 超级管理员（Super Admin）
- **角色标识**：MerchantID = 1, RoleType = 1
- **角色特性**：属于默认商户但拥有跨商户权限，是系统最高权限角色
- **权限范围**：所有现有和未来开发的API接口和菜单权限
- **数据访问**：可访问所有商户的数据，通过中间件跳过多租户限制
- **菜单权限**：拥有所有系统菜单和商户菜单的完整访问权限
- **API权限**：拥有所有以'system:'和'merchant:'开头的权限代码
- **核心职责**：
  - 商户生命周期管理：创建、编辑、删除、状态管理
  - 商户权限分配：为商户管理员分配具体的功能权限
  - 商户层级管理：调整商户父子关系和层级结构
  - 系统配置：管理系统级配置和参数
  - 跨商户数据分析：查看和分析所有商户数据

#### 2. 商户管理员（Merchant Admin）
- **角色标识**：MerchantID = 具体值, RoleType = 2
- **角色特性**：受商户限制，只能管理所属商户，不能访问其他商户数据
- **权限范围**：由超级管理员手动分配，无需商户审核流程
- **数据访问**：严格限制在所属商户范围内，数据隔离中间件强制执行
- **菜单权限**：只能访问被超级管理员分配的商户管理相关菜单
- **API权限**：只能调用被分配的以'merchant:'开头的权限代码对应的API
- **权限分配限制**：只能在已拥有的权限范围内为员工分配权限
- **核心职责**：
  - 商户信息维护：更新商户基本信息、联系方式等
  - 员工全生命周期管理：创建、编辑、删除、状态管理
  - 角色权限管理：在权限范围内创建角色并分配权限
  - 商户内部数据分析：查看商户内部运营数据

#### 3. 商户自定义角色（Custom Employee Role）
- **角色标识**：MerchantID = 具体值, RoleType = 3
- **角色特性**：受商户和具体角色双重限制
- **权限范围**：由商户管理员或超级管理员在权限范围内手动分配
- **数据访问**：只能访问所属商户的特定业务数据
- **菜单权限**：只能访问被分配的具体功能菜单
- **API权限**：只能调用被分配的具体权限代码对应的API
- **角色细分**：
  - 商户操作员（merchant_operator）：日常业务操作权限
  - 商户查看员（merchant_viewer）：只读权限，数据查看
  - 商户会计（merchant_accountant）：财务相关数据权限
  - 商户客服（merchant_service）：客户服务相关权限
  - 其他定制角色：根据商户具体业务需要定义

#### 权限继承与限制原则

```mermaid
graph TB
    subgraph "超级管理员权限"
        A1[系统所有权限]
        A2[所有商户数据访问]
        A3[商户权限分配权]
        A4[系统配置权限]
    end
    
    subgraph "商户管理员权限"
        B1[商户内所有权限]
        B2[员工权限分配权]
        B3[商户信息管理权]
        B4[商户内部数据管理]
    end
    
    subgraph "员工权限"
        C1[特定业务权限]
        C2[只读查看权限]
        C3[部分操作权限]
    end
    
    A1 --> B1
    A3 --> B2
    B1 --> C1
    B2 --> C3
    
    style A1 fill:#ff9999
    style B1 fill:#99ccff
    style C1 fill:#99ff99
```

**权限分配原则**：
1. **向下分配原则**：上级角色只能将自己拥有的权限分配给下级
2. **商户隔离原则**：商户管理员不能访问其他商户（包括子商户）数据
3. **最小权限原则**：员工角色默认只分配必要的最小权限
4. **审批流程简化**：商户管理员分配权限无需审核，但有日志记录

#### 菜单权限配置系统

为三级权限体系设计分层菜单配置：

**超级管理员菜单配置**：
| 菜单名称 | 菜单路径 | 菜单级别 | 权限代码 | 菜单显示条件 | 说明 |
|-----------|----------|-----------|----------|--------------|------|
| 系统管理 | /system | 1 | system:* | 始终显示 | 一级菜单 |
| 商户管理 | /system/merchant | 2 | system:merchant:list | 始终显示 | 商户列表管理 |
| 商户创建 | /system/merchant/create | 3 | system:merchant:create | 始终显示 | 创建商户 |
| 商户权限配置 | /system/merchant/permission | 3 | system:merchant:permission | 始终显示 | 权限分配 |
| 用户管理 | /system/user | 2 | system:user:* | 始终显示 | 系统用户管理 |
| 角色管理 | /system/authority | 2 | system:authority:* | 始终显示 | 系统角色管理 |
| 菜单管理 | /system/menu | 2 | system:menu:* | 始终显示 | 系统菜单管理 |
| 系统配置 | /system/config | 2 | system:config:* | 始终显示 | 系统参数配置 |
| 日志管理 | /system/logs | 2 | system:logs:* | 始终显示 | 系统日志查看 |

**商户管理员菜单配置**：
| 菜单名称 | 菜单路径 | 菜单级别 | 权限代码 | 菜单显示条件 | 说明 |
|-----------|----------|-----------|----------|--------------|------|
| 商户管理 | /merchant | 1 | merchant:* | 商户上下文存在 | 一级菜单 |
| 商户信息 | /merchant/info | 2 | merchant:info:view | 可分配 | 商户基本信息 |
| 信息编辑 | /merchant/info/edit | 3 | merchant:info:update | 可分配 | 编辑商户信息 |
| 员工管理 | /merchant/user | 2 | merchant:user:list | 可分配 | 商户员工管理 |
| 员工创建 | /merchant/user/create | 3 | merchant:user:create | 可分配 | 新增员工 |
| 员工编辑 | /merchant/user/edit | 3 | merchant:user:update | 可分配 | 编辑员工 |
| 角色管理 | /merchant/role | 2 | merchant:role:list | 可分配 | 商户角色管理 |
| 角色创建 | /merchant/role/create | 3 | merchant:role:create | 可分配 | 新增角色 |
| 权限分配 | /merchant/role/permission | 3 | merchant:role:permission | 可分配 | 角色权限配置 |
| 数据统计 | /merchant/statistics | 2 | merchant:statistics:view | 可分配 | 商户数据统计 |

**员工角色菜单配置**（根据具体分配权限显示）：
| 菜单名称 | 菜单路径 | 菜单级别 | 权限代码 | 菜单显示条件 | 说明 |
|-----------|----------|-----------|----------|--------------|------|
| 工作台 | /workbench | 1 | * | 始终显示 | 员工工作台 |
| 信息查看 | /workbench/info | 2 | merchant:info:view | 有相应权限 | 查看商户信息 |
| 员工列表 | /workbench/users | 2 | merchant:user:view | 有相应权限 | 查看员工列表 |
| 业务操作 | /workbench/business | 2 | business:* | 有相应权限 | 具体业务操作 |
| 个人中心 | /profile | 1 | * | 始终显示 | 个人信息管理 |

#### 菜单动态加载机制

基于用户角色和权限的动态菜单加载机制：

```mermaid
flowchart TD
    A[用户登录] --> B{检查用户角色}
    B -->|super_admin| C[加载所有系统菜单]
    B -->|merchant_admin| D[检查商户上下文]
    B -->|employee_*| E[检查具体权限]
    
    C --> F[显示完整菜单结构]
    
    D --> G{商户上下文存在？}
    G -->|存在| H[获取分配的商户权限]
    G -->|不存在| I[跳转商户选择页]
    
    H --> J[按权限过滤商户菜单]
    E --> K[按具体权限过滤菜单]
    
    J --> L[显示可访问菜单]
    K --> L
    
    style A fill:#e1f5fe
    style F fill:#ff9999
    style L fill:#99ccff
    style I fill:#fff3e0
```

#### 三级权限验证流程设计

**权限验证层级结构**：
```mermaid
sequenceDiagram
    participant Client as 前端请求
    participant MW as 权限中间件
    participant Casbin as Casbin引擎
    participant DB as 数据库
    
    Client->>MW: 发起API请求
    MW->>MW: 解析JWT获取用户信息
    
    alt 超级管理员
        MW->>MW: 跳过权限检查
        MW->>Client: 允许访问
    else 商户管理员/员工
        MW->>Casbin: 检查权限策略
        Casbin->>DB: 查询用户权限
        DB-->>Casbin: 返回权限数据
        Casbin-->>MW: 返回权限检查结果
        
        alt 权限通过
            MW->>MW: 设置商户上下文
            MW->>Client: 允许访问
        else 权限不足
            MW->>Client: 返回403错误
        end
    end
```

**三级权限验证规则**：
1. **超级管理员验证**：authority_id = 'super_admin'，直接通过所有权限检查
2. **商户管理员验证**：检查是否拥有对应的'merchant:'权限，并验证商户上下文
3. **员工角色验证**：检查具体的权限代码，并进行双重验证（角色权限+商户权限）

#### 权限分配和管理机制

**超级管理员权限分配流程**：
```mermaid
flowchart TD
    A[超级管理员登录] --> B[选择商户]
    B --> C[查看可分配权限列表]
    C --> D[选择权限项]
    D --> E[为商户管理员分配权限]
    E --> F[权限生效，记录日志]
    
    style A fill:#ff9999
    style E fill:#99ccff
    style F fill:#99ff99
```

**商户管理员权限分配流程**：
```mermaid
flowchart TD
    A[商户管理员登录] --> B[进入员工管理]
    B --> C[查看已有权限列表]
    C --> D[创建/编辑员工角色]
    D --> E[在权限范围内分配权限]
    E --> F[分配给员工用户]
    F --> G[权限生效，记录日志]
    
    style A fill:#99ccff
    style E fill:#99ff99
    style G fill:#e8f5e8
```

#### 权限资源定义表
|----------|----------|---------|----------|------------|------------|----------|------|
| system:merchant:create | 商户创建 | /api/v1/merchant | POST | ✓ | ✗ | ✗ | 仅超级管理员 |
| system:merchant:list | 商户列表查看 | /api/v1/merchant/list | GET | ✓ | ✗ | ✗ | 查看所有商户 |
| system:merchant:update | 商户信息编辑 | /api/v1/merchant/:id | PUT | ✓ | ✗ | ✗ | 跨商户编辑 |
| system:merchant:delete | 商户删除 | /api/v1/merchant/:id | DELETE | ✓ | ✗ | ✗ | 删除任意商户 |
| system:merchant:status | 商户状态管理 | /api/v1/merchant/:id/status | PUT | ✓ | ✗ | ✗ | 状态管理 |
| system:merchant:permission | 商户权限分配 | /api/v1/merchant/:id/permissions | PUT | ✓ | ✗ | ✗ | 分配商户权限 |
| merchant:info:view | 商户信息查看 | /api/v1/merchant/info | GET | ✓ | 可分配 | 可分配 | 查看所属商户信息 |
| merchant:info:update | 商户信息编辑 | /api/v1/merchant/info | PUT | ✓ | 可分配 | 可分配 | 编辑所属商户信息 |
| merchant:user:list | 商户员工查看 | /api/v1/merchant/user/list | GET | ✓ | 可分配 | 可分配 | 查看商户员工 |
| merchant:user:create | 商户员工创建 | /api/v1/merchant/user | POST | ✓ | 可分配 | 可分配 | 创建商户员工 |
| merchant:user:update | 商户员工编辑 | /api/v1/merchant/user/:id | PUT | ✓ | 可分配 | 可分配 | 编辑商户员工 |
| merchant:user:delete | 商户员工删除 | /api/v1/merchant/user/:id | DELETE | ✓ | 可分配 | ✗ | 删除商户员工 |
| merchant:role:list | 商户角色查看 | /api/v1/merchant/role/list | GET | ✓ | 可分配 | 可分配 | 查看商户角色 |
| merchant:role:create | 商户角色创建 | /api/v1/merchant/role | POST | ✓ | 可分配 | ✗ | 创建商户角色 |
| merchant:role:update | 商户角色编辑 | /api/v1/merchant/role/:id | PUT | ✓ | 可分配 | ✗ | 编辑商户角色 |
| merchant:role:delete | 商户角色删除 | /api/v1/merchant/role/:id | DELETE | ✓ | 可分配 | ✗ | 删除商户角色 |

**权限分类说明**：
- **✓**：默认拥有该权限
- **可分配**：由上级角色手动分配决定是否拥有
- **✗**：不能拥有该权限

#### 三级权限体系分配策略

```mermaid
graph TD
    subgraph "超级管理员（Super Admin）"
        A1[拥有所有权限]
        A2[不受商户限制]
        A3[system:*]
        A4[merchant:*]
    end
    
    subgraph "商户管理员（Merchant Admin）"
        B1[受商户限制]
        B2[由超级管理员分配]
        B3[merchant:info:*]
        B4[merchant:user:*]
        B5[merchant:role:*]
    end
    
    subgraph "员工角色（Employee Role）"
        C1[受商户+角色限制]
        C2[由商户管理员/超级管理员分配]
        C3[merchant:info:view]
        C4[merchant:user:view]
        C5[其他业务权限]
    end
    
    A1 --> B1
    A1 --> C1
    B2 --> C2
    
    style A1 fill:#ff9999
    style B1 fill:#99ccff
    style C1 fill:#99ff99
```

### 三级权限中间件设计

#### 超级管理员权限中间件
``go
// middleware/super_admin.go
func SuperAdminMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        claims, exists := c.Get("claims")
        if !exists {
            response.FailWithMessage("认证失败", c)
            c.Abort()
            return
        }
        
        waitUse := claims.(*systemReq.CustomClaims)
        
        // 检查是否为超级管理员角色
        if waitUse.AuthorityId != "super_admin" {
            response.FailWithCodeMessage(40301, "仅超级管理员可访问", c)
            c.Abort()
            return
        }
        
        // 设置超级管理员标识
        c.Set("isSuperAdmin", true)
        c.Set("bypassmerchantIsolation", true) // 绕过商户隔离
        
        c.Next()
    }
}
```

#### 商户权限校验中间件
``go
// middleware/merchant_permission.go
func MerchantPermissionMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        claims, exists := c.Get("claims")
        if !exists {
            response.FailWithMessage("认证失败", c)
            c.Abort()
            return
        }
        
        waitUse := claims.(*systemReq.CustomClaims)
        
        // 超级管理员绕过权限检查
        if waitUse.AuthorityId == "super_admin" {
            c.Set("isSuperAdmin", true)
            c.Next()
            return
        }
        
        // 获取当前用户的商户ID
        merchantID := waitUse.MerchantID
        if merchantID == 0 {
            response.FailWithCodeMessage(40302, "商户上下文缺失", c)
            c.Abort()
            return
        }
        
        // 设置商户上下文
        c.Set("merchantID", merchantID)
        c.Set("userRole", waitUse.AuthorityId)
        
        c.Next()
    }
}
```

### 前端三级权限控制实现

#### 三级权限指令系统
``javascript
// directive/role-auth.js
import { useUserStore } from '@/stores/user'
import { usemerchantStore } from '@/stores/merchant'

export default {
  mounted(el, binding) {
    const userStore = useUserStore()
    const merchantStore = usemerchantStore()
    
    const requiredPermission = binding.value
    const userRole = userStore.userInfo.authorityId
    const userPermissions = userStore.userInfo.permissions || []
    
    // 超级管理员拥有所有权限
    if (userRole === 'super_admin') {
      return
    }
    
    // 检查用户是否拥有该权限
    if (!userPermissions.includes(requiredPermission)) {
      el.style.display = 'none'
      return
    }
    
    // 商户维度权限检查
    const merchantId = merchantStore.currentMerchantId
    if (!merchantId && requiredPermission.startsWith('merchant:')) {
      el.style.display = 'none'
      return
    }
  },
  
  updated(el, binding) {
    this.mounted(el, binding)
  }
}
```

#### 按钮级权限控制实例
``vue
<template>
  <div class="permission-demo">
    <!-- 超级管理员专属按钮 -->
    <el-button 
      v-role-auth="'system:merchant:create'"
      type="primary" 
      @click="createMerchant">
      创建商户
    </el-button>
    
    <!-- 商户管理员可见按钮 -->
    <el-button 
      v-role-auth="'merchant:user:create'"
      type="success" 
      @click="createUser">
      新增员工
    </el-button>
    
    <!-- 员工角色可见按钮 -->
    <el-button 
      v-role-auth="'merchant:info:view'"
      type="info" 
      @click="viewInfo">
      查看信息
    </el-button>
  </div>
</template>

<script setup>
// 业务逻辑
</script>
```

#### 多租户路由配置
``javascript
// router/modules/merchant.js
const merchantRoutes = [
  {
    path: '/merchant-selector',
    name: 'MerchantSelector',
    component: () => import('@/views/merchant/MerchantSelector.vue'),
    meta: {
      title: '商户选择',
      hidden: true // 隐藏在菜单中
    }
  },
  {
    path: '/merchant',
    name: 'Merchant',
    component: () => import('@/views/layout/index.vue'),
    meta: {
      title: '商户管理',
      icon: 'merchant',
      roles: ['admin', 'platform_admin']
    },
    children: [
      {
        path: 'list',
        name: 'MerchantList',
        component: () => import('@/views/merchant/list.vue'),
        meta: {
          title: '商户列表',
          permission: 'merchant:list'
        }
      },
      {
        path: 'merchant-user',
        name: 'merchantUser',
        component: () => import('@/views/merchant/merchant-user.vue'),
        meta: {
          title: '商户员工',
          permission: 'merchant:user:list',
          requireMerchant: true // 需要商户上下文
        }
      },
      {
        path: 'merchant-role',
        name: 'merchantRole',
        component: () => import('@/views/merchant/merchant-role.vue'),
        meta: {
          title: '商户角色',
          permission: 'merchant:role:list',
          requireMerchant: true
        }
      }
    ]
  }
]

export default merchantRoutes
```

#### 按钮级权限控制
``vue
<template>
  <div class="merchant-user-list">
    <!-- 创建按钮 -->
    <el-button 
      v-auth="'merchant:user:create'"
      v-merchant-auth="currentMerchantId"
      type="primary" 
      @click="handleCreate">
      新建员工
    </el-button>
    
    <!-- 操作列 -->
    <el-table-column label="操作" width="200">
      <template #default="{ row }">
        <el-button 
          v-auth="'merchant:user:update'"
          v-merchant-auth="row.merchantId"
          size="small" 
          @click="handleEdit(row)">
          编辑
        </el-button>
        <el-button 
          v-auth="'merchant:user:delete'"
          v-merchant-auth="row.merchantId"
          size="small" 
          type="danger" 
          @click="handleDelete(row)">
          删除
        </el-button>
      </template>
    </el-table-column>
  </div>
</template>

<script setup>
import { usemerchantStore } from '@/stores/merchant'

const merchantStore = usemerchantStore()
const currentMerchantId = computed(() => merchantStore.currentMerchantId)
</script>
```

#### 多租户权限指令
``javascript
// directive/merchant-auth.js
import { usemerchantStore } from '@/stores/merchant'
import { useUserStore } from '@/stores/user'

export default {
  mounted(el, binding) {
    const merchantStore = usemerchantStore()
    const userStore = useUserStore()
    
    const requiredMerchantId = binding.value
    const currentMerchantId = merchantStore.currentMerchantId
    const userRole = userStore.userInfo.authorityId
    
    // 系统管理员不受商户限制
    if (userRole === 'admin' || userRole === 'platform_admin') {
      return
    }
    
    // 检查商户权限
    if (requiredMerchantId && requiredMerchantId !== currentMerchantId) {
      el.style.display = 'none'
      // 或者移除元素
      // el.parentNode?.removeChild(el)
    }
  },
  
  updated(el, binding) {
    // 商户切换时重新检查权限
    this.mounted(el, binding)
  }
}
```

## 三级权限体系实现细节

### 商户权限分配界面设计

**超级管理员权限分配界面**：
- 商户选择器：显示所有商户列表
- 权限分类显示：按功能模块分类展示可分配权限
- 批量分配：支持选中多个权限同时分配
- 权限预览：分配前预览将要分配的权限列表
- 分配日志：记录所有权限分配操作的历史

**商户管理员权限管理界面**：
- 已有权限显示：显示当前商户已获得的所有权限
- 员工角色管理：在权限范围内创建和管理员工角色
- 权限分配限制：只能在已拥有的权限范围内分配给员工
- 权限使用统计：查看各种权限的使用情况和频次

### 前端权限组件设计

**权限控制组件集**：

1. **RoleAuthButton.vue** - 权限按钮组件
2. **merchantSwitcher.vue** - 商户切换组件
3. **PermissionTable.vue** - 权限表格组件
4. **MerchantPermissionAssign.vue** - 商户权限分配组件
5. **RolePermissionMatrix.vue** - 角色权限矩阵组件

### 权限检查流程优化

**高效权限检查策略**：
- 缓存权限数据：在前端缓存用户权限信息，减少服务器请求
- 权限预计算：用户登录时预计算所有权限，提高检查效率
- 懒加载验证：权限检查只在必要时执行，避免无效计算
- 权限变更通知：权限变更时实时清理相关缓存

### 数据隔离安全加强

**严格的数据隔离机制**：
- SQL注入防护：所有查询都使用参数化查询，防止SQL注入
- 行级安全：所有数据访问都必须包含merchant_id条件
- 接口参数校验：严格校验所有输入参数，防止越权访问
- 数据返回过滤：返回数据前再次过滤，确保数据安全

### 日志审计强化

**全面的审计日志系统**：
- 权限操作日志：记录所有权限分配、变更、删除操作
- 商户切换日志：记录用户的商户切换历史
- 敏感操作监控：重点监控跨商户操作和管理员权限变更
- 实时报警：异常权限操作实时报警通知

### 性能优化建议

**数据库优化**：
- 建立复合索引：(merchant_id, user_id)、(merchant_id, authority_id)等
- 分区表设计：按商户ID对大表进行分区，提高查询效率
- 读写分离：查询操作使用只读实例，减少主库压力
- 数据归档：定期归档历史数据，保持表的高效性

**前端优化**：
- 组件懒加载：大型组件按需加载，减少初始打包体积
- 数据虚拟化：大列表使用虚拟滚动，提高渲染性能
- 权限缓存：在localStorage/sessionStorage中缓存用户权限
- 防抖处理：搜索和筛选功能使用防抖，减少请求频次

### 系统集成建议

**与现有gin-vue-admin系统集成**：
- 保持现有API接口兼容：渐进式升级，不影响现有功能
- 统一认证系统：继续使用JWT认证，扩展支持商户上下文
- 数据库升级脚本：提供自动化脚本，安全升级现有数据库
- 配置参数管理：增加多租户相关的系统配置参数

### 系统测试策略

#### 单元测试设计

**后端单元测试**：

```go
// service/merchant_service_test.go
package service

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    "gorm.io/gorm"
)

type MerchantServiceTestSuite struct {
    suite.Suite
    db      *gorm.DB
    service *MerchantService
}

func (suite *MerchantServiceTestSuite) SetupTest() {
    // 初始化测试数据库
    suite.db = setupTestDB()
    suite.service = &MerchantService{db: suite.db}
}

func (suite *MerchantServiceTestSuite) TearDownTest() {
    // 清理测试数据
    cleanupTestDB(suite.db)
}

// 测试商户创建
func (suite *MerchantServiceTestSuite) TestCreateMerchant() {
    req := &request.CreateMerchantRequest{
        MerchantName:    "Test Merchant",
        MerchantType:    "ENTERPRISE",
        ContactName:     "Test Contact",
        ContactPhone:    "13800138000",
        ContactEmail:    "test@example.com",
        CreatedBy:       1,
    }
    
    merchant, err := suite.service.CreateMerchant(req)
    
    assert.NoError(suite.T(), err)
    assert.NotNil(suite.T(), merchant)
    assert.Equal(suite.T(), "Test Merchant", merchant.MerchantName)
    assert.Equal(suite.T(), 1, merchant.Level)
    assert.NotEmpty(suite.T(), merchant.MerchantCode)
}

// 测试商户数据隔离
func (suite *MerchantServiceTestSuite) TestDataIsolation() {
    // 创建两个商户
    merchant1 := suite.createTestMerchant("Merchant 1")
    merchant2 := suite.createTestMerchant("Merchant 2")
    
    // 为每个商户创建用户
    user1 := suite.createTestUser("user1", merchant1.ID)
    user2 := suite.createTestUser("user2", merchant2.ID)
    
    // 测试数据隔离
    users1, err1 := suite.service.GetMerchantUsers(merchant1.ID)
    users2, err2 := suite.service.GetMerchantUsers(merchant2.ID)
    
    assert.NoError(suite.T(), err1)
    assert.NoError(suite.T(), err2)
    assert.Len(suite.T(), users1, 1)
    assert.Len(suite.T(), users2, 1)
    assert.Equal(suite.T(), user1.ID, users1[0].ID)
    assert.Equal(suite.T(), user2.ID, users2[0].ID)
}

// 测试商户层级结构
func (suite *MerchantServiceTestSuite) TestMerchantHierarchy() {
    // 创建父商户
    parentReq := &request.CreateMerchantRequest{
        MerchantName: "Parent Merchant",
        CreatedBy:    1,
    }
    parent, err := suite.service.CreateMerchant(parentReq)
    assert.NoError(suite.T(), err)
    
    // 创建子商户
    childReq := &request.CreateMerchantRequest{
        MerchantName: "Child Merchant",
        ParentID:     &parent.ID,
        CreatedBy:    1,
    }
    child, err := suite.service.CreateMerchant(childReq)
    assert.NoError(suite.T(), err)
    
    // 验证层级结构
    assert.Equal(suite.T(), 1, parent.Level)
    assert.Equal(suite.T(), 2, child.Level)
    assert.Equal(suite.T(), fmt.Sprintf("%d", parent.ID), parent.Path)
    assert.Equal(suite.T(), fmt.Sprintf("%d/%d", parent.ID, child.ID), child.Path)
}

func TestMerchantServiceTestSuite(t *testing.T) {
    suite.Run(t, new(MerchantServiceTestSuite))
}
```

**前端单元测试**：

```javascript
// tests/unit/stores/merchant.test.js
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { usemerchantStore } from '@/stores/merchant'
import * as api from '@/api/merchant'

// Mock API
vi.mock('@/api/merchant')

describe('merchantStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  describe('fetchUserMerchants', () => {
    it('should fetch user merchants successfully', async () => {
      const mockMerchants = [
        { merchantId: 1, merchantName: 'Merchant 1', isDefault: true },
        { merchantId: 2, merchantName: 'Merchant 2', isDefault: false }
      ]
      
      api.getUserMerchants.mockResolvedValue({ data: { merchants: mockMerchants } })
      
      const store = usemerchantStore()
      const result = await store.fetchUserMerchants()
      
      expect(result.merchants).toEqual(mockMerchants)
      expect(store.userMerchants).toEqual(mockMerchants)
      expect(api.getUserMerchants).toHaveBeenCalledOnce()
    })

    it('should handle fetch error', async () => {
      api.getUserMerchants.mockRejectedValue(new Error('Network error'))
      
      const store = usemerchantStore()
      
      await expect(store.fetchUserMerchants()).rejects.toThrow('Network error')
    })
  })

  describe('switchMerchant', () => {
    it('should switch merchant successfully', async () => {
      const mockResponse = {
        data: {
          token: 'new-token',
          merchantInfo: {
            merchantId: 2,
            merchantName: 'Merchant 2',
            userRole: 'merchant_admin'
          }
        }
      }
      
      api.switchMerchant.mockResolvedValue(mockResponse)
      
      const store = usemerchantStore()
      const result = await store.switchMerchant(2)
      
      expect(result).toEqual(mockResponse.data)
      expect(store.currentMerchant).toEqual(mockResponse.data.merchantInfo)
      expect(localStorage.getItem('token')).toBe('new-token')
    })
  })
})
```

#### 集成测试设计

**API集成测试**：

```go
// test/integration/merchant_api_test.go
package integration

import (
    "bytes"
    "encoding/json"
    "net/http"
    "testing"
    "github.com/stretchr/testify/assert"
)

type MerchantAPITestSuite struct {
    suite.Suite
    server     *gin.Engine
    superToken string
    adminToken string
}

func (suite *MerchantAPITestSuite) SetupSuite() {
    // 初始化测试服务器
    suite.server = setupTestServer()
    
    // 获取测试Token
    suite.superToken = suite.getSuperAdminToken()
    suite.adminToken = suite.getMerchantAdminToken()
}

func (suite *MerchantAPITestSuite) TestCreateMerchant() {
    payload := map[string]interface{}{
        "merchantName":  "Test Merchant API",
        "merchantType":  "ENTERPRISE",
        "contactName":   "API Test",
        "contactPhone":  "13800138001",
        "contactEmail":  "api@test.com",
    }
    
    body, _ := json.Marshal(payload)
    req := httptest.NewRequest("POST", "/api/v1/merchant", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+suite.superToken)
    
    w := httptest.NewRecorder()
    suite.server.ServeHTTP(w, req)
    
    assert.Equal(suite.T(), http.StatusOK, w.Code)
    
    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), float64(0), response["code"])
}

func (suite *MerchantAPITestSuite) TestDataIsolationAPI() {
    // 创建两个商户及其用户
    merchant1ID := suite.createTestMerchant("Merchant 1")
    merchant2ID := suite.createTestMerchant("Merchant 2")
    
    user1Token := suite.createMerchantUser(merchant1ID, "user1")
    user2Token := suite.createMerchantUser(merchant2ID, "user2")
    
    // 测试用户1无法访问商户2的数据
    req := httptest.NewRequest("GET", "/api/v1/merchant/user/list", nil)
    req.Header.Set("Authorization", "Bearer "+user1Token)
    
    w := httptest.NewRecorder()
    suite.server.ServeHTTP(w, req)
    
    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    
    // 验证只能看到自己商户的数据
    data := response["data"].(map[string]interface{})
    users := data["list"].([]interface{})
    
    for _, user := range users {
        userMap := user.(map[string]interface{})
        assert.Equal(suite.T(), float64(merchant1ID), userMap["merchantId"])
    }
}
```

#### E2E测试设计

**Cypress E2E测试**：

```javascript
// cypress/e2e/merchant-management.cy.js
describe('Merchant Management E2E', () => {
  beforeEach(() => {
    // 清理数据库
    cy.task('db:seed')
    
    // 登录超级管理员
    cy.loginAsSuperAdmin()
  })

  describe('Merchant CRUD Operations', () => {
    it('should create a new merchant successfully', () => {
      cy.visit('/system/merchant')
      
      // 点击新增按钮
      cy.get('[data-cy="create-merchant-btn"]').click()
      
      // 填写表单
      cy.get('[data-cy="merchant-name"]').type('E2E Test Merchant')
      cy.get('[data-cy="merchant-type"]').select('ENTERPRISE')
      cy.get('[data-cy="contact-name"]').type('E2E Contact')
      cy.get('[data-cy="contact-phone"]').type('13800138888')
      cy.get('[data-cy="contact-email"]').type('e2e@test.com')
      
      // 提交表单
      cy.get('[data-cy="submit-btn"]').click()
      
      // 验证结果
      cy.get('.el-message--success').should('contain', '创建成功')
      cy.get('[data-cy="merchant-table"]').should('contain', 'E2E Test Merchant')
    })

    it('should handle merchant hierarchy correctly', () => {
      cy.visit('/system/merchant')
      
      // 切换到树形视图
      cy.get('[data-cy="tree-view-toggle"]').click()
      
      // 选择父节点，创建子商户
      cy.get('[data-cy="merchant-tree-node-1"]').rightclick()
      cy.get('[data-cy="add-child-merchant"]').click()
      
      // 填写子商户信息
      cy.get('[data-cy="merchant-name"]').type('Child Merchant')
      cy.get('[data-cy="submit-btn"]').click()
      
      // 验证层级结构
      cy.get('[data-cy="merchant-tree"]')
        .find('[data-level="2"]')
        .should('contain', 'Child Merchant')
    })
  })

  describe('Permission Assignment', () => {
    it('should assign permissions to merchant admin', () => {
      // 进入权限管理页面
      cy.visit('/system/merchant/permission')
      
      // 选择商户
      cy.get('[data-cy="merchant-selector"]').click()
      cy.get('.el-select-dropdown__item').first().click()
      
      // 选择权限
      cy.get('[data-cy="permission-merchant-info-view"]').check()
      cy.get('[data-cy="permission-merchant-user-list"]').check()
      cy.get('[data-cy="permission-merchant-user-create"]').check()
      
      // 提交权限分配
      cy.get('[data-cy="assign-permissions-btn"]').click()
      
      // 验证结果
      cy.get('.el-message--success').should('contain', '权限分配成功')
    })
  })

  describe('Multi-merchant Login Flow', () => {
    it('should handle multi-merchant user login', () => {
      // 登出超级管理员
      cy.logout()
      
      // 登录多商户用户
      cy.login('multiuser@test.com', 'password123')
      
      // 应该显示商户选择页面
      cy.url().should('include', '/merchant-selector')
      cy.get('[data-cy="merchant-list"]').should('be.visible')
      
      // 选择商户
      cy.get('[data-cy="merchant-item-1"]').click()
      
      // 验证跳转到对应页面
      cy.url().should('include', '/merchant')
      cy.get('[data-cy="current-merchant-name"]').should('contain', 'Test Merchant')
    })
  })

  describe('Data Isolation Verification', () => {
    it('should enforce strict data isolation', () => {
      // 登录商户A的管理员
      cy.loginAsMerchantAdmin('merchant-a@test.com', 'password123')
      
      // 尝试访问其他商户的数据
      cy.request({
        url: '/api/v1/merchant/user/list',
        headers: {
          'Authorization': `Bearer ${Cypress.env('merchantAToken')}`
        },
        failOnStatusCode: false
      }).then((response) => {
        // 验证只能看到自己商户的数据
        expect(response.body.data.list).to.have.length.greaterThan(0)
        response.body.data.list.forEach(user => {
          expect(user.merchantId).to.equal(1) // 商户A的ID
        })
      })
    })
  })
})
```

#### 性能测试设计

**JMeter性能测试计划**：

```xml
<!-- performance-test-plan.jmx -->
<?xml version="1.0" encoding="UTF-8"?>
<jmeterTestPlan version="1.2">
  <hashTree>
    <TestPlan testname="Merchant Management Performance Test">
      <!-- 线程组配置 -->
      <ThreadGroup testname="API Load Test">
        <elementProp name="ThreadGroup.arguments" elementType="Arguments" guiclass="ArgumentsPanel">
          <collectionProp name="Arguments.arguments">
            <elementProp name="baseUrl" elementType="Argument">
              <stringProp name="Argument.name">baseUrl</stringProp>
              <stringProp name="Argument.value">http://localhost:8888</stringProp>
            </elementProp>
          </collectionProp>
        </elementProp>
        
        <!-- 并发用户数：100 -->
        <stringProp name="ThreadGroup.num_threads">100</stringProp>
        <!-- 启动时间：30秒 -->
        <stringProp name="ThreadGroup.ramp_time">30</stringProp>
        <!-- 循环次数：10 -->
        <stringProp name="ThreadGroup.loops">10</stringProp>
        
        <hashTree>
          <!-- 登录接口测试 -->
          <HTTPSamplerProxy testname="Login API">
            <stringProp name="HTTPSampler.path">/api/v1/login</stringProp>
            <stringProp name="HTTPSampler.method">POST</stringProp>
            <stringProp name="HTTPSampler.postBodyRaw">{
              "username": "testuser${__threadNum}",
              "password": "password123"
            }</stringProp>
          </HTTPSamplerProxy>
          
          <!-- 商户列表接口测试 -->
          <HTTPSamplerProxy testname="Merchant List API">
            <stringProp name="HTTPSampler.path">/api/v1/merchant/list</stringProp>
            <stringProp name="HTTPSampler.method">GET</stringProp>
          </HTTPSamplerProxy>
          
          <!-- 商户创建接口测试 -->
          <HTTPSamplerProxy testname="Create Merchant API">
            <stringProp name="HTTPSampler.path">/api/v1/merchant</stringProp>
            <stringProp name="HTTPSampler.method">POST</stringProp>
            <stringProp name="HTTPSampler.postBodyRaw">{
              "merchantName": "Performance Test Merchant ${__time()}",
              "merchantType": "ENTERPRISE",
              "contactName": "Performance Test",
              "contactPhone": "13800${__Random(100000,999999)}",
              "contactEmail": "perf${__threadNum}@test.com"
            }</stringProp>
          </HTTPSamplerProxy>
          
          <!-- 响应时间断言 -->
          <ResponseAssertion testname="Response Time Assertion">
            <stringProp name="Assertion.test_field">Assertion.response_time</stringProp>
            <stringProp name="Assertion.test_type">Assertion.response_time</stringProp>
            <stringProp name="Assertion.assume_success">false</stringProp>
            <stringProp name="Assertion.response_time">2000</stringProp> <!-- 2秒内 -->
          </ResponseAssertion>
          
          <!-- 结果监听器 -->
          <ResultCollector testname="View Results Tree" enabled="true">
            <stringProp name="filename">performance-test-results.jtl</stringProp>
            <objProp>
              <name>saveConfig</name>
              <value class="SampleSaveConfiguration">
                <time>true</time>
                <latency>true</latency>
                <timestamp>true</timestamp>
                <success>true</success>
                <label>true</label>
                <code>true</code>
                <message>true</message>
                <threadName>true</threadName>
                <dataType>true</dataType>
                <encoding>false</encoding>
                <assertions>true</assertions>
                <subresults>true</subresults>
                <responseData>false</responseData>
                <samplerData>false</samplerData>
                <xml>false</xml>
                <fieldNames>true</fieldNames>
                <responseHeaders>false</responseHeaders>
                <requestHeaders>false</requestHeaders>
                <responseDataOnError>false</responseDataOnError>
                <saveAssertionResultsFailureMessage>true</saveAssertionResultsFailureMessage>
                <assertionsResultsToSave>0</assertionsResultsToSave>
                <bytes>true</bytes>
                <sentBytes>true</sentBytes>
                <url>true</url>
                <threadCounts>true</threadCounts>
                <idleTime>true</idleTime>
                <connectTime>true</connectTime>
              </value>
            </objProp>
          </ResultCollector>
        </hashTree>
      </ThreadGroup>
    </TestPlan>
  </hashTree>
</jmeterTestPlan>
```

**分阶段上线策略**：

**第一阶段**：基础多租户架构
- 商户数据模型创建
- 基础数据隔离中间件
- 商户管理CRUD接口

**第二阶段**：权限体系集成
- 三级权限体系实现
- Casbin权限模型扩展
- 商户切换功能

**第三阶段**：高级功能完善
- 商户层级管理
- 认证流程系统
- 日志审计系统

**第四阶段**：优化与监控
- 性能优化实施
- 监控报警系统
- 自动化运维工具

### 多租户文件存储策略
- **存储路径隔离**：按商户ID组织文件目录结构
- **访问权限控制**：只能访问所属商户的文件
- **存储方式集成**：兼容本地存储、阿里云OSS、腾讯云COS、MinIO等
- **文件安全**：支持文件加密和数字签名

### 多租户存储路径规划
```
uploads/
├── merchants/                    # 多租户文件根目录
│   ├── merchant_1/           # 商户ID目录
│   │   ├── business_license/   # 营业执照
│   │   ├── tax_registration/   # 税务登记证
│   │   ├── legal_person_id/    # 法人身份证
│   │   └── other_docs/         # 其他证件
│   └── merchant_2/
└── system/                   # 系统级文件（不受商户限制）
```

### 文件上传安全控制
- **文件类型限制**：只允许PDF、JPG、PNG、DOC、DOCX格式
- **文件大小限制**：单文件不超过10MB，单次上传总大小不超过50MB
- **病毒扫描**：集成病毒扫描引擎，拦截恶意文件
- **数字水印**：重要证件文件支持数字水印验证

## 数据安全与验证

### 多租户数据验证规则

#### 后端验证规则
- **商户名称**：2-100个字符，全平台唯一
- **联系电话**：符合中国大陆手机号格式
- **邮箱地址**：符合标准邮箱格式
- **员工手机号**：在同一商户内唯一
- **商户编码**：系统自动生成，全平台唯一
- **角色名称**：在同一商户内唯一

#### 前端验证规则
``javascript
const merchantValidationRules = {
  // 商户验证规则
  merchantName: [
    { required: true, message: '请输入商户名称', trigger: 'blur' },
    { min: 2, max: 100, message: '商户名称长度为2-100个字符', trigger: 'blur' },
    { 
      validator: async (rule, value) => {
        // 检查商户名称全平台唯一性
        const exists = await api.checkMerchantNameExists(value)
        if (exists) {
          throw new Error('商户名称已存在')
        }
      }, 
      trigger: 'blur' 
    }
  ],
  
  // 员工验证规则
  userPhone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' },
    { 
      validator: async (rule, value) => {
        const merchantStore = usemerchantStore()
        // 检查在当前商户内的唯一性
        const exists = await api.checkUserPhoneInmerchant(value, merchantStore.currentMerchantId)
        if (exists) {
          throw new Error('该手机号在当前商户内已存在')
        }
      }, 
      trigger: 'blur' 
    }
  ],
  
  // 角色验证规则
  authorityName: [
    { required: true, message: '请输入角色名称', trigger: 'blur' },
    { min: 2, max: 50, message: '角色名称长度为2-50个字符', trigger: 'blur' },
    { 
      validator: async (rule, value) => {
        const merchantStore = usemerchantStore()
        // 检查在当前商户内的唯一性
        const exists = await api.checkAuthorityNameInmerchant(value, merchantStore.currentMerchantId)
        if (exists) {
          throw new Error('该角色名称在当前商户内已存在')
        }
      }, 
      trigger: 'blur' 
    }
  ]
}
```

### 多租户数据脱敏处理
- **手机号脱敏**：显示为138****8000格式
- **邮箱脱敏**：显示为user***@example.com格式
- **身份证号脱敏**：显示为110***********1234格式
- **银行账号脱敏**：显示为622***********1234格式
- **商户间数据隔离**：严格防止跨商户数据泄露

## 日志审计设计

### 多租户操作日志记录
基于系统现有的操作审计中间件，扩展记录商户维度的操作日志：

- **商户创建**：记录创建人、创建时间、商户基本信息、初始管理员信息
- **商户修改**：记录修改人、所属商户、修改时间、修改前后数据对比
- **商户切换**：记录用户切换商户的操作日志和时间
- **员工管理**：记录商户员工的创建、修改、删除、角色变更操作
- **角色权限**：记录商户角色的创建、修改、权限变更操作
- **跨商户操作**：重点监控系统管理员的跨商户操作

### 多租户审计日志格式
``json
{
  "operation": "merchant_user_create",
  "operator_id": 1,
  "operator_name": "merchant_admin",
  "operator_merchant_id": 1,
  "operator_merchant_name": "XX科技有限公司",
  "target_merchant_id": 1,
  "target_user_id": 123,
  "details": {
    "user_name": "张三",
    "user_phone": "138****8000",
    "assigned_roles": ["merchant_operator"]
  },
  "data_isolation": {
    "is_cross_merchant": false,
    "source_merchant": 1,
    "target_merchant": 1
  },
  "timestamp": "2024-01-01T10:00:00Z",
  "ip_address": "192.168.1.100",
  "user_agent": "Mozilla/5.0..."
}
```

## 性能优化设计

### 多租户中间件详细实现

#### 数据隔离中间件增强版

```go
// middleware/merchant_isolation.go
func merchantIsolationMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        claims, exists := c.Get("claims")
        if !exists {
            response.FailWithMessage("认证失败", c)
            c.Abort()
            return
        }
        
        waitUse := claims.(*systemReq.CustomClaims)
        
        // 超级管理员绕过租户隔离
        if waitUse.AuthorityId == "super_admin" {
            c.Set("isSuperAdmin", true)
            c.Set("bypassmerchantIsolation", true)
            c.Next()
            return
        }
        
        // 检查商户上下文
        merchantID := waitUse.MerchantID
        if merchantID == 0 {
            response.FailWithCodeMessage(40302, "商户上下文缺失", c)
            c.Abort()
            return
        }
        
        // 验证商户状态
        merchantService := service.ServiceGroupApp.MerchantServiceGroup.MerchantService
        merchant, err := merchantService.GetMerchantByID(merchantID)
        if err != nil {
            response.FailWithCodeMessage(40303, "商户信息获取失败", c)
            c.Abort()
            return
        }
        
        if merchant.Status != "ACTIVE" {
            response.FailWithCodeMessage(40303, "商户状态异常", c)
            c.Abort()
            return
        }
        
        // 设置租户上下文
        c.Set("merchantID", merchantID)
        c.Set("merchantInfo", merchant)
        c.Set("userRole", waitUse.AuthorityId)
        
        c.Next()
    }
}
```

#### 权限验证中间件

```go
// middleware/permission_check.go
func PermissionCheckMiddleware(requiredPermission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 超级管理员直接通过
        if isSuperAdmin, exists := c.Get("isSuperAdmin"); exists && isSuperAdmin.(bool) {
            c.Next()
            return
        }
        
        claims, exists := c.Get("claims")
        if !exists {
            response.FailWithCodeMessage(403, "权限验证失败", c)
            c.Abort()
            return
        }
        
        waitUse := claims.(*systemReq.CustomClaims)
        
        // 使用Casbin进行权限检查
        casbinService := service.ServiceGroupApp.SystemServiceGroup.CasbinService
        success := casbinService.Enforce(waitUse.AuthorityId, requiredPermission, "execute")
        
        if !success {
            response.FailWithCodeMessage(403, "权限不足", c)
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```
### 商户数据服务层设计

#### 商户管理服务

```go
// service/merchant_service.go
type MerchantService struct {
    db *gorm.DB
}

// 创建商户（支持层级结构）
func (m *MerchantService) CreateMerchant(req *request.CreateMerchantRequest) (*model.SysMerchant, error) {
    // 事务处理
    tx := m.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    
    // 生成商户编码
    merchantCode := m.generateMerchantCode()
    
    // 计算层级信息
    level := 1
    path := ""
    if req.ParentID != nil {
        parentMerchant, err := m.GetMerchantByID(*req.ParentID)
        if err != nil {
            tx.Rollback()
            return nil, err
        }
        level = parentMerchant.Level + 1
        path = parentMerchant.Path
    }
    
    // 创建商户记录
    merchant := &model.SysMerchant{
        MerchantCode:      merchantCode,
        MerchantName:      req.MerchantName,
        ParentID:          req.ParentID,
        MerchantType:      req.MerchantType,
        Level:            level,
        ContactName:       req.ContactName,
        ContactPhone:      req.ContactPhone,
        ContactEmail:      req.ContactEmail,
        BusinessLicense:   req.BusinessLicense,
        LegalPerson:      req.LegalPerson,
        RegisteredAddress: req.RegisteredAddress,
        BusinessScope:    req.BusinessScope,
        Status:           "PENDING",
        MerchantLevel:    "BASIC",
        CreatedBy:        req.CreatedBy,
    }
    
    if err := tx.Create(merchant).Error; err != nil {
        tx.Rollback()
        return nil, err
    }
    
    // 更新路径信息
    if path != "" {
        merchant.Path = fmt.Sprintf("%s/%d", path, merchant.ID)
    } else {
        merchant.Path = fmt.Sprintf("%d", merchant.ID)
    }
    
    if err := tx.Save(merchant).Error; err != nil {
        tx.Rollback()
        return nil, err
    }
    
    // 创建商户管理员用户
    if req.AdminUserInfo != nil {
        adminUser := &model.SysUser{
            Username:     req.AdminUserInfo.Username,
            Password:     req.AdminUserInfo.Password, // 需要加密
            NickName:     req.AdminUserInfo.NickName,
            Name:         req.AdminUserInfo.Name,
            Phone:        req.AdminUserInfo.Phone,
            Email:        req.AdminUserInfo.Email,
            MerchantID:   merchant.ID,
            AuthorityId:  fmt.Sprintf("merchant_admin_%d", merchant.ID),
            IsMainAccount: true,
        }
        
        if err := tx.Create(adminUser).Error; err != nil {
            tx.Rollback()
            return nil, err
        }
        
        merchant.AdminUserID = &adminUser.ID
        if err := tx.Save(merchant).Error; err != nil {
            tx.Rollback()
            return nil, err
        }
    }
    
    tx.Commit()
    return merchant, nil
}

// 获取商户树形结构
func (m *MerchantService) GetMerchantTree(rootID *uint, maxLevel int) ([]*response.MerchantTreeNode, error) {
    var merchants []model.SysMerchant
    query := m.db.Model(&model.SysMerchant{})
    
    if rootID != nil {
        // 查询指定根节点及其子节点
        rootMerchant, err := m.GetMerchantByID(*rootID)
        if err != nil {
            return nil, err
        }
        query = query.Where("path LIKE ? OR id = ?", rootMerchant.Path+"/%", *rootID)
    }
    
    if maxLevel > 0 {
        query = query.Where("level <= ?", maxLevel)
    }
    
    query = query.Order("level, parent_id, id")
    
    if err := query.Find(&merchants).Error; err != nil {
        return nil, err
    }
    
    return m.buildMerchantTree(merchants, rootID), nil
}

// 移动商户位置
func (m *MerchantService) MoveMerchant(merchantID uint, newParentID *uint) error {
    tx := m.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    
    // 获取当前商户
    merchant, err := m.GetMerchantByID(merchantID)
    if err != nil {
        tx.Rollback()
        return err
    }
    
    // 验证是否形成循环引用
    if newParentID != nil {
        if err := m.validateMoveOperation(merchantID, *newParentID); err != nil {
            tx.Rollback()
            return err
        }
    }
    
    // 计算新的level和path
    newLevel := 1
    newPath := fmt.Sprintf("%d", merchantID)
    
    if newParentID != nil {
        parentMerchant, err := m.GetMerchantByID(*newParentID)
        if err != nil {
            tx.Rollback()
            return err
        }
        newLevel = parentMerchant.Level + 1
        newPath = fmt.Sprintf("%s/%d", parentMerchant.Path, merchantID)
    }
    
    // 更新当前商户
    oldPath := merchant.Path
    oldLevel := merchant.Level
    
    merchant.ParentID = newParentID
    merchant.Level = newLevel
    merchant.Path = newPath
    
    if err := tx.Save(merchant).Error; err != nil {
        tx.Rollback()
        return err
    }
    
    // 更新所有子孙节点
    if err := m.updateChildrenPaths(tx, oldPath, newPath, oldLevel, newLevel); err != nil {
        tx.Rollback()
        return err
    }
    
    tx.Commit()
    return nil
}
```
### 数据库迁移与初始化脚本

#### 数据库迁移脚本

```
-- 创建商户主表
CREATE TABLE `sys_merchant` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `merchant_code` varchar(50) NOT NULL COMMENT '商户编码',
  `merchant_name` varchar(100) NOT NULL COMMENT '商户名称',
  `parent_id` bigint(20) unsigned DEFAULT NULL COMMENT '父商户ID',
  `merchant_type` varchar(20) NOT NULL DEFAULT 'ENTERPRISE' COMMENT '商户类型',
  `level` int(11) NOT NULL DEFAULT '1' COMMENT '商户层级',
  `path` varchar(500) NOT NULL DEFAULT '' COMMENT '层级路径',
  `contact_name` varchar(50) NOT NULL COMMENT '联系人姓名',
  `contact_phone` varchar(20) NOT NULL COMMENT '联系电话',
  `contact_email` varchar(100) NOT NULL COMMENT '联系邮箱',
  `business_license` varchar(50) DEFAULT NULL COMMENT '营业执照号',
  `legal_person` varchar(50) DEFAULT NULL COMMENT '法人代表',
  `registered_address` varchar(255) DEFAULT NULL COMMENT '注册地址',
  `business_scope` text COMMENT '经营范围',

  `merchant_level` varchar(20) NOT NULL DEFAULT 'BASIC' COMMENT '商户等级',


  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_merchant_code` (`merchant_code`),
  KEY `idx_merchant_parent_id` (`parent_id`),
  KEY `idx_merchant_status` (`status`),
  KEY `idx_merchant_path` (`path`),
  KEY `idx_merchant_level` (`level`),
  KEY `idx_merchant_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商户信息表';

-- 扩展sys_user表，增加商户相关字段
ALTER TABLE `sys_user` 
ADD COLUMN `merchant_id` bigint(20) unsigned DEFAULT NULL COMMENT '所属商户ID' AFTER `authority_id`,
ADD COLUMN `is_main_account` tinyint(1) DEFAULT '0' COMMENT '是否为主账号' AFTER `merchant_id`,
ADD INDEX `idx_user_merchant_id` (`merchant_id`),
ADD INDEX `idx_user_phone_merchant` (`phone`, `merchant_id`, `deleted_at`);




-- 创建商户用户关联表
CREATE TABLE `sys_merchant_user` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` bigint(20) unsigned NOT NULL COMMENT '用户ID',
  `merchant_id` bigint(20) unsigned NOT NULL COMMENT '商户ID',
  `is_default` tinyint(1) DEFAULT '0' COMMENT '是否为默认商户',
  `joined_at` datetime(3) NOT NULL COMMENT '加入时间',
  `status` varchar(20) NOT NULL DEFAULT 'ACTIVE' COMMENT '关联状态',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_merchant_user_unique` (`user_id`, `merchant_id`),
  KEY `idx_merchant_user_merchant_id` (`merchant_id`),
  KEY `idx_merchant_user_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商户用户关联表';



-- 创建商户状态变更日志表
CREATE TABLE `sys_merchant_status_log` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `merchant_id` bigint(20) unsigned NOT NULL COMMENT '商户ID',
  `previous_status` varchar(20) NOT NULL COMMENT '变更前状态',
  `new_status` varchar(20) NOT NULL COMMENT '变更后状态',
  `change_reason` varchar(255) NOT NULL COMMENT '变更原因',
  `operator_id` bigint(20) unsigned NOT NULL COMMENT '操作人 ID',
  `operator_merchant_id` bigint(20) unsigned DEFAULT NULL COMMENT '操作人所属商户ID',
  `created_at` datetime(3) DEFAULT NULL COMMENT '变更时间',
  PRIMARY KEY (`id`),
  KEY `idx_merchant_status_log_merchant_id` (`merchant_id`),
  KEY `idx_merchant_status_log_operator_id` (`operator_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商户状态变更日志表';
```

#### 数据初始化脚本

```
-- 初始化商户权限数据
INSERT INTO `sys_api` (`path`, `description`, `api_group`, `method`) VALUES 
('/api/v1/merchant', '创建商户', '商户管理', 'POST'),
('/api/v1/merchant/list', '商户列表', '商户管理', 'GET'),
('/api/v1/merchant/:id', '更新商户', '商户管理', 'PUT'),
('/api/v1/merchant/:id', '删除商户', '商户管理', 'DELETE'),
('/api/v1/merchant/:id/status', '更新商户状态', '商户管理', 'PUT'),
('/api/v1/merchant/:id/permissions', '分配商户权限', '商户管理', 'PUT'),
('/api/v1/merchant/:id/move', '移动商户', '商户管理', 'PUT'),
('/api/v1/merchant/tree', '商户树结构', '商户管理', 'GET'),
('/api/v1/merchant/info', '商户信息', '商户管理', 'GET'),
('/api/v1/merchant/user', '商户员工管理', '商户管理', 'POST'),
('/api/v1/merchant/user/list', '商户员工列表', '商户管理', 'GET'),
('/api/v1/merchant/user/:id', '更新商户员工', '商户管理', 'PUT'),
('/api/v1/merchant/user/:id', '删除商户员工', '商户管理', 'DELETE'),
('/api/v1/user/merchants', '用户商户列表', '用户管理', 'GET'),
('/api/v1/user/switch-merchant', '切换商户', '用户管理', 'POST');

-- 初始化菜单数据
INSERT INTO `sys_base_menu` (`name`, `path`, `hidden`, `component`, `sort`, `meta`, `parent_id`) VALUES 
('系统管理', '/system', 0, 'view/layout/index.vue', 1, '{"icon":"setting","title":"系统管理"}', 0),
('商户管理', '/system/merchant', 0, 'view/merchant/index.vue', 10, '{"icon":"shop","title":"商户管理"}', 1),
('商户列表', '/system/merchant/list', 0, 'view/merchant/list.vue', 1, '{"icon":"list","title":"商户列表"}', 2),
('商户权限', '/system/merchant/permission', 0, 'view/merchant/permission.vue', 2, '{"icon":"key","title":"商户权限"}', 2);

-- 初始化Casbin权限策略
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES 
('p', 'super_admin', '/api/v1/merchant', 'POST'),
('p', 'super_admin', '/api/v1/merchant/list', 'GET'),
('p', 'super_admin', '/api/v1/merchant/:id', 'PUT'),
('p', 'super_admin', '/api/v1/merchant/:id', 'DELETE'),
('p', 'super_admin', '/api/v1/merchant/:id/status', 'PUT'),
('p', 'super_admin', '/api/v1/merchant/:id/permissions', 'PUT'),
('p', 'super_admin', '/api/v1/merchant/:id/move', 'PUT');
```
### 前端完整组件实现

#### 商户选择器组件

```
<!-- components/MerchantSelector.vue -->
<template>
  <div class="merchant-selector">
    <el-card class="selector-card">
      <template #header>
        <div class="card-header">
          <span>选择商户</span>
          <el-button type="text" @click="refreshMerchants">
            <el-icon><Refresh /></el-icon>
          </el-button>
        </div>
      </template>
      
      <div class="merchant-list">
        <el-empty v-if="merchants.length === 0" description="暂无可用商户" />
        
        <div 
          v-for="merchant in merchants" 
          :key="merchant.merchantId"
          class="merchant-item"
          :class="{ 'is-default': merchant.isDefault }"
          @click="selectMerchant(merchant)">
          
          <div class="merchant-info">
            <div class="merchant-name">{{ merchant.merchantName }}</div>
            <div class="merchant-code">{{ merchant.merchantCode }}</div>
            <div class="merchant-status">
              <el-tag :type="getStatusType(merchant.status)">{{ getStatusText(merchant.status) }}</el-tag>
```

## 现有数据处理与数据库迁移

### 现有系统数据处理策略

为了保证与现有gin-vue-admin系统的兼容性，在实施多租户改造时，需要特别处理现有数据：

#### 数据迁移原则

1. **默认商户创建**：创建一个ID为1的默认商户，作为现有数据的归属商户
2. **员工数据迁移**：所有现有员工设置为商户ID=1（超级管理员除外）
3. **角色数据迁移**：所有现有角色设置为商户ID=1（系统角色除外）
4. **超级管理员特殊处理**：超级管理员不属于任何商户，MerchantID为NULL
5. **主账号标识**：为默认商户指定主管理员账号

#### 完整数据库迁移脚本

```
-- ====================================
-- 商户管理系统数据库迁移脚本
-- 适用于现有gin-vue-admin系统升级为多租户架构
-- ====================================

-- 第一步：创建商户主表
CREATE TABLE `sys_merchant` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `merchant_code` varchar(50) NOT NULL COMMENT '商户编码',
  `merchant_name` varchar(100) NOT NULL COMMENT '商户名称',
  `merchant_icon` varchar(255) DEFAULT NULL COMMENT '商户图标URL',
  `parent_id` bigint(20) unsigned DEFAULT NULL COMMENT '父商户ID',
  `merchant_type` varchar(20) NOT NULL DEFAULT 'ENTERPRISE' COMMENT '商户类型',
  `level` int(11) NOT NULL DEFAULT '1' COMMENT '商户层级',
  `path` varchar(500) NOT NULL DEFAULT '' COMMENT '层级路径',
  `contact_name` varchar(50) NOT NULL COMMENT '联系人姓名',
  `contact_phone` varchar(20) NOT NULL COMMENT '联系电话',
  `contact_email` varchar(100) NOT NULL COMMENT '联系邮箱',
  `business_license` varchar(50) DEFAULT NULL COMMENT '营业执照号',
  `legal_person` varchar(50) DEFAULT NULL COMMENT '法人代表',
  `registered_address` varchar(255) DEFAULT NULL COMMENT '注册地址',
  `business_scope` text COMMENT '经营范围',

  `merchant_level` varchar(20) NOT NULL DEFAULT 'BASIC' COMMENT '商户等级',

  `operator_id` bigint(20) unsigned NOT NULL COMMENT '操作者用户ID',
  `operator_name` varchar(50) NOT NULL COMMENT '操作者姓名',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_merchant_code` (`merchant_code`),
  KEY `idx_merchant_parent_id` (`parent_id`),
  KEY `idx_merchant_status` (`status`),
  KEY `idx_merchant_path` (`path`),
  KEY `idx_merchant_level` (`level`),
  KEY `idx_merchant_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商户信息表';

-- 第二步：创建默认商户（ID=1）
INSERT INTO `sys_merchant` 
(`id`, `merchant_code`, `merchant_name`, `merchant_type`, `level`, `path`, 
 `contact_name`, `contact_phone`, `contact_email`, `status`, `merchant_level`, 
 `operator_id`, `operator_name`, `created_at`, `updated_at`) 
VALUES 
(1, 'DEFAULT_MERCHANT', '默认商户', 'ENTERPRISE', 1, '1', 
 '系统管理员', '13800000000', 'admin@system.com', 'ACTIVE', 'VIP', 
 1, '系统', NOW(), NOW());

-- 第三步：扩展sys_user表，增加商户相关字段
ALTER TABLE `sys_user` 
ADD COLUMN `merchant_id` bigint(20) unsigned DEFAULT NULL COMMENT '所属商户ID' AFTER `authority_id`,
ADD COLUMN `is_main_account` tinyint(1) DEFAULT '0' COMMENT '是否为主账号' AFTER `merchant_id`;

-- 第四步：创建商户相关索引
ALTER TABLE `sys_user` 
ADD INDEX `idx_user_merchant_id` (`merchant_id`);

-- 第五步：处理手机号唯一约束（先删除后创建）
ALTER TABLE `sys_user` DROP INDEX IF EXISTS `idx_user_phone`;
ALTER TABLE `sys_user` 
ADD UNIQUE INDEX `idx_user_phone_merchant_unique` (`phone`, `merchant_id`, `deleted_at`);

-- 第六步：更新现有员工数据，设置商户ID为1（排除超级管理员）
UPDATE `sys_user` 
SET `merchant_id` = 1 
WHERE `authority_id` != 'super_admin' AND `merchant_id` IS NULL;

-- 第七步：设置默认商户的主管理员账号（取第一个非超级管理员用户）
UPDATE `sys_user` 
SET `is_main_account` = 1 
WHERE `merchant_id` = 1 
AND `authority_id` != 'super_admin'
AND `id` = (
    SELECT temp.min_id FROM (
        SELECT MIN(id) as min_id 
        FROM `sys_user` 
        WHERE `merchant_id` = 1 AND `authority_id` != 'super_admin'
    ) temp
);

-- 第八步：更新默认商户的管理员用户ID
UPDATE `sys_merchant` 
SET `admin_user_id` = (
    SELECT `id` FROM `sys_user` 
    WHERE `merchant_id` = 1 AND `is_main_account` = 1 
    LIMIT 1
) 
WHERE `id` = 1;

-- 第九步：扩展sys_authority表，增加商户相关字段
ALTER TABLE `sys_authority` 
ADD COLUMN `merchant_id` bigint(20) unsigned DEFAULT NULL COMMENT '所属商户ID' AFTER `parent_id`,
ADD COLUMN `is_system_role` tinyint(1) DEFAULT '0' COMMENT '是否为系统角色' AFTER `merchant_id`;

-- 第十步：更新所有现有角色的merchant_id为1（排除系统角色）
UPDATE `sys_authority` 
SET `merchant_id` = 1 
WHERE `authority_id` != 'super_admin' AND `merchant_id` IS NULL;

-- 第十一步：标识系统角色
UPDATE `sys_authority` 
SET `is_system_role` = 1 
WHERE `authority_id` = 'super_admin';

-- 迁移完成提示
SELECT '数据库迁移已完成！现有员工和角色都已归属于商户ID=1。' as '提示信息';
```

#### 迁移后验证清单

执行以上迁移脚本后，请验证以下内容：

1. **默认商户创建成功**：
   - sys_merchant表中存在ID=1的记录
   - 商户状态为ACTIVE
   - 已指定管理员用户

2. **现有员工数据迁移**：
   - 所有非超级管理员的merchant_id = 1
   - 超级管理员的merchant_id = NULL
   - 已指定主管理员账号

3. **现有角色数据迁移**：
   - 所有非系统角色的merchant_id = 1
   - 超级管理员角色标记为系统角色

4. **数据一致性验证**：
   - 手机号唯一约束已更新为商户级别
   - 所有现有业务数据正常运行
   - 权限控制仍然有效

#### 重要说明

- **向后兼容**：所有现有功能在迁移后仍然正常工作，用户无需重新登录
- **数据安全**：迁移过程中不会丢失任何现有数据
- **权限保持**：现有用户的所有权限在迁移后保持不变
- **渐进升级**：系统可以在运行中逐步添加新商户，无需停服升级
              <el-tag v-if="merchant.isDefault" type="warning" size="small">default</el-tag>
            </div>
          </div>
          
          <div class="merchant-actions">
            <el-icon><ArrowRight /></el-icon>
          </div>
        </div>
      </div>
      
      <div class="selector-footer">
        <el-button @click="handleLogout" type="text">退出登录</el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

import { useRouter } from 'vue-router'
import { usemerchantStore } from '@/stores/merchant'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import { Refresh, ArrowRight } from '@element-plus/icons-vue'

const router = useRouter()
const merchantStore = usemerchantStore()
const userStore = useUserStore()

const merchants = ref([])
const loading = ref(false)

onMounted(async () => {
  await loadMerchants()
})

async function loadMerchants() {
  try {
    loading.value = true
    const response = await merchantStore.fetchUserMerchants()
    merchants.value = response.merchants
  } catch (error) {
    ElMessage.error('获取商户列表失败')
  } finally {
    loading.value = false
  }
}

async function selectMerchant(merchant) {
  try {
    await merchantStore.switchMerchant(merchant.merchantId)
    ElMessage.success(`切换到 ${merchant.merchantName}`)
    
    // 根据用户角色跳转不同页面
    const userRole = userStore.userInfo.authorityId
    if (userRole === 'super_admin') {
      router.push('/system/merchant')
    } else if (userRole.startsWith('merchant_admin')) {
      router.push('/merchant/info')
    } else {
      router.push('/workbench')
    }
  } catch (error) {
    ElMessage.error('切换商户失败')
  }
}

function refreshMerchants() {
  loadMerchants()
}

function handleLogout() {
  userStore.logout()
  router.push('/login')
}

function getStatusType(status) {
  const typeMap = {
    'ACTIVE': 'success',
    'PENDING': 'warning',
    'SUSPENDED': 'info',
    'DISABLED': 'danger'
  }
  return typeMap[status] || 'info'
}

function getStatusText(status) {
  const textMap = {
    'ACTIVE': '正常',
    'PENDING': '待审核',
    'SUSPENDED': '暂停',
    'DISABLED': '已禁用'
  }
  return textMap[status] || status
}
</script>

<style scoped>
.merchant-selector {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.selector-card {
  width: 400px;
  max-height: 600px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.merchant-list {
  max-height: 400px;
  overflow-y: auto;
}

.merchant-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  margin-bottom: 10px;
  border: 1px solid #e4e7ed;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s;
}

.merchant-item:hover {
  border-color: #409eff;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.merchant-item.is-default {
  border-color: #e6a23c;
  background-color: #fdf6ec;
}

.merchant-info {
  flex: 1;
}

.merchant-name {
  font-weight: bold;
  margin-bottom: 5px;
}

.merchant-code {
  color: #909399;
  font-size: 12px;
  margin-bottom: 5px;
}

.merchant-status {
  display: flex;
  gap: 5px;
}

.merchant-actions {
  color: #c0c4cc;
}

.selector-footer {
  text-align: center;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #e4e7ed;
}
</style>
```

#### 商户管理列表组件

```vue
<!-- views/merchant/list.vue -->
<template>
  <div class="merchant-list-container">
    <!-- 搜索区域 -->
    <el-card class="search-card">
      <el-form :model="searchForm" :inline="true">
        <el-form-item label="商户名称">
          <el-input 
            v-model="searchForm.merchantName" 
            placeholder="请输入商户名称"
            clearable />
        </el-form-item>
        
        <el-form-item label="商户状态">
          <el-select v-model="searchForm.status" placeholder="请选择状态" clearable>
            <el-option label="正常" value="ACTIVE" />
            <el-option label="待审核" value="PENDING" />
            <el-option label="暂停" value="SUSPENDED" />
            <el-option label="已禁用" value="DISABLED" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="商户类型">
          <el-select v-model="searchForm.merchantType" placeholder="请选择类型" clearable>
            <el-option label="企业" value="ENTERPRISE" />
            <el-option label="个人" value="INDIVIDUAL" />
          </el-select>
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
    
    <!-- 操作区域 -->
    <el-card class="action-card">
      <el-row :gutter="20">
        <el-col :span="12">
          <el-button 
            v-role-auth="'system:merchant:create'"
            type="primary" 
            @click="handleCreate">
            <el-icon><Plus /></el-icon>
            新增商户
          </el-button>
          
          <el-button 
            v-role-auth="'system:merchant:permission'"
            type="warning" 
            @click="handleBatchPermission"
            :disabled="selectedMerchants.length === 0">
            <el-icon><Key /></el-icon>
            批量分配权限
          </el-button>
        </el-col>
        
        <el-col :span="12" style="text-align: right;">
          <el-switch 
            v-model="treeView"
            active-text="树形视图"
            inactive-text="列表视图" />
        </el-col>
      </el-row>
    </el-card>
    
    <!-- 数据展示区域 -->
    <el-card class="data-card">
      <!-- 树形视图 -->
      <div v-if="treeView" class="tree-view">
        <el-tree
          :data="merchantTree"
          :props="treeProps"
          node-key="merchantId"
          default-expand-all
          :expand-on-click-node="false">
          
          <template #default="{ node, data }">
            <div class="tree-node">
              <div class="node-info">
                <span class="node-name">{{ data.merchantName }}</span>
                <el-tag :type="getStatusType(data.status)" size="small">{{ getStatusText(data.status) }}</el-tag>
              </div>
              
              <div class="node-actions">
                <el-button 
                  v-role-auth="'system:merchant:update'"
                  type="text" 
                  size="small" 
                  @click="handleEdit(data)">
                  编辑
                </el-button>
                
                <el-button 
                  v-role-auth="'system:merchant:permission'"
                  type="text" 
                  size="small" 
                  @click="handlePermission(data)">
                  权限
                </el-button>
                
                <el-dropdown trigger="click">
                  <el-button type="text" size="small">
                    更多<el-icon><ArrowDown /></el-icon>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item @click="handleAddChild(data)">添加子商户</el-dropdown-item>
                      <el-dropdown-item @click="handleMove(data)">移动位置</el-dropdown-item>
                      <el-dropdown-item divided @click="handleDelete(data)">删除</el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </div>
            </div>
          </template>
        </el-tree>
      </div>
      
      <!-- 列表视图 -->
      <div v-else class="table-view">
        <el-table 
          :data="merchantList" 
          v-loading="loading"
          @selection-change="handleSelectionChange">
          
          <el-table-column type="selection" width="55" />
          
          <el-table-column prop="merchantCode" label="商户编码" width="150" />
          
          <el-table-column prop="merchantName" label="商户名称" min-width="200" />
          
          <el-table-column prop="merchantType" label="类型" width="100">
            <template #default="{ row }">
              {{ row.merchantType === 'ENTERPRISE' ? '企业' : '个人' }}
            </template>
          </el-table-column>
          
          <el-table-column prop="level" label="层级" width="80" />
          
          <el-table-column prop="contactName" label="联系人" width="120" />
          
          <el-table-column prop="contactPhone" label="联系电话" width="150" />
          
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
            </template>
          </el-table-column>
          
          <el-table-column prop="createdAt" label="创建时间" width="180">
            <template #default="{ row }">
              {{ formatDate(row.createdAt) }}
            </template>
          </el-table-column>
          
          <el-table-column label="操作" width="300" fixed="right">
            <template #default="{ row }">
              <el-button 
                v-role-auth="'system:merchant:update'"
                type="primary" 
                size="small" 
                @click="handleEdit(row)">
                编辑
              </el-button>
              
              <el-button 
                v-role-auth="'system:merchant:permission'"
                type="warning" 
                size="small" 
                @click="handlePermission(row)">
                权限
              </el-button>
              
              <el-button 
                v-role-auth="'system:merchant:status'"
                :type="row.status === 'ACTIVE' ? 'info' : 'success'" 
                size="small" 
                @click="handleToggleStatus(row)">
                {{ row.status === 'ACTIVE' ? '禁用' : '启用' }}
              </el-button>
              
              <el-button 
                v-role-auth="'system:merchant:delete'"
                type="danger" 
                size="small" 
                @click="handleDelete(row)">
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
        
        <!-- 分页 -->
        <div class="pagination-container">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="pagination.total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange" />
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { usemerchantStore } from '@/stores/merchant'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Key, ArrowDown } from '@element-plus/icons-vue'
import { formatDate } from '@/utils/date'

const merchantStore = usemerchantStore()

const loading = ref(false)
const treeView = ref(false)
const merchantList = ref([])
const merchantTree = ref([])
const selectedMerchants = ref([])

const searchForm = reactive({
  merchantName: '',
  status: '',
  merchantType: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const treeProps = {
  children: 'children',
  label: 'merchantName'
}

onMounted(() => {
  loadData()
})

async function loadData() {
  if (treeView.value) {
    await loadMerchantTree()
  } else {
    await loadMerchantList()
  }
}

async function loadMerchantList() {
  try {
    loading.value = true
    const params = {
      ...searchForm,
      page: pagination.page,
      pageSize: pagination.pageSize
    }
    const response = await merchantStore.fetchMerchantList(params)
    merchantList.value = response.list
    pagination.total = response.total
  } catch (error) {
    ElMessage.error('获取商户列表失败')
  } finally {
    loading.value = false
  }
}

async function loadMerchantTree() {
  try {
    loading.value = true
    const response = await merchantStore.fetchMerchantTree()
    merchantTree.value = response.merchantTree
  } catch (error) {
    ElMessage.error('获取商户树结构失败')
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  pagination.page = 1
  loadData()
}

function handleReset() {
  Object.assign(searchForm, {
    merchantName: '',
    status: '',
    merchantType: ''
  })
  pagination.page = 1
  loadData()
}

function handleCreate() {
  // 跳转到创建页面
}

function handleEdit(row) {
  // 跳转到编辑页面
}

function handlePermission(row) {
  // 打开权限分配对话框
}

function handleDelete(row) {
  ElMessageBox.confirm(`确定要删除商户「${row.merchantName}」吗？`, '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await merchantStore.deleteMerchant(row.merchantId)
      ElMessage.success('删除成功')
      loadData()
    } catch (error) {
      ElMessage.error('删除失败')
    }
  })
}

function handleSelectionChange(selection) {
  selectedMerchants.value = selection
}

function handleSizeChange(size) {
  pagination.pageSize = size
  pagination.page = 1
  loadData()
}

function handleCurrentChange(page) {
  pagination.page = page
  loadData()
}

function getStatusType(status) {
  const typeMap = {
    'ACTIVE': 'success',
    'PENDING': 'warning',
    'SUSPENDED': 'info',
    'DISABLED': 'danger'
  }
  return typeMap[status] || 'info'
}

function getStatusText(status) {
  const textMap = {
    'ACTIVE': '正常',
    'PENDING': '待审核',
    'SUSPENDED': '暂停',
    'DISABLED': '已禁用'
  }
  return textMap[status] || status
}
</script>

<style scoped>
.merchant-list-container {
  padding: 20px;
}

.search-card, .action-card, .data-card {
  margin-bottom: 20px;
}

.tree-node {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  padding-right: 20px;
}

.node-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.node-name {
  font-weight: bold;
}

.pagination-container {
  margin-top: 20px;
  text-align: right;
}
</style>
```

### 监控和运维设计

#### 系统监控指标

**关键性能指标（KPI）**：
- **商户管理指标**：
  - 商户注册量（日/周/月）
  - 商户激活率（有效登录的商户比例）
  - 平均商户员工数量
  - 商户认证通过率

- **性能指标**：
  - API响应时间（P95/P99）
  - 数据库连接池使用率
  - Redis缓存命中率
  - 并发用户数

- **安全指标**：
  - 登录失败次数
  - 越权访问尝试次数
  - 异常权限操作次数
  - SQL注入尝试次数



### JWT令牌扩展设计

#### 扩展JWT Claims结构

```go
// model/jwt_claims.go
type CustomClaims struct {
    BaseClaims
    BufferTime   int64  `json:"bufferTime"`
    AuthorityId  string `json:"authorityId"`
    MerchantID   uint   `json:"merchantId"`   // 新增：当前商户ID
    MerchantName string `json:"merchantName"` // 新增：当前商户名称
    UserRole     string `json:"userRole"`     // 新增：用户角色类型
    Permissions  []string `json:"permissions"` // 新增：用户权限列表
    jwt.RegisteredClaims
}

// 生成包含商户信息的JWT Token
func (j *JWT) CreateTokenByOldToken(oldToken string, merchantID uint) (string, error) {
    // 解析旧Token
    claims, err := j.ParseToken(oldToken)
    if err != nil {
        return "", err
    }
    
    // 获取商户信息
    merchantService := service.ServiceGroupApp.MerchantServiceGroup.MerchantService
    merchant, err := merchantService.GetMerchantByID(merchantID)
    if err != nil {
        return "", err
    }
    
    // 获取用户在该商户下的权限
    userService := service.ServiceGroupApp.SystemServiceGroup.UserService
    permissions, err := userService.GetUserPermissionsInMerchant(claims.BaseClaims.ID, merchantID)
    if err != nil {
        return "", err
    }
    
    // 更新Claims
    claims.MerchantID = merchantID
    claims.MerchantName = merchant.MerchantName
    claims.Permissions = permissions
    claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour))
    
    // 生成新Token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(j.SigningKey))
}
```

#### 商户切换服务实现

```go
// service/merchant_switch_service.go
type MerchantSwitchService struct {
    db *gorm.DB
}

// 切换商户
func (m *MerchantSwitchService) SwitchMerchant(userID uint, merchantID uint) (*response.SwitchMerchantResponse, error) {
    // 验证用户是否属于该商户
    var merchantUser model.SysMerchantUser
    err := m.db.Where("user_id = ? AND merchant_id = ? AND status = ?", userID, merchantID, "ACTIVE").First(&merchantUser).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("用户不属于该商户")
        }
        return nil, err
    }
    
    // 获取商户信息
    var merchant model.SysMerchant
    err = m.db.Where("id = ? AND status = ?", merchantID, "ACTIVE").First(&merchant).Error
    if err != nil {
        return nil, errors.New("商户不存在或已禁用")
    }
    
    // 获取用户信息
    var user model.SysUser
    err = m.db.Where("id = ?", userID).First(&user).Error
    if err != nil {
        return nil, err
    }
    
    // 确定用户在该商户下的角色
    userRole := m.determineUserRole(user, merchant)
    
    // 获取用户在该商户下的权限
    permissions, err := m.getUserPermissionsInMerchant(userID, merchantID, userRole)
    if err != nil {
        return nil, err
    }
    
    // 生成新的JWT Token
    jwtService := utils.JWT{
        SigningKey: []byte(global.GVA_CONFIG.JWT.SigningKey),
    }
    
    claims := systemReq.CustomClaims{
        BaseClaims: systemReq.BaseClaims{
            ID:       user.ID,
            Username: user.Username,
            NickName: user.NickName,
        },
        BufferTime:   global.GVA_CONFIG.JWT.BufferTime,
        AuthorityId:  userRole,
        MerchantID:   merchantID,
        MerchantName: merchant.MerchantName,
        UserRole:     userRole,
        Permissions:  permissions,
        RegisteredClaims: jwt.RegisteredClaims{
            Issuer:    global.GVA_CONFIG.JWT.Issuer,
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
            NotBefore: jwt.NewNumericDate(time.Now().Add(-1000)),
        },
    }
    
    token, err := jwtService.CreateToken(claims)
    if err != nil {
        return nil, err
    }
    
    // 记录切换日志
    m.logMerchantSwitch(userID, merchantID, "SUCCESS")
    
    return &response.SwitchMerchantResponse{
        Token: token,
        MerchantInfo: response.MerchantInfo{
            MerchantId:   merchantID,
            MerchantName: merchant.MerchantName,
            UserRole:     userRole,
        },
    }, nil
}

// 确定用户角色
func (m *MerchantSwitchService) determineUserRole(user model.SysUser, merchant model.SysMerchant) string {
    // 如果是商户管理员
    if merchant.AdminUserID != nil && *merchant.AdminUserID == user.ID {
        return fmt.Sprintf("merchant_admin_%d", merchant.ID)
    }
    
    // 如果有指定的角色
    if user.AuthorityId != "" {
        return user.AuthorityId
    }
    
    // 默认为员工角色
    return fmt.Sprintf("employee_default_%d", merchant.ID)
}
```
- **分层缓存**：
  - L1: 当前商户信息缓存（内存级）
  - L2: 商户用户列表缓存（Redis）
  - L3: 商户角色权限缓存（Redis）
  - L4: 商户树结构缓存（Redis）
- **缓存Key设计**：包含商户ID的分层Key结构
  - `merchant:{merchantId}:user:list`
  - `merchant:{merchantId}:role:list`
  - `merchant:{merchantId}:permission:cache`
  - `merchant:tree:structure` （全局树结构）
  - `merchant:{merchantId}:children` （子商户列表）
- **缓存失效**：商户数据变更时精准失效相关缓存
- **缓存隔离**：不同商户的缓存数据完全隔离，防止跨租户数据泄露
- **层级缓存优化**：商户移动后自动刷新影响的所有缓存

### 前端性能优化
- **懒加载组件**：商户相关组件按需加载，减少初始打包体积
- **状态管理优化**：使用Pinia的持久化功能缓存商户上下文
- **防抖处理**：商户切换和搜索功能使用防抖避免频繁请求
- **虚拟列表**：大数据量列表使用虚拟滚动技术
- **数据预加载**：预加载用户常用功能数据，提升交互体验

## 错误处理设计

### 多租户错误处理
``json
{
  "code": 40301,
  "message": "商户数据访问被拒绝",
  "data": {
    "requested_merchant_id": 2,
    "current_merchant_id": 1,
    "user_role": "merchant_admin"
  },
  "timestamp": "2024-01-01T10:00:00Z",
  "path": "/api/v1/merchant/user/list"
}
```

### 多租户错误码定义

| 错误码 | 错误信息 | 说明 |
|--------|----------|------|
| 40301 | 商户数据访问被拒绝 | 跨商户数据访问被拦截 |
| 40302 | 商户上下文缺失| JWT中缺少商户ID信息 |
| 40303 | 商户状态异常 | 商户已停用或禁用 |
| 40304 | 商户切换失败 | 用户不属于目标商户 |
| 40305 | 商户员工数量超限 | 超出商户员工数量限制 |
| 40306 | 商户角色数量超限 | 超出商户角色数量限制 |
| 40307 | 手机号在当前商户已存在 | 同一商户内手机号重复 |
| 40308 | 角色名称在当前商户已存在 | 同一商户内角色名称重复 |
| 40310 | 商户层级结构异常 | 商户层级数据不一致 |
| 40311 | 商户移动操作失败 | 不能移动到子节点或形成循环 |
| 40312 | 商户层级超出限制 | 超出最大允许的层级深度 |

### 前端错误处理策略
- **全局错误拦截**：在Axios响应拦截器中统一处理多租户错误
- **商户上下文失效**：当检测到商户上下文失效时，自动重定向到商户选择页面
- **权限不足处理**：提供明确的权限不足提示和申请权限入口
- **网络重试机制**：商户切换失败时提供重试选项
- **降级处理**：服务异常时保留基础查看功能

## 测试策略

### 多租户单元测试
- **数据隔离测试**：验证不同商户数据的完全隔离
- **权限控制测试**：测试跨商户访问被正确拦截
- **商户切换测试**：验证商户切换功能的正确性
- **多租户业务逻辑测试**：测试商户管理、员工管理、角色管理的核心业务

### 多租户集成测试
- **数据库集成测试**：测试多租户数据持久化和查询功能
- **缓存集成测试**：测试多租户缓存的隔离性和一致性
- **文件存储测试**：测试多租户文件的存储隔离和访问控制
- **权限集成测试**：测试与扩展Casbin多租户权限系统的集成

### 前端多租户测试
- **组件单元测试**：使用Vue Test Utils测试多租户组件功能
- **商户切换E2E测试**：使用Cypress测试完整的商户切换流程
- **权限界面测试**：验证不同商户角色下的界面展示
- **多租户表单验证测试**：测试多租户场景下的表单验证规则

### 测试用例示例
```
// 多租户数据隔离测试
describe('merchantDataIsolation', () => {
  test('应该只返回当前商户的用户数据', async () => {
    const merchantId = 1
    const userService = new merchantUserService()
    
    // 设置商户上下文
    userService.setmerchantContext({ merchantId })
    
    const users = await userService.getUserList()
    
    // 验证所有返回的用户都属于当前商户
    users.forEach(user => {
      expect(user.merchantId).toBe(merchantId)
    })
  })
  
  test('跨商户数据访问应该被拦截', async () => {
    const currentMerchantId = 1
    const targetMerchantId = 2
    const userService = new merchantUserService()
    
    userService.setmerchantContext({ merchantId: currentMerchantId })
    
    // 尝试访问其他商户的用户数据
    await expect(userService.getUserById(targetMerchantId, 123))
      .rejects.toThrow('商户数据访问被拒绝')
  })
})

// 商户切换功能测试
describe('MerchantSwitching', () => {
  test('用户应该能够成功切换到有权限的商户', async () => {
    const merchantStore = usemerchantStore()
    const targetMerchantId = 2
    
    // 模拟用户有多个商户权限
    merchantStore.userMerchants = [
      { merchantId: 1, merchantName: '商户A' },
      { merchantId: 2, merchantName: '商户B' }
    ]
    
    const result = await merchantStore.switchMerchant(targetMerchantId)
    
    expect(result.merchantInfo.merchantId).toBe(targetMerchantId)
    expect(merchantStore.currentMerchantId).toBe(targetMerchantId)
  })
  
  test('切换到无权限商户应该失败', async () => {
    const merchantStore = usemerchantStore()
    const unauthorizedMerchantId = 999
    
    await expect(merchantStore.switchMerchant(unauthorizedMerchantId))
      .rejects.toThrow('商户切换失败')
  })
})