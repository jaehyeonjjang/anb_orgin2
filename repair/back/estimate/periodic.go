package estimate

import (
	"encoding/json"
	"fmt"
	"log"
	"repair/global"
	"repair/global/config"
	"repair/global/time"
	"repair/models"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/xuri/excelize/v2"
)

// MultiyearPeriod 다회용 기간 구조체
type MultiyearPeriod struct {
	Year    int   `json:"year"`
	Periods []int `json:"periods"` // 1=상반기, 2=하반기
}

func Periodic(id int64, typeid int, conn *models.Connection, estimate *models.Estimate, compareestimates []models.Compareestimate, apt *models.Apt) string {
	// 계약 정보 조회
	contractManager := models.NewContractManager(conn)
	contract := contractManager.GetByEstimate(estimate.Id)

	// estimate.Subtype: 1=상반기, 2=하반기, 3=연간, 4=다회용
	periodicType := estimate.Subtype

	// 템플릿 파일 매핑: 1,2 → periodic1 / 3 → periodic2 / 4 → periodic3
	var templateNum int
	if periodicType == 1 || periodicType == 2 {
		templateNum = 1
	} else if periodicType == 3 {
		templateNum = 2
	} else {
		templateNum = 3
	}

	excelFilename := fmt.Sprintf("periodic%v.xlsx", templateNum)

	var complex int64 = 1
	if periodicType == 4 {
		complex = 2
	}

	switch typeid {
	case 1:
		if complex == 1 {
			excelFilename = "detail-compare1.xlsx"
		} else {
			excelFilename = "detail-compare1-3.xlsx"
		}
	case 2:
		excelFilename = "detail-compare2.xlsx"
	}

	f, err := excelize.OpenFile(fmt.Sprintf("./doc/estimate/%v", excelFilename))
	if err != nil {
		log.Println(err)
		return ""
	}

	sheet := "갑지"

	t := time.ParseDate(estimate.Writedate)

	str := strings.Split(apt.Completeyear, "-")
	complateyear := ""

	if len(str) == 3 {
		complateyear = fmt.Sprintf("%v년 %v월", str[0], str[1])
	} else if len(str) == 2 {
		complateyear = fmt.Sprintf("%v년 %v월", str[0], str[1])
	} else {
		str := strings.Split(apt.Completeyear, " ")

		if len(str) == 3 {
			complateyear = fmt.Sprintf("%v %v", str[0], str[1])
		} else if len(str) == 2 {
			complateyear = fmt.Sprintf("%v %v", str[0], str[1])
		} else {
			complateyear = apt.Completeyear
		}
	}

	buildingSize := ""

	if apt.Type == "아파트" || apt.Familycount3 > 0 {
		flatcount := strings.Split(apt.Flatcount, "(")
		if len(flatcount) == 2 {
			all := strings.TrimSpace(flatcount[0])
			target := strings.TrimSpace(strings.ReplaceAll(flatcount[1], ")", ""))

			buildingSize = fmt.Sprintf("아파트 %v개동 %v세대 (해당동 : %v개동)", all, apt.Familycount, target)
		} else {
			buildingSize = fmt.Sprintf("아파트 %v개동 %v세대", flatcount[0], apt.Familycount)
		}
	} else {
		floor := ""
		if apt.Undergroundfloor > 0 {
			floor = fmt.Sprintf("지하%v층~", apt.Undergroundfloor)
		}

		if apt.Groundfloor != 0 {
			apt.Floor = fmt.Sprintf("%v", apt.Groundfloor)
		}

		floor += fmt.Sprintf("지상%v층", apt.Floor)

		if apt.Area != "" {
			floor += fmt.Sprintf("(연면적 : %v)", apt.Area)
		}

		buildingSize = floor
	}

	tel := ""
	fax := ""

	if apt.Tel != "" {
		tel = fmt.Sprintf("전화 : %v", apt.Tel)
	}

	if apt.Fax != "" {
		fax = fmt.Sprintf("팩스 : %v", apt.Fax)
	}

	if fax != "" {
		tel += ",      " + fax
	}

	if estimate.Days == 0 {
		estimate.Days = 1
	}

	no := GetEstimateNo(typeid, t, estimate.Date, conn)

	part := ""
	switch estimate.Subtype {
	case 1:
		part = "상반기 "
	case 2:
		part = "하반기 "
	default:
		part = "연간 "
	}

	if typeid == 0 {
		// H10셀 텍스트 생성
		h10Text := ""

		// periodicType이 4(다회용)이고 estimate의 multiyear_periods 필드가 있으면 사용
		if periodicType == 4 && estimate.Multiyear_periods != "" {
			var periods []MultiyearPeriod
			if err := json.Unmarshal([]byte(estimate.Multiyear_periods), &periods); err == nil && len(periods) > 0 {
				// 선택된 기간들로 H10 텍스트 생성
				var periodTexts []string
				for _, p := range periods {
					// 상반기(1)와 하반기(2)가 모두 선택된 경우 연간으로 표시
					hasFirst := false
					hasSecond := false
					for _, period := range p.Periods {
						if period == 1 {
							hasFirst = true
						} else if period == 2 {
							hasSecond = true
						}
					}

					if hasFirst && hasSecond {
						periodTexts = append(periodTexts, fmt.Sprintf("%v년 연간", p.Year))
					} else {
						for _, period := range p.Periods {
							if period == 1 {
								periodTexts = append(periodTexts, fmt.Sprintf("%v년 상반기", p.Year))
							} else if period == 2 {
								periodTexts = append(periodTexts, fmt.Sprintf("%v년 하반기", p.Year))
							}
						}
					}
				}
				h10Text = strings.Join(periodTexts, "/")
			}
		} else if contract != nil {
			startdate := time.ParseDate(contract.Contractstartdate)
			enddate := time.ParseDate(contract.Contractenddate)

			if startdate != nil && enddate != nil {
				startYear := startdate.Year()
				endYear := enddate.Year()

				switch periodicType {
				case 1:
					// 상반기
					h10Text = fmt.Sprintf("%v년 상반기", startYear)
				case 2:
					// 하반기
					h10Text = fmt.Sprintf("%v년 하반기", startYear)
				case 3:
					// 연간 (1년)
					h10Text = fmt.Sprintf("%v년 연간", startYear)
				case 4:
					// 다회용 (여러 년도) - 각 년도별 상반기/하반기 판단
					type YearPeriod struct {
						Year       int
						FirstHalf  bool
						SecondHalf bool
					}

					var periods []YearPeriod

					for year := startYear; year <= endYear; year++ {
						var hasFirst, hasSecond bool

						startMonth := 1
						endMonth := 12

						if year == startYear {
							startMonth = int(startdate.Month())
						}
						if year == endYear {
							endMonth = int(enddate.Month())
						}

						// 상반기: 1-6월
						if startMonth <= 6 {
							hasFirst = true
						}
						// 하반기: 7-12월
						if endMonth >= 7 {
							hasSecond = true
						}
						// 년도가 상반기/하반기 경계를 넘나드는 경우
						if startMonth <= 6 && endMonth >= 7 {
							hasFirst = true
							hasSecond = true
						}

						periods = append(periods, YearPeriod{
							Year:       year,
							FirstHalf:  hasFirst,
							SecondHalf: hasSecond,
						})
					}

					// 모든 년도가 같은 패턴인지 확인
					allSame := true
					firstPattern := ""
					if periods[0].FirstHalf && periods[0].SecondHalf {
						firstPattern = "both"
					} else if periods[0].FirstHalf {
						firstPattern = "first"
					} else if periods[0].SecondHalf {
						firstPattern = "second"
					}

					for _, p := range periods {
						pattern := ""
						if p.FirstHalf && p.SecondHalf {
							pattern = "both"
						} else if p.FirstHalf {
							pattern = "first"
						} else if p.SecondHalf {
							pattern = "second"
						}
						if pattern != firstPattern {
							allSame = false
							break
						}
					}

					// 년도 문자열 생성
					yearList := ""
					for _, p := range periods {
						if yearList != "" {
							yearList += "/"
						}
						if allSame {
							// 모든 년도가 같은 패턴이면 년도만
							yearList += fmt.Sprintf("%v년", p.Year)
						} else {
							// 다른 패턴이면 각각 표시
							if p.FirstHalf && p.SecondHalf {
								yearList += fmt.Sprintf("%v년", p.Year)
							} else if p.FirstHalf {
								yearList += fmt.Sprintf("%v년 상반기", p.Year)
							} else if p.SecondHalf {
								yearList += fmt.Sprintf("%v년 하반기", p.Year)
							}
						}
					}

					// 최종 텍스트
					if allSame {
						if firstPattern == "both" {
							h10Text = fmt.Sprintf("%v 연간", yearList)
						} else if firstPattern == "first" {
							h10Text = fmt.Sprintf("%v 상반기", yearList)
						} else if firstPattern == "second" {
							h10Text = fmt.Sprintf("%v 하반기", yearList)
						}
					} else {
						h10Text = yearList
					}
				}
			}
		} else {
			// 계약 날짜가 없으면 기본값
			h10Text = fmt.Sprintf("%v년 %v", t.Year(), strings.TrimSpace(part))
		}

		if h10Text == "" {
			// 계약 정보도 없고 multiyear_periods도 없으면 기본값
			h10Text = fmt.Sprintf("%v년 %v", t.Year(), strings.TrimSpace(part))
		}

		f.SetCellStr(sheet, "H6", no)
		f.SetCellStr(sheet, "H7", t.Humandate())
		f.SetCellStr(sheet, "H8", fmt.Sprintf("%v 입주자대표회의", apt.Name))
		f.SetCellStr(sheet, "H10", fmt.Sprintf("%v 정기안전점검 견적 건", h10Text))
		f.SetCellStr(sheet, "N21", apt.Name)
		f.SetCellStr(sheet, "N23", apt.Address)
		f.SetCellStr(sheet, "N22", buildingSize)
		f.SetCellStr(sheet, "N24", complateyear)

		// N25: 전화 / 팩스 형식으로 출력
		telFaxStr := ""
		if apt.Tel != "" && apt.Fax != "" {
			telFaxStr = fmt.Sprintf("%v / %v", apt.Tel, apt.Fax)
		} else if apt.Tel != "" {
			telFaxStr = apt.Tel
		} else if apt.Fax != "" {
			telFaxStr = apt.Fax
		}
		f.SetCellStr(sheet, "N25", telFaxStr)
		f.SetCellStr(sheet, "N31", fmt.Sprintf("일금%v원정(₩%v)", global.HumanMoney(int64(estimate.Price)), humanize.Comma(int64(estimate.Price))))

		// I30~I34: 계약 기간의 모든 상반기/하반기 나열 (periodic3, periodic4만)
		if periodicType == 3 || periodicType == 4 {
			var periodList []string

			// periodicType 4(다회용)이면 estimate의 multiyear_periods 사용
			if periodicType == 4 && estimate.Multiyear_periods != "" {
				var periods []MultiyearPeriod
				if err := json.Unmarshal([]byte(estimate.Multiyear_periods), &periods); err == nil {
					for _, p := range periods {
						for _, period := range p.Periods {
							if period == 1 {
								periodList = append(periodList, fmt.Sprintf("%v년 상반기 :", p.Year))
							} else if period == 2 {
								periodList = append(periodList, fmt.Sprintf("%v년 하반기 :", p.Year))
							}
						}
					}
				}
			} else if contract != nil {
				// 그 외에는 계약 정보 사용
				startdate := time.ParseDate(contract.Contractstartdate)
				enddate := time.ParseDate(contract.Contractenddate)

				if startdate != nil && enddate != nil {
					startYear := startdate.Year()
					endYear := enddate.Year()

					for year := startYear; year <= endYear; year++ {
						startMonth := 1
						endMonth := 12

						if year == startYear {
							startMonth = int(startdate.Month())
						}
						if year == endYear {
							endMonth = int(enddate.Month())
						}

						// 상반기 포함 여부 (1-6월)
						if startMonth <= 6 {
							periodList = append(periodList, fmt.Sprintf("%v년 상반기 :", year))
						}

						// 하반기 포함 여부 (7-12월)
						if endMonth >= 7 {
							periodList = append(periodList, fmt.Sprintf("%v년 하반기 :", year))
						}
					}
				}
			}

			if len(periodList) > 0 {

				// H32부터 H36까지 최대 5개 출력, N열에는 대가산출 H29 참조, Z열에는 '-VAT별도'
				cellsH := []string{"H32", "H33", "H34", "H35", "H36"}
				cellsN := []string{"N32", "N33", "N34", "N35", "N36"}
				cellsY := []string{"Y32", "Y33", "Y34", "Y35", "Y36"}
				cellsZ := []string{"Z32", "Z33", "Z34", "Z35", "Z36"}
				for i := 0; i < len(cellsH) && i < len(periodList); i++ {
					f.SetCellStr(sheet, cellsH[i], periodList[i])
					f.SetCellFormula(sheet, cellsN[i], "대가산출!H29")
					f.SetCellStr(sheet, cellsY[i], "원")
					f.SetCellStr(sheet, cellsZ[i], "-VAT별도")
				}
				// 총액 행 추가
				totalRow := 32 + len(periodList)
				if totalRow <= 39 { // 안전 범위 체크
					totalCellH := fmt.Sprintf("H%d", totalRow)
					totalCellN := fmt.Sprintf("N%d", totalRow)
					totalCellY := fmt.Sprintf("Y%d", totalRow)
					totalCellZ := fmt.Sprintf("Z%d", totalRow)

					// periodicType에 따라 텍스트 구분: 3=연간(연n회), 4=연속(총n회)
					totalText := ""
					if periodicType == 3 {
						totalText = fmt.Sprintf("※총액(연%d회)", len(periodList))
					} else {
						totalText = fmt.Sprintf("※총액(총%d회)", len(periodList))
					}
					f.SetCellStr(sheet, totalCellH, totalText)

					// R열 총합 계산 (1회 금액 * 회수)
					totalAmount := estimate.Price * len(periodList)
					humanAmount := global.HumanMoney(totalAmount)

					f.SetCellStr(sheet, totalCellN, fmt.Sprintf("일금%v원정 (₩%v)", humanAmount, humanize.Comma(int64(totalAmount))))
					f.SetCellStr(sheet, totalCellY, "")
					f.SetCellStr(sheet, totalCellZ, "-VAT별도※")

					// 총액 행 스타일 (굵게 + 정렬)
					boldCenterStyle, _ := f.NewStyle(&excelize.Style{
						Font:      &excelize.Font{Bold: true, Family: "바탕체", Size: 12},
						Alignment: &excelize.Alignment{Horizontal: "center"},
					})
					boldRightStyle, _ := f.NewStyle(&excelize.Style{
						Font:      &excelize.Font{Bold: true, Family: "바탕체", Size: 12},
						Alignment: &excelize.Alignment{Horizontal: "right"},
					})
					f.SetCellStyle(sheet, totalCellH, totalCellH, boldCenterStyle)
					f.SetCellStyle(sheet, totalCellN, totalCellN, boldRightStyle)
					f.SetCellStyle(sheet, totalCellZ, totalCellZ, boldCenterStyle)
				}
			}
		}

		sheet = "대가산출"

		f.SetCellValue(sheet, "E9", estimate.Personprice7)
		f.SetCellValue(sheet, "E10", estimate.Personprice8)
		f.SetCellValue(sheet, "E11", estimate.Personprice9)
		f.SetCellValue(sheet, "E12", estimate.Personprice10)

		f.SetCellValue(sheet, "E13", estimate.Personprice2)
		f.SetCellValue(sheet, "E14", estimate.Personprice3)
		f.SetCellValue(sheet, "E15", estimate.Personprice4)
		f.SetCellValue(sheet, "E16", estimate.Personprice5)

		f.SetCellValue(sheet, "G9", estimate.Person7*estimate.Days)
		f.SetCellValue(sheet, "G10", estimate.Person8*estimate.Days)
		f.SetCellValue(sheet, "G11", estimate.Person9*estimate.Days)
		f.SetCellValue(sheet, "G12", estimate.Person10*estimate.Days)

		f.SetCellValue(sheet, "G13", estimate.Person2)
		f.SetCellValue(sheet, "G14", estimate.Person3)
		f.SetCellValue(sheet, "G15", estimate.Person4)
		f.SetCellValue(sheet, "G16", estimate.Person5)

		f.SetCellValue(sheet, "F21", 1)
		f.SetCellValue(sheet, "F22", 1)

		outPersons := (estimate.Person7 + estimate.Person8 + estimate.Person9 + estimate.Person10) * estimate.Days
		f.SetCellValue(sheet, "G21", outPersons)
		f.SetCellValue(sheet, "G22", outPersons)

		if estimate.Person7 > 0 {
			f.SetCellValue(sheet, "I9", fmt.Sprintf("%v인 * %v일", estimate.Person7, estimate.Days))
		}
		if estimate.Person8 > 0 {
			f.SetCellValue(sheet, "I10", fmt.Sprintf("%v인 * %v일", estimate.Person8, estimate.Days))
		}
		if estimate.Person9 > 0 {
			f.SetCellValue(sheet, "I11", fmt.Sprintf("%v인 * %v일", estimate.Person9, estimate.Days))
		}
		if estimate.Person10 > 0 {
			f.SetCellValue(sheet, "I12", fmt.Sprintf("%v인 * %v일", estimate.Person10, estimate.Days))
		}

		// if estimate.Person2 > 0 {
		// 	f.SetCellValue(sheet, "I13", fmt.Sprintf("%v인 * 1일", estimate.Person2))
		// }
		// if estimate.Person3 > 0 {
		// 	f.SetCellValue(sheet, "I14", fmt.Sprintf("%v인 * 1일", estimate.Person3))
		// }
		// if estimate.Person4 > 0 {
		// 	f.SetCellValue(sheet, "I15", fmt.Sprintf("%v인 * 1일", estimate.Person4))
		// }
		// if estimate.Person5 > 0 {
		// 	f.SetCellValue(sheet, "I16", fmt.Sprintf("%v인 * 1일", estimate.Person5))
		// }

		f.SetCellValue(sheet, "G17", estimate.Financialprice)
		f.SetCellValue(sheet, "G18", estimate.Techprice)

		f.SetCellValue(sheet, "I17", fmt.Sprintf("직접인건비 * %v%%", estimate.Financialprice))
		f.SetCellValue(sheet, "I18", fmt.Sprintf("(직접인건비 + 제경비) * %v%%", estimate.Techprice))

		f.SetCellValue(sheet, "E21", estimate.Travelprice)
		f.SetCellValue(sheet, "E22", estimate.Carprice)

		f.SetCellValue(sheet, "G23", estimate.Danger)
		f.SetCellValue(sheet, "I23", fmt.Sprintf("외업인건비의 %v%%", estimate.Danger))

		f.SetCellValue(sheet, "G24", estimate.Machine)
		f.SetCellValue(sheet, "I24", fmt.Sprintf("직접인건비의 %v%%", estimate.Machine))

		f.SetCellValue(sheet, "E25", estimate.Printprice)
		f.SetCellValue(sheet, "G25", estimate.Print)

		f.SetCellValue(sheet, "H28", estimate.Saleprice)

		//f.SetCellStr(sheet, "D3", fmt.Sprintf("일금 %v원정(₩%v) - VAT 별도", global.HumanMoney(int64(estimate.Price)), humanize.Comma(int64(estimate.Price))))

		// G32: 상반기/하반기 횟수 (periodic3, periodic4만)
		if periodicType == 3 || periodicType == 4 {
			periodCount := 0

			// periodicType 4(다회용)이면 estimate의 multiyear_periods 사용
			if periodicType == 4 && estimate.Multiyear_periods != "" {
				var periods []MultiyearPeriod
				if err := json.Unmarshal([]byte(estimate.Multiyear_periods), &periods); err == nil {
					for _, p := range periods {
						periodCount += len(p.Periods)
					}
				}
			} else if contract != nil {
				// 그 외에는 계약 정보 사용
				startdate := time.ParseDate(contract.Contractstartdate)
				enddate := time.ParseDate(contract.Contractenddate)

				if startdate != nil && enddate != nil {
					startYear := startdate.Year()
					endYear := enddate.Year()

					for year := startYear; year <= endYear; year++ {
						startMonth := 1
						endMonth := 12

						if year == startYear {
							startMonth = int(startdate.Month())
						}
						if year == endYear {
							endMonth = int(enddate.Month())
						}

						if startMonth <= 6 {
							periodCount++
						}
						if endMonth >= 7 {
							periodCount++
						}
					}
				}
			}

			if periodCount > 0 {
				f.SetCellValue(sheet, "G32", periodCount)
			}
		}

		sheet = "계약서"

		//priceStr := fmt.Sprintf("일금 %v원정(₩%v) - VAT 별도", global.HumanMoney(int64(estimate.Price)), humanize.Comma(int64(estimate.Price)))
		//priceStr2 := fmt.Sprintf("일금 %v원정(₩%v)", global.HumanMoney(int64(estimate.Price)), humanize.Comma(int64(estimate.Price)))
		//priceStr3 := fmt.Sprintf("일금 %v원정(₩%v) - VAT 별도", global.HumanMoney(int64(estimate.Price)), humanize.Comma(int64(estimate.Price)))

		f.SetCellStr(sheet, "B9", apt.Name)

		// B16: h10Text 내용
		if h10Text != "" {
			f.SetCellStr(sheet, "B16", h10Text)
		}

		/*switch periodicType {
		case 1:
			f.SetCellStr(sheet, "B16", fmt.Sprintf("%v년 상반기", t.Year()))
			subtitle = fmt.Sprintf("%v년 상반기", t.Year())
		case 2:
			f.SetCellStr(sheet, "B16", fmt.Sprintf("%v년 하반기", t.Year()))
			subtitle = fmt.Sprintf("%v년 하반기", t.Year())
		case 3:
			f.SetCellStr(sheet, "B16", fmt.Sprintf("%v년 연간", t.Year()))
			subtitle = fmt.Sprintf("%v년 연간", t.Year())
			priceStr = fmt.Sprintf("일금 %v원정(₩%v) - VAT 별도", global.HumanMoney(int64(estimate.Price)*2), humanize.Comma(int64(estimate.Price)*2))
			priceStr2 = fmt.Sprintf("일금 %v원정(₩%v)", global.HumanMoney(int64(estimate.Price)*2), humanize.Comma(int64(estimate.Price)*2))
		default:
			f.SetCellStr(sheet, "B16", fmt.Sprintf("%v년", t.Year()))
			subtitle = fmt.Sprintf("%v년", t.Year())
		}*/

		//f.SetCellStr(sheet, "G37", fmt.Sprintf("%v 정기안전점검 용역", apt.Name))
		//f.SetCellStr(sheet, "G38", priceStr)
		//f.SetCellStr(sheet, "G39", priceStr2)

		// 계약일자 셀 (파일별 위치 다름)
		hCell := "E52" // periodic3, periodic4 기본값
		if periodicType == 1 || periodicType == 2 {
			hCell = "E51" // periodic1, periodic2
		}

		if contract != nil {
			startdate := time.ParseDate(contract.Contractstartdate)
			enddate := time.ParseDate(contract.Contractenddate)
			contractDate := time.ParseDate(contract.Contractdate)

			// 계약 기간 (파일별 셀 위치 다름)
			periodCell := "C74" // periodic3, periodic4 기본값
			if periodicType == 1 || periodicType == 2 {
				periodCell = "C73" // periodic1, periodic2
			}

			if startdate != nil && enddate != nil {
				f.SetCellStr(sheet, periodCell, fmt.Sprintf("① 계약기간은 %04d년 %2d월 %2d일부터 %04d년 %2d월 %2d일로 종료한다.",
					startdate.Year(), startdate.Month(), startdate.Day(),
					enddate.Year(), enddate.Month(), enddate.Day()))
				f.SetCellStr(sheet, "G38", fmt.Sprintf("%04d .  %02d .  %02d .   ~   %04d .  %02d .  %02d . ", startdate.Year(), startdate.Month(), startdate.Day(), enddate.Year(), enddate.Month(), enddate.Day()))
			} else if enddate != nil {
				f.SetCellStr(sheet, "G38", fmt.Sprintf("%04d .     .     .   ~   %04d .  %02d .  %02d . ", t.Year(), enddate.Year(), enddate.Month(), enddate.Day()))
			} else {
				f.SetCellStr(sheet, "G38", fmt.Sprintf("%04d .     .     .   ~   %04d .     .     . ", t.Year(), t.Year()))
			}

			// J42: 년도 정보 (periodic1,2는 "202x년 상반기/하반기", periodic3,4는 h10Text 그대로)
			if h10Text != "" {
				f.SetCellStr(sheet, "G42", fmt.Sprintf("상기 금액은 %v 정기안전점검 용역대가임", h10Text))
			}

			if contractDate != nil {
				f.SetCellStr(sheet, hCell, contractDate.Humandate())
			} else {
				f.SetCellStr(sheet, hCell, fmt.Sprintf("%v년", t.Year()))
			}
		} else {
			f.SetCellStr(sheet, "G39", fmt.Sprintf("%04d .     .     .   ~   %04d .     .     . ", t.Year(), t.Year()))
			f.SetCellStr(sheet, hCell, fmt.Sprintf("%v년", t.Year()))
		}

		//f.SetCellStr(sheet, "G42", apt.Address)

		//f.SetCellStr(sheet, "G43", tel)
		f.SetCellStr(sheet, "I40", apt.Tel)
		f.SetCellStr(sheet, "Q40", apt.Fax)
		/*if periodicType == 4 {
			f.SetCellStr(sheet, "G43", fmt.Sprintf("상기 금액은 %v 정기안전점검 용역대가임. (연 2회)\n  - 1회 : %v", h10Text, priceStr3))
		} else {
			f.SetCellStr(sheet, "G43", fmt.Sprintf("상기 금액은 %v 정기안전점검 용역대가임.", h10Text))
		}*/

		// 공동주택 타입 셀 (파일별 위치 다름)
		buildingTypeCell := "C70" // periodic3, periodic4 기본값
		if periodicType == 1 || periodicType == 2 {
			buildingTypeCell = "C69" // periodic1, periodic2
		}

		if apt.Type == "아파트" || apt.Familycount3 > 0 {
			f.SetCellStr(sheet, buildingTypeCell, "③ 용      도 : 공동주택")
		} else {
			f.SetCellStr(sheet, buildingTypeCell, "③ 용      도 : 공동주택외 건축물")
		}
		//f.SetCellStr(sheet, "K72", buildingSize)

		// 계약서 시트 행 offset (periodic1, periodic2는 한 행씩 위로)
		rowOffset := 0
		if periodicType == 1 || periodicType == 2 {
			rowOffset = -1
		}

		f.SetCellValue(sheet, fmt.Sprintf("I%d", 139+rowOffset), estimate.Personprice7)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", 140+rowOffset), estimate.Personprice8)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", 141+rowOffset), estimate.Personprice9)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", 142+rowOffset), estimate.Personprice10)

		f.SetCellValue(sheet, fmt.Sprintf("I%d", 143+rowOffset), estimate.Personprice2)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", 144+rowOffset), estimate.Personprice3)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", 145+rowOffset), estimate.Personprice4)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", 146+rowOffset), estimate.Personprice5)

		f.SetCellValue(sheet, fmt.Sprintf("P%d", 139+rowOffset), estimate.Person7*estimate.Days)
		f.SetCellValue(sheet, fmt.Sprintf("P%d", 140+rowOffset), estimate.Person8*estimate.Days)
		f.SetCellValue(sheet, fmt.Sprintf("P%d", 141+rowOffset), estimate.Person9*estimate.Days)
		f.SetCellValue(sheet, fmt.Sprintf("P%d", 142+rowOffset), estimate.Person10*estimate.Days)

		f.SetCellValue(sheet, fmt.Sprintf("P%d", 143+rowOffset), estimate.Person2)
		f.SetCellValue(sheet, fmt.Sprintf("P%d", 144+rowOffset), estimate.Person3)
		f.SetCellValue(sheet, fmt.Sprintf("P%d", 145+rowOffset), estimate.Person4)
		f.SetCellValue(sheet, fmt.Sprintf("P%d", 146+rowOffset), estimate.Person5)

		f.SetCellValue(sheet, fmt.Sprintf("P%d", 147+rowOffset), estimate.Financialprice)
		f.SetCellValue(sheet, fmt.Sprintf("P%d", 148+rowOffset), estimate.Techprice)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", 151+rowOffset), estimate.Carprice)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", 152+rowOffset), estimate.Travelprice)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", 153+rowOffset), fmt.Sprintf("외업인건비의 %v%%", estimate.Danger))
		f.SetCellFormula(sheet, fmt.Sprintf("R%d", 153+rowOffset), fmt.Sprintf("=ROUND(SUM(R%d:R%d)*%v%%, 0)", 139+rowOffset, 142+rowOffset, estimate.Danger))
		f.SetCellValue(sheet, fmt.Sprintf("I%d", 154+rowOffset), fmt.Sprintf("직접인건비의 %v%%", estimate.Machine))
		f.SetCellFormula(sheet, fmt.Sprintf("R%d", 154+rowOffset), fmt.Sprintf("=ROUND(R%d*%v%%, 0)", 138+rowOffset, estimate.Machine))
		f.SetCellValue(sheet, fmt.Sprintf("R%d", 155+rowOffset), estimate.Printprice)

		f.SetCellValue(sheet, fmt.Sprintf("R%d", 158+rowOffset), estimate.Saleprice)
		f.SetCellValue(sheet, fmt.Sprintf("R%d", 159+rowOffset), estimate.Price)
		f.SetCellValue(sheet, fmt.Sprintf("R%d", 161+rowOffset), estimate.Price)

		f.UpdateLinkedValue()
	} else if typeid == 1 {

		// 비교 견적
		var compareestimate models.Compareestimate
		for _, v := range compareestimates {
			if v.Comparecompany == 1 {
				compareestimate = v
				break
			}
		}
		sheet = "엘림공문-비교견적"

		title := ""
		switch estimate.Subtype {
		case 1:
			title = fmt.Sprintf("%v %v년 상반기 정기안전점검", apt.Name, t.Year())
		case 2:
			title = fmt.Sprintf("%v %v년 하반기 정기안전점검", apt.Name, t.Year())
		case 3:
			title = fmt.Sprintf("%v %v년 연간 정기안전점검", apt.Name, t.Year())
		default:
			title = fmt.Sprintf("%v %v년 정기안전점검", apt.Name, t.Year())
		}

		priceStr := fmt.Sprintf("일금 %v원정(₩%v) - VAT 별도", global.HumanMoney(int64(compareestimate.Price)*complex), humanize.Comma(int64(compareestimate.Price)*complex))
		price1Str := fmt.Sprintf("일금 %v원정", global.HumanMoney(int64(compareestimate.Price)*complex))
		price2Str := fmt.Sprintf("(₩%v) - VAT 별도", humanize.Comma(int64(compareestimate.Price)*complex))

		t = time.ParseDate(compareestimate.Writedate)
		f.SetCellStr(sheet, "C8", fmt.Sprintf("ELIM-%04d%02d%02d", t.Year(), t.Month(), t.Day()))
		f.SetCellStr(sheet, "C9", t.Humandate())
		f.SetCellStr(sheet, "C11", apt.Name)
		f.SetCellStr(sheet, "C18", title)
		f.SetCellStr(sheet, "C29", title)

		f.SetCellStr(sheet, "C30", "- 육안조사 통한 현장조사\n- 시특법에 의한 정기점검 실시 및 보고서 작성\n- FMS등록 제출업무")

		sheet = "엘림대가내역서 갑지"

		f.SetCellStr(sheet, "B7", title)
		f.SetCellStr(sheet, "B10", buildingSize)
		f.SetCellStr(sheet, "B13", price1Str)
		f.SetCellStr(sheet, "D13", price2Str)
		f.SetCellValue(sheet, "F18", compareestimate.Price)

		if complex == 1 {
			sheet = "엘림대가내역서"
		} else {
			sheet = "엘림 상하반기 대가내역서"

			f.SetCellStr(sheet, "B3", fmt.Sprintf("%v %v년 건축물 정기안전점검(상반기) 용역", apt.Name, t.Year()))
			f.SetCellStr(sheet, "B31", fmt.Sprintf("%v %v년 건축물 정기안전점검(하반기) 용역", apt.Name, t.Year()))
		}

		priceStr = fmt.Sprintf("일금 %v원정(₩%v) - VAT 별도", global.HumanMoney(int64(compareestimate.Price)), humanize.Comma(int64(compareestimate.Price)))
		f.SetCellStr(sheet, "B5", priceStr)

		f.SetCellValue(sheet, "F9", compareestimate.Personprice7)
		f.SetCellValue(sheet, "F10", compareestimate.Personprice8)
		f.SetCellValue(sheet, "F11", compareestimate.Personprice9)
		f.SetCellValue(sheet, "F12", compareestimate.Personprice10)

		f.SetCellValue(sheet, "F13", compareestimate.Personprice2)
		f.SetCellValue(sheet, "F14", compareestimate.Personprice3)
		f.SetCellValue(sheet, "F15", compareestimate.Personprice4)
		f.SetCellValue(sheet, "F16", compareestimate.Personprice5)

		if compareestimate.Person7 > 0 {
			f.SetCellValue(sheet, "D9", compareestimate.Person7)
			f.SetCellValue(sheet, "E9", estimate.Days)
		}
		if compareestimate.Person8 > 0 {
			f.SetCellValue(sheet, "D10", compareestimate.Person8)
			f.SetCellValue(sheet, "E10", estimate.Days)
		}
		if compareestimate.Person9 > 0 {
			f.SetCellValue(sheet, "D11", compareestimate.Person9)
			f.SetCellValue(sheet, "E11", estimate.Days)
		}
		if compareestimate.Person10 > 0 {
			f.SetCellValue(sheet, "D12", compareestimate.Person10)
			f.SetCellValue(sheet, "E12", estimate.Days)
		}

		if compareestimate.Person2 > 0 {
			f.SetCellValue(sheet, "E13", compareestimate.Person2)
			f.SetCellValue(sheet, "D13", 1)
		}
		if compareestimate.Person3 > 0 {
			f.SetCellValue(sheet, "E14", compareestimate.Person3)
			f.SetCellValue(sheet, "D14", 1)
		}
		if compareestimate.Person4 > 0 {
			f.SetCellValue(sheet, "E15", compareestimate.Person4)
			f.SetCellValue(sheet, "D15", 1)
		}
		if compareestimate.Person5 > 0 {
			f.SetCellValue(sheet, "E16", compareestimate.Person5)
			f.SetCellValue(sheet, "D16", 1)
		}

		f.SetCellValue(sheet, "E20", estimate.Days)
		f.SetCellValue(sheet, "E21", estimate.Days)

		f.SetCellValue(sheet, "C18", compareestimate.Financialprice)
		f.SetCellValue(sheet, "C19", compareestimate.Techprice)
		f.SetCellValue(sheet, "F20", compareestimate.Carprice)
		f.SetCellValue(sheet, "F21", compareestimate.Travelprice)
		f.SetCellValue(sheet, "C22", compareestimate.Danger)
		f.SetCellValue(sheet, "C23", compareestimate.Machine)
		f.SetCellValue(sheet, "G24", compareestimate.Printprice)

		f.SetCellValue(sheet, "G28", compareestimate.Price)

		f.UpdateLinkedValue()
	} else if typeid == 2 {
		// 비교 견적
		var compareestimate models.Compareestimate
		for _, v := range compareestimates {
			if v.Comparecompany == 2 {
				compareestimate = v
				break
			}
		}
		sheet = "갑지"

		title := ""

		switch estimate.Subtype {
		case 1:
			title = fmt.Sprintf("%v년 상반기 정기안전점검", t.Year())
		case 2:
			title = fmt.Sprintf("%v년 하반기 정기안전점검", t.Year())
		case 3:
			title = fmt.Sprintf("%v년 연간 정기안전점검", t.Year())
		default:
			title = fmt.Sprintf("%v년 정기안전점검", t.Year())
		}

		t = time.ParseDate(compareestimate.Writedate)
		f.SetCellStr(sheet, "A6", apt.Name)

		if estimate.Subtype == 3 {
			title = fmt.Sprintf("%v년 상반기 정기안전점검", t.Year())
		}

		f.SetCellStr(sheet, "A9", title)
		f.SetCellStr(sheet, "D24", buildingSize)
		f.SetCellStr(sheet, "I10", t.Humandate())

		sheet = "을지"

		f.SetCellValue(sheet, "G6", compareestimate.Personprice7)
		f.SetCellValue(sheet, "G7", compareestimate.Personprice8)
		f.SetCellValue(sheet, "G8", compareestimate.Personprice9)
		f.SetCellValue(sheet, "G9", compareestimate.Personprice10)

		if compareestimate.Person7 > 0 {
			f.SetCellValue(sheet, "C6", compareestimate.Person7*estimate.Days)
		}
		if compareestimate.Person8 > 0 {
			f.SetCellValue(sheet, "C7", compareestimate.Person8*estimate.Days)
		}
		if compareestimate.Person9 > 0 {
			f.SetCellValue(sheet, "C8", compareestimate.Person9*estimate.Days)
		}
		if compareestimate.Person10 > 0 {
			f.SetCellValue(sheet, "C9", compareestimate.Person10*estimate.Days)
		}

		if compareestimate.Person2 > 0 {
			f.SetCellValue(sheet, "E6", compareestimate.Person2)
		}
		if compareestimate.Person3 > 0 {
			f.SetCellValue(sheet, "E7", compareestimate.Person3)
		}
		if compareestimate.Person4 > 0 {
			f.SetCellValue(sheet, "E8", compareestimate.Person4)
		}
		if compareestimate.Person5 > 0 {
			f.SetCellValue(sheet, "E9", compareestimate.Person5)
		}

		f.SetCellValue(sheet, "G11", compareestimate.Travelprice)
		f.SetCellValue(sheet, "G13", compareestimate.Carprice)
		f.SetCellValue(sheet, "F15", compareestimate.Danger)
		f.SetCellValue(sheet, "F16", compareestimate.Machine)
		f.SetCellValue(sheet, "G17", compareestimate.Printprice)

		f.SetCellValue(sheet, "B19", compareestimate.Financialprice)
		f.SetCellValue(sheet, "B20", compareestimate.Techprice)

		f.SetCellValue(sheet, "F27", compareestimate.Saleprice)
		f.SetCellValue(sheet, "H27", compareestimate.Price)

		f.UpdateLinkedValue()
	}

	filename := fmt.Sprintf("%v.xlsx", global.UniqueId())
	fullFilename := fmt.Sprintf("%v/%v", config.UploadPath, filename)
	e := f.SaveAs(fullFilename)
	if e != nil {
		log.Println(e)
	}
	f.Close()

	if estimate.Subtype == 3 && typeid == 2 {
		f, err = excelize.OpenFile(fmt.Sprintf("./doc/estimate/%v", excelFilename))
		if err != nil {
			log.Println(err)
			return ""
		}

		var compareestimate models.Compareestimate
		for _, v := range compareestimates {
			if v.Comparecompany == 2 {
				compareestimate = v
				break
			}
		}
		sheet = "갑지"

		title := ""

		switch estimate.Subtype {
		case 1:
			title = fmt.Sprintf("%v년 상반기 정기안전점검", t.Year())
		case 2:
			title = fmt.Sprintf("%v년 하반기 정기안전점검", t.Year())
		case 3:
			title = fmt.Sprintf("%v년 연간 정기안전점검", t.Year())
		default:
			title = fmt.Sprintf("%v년 정기안전점검", t.Year())
		}

		t = time.ParseDate(compareestimate.Writedate)
		f.SetCellStr(sheet, "A6", apt.Name)

		if estimate.Subtype == 3 {
			title = fmt.Sprintf("%v년 하반기 정기안전점검", t.Year())
		}

		f.SetCellStr(sheet, "A9", title)
		f.SetCellStr(sheet, "D24", buildingSize)
		f.SetCellStr(sheet, "I10", t.Humandate())

		sheet = "을지"

		f.SetCellValue(sheet, "G6", compareestimate.Personprice7)
		f.SetCellValue(sheet, "G7", compareestimate.Personprice8)
		f.SetCellValue(sheet, "G8", compareestimate.Personprice9)
		f.SetCellValue(sheet, "G9", compareestimate.Personprice10)

		if estimate.Person7 > 0 {
			f.SetCellValue(sheet, "C6", compareestimate.Person7*estimate.Days)
		}
		if estimate.Person8 > 0 {
			f.SetCellValue(sheet, "C7", compareestimate.Person8*estimate.Days)
		}
		if estimate.Person9 > 0 {
			f.SetCellValue(sheet, "C8", compareestimate.Person9*estimate.Days)
		}
		if estimate.Person10 > 0 {
			f.SetCellValue(sheet, "C9", compareestimate.Person10*estimate.Days)
		}

		if estimate.Person2 > 0 {
			f.SetCellValue(sheet, "E6", compareestimate.Person2)
		}
		if estimate.Person3 > 0 {
			f.SetCellValue(sheet, "E7", compareestimate.Person3)
		}
		if estimate.Person4 > 0 {
			f.SetCellValue(sheet, "E8", compareestimate.Person4)
		}
		if estimate.Person5 > 0 {
			f.SetCellValue(sheet, "E9", compareestimate.Person5)
		}

		f.SetCellValue(sheet, "G11", compareestimate.Travelprice)
		f.SetCellValue(sheet, "G13", compareestimate.Carprice)
		f.SetCellValue(sheet, "F15", compareestimate.Danger)
		f.SetCellValue(sheet, "F16", compareestimate.Machine)
		f.SetCellValue(sheet, "G17", compareestimate.Printprice)

		f.SetCellValue(sheet, "B19", compareestimate.Financialprice)
		f.SetCellValue(sheet, "B20", compareestimate.Techprice)

		f.SetCellValue(sheet, "F27", compareestimate.Saleprice)
		f.SetCellValue(sheet, "H27", compareestimate.Price)

		f.UpdateLinkedValue()

		filename = fmt.Sprintf("%v.xlsx", global.UniqueId())
		fullFilename2 := fmt.Sprintf("%v/%v", config.UploadPath, filename)
		e := f.SaveAs(fullFilename2)
		if e != nil {
			log.Println(e)
		}
		f.Close()

		fname := []string{"상반기 정기안전점검.xlsx", "하반기 정기안전점검.xlsx"}
		files := []string{fullFilename, fullFilename2}
		filename = fmt.Sprintf("%v.zip", global.UniqueId())
		fullFilename = fmt.Sprintf("%v/%v", config.UploadPath, filename)
		global.MakeZipfile(fullFilename, fname, files)
	}

	return filename
}
