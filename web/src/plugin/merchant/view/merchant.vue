
<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
        <el-form-item label="更新日期" prop="updatedAtRange">
          <template #label>
            <span>
              更新日期
              <el-tooltip content="搜索范围是开始日期(包含)至结束日期(不包含)">
                <el-icon><QuestionFilled /></el-icon>
              </el-tooltip>
            </span>
          </template>
          <el-date-picker
            v-model="searchInfo.updatedAtRange"
            class="!w-380px"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
          />
        </el-form-item>
        
        <el-form-item label="商户名称" prop="merchantName">
          <el-input v-model="searchInfo.merchantName" placeholder="搜索条件" />
        </el-form-item>
        
        <el-form-item label="商户图标URL" prop="merchantIcon">
          <el-input v-model="searchInfo.merchantIcon" placeholder="搜索条件" />
        </el-form-item>
        
        <el-form-item label="商户类型" prop="merchantType">
          <el-select v-model="searchInfo.merchantType" placeholder="搜索条件">
            <el-option label="企业" value="1" />
            <el-option label="个体" value="2" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="父商户ID" prop="parentID">
          <el-input v-model.number="searchInfo.parentID" placeholder="搜索条件" />
        </el-form-item>
        
        <el-form-item label="营业执照号" prop="businessLicense">
          <el-input v-model="searchInfo.businessLicense" placeholder="搜索条件" />
        </el-form-item>
        
        <el-form-item label="法人代表" prop="legalPerson">
          <el-input v-model="searchInfo.legalPerson" placeholder="搜索条件" />
        </el-form-item>
        
        <el-form-item label="注册地址" prop="registeredAddress">
          <el-input v-model="searchInfo.registeredAddress" placeholder="搜索条件" />
        </el-form-item>
        
        <el-form-item label="经营范围" prop="businessScope">
          <el-input v-model="searchInfo.businessScope" placeholder="搜索条件" />
        </el-form-item>
        
        <el-form-item label="商户状态" prop="isEnabled">
          <el-select v-model="searchInfo.isEnabled" placeholder="搜索条件" clearable>
            <el-option label="启用" value="1" />
            <el-option label="禁用" value="0" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="商户等级" prop="merchantLevel">
          <el-select v-model="searchInfo.merchantLevel" placeholder="搜索条件" clearable>
            <el-option label="普通商户" value="1" />
            <el-option label="高级商户" value="2" />
            <el-option label="VIP商户" value="3" />
          </el-select>
        </el-form-item>
        
        <template v-if="showAllQuery">
          <el-form-item label="有效开始时间" prop="validStartTime">
            <el-date-picker
              v-model="searchInfo.validStartTime"
              type="datetime"
              placeholder="开始时间"
              class="!w-280px"
            />
          </el-form-item>
          
          <el-form-item label="有效结束时间" prop="validEndTime">
            <el-date-picker
              v-model="searchInfo.validEndTime"
              type="datetime"
              placeholder="结束时间"
              class="!w-280px"
            />
          </el-form-item>
          
          <el-form-item label="地址" prop="address">
            <el-input v-model="searchInfo.address" placeholder="搜索条件" />
          </el-form-item>
          
          <el-form-item label="状态" prop="status">
            <el-input v-model="searchInfo.status" placeholder="搜索条件" />
          </el-form-item>
        </template>
        
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
          <el-button link type="primary" @click="showAllQuery = !showAllQuery">
            {{ showAllQuery ? '收起' : '展开' }}
            <el-icon>{{ showAllQuery ? 'arrow-up' : 'arrow-down' }}</el-icon>
          </el-button>
        </el-form-item>
      </el-form>
    </div>
    
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button v-auth="btnAuth.add" type="primary" icon="plus" @click="openDialog">新增</el-button>
        <el-button v-auth="btnAuth.batchDelete" icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length" @click="onDelete">删除</el-button>
        <ExportTemplate v-auth="btnAuth.exportTemplate" template-id="merchant_Merchant" />
        <ExportExcel v-auth="btnAuth.exportExcel" template-id="merchant_Merchant" filterDeleted/>
        <ImportExcel v-auth="btnAuth.importExcel" template-id="merchant_Merchant" @on-success="getTableData" />
      </div>
      
      <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
        @sort-change="sortChange"
      >
        <el-table-column type="selection" width="55" />
        
        <el-table-column align="left" label="商户ID" prop="ID" width="120" />

        <el-table-column align="left" label="商户名称" prop="merchantName" width="120" />

        <el-table-column align="left" label="商户图标" prop="merchantIcon" width="120">
          <template #default="scope">
            <img v-if="scope.row.merchantIcon" :src="scope.row.merchantIcon" style="width: 40px; height: 40px; object-fit: cover;" />
            <span v-else>无</span>
          </template>
        </el-table-column>

        <el-table-column align="left" label="商户类型" prop="merchantType" width="120">
          <template #default="scope">
            <span>{{ scope.row.merchantType === 1 ? '企业' : scope.row.merchantType === 2 ? '个体' : '未知' }}</span>
          </template>
        </el-table-column>

        <el-table-column align="left" label="父商户ID" prop="parentID" width="120" />

        <el-table-column align="left" label="营业执照号" prop="businessLicense" width="150" />

        <el-table-column align="left" label="法人代表" prop="legalPerson" width="120" />

        <el-table-column align="left" label="注册地址" prop="registeredAddress" width="200" show-overflow-tooltip />

        <el-table-column align="left" label="经营范围" prop="businessScope" width="150" show-overflow-tooltip />

        <el-table-column sortable align="left" label="状态" prop="isEnabled" width="120">
          <template #default="scope">
            <el-tag :type="scope.row.isEnabled ? 'success' : 'danger'">
              {{ scope.row.isEnabled ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column align="left" label="商户等级" prop="merchantLevel" width="120">
          <template #default="scope">
            <el-tag :type="scope.row.merchantLevel === 3 ? 'primary' : scope.row.merchantLevel === 2 ? 'success' : 'info'">
              {{ scope.row.merchantLevel === 1 ? '普通商户' : scope.row.merchantLevel === 2 ? '高级商户' : scope.row.merchantLevel === 3 ? 'VIP商户' : '未知' }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column align="left" label="有效时间" width="200" v-if="showAllQuery">
          <template #default="scope">
            <div>
              <div>开始: {{ formatDate(scope.row.validStartTime) }}</div>
              <div>结束: {{ formatDate(scope.row.validEndTime) }}</div>
            </div>
          </template>
        </el-table-column>

        <el-table-column align="left" label="地址" prop="address" width="120" />

        <el-table-column sortable align="left" label="更新时间" prop="updatedAt" width="180">
          <template #default="scope">{{ formatDate(scope.row.updatedAt) }}</template>
        </el-table-column>
        
        <el-table-column align="left" label="操作" fixed="right" min-width="240">
      <template #default="scope">
        <el-button v-auth="btnAuth.info" type="primary" link class="table-button" @click="goToDetail(scope.row)"><el-icon style="margin-right: 5px"><InfoFilled /></el-icon>查看</el-button>
        <el-button v-auth="btnAuth.edit" type="primary" link icon="edit" class="table-button" @click="updateMerchantFunc(scope.row)">编辑</el-button>
        <el-button v-auth="btnAuth.delete" type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
      </template>
    </el-table-column>
      </el-table>
      
      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>
    
    <el-drawer destroy-on-close size="800" v-model="dialogFormVisible" :show-close="false" :before-close="closeDialog">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{type==='create'?'新增':'编辑'}}</span>
          <div>
            <el-button :loading="btnLoading" type="primary" @click="enterDialog">确 定</el-button>
            <el-button @click="closeDialog">取 消</el-button>
          </div>
        </div>
      </template>

      <el-form :model="formData" label-position="top" ref="elFormRef" :rules="rule" label-width="80px">
        <el-form-item label="商户名称:" prop="merchantName">
          <el-input v-model="formData.merchantName" :clearable="true" placeholder="请输入商户名称" />
        </el-form-item>
        <el-form-item label="商户图标URL:" prop="merchantIcon">
          <el-input v-model="formData.merchantIcon" :clearable="true" placeholder="请输入商户图标URL" />
        </el-form-item>
        <el-form-item label="商户类型:" prop="merchantType">
          <el-select v-model="formData.merchantType" placeholder="请选择商户类型">
            <el-option label="企业" value="1" />
            <el-option label="个体" value="2" />
          </el-select>
        </el-form-item>
        <el-form-item label="父商户ID:" prop="parentID">
          <el-input v-model.number="formData.parentID" :clearable="true" placeholder="请输入父商户ID" />
        </el-form-item>
        <el-form-item label="营业执照号:" prop="businessLicense">
          <el-input v-model="formData.businessLicense" :clearable="true" placeholder="请输入营业执照号" />
        </el-form-item>
        <el-form-item label="法人代表:" prop="legalPerson">
          <el-input v-model="formData.legalPerson" :clearable="true" placeholder="请输入法人代表" />
        </el-form-item>
        <el-form-item label="注册地址:" prop="registeredAddress">
          <el-input v-model="formData.registeredAddress" :clearable="true" placeholder="请输入注册地址" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="经营范围:" prop="businessScope">
          <el-input v-model="formData.businessScope" :clearable="true" placeholder="请输入经营范围" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="商户开关状态:" prop="isEnabled">
          <el-select v-model="formData.isEnabled" placeholder="请选择商户状态">
            <el-option label="正常" value="1" />
            <el-option label="关闭" value="0" />
          </el-select>
        </el-form-item>
        <el-form-item label="有效开始时间:" prop="validStartTime">
          <el-date-picker
            v-model="formData.validStartTime"
            type="datetime"
            placeholder="选择开始时间"
            class="!w-280px"
          />
        </el-form-item>
        <el-form-item label="有效结束时间:" prop="validEndTime">
          <el-date-picker
            v-model="formData.validEndTime"
            type="datetime"
            placeholder="选择结束时间"
            class="!w-280px"
          />
        </el-form-item>
        <el-form-item label="商户等级:" prop="merchantLevel">
          <el-select v-model="formData.merchantLevel" placeholder="请选择商户等级">
            <el-option label="普通商户" value="1" />
            <el-option label="高级商户" value="2" />
            <el-option label="VIP商户" value="3" />
          </el-select>
        </el-form-item>
      </el-form>
    </el-drawer>

    <!-- 详情页面已独立到专门的详情路由 -->
  </div>
</template>

<script setup>
import { 
  createMerchant, 
  deleteMerchant, 
  deleteMerchantByIds, 
  updateMerchant, 
  findMerchant, 
  getMerchantList 
} from '@/plugin/merchant/api/merchant'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict ,filterDataSource, returnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
// 引入按钮权限标识
import { useBtnAuth } from '@/utils/btnAuth'
// 引入路由
import { useRouter } from 'vue-router'

// 导出组件
import ExportExcel from '@/components/exportExcel/exportExcel.vue'
// 导入组件
import ImportExcel from '@/components/exportExcel/importExcel.vue'
// 导出模板组件
import ExportTemplate from '@/components/exportExcel/exportTemplate.vue'

// 获取按钮权限
const btnAuth = useBtnAuth()

// 表格数据
const tableData = ref([])
// 分页数据
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
// 按钮加载状态
const btnLoading = ref(false)
// 搜索表单
const elSearchFormRef = ref()
const showAllQuery = ref(false)
const searchInfo = reactive({
  merchantName: '',
  merchantIcon: '',
  merchantType: '',
  parentID: '',
  businessLicense: '',
  legalPerson: '',
  registeredAddress: '',
  businessScope: '',
  isEnabled: '',
  validStartTime: '',
  validEndTime: '',
  merchantLevel: '',
  address: '',
  status: '',
  updatedAtRange: []
})

// 表单数据
const elFormRef = ref()
const dialogFormVisible = ref(false)
const detailShow = ref(false)
const type = ref('create')
const formData = reactive({
  ID: '',
  merchantName: '',
  merchantIcon: '',
  merchantType: '',
  parentID: '',
  businessLicense: '',
  legalPerson: '',
  registeredAddress: '',
  businessScope: '',
  isEnabled: true,
  validStartTime: '',
  validEndTime: '',
  merchantLevel: '1',
  address: ''
})

// 详情表单数据
const detailForm = ref({})

// 表格选中数据
const multipleSelection = ref([])
const multipleTable = ref()

// 表单验证规则
const rule = {
  merchantName: [
    {
      required: true,
      message: '请输入商户名称',
      trigger: 'blur'
    },
    {
      pattern: /^\S+$/,  // 不允许有空格
      message: '商户名称不能包含空格',
      trigger: 'blur'
    }
  ],
  merchantType: [
    {
      required: true,
      message: '请选择商户类型',
      trigger: 'change'
    }
  ],
  isEnabled: [
    {
      required: true,
      message: '请选择商户状态',
      trigger: 'change'
    }
  ],
  merchantLevel: [
    {
      required: true,
      message: '请选择商户等级',
      trigger: 'change'
    }
  ]
}

// 获取表格数据
const getTableData = async () => {
  try {
    // 创建提交数据的副本，进行类型转换
    const submitData = {
      page: page.value,
      pageSize: pageSize.value,
      ...searchInfo
    }
    
    // 统一处理数值类型转换
    const processValue = (value, type) => {
      if (value === '' || value === null || value === undefined) return undefined
      
      switch (type) {
        case 'int':
          return parseInt(value)
        case 'bool':
          return value === '1'
        default:
          return value
      }
    }
    
    // 转换数值类型
    submitData.merchantType = processValue(searchInfo.merchantType, 'int')
    submitData.parentID = processValue(searchInfo.parentID, 'int')
    submitData.merchantLevel = processValue(searchInfo.merchantLevel, 'int')
    submitData.isEnabled = processValue(searchInfo.isEnabled, 'bool')
    
    const res = await getMerchantList(submitData)
    if (res.code === 0) {
      tableData.value = res.data.list || []
      total.value = res.data.total || 0
    } else {
      ElMessage.error(res.msg || '获取数据失败')
    }
  } catch (error) {
    ElMessage.error('获取数据失败')
  }
}

// 查询
const onSubmit = () => {
  page.value = 1
  getTableData()
}

// 重置
const onReset = () => {
  if (elSearchFormRef.value) {
    elSearchFormRef.value.resetFields()
  }
  // 清空自定义重置的值
  Object.keys(searchInfo).forEach(key => {
    searchInfo[key] = ''
  })
  searchInfo.updatedAtRange = []
  page.value = 1
  getTableData()
}

// 分页
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

const handleSizeChange = (val) => {
  pageSize.value = val
  page.value = 1
  getTableData()
}

// 表格排序
const sortChange = (obj) => {
  if (obj.order === 'ascending') {
    searchInfo.sort = 'asc'
    searchInfo.field = obj.prop
  } else if (obj.order === 'descending') {
    searchInfo.sort = 'desc'
    searchInfo.field = obj.prop
  } else {
    searchInfo.sort = ''
    searchInfo.field = ''
  }
  getTableData()
}

// 表格选择
const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

// 打开创建弹窗
const openDialog = () => {
  type.value = 'create'
  // 重置表单
    Object.keys(formData).forEach(key => {
      if (key === 'isEnabled') {
        formData[key] = true
      } else if (key === 'merchantLevel') {
        formData[key] = '1'
      } else if (key === 'status') {
        formData[key] = 1
      } else {
        formData[key] = ''
      }
    })
  dialogFormVisible.value = true
}

// 打开编辑弹窗
const updateMerchantFunc = async (row) => {
  try {
    type.value = 'edit'
    const res = await findMerchant({ ID: row.ID })
    if (res.code === 0) {
      // 赋值
      Object.keys(formData).forEach(key => {
        formData[key] = res.data[key] !== undefined ? res.data[key] : ''
      })
      dialogFormVisible.value = true
    } else {
      ElMessage.error(res.msg || '获取详情失败')
    }
  } catch (error) {
    ElMessage.error('获取详情失败')
  }
}

// 提交表单
const enterDialog = async () => {
  if (!elFormRef.value) {
    return
  }
  try {
    await elFormRef.value.validate()
    btnLoading.value = true
    let res
    
    // 创建提交数据的副本，进行类型转换
    const submitData = { ...formData }
    
    // 统一处理数值类型转换
    const processValue = (value, type) => {
      if (value === '' || value === null || value === undefined) {
        switch (type) {
          case 'int':
            return 0
          case 'bool':
            return false
          default:
            return ''
        }
      }
      
      switch (type) {
        case 'int':
          return parseInt(value)
        case 'bool':
          return value === '1' || value === true
        default:
          return value
      }
    }
    
    // 转换数值类型
    submitData.merchantType = processValue(formData.merchantType, 'int')
    submitData.parentID = processValue(formData.parentID, 'int')
    submitData.merchantLevel = processValue(formData.merchantLevel, 'int')
    submitData.isEnabled = processValue(formData.isEnabled, 'bool')
    
    if (type.value === 'create') {
      res = await createMerchant(submitData)
    } else {
      res = await updateMerchant(submitData)
    }
    if (res.code === 0) {
      ElMessage.success('操作成功')
      dialogFormVisible.value = false
      getTableData()
    } else {
      ElMessage.error(res.msg || '操作失败')
    }
  } catch (error) {
    ElMessage.error('操作失败')
  } finally {
    btnLoading.value = false
  }
}

// 关闭弹窗
const closeDialog = () => {
  dialogFormVisible.value = false
  if (elFormRef.value) {
    elFormRef.value.resetFields()
  }
}

// 初始化路由
const router = useRouter()

// 跳转到详情页面
const goToDetail = (row) => {
  // 跳转到详情页面
  router.push({ path: `/layout/merchant/detail/${row.ID}` })
}

// 删除单条数据
const deleteRow = async (row) => {
  try {
    await ElMessageBox.confirm(`确定要删除 ${row.merchantName} 吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    const res = await deleteMerchant({ ID: row.ID })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      getTableData()
    } else {
      ElMessage.error(res.msg || '删除失败')
    }
  } catch (error) {
    if (error === 'cancel') {
      return
    }
    ElMessage.error('删除失败')
  }
}

// 批量删除
const onDelete = async () => {
  if (multipleSelection.value.length === 0) {
    ElMessage.warning('请选择要删除的数据')
    return
  }
  try {
    await ElMessageBox.confirm(`确定要删除选中的 ${multipleSelection.value.length} 条数据吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    const ids = multipleSelection.value.map(item => item.ID)
    const res = await deleteMerchantByIds({ "IDs[]": ids })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      getTableData()
      multipleSelection.value = []
    } else {
      ElMessage.error(res.msg || '删除失败')
    }
  } catch (error) {
    if (error === 'cancel') {
      return
    }
    ElMessage.error('删除失败')
  }
}

// 初始化数据
// 注意：默认排序已移至后端实现
getTableData()
</script>
