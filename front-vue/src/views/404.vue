<template>
	<div class="not-found-container">
	  <div class="not-found-content">
		<!-- 使用Ant Design Vue的图标组件 -->
		<a-icon class="not-found-icon" type="frown" />
		<h1 class="not-found-title">404</h1>
		<p class="not-found-message">抱歉，您访问的页面不存在</p>
		<a-button type="primary" class="back-home-btn" href="/">
		  返回首页
		</a-button>
	  </div>
	</div>
  </template>
  
  <script setup>
  // Vue 3 的 script setup 语法，不需要额外的内容
  </script>
  
  <style scoped>
  .not-found-container {
	display: flex;
	justify-content: center;
	align-items: center;
	height: 100vh;
	background-color: #f7f7f7;
  }
  
  .not-found-content {
	text-align: center;
	background-color: #fff;
	padding: 30px 50px;
	border-radius: 10px;
	box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
  }
  
  .not-found-icon {
	font-size: 100px;
	color: #ff4d4f;
  }
  
  .not-found-title {
	font-size: 50px;
	font-weight: bold;
	margin: 20px 0;
  }
  
  .not-found-message {
	font-size: 18px;
	color: #888;
	margin-bottom: 30px;
  }
  
  .back-home-btn {
	font-size: 16px;
  }
  </style>
  