package equip

import (
	"jd_workout_golang/app/middleware"
	"jd_workout_golang/app/models"
	db "jd_workout_golang/lib/database"

	"github.com/gin-gonic/gin"
)

// create personal equip with name and note
// @Summary create equip
// @Description create equip for personal user
// @Tags Equip
// @Accept json
// @Produce json
// @Param name formData string true "equip name"
// @Param note formData string false "note for equip"
// @Success 200 {string} string "{'message': 'create success'}"
// @Failure 422 {string} string "{'message': '缺少必要欄位', 'error': 'error message'}"
// @Failure 403 {string} string "{'message': 'jwt token error', 'error': 'error message'}"
// @Router /equip [post]
func CreateEquip(c *gin.Context) {
	createBody := struct {
		Name    string `json:"name" form:"name" binding:"required"`
		Note string `json:"note" form:"note"`
	}{}

	if err := c.ShouldBind(&createBody); err != nil {
		c.JSON(422, gin.H{
			"message": "缺少必要欄位",
			"error":   err.Error(),
		})

		return
	}

	db := db.InitDatabase()

	equip := models.Equip{
		UserId: middleware.Uid,
		Name: createBody.Name,
		Note: createBody.Note,
	}

	tx := db.Create(&equip)

	if tx.Error != nil {
		c.JSON(422, gin.H{
			"message": "create error",
			"error":   tx.Error.Error(),
		})

		c.Abort()

		return
	}

	c.JSON(200, gin.H{
		"message": "create success",
	})
}
