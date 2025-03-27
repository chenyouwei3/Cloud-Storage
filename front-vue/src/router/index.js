import { createRouter, createWebHashHistory, createWebHistory } from "vue-router"

const routes = [
	//默认目录
	{
		path:'/',
		redirect:'/cloud_storage/login'
	},
	//页面路由
	{
		name:"login",
		path:"/cloud_storage/login",
		meta:{
			title:"登录页面"
		},
		component: () => import("@/views/login.vue")
	},
	{
		name:"register",
		path:"/cloud_storage/register",
		meta:{
			title:"注册页面"
		},
		component: () => import("@/views/register.vue")
	},
	{
		name:"404",
		path:"/404",
		meta:{
			title:"404"
		},
		component: () => import("@/views/404.vue")
	},
	// 捕获所有未匹配的路由，重定向到 404 页面
	{
		path: "/:pathMatch(.*)*",
		redirect: "/404"
	}
]

const router = createRouter({
	//使用url的#符号之后的部分模拟url路径的变化,因为不会触发页面刷新,所以不需要服务端支持
	//history: createWebHashHistory(), 
	history: createWebHistory(),
	routes
})

export default router