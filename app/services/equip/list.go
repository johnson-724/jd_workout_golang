package equip

import (
	"fmt"
	"jd_workout_golang/app/middleware"
	"jd_workout_golang/app/models"
	db "jd_workout_golang/lib/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type equipListRequest struct {
	Page    int ` json:"currentPage" form:"currentPage"`
	PerPage int ` json:"perPage" form:"perPage"`
}

type equipListResponse struct {
	Page    int            `json:"currentPage" form:"currentPage"`
	PerPage int            `json:"perPage" form:"perPage"`
	Data    []models.Equip `json:"data"`
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

	db := db.InitDatabase()

	data, err := getEqupis(paginate, middleware.Uid, db)

	fmt.Println("test")

	if err != nil {
		c.JSON(422, gin.H{
			"message": "get equip list error",
			"error":  err.Error(),
		})

		c.Abort()

		return
	}

	c.JSON(200, data)
}

func getEqupis(equipListRequest equipListRequest, uid uint, db *gorm.DB) (*equipListResponse, error) {
	data := []models.Equip{}

	result := db.Where("user_id = ?", uid).Order("name asc").Scopes(Paginate(equipListRequest.Page, equipListRequest.PerPage)).Find(&data)

	if result.Error != nil {
		return nil, result.Error
	}

	return &equipListResponse{
		Page:    equipListRequest.Page,
		PerPage: equipListRequest.PerPage,
		Data:    data,
	}, nil
}

func Paginate(currentPage int, perPage int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		fmt.Printf("currentPage: %d, perPage: %d", currentPage, perPage)
		offset := (currentPage - 1) * perPage

		return db.Offset(offset).Limit(perPage)
	}
}
