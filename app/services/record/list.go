package record

import (
	"jd_workout_golang/app/middleware"
	"strconv"
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
	Page    int                    `json:"currentPage" form:"currentPage"`
	PerPage int                    `json:"perPage" form:"perPage"`
	Data    map[string]*dateGroup `json:"data"`
	Total   int64                  `json:"total"`
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
	ptr := &groupByData
	group := *ptr

	c.JSON(200, recordListResponse{
		Page:    paginate.Page,
		PerPage: paginate.PerPage,
		Data:    group,
		Total:   *count,
	})
}

type dateGroup struct {
	Date   string `json:"date"`
	Start  time.Time `json:"start`
	End    time.Time `json:"end`
	Equips map[int]equipGroup `json:"equips"`
}

type equipGroup struct {
	ID      uint `json:"id"`
	Name    string `json:"name"`
	Note    string `json:"note"`
	Records map[string]recordDetail `json:"records"`
}

type recordDetail struct {
	ID     uint `json:"id"`
	Weight float64 `json:"weight"`
	Reps   int `json:"reps"`
	Sets   float64 `json:sets` // Weight * Reps * Count
	Note   []string `json:"note"`
}

func groupBy(data []repo.RecordByDate) map[string]*dateGroup {
	group := make(map[string]*dateGroup)
	for _, v := range data {
		if _, ok := group[v.Date]; !ok {
			equipGroupMap := make(map[int]equipGroup)
			group[v.Date] = &dateGroup{
				Date:   v.Date,
				Start:  v.CreatedAt,
				End:    v.CreatedAt,
				Equips: equipGroupMap,
			}
		}

		recordGroupMap := make(map[string]recordDetail)

		if _, ok := (*group[v.Date]).Equips[int(v.EquipId)]; !ok {
			(*group[v.Date]).Equips[int(v.EquipId)] = equipGroup{
				ID:      v.EquipId,
				Name:    v.Name,
				Note:    v.Note,
				Records: recordGroupMap,
			}
		}

		setsKey := strconv.FormatFloat(float64(v.Weight), 'E', -1, 64) + "-" + strconv.Itoa(int(v.Reps))

		(*group[v.Date]).Start = v.CreatedAt

		if _, ok := (*group[v.Date]).Equips[int(v.EquipId)].Records[setsKey]; !ok {
			(*group[v.Date]).Equips[int(v.EquipId)].Records[setsKey] = recordDetail{
				ID:     v.ID,
				Weight: float64(v.Weight),
				Reps:   int(v.Reps),
				Sets:   0,
				Note:   make([]string, 0),
			}
		}

		record := (*group[v.Date]).Equips[int(v.EquipId)].Records[setsKey]
		record.Sets+=1
		record.Note = append(record.Note, v.Note)
		(*group[v.Date]).Equips[int(v.EquipId)].Records[setsKey] = record
	}

	return group
}
