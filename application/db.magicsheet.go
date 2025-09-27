package application

import (
	"fmt"
	"log"
	"github.com/spo-iitk/ras-backend/student"

	"github.com/gin-gonic/gin"
)

func CreateMagicSheetData(ctx *gin.Context, magicsheetdata *MagicSheet) error {
	tx := db.WithContext(ctx).Create(magicsheetdata)
	return tx.Error
}

// ///////////////////////////
func FetchMagicSheetDataForCoco(ctx *gin.Context, id uint, CocoData *[]MagicSheet) error {

	tx := db.WithContext(ctx).Where("coco_id = ?", id).Find(CocoData)
	return tx.Error

}
func FetchMagicSheetData(ctx *gin.Context, id uint, Data *[]MagicSheet) error {

	tx := db.WithContext(ctx).Where("recruitment_cycle_id=?", id).Find(Data)
	return tx.Error

}
func FetchComanyMagicSheetData(ctx *gin.Context, pids PIDRequest, data *[]MagicSheet) error {
	tx := db.WithContext(ctx).Where("proforma_id IN ?", pids.IDs).Find(data)
	return tx.Error
}

func FetchStudentInfo(ctx *gin.Context, ids []uint) ([]student.Student, error) {
	var students []student.Student
	if len(ids) == 0 {
		return students, nil
	}
	studentDB := student.GetDB()
	err := studentDB.WithContext(ctx).
		Select("id", "roll_no", "name", "specialization", "dob", "iitk_email", "personal_email", "phone", "alternate_phone", "current_cpi", "friend_name", "friend_phone", "program_department_id", "secondary_program_department_id").
		Where("id IN ?", ids).
		Find(&students).Error

	return students, err
}

func getRollNumberByID(sid uint) (string, error) {
	db := student.GetDB()
	var student student.Student

	err := db.Select("roll_no").First(&student, sid).Error
	if err != nil {
		return "", fmt.Errorf("failed to fetch roll number for student ID %d: %w", sid, err)
	}
	return student.RollNo, nil
}

func AddMagicSheetdata(ctx *gin.Context, rollNumbers []string, rid uint, pid uint) error {
	var sheets []MagicSheet
	studentDB := student.GetDB()

	for _, roll := range rollNumbers {
		var student student.Student
		if err := studentDB.WithContext(ctx).
			Where("roll_no = ?", roll).
			First(&student).Error; err != nil {
			return fmt.Errorf("student with roll number %s not found: %w", roll, err)
		}

		sheets = append(sheets, MagicSheet{
			ProformaID:         pid,
			RecruitmentCycleID: rid,
			StudentID:          student.ID,
			CocoID:             0,
			Status:             "Pending",
		})
	}

	if len(sheets) == 0 {
		return fmt.Errorf("empty slice found")
	}

	return db.WithContext(ctx).Create(&sheets).Error
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

func UpdateMagicSheetFields(ctx *gin.Context, data *EditableMagicSheetFields) error {
	tx := db.WithContext(ctx).
		Model(&MagicSheet{}).
		Where("id = ?", *data.ID).
		Updates(map[string]interface{}{
			"status":      data.Status,
			"r1_in_time":  data.R1InTime,
			"r1_out_time": data.R1OutTime,
			"comments":    data.Comments,
		})

	return tx.Error
}

func AssignCocoToMagicSheets(ctx *gin.Context, cocoID uint, rcID uint, req PIDRequestt) error {
	log.Printf("➡️ Assigning coco_id: %d to RC: %d, Proforma IDs: %v", cocoID, rcID, req.IDs)

	tx := db.WithContext(ctx).Model(&MagicSheet{}).
		Where("recruitment_cycle_id = ? AND proforma_id IN ?", rcID, req.IDs).
		Select("coco_id"). // 👈 ensures this field is included even if 0 or NULL
		Updates(map[string]interface{}{"coco_id": cocoID})

	log.Printf("⬆️ Rows updated: %d", tx.RowsAffected)
	return tx.Error
}

func DeleteMagicSheetData(ctx *gin.Context, id uint) error {
	tx := db.WithContext(ctx).Where("id = ?", id).Delete(&MagicSheet{})
	return tx.Error
}
