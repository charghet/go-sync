<script setup lang="ts">
import { ref } from 'vue'
import { fetchRepos, fetchCommits, fetchRevert, fetchChanges,type ChangesRes } from '../api/index'
import type { CommitsRes, Commit } from '../api/index'
import { NButton, NSpace, type DataTableColumns, type PaginationProps } from 'naive-ui'

const repos = ref()
const id = ref(1)
const loading = ref(false)
const drawer = reactive({
  show:false,
  title:'',
})
const commits = ref<CommitsRes>({
  total: 0,
  list: []
})
const changes = ref<ChangesRes[]>([])
const page = reactive<PaginationProps>({
  page: 1,
  itemCount: 50,
  pageSize: 10,
  prefix({ itemCount }) {
    return `共 ${itemCount} 条`
  }
})

const columns: DataTableColumns<Commit> = [
  {
    title: 'Hash',
    key: 'hash',
    render(row) {
      return h(
        'p',
        shortHash(row.hash)
      )
    }
  },
  {
    title: 'Message',
    key: 'message',
  },
  {
    title: 'Date',
    key: 'date',
  },
  {
    title: 'Action',
    key: 'actions',
    align: 'center',
    render(row) {
      return h(
        NSpace,
        { size: 8, justify: "center" },
        () => [
          h(
            NButton,
            {
              strong: true,
              onClick: () => getChanges(row)
            }, {
            default: () => '查看'
          }
          ),
          h(
            NButton,
            {
              strong: true,
              onClick: () => toRevert(row)
            }, {
            default: () => '恢复'
          }
          ),
        ]
      )

    }
  },
]

function shortHash(hash:string) {
  const len = 7
  if(hash.length <= len) {
    return hash
  }
  return hash.substring(0, len+1)
}

async function getRepos() {
  repos.value = await fetchRepos()
}

async function getCommits() {
  loading.value = true
  commits.value = await fetchCommits({
    id: id.value,
    pager: {
      index: page.page,
      size: page.pageSize
    }
  })
  page.itemCount = commits.value.total
  loading.value = false
}

async function toRevert(row: Commit) {
  await fetchRevert({
    id: id.value,
    hash: row.hash,
    file: []
  })
  window.$message?.success("恢复成功")
}

async function update(value: number) {
  id.value = value + 1
  getCommits()
}

async function updatePage(p: number) {
  page.page = p
  getCommits()
}

async function getChanges(row: Commit) {
  const res = await fetchChanges({
    id: id.value,
    hash: row.hash
  })
  drawer.show = true
  drawer.title = shortHash(row.hash)
  changes.value = res
}

getRepos().then(() => {
  getCommits()
})
</script>

<template>
<div class="home">
  <n-card title="仓库">
    <template #header-extra>
    </template>
    <n-tabs type="line" animated :default-value="0" @update:value="update">
      <n-tab-pane v-for="item, i in repos" :key="i" :name="i" :tab="item.name">
        <n-data-table remote row-class-name="row" :loading="loading" :data="commits.list" :columns="columns" :row-key="(r) => r.hash"
          :pagination="page" @update:page="updatePage" />
      </n-tab-pane>
      <template #prefix>
        <n-button @click="getCommits">刷新</n-button>
      </template>
    </n-tabs>
  </n-card>
  <n-drawer v-model:show="drawer.show" placement="right" width="400px">
    <n-drawer-content :title="drawer.title">
      <n-list>
        <n-list-item v-for="item in changes">
          <n-thing>
            {{ item.action }} {{ item.name }}
          </n-thing>
        </n-list-item> 
      </n-list>
    </n-drawer-content>
  </n-drawer>
</div>
</template>

<style scoped>
.home {
  margin: 20px
}
:deep(.row td) {
  text-align: left !important;
}
</style>
