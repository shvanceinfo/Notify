<template>
  <div class="main-container">
    <!-- 顶部导航栏 -->
    <div class="top-nav">
      <div class="logo-area">
        <img src="https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=RQ%20Center%20Client%20logo%20with%20purple%20color%20scheme&image_size=square" alt="RQ Center Client" class="logo-image">
        <span class="logo-text">Client</span>
        <span class="version">BETA v2024.2.9</span>
      </div>
      <div class="nav-tabs">
        <el-button type="text" :class="['nav-tab', { active: activeTab === 'control' }]" @click="switchTab('control')">控制中心</el-button>
        <el-button type="text" :class="['nav-tab', { active: activeTab === 'remote' }]" @click="switchTab('remote')">远程推送</el-button>
        <el-button type="text" :class="['nav-tab', { active: activeTab === 'command' }]" @click="switchTab('command')">命令管理</el-button>
      </div>

    </div>

    <!-- 主内容区域 -->
    <div class="content-container">
      <!-- 左侧侧边栏 -->
      <div class="sidebar">
        <div class="user-info">
          <div class="user-badge">
            <el-icon><user /></el-icon>
            <span>{{ username }}（管理员）</span>
          </div>
        </div>
        <el-menu
          :default-active="activeMenu"
          class="sidebar-menu"
          @select="handleMenuSelect"
        >
          <el-menu-item index="home">
            <el-icon><house /></el-icon>
            <span>首页</span>
          </el-menu-item>
          <el-menu-item index="control">
            <el-icon><setting /></el-icon>
            <span>控制中心</span>
          </el-menu-item>
          <el-menu-item index="app">
            <el-icon><grid /></el-icon>
            <span>应用中心</span>
          </el-menu-item>
          <el-menu-item index="tools">
            <el-icon><tools /></el-icon>
            <span>快捷工具</span>
          </el-menu-item>
          <el-menu-item index="project">
            <el-icon><document /></el-icon>
            <span>项目启动器</span>
          </el-menu-item>
        </el-menu>
      </div>

      <!-- 右侧内容区域 -->
      <div class="main-panel">
        <!-- 控制中心内容 -->
        <div v-if="activeTab === 'control'" class="control-center">
          <div class="panel-header">
            <h3>控制中心</h3>
            <el-button type="primary" @click="fetchMachines">刷新机器列表</el-button>
          </div>
          
          <div v-if="loading" class="loading">
            <el-loading :fullscreen="false" text="加载机器列表中..." />
          </div>
          
          <div v-else-if="machines.length > 0" class="machine-list">
            <el-table :data="machines" style="width: 100%">
              <el-table-column prop="name" label="机器名" width="180">
                <template #default="scope">
                  <div class="machine-name">
                    {{ scope.row.name }}
                    <el-tag v-if="isOwnMachine(scope.row)" size="small" type="info" class="own-machine-tag">
                      本机器
                    </el-tag>
                  </div>
                </template>
              </el-table-column>
              <el-table-column prop="ip" label="机器IP" width="180" />
              <el-table-column prop="online" label="状态">
                <template #default="scope">
                  <el-tag :type="scope.row.online ? 'success' : 'danger'">
                    {{ scope.row.online ? '在线' : '离线' }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="lastSeen" label="最后在线" />
              <el-table-column label="操作" width="150">
                <template #default="scope">
                  <el-button type="primary" size="small" :disabled="!scope.row.online" @click="handleRemoteConnect(scope.row)">
                    远程连接
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>
          
          <div v-else class="empty-state">
            <el-empty description="暂无机器数据" />
          </div>
        </div>

        <!-- 远程推送内容 -->
        <div v-else-if="activeTab === 'remote'" class="remote-push">
          <!-- 横幅区域 -->
          <div class="banner">
            <img src="https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=abstract%20blue%20and%20white%20curved%20shapes%20background%20for%20remote%20desktop%20application&image_size=landscape_16_9" alt="远程推送" class="banner-image">
            <div class="banner-text">
              <h3>远程推送</h3>
              <p>更方便的管理客户机，只需一键推送</p>
            </div>
          </div>

          <!-- 功能卡片区域 -->
          <div class="feature-cards">
            <el-card class="feature-card">
              <div class="card-content">
                <el-icon class="card-icon"><video-camera /></el-icon>
                <h4>推送向导</h4>
                <p>选择客户机、命令进行推送</p>
              </div>
            </el-card>
            <el-card class="feature-card">
              <div class="card-content">
                <el-icon class="card-icon"><list /></el-icon>
                <h4>推送列表</h4>
                <p>检查添加到推送列表的任务</p>
              </div>
            </el-card>
          </div>
        </div>

        <!-- 命令管理内容 -->
        <div v-else-if="activeTab === 'command'" class="command-management">
          <div class="panel-header">
            <h3>命令管理</h3>
          </div>
          <div class="empty-state">
            <el-empty description="命令管理功能开发中" />
          </div>
        </div>

        <!-- 首页内容 -->
        <div v-else-if="activeTab === 'home'" class="home">
          <div class="panel-header">
            <h3>首页</h3>
            <el-button type="danger" @click="handleLogout">退出登录</el-button>
          </div>
          <div class="empty-state">
            <el-empty description="首页功能开发中" />
          </div>
        </div>

        <!-- 应用中心内容 -->
        <div v-else-if="activeTab === 'app'" class="app-center">
          <!-- 横幅区域 -->
          <div class="banner">
            <img src="https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=app%20store%20banner%20with%20purple%20color%20scheme&image_size=landscape_16_9" alt="应用中心" class="banner-image">
            <div class="banner-text">
              <h3>应用中心</h3>
              <p>浏览和管理可用的应用程序</p>
            </div>
          </div>

          <!-- 功能卡片区域 -->
          <div class="feature-cards">
            <el-card class="feature-card">
              <div class="card-content">
                <el-icon class="card-icon"><download /></el-icon>
                <h4>应用安装</h4>
                <p>安装新的应用程序到客户端</p>
                <el-button type="primary" size="small">浏览应用</el-button>
              </div>
            </el-card>
            <el-card class="feature-card">
              <div class="card-content">
                <el-icon class="card-icon"><upload /></el-icon>
                <h4>应用更新</h4>
                <p>检查和更新已安装的应用</p>
                <el-button type="primary" size="small">检查更新</el-button>
              </div>
            </el-card>
            <el-card class="feature-card">
              <div class="card-content">
                <el-icon class="card-icon"><delete /></el-icon>
                <h4>应用卸载</h4>
                <p>卸载不需要的应用程序</p>
                <el-button type="primary" size="small">管理应用</el-button>
              </div>
            </el-card>
            <el-card class="feature-card">
              <div class="card-content">
                <el-icon class="card-icon"><star /></el-icon>
                <h4>推荐应用</h4>
                <p>查看为您推荐的应用</p>
                <el-button type="primary" size="small">查看推荐</el-button>
              </div>
            </el-card>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- 远程桌面对话框 -->
  <el-dialog
    v-model="remoteDialogVisible"
    :show-header="false"
    width="90%"
    top="10px"
    :close-on-click-modal="false"
  >
    <div class="remote-desktop-container">
      <!-- 控制栏显示/隐藏按钮 -->
      <div class="control-toggle" @click="toggleControlBar">
        <el-button type="text" class="toggle-btn">
          <el-icon><arrow-down /></el-icon>
        </el-button>
      </div>
      
      <!-- 远程桌面控制栏 -->
      <div 
        class="remote-control-bar" 
        v-if="showControlBar"
        @mouseenter="clearHideTimer"
        @mouseleave="startHideTimer"
      >
        <div class="machine-info">
          <span>机器: {{ currentMachine?.name }}</span>
          <span>IP: {{ currentMachine?.ip }}</span>
        </div>
        <div class="remote-controls">
          <el-button type="text" size="small" @click="handleRemoteDesktopScale('zoomOut')">
            <el-icon><zoom-out /></el-icon> 缩小
          </el-button>
          <el-button type="text" size="small" @click="handleRemoteDesktopScale('reset')">
            <el-icon><refresh /></el-icon> 重置
          </el-button>
          <el-button type="text" size="small" @click="handleRemoteDesktopScale('zoomIn')">
            <el-icon><zoom-in /></el-icon> 放大
          </el-button>
          <el-button type="text" size="small" @click="closeRemoteDialog">
            <el-icon><close /></el-icon> 关闭
          </el-button>
        </div>
      </div>
      
      <!-- 远程桌面显示区域 -->
      <div class="remote-desktop-display" :style="{ transform: `scale(${remoteDesktopScale})` }">
        <!-- 远程桌面视频显示 -->
        <div class="remote-desktop-screen">
          <video 
            ref="remoteVideo" 
            class="desktop-video" 
            autoplay 
            playsinline 
            muted
          ></video>
          <div class="connection-status">
            <el-tag type="success">已连接</el-tag>
            <span>{{ new Date().toLocaleTimeString() }}</span>
          </div>
        </div>
      </div>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { ZoomIn, ZoomOut, Refresh, Download, Upload, Delete, Star, ArrowDown, Close } from '@element-plus/icons-vue'

// WebRTC相关变量
let peerConnection: RTCPeerConnection | null = null
let dataChannel: RTCDataChannel | null = null
let localStream: MediaStream | null = null
let remoteStream: MediaStream | null = null
let isInitiator = false

// 屏幕捕获相关变量
let screenCaptureStream: MediaStream | null = null

// WebSocket信令连接
let ws: WebSocket | null = null
let clientID: string = ''

const username = ref('')
const activeMenu = ref('control')
const activeTab = ref('control') // 默认显示控制中心
const machines = ref<any[]>([
  {
    id: 'default_machine',
    name: '默认机器',
    ip: '127.0.0.1',
    online: true,
    lastSeen: new Date().toLocaleString()
  }
])
const loading = ref(false)
const remoteDialogVisible = ref(false)
const currentMachine = ref<any>(null)
const remoteDesktopScale = ref(1)
const showControlBar = ref(false)
const hideTimer = ref<number | null>(null)
const remoteVideo = ref<HTMLVideoElement | null>(null)
const clientInfo = ref<any>(null)

const handleLogout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('username')
  localStorage.removeItem('password')
  localStorage.removeItem('remember')
  ElMessage.success('已退出登录')
  window.location.reload()
}

const switchTab = (tab: string) => {
  activeTab.value = tab
  activeMenu.value = tab
  
  // 如果切换到控制中心，刷新机器列表
  if (tab === 'control') {
    fetchMachines()
  }
}

const handleMenuSelect = (key: string) => {
  activeMenu.value = key
  
  // 根据菜单选择切换标签页
  if (key === 'control') {
    activeTab.value = 'control'
    fetchMachines()
  } else if (key === 'remote') {
    activeTab.value = 'remote'
  } else if (key === 'command') {
    activeTab.value = 'command'
  } else if (key === 'home') {
    activeTab.value = 'home'
  } else if (key === 'app') {
    activeTab.value = 'app'
    fetchMachines()
  }
}

const fetchMachines = async () => {
  loading.value = true
  try {
    console.log('开始获取机器列表...')
    const response = await fetch('http://localhost:8081/api/machines')
    console.log('响应状态:', response.status)
    const data = await response.json()
    console.log('响应数据:', data)
    
    if (data.success) {
      console.log('机器列表数据:', data.data)
      console.log('机器列表数据长度:', data.data.length)
      machines.value = data.data
      
      // 如果机器列表为空，添加默认机器
      if (machines.value.length === 0) {
        machines.value = [{
          id: 'default_machine',
          name: '默认机器',
          ip: '127.0.0.1',
          online: true,
          lastSeen: new Date().toLocaleString()
        }]
        console.log('添加默认机器')
      }
      
      console.log('机器列表已更新:', machines.value)
      console.log('机器列表长度:', machines.value.length)
      ElMessage.success('机器列表加载成功')
    } else {
      // 如果加载失败，添加默认机器
      machines.value = [{
        id: 'default_machine',
        name: '默认机器',
        ip: '127.0.0.1',
        online: true,
        lastSeen: new Date().toLocaleString()
      }]
      ElMessage.error('加载机器列表失败: ' + data.message)
    }
  } catch (error) {
    console.error('获取机器列表失败:', error)
    // 如果网络错误，添加默认机器
    machines.value = [{
      id: 'default_machine',
      name: '默认机器',
      ip: '127.0.0.1',
      online: true,
      lastSeen: new Date().toLocaleString()
    }]
    ElMessage.error('网络错误，无法加载机器列表')
  } finally {
    loading.value = false
  }
}

const handleRemoteDesktopScale = (action: string) => {
  if (action === 'zoomIn' && remoteDesktopScale.value < 1.5) {
    remoteDesktopScale.value += 0.1
  } else if (action === 'zoomOut' && remoteDesktopScale.value > 0.5) {
    remoteDesktopScale.value -= 0.1
  } else if (action === 'reset') {
    remoteDesktopScale.value = 1
  }
}

const closeRemoteDialog = () => {
  remoteDialogVisible.value = false
  currentMachine.value = null
  
  // 关闭WebRTC连接
  if (peerConnection) {
    peerConnection.close()
    peerConnection = null
  }
  
  // 停止屏幕捕获
  if (screenCaptureStream) {
    screenCaptureStream.getTracks().forEach(track => track.stop())
    screenCaptureStream = null
  }
  
  ElMessage.info('远程连接已关闭')
}

// 初始化WebRTC连接
const initWebRTC = async () => {
  try {
    // 创建RTCPeerConnection
    peerConnection = new RTCPeerConnection({
      iceServers: [
        {
          urls: 'stun:stun.l.google.com:19302'
        }
      ]
    })
    
    // 监听ICE候选
    peerConnection.onicecandidate = (event) => {
      if (event.candidate) {
        console.log('ICE候选:', event.candidate)
        // 这里应该将ICE候选发送给对方
      }
    }
    
    // 监听远程流
    peerConnection.ontrack = (event) => {
      remoteStream = event.streams[0]
      console.log('收到远程流:', remoteStream)
      // 将远程流渲染到视频元素
      if (remoteVideo.value && remoteStream) {
        remoteVideo.value.srcObject = remoteStream
        console.log('远程流已绑定到视频元素')
        ElMessage.success('远程桌面已显示')
      } else {
        console.error('remoteVideo.value或remoteStream为null')
        ElMessage.error('无法显示远程桌面')
      }
    }
    
    // 监听连接状态变化
    peerConnection.onconnectionstatechange = () => {
      const state = peerConnection?.connectionState
      console.log('连接状态:', state)
      if (state === 'connected') {
        ElMessage.success('远程连接已建立')
      } else if (state === 'disconnected' || state === 'failed') {
        ElMessage.error('远程连接已断开')
      }
    }
    
    // 监听ICE连接状态变化
    peerConnection.oniceconnectionstatechange = () => {
      const iceState = peerConnection?.iceConnectionState
      console.log('ICE连接状态:', iceState)
    }
    
    // 监听ICE收集状态变化
    peerConnection.onicegatheringstatechange = () => {
      const gatheringState = peerConnection?.iceGatheringState
      console.log('ICE收集状态:', gatheringState)
    }
    
    // 创建数据通道
    dataChannel = peerConnection.createDataChannel('desktop-data')
    
    // 监听数据通道事件
    dataChannel.onopen = () => {
      console.log('数据通道已打开')
      ElMessage.success('远程连接已建立')
    }
    
    dataChannel.onmessage = (event) => {
      console.log('收到数据:', event.data)
      // 处理收到的数据
    }
    
    dataChannel.onclose = () => {
      console.log('数据通道已关闭')
    }
    
    dataChannel.onerror = (error) => {
      console.error('数据通道错误:', error)
    }
    
  } catch (error) {
    console.error('初始化WebRTC失败:', error)
    ElMessage.error('建立远程连接失败')
  }
}

// 开始屏幕捕获
const startScreenCapture = async () => {
  try {
    // 检查浏览器是否支持屏幕捕获
    if (!navigator.mediaDevices || !navigator.mediaDevices.getDisplayMedia) {
      console.error('浏览器不支持屏幕捕获')
      ElMessage.error('浏览器不支持屏幕捕获，请使用Chrome或Firefox浏览器')
      return
    }
    
    // 请求屏幕捕获（捕捉整个显示器的界面）
    console.log('请求捕获整个显示器的界面...')
    ElMessage.info('请在弹出的对话框中选择要捕获的显示器，然后点击"共享"按钮')
    
    screenCaptureStream = await navigator.mediaDevices.getDisplayMedia({
      video: {
        cursor: 'always',
        displaySurface: 'monitor' // 优先选择显示器作为捕获源
      },
      audio: false
    })
    
    console.log('显示器界面捕获已开始')
    ElMessage.success('显示器界面捕获已开始')
    
    // 将捕获的显示器界面流添加到peerConnection
    if (peerConnection && screenCaptureStream) {
      console.log('将显示器界面流添加到peerConnection')
      screenCaptureStream.getTracks().forEach(track => {
        peerConnection?.addTrack(track, screenCaptureStream!)
        console.log('添加了轨道:', track.kind)
      })
    } else {
      console.error('peerConnection或screenCaptureStream为null')
      ElMessage.warning('无法将显示器界面流添加到远程连接')
    }
    
    // 监听显示器捕获停止
    if (screenCaptureStream && screenCaptureStream.getVideoTracks().length > 0) {
      screenCaptureStream.getVideoTracks()[0].onended = () => {
        console.log('显示器界面捕获已停止')
        ElMessage.info('显示器界面捕获已停止')
      }
    }
    
  } catch (error) {
    console.error('显示器界面捕获失败:', error)
    ElMessage.error('显示器界面捕获失败，请检查浏览器权限设置')
  }
}

// 处理远程连接
const handleRemoteConnect = async (machine: any) => {
  currentMachine.value = machine
  remoteDesktopScale.value = 1
  remoteDialogVisible.value = true
  ElMessage.success(`正在连接到 ${machine.name} (${machine.ip})`)
  
  // 初始化WebRTC连接
  await initWebRTC()
  
  // 开始屏幕捕获
  await startScreenCapture()
  
  // 创建offer
  if (peerConnection) {
    try {
      const offer = await peerConnection.createOffer()
      await peerConnection.setLocalDescription(offer)
      console.log('创建offer:', offer)
      
      // 这里应该将offer发送给对方
      // 由于是P2P连接，这里需要通过某种方式交换SDP信息
      // 可以使用WebSocket服务器作为信令服务器
      
      // 模拟接收方的响应，创建一个answer并设置为远程描述
      // 注意：在实际应用中，这里应该通过WebSocket发送offer给对方，然后接收对方的answer
      setTimeout(async () => {
        if (peerConnection) {
          try {
            // 模拟远程peerConnection
            const remotePeerConnection = new RTCPeerConnection({
              iceServers: [
                {
                  urls: 'stun:stun.l.google.com:19302'
                }
              ]
            })
            
            // 模拟远程屏幕捕获
            let remoteScreenStream: MediaStream | null = null
            try {
              remoteScreenStream = await navigator.mediaDevices.getDisplayMedia({
                video: {
                  cursor: 'always',
                  displaySurface: 'monitor'
                },
                audio: false
              })
              console.log('模拟远程屏幕捕获已开始')
              
              // 将远程屏幕捕获流添加到远程peerConnection
              if (remoteScreenStream) {
                remoteScreenStream.getTracks().forEach(track => {
                  remotePeerConnection.addTrack(track, remoteScreenStream!)
                  console.log('远程屏幕捕获流已添加到远程peerConnection')
                })
              }
            } catch (error) {
              console.error('模拟远程屏幕捕获失败:', error)
            }
            
            // 监听本地ICE候选并添加到远程peerConnection
            peerConnection.onicecandidate = (event) => {
              if (event.candidate) {
                remotePeerConnection.addIceCandidate(event.candidate)
                console.log('本地ICE候选已添加到远程peerConnection')
              }
            }
            
            // 监听远程ICE候选并添加到本地peerConnection
            remotePeerConnection.onicecandidate = (event) => {
              if (event.candidate) {
                peerConnection.addIceCandidate(event.candidate)
                console.log('远程ICE候选已添加到本地peerConnection')
              }
            }
            
            // 监听远程流并添加到本地peerConnection的ontrack事件
            remotePeerConnection.ontrack = (event) => {
              const remoteStream = event.streams[0]
              console.log('远程peerConnection收到远程流:', remoteStream)
              // 手动触发本地peerConnection的ontrack事件
              if (peerConnection.ontrack) {
                peerConnection.ontrack({ streams: [remoteStream] } as any)
              }
            }
            
            // 设置远程描述并创建answer
            await remotePeerConnection.setRemoteDescription(offer)
            const answer = await remotePeerConnection.createAnswer()
            await remotePeerConnection.setLocalDescription(answer)
            console.log('创建answer:', answer)
            
            // 设置远程描述到本地peerConnection
            await peerConnection.setRemoteDescription(answer)
            console.log('远程描述已设置到本地peerConnection')
            
            // 模拟远程流
            if (!remoteScreenStream && navigator.mediaDevices) {
              // 如果无法获取屏幕捕获，使用摄像头作为替代
              try {
                remoteScreenStream = await navigator.mediaDevices.getUserMedia({ video: true, audio: false })
                console.log('使用摄像头作为远程流')
                if (remoteScreenStream) {
                  remoteScreenStream.getTracks().forEach(track => {
                    remotePeerConnection.addTrack(track, remoteScreenStream!)
                    console.log('摄像头流已添加到远程peerConnection')
                  })
                }
              } catch (error) {
                console.error('获取摄像头失败:', error)
              }
            }
            
          } catch (error) {
            console.error('模拟远程连接失败:', error)
            ElMessage.error('建立远程连接失败')
          }
        }
      }, 1000)
      
    } catch (error) {
      console.error('创建offer失败:', error)
      ElMessage.error('建立连接失败')
    }
  }
}

const toggleControlBar = () => {
  showControlBar.value = !showControlBar.value
  if (showControlBar.value) {
    clearHideTimer()
  }
}

const startHideTimer = () => {
  clearHideTimer()
  hideTimer.value = window.setTimeout(() => {
    showControlBar.value = false
  }, 3000) // 3秒后自动隐藏
}

const clearHideTimer = () => {
  if (hideTimer.value) {
    clearTimeout(hideTimer.value)
    hideTimer.value = null
  }
}

onMounted(() => {
  console.log('Main组件已挂载')
  const storedUsername = localStorage.getItem('username')
  if (storedUsername) {
    username.value = storedUsername
  }

  // 从本地存储加载客户端信息
  const storedClientInfo = localStorage.getItem('clientInfo')
  if (storedClientInfo) {
    clientInfo.value = JSON.parse(storedClientInfo)
    clientID = clientInfo.value.clientId
    console.log('加载客户端信息:', clientInfo.value)
  } else {
    // 生成客户端ID
    clientID = 'client_' + Math.random().toString(36).substr(2, 9)
    console.log('生成客户端ID:', clientID)
  }
  
  // 初始化WebSocket连接
  console.log('初始化WebSocket连接')
  initWebSocket()

  // 初始加载机器列表
  console.log('初始加载机器列表')
  fetchMachines()
})

// 初始化WebSocket连接
const initWebSocket = () => {
  try {
    const wsUrl = `ws://localhost:8081/ws?clientId=${clientID}&machineId=${clientID}`
    ws = new WebSocket(wsUrl)
    
    ws.onopen = () => {
      console.log('WebSocket信令连接已建立')
    }
    
    ws.onmessage = (event) => {
      console.log('收到信令消息:', event.data)
      handleSignalingMessage(JSON.parse(event.data))
    }
    
    ws.onclose = () => {
      console.log('WebSocket信令连接已关闭')
    }
    
    ws.onerror = (error) => {
      console.error('WebSocket信令连接错误:', error)
    }
  } catch (error) {
    console.error('初始化WebSocket失败:', error)
  }
}

// 处理信令消息
const handleSignalingMessage = (message: any) => {
  if (!peerConnection) return
  
  switch (message.type) {
    case 'offer':
      handleOffer(message.data)
      break
    case 'answer':
      handleAnswer(message.data)
      break
    case 'ice-candidate':
      handleICECandidate(message.data)
      break
    default:
      console.log('未知的信令消息类型:', message.type)
  }
}

// 处理offer消息
const handleOffer = async (message: any) => {
  if (!peerConnection || !message.data) return
  
  try {
    const offer = message.data
    await peerConnection.setRemoteDescription(new RTCSessionDescription(offer))
    const answer = await peerConnection.createAnswer()
    await peerConnection.setLocalDescription(answer)
    
    // 发送answer消息
    sendSignalingMessage('answer', answer, message.senderId)
  } catch (error) {
    console.error('处理offer失败:', error)
  }
}

// 处理answer消息
const handleAnswer = async (answer: RTCSessionDescriptionInit) => {
  if (!peerConnection) return
  
  try {
    await peerConnection.setRemoteDescription(new RTCSessionDescription(answer))
  } catch (error) {
    console.error('处理answer失败:', error)
  }
}

// 处理ICE候选
const handleICECandidate = async (candidate: RTCIceCandidateInit) => {
  if (!peerConnection || !candidate) return
  
  try {
    await peerConnection.addIceCandidate(new RTCIceCandidate(candidate))
  } catch (error) {
    console.error('处理ICE候选失败:', error)
  }
}

// 发送信令消息
const sendSignalingMessage = (type: string, data: any, receiverID: string) => {
  if (!ws || ws.readyState !== WebSocket.OPEN) {
    console.error('WebSocket未连接')
    return
  }
  
  const message = {
    type,
    senderId: clientID,
    receiverId: receiverID,
    data
  }
  
  ws.send(JSON.stringify(message))
  console.log('发送信令消息:', message)
}

// 判断是否是自己的机器
const isOwnMachine = (machine: any) => {
  if (!clientInfo.value) return false
  return machine.id === clientInfo.value.machineId
}
</script>

<style scoped>
.main-container {
  width: 100%;
  height: 100vh;
  background-color: #ffffff;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* 顶部导航栏 */
.top-nav {
  height: 48px;
  background-color: #ffffff;
  border-bottom: 1px solid #eaeaea;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
}

.logo-area {
  display: flex;
  align-items: center;
  gap: 10px;
}

.logo-image {
  width: 32px;
  height: 32px;
}

.logo-text {
  font-size: 16px;
  font-weight: bold;
  color: #333333;
}

.version {
  font-size: 12px;
  color: #999999;
}

.nav-tabs {
  display: flex;
  gap: 10px;
}

.nav-tab {
  padding: 8px 16px;
  border-radius: 4px;
  font-size: 14px;
  color: #666666;
}

.nav-tab.active {
  background-color: #f0f0f0;
  color: #333333;
  font-weight: 500;
}

.window-controls {
  display: flex;
  gap: 5px;
}

.window-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  color: #666666;
  border-radius: 4px;
}

.window-btn:hover {
  background-color: #f0f0f0;
}

/* 主内容区域 */
.content-container {
  flex: 1;
  display: flex;
  overflow: hidden;
}

/* 左侧侧边栏 */
.sidebar {
  width: 200px;
  background-color: #f9f9f9;
  border-right: 1px solid #eaeaea;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.user-info {
  padding: 20px;
  border-bottom: 1px solid #eaeaea;
}

.user-badge {
  background-color: #e9d5ff;
  padding: 10px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #6b21a8;
  font-weight: 500;
}

.sidebar-menu {
  flex: 1;
  border-right: none;
}

.menu-item-active {
  background-color: #f0f0f0;
}

/* 右侧内容区域 */
.main-panel {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
  background-color: #ffffff;
}

/* 横幅区域 */
.banner {
  position: relative;
  margin-bottom: 24px;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.banner-image {
  width: 100%;
  height: 180px;
  object-fit: cover;
}

.banner-text {
  position: absolute;
  top: 30px;
  left: 30px;
  color: #ffffff;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

.banner-text h3 {
  font-size: 24px;
  font-weight: bold;
  margin: 0 0 8px 0;
}

.banner-text p {
  font-size: 14px;
  margin: 0;
  opacity: 0.9;
}

/* 功能卡片区域 */
.feature-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
}

.feature-card {
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.feature-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  transform: translateY(-2px);
}

.card-content {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 12px;
  padding: 24px;
}

.card-icon {
  font-size: 32px;
  color: #6b21a8;
}

.card-content h4 {
  font-size: 16px;
  font-weight: bold;
  margin: 0;
  color: #333333;
}

.card-content p {
  font-size: 14px;
  color: #666666;
  margin: 0;
  line-height: 1.4;
}

/* 面板头部 */
.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 1px solid #eaeaea;
}

.panel-header h3 {
  font-size: 18px;
  font-weight: bold;
  color: #333333;
  margin: 0;
}

/* 机器列表 */
.machine-list {
  margin-top: 20px;
}

.machine-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.own-machine-tag {
  font-size: 12px;
  padding: 2px 6px;
  height: 20px;
  line-height: 16px;
}

/* 加载状态 */
.loading {
  min-height: 300px;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 空状态 */
.empty-state {
  min-height: 300px;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 控制中心 */
.control-center {
  height: 100%;
  display: flex;
  flex-direction: column;
}

/* 远程推送 */
.remote-push {
  height: 100%;
  display: flex;
  flex-direction: column;
}

/* 命令管理 */
.command-management {
  height: 100%;
  display: flex;
  flex-direction: column;
}

/* 首页 */
.home {
  height: 100%;
  display: flex;
  flex-direction: column;
}

/* 远程桌面 */
.remote-desktop-container {
  height: 70vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  position: relative;
}

/* 控制栏显示/隐藏按钮 */
.control-toggle {
  position: absolute;
  top: 0;
  left: 50%;
  transform: translateX(-50%);
  z-index: 10;
  background-color: rgba(255, 255, 255, 0.8);
  border-radius: 0 0 4px 4px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.toggle-btn {
  width: 40px;
  height: 30px;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 远程桌面控制栏 */
.remote-control-bar {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px;
  background-color: rgba(245, 245, 245, 0.95);
  border-bottom: 1px solid #eaeaea;
  z-index: 5;
  transition: all 0.3s ease;
}

.machine-info {
  display: flex;
  gap: 20px;
  font-size: 14px;
  color: #333;
}

.remote-controls {
  display: flex;
  gap: 10px;
}

/* 远程桌面显示区域 */
.remote-desktop-display {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: auto;
  transition: transform 0.3s ease;
}

.remote-desktop-screen {
  position: relative;
  max-width: 100%;
  max-height: 100%;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  border-radius: 4px;
  overflow: hidden;
}

.desktop-image {
  width: 100%;
  height: auto;
  display: block;
}

.desktop-video {
  width: 100%;
  height: auto;
  display: block;
  background-color: #000;
}

.connection-status {
  position: absolute;
  bottom: 10px;
  right: 10px;
  display: flex;
  align-items: center;
  gap: 10px;
  background-color: rgba(255, 255, 255, 0.8);
  padding: 5px 10px;
  border-radius: 4px;
  font-size: 12px;
  color: #333;
}

/* 应用中心 */
.app-center {
  height: 100%;
  display: flex;
  flex-direction: column;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .sidebar {
    width: 180px;
  }
  
  .feature-cards {
    grid-template-columns: 1fr;
  }
  
  .banner-text {
    top: 20px;
    left: 20px;
  }
  
  .banner-text h3 {
    font-size: 20px;
  }
}
</style>