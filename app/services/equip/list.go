package equip

import (
	"github.com/gin-gonic/gin"
	"jd_workout_golang/app/middleware"
	"jd_workout_golang/app/models"
	repo "jd_workout_golang/app/repositories/equip"
)

type equipListRequest struct {
	Page    int ` json:"currentPage" form:"currentPage"`
	PerPage int ` json:"perPage" form:"perPage"`
}

type equipListResponse struct {
	Page    int            `json:"currentPage" form:"currentPage"`
	PerPage int            `json:"perPage" form:"perPage"`
	Data    []models.Equip `json:"data"`
	Total   int64          `json:"total"`
}

// get personal equip list
// @Summary equip list
// @Description equip list for personal user
// @Tags Equip
// @Accept x-www-form-urlencoded
// @Produce json
// @Param equipList query equipListRequest true "equipList"
// @Success 200 {object} equipListResponse
// @Failure 422 {string} string "{'message': '缺少必要欄位', 'error': 'error message'}"
// @Failure 403 {string} string "{'message': 'jwt token error', 'error': 'error message'}"
// @Router /equip [get]
// @Security Bearer
func List(c *gin.Context) {
	paginate := equipListRequest{
		Page:    1,
		PerPage: 10,
	}

	if err := c.ShouldBind(&paginate); err != nil {
		c.JSON(422, gin.H{
			"message": "缺少必要欄位",
			"error":   err.Error(),
		})

		return
	}

	paginateCondition := repo.PaginateCondition{
		Page:    paginate.Page,
		PerPage: paginate.PerPage,
	}

	data, count, err := repo.GetEqupis(paginateCondition, middleware.Uid)

	if err != nil {
		c.JSON(422, gin.H{
			"message": "get equip list error",
			"error":   err.Error(),
		})

		c.Abort()

		return
	}

	c.JSON(200, equipListResponse{
		Page:    paginate.Page,
		PerPage: paginate.PerPage,
		Data:    *data,
		Total:   *count,
	})
}
