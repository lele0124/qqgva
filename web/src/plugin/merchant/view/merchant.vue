
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
          <el-form-item label="地址" prop="address">
            <el-input v-model="searchInfo.address" placeholder="搜索条件" />
          </el-form-item>
          
          <el-form-item label="状态" prop="status">
            <el-input v-model="searchInfo.status" placeholder="搜索条件" />
          </el-form-item>
        </template>
        
        <el-form-item>
          <el-button type="primary" icon="Search" @click="onSubmit">查询</el-button>
          <el-button icon="Refresh" @click="onReset">重置</el-button>
          <el-button link type="primary" @click="showAllQuery = !showAllQuery">
            {{ showAllQuery ? '收起' : '展开' }}
            <el-icon>
              <!-- 使用计算属性或方法来决定显示哪个图标 -->
              <component :is="showAllQuery ? 'ArrowUp' : 'ArrowDown'" />
            </el-icon>
          </el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="Plus" @click="openDialog('create')">新增</el-button>
        <el-button icon="Delete" :disabled="!selectedIds.length" @click="onDeleteBatch">删除</el-button>
      </div>
      <el-table
        ref="multipleTable"
        :data="tableData"
        style="width: 100%"
        tooltip-effect="dark"
        row-key="ID"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="ID" prop="ID" width="60" />
        <el-table-column align="left" label="商户名称" prop="merchantName" min-width="120" />
        <el-table-column align="left" label="商户类型" prop="merchantType" width="80">
          <template #default="scope">
            <el-tag :type="scope.row.merchantType === 1 ? 'primary' : 'success'">
              {{ formatMerchantType(scope.row.merchantType) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="法人代表" prop="legalPerson" width="100" />
        <el-table-column align="left" label="营业执照号" prop="businessLicense" min-width="150" />
        <el-table-column align="left" label="商户等级" prop="merchantLevel" width="100">
          <template #default="scope">
            <el-tag :type="formatMerchantLevelType(scope.row.merchantLevel)">
              {{ formatMerchantLevel(scope.row.merchantLevel) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="状态" prop="isEnabled" width="80">
          <template #default="scope">
            <el-tag :type="scope.row.isEnabled ? 'success' : 'danger'">
              {{ scope.row.isEnabled ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="创建时间" prop="CreatedAt" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" width="180">
          <template #default="scope">
            <el-button link type="primary" @click="openDetailDialog(scope.row)">详情</el-button>
            <el-button link type="primary" @click="openDialog('edit', scope.row)">编辑</el-button>
            <el-button link type="danger" @click="onDelete(scope.row)">删除</el-button>
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
          @current-change="handlePageChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>
    
    <!-- 表单弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="60%"
      destroy-on-close
    >
      <MerchantForm ref="merchantFormRef" :type="formType" :data="formData" @submit="handleSubmit" @cancel="dialogVisible = false" />
    </el-dialog>
    
    <!-- 详情弹窗 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="商户详情"
      width="60%"
      destroy-on-close
    >
      <MerchantDetail :data="detailData" />
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import MerchantForm from '@/plugin/merchant/form/merchant.vue'
import MerchantDetail from '@/plugin/merchant/view/detail.vue'
import { useMerchantStore } from '@/plugin/merchant/store/merchant'
import { formatTimeToStr } from '@/utils/date'
import { processSearchData } from '@/plugin/merchant/utils/dataProcessor'

// 引入图标组件
import { ArrowUp, ArrowDown, QuestionFilled } from '@element-plus/icons-vue'

const merchantStore = useMerchantStore()

// 响应式数据
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

const showAllQuery = ref(false)
const elSearchFormRef = ref()
const multipleTable = ref()
const dialogVisible = ref(false)
const detailDialogVisible = ref(false)
const formType = ref('create')
const formData = ref({})
const detailData = ref({})

// 表格相关数据
const tableData = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const selectedIds = ref([])

// 组件引用
const merchantFormRef = ref()

// 计算属性
const dialogTitle = computed(() => {
  return formType.value === 'create' ? '创建商户' : '编辑商户'
})

// 格式化函数
const formatDate = (date) => {
  return date ? formatTimeToStr(date, 'yyyy-MM-dd HH:mm:ss') : ''
}

const formatMerchantType = (type) => {
  const types = { 1: '企业', 2: '个体' }
  return types[type] || '未知'
}

const formatMerchantLevel = (level) => {
  const levels = { 1: '普通商户', 2: '高级商户', 3: 'VIP商户' }
  return levels[level] || '未知'
}

const formatMerchantLevelType = (level) => {
  const types = { 1: 'info', 2: 'warning', 3: 'danger' }
  return types[level] || 'info'
}

// 方法
const onSubmit = () => {
  page.value = 1
  getTableData()
}

const onReset = () => {
  elSearchFormRef.value?.resetFields()
  page.value = 1
  getTableData()
}

const handleSelectionChange = (val) => {
  selectedIds.value = val.map(item => item.ID)
}

const handlePageChange = (val) => {
  page.value = val
  getTableData()
}

const handleSizeChange = (val) => {
  pageSize.value = val
  page.value = 1
  getTableData()
}

const openDialog = (type, row) => {
  formType.value = type
  if (type === 'edit' && row) {
    formData.value = { ...row }
  } else {
    formData.value = {}
  }
  dialogVisible.value = true
}

const openDetailDialog = (row) => {
  detailData.value = { ...row }
  detailDialogVisible.value = true
}

const handleSubmit = () => {
  dialogVisible.value = false
  getTableData()
}

const onDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该商户吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    const res = await merchantStore.deleteMerchant(row.ID)
    if (res) {
      ElMessage.success('删除成功')
      getTableData()
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败: ' + error.message)
    } else {
      ElMessage.info('已取消删除')
    }
  }
}

const onDeleteBatch = async () => {
  try {
    await ElMessageBox.confirm(`确定要删除选中的 ${selectedIds.value.length} 个商户吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    const res = await merchantStore.deleteMerchantBatch(selectedIds.value)
    if (res) {
      ElMessage.success('删除成功')
      getTableData()
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败: ' + error.message)
    } else {
      ElMessage.info('已取消删除')
    }
  }
}

// 获取表格数据
const getTableData = async () => {
  // 准备搜索参数
  const params = {
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo
  }
  
  // 使用统一的处理函数处理搜索数据
  const processedParams = processSearchData(params)
  
  // 处理时间范围
  if (searchInfo.updatedAtRange && searchInfo.updatedAtRange.length === 2) {
    processedParams.updatedAtRange = [
      searchInfo.updatedAtRange[0],
      searchInfo.updatedAtRange[1]
    ]
  }
  
  try {
    const res = await merchantStore.fetchMerchantList(processedParams)
    if (res.code === 0) {
      tableData.value = res.data.list
      total.value = res.data.total
    }
  } catch (error) {
    ElMessage.error('获取数据失败: ' + error.message)
  }
}

// 初始化
onMounted(() => {
  getTableData()
})
</script>

<style scoped>
.gva-search-box {
  padding: 20px;
  background-color: #fff;
  border-radius: 4px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  margin-bottom: 20px;
}

.gva-table-box {
  padding: 20px;
  background-color: #fff;
  border-radius: 4px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.gva-btn-list {
  margin-bottom: 20px;
}

.gva-pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
