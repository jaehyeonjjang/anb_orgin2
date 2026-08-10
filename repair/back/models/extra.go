package models

import log "github.com/sirupsen/logrus"

type EstimateExtra struct {
	Estimate Estimate          `json:"estimate"`
	Compares []Compareestimate `json:"compares"`
}

func (p *PeriodicdataManager) FakeDelete(periodic int64, blueprint int64) {
	query := "update periodicdata_tb set pd_blueprint = ?, pd_periodic = ? where pd_blueprint = ? and pd_periodic = ?"
	_, err := p.Exec(query, blueprint*-1, periodic*-1, blueprint, periodic)

	log.Println(err)
}

// FakeDeleteByTypeRanges 는 지정한 type 범위(각 원소 = [min, max), max 는 포함하지 않음)에
// 해당하는 periodicdata 만 소프트 삭제한다. 클라이언트가 특정 아이콘셋만 부분 업로드할 때 사용한다.
func (p *PeriodicdataManager) FakeDeleteByTypeRanges(periodic int64, blueprint int64, ranges [][2]int) {
	if len(ranges) == 0 {
		return
	}

	conds := make([]string, 0, len(ranges))
	args := make([]any, 0, len(ranges)*2+4)
	args = append(args, blueprint*-1, periodic*-1, blueprint, periodic)
	for _, r := range ranges {
		conds = append(conds, "(pd_type >= ? and pd_type < ?)")
		args = append(args, r[0], r[1])
	}

	query := "update periodicdata_tb set pd_blueprint = ?, pd_periodic = ? where pd_blueprint = ? and pd_periodic = ? and (" +
		joinStrings(conds, " or ") + ")"
	_, err := p.Exec(query, args...)

	log.Println(err)
}

func joinStrings(items []string, sep string) string {
	out := ""
	for i, v := range items {
		if i > 0 {
			out += sep
		}
		out += v
	}
	return out
}

func (p *PeriodicdataimageManager) FakeDelete(periodic int64, periodicdata int64) {
	query := "update periodicdataimage_tb set pi_periodicdata = ?, pi_periodic = ? where pi_periodicdata = ?"
	_, err := p.Exec(query, periodicdata*-1, periodic*-1, periodicdata)

	log.Println(err)
}
