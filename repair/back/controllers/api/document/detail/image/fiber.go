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
	"strings"

	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"golang.org/x/image/font/gofont/goregular"
)

func MakeFiberImage(periodic int64, blueprint models.Blueprint, items []models.Periodicdata, iconZoom float64, numberZoom float64) {
	draw2d.SetFontFolder("./doc")

	os.Mkdir(fmt.Sprintf("%v/periodicresult/%v", config.UploadPath, periodic), 0755)

	filename := fmt.Sprintf("%v/%v", config.UploadPath, blueprint.Filename)
	img, _ := global.LoadImageFile(filename)

	if img == nil {
		log.Println("not found", filename)
		return
	}

	extra := 260
	dest, width, height := global.ImageToRGBAWithResize(img, extra, 0)
	gc := draw2dimg.NewGraphicContext(dest)

	step := iconZoom / 2.0
	stepMiddle := step * 0.66

	white := color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	blue := color.RGBA{0x00, 0x00, 0xff, 0xff}
	red := color.RGBA{0xf4, 0x43, 0x37, 0xff}

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
	gc.SetFillColor(red)
	gc.SetStrokeColor(red)

	fontData := draw2d.FontData(draw2d.FontData{
		Name:  "NotoSansKR",
		Style: draw2d.FontStyleNormal,
	})

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

	gc.SetFontData(fontData)

	for _, v := range items {
		// 강도/탄산화 데이터만 처리 (type 300-399)
		if v.Type < 300 || v.Type >= 400 {
			continue
		}

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

		if v.Type == 301 {
			gc.BeginPath()

			gc.SetStrokeColor(red)
			gc.SetFillColor(red)

			gc.ArcTo(x, y, stepMiddle, stepMiddle, 0, math.Pi*2)

			gc.FillStroke()
		} else if v.Type == 302 {
			gc.BeginPath()

			gc.SetStrokeColor(blue)
			gc.SetFillColor(blue)

			gc.ArcTo(x, y, stepMiddle, stepMiddle, 0, math.Pi*2)

			gc.FillStroke()
		}

	}

	// 현재 층 자신의 SH/N 값을 위치도 옆 여백에 표로 표시
	vertical, horizontal := getFiberSummaryRows(periodic, blueprint)
	drawFiberSummary(gc, width, height, extra, vertical, horizontal)

	targetFilename := fmt.Sprintf("%v/periodicresult/%v/%v_300.jpg", config.UploadPath, periodic, blueprint.Id)
	global.SaveToJpegFile(targetFilename, dest)
}

type fiberSummaryRow struct {
	Floor string
	SH    string
	N     string
}

// 현재 층(blueprint) 자신의 강도/탄산화 데이터만 조회하여
// 수직(301)/수평(302)으로 나눠 SH, N 값을 반환한다.
func getFiberSummaryRows(periodic int64, blueprint models.Blueprint) ([]fiberSummaryRow, []fiberSummaryRow) {
	periodicdataManager := models.NewPeriodicdataManager(nil)
	defer periodicdataManager.Close()

	conditions := []interface{}{
		models.Where{Column: "periodic", Value: periodic, Compare: "="},
		models.Where{Column: "blueprint", Value: blueprint.Id, Compare: "="},
		models.Where{Column: "type", Value: 300, Compare: ">"},
		models.Where{Column: "type", Value: 400, Compare: "<"},
		models.Ordering("pd_group"),
	}

	items := periodicdataManager.Find(conditions)

	vertical := make([]fiberSummaryRow, 0)
	horizontal := make([]fiberSummaryRow, 0)

	floor := formatFloorName(blueprint.Name)

	for _, v := range items {
		sh := ""
		if v.Shape != "" {
			sh = fmt.Sprintf("%v, %v", v.Shape, v.Length)
		}

		row := fiberSummaryRow{Floor: floor, SH: sh, N: v.Width}

		if v.Type == 301 {
			vertical = append(vertical, row)
		} else if v.Type == 302 {
			horizontal = append(horizontal, row)
		}
	}

	return vertical, horizontal
}

// 위치도 이미지 우측 여백에 수직/수평 SH, N 요약 표를 그린다.
func drawFiberSummary(gc *draw2dimg.GraphicContext, width int, height int, extra int, vertical []fiberSummaryRow, horizontal []fiberSummaryRow) {
	blue := color.RGBA{0x00, 0x00, 0xff, 0xff}
	red := color.RGBA{0xf4, 0x43, 0x37, 0xff}

	x0 := float64(width) + 15
	xFloor := x0
	xSH := x0 + 90
	xN := x0 + 190

	rowHeight := 35.0
	y := 200.0

	gc.SetFontSize(25)

	drawTable := func(title string, titleColor color.RGBA, rows []fiberSummaryRow) {
		if y > float64(height)-10 {
			return
		}

		gc.SetFillColor(titleColor)
		gc.SetStrokeColor(titleColor)
		fillStringAtCenter(gc, title, xFloor, y)
		y += rowHeight

		fillStringAtCenter(gc, "층", xFloor, y)
		fillStringAtCenter(gc, "SH", xSH, y)
		fillStringAtCenter(gc, "N", xN, y)

		y += rowHeight

		for _, row := range rows {
			if y > float64(height)-10 {
				break
			}

			fillStringAtCenter(gc, row.Floor, xFloor, y)
			fillStringAtCenter(gc, row.SH, xSH, y)
			fillStringAtCenter(gc, row.N, xN, y)
			y += rowHeight
		}

		y += 10
	}

	drawTable("수직", red, vertical)
	drawTable("수평", blue, horizontal)
}

// 지정한 x 좌표를 중심으로 텍스트를 가운데 정렬해서 그린다.
func fillStringAtCenter(gc *draw2dimg.GraphicContext, text string, centerX float64, y float64) {
	left, _, right, _ := gc.GetStringBounds(text)
	width := right - left
	gc.FillStringAt(text, centerX-width/2-left, y)
}

// 위치도 데이터표에 표시할 층 이름을 다듬는다.
// "지하1층" -> "지하1" -> "B1" 처럼 지하층은 B# 표기로 바꾼다.
func formatFloorName(name string) string {
	floor := strings.TrimSuffix(name, "층")
	if strings.HasPrefix(floor, "지하") {
		floor = "B" + strings.TrimPrefix(floor, "지하")
	}
	return floor
}
