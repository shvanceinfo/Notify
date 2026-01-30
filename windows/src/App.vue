<template>
  <div class="app">
    <Login v-if="!isLoggedIn" @login-success="handleLoginSuccess" />
    <Main v-else />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import Login from './components/Login.vue'
import Main from './components/Main.vue'

const isLoggedIn = ref(false)

const handleLoginSuccess = () => {
  isLoggedIn.value = true
}

onMounted(() => {
  const token = localStorage.getItem('token')
  if (token) {
    isLoggedIn.value = true
  }
})
</script>

<style scoped>
.app {
  width: 100%;
  height: 100vh;
}
</style>