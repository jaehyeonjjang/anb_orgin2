<template>
  <Title title="결함현황표" />

  <el-tabs v-model="data.menu">

    <el-tab-pane label="결함도" name="defect">
      <DataDefectInsert />
    </el-tab-pane>

    <el-tab-pane label="기울기" name="slope">
      <DataSlopeInsert />
    </el-tab-pane>

    <el-tab-pane label="강도/탄산화" name="strength">
      <DataStrengthInsert />
    </el-tab-pane>

    <el-tab-pane label="부재" name="member">
      <DataMemberInsert />
    </el-tab-pane>

  </el-tabs>

</template>


<script setup lang="ts">

import { reactive, onMounted } from "vue"
import router from '~/router'
import { util, size }  from "~/global"
import { useStore } from 'vuex'
import { useRoute } from 'vue-router'
import DataDefectInsert from '~/components/periodic/DataDefectInsert.vue'
import DataSlopeInsert from '~/components/periodic/DataSlopeInsert.vue'
import DataStrengthInsert from '~/components/periodic/DataStrengthInsert.vue'
import DataMemberInsert from '~/components/periodic/DataMemberInsert.vue'

const { width, height } = size()

const store = useStore()
const route = useRoute()

const data = reactive({
  id: 0,
  apt: 0,
  menu: 'defect'  
})

onMounted(async () => {
  const apt = parseInt(route.params.apt)
  const id = parseInt(route.params.id)
  
  data.apt = apt
  data.id = id
})

</script>
