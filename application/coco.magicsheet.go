package application

import (
	"net/http"
	"github.com/spo-iitk/ras-backend/util"

	"github.com/gin-gonic/gin"

)

//1 route: /api/coco/:cocoID/magic-sheet
func GetMagicSheetData(ctx *gin.Context) {
	cocoIDParam := ctx.Param("cocoID")
	cocoID, err := util.ParseUint(cocoIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid COCO ID"})
		return
	}

	var cocoData []MagicSheet
	err = FetchMagicSheetDataForCoco(ctx, uint(cocoID), &cocoData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch magic sheet data"})
		return
	}

	ctx.JSON(http.StatusOK, cocoData)
}
type MagicSheetUpdateInput struct {
	R1InTime  string `json:"r1_in_time"`
	R1OutTime string `json:"r1_out_time"`
	Status    string `json:"status"`
}
//2nd route: /api/coco/magic-sheet/:id
func UpdateMagicSheetData(ctx *gin.Context) {
	idParam := ctx.Param("id")
	idUint64, err := util.ParseUint(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	id := uint(idUint64)


	var input MagicSheetUpdateInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err = UpdateMagicSheetTimes(ctx, input, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch magic sheet data"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Updated successfully"})
}
