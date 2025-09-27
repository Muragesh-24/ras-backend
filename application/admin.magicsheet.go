package application

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spo-iitk/ras-backend/auth"
	"github.com/spo-iitk/ras-backend/student"
	"github.com/spo-iitk/ras-backend/util"

	"github.com/gin-gonic/gin"
)
// GET /admin/magic-sheets/:rid
func getAllMagicSheetsHandler(ctx *gin.Context) {
	rc_id, err := util.ParseUint(ctx.Param("rid"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid rc_id"})
		return
	}

	var sheets []MagicSheet
	err = FetchMagicSheetData(ctx, rc_id, &sheets)
	if err != nil {
		fmt.Println("Error fetching magic sheets:", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			
			"error":   "Failed to fetch magic sheets",
			"details": err.Error(),
		})
		return
	}

	// Collect unique student IDs
	studentIDs := make([]uint, 0, len(sheets))
	for _, sheet := range sheets {
		studentIDs = append(studentIDs, sheet.StudentID)
	}

	// Fetch student info from separate student DB
	students, err := FetchStudentInfo(ctx, studentIDs)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create lookup map
	studentMap := make(map[uint]student.Student)
	for _, student := range students {
		studentMap[student.ID] = student
	}

	// Merge and format response
	var result []gin.H
	for _, sheet := range sheets {
		student := studentMap[sheet.StudentID]
		var cocoName string = "-" // default value

		var coco auth.User
		if err := auth.FetchUserByid(ctx, &coco, sheet.CocoID); err == nil {
			cocoName = coco.Name
		} else {
			fmt.Printf("Failed to fetch coco: %s\n", err.Error())
			// Don't return, just log and continue
		}

		// Continue with building the result
		result = append(result, gin.H{
			"id":                   sheet.ID,
			"created_at":           sheet.CreatedAt,
			"updated_at":           sheet.UpdatedAt,
			"student_id":           sheet.StudentID,
			"proforma_id":          sheet.ProformaID,
			"recruitment_cycle_id": sheet.RecruitmentCycleID,
			"coco_id":              cocoName, // safe value
			"r1_in_time":           sheet.R1InTime,
			"r1_out_time":          sheet.R1OutTime,
			"comments":             sheet.Comments,
			"status":               sheet.Status,

			"roll_no":                         student.RollNo,
			"name":                            student.Name,
			"specialization":                  student.Specialization,
			"dob":                             student.DOB,
			"iitk_email":                      student.IITKEmail,
			"personal_email":                  student.PersonalEmail,
			"phone":                           student.Phone,
			"alternate_phone":                 student.AlternatePhone,
			"secondary_program_department_id": student.SecondaryProgramDepartmentID,
			"program_department_id":           student.ProgramDepartmentID,
			"current_cpi":                     student.CurrentCPI,
			"friend_name":                     student.FriendName,
			"friend_phone":                    student.FriendPhone,
		})

	}

	ctx.JSON(http.StatusOK, result)
}

type PIDRequest struct {
	IDs []uint `json:"pids"`
}

func getAllCompanyMagicSheetsHandler(ctx *gin.Context) {
	var req PIDRequest
	var sheets []MagicSheet

	
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := FetchComanyMagicSheetData(ctx, req, &sheets); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch Company magic sheets",
		})
		return
	}

	studentIDs := make([]uint, 0, len(sheets))
	for _, sheet := range sheets {
		studentIDs = append(studentIDs, sheet.StudentID)
	}

	students, err := FetchStudentInfo(ctx, studentIDs)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	studentMap := make(map[uint]student.Student)
	for _, student := range students {
		studentMap[student.ID] = student
	}

	var result []gin.H
	for _, sheet := range sheets {
		student := studentMap[sheet.StudentID]

		var cocoName string = "-" // default value

		var coco auth.User
		if err := auth.FetchUserByid(ctx, &coco, sheet.CocoID); err == nil {
			cocoName = coco.Name
		} else {
			fmt.Printf("Failed to fetch coco: %s\n", err.Error())
			// Don't return, just log and continue
		}

		// Continue with building the result
		result = append(result, gin.H{
			"id":                   sheet.ID,
			"created_at":           sheet.CreatedAt,
			"updated_at":           sheet.UpdatedAt,
			"student_id":           sheet.StudentID,
			"proforma_id":          sheet.ProformaID,
			"recruitment_cycle_id": sheet.RecruitmentCycleID,
			"coco_id":              cocoName,
			"r1_in_time":           sheet.R1InTime,
			"r1_out_time":          sheet.R1OutTime,
			"comments":             sheet.Comments,
			"status":               sheet.Status,

			"roll_no":                         student.RollNo,
			"name":                            student.Name,
			"specialization":                  student.Specialization,
			"dob":                             student.DOB,
			"iitk_email":                      student.IITKEmail,
			"personal_email":                  student.PersonalEmail,
			"phone":                           student.Phone,
			"alternate_phone":                 student.AlternatePhone,
			"secondary_program_department_id": student.SecondaryProgramDepartmentID,
			"program_department_id":           student.ProgramDepartmentID,
			"current_cpi":                     student.CurrentCPI,
			"friend_name":                     student.FriendName,
			"friend_phone":                    student.FriendPhone,
		})
	}

	ctx.JSON(http.StatusOK, result)
}

type EditableMagicSheetFields struct {
	ID        *uint      `json:"ID" binding:"required"`
	Status    string     `json:"status"`
	R1InTime  *time.Time `json:"r1_in_time"`
	R1OutTime *time.Time `json:"r1_out_time"`
	Comments  string     `json:"comments"`
}

func addMagicsheetdataHandler(ctx *gin.Context) {
	var rollNumbers []string

	rc_id, _ := util.ParseUint(ctx.Param("rid"))
	p_id, _ := util.ParseUint(ctx.Param("pid"))

	if err := ctx.ShouldBindJSON(&rollNumbers); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := AddMagicSheetdata(ctx, rollNumbers, rc_id, p_id); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Magic sheets added successfully"})
}

type PIDRequestt struct {
	IDs []uint `json:"ids"`
}

type AssignCocoRequest struct {
	EmailId string      `json:"email_id"`
	PIDs    PIDRequestt `json:"pids"`
}

func assignCocoHandler(ctx *gin.Context) {
	rc_id, err := util.ParseUint(ctx.Param("rid"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recruitment cycle ID"})
		return
	}

	var req AssignCocoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var coco auth.User
	if err := auth.FetchUser(ctx, &coco, req.EmailId); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("Failed to fetch coco: %s", err.Error()),
		})
		return
	}

	if err := AssignCocoToMagicSheets(ctx, coco.ID, rc_id, req.PIDs); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Coco assigned successfully"})
}

// // POST /admin/magic-sheets
func createMagicSheetHandler(ctx *gin.Context) {
	var sheet MagicSheet
	if err := ctx.ShouldBindJSON(&sheet); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := CreateMagicSheetData(ctx, &sheet)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, sheet)
}

// PUT /admin/magic-sheets/:id
func updateMagicSheetHandler(ctx *gin.Context) {
	var sheet EditableMagicSheetFields
	if err := ctx.ShouldBindJSON(&sheet); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := UpdateMagicSheetFields(ctx, &sheet)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update magic sheet"})
		fmt.Println(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "Updated"})
}

// DELETE /admin/magic-sheets/:id
func deleteMagicSheetHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	// id, err := strconv.Atoi(idStr)
	id, err := util.ParseUint(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = DeleteMagicSheetData(ctx, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Could not delete magic sheet"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "Deleted"})
}
