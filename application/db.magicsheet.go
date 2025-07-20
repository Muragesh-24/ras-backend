package application

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)



func CreateMagicSheetData(ctx *gin.Context,magicsheetdata * MagicSheet) error{
tx := db.WithContext(ctx).Create(magicsheetdata)
	return tx.Error
}

/////////////////////////////
func FetchMagicSheetDataForCoco(ctx *gin.Context,id uint,CocoData * []MagicSheet) error  {

	tx := db.WithContext(ctx).Where("coco_id = ?", id).Find(CocoData)
	return tx.Error

}
func FetchMagicSheetData(ctx *gin.Context,id uint,Data * []MagicSheet) error  {

	tx := db.WithContext(ctx).Where("rc_id=?",id).Find(Data)
	return tx.Error

}
func FetchComanyMagicSheetData(ctx *gin.Context, pids []uint, data *[]MagicSheet) error {
	tx := db.WithContext(ctx).Where("proforma_id IN ?", pids).Find(data)
	return tx.Error
}

func UpdateMagicSheetTimes(ctx *gin.Context, data MagicSheetUpdateInput, id uint) error {
	fmt.Println(data.R1InTime, data.R1OutTime, data.Status, id)
tx := db.WithContext(ctx).
		Model(&MagicSheet{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"r1_in_time":  data.R1InTime,
			"r1_out_time": data.R1OutTime,
			"status":      data.Status,
		})


	return tx.Error
}

func UpdateMagicSheetFull(ctx *gin.Context, data *MagicSheet) error {
	tx := db.WithContext(ctx).
		Clauses(clause.Returning{}). 
		Where("id = ?", data.ID).
		Updates(data)

	return tx.Error
}

func DeleteMagicSheetData(ctx *gin.Context, id uint) error {
	tx := db.WithContext(ctx).Where("id = ?", id).Delete(&MagicSheet{})
	return tx.Error
}


