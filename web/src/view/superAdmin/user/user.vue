<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="用户名">
          <el-input v-model="searchInfo.username" placeholder="用户名" />
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="searchInfo.nickname" placeholder="昵称" />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="searchInfo.phone" placeholder="手机号" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="searchInfo.email" placeholder="邮箱" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">
            查询
          </el-button>
          <el-button icon="refresh" @click="onReset"> 重置 </el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="addUser"
          >新增用户</el-button
        >
      </div>
      <el-table :data="tableData" row-key="ID">
        <el-table-column align="center" label="头像" width="70">
          <template #default="scope">
            <CustomPic style="display: flex; align-items: center; justify-content: center; " :pic-src="scope.row.headerImg" :size="40" />
          </template>
        </el-table-column>

        <el-table-column
          align="left"
          label="姓名/ID"
          min-width="150"
        >
          <template #default="scope">
            <div>
              <div><strong>{{ scope.row.name }}</strong></div>
              <div class="text-xs text-gray-400">ID: {{ scope.row.ID }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="手机号/用户名"
          min-width="201"
        >
          <template #default="scope">
            <div>
              <div>{{ scope.row.phone }}</div>
              <div class="text-xs text-gray-400">{{ scope.row.userName }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="昵称/邮箱"
          min-width="200"
        >
          <template #default="scope">
            <div>
              <div>{{ scope.row.nickName }}</div>
              <div class="text-xs text-gray-400">{{ scope.row.email }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="用户角色" min-width="200">
          <template #default="scope">
            <div>
              {{ scope.row.authorities?.map(auth => auth.authorityName).join(', ') || '无' }}
            </div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="状态" width="70">
          <template #default="scope">
            <el-switch
              v-model="scope.row.enable"
              inline-prompt
              :active-value="1"
              :inactive-value="2"
              :disabled="true"
              @change="
                () => {
                  switchEnable(scope.row)
                }
              "
            />
          </template>
        </el-table-column>

        <el-table-column align="left" label="操作者/时间" min-width="190">
          <template #default="scope">
            <div>
              <div>{{ scope.row.operatorName || '-' }} <span class="text-xs text-gray-400">ID:{{ scope.row.operatorId }}</span></div>
              <div class="text-xs text-gray-400">{{ formatDate(scope.row.updatedAt) }}</div>
            </div>
          </template>
        </el-table-column>

        <el-table-column align="left" label="操作"  fixed="right" width="280 ">
          <template #default="scope">
            
            <el-button 
            type="primary" 
            link 
            icon="View"
            @click="openDetailDialog(scope.row)"
            >查看</el-button
            >
              <el-button
              type="primary"
              link
              icon="edit"
              @click="openEdit(scope.row)"
              >编辑</el-button
            >

            <el-button
              type="primary"
              link
              icon="delete"
              @click="deleteUserFunc(scope.row)"
              >删除</el-button
            >

            <el-button
              type="primary"
              link
              icon="magic-stick"
              @click="resetPasswordFunc(scope.row)"
              >密码</el-button
            >
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>
    <!-- 重置密码对话框 -->
    <el-dialog
      v-model="resetPwdDialog"
      title="重置密码"
      width="500px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
    >
      <el-form :model="resetPwdInfo" ref="resetPwdForm" label-width="100px">
        <el-form-item label="用户账号">
          <el-input v-model="resetPwdInfo.userName" disabled />
        </el-form-item>
        <el-form-item label="用户昵称">
          <el-input v-model="resetPwdInfo.nickName" disabled />
        </el-form-item>
        <el-form-item label="新密码">
          <div class="flex w-full">
            <el-input class="flex-1" v-model="resetPwdInfo.password" placeholder="请输入新密码" show-password />
            <el-button type="primary" @click="generateRandomPassword" style="margin-left: 10px">
              生成随机密码
            </el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeResetPwdDialog">取 消</el-button>
          <el-button type="primary" @click="confirmResetPassword">确 定</el-button>
        </div>
      </template>
    </el-dialog>
    
    <!-- 详情抽屉 -->
    <el-drawer
      v-model="detailDialogVisible"
      :size="appStore.drawerSize || '50%'"
      :show-close="false"
      :close-on-press-escape="false"
      :close-on-click-modal="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">用户详情</span>
          <div>
            <el-button @click="detailDialogVisible = false">关 闭</el-button>
          </div>
        </div>
      </template>

      <el-form
        ref="detailForm"
        label-width="80px"
        :model="detailData"
      >
        <!-- ID和UUID字段 -->
        <div class="form-row">
          <div class="flex items-start space-x-4">
            <el-form-item label="ID" class="id-field-small flex-1">
              <div class="flex items-center">
                <el-input v-model="detailData.id" disabled style="margin-right: 8px; min-width: 150px;" />
                <el-button 
                  type="text" 
                  size="small" 
                  @click="copyToClipboard(detailData.id, 'ID已复制')"
                  title="复制ID"
                >
                  <el-icon><copy-document /></el-icon>
                </el-button>
              </div>
            </el-form-item>
            <el-form-item label="UUID" class="uuid-field-full flex-1">
              <div class="flex items-center">
                <el-input v-model="detailData.uuid" disabled style="margin-right: 8px; min-width: 350px;" />
                <el-button 
                  type="text" 
                  size="small" 
                  @click="copyToClipboard(detailData.uuid, 'UUID已复制')"
                  title="复制UUID"
                >
                  <el-icon><copy-document /></el-icon>
                </el-button>
              </div>
            </el-form-item>
          </div>
        </div>
        
        <!-- 用户信息字段 -->
        <div class="form-row">
          <el-form-item
            label="用户名"
            class="form-half bold-label"
          >
            <el-input v-model="detailData.userName" disabled />
          </el-form-item>
          <el-form-item label="昵称" class="form-half bold-label">
            <el-input v-model="detailData.nickName" disabled />
          </el-form-item>
        </div>
        <div class="form-row">
          <el-form-item label="姓名" class="form-half bold-label">
            <el-input v-model="detailData.name" disabled />
          </el-form-item>
          <el-form-item label="邮箱" class="form-half bold-label">
            <el-input v-model="detailData.email" disabled />
          </el-form-item>
        </div>
        <div class="form-row">
          <el-form-item label="手机号" class="form-half bold-label">
            <el-input v-model="detailData.phone" disabled />
          </el-form-item>
          <el-form-item label="状态" class="form-half bold-label">
            <el-switch
              v-model="detailData.enable"
              inline-prompt
              :active-value="1"
              :inactive-value="2"
              disabled
            />
          </el-form-item>
        </div>
        <el-form-item label="用户角色" class="bold-label">
          <el-cascader
            v-model="detailData.authorityIds"
            style="width: 100%"
            :options="authOptions"
            :show-all-levels="false"
            :props="{
              multiple: true,
              checkStrictly: true,
              label: 'authorityName',
              value: 'authorityId',
              disabled: 'disabled',
              emitPath: false
            }"
            :clearable="false"
            disabled
          />
        </el-form-item>
        <div class="form-row">
          <el-form-item label="头像" label-width="80px" class="form-half bold-label">
            <SelectImage v-model="detailData.headerImg" disabled />
          </el-form-item>
          <el-form-item label="配置" class="form-half bold-label">
            <el-input
              v-model="detailData.originSetting"
              type="textarea"
              placeholder="JSON格式的配置信息"
              :rows="10"
              disabled
            />
          </el-form-item>
        </div>
        
        <!-- 操作者相关信息 -->
        <div class="form-row">
          <el-form-item label="操作者" class="form-half bold-label">
            <el-input v-model="detailData.operatorName" disabled />
          </el-form-item>
          <el-form-item label="操作者ID" class="form-half bold-label">
            <div class="flex items-center">
              <el-input v-model="detailData.operatorId" disabled style="margin-right: 8px; min-width: 200px;" />
              <el-button 
                type="text" 
                size="small" 
                @click="copyToClipboard(detailData.operatorId, '操作者ID已复制')"
                title="复制操作者ID"
              >
                <el-icon><copy-document /></el-icon>
              </el-button>
            </div>
          </el-form-item>
        </div>
        
        <!-- 时间相关字段 -->
        <div class="form-row">
          <el-form-item label="更新时间" class="form-half bold-label">
            <el-input :value="formatDate(detailData.updatedAt)" disabled />
          </el-form-item>
          <el-form-item label="创建时间" class="form-half bold-label">
            <el-input :value="formatDate(detailData.createdAt)" disabled />
          </el-form-item>
        </div>
      </el-form>
    </el-drawer>
    
    <el-drawer
      v-model="addUserDialog"
      :size="appStore.drawerSize || '50%'"
      :show-close="false"
      :close-on-press-escape="false"
      :close-on-click-modal="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">用户</span>
          <div>
            <el-button @click="closeAddUserDialog">取 消</el-button>
            <el-button type="primary" @click="enterAddUserDialog"
              >确 定</el-button
            >
          </div>
        </div>
      </template>

      <el-form
        ref="userForm"
        :rules="rules"
        :model="userInfo"
        label-width="80px"
      >
        <!-- 统一显示ID和UUID字段,仅在编辑模式下显示 -->
        <div v-if="dialogFlag === 'edit'" class="form-row">
          <div class="flex items-start space-x-4">
            <el-form-item label="ID" class="id-field-small flex-1">
              <div class="flex items-center">
                <el-input v-model="userInfo.id" disabled style="margin-right: 8px; min-width: 150px;" />
                <el-button 
                  type="text" 
                  size="small" 
                  @click="copyToClipboard(userInfo.id, 'ID已复制')"
                  title="复制ID"
                >
                  <el-icon><copy-document /></el-icon>
                </el-button>
              </div>
            </el-form-item>
            <el-form-item label="UUID" class="uuid-field-full flex-1">
              <div class="flex items-center">
                <el-input v-model="userInfo.uuid" disabled style="margin-right: 8px; min-width: 350px;" />
                <el-button 
                  type="text" 
                  size="small" 
                  @click="copyToClipboard(userInfo.uuid, 'UUID已复制')"
                  title="复制UUID"
                >
                  <el-icon><copy-document /></el-icon>
                </el-button>
              </div>
            </el-form-item>
          </div>
        </div>
        
        <!-- 可编辑字段 -->
        <div class="form-row">
          <el-form-item
            label="用户名"
            prop="userName"
            class="form-half bold-label"
          >
            <el-input v-model="userInfo.userName" />
          </el-form-item>
          <el-form-item label="昵称" prop="nickName" class="form-half bold-label">
            <el-input v-model="userInfo.nickName" />
          </el-form-item>
        </div>
        <el-form-item v-if="dialogFlag === 'add'" label="密码" prop="password" class="form-row bold-label">
          <el-input v-model="userInfo.password" />
        </el-form-item>
        <div class="form-row">
          <el-form-item label="姓名" prop="name" class="form-half bold-label">
            <el-input v-model="userInfo.name" />
          </el-form-item>
          <el-form-item label="邮箱" prop="email" class="form-half bold-label">
            <el-input v-model="userInfo.email" />
          </el-form-item>
        </div>
        <div class="form-row">
          <el-form-item label="手机号" prop="phone" class="form-half bold-label">
            <el-input v-model="userInfo.phone" />
          </el-form-item>
          <el-form-item label="状态" prop="disabled" class="form-half bold-label">
            <el-switch
              v-model="userInfo.enable"
              inline-prompt
              :active-value="1"
              :inactive-value="2"
              :disabled="dialogFlag === 'edit'"
            />
          </el-form-item>
        </div>
        <el-form-item label="用户角色" prop="authorityId" class="bold-label">
          <el-cascader
            v-model="userInfo.authorityIds"
            style="width: 100%"
            :options="authOptions"
            :show-all-levels="false"
            :props="{
              multiple: true,
              checkStrictly: true,
              label: 'authorityName',
              value: 'authorityId',
              disabled: 'disabled',
              emitPath: false
            }"
            :clearable="false"
          />
        </el-form-item>
        <div class="form-row">
          <el-form-item label="头像" label-width="80px" class="form-half bold-label">
            <SelectImage v-model="userInfo.headerImg" />
          </el-form-item>
          <el-form-item label="配置" prop="originSetting" class="form-half bold-label">
            <el-input
              v-model="userInfo.originSetting"
              type="textarea"
              placeholder="JSON格式的配置信息"
              :rows="10"
            />
          </el-form-item>
        </div>
        
        <!-- 操作者相关信息 -->
        <div v-if="dialogFlag === 'edit'" class="form-row">
          <el-form-item label="操作者" prop="operatorName" class="form-half bold-label">
            <el-input v-model="userInfo.operatorName" disabled />
          </el-form-item>
          <el-form-item label="操作者ID" prop="operatorId" class="form-half bold-label">
            <div class="flex items-center">
              <el-input v-model="userInfo.operatorId" disabled style="margin-right: 8px; min-width: 200px;" />
              <el-button 
                type="text" 
                size="small" 
                @click="copyToClipboard(userInfo.operatorId, '操作者ID已复制')"
                title="复制操作者ID"
              >
                <el-icon><copy-document /></el-icon>
              </el-button>
            </div>
          </el-form-item>
        </div>
        
        <!-- 将时间相关的只读字段移到表单最下方,并统一显示方式 -->
        <div v-if="dialogFlag === 'edit'" class="form-row">
          <el-form-item label="更新时间" prop="updatedAt" class="form-half bold-label">
            <el-input :value="formatDate(userInfo.updatedAt)" disabled />
          </el-form-item>
          <el-form-item label="创建时间" prop="createdAt" class="form-half bold-label">
            <el-input :value="formatDate(userInfo.createdAt)" disabled />
          </el-form-item>
        </div>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
  import {
    getUserList,
    setUserAuthorities,
    register,
    deleteUser
  } from '@/api/user'

  import { getAuthorityList } from '@/api/authority'
  import CustomPic from '@/components/customPic/index.vue'
  import WarningBar from '@/components/warningBar/warningBar.vue'
  import { setUserInfo, resetPassword } from '@/api/user.js'

  import { nextTick, ref, watch } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import SelectImage from '@/components/selectImage/selectImage.vue'
  import { useAppStore } from "@/pinia";
  
  // 日期格式化函数
  const formatDate = (dateString) => {
    if (!dateString) return '-';
    const date = new Date(dateString);
    if (isNaN(date.getTime())) return '-';
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    }).replace(/\//g, '-');
  };

  defineOptions({
    name: 'User'
  })

  const appStore = useAppStore()

  const searchInfo = ref({
    username: '',
    nickname: '',
    phone: '',
    email: ''
  })

  const onSubmit = () => {
    page.value = 1
    getTableData()
  }

  const onReset = () => {
    searchInfo.value = {
      username: '',
      nickname: '',
      phone: '',
      email: ''
    }
    getTableData()
  }
  // 初始化相关
  const setAuthorityOptions = (AuthorityData, optionsData) => {
    AuthorityData &&
      AuthorityData.forEach((item) => {
        if (item.children && item.children.length) {
          const option = {
            authorityId: item.authorityId,
            authorityName: item.authorityName,
            children: []
          }
          setAuthorityOptions(item.children, option.children)
          optionsData.push(option)
        } else {
          const option = {
            authorityId: item.authorityId,
            authorityName: item.authorityName
          }
          optionsData.push(option)
        }
      })
  }

  const page = ref(1)
  const total = ref(0)
  const pageSize = ref(10)
  const tableData = ref([])
  // 分页
  const handleSizeChange = (val) => {
    pageSize.value = val
    getTableData()
  }

  const handleCurrentChange = (val) => {
    page.value = val
    getTableData()
  }

  // 查询
  const getTableData = async () => {
    const table = await getUserList({
      page: page.value,
      pageSize: pageSize.value,
      ...searchInfo.value
    })
    if (table.code === 0) {
      tableData.value = table.data.list.map(user => {
        // 将 authority 字段转换为 authorities 数组
        if (user.authority && user.authority.authorityName) {
          return {
            ...user,
            authorities: [
              {
                authorityId: user.authority.authorityId,
                authorityName: user.authority.authorityName
              }
            ]
          }
        }
        return user
      })
      total.value = table.data.total
      page.value = table.data.page
      pageSize.value = table.data.pageSize
    }
  }

  watch(
    () => tableData.value,
    () => {
      setAuthorityIds()
    }
  )

  const initPage = async () => {
    getTableData()
    const res = await getAuthorityList()
    setOptions(res.data)
  }

  initPage()

  // 重置密码对话框相关
  const resetPwdDialog = ref(false)
  const resetPwdForm = ref(null)
  const resetPwdInfo = ref({
    ID: '',
    userName: '',
    nickName: '',
    password: ''
  })
  
  // 详情抽屉相关
  const detailDialogVisible = ref(false)
  const detailForm = ref(null)
  const detailData = ref({
    id: '',
    uuid: '',
    updatedAt: '',
    deletedAt: '',
    userName: '',
    password: '',
    nickName: '',
    name: '',
    headerImg: '',
    authorityId: '',
    authorityIds: [],
    enable: 1,
    originSetting: '',
    operatorName: '',
    operatorId: '',
    createdAt: ''
  })
  
  const openDetailDialog = (row) => {
    // 设置详情数据
    detailData.value = JSON.parse(JSON.stringify(row))
    // 处理authorityIds字段
    if (row.authorities && Array.isArray(row.authorities)) {
      detailData.value.authorityIds = row.authorities.map(i => i.authorityId)
    }
    // 处理originSetting字段,如果是对象则转换为格式化的JSON字符串
    if (detailData.value.originSetting && typeof detailData.value.originSetting === 'object') {
      detailData.value.originSetting = JSON.stringify(detailData.value.originSetting, null, 2)
    }
    // 显示抽屉
    detailDialogVisible.value = true
  }
  
  // 生成随机密码
  const generateRandomPassword = () => {
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*'
    let password = ''
    for (let i = 0; i < 12; i++) {
      password += chars.charAt(Math.floor(Math.random() * chars.length))
    }
    resetPwdInfo.value.password = password
    // 复制到剪贴板
    navigator.clipboard.writeText(password).then(() => {
      ElMessage({
        type: 'success',
        message: '密码已复制到剪贴板'
      })
    }).catch(() => {
      ElMessage({
        type: 'error',
        message: '复制失败,请手动复制'
      })
    })
  }
  
  // 打开重置密码对话框
  const resetPasswordFunc = (row) => {
    resetPwdInfo.value.ID = row.ID
    resetPwdInfo.value.userName = row.userName
    resetPwdInfo.value.nickName = row.nickName
    resetPwdInfo.value.password = ''
    resetPwdDialog.value = true
  }
  
  // 确认重置密码
  const confirmResetPassword = async () => {
    if (!resetPwdInfo.value.password) {
      ElMessage({
        type: 'warning',
        message: '请输入或生成密码'
      })
      return
    }
    
    const res = await resetPassword({
      ID: resetPwdInfo.value.ID,
      password: resetPwdInfo.value.password
    })
    
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: res.msg || '密码重置成功'
      })
      resetPwdDialog.value = false
    } else {
      ElMessage({
        type: 'error',
        message: res.msg || '密码重置失败'
      })
    }
  }
  
  // 关闭重置密码对话框
  const closeResetPwdDialog = () => {
    resetPwdInfo.value.password = ''
    resetPwdDialog.value = false
  }
  const setAuthorityIds = () => {
    tableData.value &&
      tableData.value.forEach((user) => {
        user.authorityIds =
          user.authorities &&
          user.authorities.map((i) => {
            return i.authorityId
          })
      })
  }

  const authOptions = ref([])
  const setOptions = (authData) => {
    authOptions.value = []
    setAuthorityOptions(authData, authOptions.value)
  }

  const deleteUserFunc = async (row) => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      const res = await deleteUser({ id: row.ID })
      if (res.code === 0) {
        ElMessage.success('删除成功')
        await getTableData()
      }
    })
  }

  // 复制文本到剪贴板
  const copyToClipboard = (text, message = '复制成功') => {
    navigator.clipboard.writeText(text).then(() => {
      ElMessage.success(message)
    }).catch(() => {
      ElMessage.error('复制失败')
    })
  }

  // 弹窗相关
  const userInfo = ref({
    id: '',
    uuid: '',
    updatedAt: '',
    deletedAt: '',
    userName: '',
    password: '',
    nickName: '',
    name: '',
    headerImg: '',
    authorityId: '',
    authorityIds: [],
    enable: 1, // 默认启用
    originSetting: '',
    operatorName: '',
    operatorId: '' // 修改为小写id,与后端保持一致
  })

  const rules = ref({
    userName: [
      { required: true, message: '请输入用户名', trigger: 'blur' },
      { min: 5, message: '最低5位字符', trigger: 'blur' }
    ],
    password: [
      { required: true, message: '请输入用户密码', trigger: 'blur' },
      { min: 6, message: '最低6位字符', trigger: 'blur' }
    ],
    nickName: [{ required: false, message: '请输入用户昵称', trigger: 'blur' }],
    name: [{ required: true, message: '请输入用户姓名', trigger: 'blur' }],
    phone: [
      { required: true, message: '请输入手机号', trigger: 'blur' },
      {
        pattern: /^1([38][0-9]|4[014-9]|[59][0-35-9]|6[2567]|7[0-8])\d{8}$/,
        message: '请输入合法手机号',
        trigger: 'blur'
      }
    ],
    email: [
      {
        pattern: /^([0-9A-Za-z\-_.]+)@([0-9a-z]+\.[a-z]{2,3}(\.[a-z]{2})?)$/g,
        message: '请输入正确的邮箱',
        trigger: 'blur'
      }
    ],
    authorityId: [
      { required: true, message: '请选择用户角色', trigger: 'blur' }
    ]
  })
  const userForm = ref(null)
  const enterAddUserDialog = async () => {
    userInfo.value.authorityId = userInfo.value.authorityIds[0]
    userForm.value.validate(async (valid) => {
      if (valid) {
        const req = {
          ...userInfo.value
        }
        if (dialogFlag.value === 'add') {
          const res = await register(req)
          if (res.code === 0) {
            ElMessage({ type: 'success', message: '创建成功' })
            await getTableData()
            closeAddUserDialog()
          }
        }
        if (dialogFlag.value === 'edit') {
          const res = await setUserInfo(req)
          if (res.code === 0) {
            ElMessage({ type: 'success', message: '编辑成功' })
            await getTableData()
            closeAddUserDialog()
          }
        }
      }
    })
  }

  const addUserDialog = ref(false)
  const closeAddUserDialog = () => {
    userForm.value.resetFields()
    userInfo.value.headerImg = ''
    userInfo.value.authorityIds = []
    addUserDialog.value = false
  }

  const dialogFlag = ref('add')

  const addUser = () => {
    dialogFlag.value = 'add'
    // 清除所有默认值
    userInfo.value = {
      ID: '',
      uuid: '',
      updatedAt: '',
      deletedAt: '',
      userName: '',
      password: '',
      nickName: '',
      name: '',
      headerImg: '',
      authorityId: '',
      authorityIds: [],
      enable: 1, // 默认启用
      originSetting: '',
      operatorName: '',
      operatorId: '' // 修改为小写id,与后端保持一致
    }
    addUserDialog.value = true
  }

  const tempAuth = {}
  const changeAuthority = async (row, flag, removeAuth) => {
    if (flag) {
      if (!removeAuth) {
        tempAuth[row.ID] = [...row.authorityIds]
      }
      return
    }
    await nextTick()
    const res = await setUserAuthorities({
      ID: row.ID,
      authorityIds: row.authorityIds
    })
    if (res.code === 0) {
      ElMessage({ type: 'success', message: '角色设置成功' })
    } else {
      if (!removeAuth) {
        row.authorityIds = [...tempAuth[row.ID]]
        delete tempAuth[row.ID]
      } else {
        row.authorityIds = [removeAuth, ...row.authorityIds]
      }
    }
  }

  const openEdit = (row) => {
    dialogFlag.value = 'edit'
    userInfo.value = JSON.parse(JSON.stringify(row))
    // 处理originSetting字段,如果是对象则转换为格式化的JSON字符串
    if (userInfo.value.originSetting && typeof userInfo.value.originSetting === 'object') {
      userInfo.value.originSetting = JSON.stringify(userInfo.value.originSetting, null, 2)
    }
    addUserDialog.value = true
  }

  const switchEnable = async (row) => {
    userInfo.value = JSON.parse(JSON.stringify(row))
    await nextTick()
    const req = {
      ...userInfo.value
    }
    const res = await setUserInfo(req)
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: `${req.enable === 2 ? '禁用' : '启用'}成功`
      })
      await getTableData()
      userInfo.value.headerImg = ''
      userInfo.value.authorityIds = []
    }
  }
</script>

<style lang="scss">
  .header-img-box {
    @apply w-52 h-52 border border-solid border-gray-300 rounded-xl flex justify-center items-center cursor-pointer;
  }
  .form-row {
    @apply flex flex-wrap;
  }
  .form-half {
    @apply w-[48%] mr-2 mb-4;
  }
  .id-field-small {
    @apply w-[30%] !important;
  }
  .uuid-field-full {
    @apply mb-4;
    width: calc(100% + 150px);
    flex: none;
  }
  .bold-label .el-form-item__label {
    font-weight: bold;
  }
</style>
