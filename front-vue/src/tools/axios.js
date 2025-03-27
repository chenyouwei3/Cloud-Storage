const VueAxios = {
    //app(vue3实例) install(传入的axios实例)
    install(app, instance) {
        //是否传入
      if (!instance) {
        console.error('You have to install axios')
        return
      }
  
      //在全局属性当中注册
      app.config.globalProperties.axios = instance
      app.config.globalProperties.$http = instance
    }
  }
  
  //导出vueAxios组件
  export { VueAxios }
  