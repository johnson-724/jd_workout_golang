package equip

import (
	"github.com/gin-gonic/gin"
	"jd_workout_golang/app/middleware"
	repo "jd_workout_golang/app/repositories/equip"
	"strconv"
)

// delete personal equip
// @Summary delete equip
// @Description delete equip for personal user
// @Tags Equip
// @Accept json
// @Produce json
// @Param id path integer true "equip id"
// @Success 200 {string} string "{'message': 'deleted success'}"
// @Failure 422 {string} string "{'message': '缺少必要欄位', 'error': 'error message'}"
// @Failure 403 {string} string "{'message': 'jwt token error', 'error': 'error message'}"
// @Router /equip/{id} [delete]
// @Security Bearer
func DeleteEquip(c *gin.Context) {

	id := c.Param("id")

	weightId, err := strconv.ParseUint(id, 10, 32)

	if err != nil {
		c.JSON(422, gin.H{
			"message": "uri id error",
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

	repo.Delete(equip)

	c.JSON(200, gin.H{
		"message": "equip deleted",
	})
}
