import axios from 'axios'
import { VueAxios } from './axios'
import message from 'ant-design-vue/es/message'
import storage from 'store'

// 创建 axios 实例
const request = axios.create({
  // baseURL 需要从全局 `window.g.baseUrl` 读取
  baseURL: window.g.baseUrl,
  timeout: 6000 // 请求超时时间 6000ms
})

// 异常处理器
const errorHandler = (error) => {
  if (error.response) {
    if (error.response.status === 401) { 
      // 401 表示未授权，可能是 token 失效
      message.error('授权验证失败')
      const token = storage.get("Access-Token")
      
      if (token) {
        // 移除无效的 token
        storage.remove("Access-Token")
      }

      // 1 秒后刷新页面，触发登录
      setTimeout(() => {
        window.location.reload()
      }, 1000)
    } else {
      // 其他错误直接提示
      message.error(error.response.statusText)
    }
  }
  return Promise.reject(error) // 让 Promise 链终止并抛出错误
}

// 请求拦截器
request.interceptors.request.use(config => {
  const token = storage.get("Access-Token")
  if (token) {
    // 让每个请求都带上 token
    config.headers["Access-Token"] = token
  }
  return config
}, errorHandler)

// 响应拦截器
request.interceptors.response.use(response => {
  if ('success' in response.data) { // 判断是否有 success 字段
    if (!response.data.success) { 
      message.error(response.data.message)
      return Promise.reject(response)
    } else {
      return response.data.data
    }
  } else {
    return response // 直接返回原始响应
  }
}, errorHandler)

// Vue 3 安装插件的方式
const installer = {
  install(app) {
    app.use(VueAxios, request) // 让 VueAxios 安装到 Vue 应用中
  }
}

export default request

export {
  installer as VueAxios,  // 让 VueAxios 插件可以 `app.use(VueAxios)`
  request as axios  // 让外部模块可以 `import { axios }`
}
