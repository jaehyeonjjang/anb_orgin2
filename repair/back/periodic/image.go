package periodic

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

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
)

const (
	basicVerticalLine   = 3
	basicHorizontalLine = 4

	basicVerticalBreak   = 5
	basicHorizontalBreak = 6

	basicVerticalLineV = 7 // 세로 시작 직선

	basicVerticalBreakV = 9 // 세로 시작 꺾은선

	basicVerticalLineFree  = 11 // 자유 직선(정)
	basicVerticalBreakFree = 13 // 자유 꺾은선(역)

	curveBlue   = 31
	curveRed    = 32
	curveGreen  = 33
	curveViolet = 34
	lineBlue    = 41
	lineRed     = 42
	lineGreen   = 43
	lineViolet  = 44

	dashedRectRed  = 51
	dashedRectBlue = 52

	inclinationLine       = 201
	inclinationHorizontal = 202
	inclinationVertical   = 203

	fiberVertical   = 301
	fiberHorizontal = 302

	materialVertical   = 401
	materialHorizontal = 402

	crackLineRed    = 121
	crackLineBlue   = 122
	crackLineViolet = 123
	crackLineGreen  = 124

	crackCurveRed    = 126
	crackCurveBlue   = 127
	crackCurveViolet = 128
	crackCurveGreen  = 129

	rebarXViolet = 135

	corrosionGreen  = 136
	corrosionViolet = 137

	beamGreen  = 138
	beamViolet = 139

	otherGreen  = 140
	otherViolet = 141

	pipeLeakGreen  = 142
	pipeLeakViolet = 143

	leakRed   = 144
	leakGreen = 145
)

func IsCrack(v int) bool {
	if v == crackLineRed || v == crackLineBlue || v == crackLineViolet || v == crackLineGreen ||
		v == crackCurveRed || v == crackCurveBlue || v == crackCurveViolet || v == crackCurveGreen {
		return true
	}

	return false
}

func IsCrackCurve(v int) bool {
	if v == crackCurveRed || v == crackCurveBlue || v == crackCurveViolet || v == crackCurveGreen {
		return true
	}

	return false
}

func MakeImage(periodic int64, blueprint models.Blueprint, items []models.Periodicdata, iconZoom float64, numberZoom float64, crackZoom float64) {
	draw2d.SetFontFolder("./doc")

	os.MkdirAll(fmt.Sprintf("%v/periodicresult/%v", config.UploadPath, periodic), 0755)

	filename := fmt.Sprintf("%v/%v", config.UploadPath, blueprint.Filename)
	img, _ := global.LoadImageFile(filename)

	if img == nil {
		log.Println("not found", filename)
		return
	}

	dest := global.ImageToRGBA(img)
	gc := draw2dimg.NewGraphicContext(dest)

	step := iconZoom / 2.0
	stepMiddle := step * 0.66
	stepSmall := step * 0.33

	if numberZoom <= 0.0 {
		numberZoom = iconZoom
	}

	stepNumber := numberZoom / 2.0

	if crackZoom <= 0.0 {
		crackZoom = iconZoom
	}

	log.Println("crackzoom", crackZoom)

	stepCrack := crackZoom / 2.0
	// stepCrackMiddle := stepCrack * 0.66
	stepCrackSmall := stepCrack * 0.33

	// lightblue := color.RGBA{0x1f, 0x96, 0xf3, 0xff}
	lightblue := color.RGBA{0x00, 0x00, 0xff, 0xff}
	blue := color.RGBA{0x00, 0x00, 0xff, 0xff}
	// red := color.RGBA{0xf4, 0x43, 0x37, 0xff}
	red := color.RGBA{0xff, 0x00, 0x00, 0xff}
	green := color.RGBA{0x22, 0x8B, 0x22, 0xff}
	violet := color.RGBA{0xa0, 0x00, 0xa0, 0xff}

	for _, v := range items {
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

		if v.Type == 1 || v.Type == 2 {
			w := 4.0 * stepNumber / 50.0

			if w < 1.5 {
				w = 1.5
			}

			w = 1.5
			gc.SetLineWidth(w)

			gc.SetFillColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
			if v.Type == 1 {
				gc.SetStrokeColor(red)
			} else {
				gc.SetStrokeColor(blue)
			}

			gc.BeginPath()
			gc.ArcTo(x, y, stepNumber, stepNumber, 0, math.Pi*2)
			gc.FillStroke()

			if v.Type == 1 {
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

			if v.Group >= 10 {
				gc.FillStringAt(fmt.Sprintf("%v", v.Group), point.Dx-stepNumber*0.7, point.Dy+stepNumber/2)
			} else {
				gc.FillStringAt(fmt.Sprintf("%v", v.Group), point.Dx-stepNumber*0.35, point.Dy+stepNumber/2)
			}
			gc.Stroke()
		} else if v.Type == basicVerticalLine || v.Type == basicHorizontalLine || v.Type == basicVerticalBreak || v.Type == basicHorizontalBreak || v.Type == basicVerticalLineV || v.Type == basicVerticalBreakV || v.Type == basicVerticalLineFree || v.Type == basicVerticalBreakFree {
			w := 4.0 * stepNumber / 50.0

			if w < 1.5 {
				w = 1.5
			}

			w = 1.0
			gc.SetLineWidth(w)

			if v.Type == basicVerticalLine || v.Type == basicVerticalBreak || v.Type == basicVerticalLineV || v.Type == basicVerticalBreakV || v.Type == basicVerticalLineFree || v.Type == basicVerticalBreakFree {
				gc.SetStrokeColor(red)
				gc.SetFillColor(red)
			} else {
				gc.SetStrokeColor(blue)
				gc.SetFillColor(blue)
			}

			y1 := results[0].Dy
			y2 := results[0].Dy
			x1 := results[0].Dx
			x2 := results[0].Dx

			if len(results) > 1 {
				y2 = results[1].Dy
				x2 = results[1].Dx
			}

			var dx float64 = 0.0
			var dy float64 = 0.0

			for i := 2; i < len(results); i++ {
				if x1 == x2 && y1 == y2 {
					y2 = results[i].Dy
					x2 = results[i].Dx
				}
			}

			sy := y2 - y1
			sx := x2 - x1

			angle := math.Atan(sy/sx) * (180.0 / 3.14)

			angle += 180

			if sx < 0.0 {
				angle += 180.0
			} else {
				if sy < 0.0 {
					angle += 360.0
				}
			}

			// 기본 직선/꺾은선 타입들은 가로 방향으로 고정
			if v.Type == basicVerticalLine || v.Type == basicVerticalBreak {
				if len(results) >= 2 {
					if results[0].Dx > results[len(results)-1].Dx {
						angle = 0
					} else {
						angle = 180
					}
				}
			} else if v.Type == basicVerticalLineV || v.Type == basicVerticalBreakV {
				// 세로 시작 타입들은 세로 방향으로 고정
				if len(results) >= 2 {
					if results[0].Dy > results[len(results)-1].Dy {
						angle = 90
					} else {
						angle = 270
					}
				}
			} else if v.Type == basicHorizontalLine {
				// 수평부재 가로: 가로 방향 화살표
				if len(results) >= 2 {
					if results[0].Dx > results[len(results)-1].Dx {
						angle = 0
					} else {
						angle = 180
					}
				}
			} else if v.Type == basicHorizontalBreak {
				// 수평부재 세로: 세로 방향 화살표
				if len(results) >= 2 {
					if results[0].Dy > results[len(results)-1].Dy {
						angle = 90
					} else {
						angle = 270
					}
				}
			}
			// 자유선(11, 13)은 자동 각도 계산 사용 (위의 angle 값 그대로 사용)

			r := (numberZoom / 2.0) * 0.5 / 2.0
			zoom := 1.0

			xinc := math.Cos(angle*3.14/180) * r
			yinc := math.Sin(angle*3.14/180) * r

			originalX := x1
			originalY := y1

			x1 -= xinc
			y1 -= yinc

			if v.Type == basicVerticalBreak || v.Type == basicVerticalBreakV || v.Type == basicVerticalBreakFree {
				// 반대선
				gc.BeginPath()
				cx := math.Cos((angle+90)*3.14/180) * r
				cy := math.Sin((angle+90)*3.14/180) * r
				gc.MoveTo((x1+xinc+cx)*zoom+dx, (y1+yinc+cy)*zoom+dy)

				cx = math.Cos((angle+90+180)*3.14/180) * r
				cy = math.Sin((angle+90+180)*3.14/180) * r
				gc.LineTo((x1+xinc+cx)*zoom+dx, (y1+yinc+cy)*zoom+dy)

				// 연장선
				cx = math.Cos((angle)*3.14/180) * r * 2
				cy = math.Sin((angle)*3.14/180) * r * 2
				gc.MoveTo((x1+xinc+cx)*zoom+dx, (y1+yinc+cy)*zoom+dy)

				cx = math.Cos((angle+180)*3.14/180) * r * 2
				cy = math.Sin((angle+180)*3.14/180) * r * 2
				gc.LineTo((x1+xinc+cx)*zoom+dx, (y1+yinc+cy)*zoom+dy)

				gc.Stroke()

				gc.BeginPath()

				angle += 180

				xinc = math.Cos(angle*3.14/180) * r
				yinc = math.Sin(angle*3.14/180) * r

				x1 = originalX - xinc
				y1 = originalY - yinc

				dx = 0
				dy = 0
				gc.MoveTo((x1+xinc)*zoom+dx, (y1+yinc)*zoom+dy)

				angle += 120

				xinc = math.Cos(angle*3.14/180) * r
				yinc = math.Sin(angle*3.14/180) * r

				gc.LineTo((x1+xinc)*zoom+dx, (y1+yinc)*zoom+dy)

				gc.LineTo((x1)*zoom+dx, (y1)*zoom+dy)

				angle += 120

				xinc = math.Cos(angle*3.14/180) * r
				yinc = math.Sin(angle*3.14/180) * r

				gc.LineTo((x1+xinc)*zoom+dx, (y1+yinc)*zoom+dy)

				angle += 120

				xinc = math.Cos(angle*3.14/180) * r
				yinc = math.Sin(angle*3.14/180) * r

				gc.LineTo((x1+xinc)*zoom+dx, (y1+yinc)*zoom+dy)

				gc.FillStroke()
			} else {
				gc.BeginPath()

				dx = 0
				dy = 0
				gc.MoveTo((x1+xinc)*zoom+dx, (y1+yinc)*zoom+dy)

				angle += 120

				xinc = math.Cos(angle*3.14/180) * r
				yinc = math.Sin(angle*3.14/180) * r

				gc.LineTo((x1+xinc)*zoom+dx, (y1+yinc)*zoom+dy)

				gc.LineTo((x1)*zoom+dx, (y1)*zoom+dy)

				angle += 120

				xinc = math.Cos(angle*3.14/180) * r
				yinc = math.Sin(angle*3.14/180) * r

				gc.LineTo((x1+xinc)*zoom+dx, (y1+yinc)*zoom+dy)

				angle += 120

				xinc = math.Cos(angle*3.14/180) * r
				yinc = math.Sin(angle*3.14/180) * r

				gc.LineTo((x1+xinc)*zoom+dx, (y1+yinc)*zoom+dy)

				gc.FillStroke()
			}

			gc.BeginPath()

			// 타입별로 선 그리기
			if v.Type == basicVerticalLine || v.Type == basicVerticalBreak {
				// 타입 3, 5: 항상 L자 (가로→세로) - 첫점과 끝점만 사용
				endPoint := results[len(results)-1]
				gc.MoveTo(results[0].Dx, results[0].Dy)
				gc.LineTo(endPoint.Dx, results[0].Dy) // 가로선
				gc.LineTo(endPoint.Dx, endPoint.Dy)   // 세로선
			} else if v.Type == basicVerticalLineV || v.Type == basicVerticalBreakV {
				// 타입 7, 9: 항상 L자 (세로→가로) - 첫점과 끝점만 사용
				endPoint := results[len(results)-1]
				gc.MoveTo(results[0].Dx, results[0].Dy)
				gc.LineTo(results[0].Dx, endPoint.Dy) // 세로선
				gc.LineTo(endPoint.Dx, endPoint.Dy)   // 가로선
			} else if v.Type == basicHorizontalLine {
				// 타입 4: 항상 L자 (가로→세로) - 첫점과 끝점만 사용
				endPoint := results[len(results)-1]
				gc.MoveTo(results[0].Dx, results[0].Dy)
				gc.LineTo(endPoint.Dx, results[0].Dy) // 가로선
				gc.LineTo(endPoint.Dx, endPoint.Dy)   // 세로선
			} else if v.Type == basicHorizontalBreak {
				// 타입 6: 항상 L자 (세로→가로) - 첫점과 끝점만 사용
				endPoint := results[len(results)-1]
				gc.MoveTo(results[0].Dx, results[0].Dy)
				gc.LineTo(results[0].Dx, endPoint.Dy) // 세로선
				gc.LineTo(endPoint.Dx, endPoint.Dy)   // 가로선
			} else {
				// 자유선(11, 13): 모든 점을 순서대로 연결
				for i, point := range results {
					if i == 0 {
						gc.MoveTo(point.Dx, point.Dy)
					} else {
						gc.LineTo(point.Dx, point.Dy)
					}
				}
			}

			if v.Type == basicVerticalLine || v.Type == basicVerticalBreak || v.Type == basicVerticalLineV || v.Type == basicVerticalBreakV || v.Type == basicVerticalLineFree || v.Type == basicVerticalBreakFree {
				gc.SetStrokeColor(red)
			} else if v.Type == basicHorizontalLine || v.Type == basicHorizontalBreak {
				gc.SetStrokeColor(blue)
			} else {
				gc.SetStrokeColor(blue)
			}

			gc.Stroke() // 선 그리기!

			// 끝점에 숫자가 들어간 원 그리기
			endPoint := results[len(results)-1]

			// 원 그리기
			gc.BeginPath()
			gc.SetFillColor(color.RGBA{0xff, 0xff, 0xff, 0xff}) // 흰색 배경
			if v.Type == basicVerticalLine || v.Type == basicVerticalBreak || v.Type == basicVerticalLineV || v.Type == basicVerticalBreakV || v.Type == basicVerticalLineFree || v.Type == basicVerticalBreakFree {
				gc.SetStrokeColor(red)
			} else if v.Type == basicHorizontalLine || v.Type == basicHorizontalBreak {
				gc.SetStrokeColor(blue)
			} else {
				gc.SetStrokeColor(blue)
			}
			gc.ArcTo(endPoint.Dx, endPoint.Dy, stepNumber, stepNumber, 0, math.Pi*2)
			gc.FillStroke()

			// 숫자 표시
			if v.Type == basicVerticalLine || v.Type == basicVerticalBreak || v.Type == basicVerticalLineV || v.Type == basicVerticalBreakV || v.Type == basicVerticalLineFree || v.Type == basicVerticalBreakFree {
				gc.SetFillColor(red)
			} else if v.Type == basicHorizontalLine || v.Type == basicHorizontalBreak {
				gc.SetFillColor(blue)
			} else {
				gc.SetFillColor(blue)
			}

			gc.SetFontSize(50 * numberZoom / 100)
			gc.SetFontData(draw2d.FontData{
				Name:   "Noto Sans KR",
				Family: draw2d.FontFamilyMono,
				Style:  draw2d.FontStyleBold,
			})

			if v.Group >= 10 {
				gc.FillStringAt(fmt.Sprintf("%v", v.Group), endPoint.Dx-stepNumber*0.7, endPoint.Dy+stepNumber/2)
			} else {
				gc.FillStringAt(fmt.Sprintf("%v", v.Group), endPoint.Dx-stepNumber*0.35, endPoint.Dy+stepNumber/2)
			}

			// 텍스트 표시 (화살표 반대편)
			txt := v.Member
			if v.Shape != "" {
				txt = fmt.Sprintf("%v(%v)", txt, v.Shape)
			}

			// Flutter 앱과 동일하게 끝점 기준으로 텍스트 배치
			// Flutter: if (points.items[0].dx > points.items[points.items.length - 1].dx)
			if results[0].Dx > results[len(results)-1].Dx {
				// 오른쪽에서 왼쪽으로: 텍스트는 왼쪽에
				gc.FillStringAt(txt, endPoint.Dx-stepNumber*1.3-float64(len(txt))*stepNumber/2.5, endPoint.Dy+stepNumber/2)
			} else {
				// 왼쪽에서 오른쪽으로: 텍스트는 오른쪽에
				gc.FillStringAt(txt, endPoint.Dx+stepNumber*1.3, endPoint.Dy+stepNumber/2)
			}

			gc.Stroke()
		} else if v.Type >= 31 && v.Type <= 44 {
			if v.Type == 31 || v.Type == 41 {
				gc.SetStrokeColor(lightblue)
			} else if v.Type == 32 || v.Type == 42 {
				gc.SetStrokeColor(red)
			} else if v.Type == 33 || v.Type == 43 {
				gc.SetStrokeColor(green)
			} else if v.Type == 34 || v.Type == 44 {
				gc.SetStrokeColor(violet)
			}

			gc.BeginPath()

			for i, point := range results {
				if i == 0 {
					gc.MoveTo(point.Dx, point.Dy)
				} else {
					gc.LineTo(point.Dx, point.Dy)
				}
			}

			gc.Stroke()
		} else if v.Type == dashedRectRed || v.Type == dashedRectBlue {
			// 크기 조절 가능한 점선 사각형: 두 끝점이 대각선 모서리
			if len(results) >= 2 {
				p1 := results[0]
				p2 := results[len(results)-1]

				if v.Type == dashedRectBlue {
					gc.SetStrokeColor(blue)
				} else {
					gc.SetStrokeColor(red)
				}
				gc.SetLineWidth(2.0)

				x1 := p1.Dx
				y1 := p1.Dy
				x2 := p2.Dx
				y2 := p2.Dy
				if x1 > x2 {
					x1, x2 = x2, x1
				}
				if y1 > y2 {
					y1, y2 = y2, y1
				}

				dashLen := 10.0
				gapLen := 6.0

				drawDashed := func(sx, sy, ex, ey float64) {
					ddx := ex - sx
					ddy := ey - sy
					dist := math.Sqrt(ddx*ddx + ddy*ddy)
					if dist <= 0 {
						return
					}
					ux := ddx / dist
					uy := ddy / dist
					drawn := 0.0
					for drawn < dist {
						segLen := dashLen
						if dist-drawn < dashLen {
							segLen = dist - drawn
						}
						gc.BeginPath()
						gc.MoveTo(sx+ux*drawn, sy+uy*drawn)
						gc.LineTo(sx+ux*(drawn+segLen), sy+uy*(drawn+segLen))
						gc.Stroke()
						drawn += dashLen + gapLen
					}
				}

				drawDashed(x1, y1, x2, y1) // top
				drawDashed(x2, y1, x2, y2) // right
				drawDashed(x1, y2, x2, y2) // bottom
				drawDashed(x1, y1, x1, y2) // left
			}
		} else if v.Type == 101 || v.Type == 130 || v.Type == 131 || v.Type == rebarXViolet {
			gc.BeginPath()

			if v.Type == 101 {
				gc.SetStrokeColor(lightblue)
				gc.SetFillColor(lightblue)
			} else if v.Type == 130 {
				gc.SetStrokeColor(red)
				gc.SetFillColor(red)
			} else if v.Type == 131 {
				gc.SetStrokeColor(green)
				gc.SetFillColor(green)
			} else if v.Type == rebarXViolet {
				gc.SetStrokeColor(violet)
				gc.SetFillColor(violet)
			}

			gc.MoveTo(x-step, y-step)
			gc.LineTo(x+step, y+step)

			gc.MoveTo(x+step, y-step)
			gc.LineTo(x-step, y+step)

			gc.Stroke()
		} else if v.Type == 102 || v.Type == 132 || v.Type == corrosionGreen || v.Type == corrosionViolet {
			gc.BeginPath()

			if v.Type == 102 {
				gc.SetStrokeColor(lightblue)
				gc.SetFillColor(lightblue)
			} else if v.Type == corrosionGreen {
				gc.SetStrokeColor(green)
				gc.SetFillColor(green)
			} else if v.Type == corrosionViolet {
				gc.SetStrokeColor(violet)
				gc.SetFillColor(violet)
			} else {
				gc.SetStrokeColor(red)
				gc.SetFillColor(red)
			}

			gc.MoveTo(x-step, y)
			gc.LineTo(x, y-step)
			gc.LineTo(x+step, y)
			gc.LineTo(x, y+step)
			gc.LineTo(x-step, y)

			gc.Stroke()
		} else if v.Type == 103 || v.Type == 133 || v.Type == beamGreen || v.Type == beamViolet {
			gc.BeginPath()

			if v.Type == 103 {
				gc.SetStrokeColor(lightblue)
				gc.SetFillColor(lightblue)
			} else if v.Type == beamGreen {
				gc.SetStrokeColor(green)
				gc.SetFillColor(green)
			} else if v.Type == beamViolet {
				gc.SetStrokeColor(violet)
				gc.SetFillColor(violet)
			} else {
				gc.SetStrokeColor(red)
				gc.SetFillColor(red)
			}

			gc.MoveTo(x+math.Cos(30.0*3.14/180.0)*step, y+math.Sin(30.0*3.14/180.0)*step)
			gc.LineTo(x+math.Cos(150.0*3.14/180.0)*step, y+math.Sin(150.0*3.14/180.0)*step)
			gc.LineTo(x+math.Cos(270.0*3.14/180.0)*step, y+math.Sin(270.0*3.14/180.0)*step)
			gc.LineTo(x+math.Cos(30.0*3.14/180.0)*step, y+math.Sin(30.0*3.14/180.0)*step)

			gc.Fill()
		} else if v.Type == 104 || v.Type == 134 || v.Type == otherGreen || v.Type == otherViolet {
			gc.BeginPath()

			if v.Type == 104 {
				gc.SetStrokeColor(red)
				gc.SetFillColor(red)
			} else if v.Type == otherGreen {
				gc.SetStrokeColor(green)
				gc.SetFillColor(green)
			} else if v.Type == otherViolet {
				gc.SetStrokeColor(violet)
				gc.SetFillColor(violet)
			} else {
				gc.SetStrokeColor(lightblue)
				gc.SetFillColor(lightblue)
			}

			gc.MoveTo(x-step, y-step)
			gc.LineTo(x+step, y-step)
			gc.LineTo(x+step, y+step)
			gc.LineTo(x-step, y+step)
			gc.LineTo(x-step, y-step)

			gc.Fill()
		} else if v.Type == 105 {
			gc.BeginPath()

			gc.SetStrokeColor(blue)
			gc.SetFillColor(blue)

			gc.ArcTo(x, y, stepMiddle, stepMiddle, 0, math.Pi*2)

			gc.Stroke()

		} else if v.Type == 106 {
			gc.BeginPath()

			gc.SetStrokeColor(blue)
			gc.SetFillColor(blue)

			gc.ArcTo(x, y, step, step, 0, math.Pi*2)

			gc.MoveTo(x+math.Cos(30.0*3.14/180.0)*step, y+math.Sin(30.0*3.14/180.0)*step)
			gc.LineTo(x+math.Cos(150.0*3.14/180.0)*step, y+math.Sin(150.0*3.14/180.0)*step)
			gc.LineTo(x+math.Cos(270.0*3.14/180.0)*step, y+math.Sin(270.0*3.14/180.0)*step)
			gc.LineTo(x+math.Cos(30.0*3.14/180.0)*step, y+math.Sin(30.0*3.14/180.0)*step)

			gc.Stroke()
		} else if v.Type == 107 || v.Type == 108 || v.Type == 109 || v.Type == 110 {
			gc.BeginPath()

			gc.SetStrokeColor(blue)
			gc.SetFillColor(blue)

			gc.ArcTo(x, y, stepSmall, stepSmall, 0, math.Pi*2)

			var angles = []float64{135.0, 90.0, 45.0, 0.0}
			var angle = angles[v.Type-107]

			gc.MoveTo(x+math.Cos(angle*3.14/180.0)*step,
				y+math.Sin(angle*3.14/180.0)*step)
			angle += 180.0
			gc.LineTo(x+math.Cos(angle*3.14/180.0)*step,
				y+math.Sin(angle*3.14/180.0)*step)

			gc.Stroke()
		} else if v.Type == 115 {
			gc.BeginPath()

			gc.SetStrokeColor(violet)
			gc.SetFillColor(violet)

			gc.ArcTo(x, y, stepMiddle, stepMiddle, 0, math.Pi*2)

			gc.Stroke()

		} else if v.Type == 116 {
			gc.BeginPath()

			gc.SetStrokeColor(violet)
			gc.SetFillColor(violet)

			gc.ArcTo(x, y, step, step, 0, math.Pi*2)

			gc.MoveTo(x+math.Cos(30.0*3.14/180.0)*step, y+math.Sin(30.0*3.14/180.0)*step)
			gc.LineTo(x+math.Cos(150.0*3.14/180.0)*step, y+math.Sin(150.0*3.14/180.0)*step)
			gc.LineTo(x+math.Cos(270.0*3.14/180.0)*step, y+math.Sin(270.0*3.14/180.0)*step)
			gc.LineTo(x+math.Cos(30.0*3.14/180.0)*step, y+math.Sin(30.0*3.14/180.0)*step)

			gc.Stroke()
		} else if v.Type == 117 || v.Type == 118 || v.Type == 119 || v.Type == 120 {
			gc.BeginPath()

			gc.SetStrokeColor(violet)
			gc.SetFillColor(violet)

			gc.ArcTo(x, y, stepSmall, stepSmall, 0, math.Pi*2)

			var angles = []float64{135.0, 90.0, 45.0, 0.0}
			var angle = angles[v.Type-117]

			gc.MoveTo(x+math.Cos(angle*3.14/180.0)*step,
				y+math.Sin(angle*3.14/180.0)*step)
			angle += 180.0
			gc.LineTo(x+math.Cos(angle*3.14/180.0)*step,
				y+math.Sin(angle*3.14/180.0)*step)

			gc.Stroke()
		} else if IsCrack(v.Type) {
			point1 := results[0]
			x1 := point1.Dx
			y1 := point1.Dy
			point2 := results[0]
			if len(results) > 1 {
				point2 = results[len(results)-1]
			}
			x2 := point2.Dx
			y2 := point2.Dy

			gc.BeginPath()

			if v.Type == crackLineRed || v.Type == crackCurveRed {
				gc.SetStrokeColor(red)
				gc.SetFillColor(red)
			} else if v.Type == crackLineBlue || v.Type == crackCurveBlue {
				gc.SetStrokeColor(blue)
				gc.SetFillColor(blue)
			} else if v.Type == crackLineViolet || v.Type == crackCurveViolet {
				gc.SetStrokeColor(violet)
				gc.SetFillColor(violet)
			} else if v.Type == crackLineGreen || v.Type == crackCurveGreen {
				gc.SetStrokeColor(green)
				gc.SetFillColor(green)
			}

			for i, v := range results {
				if i == 0 {
					gc.MoveTo(v.Dx, v.Dy)
				} else {
					gc.LineTo(v.Dx, v.Dy)
				}
			}

			if x1 > x2 {
				x = x2 + (x1-x2)/2
			} else {
				x = x1 + (x2-x1)/2
			}

			if y1 > y2 {
				y = y2 + (y1-y2)/2
			} else {
				y = y1 + (y2-y1)/2
			}

			if IsCrackCurve(v.Type) && len(results) > 2 {
				lengths := make([]float64, 0)
				var total float64 = 0.0
				old := results[0]
				for i, cu := range results {
					if i == 0 {
						continue
					}
					length := math.Sqrt(math.Abs(old.Dx-cu.Dx) + math.Abs(old.Dy-cu.Dy))
					lengths = append(lengths, length)
					total += length

					old = cu
				}

				var half float64 = total / 2
				var current float64 = 0.0

				for i, v := range lengths {
					current += v

					if current >= half {
						x1 = results[i].Dx
						y1 = results[i].Dy
						x2 = results[i+1].Dx
						y2 = results[i+1].Dy

						x = x1
						y = y1
						if x1 > x2 {
							x = x2 + (x1-x2)/2
						} else {
							x = x1 + (x2-x1)/2
						}

						if y1 > y2 {
							y = y2 + (y1-y2)/2
						} else {
							y = y1 + (y2-y1)/2
						}

						break
					}
				}
			}

			gc.Stroke()

			gc.BeginPath()
			gc.ArcTo(x, y, stepCrackSmall, stepCrackSmall, 0, math.Pi*2)
			gc.Stroke()
		} else if v.Type == 111 {
			gc.BeginPath()

			gc.SetStrokeColor(red)
			gc.SetFillColor(red)

			gc.ArcTo(x, y, stepMiddle, stepMiddle, 0, math.Pi*2)

			gc.Fill()
		} else if v.Type == 112 {
			gc.BeginPath()

			gc.SetStrokeColor(blue)
			gc.SetFillColor(blue)

			gc.ArcTo(x, y, stepMiddle, stepMiddle, 0, math.Pi*2)

			gc.Fill()
		} else if v.Type == pipeLeakGreen {
			gc.BeginPath()

			gc.SetStrokeColor(green)
			gc.SetFillColor(green)

			gc.ArcTo(x, y, stepMiddle, stepMiddle, 0, math.Pi*2)

			gc.Fill()
		} else if v.Type == pipeLeakViolet {
			gc.BeginPath()

			gc.SetStrokeColor(violet)
			gc.SetFillColor(violet)

			gc.ArcTo(x, y, stepMiddle, stepMiddle, 0, math.Pi*2)

			gc.Fill()
		} else if v.Type == leakRed {
			gc.BeginPath()

			gc.SetStrokeColor(red)
			gc.SetFillColor(red)

			gc.ArcTo(x, y, stepMiddle, stepMiddle, 0, math.Pi*2)

			gc.Stroke()
		} else if v.Type == leakGreen {
			gc.BeginPath()

			gc.SetStrokeColor(green)
			gc.SetFillColor(green)

			gc.ArcTo(x, y, stepMiddle, stepMiddle, 0, math.Pi*2)

			gc.Stroke()
		}

	}

	targetFilename := fmt.Sprintf("%v/periodicresult/%v/%v.jpg", config.UploadPath, periodic, blueprint.Id)

	// 기존 이미지 파일 삭제 (캐시 방지)
	os.Remove(targetFilename)

	err := global.SaveToJpegFile(targetFilename, dest)
	if err != nil {
		log.Println("Failed to save image:", targetFilename, err)
	}
}
