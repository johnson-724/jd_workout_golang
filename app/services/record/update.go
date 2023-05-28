package record

import (
	"jd_workout_golang/app/middleware"
	equipRepo "jd_workout_golang/app/repositories/equip"
	recordRepo "jd_workout_golang/app/repositories/record"
	"strconv"

	"github.com/gin-gonic/gin"
)

type updateBody struct {
	EquipId uint    `json:"equip_id" form:"equip_id" binding:"required"`
	Weight  float32 `json:"weight" form:"weight" binding:"required"`
	Reps    uint    `json:"reps" form:"reps" binding:"required"`
	Note    string  `json:"note" form:"note"`
}

// update record
// @Summary update record
// @Description update record
// @Tags Record
// @Accept json
// @Produce json
// @Param id path string true "record id"
// @Param updateBody body updateBody true "updateBody"
// @Success 200 {string} string "{'message': 'update success', 'id' : '1'}"
// @Failure 422 {string} string "{'message': '缺少必要欄位', 'error': 'error message'}"
// @Failure 403 {string} string "{'message': 'jwt token error', 'error': 'error message'}"
// @Router /record/{id} [patch]
// @Security Bearer
func UpdateRecord(c *gin.Context) {
	id := c.Param("id")

	recordId, err := strconv.ParseUint(id, 10, 32)

	if err != nil {
		c.JSON(422, gin.H{
			"message": "uri id error",
		})

		c.Abort()

		return
	}

	updateBody := &updateBody{}
	if err := c.ShouldBind(&updateBody); err != nil {
		c.JSON(422, gin.H{
			"message": "缺少必要欄位",
			"error":   err.Error(),
		})

		return
	}

	record, err := recordRepo.GetRecord(uint64(recordId), middleware.Uid)

	if err != nil {
		c.JSON(422, gin.H{
			"message": "record not found",
			"error":   err.Error(),
		})

		return
	}

	equip, err := equipRepo.GetEquip(uint64(updateBody.EquipId), middleware.Uid)

	if err != nil {
		c.JSON(422, gin.H{
			"message": "equip not found",
			"error":   err.Error(),
		})

		return
	}

	record.EquipId = updateBody.EquipId
	record.Name = equip.Name
	record.Weight = updateBody.Weight
	record.Reps = updateBody.Reps
	record.Note = updateBody.Note

	recordRepo.Update(record)

	if err != nil {
		c.JSON(422, gin.H{
			"message": "update error",
			"error":   err.Error(),
		})

		c.Abort()

		return
	}

	c.JSON(200, gin.H{
		"message": "update success",
		"id":      record.ID,
	})
}
