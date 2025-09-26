# 用户操作审计功能开发文档

## 1. 功能概述

为了加强系统的安全性和可追溯性，需要在用户进行新增、修改、删除操作时，记录执行操作的用户信息，具体包括：
- 在用户表中增加两个新字段：操作者ID和操作者姓名
- 在用户执行新增、修改、删除操作时，自动更新这两个字段
- 确保所有用户相关操作都能正确记录操作人信息，便于后续审计和责任追溯

## 2. 需求分析

### 2.1 业务需求
- 系统需要能够记录每个用户操作的执行者信息
- 支持对历史操作进行审计和查询
- 确保操作责任可追溯

### 2.2 技术需求
- 在现有用户表结构中增加两个新字段
- 确保所有用户相关的CRUD操作都能正确更新操作人信息
- 不影响现有业务流程和功能

## 3. 技术方案

### 3.1 总体架构
本功能将在现有系统架构基础上进行扩展，主要涉及以下几个层面：
- 数据模型层：在用户表中增加操作人相关字段
- 服务层：修改用户相关服务方法，添加操作人参数
- 控制器层：从请求上下文中获取当前登录用户信息，并传递给服务层

### 3.2 技术栈
- 后端框架：Gin
- 数据库：MySQL
- ORM框架：GORM
- 认证：JWT

## 4. 数据库设计

### 4.1 表结构变更
需要在`sys_users`表中增加两个新字段：

| 字段名 | 数据类型 | 长度 | 约束 | 描述 |
|-------|---------|-----|-----|------|
| operator_id | INT | 10 | 可为空 | 执行操作的用户ID |
| operator_name | VARCHAR | 100 | 可为空 | 执行操作的用户名 |

## 5. 代码实现

### 5.1 数据模型层修改

修改`server/model/system/sys_user.go`文件，在`SysUser`结构体中添加两个新字段：

```go
// SysUser 用户表
type SysUser struct {
    global.GVA_MODEL
    UUID          uuid.UUID      `json:"uuid" gorm:"index;comment:用户UUID"`
    Username      string         `json:"userName" gorm:"index;comment:用户登录名"`
    Password      string         `json:"-"  gorm:"comment:用户登录密码"`
    NickName      string         `json:"nickName" gorm:"default:系统用户;comment:用户昵称"`
    Name          string         `json:"name" gorm:"default:'';comment:用户姓名"`
    HeaderImg     string         `json:"headerImg" gorm:"default:https://qmplusimg.henrongyi.top/gva_header.jpg;comment:用户头像"`
    AuthorityId   uint           `json:"authorityId" gorm:"default:888;comment:用户角色ID"`
    Authority     SysAuthority   `json:"authority" gorm:"foreignKey:AuthorityId;references:AuthorityId;comment:用户角色"`
    Authorities   []SysAuthority `json:"authorities" gorm:"many2many:sys_user_authority;"`
    Phone         string         `json:"phone"  gorm:"index:idx_phone_not_deleted,unique,condition:deleted_at IS NULL;comment:用户手机号"`
    Email         string         `json:"email"  gorm:"comment:用户邮箱"`
    Enable        int            `json:"enable" gorm:"default:1;comment:用户是否被冻结 1正常 2冻结"`
    OriginSetting common.JSONMap `json:"originSetting" form:"originSetting" gorm:"type:text;default:null;column:origin_setting;comment:配置;"`
    OperatorId    uint           `json:"operatorId" gorm:"comment:操作者ID"`
    OperatorName  string         `json:"operatorName" gorm:"comment:操作者姓名;size:100"`
}
```

### 5.2 服务层修改

修改`server/service/system/sys_user.go`文件中的相关方法，添加操作人参数：

#### 5.2.1 Register方法
```go
// Register 用户注册
def (s *UserService) Register(u system.SysUser, operatorId uint, operatorName string) (system.SysUser, error) {
    // 设置操作人信息
    u.OperatorId = operatorId
    u.OperatorName = operatorName
    
    // 原有注册逻辑...
    // 唯一性检查
    // 密码加密
    // 保存用户
}
```

#### 5.2.2 SetUserInfo方法
```go
// SetUserInfo 设置用户信息
def (s *UserService) SetUserInfo(u system.SysUser, operatorId uint, operatorName string) error {
    // 设置操作人信息
    var user system.SysUser
    err := global.GVA_DB.Where("id = ?", u.ID).First(&user).Error
    if err != nil {
        return err
    }
    
    // 更新用户信息
    user.OperatorId = operatorId
    user.OperatorName = operatorName
    // 其他字段更新...
    
    return global.GVA_DB.Save(&user).Error
}
```

#### 5.2.3 DeleteUser方法
```go
// DeleteUser 删除用户
def (s *UserService) DeleteUser(id uint, operatorId uint, operatorName string) error {
    // 事务处理
    return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
        // 先查询用户信息
        var user system.SysUser
        err := tx.Where("id = ?", id).First(&user).Error
        if err != nil {
            return err
        }
        
        // 更新操作人信息
        user.OperatorId = operatorId
        user.OperatorName = operatorName
        err = tx.Save(&user).Error
        if err != nil {
            return err
        }
        
        // 删除用户
        if err := tx.Where("id = ?", id).Unscoped().Delete(&system.SysUser{}).Error; err != nil {
            return err
        }
        
        // 删除用户与角色关联
        if err := tx.Where("sys_user_id = ?", id).Delete(&system.SysUseAuthority{}).Error; err != nil {
            return err
        }
        
        return nil
    })
}
```

### 5.3 控制器层修改

修改`server/api/v1/system/sys_user.go`文件中的相关方法，从请求上下文中获取当前登录用户信息并传递给服务层：

#### 5.3.1 Register方法
```go
// Register
// @Tags     SysUser
// @Summary  用户注册账号
// @Produce   application/json
// @Param    data  body      systemReq.Register                                            true  "用户名, 昵称, 密码, 角色ID"
// @Success  200   {object}  response.Response{data=systemRes.SysUserResponse,msg=string}  "用户注册账号,返回包括用户信息"
// @Router   /user/admin_register [post]
func (b *BaseApi) Register(c *gin.Context) {
    var r systemReq.Register
    err := c.ShouldBindJSON(&r)
    if err != nil {
        response.FailWithMessage("请求参数格式不正确，请检查数据格式", c)
        return
    }
    err = utils.Verify(r, utils.RegisterVerify)
    if err != nil {
        response.FailWithMessage(err.Error(), c)
        return
    }
    
    // 获取当前登录用户信息作为操作人
    operatorId := utils.GetUserID(c)
    operatorName := utils.GetUserName(c)
    
    var authorities []system.SysAuthority
    for _, v := range r.AuthorityIds {
        authorities = append(authorities, system.SysAuthority{
            AuthorityId: v,
        })
    }
    // 获取Enable值，默认为1(启用)
    enable := r.Enable
    if enable == 0 {
        enable = 1
    }
    user := &system.SysUser{Username: r.Username, NickName: r.NickName, Name: r.Name, Password: r.Password, HeaderImg: r.HeaderImg, AuthorityId: r.AuthorityId, Authorities: authorities, Enable: enable, Phone: r.Phone, Email: r.Email}
    userReturn, err := userService.Register(*user, operatorId, operatorName)
    
    // 后续错误处理和响应逻辑...
}
```

#### 5.3.2 SetUserInfo方法
```go
// SetUserInfo
// @Tags      SysUser
// @Summary   设置用户信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysUser                                             true  "ID, 用户名, 昵称, 头像链接"
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "设置用户信息"
// @Router    /user/setUserInfo [put]
func (b *BaseApi) SetUserInfo(c *gin.Context) {
    var user systemReq.ChangeUserInfo
    err := c.ShouldBindJSON(&user)
    if err != nil {
        response.FailWithMessage(err.Error(), c)
        return
    }
    err = utils.Verify(user, utils.IdVerify)
    if err != nil {
        response.FailWithMessage(err.Error(), c)
        return
    }
    
    // 获取当前登录用户信息作为操作人
    operatorId := utils.GetUserID(c)
    operatorName := utils.GetUserName(c)
    
    if len(user.AuthorityIds) != 0 {
        authorityID := utils.GetUserAuthorityId(c)
        err = userService.SetUserAuthorities(authorityID, user.ID, user.AuthorityIds)
        if err != nil {
            global.GVA_LOG.Error("设置失败!", zap.Error(err))
            response.FailWithMessage("设置失败", c)
            return
        }
    }
    err = userService.SetUserInfo(system.SysUser{
        GVA_MODEL: global.GVA_MODEL{
            ID: user.ID,
        },
        NickName:  user.NickName,
        Name:      user.Name,
        Username:  user.UserName,
        HeaderImg: user.HeaderImg,
        Phone:     user.Phone,
        Email:     user.Email,
        Enable:    user.Enable,
    }, operatorId, operatorName)
    
    // 后续错误处理和响应逻辑...
}
```

#### 5.3.3 DeleteUser方法
```go
// DeleteUser
// @Tags      SysUser
// @Summary   删除用户
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.GetById                true  "用户ID"
// @Success   200   {object}  response.Response{msg=string}  "删除用户"
// @Router    /user/deleteUser [delete]
func (b *BaseApi) DeleteUser(c *gin.Context) {
    var reqId request.GetById
    err := c.ShouldBindJSON(&reqId)
    if err != nil {
        response.FailWithMessage(err.Error(), c)
        return
    }
    err = utils.Verify(reqId, utils.IdVerify)
    if err != nil {
        response.FailWithMessage(err.Error(), c)
        return
    }
    jwtId := utils.GetUserID(c)
    if jwtId == uint(reqId.ID) {
        response.FailWithMessage("删除失败, 无法删除自己。", c)
        return
    }
    
    // 获取当前登录用户信息作为操作人
    operatorId := utils.GetUserID(c)
    operatorName := utils.GetUserName(c)
    
    err = userService.DeleteUser(reqId.ID, operatorId, operatorName)
    
    // 后续错误处理和响应逻辑...
}
```

## 6. 测试计划

### 6.1 单元测试
需要编写单元测试来验证以下场景：
- 用户注册时，操作人字段是否正确设置
- 修改用户信息时，操作人字段是否正确更新
- 删除用户时，操作人字段是否正确更新

### 6.2 集成测试
需要进行集成测试来验证完整的业务流程：
- 登录系统
- 执行用户注册、修改、删除操作
- 查询数据库，验证操作人字段是否正确记录

### 6.3 测试案例
| 测试场景 | 预期结果 |
|---------|---------|
| 管理员注册新用户 | 新用户记录中的operator_id和operator_name字段应包含管理员信息 |
| 管理员修改用户信息 | 被修改用户记录中的operator_id和operator_name字段应包含管理员信息 |
| 管理员删除用户 | 被删除用户记录中的operator_id和operator_name字段应包含管理员信息 |

## 7. 部署说明

### 7.1 数据库迁移
需要执行数据库迁移脚本，在`sys_users`表中添加两个新字段：

```sql
ALTER TABLE `sys_users` 
ADD COLUMN `operator_id` INT(10) NULL COMMENT '操作者ID',
ADD COLUMN `operator_name` VARCHAR(100) NULL COMMENT '操作者姓名';
```

### 7.2 代码部署
- 更新代码库中的相关文件
- 重新构建和部署应用程序
- 验证功能是否正常工作

## 8. 风险评估

- **数据一致性风险**：确保在所有用户相关操作中都正确更新操作人字段
- **性能风险**：新增字段对数据库查询性能的影响很小，可以忽略不计
- **兼容性风险**：需要确保现有功能不受新增字段的影响

## 9. 后续优化

- 考虑为更多业务表添加操作审计功能
- 实现操作日志查询界面，方便管理员查看和审计操作记录
- 添加操作时间戳字段，记录具体的操作时间