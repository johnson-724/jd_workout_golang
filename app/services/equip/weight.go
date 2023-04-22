package equip

import (
	"encoding/json"
	"fmt"
	"jd_workout_golang/app/middleware"
	"jd_workout_golang/app/models"
	db "jd_workout_golang/lib/database"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type updateFrom struct {
	Weights  []float32 `json:"weights" form:"weights" binding:"required"`
}

// update personal equip weight
// @Summary update equip weight
// @Description update equip weight for personal user
// @Tags Equip
// @Accept json
// @Produce json
// @Param id path integer true "equip id"
// @Param weights body updateFrom false "note for equip"
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

	updateFrom := updateFrom{}
	if err := c.ShouldBind(&updateFrom); err != nil {
		c.JSON(422, gin.H{
			"message": "缺少必要欄位",
			"error":   err.Error(),
		})

		c.Abort()

		return
	}
	
	db := db.InitDatabase()

	equip, err := getEquip(weightId, middleware.Uid, db)

	if err != nil {
		c.JSON(422, gin.H{
			"message": "equip not found",
			"error":   err.Error(),
		})

		c.Abort()

		return
	}

	json,_ := json.Marshal(&updateFrom.Weights)
	equip.Weights = string(json)

	db.Save(&equip)

	c.JSON(200, gin.H{
		"message": "weights updated",
	})
}

func getEquip(equipId uint64, uid uint, db *gorm.DB) (*models.Equip, error)  {
	equip := models.Equip{}
	
	result := db.Where("user_id = ?", uid).Where("id = ? ", equipId).First(&equip)

	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		return &models.Equip{}, fmt.Errorf("equip not found : %w", result.Error)
	}

	return &equip, nil
}
