<template>
  <div style="display:flex;justify-content: space-between;gap:5px;margin-bottom:10px;">    
    <div style="flex:1;text-align:right;gap:5;">
      <el-button size="small" type="danger" @click="clickDeleteMulti" style="margin-right:-5px;">삭제</el-button>
      <el-button size="small" type="success" @click="clickInsert" style="margin-right:-5px;">등록</el-button>
      <el-button size="small" type="warning" @click="clickBatch">일괄처리</el-button>
    </div>
  </div>  

  
  <el-table :data="data.items" border :height="height(170)" @row-click="clickUpdate"  ref="listRef" @selection-change="changeList" v-show="data.preview == false">
    <el-table-column type="selection" width="40" align="center" />
    <el-table-column label="위치" align="center">
      <template #default="scope">
        {{getBlueprint(scope.row)}}
      </template>
    </el-table-column>
    <el-table-column prop="group" label="순번" align="center" width="50" />
    <el-table-column prop="part" label="변위방향" align="center" width="100" />
    <el-table-column prop="member" label="변위량(mm)" align="center" width="100" />
    <el-table-column prop="shape" label="측정높이(mm)" align="center" width="120" />    
    <el-table-column prop="width" label="기울기" align="center" width="100" />
    <el-table-column prop="remark" label="비고" align="center" />
    <el-table-column label="이미지" align="center">
      <template #default="scope">
        <div v-if="scope.row.filename != ''">
          <el-image v-for="(item, index) in scope.row.filename.split(',')" style="width: 20px; height: 20px; top:4px;left:0px;position:relative;margin-right:5px;"
                    :src="util.getImagePath(item)"                   
                    fit="cover"
                    @click.stop="clickPreviews(scope.row.filename, index)"
          />
        </div>
      </template>
    </el-table-column>
    <el-table-column label="위치도" align="center" width="50">
      <template #default="scope">
        <el-image v-if="scope.row.extra.resultimage != ''"          
                  style="width: 20px; height: 20px; top:4px;left:0px;position:relative;z-index:999999999999999 !important;"
                  :src="util.getImagePath(scope.row.extra.resultimage)"                  
                  fit="cover"
                  @click.stop="clickPreview(scope.row.extra.resultimage)"
        />
        
      </template>
    </el-table-column>
  </el-table>  

  
  <el-dialog
    v-model="data.visible"
    :before-close="handleClose"
    width="1050px"
  >

    <el-form label-width="100px">      
      <el-table :data="data.batchs" border style="margin-top:15px;">
        <el-table-column label="" align="center" width="35" v-if="data.mode == 'batch'">
          <template #default="scope">
            <el-icon @click="clickRegistDelete(scope.$index)"><Delete /></el-icon>
          </template>
        </el-table-column>
        <el-table-column label="위치" align="center" width="180">
          <template #default="scope">
            <el-tree-select style="width:170px;" v-model="data.batchs[scope.$index].blueprint" :data="data.blueprints" :default-expand-all="true" :render-after-expand="false" placeholder="위치" />
          </template>
        </el-table-column>
        <el-table-column label="순번" align="center" width="50">
          <template #default="scope">
            <el-input v-model.number="data.batchs[scope.$index].group" />
          </template>
        </el-table-column>
        <el-table-column label="변위방향" align="center" width="90">
          <template #default="scope">
            <el-select v-model="data.batchs[scope.$index].part" placeholder="선택" style="width:100%;">
              <el-option label="우측" value="우측" />
              <el-option label="좌측" value="좌측" />
            </el-select>
          </template>
        </el-table-column>
        <el-table-column label="변위량(mm)" align="center" width="100">
          <template #default="scope">
            <el-input v-model="data.batchs[scope.$index].member" @input="calculateSlope(scope.$index)" />
          </template>
        </el-table-column>
        <el-table-column label="측정높이(mm)" align="center" width="120">
          <template #default="scope">
            <el-input v-model="data.batchs[scope.$index].shape" @input="calculateSlope(scope.$index)" />
          </template>
        </el-table-column>
        <el-table-column label="기울기" align="center" width="100">
          <template #default="scope">
            <el-input v-model="data.batchs[scope.$index].width" readonly />
          </template>
        </el-table-column>
        <el-table-column label="비고" align="center">
          <template #default="scope">
            <el-input v-model="data.batchs[scope.$index].remark" />
          </template>
        </el-table-column>
        
      </el-table>


    </el-form>

      <template #footer>
        <el-button size="small" type="danger" v-if="data.mode != 'batch' && (data.batchs.length > 0 && data.batchs[0].id > 0)" style="float:left;" @click="clickDelete">삭제</el-button>
        <el-button size="small" v-if="data.mode == 'batch'" style="float:left;" @click="clickAdd(1)"><el-icon><Plus /></el-icon></el-button>
        <el-button size="small" v-if="data.mode == 'batch'" style="float:left;" @click="clickAdd(10)"><el-icon><Plus /></el-icon> &nbsp;10</el-button>
        <el-button size="small" @click="clickCancel">취소</el-button>
        <el-button size="small" type="primary" @click="clickSubmit">등록</el-button>
      </template>
  </el-dialog>

</template>


<script setup lang="ts">

import { ref, reactive, onMounted, onUnmounted } from "vue"
import router from '~/router'
import { util, size }  from "~/global"
import { Periodicdata, Blueprint } from "~/models"
import { useStore } from 'vuex'
import { useRoute } from 'vue-router'
import { ElTable } from 'element-plus'
import { v3ImgPreviewFn } from 'v3-img-preview'

const { width, height } = size()

const store = useStore()
const route = useRoute()

const model = Periodicdata

const item = {
  id: 0,
  group: 0,
  part: '',
  member: '',
  shape: '',
  width: '',
  length: '',
  count: 0,
  progress: '',
  remark: '',
  order: 0,
  content: '',
  blueprint: 0,
  periodic: 0,
  date: '',
  type: 201  // inclinationLine
}

const data = reactive({
  apt: 0,
  id: 0,
  mode: 'normal',
  items: [],
  total: 0,
  page: 1,
  pagesize: 0,
  item: util.clone(item),
  visible: false,
  allblueprints: [],
  blueprints: [],
  preview: false,
  batchs: []
})

async function initData() {
  let res = await Blueprint.find({apt: data.apt, orderby: 'bp_parentorder,bp_order desc,bp_id'})

  data.allblueprints = res.items
  let blueprints = [{label:'위치', value: 0}]
  let items = res.items
  if (items == null) {
    items = []
  }
  for (let i = 0; i < items.length; i++) {
    let item = items[i]

    if (item.level != 1) {
      continue
    }

    let children = []
    for (let i = 0; i < items.length; i++) {
      let item2 = items[i]

      if (item2.parent != item.id) {
        continue
      }

      children.push({label: `${item.name} ${item2.name}`, value: item2.id})
    }
    
    blueprints.push({label: item.name, value: item.id, children: children})
  }

  data.blueprints = blueprints
}

async function getItems() {
  let res = await model.find({
    page: data.page,
    pagesize: data.pagesize,
    periodic: data.id,
    orderby: 'bp_parentorder,bp_order desc,bp_id,pd_order,pd_id'    
  })

  let items = []

  if (res.items == null) {
    res.items = []
  }

  for (let i = 0; i < res.items.length; i++) {
    let item = res.items[i]

    if (item.group == 0) {
      continue
    }

    // 기울기 데이터만 필터링 (type 200-299)
    if (item.type < 200 || item.type >= 300) {
      continue
    }
    
    item.index = i + 1
    items.push(item)
  }  

  data.total = res.total  
  data.items = items

  console.log(data.items)
}

function clickInsert() {  
  data.item = util.clone(item)

  let items = [data.item]
  
  data.mode = 'normal'
  data.batchs = items
  data.visible = true  
}

function clickUpdate(item, index) {
  if (index == undefined) {
    return
  }

  if (index.no == 0 || index.no > 10) {
    return
  }
  
  let items = [util.clone(item)]

  data.mode = 'normal'
  data.batchs = items
  data.visible = true  
}

function clickDelete() {
  let item = data.batchs[0]
  
  util.confirm('삭제하시겠습니까', async function() {
    let res = await model.remove(item)
    if (res.code === 'ok') {
      util.info('삭제되었습니다')
      data.visible = false
      await getItems()
    }
  })
}

const handleClose = (done: () => void) => {
  if (data.mode == 'batch') {
    util.confirm('팝업창을 닫으시겠습니까', function() {
      done()
    })
  } else {
    done()
  }
}

onMounted(async () => {
  data.id = parseInt(route.params.id)
  data.apt = parseInt(route.params.apt)  
  
  util.loading(true)
  
  await initData()
  await getItems()

  data.visible = false
  util.loading(false)
})

function clickBatch() {
  let items = util.clone(data.items)

  if (items == null) {
    items = []
  }

  if (items.length == 0) {
    for (let i = 0; i < 5; i++) {
      items.push(util.clone(data.item))
    }
  }

  data.mode = 'batch'
  data.batchs = items
  data.visible = true  
}

function clickCancel() {
  if (data.mode == 'batch') {
    util.confirm('팝업창을 닫으시겠습니까', function() {
      data.visible = false
    })
  } else {
    data.visible = false
  }
}

const listRef = ref<InstanceType<typeof ElTable>>()
const listSelection = ref([])
const toggleListSelection = (rows) => {
  if (rows) {
    rows.forEach((row) => {
      listRef.value!.toggleRowSelection(row, undefined)
    })
  } else {
    listRef.value!.clearSelection()
  }
}
const changeList = (val) => {
  listSelection.value = val
}

function clickDeleteMulti() {
  util.confirm('삭제하시겠습니까', async function() {
    util.loading(true)
    
    for (let i = 0; i < listSelection.value.length; i++) {
      let value = listSelection.value[i]

      let item = {
        id: value.id
      }

      await model.remove(item)
    }

    util.info('삭제되었습니다')
    await getItems()

    util.loading(false)
  })
}

async function clickSubmit() {
  util.loading(true)

  if (data.mode == 'batch') {
    for (let i = 0; i < data.items.length; i++) {
      let item = data.items[i]
      let flag = false;
      for (let j = 0; j < data.batchs.length; j++) {
        if (data.items[i].id == data.batchs[j].id) {
          flag = true
          break
        }
      }

      if (flag == false) {      
        await model.remove(item)
      }
    }
  } else {
    // 신규 등록인 경우에만 order를 최대값+1로 설정
    if (data.batchs[0].id == 0) {
      let max = 0
      for (let i = 0; i < data.items.length; i++) {
        let item = data.items[i]

        if (item.order > max) {
          max = item.order
        }
      }

      max++

      data.batchs[0].order = max
    }
    // 기존 항목 수정인 경우 order 값 유지 (변경하지 않음)
  }

  let blueprint = 0
  
  for (let i = 0; i < data.batchs.length; i++) {
    let item = data.batchs[i]

    if (item.blueprint == 0) {
      item.blueprint = blueprint
    } else {
      blueprint = util.getInt(item.blueprint)
    }
    
    item.periodic = data.id
    item.type = 201  // inclinationLine (기울기 지시선)

    item.progress = util.getInt(item.progress)
    item.group = util.getInt(item.group)
    item.count = util.getInt(item.count)

    if (item.blueprint == 0 && item.group == 0) {
      continue
    }
    
    if (data.mode == 'batch' && item.id == 0) {
      item.order = i + 1
    }
    
    if (item.id > 0) {
      await model.update(item)
    } else { 
      await model.insert(item)
    }
  }

  util.info('등록되었습니다')
  
  await getItems()
  data.visible = false  
  util.loading(false)  
}

function clickRegistDelete(index) {
  data.batchs.splice(index, 1)
}

function clickAdd(count) {
  let items = []
  for (let i = 0; i < count; i++) {
    items.push(util.clone(item))
  }

  data.batchs = data.batchs.concat(items)
}

function getBlueprint(item) {
  for (let i = 0; i < data.allblueprints.length; i++) {
    if (item.blueprint == data.allblueprints[i].id) {
      for (let j = 0; j < data.allblueprints.length; j++) {
        if (data.allblueprints[j].id == data.allblueprints[i].parent) {
          return data.allblueprints[j].name + ' ' + data.allblueprints[i].name 
        }
      }

      return data.allblueprints[i].name
    }
  }
}

function clickPreview(url) {
  const img = util.getImagePath(url)
  v3ImgPreviewFn(img)  
}

function clickPreviews(str, index) {
  const imgs = str.split(',').map(item => util.getImagePath(item)) 
  v3ImgPreviewFn({images:imgs, index: index})  
}

function calculateSlope(index) {
  const item = data.batchs[index]
  const displacement = parseFloat(item.member) || 0  // 변위량
  const height = parseFloat(item.shape) || 0  // 측정높이
  
  if (height > 0 && displacement > 0) {
    // 변위량 / 측정높이 계산
    const slope = displacement / height
    // 분수 형태로 표시 (예: 1/5610)
    item.width = `1/${Math.round(height / displacement)}`
  } else {
    item.width = ''
  }
}
</script>
