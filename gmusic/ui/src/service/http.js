import axios from 'axios'

const http = axios.create({
  baseURL: 'http://localhost:8080/api',
  timeout: 10000
})

// 可选：请求拦截器
http.interceptors.request.use(cfg => cfg, err => Promise.reject(err))

// 可选：响应拦截器
http.interceptors.response.use(
  res => res,
  err => Promise.reject(err)
)

export default http

