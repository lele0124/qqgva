import { defineAsyncComponent } from 'vue'
import Layout from '@/view/layout/index.vue'

/**
 * 商户管理模块路由配置
 * 此配置将被动态添加到系统路由中
 */
const merchantRouter = {
  path: '/merchant',
  name: 'merchant',
  component: Layout,
  meta: {
    title: '商户管理',
    icon: 'ShoppingBag',
    hidden: false,
    permissions: ['merchant:list'],
    keepAlive: false
  },
  children: [
    {
      path: '',
      name: 'MerchantList',
      component: () => import('./view/merchant.vue'),
      meta: {
        title: '商户列表',
        icon: 'List',
        hidden: false,
        permissions: ['merchant:list'],
        keepAlive: true
      }
    },
    {
      path: 'detail/:id',
      name: 'MerchantDetail',
      component: () => import('./view/detail.vue'),
      meta: {
        title: '商户详情',
        icon: 'InfoFilled',
        hidden: true, // 隐藏在菜单中，通过列表页点击进入
        permissions: ['merchant:view'],
        keepAlive: false,
        breadcrumb: [
          {
            path: '/layout/merchant',
            title: '商户列表'
          }
        ]
      }
    },
    {
      path: 'create',
      name: 'MerchantCreate',
      component: () => import('./form/merchant.vue'),
      meta: {
        title: '创建商户',
        icon: 'Plus',
        hidden: true,
        permissions: ['merchant:create'],
        keepAlive: false,
        breadcrumb: [
          {
            path: '/layout/merchant',
            title: '商户列表'
          }
        ]
      }
    },
    {
      path: 'edit/:id',
      name: 'MerchantEdit',
      component: () => import('./form/merchant.vue'),
      meta: {
        title: '编辑商户',
        icon: 'Edit',
        hidden: true,
        permissions: ['merchant:edit'],
        keepAlive: false,
        breadcrumb: [
          {
            path: '/layout/merchant',
            title: '商户列表'
          }
        ]
      }
    }
  ]
}

export default merchantRouter