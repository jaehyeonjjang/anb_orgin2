<template>
  <div style="display:flex;justify-content:space-between;margin-bottom:10px;">
    <div></div>
    <el-button size="small" type="success" @click="clickUpdate">수정</el-button>
  </div>

  <el-table :data="tableData" border :span-method="spanMethod" class="attachment-table" :row-style="rowStyle">
    <el-table-column label="구분" align="center">
      <el-table-column width="100" align="center">
        <template #header>
          <span></span>
        </template>
        <template #default="scope">
          {{ shouldShowCategory(scope.$index) ? scope.row.category : '' }}
        </template>
      </el-table-column>
      <el-table-column width="100" align="center">
        <template #header>
          <span></span>
        </template>
        <template #default="scope">
          {{ scope.row.subcategory }}
        </template>
      </el-table-column>
    </el-table-column>
    <el-table-column label="조사항목" width="180" align="center">
      <template #default="scope">
        <span style="white-space:pre-line;">{{ scope.row.item }}</span>
      </template>
    </el-table-column>
    <el-table-column label="내용">
      <template #default="scope">
        <template v-if="scope.$index < 7">
          <span style="white-space:pre-line;">{{ scope.row.content }}</span>
        </template>
        <template v-else>
          <div style="padding:8px;white-space:pre-wrap;word-break:break-word;">
            {{ data.item.content1 || '' }}
          </div>
        </template>
      </template>
    </el-table-column>
    <el-table-column label="이미지" width="280" align="center">
      <template #default="scope">
        <div style="min-height:40px;padding:8px;display:flex;flex-wrap:wrap;gap:8px;justify-content:center;align-items:center;">
          <template v-if="scope.$index < 7">
            <img 
              v-for="(img, idx) in getRowImages(scope.$index)" 
              :key="idx"
              :src="util.getImagePath(img, '')" 
              alt="이미지"
              style="width:30px;height:30px;object-fit:cover;cursor:pointer;border:1px solid #ddd;"
              @click="showImagePreview(util.getImagePath(img, ''))"
            />
          </template>
        </div>
      </template>
    </el-table-column>
  </el-table>

  <!-- 수정 다이얼로그 -->
  <el-dialog v-model="data.visible" width="1200px">
    <y-table>
      <y-tr>
        <y-th style="text-align:center;width:100px;">구분</y-th>
        <y-th style="text-align:center;width:100px;"></y-th>
        <y-th style="text-align:center;width:180px;">조사항목</y-th>
        <y-th>내용</y-th>
        <y-th style="text-align:center;width:250px;">이미지</y-th>
      </y-tr>
      <y-tr v-for="(row, index) in data.editData" :key="index">
        <y-td v-if="index === 0" rowspan="6" style="text-align:center;">{{ row.category }}</y-td>
        <y-td v-else-if="index === 6" colspan="2" style="text-align:center;">{{ row.category }}</y-td>
        <y-td v-else-if="index === 7" colspan="2" style="text-align:center;">{{ row.category }}</y-td>
        
        <y-td v-if="index < 6 && index === 0" rowspan="3" style="text-align:center;">{{ row.subcategory }}</y-td>
        <y-td v-else-if="index === 3" rowspan="2" style="text-align:center;">{{ row.subcategory }}</y-td>
        <y-td v-else-if="index === 5" style="text-align:center;">{{ row.subcategory }}</y-td>
        
        <y-td v-if="index < 7" style="white-space:pre-line;">{{ row.item }}</y-td>
        
        <y-td v-if="index < 7" style="white-space:pre-line;">{{ row.content }}</y-td>
        <y-td v-else colspan="2">
          <el-input 
            v-model="data.item.content1" 
            type="textarea" 
            :rows="3"
            placeholder="점검내용을 입력하세요"
            style="width:100%;"
          />
        </y-td>
        
        <y-td v-if="index < 7" style="text-align:center;">
          <div style="min-height:40px;display:flex;align-items:center;justify-content:flex-start;gap:8px;padding:4px 8px;">
            <el-upload
              accept="image/jpeg,image/png"
              ref="upload"
              :data="{path:'periodic'}"
              :action="data.upload"
              :headers="headers"
              :limit="10"
              :on-success="(res: any, file: any, files: any) => handleSuccess(res, file, files, index)"
              :show-file-list="false"
              :auto-upload="true"
              :before-upload="beforeImageUpload"
            >
              <el-icon style="font-size:24px;cursor:pointer;"><Plus /></el-icon>
            </el-upload>

            <template v-if="data.editData[index].images && data.editData[index].images.length > 0">
              <el-image 
                v-for="(img, pos) in data.editData[index].images" 
                :key="pos"
                style="width: 30px; height: 30px; cursor:pointer;"
                :src="util.getImagePath(img, '')"                    
                fit="cover"
                @click="clickImageDelete(index, Number(pos))"
              />
            </template>
          </div>
        </y-td>
      </y-tr>
    </y-table>

    <template #footer>
      <el-button size="small" @click="data.visible = false">취소</el-button>
      <el-button size="small" type="primary" @click="clickSubmit">저장</el-button>
    </template>
  </el-dialog>

  <!-- 이미지 미리보기 다이얼로그 -->
  <el-dialog v-model="previewDialog.visible" width="80%" :close-on-click-modal="true">
    <img :src="previewDialog.imageUrl" style="width:100%;height:auto;" alt="미리보기" />
  </el-dialog>
</template>

<script setup lang="ts">

import { reactive, onMounted, computed, ref } from "vue"
import { util, size } from "~/global"
import { Periodicotheretc } from "~/models"
import { useStore } from 'vuex'
import { useRoute } from 'vue-router'
import { ElDialog, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import type { UploadInstance, UploadProps, UploadFile, UploadFiles } from 'element-plus'

const { width, height } = size()

const store = useStore()
const route = useRoute()

const imageBaseUrl = import.meta.env.VITE_REPORT_URL

const headers = {
  Authorization: 'Bearer ' + store.state.token
}

const upload = ref<UploadInstance>()

// 앱의 imageCell 키(0,2,3,4,5,6,7)를 웹 테이블 행(0,1,2,3,4,5,6)으로 매핑
const rowKeyMap = [0, 2, 3, 4, 5, 6, 7]

const previewDialog = reactive({
  visible: false,
  imageUrl: ''
})

const data = reactive({
  id: 0,
  apt: 0,
  visible: false,
  item: {
    id: 0,
    content1: '',
    content2: '',
    periodic: 0,
    date: ''
  },
  editData: [] as any[],
  upload: `${import.meta.env.VITE_REPORT_URL}/api/upload/index`
})

function showImagePreview(url: string) {
  previewDialog.imageUrl = url
  previewDialog.visible = true
}

function getRowImages(rowIndex: number): string[] {
  if (!data.item.content2) return []
  
  try {
    const imagesData = JSON.parse(data.item.content2)
    // rowIndex(0-6)를 앱의 키로 변환
    const key = rowKeyMap[rowIndex]
    const images = imagesData[key.toString()]
    return images || []
  } catch (e) {
    return []
  }
}

function spanMethod({ rowIndex, columnIndex }: { rowIndex: number, columnIndex: number }) {
  // 대분류 컬럼(0): 부착물 등 6행 병합
  if (columnIndex === 0) {
    if (rowIndex === 0) return { rowspan: 6, colspan: 1 }
    if (rowIndex >= 1 && rowIndex <= 5) return { rowspan: 0, colspan: 0 }
    // 변위 변형, 점검내용일 때 2칸 병합
    if (rowIndex === 6 || rowIndex === 7) return { rowspan: 1, colspan: 2 }
  }
  // 소분류 컬럼(1): 정착부 3행, 연결부 2행, 보강부 1행 병합
  if (columnIndex === 1) {
    if (rowIndex === 0) return { rowspan: 3, colspan: 1 } // 정착부
    if (rowIndex === 1 || rowIndex === 2) return { rowspan: 0, colspan: 0 }
    if (rowIndex === 3) return { rowspan: 2, colspan: 1 } // 연결부
    if (rowIndex === 4) return { rowspan: 0, colspan: 0 }
    if (rowIndex === 5) return { rowspan: 1, colspan: 1 } // 보강부
    if (rowIndex === 6 || rowIndex === 7) return { rowspan: 0, colspan: 0 }
  }
  // 조사항목 컬럼(2): 점검내용은 숨김
  if (columnIndex === 2) {
    if (rowIndex === 7) return { rowspan: 0, colspan: 0 }
  }
  // 내용 컬럼(3): 점검내용은 2칸 병합
  if (columnIndex === 3) {
    if (rowIndex === 7) return { rowspan: 1, colspan: 2 }
  }
  // 이미지 컬럼(4): 점검내용은 숨김
  if (columnIndex === 4) {
    if (rowIndex === 7) return { rowspan: 0, colspan: 0 }
  }
  return { rowspan: 1, colspan: 1 }
}

function rowStyle({ rowIndex }: { rowIndex: number }) {
  // 점검내용 행의 높이 증가
  if (rowIndex === 7) {
    return { height: '70px' }
  }
  return {}
}

function shouldShowCategory(index: number) {
  if (index === 0) return true
  return tableData[index].category !== tableData[index - 1].category
}

const tableData = [
  { category: '부착물 등', subcategory: '정착부', item: '앵커 정착부\n브라켓 정착부', content: '정착부 콘크리트의 균열, 박락 여부 확인\n앵커 시공상태, 풀림 및 빠짐, 부식 여부 확인' },
  { category: '부착물 등', subcategory: '정착부', item: '용접 정착부', content: '용접 면적 적정성, 균열, 부식 등 손상 발생 여부 등 확인' },
  { category: '부착물 등', subcategory: '정착부', item: '매립 정착부', content: '정착 철물 매립 길이, 철물 여장(노출) 길이 등' },
  { category: '부착물 등', subcategory: '연결부', item: '볼트 연결부', content: '볼트 풀림 및 빠짐, 부재 변형, 부식 등 확인\n볼트의 시공상태(볼트 규격, 설치 간격 등) 확인' },
  { category: '부착물 등', subcategory: '연결부', item: '용접 연결부', content: '용접 면적 적정성, 균열 발생 여부 등 확인' },
  { category: '부착물 등', subcategory: '보강부', item: '와이어 로프', content: '손상, 부식, 변형, 고정클립 수량 및 상태 등 확인' },
  { category: '변위 변형', subcategory: '', item: '기울기 및 배부름', content: '면외방향 기울기 및 배부름 발생 유무' },
  { category: '점검내용', subcategory: '', item: '', content: '' },
]

async function getItems() {
  const res: any = await Periodicotheretc.getByPeriodic(data.id)
  if (res && res.item) {
    data.item = res.item
  } else {
    data.item.periodic = data.id
  }
}

async function clickUpdate() {
  // editData 초기화
  data.editData = tableData.map((row, index) => {
    const rowData = { ...row, images: [] as string[] }
    
    // content2에서 해당 행의 이미지 가져오기
    if (data.item.content2 && index < 7) {
      try {
        const imagesData = JSON.parse(data.item.content2)
        const key = rowKeyMap[index]
        rowData.images = imagesData[key.toString()] || []
      } catch (e) {
        rowData.images = []
      }
    }
    
    return rowData
  })
  
  data.visible = true
}

async function clickSubmit() {
  util.loading(true)
  
  // editData에서 이미지 데이터를 content2로 변환
  const imagesData: Record<string, string[]> = {}
  data.editData.forEach((row, index) => {
    if (index < 7 && row.images && row.images.length > 0) {
      const key = rowKeyMap[index]
      imagesData[key.toString()] = row.images
    }
  })
  
  const item = {
    id: data.item.id,
    content1: data.item.content1,
    content2: Object.keys(imagesData).length > 0 ? JSON.stringify(imagesData) : '',
    periodic: data.item.periodic,
    date: data.item.date
  }
  
  if (item.id > 0) {
    await Periodicotheretc.update(item)
  } else {
    await Periodicotheretc.insert(item)
  }
  
  data.visible = false
  await getItems()
  
  util.info('저장되었습니다')
  util.loading(false)
}

function handleSuccess(response: any, file: any, files: any, index: number) {
  const filename = response.filename
  
  if (!data.editData[index].images) {
    data.editData[index].images = []
  }
  
  data.editData[index].images.push(filename)

  // 업로드 파일 목록 초기화
  if (upload.value && upload.value.length) {
    for (let i = 0; i < upload.value.length; i++) {
      upload.value[i].clearFiles()
    }
  }
}

function beforeImageUpload(file: File) {
  const isImage = (file.type === 'image/jpeg' || file.type === 'image/png')

  if (!isImage) {
    util.error('이미지 파일만 업로드 가능합니다 (jpg, png)')
  }

  return isImage
}

async function clickImageDelete(rowIndex: number, imageIndex: number) {
  try {
    await ElMessageBox.confirm(
      '삭제하시겠습니까',
      '확인',
      {
        confirmButtonText: '삭제',
        cancelButtonText: '취소',
        type: 'warning',
      }
    )
    
    if (data.editData[rowIndex].images) {
      data.editData[rowIndex].images.splice(imageIndex, 1)
    }
  } catch (error) {
    // 취소 버튼 클릭 시 아무 작업도 하지 않음
  }
}

onMounted(async () => {
  util.loading(true)

  const apt = parseInt(route.params.apt as string)
  const id = parseInt(route.params.id as string)

  data.id = id
  data.apt = apt

  await getItems()

  util.loading(false)
})
</script>

<style scoped>
.attachment-table {
  width: 100%;
}

/* 구분 헤더 스타일 조정 */
.attachment-table :deep(.el-table__header-wrapper thead tr:first-child th:first-child) {
  border-bottom: 1px solid #dadada !important;
}

.attachment-table :deep(.el-table__header-wrapper thead tr:last-child th:nth-child(1)),
.attachment-table :deep(.el-table__header-wrapper thead tr:last-child th:nth-child(2)) {
  display: none;
}

/* 점검내용 행의 셀 높이와 중앙 정렬 */
.attachment-table :deep(.el-table__body tr:nth-child(8) td) {
  height: 120px;
  vertical-align: middle;
}
</style>

