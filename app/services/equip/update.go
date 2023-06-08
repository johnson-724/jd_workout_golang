package equip

import (
	"github.com/gin-gonic/gin"
	"jd_workout_golang/app/middleware"
	repo "jd_workout_golang/app/repositories/equip"
	"strconv"
)

type updateFrom struct {
	Name string `json:"name" form:"name"`
	Note string `json:"note" form:"note"`
}

// update personal equip
// @Summary update equip
// @Description update equip for personal user
// @Tags Equip
// @Accept x-www-form-urlencoded
// @Produce json
// @Param id path integer true "equip id"
// @Param weights formData updateFrom false "note for equip"
// @Param image formData file false "image for equip"
// @Success 200 {string} string "{'message': 'create success'}"
// @Failure 422 {string} string "{'message': '缺少必要欄位', 'error': 'error message'}"
// @Failure 403 {string} string "{'message': 'jwt token error', 'error': 'error message'}"
// @Router /equip/{id} [patch]
// @Security Bearer
func UpdateEquip(c *gin.Context) {

	id := c.Param("id")

	weightId, err := strconv.ParseUint(id, 10, 32)

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

	equip, err := repo.GetEquip(weightId, middleware.Uid)

	if err != nil {
		c.JSON(422, gin.H{
			"message": "equip not found",
			"error":   err.Error(),
		})

		c.Abort()

		return
	}

	if _, ok := c.GetPostForm("name"); ok {
		equip.Name = updateFrom.Name
	}

	if _, ok := c.GetPostForm("note"); ok {
		equip.Note = &updateFrom.Note
	}

	path, err := StoreFile(c)

	if err != nil {
		c.JSON(422, gin.H{
			"message": "file error",
			"error":   err.Error(),
		})

		c.Abort()

		return
	}

	equip.Image = path

	repo.Update(equip)

	c.JSON(200, gin.H{
		"message": "equip updated",
	})
}
