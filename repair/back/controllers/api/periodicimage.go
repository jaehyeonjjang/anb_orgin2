package api

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"image/jpeg"
	"log"
	"os"
	"repair/controllers"
	"repair/global"
	"repair/global/config"
	"time"

	"repair/models"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

type PeriodicimageController struct {
	controllers.Controller
}

func (c *PeriodicimageController) Pre_Delete(item *models.Periodicimage) {
	conn := c.NewConnection()

	periodicimageManager := models.NewPeriodicimageManager(conn)
	periodicimageItem := periodicimageManager.Get(item.Id)

	removePeriodicimageFile(periodicimageItem)
}

func (c *PeriodicimageController) Pre_Deletebatch(items *[]models.Periodicimage) {
	if len(*items) == 0 {
		return
	}

	for _, v := range *items {
		removePeriodicimageFile(&v)
	}
}

func removePeriodicimageFile(item *models.Periodicimage) {
	fullFilename := fmt.Sprintf("%v/periodic/%v", config.UploadPath, item.Filename)
	os.Remove(fullFilename)
}

// @Post()
func (c *PeriodicimageController) Process(item *models.Periodicimage) {
	conn := c.NewConnection()

	periodicimageManager := models.NewPeriodicimageManager(conn)
	items := periodicimageManager.Find([]interface{}{
		models.Where{Column: "periodic", Value: item.Periodic, Compare: "="},
	})

	max := 0
	for _, v := range items {
		if v.Order > max {
			max = v.Order
		}
	}

	fullFilename := fmt.Sprintf("%v/%v", config.UploadPath, item.Filename)
	
	// PDF 페이지 수 확인
	pageCount, err := api.PageCountFile(fullFilename)
	if err != nil {
		log.Println(err)
		return
	}

	// 임시 디렉토리 생성
	tmpDir := fmt.Sprintf("%v/tmp-%v", config.UploadPath, global.UniqueId())
	os.Mkdir(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	now := global.GetDatetime(time.Now())
	for i := 1; i <= pageCount; i++ {
		// PDF 페이지를 이미지로 추출
		err := api.ExtractImagesFile(fullFilename, tmpDir, []string{fmt.Sprintf("%d", i)}, nil)
		if err != nil {
			log.Println(err)
			continue
		}

		// 추출된 이미지 파일 찾기
		files, _ := os.ReadDir(tmpDir)
		if len(files) == 0 {
			continue
		}

		// 첫 번째 이미지 읽기
		tmpImgPath := fmt.Sprintf("%v/%v", tmpDir, files[0].Name())
		tmpFile, err := os.Open(tmpImgPath)
		if err != nil {
			log.Println(err)
			continue
		}

		img, _, err := image.Decode(tmpFile)
		tmpFile.Close()
		os.Remove(tmpImgPath)
		if err != nil {
			log.Println(err)
			continue
		}

		id := global.UniqueId()
		filename := fmt.Sprintf("periodic/pastimage-%v.jpg", id)
		fullFilename := fmt.Sprintf("%v/%v", config.UploadPath, filename)
		f, err := os.Create(fullFilename)
		if err != nil {
			log.Println(err)
			continue
		}

		jpeg.Encode(f, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
		f.Close()

		item := models.Periodicimage{Type: item.Type, Filename: filename, Name: item.Name, Use: item.Use, Order: max + i + 1, Periodic: item.Periodic, Date: now}
		periodicimageManager.Insert(&item)

	}

}
