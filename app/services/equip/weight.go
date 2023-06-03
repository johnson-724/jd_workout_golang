package equip

import (
	"encoding/json"
	"jd_workout_golang/app/middleware"
	"strconv"
	repo "jd_workout_golang/app/repositories/equip"
	"github.com/gin-gonic/gin"
)

type weightForm struct {
	Weights  []float32 `json:"weights" form:"weights" binding:"required"`
}

// update personal equip weight
// @Summary update equip weight
// @Description update equip weight for personal user
// @Tags Equip
// @Accept json
// @Produce json
// @Param id path integer true "equip id"
// @Param weights body weightForm false "note for equip"
// @Success 200 {string} string "{'message': 'create success'}"
// @Failure 422 {string} string "{'message': '缺少必要欄位', 'error': 'error message'}"
// @Failure 403 {string} string "{'message': 'jwt token error', 'error': 'error message'}"
// @Router /equip/{id}/weight [put]
// @Security Bearer
func UpdateWeight(c *gin.Context) {

	id := c.Param("id")
	
	weightId , err:= strconv.ParseUint(id, 10, 32)
	
	if err != nil {
		c.JSON(422, gin.H{
			"message": "uri id error",
		})

		c.Abort()

		return
	}

	weightForm := weightForm{}
	if err := c.ShouldBind(&weightForm); err != nil {
		c.JSON(422, gin.H{
			"message": "缺少必要欄位",
			"error":   err.Error(),
		})

		c.Abort()

		return
	}

	equip, err := repo.GetEquip(weightId, middleware.Uid)

	if err != nil {
		c.JSON(422, gin.H{
			"message": "equip not found",
			"error":   err.Error(),
		})

		c.Abort()

		return
	}

	json, _ := json.Marshal(&weightForm.Weights)
	jsonPtr := string(json)
	equip.Weights = &jsonPtr

	repo.Update(equip)

	c.JSON(200, gin.H{
		"message": "weights updated",
	})
}
