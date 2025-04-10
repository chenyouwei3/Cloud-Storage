<template>
  <div class="help-center-page">
    <!-- 顶部动画区域 -->
    <div class="top-section">
      <div class="clouds">
        <div class="cloud cloud1"></div>
        <div class="cloud cloud2"></div>
        <div class="cloud cloud3"></div>
        <div class="cloud cloud4"></div>
      </div>
      <div class="animation-container">
        <img src="@/assets/logo.png" alt="logo" class="logo" />
        <div class="title-container">
          <div class="title-wrapper">
            <h1 class="main-title">帮助中心</h1>
            <div class="author-tag">
              <span class="author-name">常见问题解答</span>
            </div>
          </div>
          <p class="sub-title">为您提供全面的使用指南和问题解决方案</p>
        </div>
      </div>
    </div>

    <!-- 中间内容区域 -->
    <div class="middle-section">
      <div class="help-content">
        <div class="help-left">
          <img src="@/assets/title.png" alt="title" class="title-image" />
        </div>
        <div class="help-divider"></div>
        <div class="help-right">
          <div class="faq-section">
            <h2>常见问题</h2>
            <a-collapse>
              <a-collapse-panel header="如何注册账号？" key="1">
                <p>1. 点击首页的"注册"按钮</p>
                <p>2. 填写邮箱和密码</p>
                <p>3. 点击注册完成账号创建</p>
              </a-collapse-panel>
              <a-collapse-panel header="如何上传文件？" key="2">
                <p>1. 登录后进入个人空间</p>
                <p>2. 点击"上传"按钮</p>
                <p>3. 选择要上传的文件</p>
                <p>4. 等待上传完成</p>
              </a-collapse-panel>
              <a-collapse-panel header="如何分享文件？" key="3">
                <p>1. 选择要分享的文件</p>
                <p>2. 点击"分享"按钮</p>
                <p>3. 设置分享权限和有效期</p>
                <p>4. 生成分享链接</p>
              </a-collapse-panel>
              <a-collapse-panel header="如何恢复已删除的文件？" key="4">
                <p>1. 进入回收站</p>
                <p>2. 找到要恢复的文件</p>
                <p>3. 点击"恢复"按钮</p>
              </a-collapse-panel>
            </a-collapse>
            <a-button 
              type="primary" 
              size="large" 
              class="login-button" 
              href="/"
              style="margin-top: 30px; width: 100%; height: 40px; background: #00a1d6; border: none; border-radius: 4px;"
            >
              返回首页
            </a-button>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部信息区域 -->
    <div class="bottom-section">
      <div class="bottom-content">
        <div class="developer-link">
          <a href="https://github.com/chenyouwei3" target="_blank">个人开发者作品</a>
        </div>
        <div class="footer-links">
          <a href="/help-center">帮助中心</a>
          <a href="/privacy-policy">隐私政策</a>
          <a href="/terms-of-service">服务条款</a>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
// Vue 3 的 script setup 语法
</script>

<style scoped>
.help-center-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: #f0f2f5;
}

/* 顶部动画区域样式 */
.top-section {
  position: relative;
  height: 300px;
  background: linear-gradient(135deg, #1890ff 0%, #36cfc9 100%);
  overflow: hidden;
}

.clouds {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

.cloud {
  position: absolute;
  background: rgba(255, 255, 255, 0.8);
  border-radius: 50%;
  animation: float 15s infinite linear;
}

.cloud1 {
  width: 100px;
  height: 60px;
  top: 20%;
  left: 10%;
  animation-duration: 20s;
}

.cloud2 {
  width: 150px;
  height: 90px;
  top: 40%;
  left: 30%;
  animation-duration: 25s;
}

.cloud3 {
  width: 120px;
  height: 70px;
  top: 60%;
  left: 50%;
  animation-duration: 18s;
}

.cloud4 {
  width: 80px;
  height: 50px;
  top: 30%;
  left: 70%;
  animation-duration: 22s;
}

@keyframes float {
  0% {
    transform: translateX(-100%);
  }
  100% {
    transform: translateX(100vw);
  }
}

.animation-container {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-top: 40px;
}

.logo {
  width: 80px;
  height: 80px;
  margin-bottom: 20px;
}

.title-container {
  text-align: center;
  color: white;
}

.main-title {
  font-size: 48px;
  font-weight: bold;
  margin: 0;
}

.author-tag {
  margin-top: 10px;
  font-size: 16px;
}

.author-name {
  color: rgba(255, 255, 255, 0.8);
}

.sub-title {
  margin-top: 10px;
  font-size: 18px;
  color: rgba(255, 255, 255, 0.8);
}

/* 中间内容区域样式 */
.middle-section {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 40px 20px;
}

.help-content {
  display: flex;
  background-color: white;
  border-radius: 10px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  width: 1000px;
  min-height: 500px;
}

.help-left {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 40px;
}

.title-image {
  max-width: 100%;
  max-height: 400px;
}

.help-divider {
  width: 1px;
  background-color: #f0f0f0;
  margin: 40px 0;
}

.help-right {
  flex: 2;
  padding: 40px;
}

.faq-section {
  margin-bottom: 30px;
}

.faq-section h2 {
  font-size: 24px;
  color: #1890ff;
  margin-bottom: 20px;
}

/* 底部信息区域样式 */
.bottom-section {
  background-color: white;
  padding: 20px;
  text-align: center;
}

.bottom-content {
  max-width: 1200px;
  margin: 0 auto;
}

.developer-link {
  margin-bottom: 20px;
}

.developer-link a {
  color: #666;
  text-decoration: none;
  font-size: 16px;
}

.developer-link a:hover {
  color: #1890ff;
}

.footer-links {
  margin-top: 10px;
}

.footer-links a {
  color: #666;
  margin: 0 10px;
  text-decoration: none;
}

.footer-links a:hover {
  color: #1890ff;
}

.login-button {
  background: linear-gradient(135deg, #1890ff 0%, #36cfc9 100%);
  border: none;
  border-radius: 4px;
  color: white;
  font-weight: 500;
  transition: all 0.3s;
}

.login-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(24, 144, 255, 0.4);
}
</style> 