<script setup lang="ts">
import { ref } from 'vue'
import { fetchRepos, fetchCommits,fetchRevert } from '../api/index'
import type {CommitsRes, Commit} from '../api/index'
import { NButton, type PaginationProps  } from 'naive-ui'

const repos = ref()
const id = ref(1)
const loading = ref(false)
const commits = ref<CommitsRes>({
  total: 0,
  list: []
})
const page = reactive<PaginationProps>({
  page:1,
  itemCount: 50,
  pageSize: 10,
  prefix({ itemCount }) {
          return `共 ${itemCount} 条`
        }
})

const columns = [
  {
    title: 'Hash',
    key: 'hash'
  },
  {
    title: 'Message',
    key: 'message'
  },
  {
    title: 'Date',
    key: 'date'
  },
  {
    title: 'Action',
    key:'actions',
    render(row:Commit) {
      return h(
        NButton,
        {
          strong:true,
          onClick: () => toRevert(row)
        },{
          default:() => '恢复'
        }
      )
    }
  }
]

async function getRepos() {
  repos.value = await fetchRepos()
  console.log(repos.value)
}

async function getCommits() {
  loading.value = true
  commits.value = await fetchCommits({
    id:id.value,
    pager: {
      index: page.page,
      size: page.pageSize
    }
  })
  page.itemCount = commits.value.total
  console.log(page)
  loading.value = false
}

async function toRevert(row:Commit) {
  await fetchRevert({
    id: id.value,
    hash:row.hash,
    file:[]
  })
  window.$message?.success("恢复成功")
}

async function update(value: number) {
  id.value = value + 1
  getCommits()
}

async function updatePage(p:number) {
  page.page = p
  getCommits()
}

getRepos()
getCommits()
</script>

<template>
  <div>
    <n-card title="仓库">
      <template #header-extra>
      </template>
      <n-tabs type="line" animated :default-value="0" @update:value="update">
        <n-tab-pane v-for="item, i in repos" :key="i" :name="i" :tab="item.name">
          <n-data-table remote :loading="loading" :data="commits.list" :columns="columns" :row-key="(r) => r.hash" :pagination="page" @update:page="updatePage" />
          </n-tab-pane>
          <template #prefix>
            <n-button @click="getCommits">刷新</n-button>
      </template>
    </n-tabs>
  </n-card>
</div>
</template>

<style scoped></style>
