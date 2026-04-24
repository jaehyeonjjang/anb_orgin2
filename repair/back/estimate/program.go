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

func Program(id int64, typeid int, conn *models.Connection, estimate *models.Estimate, compareestimates []models.Compareestimate, apt *models.Apt) string {
	// 계약 정보 조회하여 파일명 결정
	contractManager := models.NewContractManager(conn)
	contract := contractManager.GetByEstimate(estimate.Id)

	// estimate.Subtype: 1=상반기, 2=하반기, 3=연간, 4=연속, 5=무상1회
	periodicType := estimate.Subtype

	// 템플릿 파일 매핑: 1,2 → program-periodic1 / 3,5 → program-periodic2 / 4 → program-periodic3
	var templateNum int
	if periodicType == 1 || periodicType == 2 {
		templateNum = 1
	} else if periodicType == 3 || periodicType == 5 {
		templateNum = 2 // 연간, 무상1회
	} else {
		templateNum = 3 // 4 (연속)
	}

	excelFilename := fmt.Sprintf("program-periodic%v.xlsx", templateNum)

	f, err := excelize.OpenFile(fmt.Sprintf("./doc/estimate/%v", excelFilename))
	if err != nil {
		log.Println(err)
		return ""
	}

	sheet := "갑지"

	t := time.ParseDate(estimate.Writedate)

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

	no := GetEstimateNo(typeid, t, estimate.Date, conn)

	// buildingType := "공동주택"
	// if apt.Type == "아파트" || apt.Familycount3 > 0 {
	// 	buildingType = "공동주택"
	// } else {
	// 	buildingType = "공동주택외 건축물"
	// }

	part := ""

	switch estimate.Subtype {
	case 1:
		part = "상반기 "
	case 2:
		part = "하반기 "
	default:
		part = "연간 "
	}

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

	// J10셀 계약 정보 기반 텍스트
	j10Text := ""

	// periodicType이 4(다회용)이고 estimate의 multiyear_periods 필드가 있으면 사용
	if periodicType == 4 && estimate.Multiyear_periods != "" {
		var periods []MultiyearPeriod
		if err := json.Unmarshal([]byte(estimate.Multiyear_periods), &periods); err == nil && len(periods) > 0 {
			// 선택된 기간들로 J10 텍스트 생성
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
			j10Text = strings.Join(periodTexts, "/")
		}
	} else if contract != nil {
		startdate := time.ParseDate(contract.Contractstartdate)
		enddate := time.ParseDate(contract.Contractenddate)

		if startdate != nil && enddate != nil {
			startYear := startdate.Year()
			endYear := enddate.Year()

			switch periodicType {
			case 1:
				// 1회만 (상반기 또는 하반기)
				j10Text = fmt.Sprintf("%v년 %v", startYear, strings.TrimSpace(part))
			case 2:
				// 연간 (상반기+하반기, 1년)
				j10Text = fmt.Sprintf("%v년 연간", startYear)
			case 3:
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
						j10Text = fmt.Sprintf("%v 연간", yearList)
					} else if firstPattern == "first" {
						j10Text = fmt.Sprintf("%v 상반기", yearList)
					} else if firstPattern == "second" {
						j10Text = fmt.Sprintf("%v 하반기", yearList)
					}
				} else {
					j10Text = yearList
				}
			}
		} else {
			// 계약 날짜가 없으면 기본값
			j10Text = fmt.Sprintf("%v년 %v", t.Year(), strings.TrimSpace(part))
		}
	} else {
		// 계약 정보가 없으면 기본값
		j10Text = fmt.Sprintf("%v년 %v", t.Year(), strings.TrimSpace(part))
	}
	f.SetCellStr(sheet, "J6", no)
	f.SetCellStr(sheet, "J7", t.Humandate())
	f.SetCellStr(sheet, "J8", fmt.Sprintf("%v 입주자대표회장님", apt.Name))
	f.SetCellStr(sheet, "J10", j10Text)
	//f.SetCellStr(sheet, "T10", "점검 프로그램 사용 견적건")
	f.SetCellStr(sheet, "R20", apt.Name)
	f.SetCellStr(sheet, "R21", apt.Address)
	f.SetCellStr(sheet, "R22", buildingSize)
	f.SetCellStr(sheet, "R23", complateyear)

	// R24: 전화 / 팩스 형식으로 출력
	telFaxStr := ""
	if apt.Tel != "" && apt.Fax != "" {
		telFaxStr = fmt.Sprintf("%v / %v", apt.Tel, apt.Fax)
	} else if apt.Tel != "" {
		telFaxStr = apt.Tel
	} else if apt.Fax != "" {
		telFaxStr = apt.Fax
	}
	f.SetCellStr(sheet, "R24", telFaxStr)

	// K40: 년도 + 회수 정보
	k40Text := ""
	periodCount := 0

	// periodicType 4(다회용)이면 estimate의 multiyear_periods 사용
	if periodicType == 4 && estimate.Multiyear_periods != "" {
		var periods []MultiyearPeriod
		if err := json.Unmarshal([]byte(estimate.Multiyear_periods), &periods); err == nil {
			for _, p := range periods {
				periodCount += len(p.Periods)
			}
			if j10Text != "" {
				k40Text = fmt.Sprintf("나. 점검프로그램 사용 기간 - %v (%d회)", j10Text, periodCount)
			}
		}
	} else if contract != nil {
		startdate := time.ParseDate(contract.Contractstartdate)
		enddate := time.ParseDate(contract.Contractenddate)

		if startdate != nil && enddate != nil {
			// periodList 계산하여 회수 파악
			startYear := startdate.Year()
			endYear := enddate.Year()

			var tempPeriodList []string
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
					tempPeriodList = append(tempPeriodList, fmt.Sprintf("%v년 상반기", year))
				}
				if endMonth >= 7 {
					tempPeriodList = append(tempPeriodList, fmt.Sprintf("%v년 하반기", year))
				}
			}
			periodCount = len(tempPeriodList)

			// j10Text 기반으로 K 텍스트 생성
			k40Text = fmt.Sprintf("나.점검프로그램 사용 기간 - %v (%d회)", j10Text, periodCount)
		}
	}

	if k40Text != "" {
		// 파일별 K셀 위치 다름 (templateNum 기준)
		kCell := "K39" // program-periodic3 기본값
		if templateNum == 1 {
			kCell = "K33" // program-periodic1
		} else if templateNum == 2 {
			kCell = "K36" // program-periodic2
		}
		f.SetCellStr(sheet, kCell, k40Text)
	}

	// L30~L34: 계약 기간의 모든 상반기/하반기 나열 (periodic2, periodic3만)
	if periodicType != 1 {
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

			// L30부터 L34까지 최대 5개 출력, M열에는 대가산출 H25 참조, AC열에는 '-VAT별도'
			cellsL := []string{"L30", "L31", "L32", "L33", "L34"}
			cellsR := []string{"R30", "R31", "R32", "R33", "R34"}
			cellsAB := []string{"AB30", "AB31", "AB32", "AB33", "AB34"}
			cellsAC := []string{"AC30", "AC31", "AC32", "AC33", "AC34"}
			for i := 0; i < len(cellsL) && i < len(periodList); i++ {
				f.SetCellStr(sheet, cellsL[i], periodList[i])
				f.SetCellFormula(sheet, cellsR[i], "대가산출!H25")
				f.SetCellStr(sheet, cellsAB[i], "원")
				f.SetCellStr(sheet, cellsAC[i], "-VAT별도")
			}
			// 총액 행 추가
			totalRow := 30 + len(periodList)
			if totalRow <= 39 { // 안전 범위 체크
				totalCellL := fmt.Sprintf("L%d", totalRow)
				totalCellR := fmt.Sprintf("R%d", totalRow)
				totalCellAB := fmt.Sprintf("AB%d", totalRow)
				totalCellAC := fmt.Sprintf("AC%d", totalRow)

				// 총액 텍스트 설정 (연간/무상1회는 "연n회", 연속은 "총n회")
				var totalText string
				if periodicType == 3 || periodicType == 5 {
					totalText = fmt.Sprintf("※총액(연%d회)", len(periodList))
				} else {
					totalText = fmt.Sprintf("※총액(총%d회)", len(periodList))
				}
				f.SetCellStr(sheet, totalCellL, totalText)

				// R열 총합 계산 (1회 금액 * 회수)
				totalAmount := estimate.Price * len(periodList)
				humanAmount := global.HumanMoney(totalAmount)

				f.SetCellStr(sheet, totalCellR, fmt.Sprintf("일금 %v원정 (₩%v)", humanAmount, humanize.Comma(int64(totalAmount))))
				f.SetCellStr(sheet, totalCellAB, "")
				f.SetCellStr(sheet, totalCellAC, "-VAT별도※")

				// 총액 행 스타일 (굵게 + 정렬)
				boldCenterStyle, _ := f.NewStyle(&excelize.Style{
					Font:      &excelize.Font{Bold: true, Family: "바탕체", Size: 12},
					Alignment: &excelize.Alignment{Horizontal: "center"},
				})
				boldRightStyle, _ := f.NewStyle(&excelize.Style{
					Font:      &excelize.Font{Bold: true, Family: "바탕체", Size: 12},
					Alignment: &excelize.Alignment{Horizontal: "right"},
				})
				f.SetCellStyle(sheet, totalCellL, totalCellL, boldCenterStyle)
				f.SetCellStyle(sheet, totalCellR, totalCellR, boldRightStyle)
				f.SetCellStyle(sheet, totalCellAC, totalCellAC, boldCenterStyle)
			}
		}
	}

	/*
		switch estimate.Subtype {
		case 1, 2:
			f.SetCellStr(sheet, "M33", subtitle)
		case 3:
			f.SetCellStr(sheet, "M33", subtitle2)
		}
	*/

	sheet = "대가산출"

	/*f.SetCellValue(sheet, "E9", estimate.Personprice2)
	f.SetCellValue(sheet, "E10", estimate.Personprice3)
	f.SetCellValue(sheet, "E11", estimate.Personprice4)
	f.SetCellValue(sheet, "E12", estimate.Personprice5)

	f.SetCellValue(sheet, "G9", estimate.Person2)
	f.SetCellValue(sheet, "G10", estimate.Person3)
	f.SetCellValue(sheet, "G11", estimate.Person4)
	f.SetCellValue(sheet, "G12", estimate.Person5)

	f.SetCellValue(sheet, "F17", 1)
	f.SetCellValue(sheet, "F18", 1)

	outPersons := (estimate.Person7 + estimate.Person8 + estimate.Person9 + estimate.Person10) * estimate.Days
	f.SetCellValue(sheet, "G17", outPersons)
	f.SetCellValue(sheet, "G18", outPersons)

	f.SetCellValue(sheet, "G13", estimate.Financialprice)
	f.SetCellValue(sheet, "G14", estimate.Techprice)

	f.SetCellValue(sheet, "I13", fmt.Sprintf("직접인건비 * %v%%", estimate.Financialprice))
	f.SetCellValue(sheet, "I14", fmt.Sprintf("(직접인건비 + 제경비) * %v%%", estimate.Techprice))

	f.SetCellValue(sheet, "E17", estimate.Travelprice)
	f.SetCellValue(sheet, "E18", estimate.Carprice)

	f.SetCellValue(sheet, "G19", estimate.Danger)
	f.SetCellValue(sheet, "I19", fmt.Sprintf("외업인건비의 %v%%", estimate.Danger))

	f.SetCellValue(sheet, "G20", estimate.Machine)
	f.SetCellValue(sheet, "I20", fmt.Sprintf("직접인건비의 %v%%", estimate.Machine))

	f.SetCellValue(sheet, "E21", estimate.Printprice)
	f.SetCellValue(sheet, "G21", estimate.Print)

	f.SetCellValue(sheet, "H24", estimate.Saleprice)*/

	// H25: 웹사이트에서 작성한 견적금액
	f.SetCellValue(sheet, "H25", estimate.Price)

	// G28: 상반기/하반기 횟수 (periodic2, periodic3만)
	if periodicType != 1 {
		periodCountG28 := 0

		// periodicType 4(다회용)이면 estimate의 multiyear_periods 사용
		if periodicType == 4 && estimate.Multiyear_periods != "" {
			var periods []MultiyearPeriod
			if err := json.Unmarshal([]byte(estimate.Multiyear_periods), &periods); err == nil {
				for _, p := range periods {
					periodCountG28 += len(p.Periods)
				}
			}
		} else if contract != nil {
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
						periodCountG28++
					}
					if endMonth >= 7 {
						periodCountG28++
					}
				}
			}
		}

		if periodCountG28 > 0 {
			f.SetCellValue(sheet, "G28", periodCountG28)
		}
	}

	sheet = "계약서"

	f.SetCellStr(sheet, "L40", apt.Tel)
	f.SetCellStr(sheet, "T40", apt.Fax)

	// 계약일자 셀 위치 (파일별 위치 다름)
	hCell := "H52" // periodic2, periodic3 기본값
	if periodicType == 1 {
		hCell = "H51" // periodic1
	}

	// 계약 정보 입력
	if contract != nil {
		startdate := time.ParseDate(contract.Contractstartdate)
		enddate := time.ParseDate(contract.Contractenddate)
		contractDate := time.ParseDate(contract.Contractdate)

		// 프로그램 사용 기간 (파일별 셀 위치 다름)
		periodCell := "F79" // periodic2, periodic3 기본값
		if periodicType == 1 {
			periodCell = "F78" // periodic1
		}

		if startdate != nil && enddate != nil {
			f.SetCellStr(sheet, periodCell, fmt.Sprintf("① 프로그램 사용은 %04d년 %2d월 %2d일부터 %04d년 %2d월 %2d일로 종료한다.",
				startdate.Year(), startdate.Month(), startdate.Day(),
				enddate.Year(), enddate.Month(), enddate.Day()))
		}

		// M42: 년도 정보 (회수 없이, periodic2, periodic3만)
		if j10Text != "" && periodicType != 1 {
			f.SetCellStr(sheet, "M42", j10Text)
		}

		if startdate != nil && enddate != nil {
			f.SetCellStr(sheet, "J38", fmt.Sprintf("%04d .  %02d .  %02d .   ~   %04d .  %02d .  %02d . ", startdate.Year(), startdate.Month(), startdate.Day(), enddate.Year(), enddate.Month(), enddate.Day()))
		} else if enddate != nil {
			f.SetCellStr(sheet, "J38", fmt.Sprintf("%04d .     .     .   ~   %04d .  %02d .  %02d . ", t.Year(), enddate.Year(), enddate.Month(), enddate.Day()))
		} else {
			f.SetCellStr(sheet, "J38", fmt.Sprintf("%04d .     .     .   ~   %04d .     .     . ", t.Year(), t.Year()))
		}

		// 계약일자
		if contractDate != nil {
			f.SetCellStr(sheet, hCell, contractDate.Humandate())
		} else {
			f.SetCellStr(sheet, hCell, fmt.Sprintf("%v년", t.Year()))
		}
	} else {
		f.SetCellStr(sheet, "J38", fmt.Sprintf("%04d .     .     .   ~   %04d .     .     . ", t.Year(), t.Year()))
		f.SetCellStr(sheet, hCell, fmt.Sprintf("%v년", t.Year()))
	}

	f.UpdateLinkedValue()

	filename := fmt.Sprintf("%v.xlsx", global.UniqueId())
	fullFilename := fmt.Sprintf("%v/%v", config.UploadPath, filename)
	e := f.SaveAs(fullFilename)
	if e != nil {
		log.Println(e)
	}
	f.Close()

	return filename
}
