
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="商户名称:" prop="merchantName">
          <el-input v-model="formData.merchantName" :clearable="true"  placeholder="请输入商户名称" />
       </el-form-item>
        <el-form-item label="联系人:" prop="contactPerson">
          <el-input v-model="formData.contactPerson" :clearable="true"  placeholder="请输入联系人" />
       </el-form-item>
        <el-form-item label="联系电话:" prop="contactPhone">
          <el-input v-model="formData.contactPhone" :clearable="true"  placeholder="请输入联系电话" />
       </el-form-item>
        <el-form-item label="商户地址:" prop="address">
          <el-input v-model="formData.address" :clearable="true"  placeholder="请输入商户地址" />
       </el-form-item>
        <el-form-item label="经营范围:" prop="businessScope">
          <el-input v-model="formData.businessScope" :clearable="true"  placeholder="请输入经营范围" />
       </el-form-item>
        <el-form-item label="是否启用:" prop="isEnabled">
          <el-switch v-model="formData.isEnabled" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>
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
  createMerchant,
  updateMerchant,
  findMerchant
} from '@/plugin/merchant/api/merchant'

defineOptions({
    name: 'MerchantForm'
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
            merchantName: '',
            contactPerson: '',
            contactPhone: '',
            address: '',
            businessScope: '',
            isEnabled: false,
        })
// 验证规则
const rule = reactive({
               merchantName : [{
                   required: true,
                   message: '商户名称不能为空',
                   trigger: ['input','blur'],
               }],
               contactPerson : [{
                   required: true,
                   message: '联系人不能为空',
                   trigger: ['input','blur'],
               }],
               contactPhone : [{
                   required: true,
                   message: '联系电话不能为空',
                   trigger: ['input','blur'],
               }],
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findMerchant({ ID: route.query.id })
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
               res = await createMerchant(formData.value)
               break
             case 'update':
               res = await updateMerchant(formData.value)
               break
             default:
               res = await createMerchant(formData.value)
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
