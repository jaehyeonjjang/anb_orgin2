package api

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"os"
	"os/exec"
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

	if managebookItem != nil {
		removeFile(managebookItem)
	}
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

// convertPDFToImages converts all pages of a PDF to JPEG using pdftoppm
func convertPDFToImages(pdfPath string, outputDir string, outputPrefix string) (int, error) {
	// pdftoppm 경로 찾기 (macOS Homebrew, Linux 기본 경로)
	pdftoppmPaths := []string{
		"/opt/homebrew/bin/pdftoppm", // macOS M1/M2 Homebrew
		"/usr/local/bin/pdftoppm",    // macOS Intel Homebrew
		"/usr/bin/pdftoppm",          // Linux
		"pdftoppm",                   // PATH에서 찾기
	}

	var pdftoppmPath string
	for _, path := range pdftoppmPaths {
		if _, err := os.Stat(path); err == nil {
			pdftoppmPath = path
			break
		}
	}

	// PATH에서 찾기 시도
	if pdftoppmPath == "" {
		if path, err := exec.LookPath("pdftoppm"); err == nil {
			pdftoppmPath = path
		} else {
			return 0, fmt.Errorf("pdftoppm을 찾을 수 없습니다. poppler-utils를 설치해주세요")
		}
	}

	// PDF 페이지 수 확인
	pageCount, err := api.PageCountFile(pdfPath)
	if err != nil {
		return 0, fmt.Errorf("PDF 페이지 확인 실패: %v", err)
	}

	// pdftoppm 실행: 전체 페이지를 한 번에 변환
	fullOutputPath := fmt.Sprintf("%s/%s", outputDir, outputPrefix)
	cmd := exec.Command(pdftoppmPath, "-jpeg", "-r", "150", pdfPath, fullOutputPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("pdftoppm 실행 실패: %v, output: %s", err, string(output))
	}

	return pageCount, nil
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

	// 임시 디렉토리 생성
	tmpDir := fmt.Sprintf("%v/tmp-%v", config.UploadPath, global.UniqueId())
	os.Mkdir(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	// PDF 전체를 이미지로 변환
	pageCount, err := convertPDFToImages(fullFilename, tmpDir, "page")
	if err != nil {
		errMsg := fmt.Sprintf("PDF 변환 실패 (%s): %v", name, err)
		log.Println(errMsg)
		c.Result["code"] = "error"
		c.Result["message"] = errMsg
		c.Code = http.StatusBadRequest
		return
	}

	var errors []string
	successCount := 0

	for i := 1; i <= pageCount; i++ {
		// 페이지 수에 따라 파일명 형식이 달라짐 (10개 이상이면 zero-padding)
		var tmpImgPath string
		if pageCount >= 10 {
			tmpImgPath = fmt.Sprintf("%s/page-%02d.jpg", tmpDir, i)
		} else {
			tmpImgPath = fmt.Sprintf("%s/page-%d.jpg", tmpDir, i)
		}

		// 생성된 이미지 읽기
		tmpFile, err := os.Open(tmpImgPath)
		if err != nil {
			errMsg := fmt.Sprintf("이미지 읽기 실패 (%s, 페이지 %d): %v", name, i, err)
			log.Println(errMsg)
			errors = append(errors, errMsg)
			continue
		}

		img, _, err := image.Decode(tmpFile)
		tmpFile.Close()
		if err != nil {
			errMsg := fmt.Sprintf("이미지 디코딩 실패 (%s, 페이지 %d): %v", name, i, err)
			log.Println(errMsg)
			errors = append(errors, errMsg)
			continue
		}

		os.Mkdir(fmt.Sprintf("%v/periodicresult/%v", config.UploadPath, id), 0755)

		filename := fmt.Sprintf("managebook-%v.jpg", global.UniqueId())
		fullFilename := fmt.Sprintf("%v/periodicresult/%v/%v", config.UploadPath, id, filename)
		f, err := os.Create(fullFilename)
		if err != nil {
			errMsg := fmt.Sprintf("파일 생성 실패 (%s, 페이지 %d): %v", name, i, err)
			log.Println(errMsg)
			errors = append(errors, errMsg)
			continue
		}

		jpeg.Encode(f, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
		f.Close()

		item := models.Managebook{Filename: filename, Order: i + 1, Periodic: id, Managebookcategory: categoryId, Date: now}
		managebookManager.Insert(&item)

		successCount++
	}

	if successCount == 0 {
		c.Result["code"] = "error"
		c.Result["message"] = fmt.Sprintf("파일 변환 완전 실패: %s (0/%d 페이지 처리됨). 에러: %s", name, pageCount, strings.Join(errors, "; "))
		c.Code = http.StatusBadRequest
		return
	} else if len(errors) > 0 {
		c.Result["code"] = "warning"
		c.Result["message"] = fmt.Sprintf("파일 일부 변환 실패: %s (%d/%d 페이지만 처리됨). 에러: %s", name, successCount, pageCount, strings.Join(errors, "; "))
	}
}

// @Post()
func (c *ManagebookController) Multiprocess(id int64, filename string, originalfilename string) {
	log.Printf("=== Multiprocess START: id=%d, filename=%s, originalfilename=%s ===", id, filename, originalfilename)
	conn := c.NewConnection()

	managebookManager := models.NewManagebookManager(conn)
	managebookcategoryManager := models.NewManagebookcategoryManager(conn)

	filenames := strings.Split(filename, ",")
	originalfilenames := strings.Split(originalfilename, ",")

	var errors []string

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

		// 임시 디렉토리 생성
		tmpDir := fmt.Sprintf("%v/tmp-%v", config.UploadPath, global.UniqueId())
		os.Mkdir(tmpDir, 0755)

		// PDF 전체를 이미지로 변환
		pageCount, err := convertPDFToImages(fullFilename, tmpDir, "page")
		if err != nil {
			errMsg := fmt.Sprintf("PDF 변환 실패 (%s): %v", name, err)
			log.Println(errMsg)
			os.RemoveAll(tmpDir) // 실패 시 즉시 정리
			errors = append(errors, errMsg)
			continue
		}

		successCount := 0
		for i := 1; i <= pageCount; i++ {
			// 페이지 수에 따라 파일명 형식이 달라짐 (10개 이상이면 zero-padding)
			var tmpImgPath string
			if pageCount >= 10 {
				tmpImgPath = fmt.Sprintf("%s/page-%02d.jpg", tmpDir, i)
			} else {
				tmpImgPath = fmt.Sprintf("%s/page-%d.jpg", tmpDir, i)
			}

			// 생성된 이미지 읽기
			tmpFile, err := os.Open(tmpImgPath)
			if err != nil {
				errMsg := fmt.Sprintf("이미지 읽기 실패 (%s, 페이지 %d): %v", name, i, err)
				log.Println(errMsg)
				errors = append(errors, errMsg)
				continue
			}

			img, _, err := image.Decode(tmpFile)
			tmpFile.Close()
			if err != nil {
				errMsg := fmt.Sprintf("이미지 디코딩 실패 (%s, 페이지 %d): %v", name, i, err)
				log.Println(errMsg)
				errors = append(errors, errMsg)
				continue
			}

			os.Mkdir(fmt.Sprintf("%v/periodicresult/%v", config.UploadPath, id), 0755)

			filename := fmt.Sprintf("managebook-%v.jpg", global.UniqueId())
			fullFilename := fmt.Sprintf("%v/periodicresult/%v/%v", config.UploadPath, id, filename)
			f, err := os.Create(fullFilename)
			if err != nil {
				errMsg := fmt.Sprintf("파일 생성 실패 (%s, 페이지 %d): %v", name, i, err)
				log.Println(errMsg)
				errors = append(errors, errMsg)
				continue
			}

			jpeg.Encode(f, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
			f.Close()

			item := models.Managebook{Filename: filename, Order: i + 1, Periodic: id, Managebookcategory: categoryId, Date: now}
			managebookManager.Insert(&item)

			successCount++
		}

		// 처리 완료 후 임시 디렉토리 삭제
		os.RemoveAll(tmpDir)

		if successCount == 0 {
			errors = append(errors, fmt.Sprintf("파일 변환 완전 실패: %s (0/%d 페이지 처리됨)", name, pageCount))
		} else if successCount < pageCount {
			errors = append(errors, fmt.Sprintf("파일 일부 변환 실패: %s (%d/%d 페이지만 처리됨)", name, successCount, pageCount))
		}
	}

	if len(errors) > 0 {
		log.Printf("=== Multiprocess ERRORS: %s ===", strings.Join(errors, "; "))
		c.Result["code"] = "error"
		c.Result["message"] = strings.Join(errors, "; ")
		c.Code = http.StatusBadRequest
		return
	}

	log.Printf("=== Multiprocess SUCCESS: processed %d files ===", len(filenames))
}
