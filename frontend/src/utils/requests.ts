import { router } from '@/router'
import axios,{type AxiosProgressEvent,type AxiosResponse,type GenericAbortSignal } from 'axios'

export interface HttpOption {
  url: string
  data?: any
  method?: string
  headers?: any
  onDownloadProgress?: (progressEvent: AxiosProgressEvent) => void
  signal?: GenericAbortSignal
  beforeRequest?: () => void
  afterRequest?: () => void
}

export interface Response<T = any> {
  code: number
  data: T
  msg: string | null
}

const request = axios.create({
  baseURL: import.meta.env.VITE_GLOB_API_PREFIX,
})

request.interceptors.request.use(
  (config) => {
    return config
  },
  (error) => {
    return Promise.reject(error.response)
  },
)

request.interceptors.response.use(
  (response: AxiosResponse) => {
    if (response.status != 200) {
      window.$message?.error("网络连接异常")
      return Promise.reject(response.data)
    }
    if (response.data.code === 401) {
      window.$message?.error(response.data.msg)
      router.push('/login')
      return Promise.reject(response.data)
    }
    if (response.data.code != 200) {
      window.$message?.error(response.data.msg)
      return Promise.reject(response.data)
    }
    return response.data
  },
  (error) => {
    return Promise.reject(error)
  },
)

function http<T = any>(
  { url, data, method, headers, onDownloadProgress, signal, beforeRequest, afterRequest }: HttpOption,
):Promise<T> {
  const successHandler = (res: AxiosResponse<Response<T>>) => {
    // const authStore = useAuthStore()

    if (res.data.code && res.data.code !== 200)
      return Promise.reject(res.data)
    return res.data
  }

  const failHandler = (error: Response<Error>) => {
    afterRequest?.()
    throw new Error(error?.msg || 'Error')
  }

  beforeRequest?.()

  method = method || 'GET'

  const params = Object.assign(typeof data === 'function' ? data() : data ?? {}, {})

  return method === 'GET'
    ? request.get(url, { params, signal, onDownloadProgress }).then(successHandler, failHandler).then<T>()
    : request.post(url, params, { headers, signal, onDownloadProgress }).then(successHandler, failHandler).then<T>()
}

export function get<T = any>(
  { url, data, method = 'GET', onDownloadProgress, signal, beforeRequest, afterRequest }: HttpOption,
): Promise<T> {
  return http<T>({
    url,
    method,
    data,
    onDownloadProgress,
    signal,
    beforeRequest,
    afterRequest,
  })
}

export function post<T = any>(
  { url, data, method = 'POST', headers, onDownloadProgress, signal, beforeRequest, afterRequest }: HttpOption,
): Promise<T> {
  return http<T>({
    url,
    method,
    data,
    headers,
    onDownloadProgress,
    signal,
    beforeRequest,
    afterRequest,
  })
}

export default post
