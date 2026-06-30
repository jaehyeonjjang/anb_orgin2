<template>
  <Title title="외벽 마감재 및 부착물 등" />

  <el-tabs v-model="data.menu">

    <el-tab-pane label="외벽 마감재" name="outerwall">
      <OuterwallInsert />
      <OuterwallBottomInsert />
    </el-tab-pane>

    <el-tab-pane label="부착물 등" name="attachment">
      <AttachmentInsert />
    </el-tab-pane>

  </el-tabs>

</template>

<script setup lang="ts">

import { reactive, onMounted, ref, watch } from "vue"
import router from '~/router'
import { util, size }  from "~/global"
import { useStore } from 'vuex'
import { useRoute } from 'vue-router'
import request from '~/global/request'

const { width, height } = size()

const store = useStore()
const route = useRoute()

const data = reactive({
  id: 0,
  apt: 0,
  menu: 'outerwall'  
})

async function initData() {
}

async function getItems() {
}

onMounted(async () => {
  const apt = parseInt(route.params.apt)
  const id = parseInt(route.params.id)
  
  data.id = id
  data.apt = apt

  if (store.getters['getUser'] != null) {
    data.level = store.getters['getUser'].level
  }
  
  await initData()
  await getItems()
})

</script>
