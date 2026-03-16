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

// convertPDFPageToImage converts a single PDF page to JPEG using pdftoppm
func convertPDFPageToImage(pdfPath string, pageNum int, outputPrefix string) (string, error) {
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
			return "", fmt.Errorf("pdftoppm을 찾을 수 없습니다. poppler-utils를 설치해주세요")
		}
	}

	// pdftoppm 실행: -jpeg -r 150 -f pageNum -l pageNum input.pdf output_prefix
	cmd := exec.Command(pdftoppmPath, "-jpeg", "-r", "150", "-f", fmt.Sprintf("%d", pageNum), "-l", fmt.Sprintf("%d", pageNum), pdfPath, outputPrefix)
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("pdftoppm 실행 실패: %v, output: %s", err, string(output))
	}

	// 생성된 파일 찾기 (pdftoppm은 output_prefix-XX.jpg 형식으로 저장)
	outputFile := fmt.Sprintf("%s-%02d.jpg", outputPrefix, pageNum)
	if _, err := os.Stat(outputFile); err != nil {
		return "", fmt.Errorf("변환된 이미지 파일을 찾을 수 없습니다: %s", outputFile)
	}

	return outputFile, nil
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
		errMsg := fmt.Sprintf("PDF 페이지 확인 실패 (%s): %v", name, err)
		log.Println(errMsg)
		c.Result["code"] = "error"
		c.Result["message"] = errMsg
		c.Code = http.StatusBadRequest
		return
	}

	// 임시 디렉토리 생성
	tmpDir := fmt.Sprintf("%v/tmp-%v", config.UploadPath, global.UniqueId())
	os.Mkdir(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	var errors []string
	successCount := 0

	for i := 1; i <= pageCount; i++ {
		// pdftoppm을 사용하여 PDF 페이지를 이미지로 변환
		tmpOutputPrefix := fmt.Sprintf("%v/page", tmpDir)
		tmpImgPath, err := convertPDFPageToImage(fullFilename, i, tmpOutputPrefix)
		if err != nil {
			errMsg := fmt.Sprintf("PDF 렌더링 실패 (%s, 페이지 %d): %v", name, i, err)
			log.Println(errMsg)
			errors = append(errors, errMsg)
			continue
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
		os.Remove(tmpImgPath)
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

		// PDF 페이지 수 확인
		pageCount, err := api.PageCountFile(fullFilename)
		if err != nil {
			errMsg := fmt.Sprintf("PDF 페이지 확인 실패 (%s): %v", name, err)
			log.Println(errMsg)
			errors = append(errors, errMsg)
			continue
		}

		// 임시 디렉토리 생성
		tmpDir := fmt.Sprintf("%v/tmp-%v", config.UploadPath, global.UniqueId())
		os.Mkdir(tmpDir, 0755)
		defer os.RemoveAll(tmpDir)

		successCount := 0
		for i := 1; i <= pageCount; i++ {
			// pdftoppm을 사용하여 PDF 페이지를 이미지로 변환
			tmpOutputPrefix := fmt.Sprintf("%v/page", tmpDir)
			tmpImgPath, err := convertPDFPageToImage(fullFilename, i, tmpOutputPrefix)
			if err != nil {
				errMsg := fmt.Sprintf("PDF 렌더링 실패 (%s, 페이지 %d): %v", name, i, err)
				log.Println(errMsg)
				errors = append(errors, errMsg)
				continue
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
