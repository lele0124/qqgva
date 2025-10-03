<template>
  <div class="merchant-detail">
    <div class="detail-header">
      <h2>商户详情</h2>
      <el-button @click="goBack" type="default">返回列表</el-button>
    </div>
    
    <div class="detail-content" v-if="loading">
      <el-empty description="加载中..." />
    </div>
    
    <div class="detail-content" v-else-if="merchantInfo">
      <el-descriptions :column="2" border>
        <!-- 基本信息区域 -->
        <el-descriptions-item label="商户名称" span="2">
          <span class="text-bold">{{ merchantInfo.merchantName }}</span>
        </el-descriptions-item>
        
        <el-descriptions-item label="商户图标" :span="2">
          <img 
            v-if="merchantInfo.merchantIcon" 
            :src="merchantInfo.merchantIcon" 
            alt="商户图标"
            class="merchant-icon"
          />
          <span v-else class="text-gray">暂无图标</span>
        </el-descriptions-item>
        
        <el-descriptions-item label="商户类型">
          <el-tag type="info">
            {{ merchantInfo.merchantType === 1 ? '企业' : merchantInfo.merchantType === 2 ? '个体' : '未知' }}
          </el-tag>
        </el-descriptions-item>
        
        <el-descriptions-item label="父商户ID">
          {{ merchantInfo.parentID || '-' }}
        </el-descriptions-item>
        
        <el-descriptions-item label="营业执照号">
          {{ merchantInfo.businessLicense || '-' }}
        </el-descriptions-item>
        
        <el-descriptions-item label="法人代表">
          {{ merchantInfo.legalPerson || '-' }}
        </el-descriptions-item>
        
        <!-- 详细信息区域 -->
        <el-descriptions-item label="注册地址" :span="2">
          <div class="text-break">{{ merchantInfo.registeredAddress || '-' }}</div>
        </el-descriptions-item>
        
        <el-descriptions-item label="经营范围" :span="2">
          <div class="text-break">{{ merchantInfo.businessScope || '-' }}</div>
        </el-descriptions-item>
        
        <el-descriptions-item label="商户状态">
          <el-tag :type="merchantInfo.isEnabled ? 'success' : 'danger'">
            {{ merchantInfo.isEnabled ? '正常' : '关闭' }}
          </el-tag>
        </el-descriptions-item>
        
        <el-descriptions-item label="商户等级">
          <el-tag :type="getLevelTagType(merchantInfo.merchantLevel)">
            {{ getLevelText(merchantInfo.merchantLevel) }}
          </el-tag>
        </el-descriptions-item>
        
        <el-descriptions-item label="有效开始时间">
          {{ formatDate(merchantInfo.validStartTime) || '-' }}
        </el-descriptions-item>
        
        <el-descriptions-item label="有效结束时间">
          {{ formatDate(merchantInfo.validEndTime) || '-' }}
        </el-descriptions-item>
        
        <!-- 系统信息区域 -->
        <el-descriptions-item label="创建时间">
          {{ formatDate(merchantInfo.createdAt) || '-' }}
        </el-descriptions-item>
        
        <el-descriptions-item label="更新时间">
          {{ formatDate(merchantInfo.updatedAt) || '-' }}
        </el-descriptions-item>
        
        <el-descriptions-item label="创建人">
          {{ merchantInfo.createBy || '-' }}
        </el-descriptions-item>
        
        <el-descriptions-item label="更新人">
          {{ merchantInfo.updateBy || '-' }}
        </el-descriptions-item>
      </el-descriptions>
    </div>
    
    <div class="detail-content" v-else>
      <el-empty description="暂无数据" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { findMerchant } from '@/plugin/merchant/api/merchant'
import { formatDate } from '@/utils/format'
import { ElMessage } from 'element-plus'
import { useMerchantStore } from '@/plugin/merchant/store/merchant'

const route = useRoute()
const router = useRouter()
const merchantStore = useMerchantStore()

// 商户详情数据
const merchantInfo = ref({})
// 加载状态
const loading = ref(false)

// 获取商户详情
const getMerchantDetail = async () => {
  const id = route.params.id
  if (!id) {
    ElMessage.error('缺少商户ID')
    return
  }
  
  loading.value = true
  try {
    // 尝试从Store获取数据，避免重复请求
    if (merchantStore.detailForm && merchantStore.detailForm.ID === id) {
      merchantInfo.value = merchantStore.detailForm
    } else {
      // 直接调用API获取详情
      const res = await findMerchant({ ID: id })
      if (res.code === 0) {
        merchantInfo.value = res.data
        // 更新Store中的详情数据
        merchantStore.detailForm = res.data
      } else {
        ElMessage.error(res.msg || '获取商户详情失败')
      }
    }
  } catch (error) {
    ElMessage.error('获取商户详情失败')
    console.error('获取商户详情失败:', error)
  } finally {
    loading.value = false
  }
}

// 获取商户等级标签类型
const getLevelTagType = (level) => {
  switch (level) {
    case 3:
      return 'primary'
    case 2:
      return 'success'
    default:
      return 'info'
  }
}

// 获取商户等级文本
const getLevelText = (level) => {
  switch (level) {
    case 1:
      return '普通商户'
    case 2:
      return '高级商户'
    case 3:
      return 'VIP商户'
    default:
      return '未知'
  }
}

// 返回列表页
const goBack = () => {
  // 检查是否有回退历史，如果没有则直接跳转到列表页
  if (router.options.history.state.back) {
    router.back()
  } else {
    router.push('/layout/merchant')
  }
}

// 组件挂载时获取详情
onMounted(() => {
  getMerchantDetail()
})
</script>

<style scoped>
.merchant-detail {
  padding: 20px;
  background-color: #fff;
  min-height: 100%;
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid #ebeef5;
}

.detail-header h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 500;
  color: #303133;
}

.detail-content {
  padding: 20px;
  background-color: #fff;
}

.merchant-icon {
  width: 100px;
  height: 100px;
  object-fit: cover;
  border-radius: 8px;
  border: 1px solid #ebeef5;
}

.text-bold {
  font-weight: 500;
  font-size: 16px;
  color: #303133;
}

.text-gray {
  color: #909399;
}

.text-break {
  word-break: break-all;
}

/* 响应式布局调整 */
@media (max-width: 768px) {
  .detail-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .el-descriptions {
    --el-descriptions-column: 1 !important;
  }
}
</style>