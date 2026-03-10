package api

import (
	"fmt"
	"os"
	"repair/controllers"
	"repair/global"
	"repair/global/config"
	"repair/models"
)

type PeriodicdataController struct {
	controllers.Controller
}

func (c *PeriodicdataController) Post_Index(items []models.Periodicdata) {
	for i, v := range items {
		filename := fmt.Sprintf("%v/periodicresult/%v/%v.jpg", config.UploadPath, v.Periodic, v.Extra["blueprint"].(models.Blueprint).Id)
		_, error := os.Stat(filename)

		if os.IsNotExist(error) {
			items[i].AddExtra("resultimage", "")
		} else {
			items[i].AddExtra("resultimage", fmt.Sprintf("periodicresult/%v/%v.jpg", v.Periodic, v.Extra["blueprint"].(models.Blueprint).Id))
		}
	}
}

func (c *PeriodicdataController) Post_Delete(item *models.Periodicdata) {
	// 삭제 후 결함도 이미지 재생성 요청
	if item.Periodic > 0 && item.Blueprint > 0 {
		global.SendNotify(global.Notify{Type: global.NotifyBlueprint, Periodic: item.Periodic, Blueprint: item.Blueprint})
	}
}

func (c *PeriodicdataController) Post_Deletebatch(items *[]models.Periodicdata) {
	// 삭제된 항목들의 blueprint별로 결함도 재생성 요청
	blueprints := make(map[int64]int64) // blueprint -> periodic
	for _, v := range *items {
		if v.Periodic > 0 && v.Blueprint > 0 {
			blueprints[v.Blueprint] = v.Periodic
		}
	}

	for blueprint, periodic := range blueprints {
		global.SendNotify(global.Notify{Type: global.NotifyBlueprint, Periodic: periodic, Blueprint: blueprint})
	}
}

func (c *PeriodicdataController) Post_DeleteByPeriodicBlueprint(periodic int64, blueprint int64) {
	// 도면의 모든 데이터 삭제 후 결함도 이미지 삭제
	if periodic > 0 && blueprint > 0 {
		filename := fmt.Sprintf("%v/periodicresult/%v/%v.jpg", config.UploadPath, periodic, blueprint)
		os.Remove(filename)
	}
}
