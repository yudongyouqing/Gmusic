import axios from 'axios'

const http = axios.create({
  baseURL: 'http://localhost:8080/api',
  // 0 表示不超时，避免后端初始化音频设备或解码稍慢导致前端报 10s 超时
  timeout: 0,
  withCredentials: false,
})

// 可选：请求拦截器
http.interceptors.request.use(cfg => cfg, err => Promise.reject(err))

// 可选：响应拦截器（将后端报错直接透传给调用处显示）
http.interceptors.response.use(
  res => res,
  err => Promise.reject(err)
)

export default http
