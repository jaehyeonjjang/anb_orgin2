package periodic

import (
	"fmt"
	"repair/models"
	"strings"
)

type OtherResult struct {
	Grade    string
	Items    []string
	Head     string
	Position string
	Good     int
}

func GetOtherByCategory(category int, items []models.Periodicother) []models.Periodicother {
	lists := make([]models.Periodicother, 0)

	for _, v := range items {
		if v.Category != category {
			continue
		}

		lists = append(lists, v)
	}

	return lists
}

func GetOtherByType(typeid int, items []models.Periodicother) []models.Periodicother {
	lists := make([]models.Periodicother, 0)

	for _, v := range items {
		if v.Type != typeid {
			continue
		}

		lists = append(lists, v)
	}

	return lists
}

func CheckStatus(datas []models.Periodicother) bool {
	for _, v := range datas {
		if v.Status != "" {
			return true
		}
	}

	return false
}

func GetItemByPosition(name string, items []models.Periodicother) *models.Periodicother {
	for _, v := range items {
		if v.Position == name {
			return &v
		}
	}

	return nil
}

func GetGrade(name string, items []models.Periodicother) string {
	for _, v := range items {
		if v.Position == name {
			return v.Status
		}
	}

	return ""
}

func Other(datas []models.Periodicother) map[int]OtherResult {
	ret := make(map[int]OtherResult)

	{
		pos := 10
		items := GetOtherByCategory(pos, datas)
		typeitems := GetOtherByType(2, items)

		strs := make([]string, 0)

		// 각 항목별로 결함 확인
		for _, v := range typeitems {
			if v.Status == "" {
				continue
			}

			statusList := strings.Split(v.Status, ",")
			defects := make([]string, 0)
			for _, s := range statusList {
				s = strings.TrimSpace(s)
				if s == "" || s == "상태양호" {
					continue
				}
				defects = append(defects, s)
			}

			if len(defects) > 0 {
				strs = append(strs, fmt.Sprintf("%v - %v", v.Position, strings.Join(defects, ", ")))
			}
		}

		good := 0
		if len(strs) == 0 {
			strs = append(strs, "해당 대지면적 내 추락방지시설은 전반적으로 양호한 상태로 조사되었다.")
			good = 1
		}

		ret[pos] = OtherResult{Grade: GetGrade("추락방지시설", items), Items: strs, Good: good}
	}

	{
		pos := 11
		items := GetOtherByCategory(pos, datas)
		typeitems := GetOtherByType(2, items)

		strs := make([]string, 0)

		// 각 항목별로 결함 확인
		for _, v := range typeitems {
			if v.Status == "" {
				continue
			}

			statusList := strings.Split(v.Status, ",")
			defects := make([]string, 0)
			for _, s := range statusList {
				s = strings.TrimSpace(s)
				if s == "" || s == "상태양호" {
					continue
				}
				defects = append(defects, s)
			}

			if len(defects) > 0 {
				strs = append(strs, fmt.Sprintf("%v - %v", v.Position, strings.Join(defects, ", ")))
			}
		}

		good := 0
		if len(strs) == 0 {
			strs = append(strs, "해당 대지면적 내 도로포장은 전반적으로 양호한 상태로 조사되었다.")
			good = 1
		}

		ret[pos] = OtherResult{Grade: GetGrade("도로포장", items), Items: strs, Good: good}
	}

	{
		pos := 12
		items := GetOtherByCategory(pos, datas)
		typeitems := GetOtherByType(2, items)

		strs := make([]string, 0)

		work := GetItemByPosition("시공", items)
		if work == nil {
			work = &models.Periodicother{}
		}

		good := 0
		if work.Status == "미시공" {
			strs = append(strs, "해당 대지면적 내 시공되어있지 않아 조사에서 제외하였다.")
			good = 2
		} else {
			// 각 항목별로 결함 확인
			for _, v := range typeitems {
				if v.Status == "" {
					continue
				}

				statusList := strings.Split(v.Status, ",")
				defects := make([]string, 0)
				for _, s := range statusList {
					s = strings.TrimSpace(s)
					if s == "" || s == "상태양호" {
						continue
					}
					defects = append(defects, s)
				}

				if len(defects) > 0 {
					strs = append(strs, fmt.Sprintf("%v - %v", v.Position, strings.Join(defects, ", ")))
				}
			}

			if len(strs) == 0 {
				if work.Status == "" {
					strs = append(strs, "신축 이음부는 전반적으로 양호한 상태로 조사되었다.")
				} else {
					strs = append(strs, fmt.Sprintf("신축 이음부는 %v로 시공되어 있으며 전반적으로 양호한 상태로 조사되었다.", work.Status))
				}
				good = 1
			}
		}

		ret[pos] = OtherResult{Grade: GetGrade("도로부 신축 이음부", items), Items: strs, Good: good}
	}

	{
		pos := 13
		items := GetOtherByCategory(pos, datas)
		typeitems := GetOtherByType(2, items)

		strs := make([]string, 0)

		work := GetItemByPosition("환기구", items)
		if work == nil {
			work = &models.Periodicother{}
		}

		// 각 항목별로 결함 확인
		for _, v := range typeitems {
			if v.Status == "" {
				continue
			}

			statusList := strings.Split(v.Status, ",")
			defects := make([]string, 0)
			for _, s := range statusList {
				s = strings.TrimSpace(s)
				if s == "" || s == "상태양호" {
					continue
				}
				defects = append(defects, s)
			}

			if len(defects) > 0 {
				strs = append(strs, fmt.Sprintf("%v - %v", v.Position, strings.Join(defects, ", ")))
			}
		}

		good := 0
		if len(strs) == 0 {
			str := ""
			if work.Status != "" {
				str = fmt.Sprintf("(%v) ", strings.ReplaceAll(work.Status, ",", ", "))
			}
			strs = append(strs, fmt.Sprintf("해당 건축물 내 %v환기구의 상태는 전반적으로 양호한 상태로 조사되었다.", str))
			good = 1
		}

		ret[pos] = OtherResult{Grade: GetGrade("환기구 등의 덮개", items), Items: strs, Good: good}
	}

	{
		pos := 14
		items := GetOtherByCategory(pos, datas)
		typeitems := GetOtherByType(2, items)

		strs := make([]string, 0)

		work := GetItemByPosition("건물외부 벽체", items)
		if work == nil {
			work = &models.Periodicother{}
		}

		// 선택된 재질 확인
		selectedMaterials := make(map[string]bool)
		if work.Status != "" {
			materials := strings.Split(work.Status, ",")
			for _, m := range materials {
				material := strings.TrimSpace(m)
				selectedMaterials[material] = true
			}
		}

		// 재질별 Order 매핑
		materialOrderMap := map[string]int{
			"석재":      143,
			"A/L 판넬":  144,
			"알루미늄 판넬": 144,
			"드라이비트":   145,
			"적벽돌":     146,
			"벽돌":      146,
		}

		head := ""
		position := ""
		good := 0

		// 각 항목별로 결함 확인
		for _, v := range typeitems {
			if v.Order == 141 {
				continue
			}

			if v.Status == "" {
				continue
			}

			// 선택된 재질에 해당하는 항목만 확인
			shouldInclude := false
			if len(selectedMaterials) == 0 {
				// 재질 선택 안 됨 - 모두 포함
				shouldInclude = true
			} else {
				// 선택된 재질과 매칭되는지 확인
				for material, order := range materialOrderMap {
					if selectedMaterials[material] && v.Order == order {
						shouldInclude = true
						break
					}
				}
			}

			if shouldInclude {
				statusList := strings.Split(v.Status, ",")
				defects := make([]string, 0)
				for _, s := range statusList {
					s = strings.TrimSpace(s)
					if s == "" || s == "상태양호" {
						continue
					}
					defects = append(defects, s)
				}

				if len(defects) > 0 {
					strs = append(strs, fmt.Sprintf("%v - %v", v.Position, strings.Join(defects, ", ")))
				}
			}
		}

		if len(strs) == 0 {
			str := ""
			if work.Status != "" {
				position = fmt.Sprintf("(%v 등)", strings.ReplaceAll(work.Status, ",", ", "))
				str = fmt.Sprintf("%v ", position)
			}
			head = fmt.Sprintf("건물외부 벽체 %v마감부위 상태는 양호한 것으로 확인됨.", str)
			good = 1
		}

		ret[pos] = OtherResult{Grade: GetGrade("외벽 마감재", items), Items: strs, Head: head, Position: position, Good: good}
	}

	{
		pos := 16
		items := GetOtherByCategory(pos, datas)
		typeitems := GetOtherByType(2, items)

		strs := make([]string, 0)

		// 각 부재별로 결함 확인
		for _, v := range typeitems {
			if v.Status == "" {
				continue
			}

			statusList := strings.Split(v.Status, ",")
			defects := make([]string, 0)
			for _, s := range statusList {
				s = strings.TrimSpace(s)
				if s == "" || s == "상태양호" {
					continue
				}
				defects = append(defects, s)
			}

			if len(defects) > 0 {
				// Position에서 앞쪽 중분류 제거 (예: "정착부 - 앵커 및 브라켓" → "앵커 및 브라켓")
				position := v.Position
				if idx := strings.Index(position, " - "); idx != -1 {
					position = strings.TrimSpace(position[idx+3:])
				}

				strs = append(strs, fmt.Sprintf("%v - %v", position, strings.Join(defects, ", ")))
			}
		}

		good := 0
		if len(strs) == 0 {
			strs = append(strs, "해당 건축물 내 부착물 등의 상태는 대체적으로 양호한 상태로 조사되었다.")
			good = 1
		}

		ret[pos] = OtherResult{Grade: GetGrade("부착물 등", items), Items: strs, Good: good}
	}

	return ret
}
