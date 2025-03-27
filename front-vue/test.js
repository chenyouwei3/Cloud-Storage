import { createApp } from 'vue' // 导入 Vue 3 的 `createApp` 方法，用于创建 Vue 应用实例
import App from './App.vue' // 导入根组件 App.vue，它是整个应用的入口组件
import Antd from 'ant-design-vue' // 导入 Ant Design Vue 库，提供丰富的 UI 组件
import router from './router' // 导入路由配置，管理页面跳转
import storage from 'store' // 导入存储库，用于存储本地数据，比如 Token
import 'ant-design-vue/dist/antd.css' // 导入 Ant Design Vue 的样式文件
import VueClipboard from 'vue-clipboard2' // 导入剪贴板插件，允许操作剪贴板
import moment from 'moment' // 导入 moment.js 库，用于日期和时间的处理

moment.locale('zh-cn') // 设置 moment.js 的语言环境为中文

const allowList = ['login', 'fileshared'] // 定义免登录的路由列表
const loginRoutePath = '/filecloud/login' // 定义登录页面的路径

// 设置全局路由守卫，beforeEach 是每次路由跳转前都会执行的钩子函数
router.beforeEach((to, from, next) => {
  // 如果目标路由有 meta.title，动态修改网页的标题
  if (to.meta.title) {
    document.title = to.meta.title
  }

  // 从本地存储获取 Access-Token，判断用户是否已经登录
  const token = storage.get('Access-Token')
  if (token) {
    // 如果 token 存在，说明用户已登录，直接放行
    next()
  } else {
    // 如果没有 token，进行免登录路由判断
    if (allowList.includes(to.name)) {
      // 如果目标路由在免登录列表中，直接放行
      next()
    } else {
      // 否则跳转到登录页面，并传递当前路由路径，登录后可以跳回去
      next({ path: loginRoutePath, query: { redirect: to.fullPath } })
    }
  }
})

// 创建 Vue 应用实例，传入根组件 App.vue
const app = createApp(App)

// 使用插件
app.use(router) // 注册路由插件
app.use(Antd) // 注册 Ant Design Vue 插件
app.use(VueClipboard) // 注册剪贴板插件

// 挂载 Vue 应用到页面中的 #app 元素
app.mount('#app')
