<template>
  <div class="privacy-policy-page">
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
            <h1 class="main-title">隐私政策</h1>
            <div class="author-tag">
              <span class="author-name">保护您的数据安全</span>
            </div>
          </div>
          <p class="sub-title">我们重视您的隐私，致力于保护您的个人信息</p>
        </div>
      </div>
    </div>

    <!-- 中间内容区域 -->
    <div class="middle-section">
      <div class="privacy-content">
        <div class="privacy-left">
          <img src="@/assets/title.png" alt="title" class="title-image" />
        </div>
        <div class="privacy-divider"></div>
        <div class="privacy-right">
          <div class="policy-section">
            <h2>隐私政策</h2>
            <div class="policy-content">
              <h3>1. 信息收集</h3>
              <p>我们收集的信息包括：</p>
              <ul>
                <li>注册时提供的个人信息（邮箱、用户名等）</li>
                <li>上传的文件信息</li>
                <li>使用服务的相关数据</li>
              </ul>

              <h3>2. 信息使用</h3>
              <p>我们使用收集的信息用于：</p>
              <ul>
                <li>提供和改进服务</li>
                <li>处理您的请求</li>
                <li>发送服务通知</li>
                <li>保障账户安全</li>
              </ul>

              <h3>3. 信息保护</h3>
              <p>我们采取以下措施保护您的信息：</p>
              <ul>
                <li>数据加密存储</li>
                <li>访问权限控制</li>
                <li>定期安全审计</li>
                <li>安全漏洞修复</li>
              </ul>

              <h3>4. 信息共享</h3>
              <p>我们不会在未经您同意的情况下：</p>
              <ul>
                <li>出售您的个人信息</li>
                <li>与第三方共享您的数据</li>
                <li>用于商业推广</li>
              </ul>

              <h3>5. 您的权利</h3>
              <p>您拥有以下权利：</p>
              <ul>
                <li>访问您的个人信息</li>
                <li>更正不准确的信息</li>
                <li>删除您的账户</li>
                <li>导出您的数据</li>
              </ul>
            </div>
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
.privacy-policy-page {
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

.privacy-content {
  display: flex;
  background-color: white;
  border-radius: 10px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  width: 1000px;
  min-height: 500px;
}

.privacy-left {
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

.privacy-divider {
  width: 1px;
  background-color: #f0f0f0;
  margin: 40px 0;
}

.privacy-right {
  flex: 2;
  padding: 40px;
}

.policy-section {
  margin-bottom: 30px;
}

.policy-section h2 {
  font-size: 24px;
  color: #1890ff;
  margin-bottom: 20px;
}

.policy-content h3 {
  font-size: 18px;
  color: #333;
  margin: 20px 0 10px;
}

.policy-content p {
  color: #666;
  line-height: 1.6;
  margin-bottom: 10px;
}

.policy-content ul {
  color: #666;
  line-height: 1.6;
  padding-left: 20px;
  margin-bottom: 20px;
}

.policy-content li {
  margin-bottom: 5px;
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
</style> 