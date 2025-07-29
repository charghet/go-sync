import { post } from '@/utils/requests'

export function fetchLogin<T>(username: string, password: string) {
  return post<T>({
    url: "/login",
    data: {
      username, password
    }
  })
}

export function fetchRepos<T>() {
  return post<T>({
    url: "/repos",
  })
}

export interface CommitsReq {
  id: number,
  pager: {
    index?: number,
    size?: number
  }
}
export interface CommitsRes {
  total: number,
  list: Commit[]
}
export interface Commit {
    hash: string,
    message: string,
    author: string,
    date: string,
    email: string
}

export function fetchCommits(data: CommitsReq): Promise<CommitsRes> {
  return post<CommitsRes>({
    url: "/commits",
    data
  })
}

export interface RevertReq {
  id: number,
  hash: string,
  file: string[]
}

export function fetchRevert(data: RevertReq): Promise<any> {
  return post<CommitsRes>({
    url: "/revert",
    data
  })
}