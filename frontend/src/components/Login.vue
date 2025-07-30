<script setup lang="ts">
import { fetchLogin } from '@/api';
import { router } from '@/router'


const username = ref('')
const password = ref('')

async function toLogin() {
  await fetchLogin(username.value, password.value)
  window.$message?.success("登录成功")
  router.push('/')
}
</script>

<template>
<div class="login">
  <n-card title="Sync" hoverable>
    <n-input class="input" type="text" v-model:value="username" placeholder="账户"></n-input>
    <n-input class="input" type="password" v-model:value="password" placeholder="密码"></n-input>
    <n-button class="input" type="primary" strong block @click="toLogin">登录</n-button>
  </n-card>
</div>
</template>

<style scoped>
.login {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 500px;
  background-color: antiquewhite;
  transform: translate(-50%,-80%);
}
.input {
  margin-top: 10px;
}
/* .ncard { */
  /* text-align: center; */
  /* max-width: 300px; */
/* } */
</style>