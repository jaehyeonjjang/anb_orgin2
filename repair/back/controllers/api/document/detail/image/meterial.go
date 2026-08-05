package image

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"math"
	"os"
	"repair/global"
	"repair/global/config"
	"repair/models"

	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"golang.org/x/image/font/gofont/goregular"
)

func MakeMeterialImage(periodic int64, blueprint models.Blueprint, items []models.Periodicdata, iconZoom float64, numberZoom float64) {
	draw2d.SetFontFolder("./doc")

	os.Mkdir(fmt.Sprintf("%v/periodicresult/%v", config.UploadPath, periodic), 0755)

	filename := fmt.Sprintf("%v/%v", config.UploadPath, blueprint.Filename)
	img, _ := global.LoadImageFile(filename)

	if img == nil {
		log.Println("not found", filename)
		return
	}

	log.Println(filename)

	extra := 260
	dest, width, height := global.ImageToRGBAWithResize(img, extra, 0)
	gc := draw2dimg.NewGraphicContext(dest)

	step := numberZoom / 2.0

	white := color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	blue := color.RGBA{0x00, 0x00, 0xff, 0xff}
	red := color.RGBA{0xf4, 0x43, 0x37, 0xff}

	// 원본 이미지 오른쪽에 추가된 여백(측정위치 요약 표 자리)을 흰색으로 채움
	gc.BeginPath()
	gc.SetStrokeColor(white)
	gc.SetFillColor(white)
	gc.MoveTo(float64(width), 0)
	gc.LineTo(float64(width+extra), 0)
	gc.LineTo(float64(width+extra), float64(height))
	gc.LineTo(float64(width), float64(height))
	gc.LineTo(float64(width), 0)
	gc.FillStroke()
	gc.BeginPath()

	fontData := draw2d.FontData{
		Name:  "NotoSansKR",
		Style: draw2d.FontStyleNormal,
	}

	fontPath := "./doc/NotoSansKR-Regular.ttf"
	fontBytes, ferr := os.ReadFile(fontPath)
	if ferr != nil {
		// 폰트 파일을 읽지 못하면 라틴 폰트로 대체 (한글은 깨질 수 있음)
		fontBytes = goregular.TTF
		log.Println("NotoSansKR-Regular.ttf not found, fallback to goregular", ferr)
	}
	font, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Fatalf("Failed to parse font: %v", err)
	}
	draw2d.RegisterFont(fontData, font)

	// 표(측정위치)의 연속 번호와 별개로 도면에는 층마다 해당 층 안에서만 1부터 다시 매긴 번호를 표시
	// (요약 표의 "측정위치" 열도 이 번호를 그대로 사용해서 도면의 동그라미 번호와 일치시킨다)
	localNumber := 0
	floor := formatFloorName(blueprint.Name)
	vertical := make([]materialSummaryRow, 0)
	horizontal := make([]materialSummaryRow, 0)

	for _, v := range items {
		// 부재 데이터만 처리 (type 400-499)
		if v.Type < 400 || v.Type >= 500 {
			continue
		}

		log.Println("type", v.Type)
		var results []global.Offset
		json.Unmarshal([]byte(v.Content), &results)

		if len(results) == 0 {
			continue
		}

		point := results[0]
		x := point.Dx
		y := point.Dy

		w := 4.0 * step / 50.0

		if w < 1.5 {
			w = 1.5
		}

		w = 1.5
		gc.SetLineWidth(w)

		if v.Type == 401 || v.Type == 402 {
			localNumber++
			displayNumber := localNumber

			gc.SetFillColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
			if v.Type == 401 {
				gc.SetStrokeColor(red)
			} else {
				gc.SetStrokeColor(blue)
			}

			gc.BeginPath()
			gc.ArcTo(x, y, step, step, 0, math.Pi*2)
			gc.FillStroke()

			if v.Type == 401 {
				gc.SetFillColor(red)
			} else {
				gc.SetFillColor(blue)
			}

			gc.SetFontSize(50 * numberZoom / 100)
			gc.SetFontData(draw2d.FontData{
				Name:   "luxi",
				Family: draw2d.FontFamilyMono,
				Style:  draw2d.FontStyleBold,
			})

			if displayNumber >= 10 {
				gc.FillStringAt(fmt.Sprintf("%v", displayNumber), point.Dx-step*0.7, point.Dy+step/2)
			} else {
				gc.FillStringAt(fmt.Sprintf("%v", displayNumber), point.Dx-step*0.35, point.Dy+step/2)
			}
			gc.Stroke()

			row := materialSummaryRow{
				Position: fmt.Sprintf("%v", displayNumber),
				Floor:    floor,
				T:        fmt.Sprintf("%v", v.Group),
			}

			if v.Type == 401 {
				vertical = append(vertical, row)
			} else {
				horizontal = append(horizontal, row)
			}

		}

	}

	// 현재 층 자신의 측정위치 번호를 위치도 옆 여백에 표로 표시
	drawMaterialSummary(gc, width, height, extra, vertical, horizontal)

	log.Println("make image")

	targetFilename := fmt.Sprintf("%v/periodicresult/%v/%v_400.jpg", config.UploadPath, periodic, blueprint.Id)
	log.Println("targetFilename", targetFilename)
	global.SaveToJpegFile(targetFilename, dest)
}

type materialSummaryRow struct {
	Position string
	Floor    string
	T        string
}

// 위치도 이미지 우측 여백에 측정위치/수직/수평 측정위치(T) 요약 표를 그린다.
func drawMaterialSummary(gc *draw2dimg.GraphicContext, width int, height int, extra int, vertical []materialSummaryRow, horizontal []materialSummaryRow) {
	blue := color.RGBA{0x00, 0x00, 0xff, 0xff}
	red := color.RGBA{0xf4, 0x43, 0x37, 0xff}

	x0 := float64(width) + 15
	xPosition := x0
	xFloor := x0 + 120
	xT := x0 + 220

	rowHeight := 35.0
	y := 200.0

	gc.SetFontData(draw2d.FontData{
		Name:  "NotoSansKR",
		Style: draw2d.FontStyleNormal,
	})
	gc.SetFontSize(25)

	drawTable := func(title string, titleColor color.RGBA, rows []materialSummaryRow) {
		if y > float64(height)-10 {
			return
		}

		gc.SetFillColor(titleColor)
		gc.SetStrokeColor(titleColor)
		fillStringAtCenter(gc, title, xPosition, y)
		y += rowHeight

		fillStringAtCenter(gc, "측정위치", xPosition, y)
		fillStringAtCenter(gc, "층", xFloor, y)
		fillStringAtCenter(gc, "T", xT, y)

		y += rowHeight

		for _, row := range rows {
			if y > float64(height)-10 {
				break
			}

			fillStringAtCenter(gc, row.Position, xPosition, y)
			fillStringAtCenter(gc, row.Floor, xFloor, y)
			fillStringAtCenter(gc, row.T, xT, y)
			y += rowHeight
		}

		y += 10
	}

	drawTable("수직", red, vertical)
	drawTable("수평", blue, horizontal)
}
