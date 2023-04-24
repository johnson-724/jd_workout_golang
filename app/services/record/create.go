package record

import (
	"github.com/gin-gonic/gin"
	"jd_workout_golang/app/middleware"
	"jd_workout_golang/app/models"
	equipRepo "jd_workout_golang/app/repositories/equip"
	recordRepo "jd_workout_golang/app/repositories/record"
)

type createBody struct {
	EquipId uint    `json:"equip_id" form:"equip_id" binding:"required"`
	Weight  float32 `json:"weight" form:"weight" binding:"required"`
	Reps    uint    `json:"reps" form:"reps" binding:"required"`
	Note    string  `json:"note" form:"note"`
}

// create record
// @Summary create record
// @Description create record
// @Tags Record
// @Accept json
// @Produce json
// @Param createBody body createBody true "createBody"
// @Success 200 {string} string "{'message': 'create success', 'id' : '1'}"
// @Failure 422 {string} string "{'message': '缺少必要欄位', 'error': 'error message'}"
// @Failure 403 {string} string "{'message': 'jwt token error', 'error': 'error message'}"
// @Router /record [post]
// @Security Bearer
func CreateRecord(c *gin.Context) {
	createBody := &createBody{}
	if err := c.ShouldBind(&createBody); err != nil {
		c.JSON(422, gin.H{
			"message": "缺少必要欄位",
			"error":   err.Error(),
		})

		return
	}

	equip, err := equipRepo.GetEquip(uint64(createBody.EquipId), middleware.Uid)

	if err != nil {
		c.JSON(422, gin.H{
			"message": "equip not found",
			"error":   err.Error(),
		})

		return
	}

	record := &models.Record{
		EquipId: createBody.EquipId,
		UserId: middleware.Uid,
		Name:    equip.Name,
		Weight:  createBody.Weight,
		Reps:    createBody.Reps,
		Note:    createBody.Note,
	}

	recordRepo.Create(record)

	if err != nil {
		c.JSON(422, gin.H{
			"message": "create error",
			"error":   err.Error(),
		})

		c.Abort()

		return
	}

	c.JSON(200, gin.H{
		"message": "create success",
		"id":      record.ID,
	})
}
