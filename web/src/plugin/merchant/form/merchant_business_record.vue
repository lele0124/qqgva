
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="商户ID:" prop="merchantId">
        <el-select  multiple  v-model="formData.merchantId" placeholder="请选择商户ID" style="width:100%" :clearable="false" >
          <el-option v-for="(item,key) in dataSource.merchantId" :key="key" :label="item.label" :value="item.value" />
        </el-select>
       </el-form-item>
        <el-form-item label="记录类型:" prop="recordType">
          <el-input v-model="formData.recordType" :clearable="true"  placeholder="请输入记录类型" />
       </el-form-item>
        <el-form-item label="金额:" prop="amount">
          <el-input-number v-model="formData.amount" :precision="2" :clearable="true"></el-input-number>
       </el-form-item>
        <el-form-item label="描述:" prop="description">
          <el-input v-model="formData.description" :clearable="true"  placeholder="请输入描述" />
       </el-form-item>
        <el-form-item label="记录时间:" prop="recordTime">
          <el-date-picker v-model="formData.recordTime" type="date" placeholder="选择日期" :clearable="true"></el-date-picker>
       </el-form-item>
        <el-form-item>
          <el-button :loading="btnLoading" type="primary" @click="save">保存</el-button>
          <el-button type="primary" @click="back">返回</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import {
    getMerchantBusinessRecordDataSource,
  createMerchantBusinessRecord,
  updateMerchantBusinessRecord,
  findMerchantBusinessRecord
} from '@/plugin/merchant/api/merchant_business_record'

defineOptions({
    name: 'MerchantBusinessRecordForm'
})

// 自动获取字典
import { getDictFunc } from '@/utils/format'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'


const route = useRoute()
const router = useRouter()

// 提交按钮loading
const btnLoading = ref(false)

const type = ref('')
const formData = ref({
            merchantId: '',
            recordType: '',
            amount: 0,
            description: '',
            recordTime: new Date(),
        })
// 验证规则
const rule = reactive({
               merchantId : [{
                   required: true,
                   message: '商户ID不能为空',
                   trigger: ['input','blur'],
               }],
               recordType : [{
                   required: true,
                   message: '记录类型不能为空',
                   trigger: ['input','blur'],
               }],
               amount : [{
                   required: true,
                   message: '金额不能为空',
                   trigger: ['input','blur'],
               }],
               recordTime : [{
                   required: true,
                   message: '记录时间不能为空',
                   trigger: ['input','blur'],
               }],
})

const elFormRef = ref()
  const dataSource = ref([])
  const getDataSourceFunc = async()=>{
    const res = await getMerchantBusinessRecordDataSource()
    if (res.code === 0) {
      dataSource.value = res.data
    }
  }
  getDataSourceFunc()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findMerchantBusinessRecord({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
}

init()
// 保存按钮
const save = async() => {
      btnLoading.value = true
      elFormRef.value?.validate( async (valid) => {
         if (!valid) return btnLoading.value = false
            let res
           switch (type.value) {
             case 'create':
               res = await createMerchantBusinessRecord(formData.value)
               break
             case 'update':
               res = await updateMerchantBusinessRecord(formData.value)
               break
             default:
               res = await createMerchantBusinessRecord(formData.value)
               break
           }
           btnLoading.value = false
           if (res.code === 0) {
             ElMessage({
               type: 'success',
               message: '创建/更改成功'
             })
           }
       })
}

// 返回按钮
const back = () => {
    router.go(-1)
}

</script>

<style>
</style>
