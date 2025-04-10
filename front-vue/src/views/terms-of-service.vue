<template>
  <div class="terms-page">
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
            <h1 class="main-title">服务条款</h1>
            <div class="author-tag">
              <span class="author-name">使用条款和条件</span>
            </div>
          </div>
          <p class="sub-title">请仔细阅读以下服务条款</p>
        </div>
      </div>
    </div>

    <!-- 中间内容区域 -->
    <div class="middle-section">
      <div class="terms-content">
        <div class="terms-left">
          <img src="@/assets/title.png" alt="title" class="title-image" />
        </div>
        <div class="terms-divider"></div>
        <div class="terms-right">
          <div class="terms-section">
            <h2>服务条款</h2>
            <div class="terms-text">
              <h3>1. 服务说明</h3>
              <p>本服务提供云存储功能，允许用户：</p>
              <ul>
                <li>上传、存储和管理文件</li>
                <li>分享文件给其他用户</li>
                <li>在线预览支持的文件类型</li>
                <li>管理个人账户设置</li>
              </ul>

              <h3>2. 用户责任</h3>
              <p>用户在使用服务时应当：</p>
              <ul>
                <li>遵守所有适用的法律法规</li>
                <li>不传播违法或侵权内容</li>
                <li>保护账户安全</li>
                <li>不干扰服务正常运行</li>
              </ul>

              <h3>3. 服务限制</h3>
              <p>服务存在以下限制：</p>
              <ul>
                <li>单个文件大小限制</li>
                <li>总存储空间限制</li>
                <li>每日上传次数限制</li>
                <li>支持的文件类型限制</li>
              </ul>

              <h3>4. 服务变更</h3>
              <p>我们保留以下权利：</p>
              <ul>
                <li>修改服务条款</li>
                <li>调整服务功能</li>
                <li>更新服务政策</li>
                <li>暂停或终止服务</li>
              </ul>

              <h3>5. 免责声明</h3>
              <p>以下情况我们免责：</p>
              <ul>
                <li>不可抗力导致的服务中断</li>
                <li>用户操作导致的损失</li>
                <li>第三方服务问题</li>
                <li>法律法规变更</li>
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
.terms-page {
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

.terms-content {
  display: flex;
  background-color: white;
  border-radius: 10px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  width: 1000px;
  min-height: 500px;
  max-width: 90%;
}

.terms-left {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 40px;
  min-width: 300px;
}

.title-image {
  max-width: 100%;
  max-height: 300px;
  object-fit: contain;
}

.terms-divider {
  width: 1px;
  background-color: #f0f0f0;
  margin: 40px 0;
}

.terms-right {
  flex: 2;
  padding: 40px;
  max-width: 600px;
}

.terms-section {
  margin-bottom: 30px;
}

.terms-section h2 {
  font-size: 24px;
  color: #1890ff;
  margin-bottom: 20px;
  text-align: center;
}

.terms-text {
  line-height: 1.6;
}

.terms-text h3 {
  font-size: 18px;
  color: #333;
  margin: 20px 0 10px;
}

.terms-text p {
  color: #666;
  margin-bottom: 10px;
}

.terms-text ul {
  color: #666;
  padding-left: 20px;
  margin-bottom: 20px;
}

.terms-text li {
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