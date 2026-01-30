<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <div class="login-header">
          <h2>11111</h2>
        </div>
      </template>
      <el-form
        :model="loginForm"
        :rules="loginRules"
        ref="loginFormRef"
        label-width="80px"
        class="login-form"
      >
        <el-form-item label="账号" prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="请输入账号"
            prefix-icon="el-icon-user"
          ></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            prefix-icon="el-icon-lock"
            show-password
          ></el-input>
        </el-form-item>
        <el-form-item>
          <el-checkbox v-model="loginForm.remember">记住密码</el-checkbox>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleLogin" class="login-button">登录</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

const emit = defineEmits(['login-success'])

// 登录表单
const loginForm = reactive({
  username: '',
  password: '',
  remember: false
})

// 表单验证规则
const loginRules = {
  username: [
    { required: true, message: '请输入账号', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
}

// 表单引用
const loginFormRef = ref()

// 登录方法
const handleLogin = () => {
  if (loginFormRef.value) {
    loginFormRef.value.validate(async (valid: boolean) => {
      if (valid) {
        try {
          const response = await fetch('http://localhost:8081/api/login', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify({
              username: loginForm.username,
              password: loginForm.password
            })
          })

          const data = await response.json()

          if (data.success) {
            // 保存登录信息到本地存储
            if (loginForm.remember) {
              localStorage.setItem('username', loginForm.username)
              localStorage.setItem('password', loginForm.password)
              localStorage.setItem('remember', 'true')
            } else {
              localStorage.removeItem('username')
              localStorage.removeItem('password')
              localStorage.removeItem('remember')
            }
            localStorage.setItem('token', data.token)
            
            // 保存客户端信息
            if (data.clientInfo) {
              localStorage.setItem('clientInfo', JSON.stringify(data.clientInfo))
              console.log('客户端信息已保存:', data.clientInfo)
            }
            
            ElMessage.success('登录成功')
            
            // 建立WebSocket连接，用于远程桌面数据传输
            try {
              const clientInfo = data.clientInfo || { clientId: 'client_' + Math.random().toString(36).substr(2, 9), machineId: 'machine_' + Math.random().toString(36).substr(2, 9) }
              const wsUrl = `ws://localhost:8081/ws?clientId=${clientInfo.clientId}&machineId=${clientInfo.machineId}`
              const ws = new WebSocket(wsUrl)
              
              ws.onopen = () => {
                console.log('WebSocket连接已建立')
                localStorage.setItem('wsConnected', 'true')
              }
              
              ws.onclose = () => {
                console.log('WebSocket连接已关闭')
                localStorage.removeItem('wsConnected')
              }
              
              ws.onerror = (error) => {
                console.error('WebSocket连接错误:', error)
                ElMessage.warning('远程桌面服务连接失败，部分功能可能不可用')
              }
              
              // 保存WebSocket连接到全局变量
              window.wsConnection = ws
            } catch (error) {
              console.error('建立WebSocket连接失败:', error)
              ElMessage.warning('远程桌面服务连接失败，部分功能可能不可用')
            }
            
            emit('login-success')
          } else {
            ElMessage.error(data.message || '登录失败')
          }
        } catch (error) {
          ElMessage.error('服务器连接失败，请检查服务器是否启动')
          console.error('Login error:', error)
        }
      } else {
        ElMessage.error('请填写完整的登录信息')
        return false
      }
    })
  }
}

// 组件挂载时，从本地存储加载登录信息
onMounted(() => {
  const remember = localStorage.getItem('remember') === 'true'
  if (remember) {
    loginForm.username = localStorage.getItem('username') || ''
    loginForm.password = localStorage.getItem('password') || ''
    loginForm.remember = true
  }
})
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f5f7fa;
}

.login-card {
  width: 400px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.login-header {
  text-align: center;
}

.login-form {
  margin-top: 20px;
}

.login-button {
  width: 100%;
}
</style>