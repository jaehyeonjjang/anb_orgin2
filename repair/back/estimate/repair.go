package estimate

import (
	"fmt"
	"log"
	"os"
	"repair/global"
	"repair/global/config"
	"repair/global/time"
	"repair/models"
	"strings"

	"github.com/xuri/excelize/v2"
)

func Repair(id int64, typeid int, conn *models.Connection, estimate *models.Estimate, compareestimates []models.Compareestimate, apt *models.Apt) string {
	excelFilename := ""
	type1 := ""
	type2 := ""
	switch estimate.Subtype {
	case 1:
		excelFilename = "repair1.xlsx"
		type1 = "조정"
		type2 = "조정"
	case 2:
		excelFilename = "repair2.xlsx"
		type1 = "수립"
		type2 = "수립(조정 포함)"
	}

	log.Println("typeid", typeid)
	switch typeid {
	case 3:
		excelFilename = "repair-compare.xlsx"
	case 4:
		excelFilename = "repair-compare2.xlsx"
	}

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)

	f, err := excelize.OpenFile(fmt.Sprintf("./doc/estimate/%v", excelFilename))
	if err != nil {
		log.Println(err)
		return ""
	}
	defer f.Close()

	flatcount := strings.Split(apt.Flatcount, "(")

	sheet := "갑지"

	if typeid == 0 {
		t := time.ParseDate(estimate.Writedate)

		// estimateManager := models.NewEstimateManager(conn)
		// duration := time.Now().GetDurationArray()
		// count := estimateManager.Count([]interface{}{
		// 	// models.Where{Column: "type", Value: 1, Compare: "="},
		// 	models.Where{Column: "date", Value: duration[0], Compare: ">="},
		// 	models.Where{Column: "date", Value: estimate.Date, Compare: "<"},
		// })

		// no := fmt.Sprintf("ANB-%v-%v", t.DateAsOnlyNumber(), count+1)
		no := GetEstimateNo(typeid, t, estimate.Date, conn)

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

		str := strings.Split(estimate.Writedate, "-")
		f.SetCellStr(sheet, "K6", no)
		f.SetCellStr(sheet, "K7", fmt.Sprintf("%v년 %v월 %v일", str[0], str[1], str[2]))
		f.SetCellStr(sheet, "K8", fmt.Sprintf("%v 입주자대표회장님", apt.Name))

		if estimate.Event == 1 {
			f.SetCellStr(sheet, "K10", fmt.Sprintf("%v 장기수선계획서 %v 견적건 (이벤트 견적)", apt.Name, type2))
		} else {
			f.SetCellStr(sheet, "K10", fmt.Sprintf("%v 장기수선계획서 %v 견적건", apt.Name, type2))
		}

		f.SetCellStr(sheet, "G21", fmt.Sprintf("에서 %v하고자 하는 장기수선계획서 작성에 관한 견적서를 아래와 같이 제출하오니 검토", type1))
		f.SetCellStr(sheet, "G22", fmt.Sprintf("하시어 %v 대행업무를 위임하여 주시기 바랍니다.", type1))

		f.SetCellStr(sheet, "G43", fmt.Sprintf("※ 첨    부 :  1. 장기수선계획서 %v 견적서 1부 끝.", type1))

		f.SetCellStr(sheet, "R25", apt.Name)

		complateyear := ""
		str = strings.Split(apt.Completeyear, "-")

		if len(str) == 3 {
			complateyear = fmt.Sprintf("%v년 %v월 %v일", str[0], str[1], str[2])
		} else if len(str) == 2 {
			complateyear = fmt.Sprintf("%v년 %v월", str[0], str[1])
		} else {
			complateyear = apt.Completeyear
		}

		f.SetCellStr(sheet, "R26", apt.Address)
		f.SetCellStr(sheet, "R27", fmt.Sprintf("아파트 %v개동 %v세대", flatcount[0], apt.Familycount))
		f.SetCellStr(sheet, "R28", complateyear)
		f.SetCellStr(sheet, "R29", apt.Tel)
		f.SetCellStr(sheet, "AA29", apt.Fax)

		if estimate.Parcel == 1 {
			boldUnderlineStyle, _ := f.NewStyle(&excelize.Style{
				Font: &excelize.Font{
					Bold:      true,
					Underline: "single",
				},
			})
			f.SetCellStr(sheet, "K40", "다. 장기수선계획 관련 도면, 서류 수령 및 보고서 납품은")
			f.SetCellStyle(sheet, "K40", "K40", boldUnderlineStyle)
			f.SetCellStr(sheet, "K41", "   택배활용으로 함.")
			f.SetCellStyle(sheet, "K41", "K41", boldUnderlineStyle)
		}

		sheet = "대가산출"

		f.SetCellValue(sheet, "K13", float64(estimate.Person2))
		f.SetCellValue(sheet, "K14", float64(estimate.Person3))
		f.SetCellValue(sheet, "K15", float64(estimate.Person4))
		f.SetCellValue(sheet, "K16", float64(estimate.Person5))

		f.SetCellValue(sheet, "L9", estimate.Personprice2)
		f.SetCellValue(sheet, "L10", estimate.Personprice3)
		f.SetCellValue(sheet, "L11", estimate.Personprice4)
		f.SetCellValue(sheet, "L12", estimate.Personprice5)

		f.SetCellValue(sheet, "L13", estimate.Personprice2)
		f.SetCellValue(sheet, "L14", estimate.Personprice3)
		f.SetCellValue(sheet, "L15", estimate.Personprice4)
		f.SetCellValue(sheet, "L16", estimate.Personprice5)

		f.SetCellValue(sheet, "K17", estimate.Financialprice)
		f.SetCellValue(sheet, "K18", estimate.Techprice)

		f.SetCellValue(sheet, "M19", estimate.Directprice)
		f.SetCellValue(sheet, "M20", estimate.Printprice)
		f.SetCellValue(sheet, "M21", estimate.Extraprice)

		f.SetCellValue(sheet, "M23", estimate.Saleprice)
		//f.SetCellValue(sheet, "I27", estimate.Price)

		//f.SetCellValue(sheet, "F10", fmt.Sprintf("%v원정(₩%v)", global.HumanMoney(estimate.Price), humanize.Comma(int64(estimate.Price))))

		if estimate.Event == 1 {
			f.SetCellValue(sheet, "A26", "이벤트 할인 금액")
		}

		sheet = "계약서1"

		if estimate.Subtype == 1 {
			// repair1.xlsx (조정)
			f.SetCellStr(sheet, "B9", apt.Name)
			f.SetCellStr(sheet, "B16", fmt.Sprintf("%v년", t.Year()))
			f.SetCellStr(sheet, "G41", tel)
		} else {
			// repair2.xlsx (수립)
			f.SetCellStr(sheet, "H9", apt.Name)
			f.SetCellStr(sheet, "H16", fmt.Sprintf("%v년", t.Year()))
			f.SetCellStr(sheet, "M41", tel)
		}

		f.UpdateLinkedValue()
	} else {
		// 비교 견적 -
		compareestimate := models.Compareestimate{}
		for _, v := range compareestimates {
			if v.Comparecompany == int64(typeid) {
				compareestimate = v
			}
		}
		log.Println("COM", compareestimate.Id)
		log.Println("comparecompany", compareestimate.Comparecompany)
		log.Println("COM", compareestimate.Type)
		t := time.ParseDate(compareestimate.Writedate)
		sheet = "표지"
		f.SetCellStr(sheet, "A5", fmt.Sprintf("%v 귀하", apt.Name))
		f.SetCellStr(sheet, "A11", t.Humandate())
		f.SetCellStr(sheet, "B9", fmt.Sprintf("일금%v원정", global.HumanMoney(compareestimate.Price)))
		f.SetCellValue(sheet, "G13", compareestimate.Price)
		f.SetCellStr(sheet, "A26", fmt.Sprintf("[ 견적조건 및 특기사항 ]    %v개동 %v세대", flatcount[0], apt.Familycount))

		if estimate.Subtype == 1 {
			f.SetCellValue(sheet, "A28", "2) 장기수선계획 보고서 1권 제출")
			f.SetCellValue(sheet, "A29", "")
		}

		sheet = "내역서"
		f.SetCellValue(sheet, "F12", compareestimate.Saleprice)
		f.SetCellValue(sheet, "F10", compareestimate.Printprice)

		if compareestimate.Person2 > 0 {
			f.SetCellValue(sheet, "B4", float64(compareestimate.Person2))
			f.SetCellValue(sheet, "C4", 1)
		}
		if compareestimate.Person3 > 0 {
			f.SetCellValue(sheet, "B5", float64(compareestimate.Person3))
			f.SetCellValue(sheet, "C5", 1)
		}
		if compareestimate.Person4 > 0 {
			f.SetCellValue(sheet, "B6", float64(compareestimate.Person4))
			f.SetCellValue(sheet, "C6", 1)
		}
		if compareestimate.Person5 > 0 {
			f.SetCellValue(sheet, "B7", float64(compareestimate.Person5))
			f.SetCellValue(sheet, "C7", 1)
		}

		f.SetCellValue(sheet, "E4", compareestimate.Personprice2)
		f.SetCellValue(sheet, "E5", compareestimate.Personprice3)
		f.SetCellValue(sheet, "E6", compareestimate.Personprice4)
		f.SetCellValue(sheet, "E7", compareestimate.Personprice5)

		f.SetCellValue(sheet, "B9", fmt.Sprintf("직접인건비*%v%%", compareestimate.Financialprice))
		f.SetCellFormula(sheet, "F9", fmt.Sprintf("=F8*%v%%", compareestimate.Financialprice))
		f.UpdateLinkedValue()
	}

	filename := fmt.Sprintf("%v.xlsx", global.UniqueId())
	fullFilename := fmt.Sprintf("%v/%v", config.UploadPath, filename)
	e := f.SaveAs(fullFilename)
	if e != nil {
		log.Println(e)
	}
	f.Close()

	return filename
}
