package record

import (
	"jd_workout_golang/app/middleware"
	// "strconv"
	// "jd_workout_golang/app/models"
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
	Data    map[string]interface{} `json:"data"`
	Total   int64           `json:"total"`
}

// get personal record list
// @Summary record list
// @Description record list for personal user
// @Tags Record
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

	groupByData := groupBy(*data)

	c.JSON(200, recordListResponse{
		Page:    paginate.Page,
		PerPage: paginate.PerPage,
		Data:    *groupByData,
		Total:   *count,
	})
}

type dateGroup struct {
	Date   string
	Start  time.Time
	End    time.Time
	Equips map[int]equipGroup
}

type equipGroup struct {
	ID      uint
	Name    string
	Note    string
	Records map[int]recordDetail
}

type recordDetail struct {
	ID     uint
	Weight float64
	Reps   int
	Sets   float64 // Weight * Reps * Count
	Note   []string
}

func groupBy(data []repo.RecordByDate) *map[string]interface{} {
	group := make(map[string]interface{})
	for _, v := range data {
		if _, ok := group[v.Date]; !ok {
			equipGroupMap := make(map[int]equipGroup)
			// recordGroup := make(map[string]recordDetail)
			group[v.Date] = &dateGroup{
				Date:  v.Date,
				Start: v.CreatedAt,
				End:   v.CreatedAt,
				Equips: equipGroupMap,
			}
		}

		if _, ok := group[v.Date].(*dateGroup).Equips[int(v.EquipId)]; !ok {
			recordGroupMap := make(map[int]recordDetail)
			group[v.Date].(*dateGroup).Equips[int(v.EquipId)] = equipGroup{
				ID:      v.EquipId,
				Name:    v.Name,
				Note:    v.Note,
				Records: recordGroupMap,
			}
		}

		group[v.Date].(*dateGroup).Start = v.CreatedAt
		group[v.Date].(*dateGroup).Equips[int(v.EquipId)].Records[int(v.ID)] = recordDetail{
			ID:     v.ID,
			Weight: float64(v.Weight),
			Reps:   int(v.Reps),
			Sets:   float64(v.Weight) * float64(v.Reps),
			Note:   []string{v.Note},
		}

	}

	return &group
}
