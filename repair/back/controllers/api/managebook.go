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
	"strings"
	"time"

	"repair/models"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

type ManagebookController struct {
	controllers.Controller
}

func (c *ManagebookController) Pre_Delete(item *models.Managebook) {
	conn := c.NewConnection()

	managebookManager := models.NewManagebookManager(conn)
	managebookItem := managebookManager.Get(item.Id)

	removeFile(managebookItem)
}

func (c *ManagebookController) Pre_Deletebatch(items *[]models.Managebook) {
	if len(*items) == 0 {
		return
	}

	for _, v := range *items {
		removeFile(&v)
	}
}

func removeFile(item *models.Managebook) {
	fullFilename := fmt.Sprintf("%v/periodicresult/%v/%v", config.UploadPath, item.Periodic, item.Filename)
	os.Remove(fullFilename)
}

// @Post()
func (c *ManagebookController) Process(id int64, name string, order int, filename string) {
	conn := c.NewConnection()

	managebookManager := models.NewManagebookManager(conn)
	managebookcategoryManager := models.NewManagebookcategoryManager(conn)

	categorys := managebookcategoryManager.Find([]any{
		models.Where{Column: "periodic", Value: id, Compare: "="},
		models.Where{Column: "name", Value: name, Compare: "="},
	})

	now := global.GetDatetime(time.Now())
	var categoryId int64
	if len(categorys) > 0 {
		category := categorys[0]
		managebookManager.DeleteByManagebookcategory(category.Id)

		categoryId = category.Id
	} else {
		category := models.Managebookcategory{Name: name, Order: order, Periodic: id, Date: now}
		managebookcategoryManager.Insert(&category)
		categoryId = managebookcategoryManager.GetIdentity()
	}

	fullFilename := fmt.Sprintf("%v/%v", config.UploadPath, filename)
	
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

		os.Mkdir(fmt.Sprintf("%v/periodicresult/%v", config.UploadPath, id), 0755)

		filename := fmt.Sprintf("managebook-%v.jpg", global.UniqueId())
		fullFilename := fmt.Sprintf("%v/periodicresult/%v/%v", config.UploadPath, id, filename)
		f, err := os.Create(fullFilename)
		if err != nil {
			log.Println(err)
			continue
		}

		jpeg.Encode(f, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
		f.Close()

		item := models.Managebook{Filename: filename, Order: i + 1, Periodic: id, Managebookcategory: categoryId, Date: now}
		managebookManager.Insert(&item)

	}

}

// @Post()
func (c *ManagebookController) Multiprocess(id int64, filename string, originalfilename string) {
	conn := c.NewConnection()

	managebookManager := models.NewManagebookManager(conn)
	managebookcategoryManager := models.NewManagebookcategoryManager(conn)

	filenames := strings.Split(filename, ",")
	originalfilenames := strings.Split(originalfilename, ",")

	for j, name := range originalfilenames {
		name = strings.ReplaceAll(name, ".pdf", "")
		filename := filenames[j]

		categorys := managebookcategoryManager.Find([]any{
			models.Where{Column: "periodic", Value: id, Compare: "="},
			models.Where{Column: "name", Value: name, Compare: "="},
		})

		now := global.GetDatetime(time.Now())
		var categoryId int64
		if len(categorys) > 0 {
			category := categorys[0]
			managebookManager.DeleteByManagebookcategory(category.Id)

			categoryId = category.Id
		} else {
			category := models.Managebookcategory{Name: name, Order: 0, Periodic: id, Date: now}
			managebookcategoryManager.Insert(&category)
			categoryId = managebookcategoryManager.GetIdentity()
		}

		fullFilename := fmt.Sprintf("%v/%v", config.UploadPath, filename)
		
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

			os.Mkdir(fmt.Sprintf("%v/periodicresult/%v", config.UploadPath, id), 0755)

			filename := fmt.Sprintf("managebook-%v.jpg", global.UniqueId())
			fullFilename := fmt.Sprintf("%v/periodicresult/%v/%v", config.UploadPath, id, filename)
			f, err := os.Create(fullFilename)
			if err != nil {
				log.Println(err)
				continue
			}

			jpeg.Encode(f, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
			f.Close()

			item := models.Managebook{Filename: filename, Order: i + 1, Periodic: id, Managebookcategory: categoryId, Date: now}
			managebookManager.Insert(&item)

		}
	}
}
