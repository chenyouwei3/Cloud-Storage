<script setup>
import { UserOutlined, LockOutlined } from '@ant-design/icons-vue'
</script>

<template>
  <div id="content">
    <div style="height: 36px; line-height: 36px; text-align: center">
      <img src="@/assets/logo.png" alt="logo" style="width: 28px; height: 28px" />&nbsp;
      <span style="font-size: 24px; font-weight: 600">个人云盘</span>
    </div>
    <br />
    <a-form>
      <a-form-item>
        <a-input size="large" placeholder="用户名">
          <template #prefix>
            <UserOutlined />
          </template>
        </a-input>
      </a-form-item>
      <a-form-item>
        <a-input-password size="large" placeholder="密码">
          <template #prefix>
            <LockOutlined />
          </template>
        </a-input-password>
      </a-form-item>
      <a-form-item>
        <a-button size="large" type="primary" style="width: 100%">登录</a-button>
      </a-form-item>
    </a-form>
  </div>
</template>




<style>
#content {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 500px;
  height: 300px;
  margin-top: -150px;
  margin-left: -250px;
  box-shadow: 0 1px 20px rgba(0, 0, 0, 0.2);
  border-radius: 4px;
  padding: 20px 60px;
}
</style>
