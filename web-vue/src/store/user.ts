import { defineStore } from 'pinia'
import { ref } from 'vue'
import { authApi } from '@/api'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('casdoor_token') || '')
  const userInfo = ref<any>(null)

  async function login(application: string, username: string, password: string, organization = 'built-in') {
    const res: any = await authApi.login({
      application,
      username,
      password,
      organization,
      signinMethod: 'Password'
    })
    if (res.data?.token) {
      token.value = res.data.token
      localStorage.setItem('casdoor_token', res.data.token)
    }
    if (res.data?.user) {
      userInfo.value = res.data.user
      localStorage.setItem('casdoor_user', JSON.stringify(res.data.user))
    }
    return res
  }

  async function fetchUserInfo() {
    const res: any = await authApi.getAccount()
    userInfo.value = res.data || res
    localStorage.setItem('casdoor_user', JSON.stringify(userInfo.value))
    return userInfo.value
  }

  function logout() {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('casdoor_token')
    localStorage.removeItem('casdoor_user')
  }

  return { token, userInfo, login, fetchUserInfo, logout }
})
