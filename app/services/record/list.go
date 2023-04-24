package record

import (
	"jd_workout_golang/app/middleware"
	"jd_workout_golang/app/models"
	"jd_workout_golang/app/repositories/pageinate"
	repo "jd_workout_golang/app/repositories/record"
	"time"

	"github.com/gin-gonic/gin"
)

type recordListRequest struct {
	Page    int ` json:"currentPage" form:"currentPage"`
	PerPage int ` json:"perPage" form:"perPage"`
}

type recordListResponse struct {
	Page    int             `json:"currentPage" form:"currentPage"`
	PerPage int             `json:"perPage" form:"perPage"`
	Data    []models.Record `json:"data"`
	Total   int64           `json:"total"`
}

// get personal record list
// @Summary record list
// @Description record list for personal user
// @Tags Equip
// @Accept x-www-form-urlencoded
// @Produce json
// @Param equipList query recordListRequest true "equipList"
// @Success 200 {object} recordListResponse
// @Failure 422 {string} string "{'message': '缺少必要欄位', 'error': 'error message'}"
// @Failure 403 {string} string "{'message': 'jwt token error', 'error': 'error message'}"
// @Router /record [get]
// @Security Bearer
func List(c *gin.Context) {
	paginate := recordListRequest{
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

	paginateCondition := pageinate.PaginateCondition{
		Page:    paginate.Page,
		PerPage: paginate.PerPage,
	}

	data, count, err := repo.GetRecords(paginateCondition, middleware.Uid)

	if err != nil {
		c.JSON(422, gin.H{
			"message": "get record list error",
			"error":   err.Error(),
		})

		c.Abort()

		return
	}

	c.JSON(200, recordListResponse{
		Page:    paginate.Page,
		PerPage: paginate.PerPage,
		Data:    *data,
		Total:   *count,
	})
}

type dateGroup struct {
	Date   string
	Start  time.Time
	End    time.Time
	Equips map[string]equipGroup
}

type equipGroup struct {
	ID      uint
	Name    string
	Note    string
	Records map[string]recordDetail
}

type recordDetail struct {
	ID     uint
	Weight float64
	Reps   int
	Sets   float64 // Weight * Reps * Count
	Note   []string
}

func groupBy(data []models.Record) {
	group := make(map[string]interface{})
	for _, v := range data {
		if _, ok := group[v.CreatedAt.Format("2006-01-02")]; !ok {
			equipGroup := make(map[string]equipGroup)
			// recordGroup := make(map[string]recordDetail)

			group[v.CreatedAt.Format("2006-01-02")] = &dateGroup{
				Date:   v.CreatedAt.Format("2006-01-02"),
				Start:  v.CreatedAt,
				End:    v.CreatedAt,
				Equips: equipGroup,
			}
		}


	}
}
