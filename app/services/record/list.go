package record

import (
	"github.com/gin-gonic/gin"
	"jd_workout_golang/app/middleware"
	"jd_workout_golang/app/repositories/pageinate"
	repo "jd_workout_golang/app/repositories/record"
)

type listResponse struct {
	Page    int      `json:"currentPage" form:"currentPage"`
	PerPage int      `json:"perPage" form:"perPage"`
	Data    []record `json:"data"`
	Total   int64    `json:"total"`
}

type record struct {
	ID     uint    `json:"id"`
	Weight float32 `json:"weight"`
	Reps   int     `json:"reps"`
	Sets   int     `json:"sets"`
	Note   string  `json:"note"`
	Equip  equip   `json:"equip"`
}

type equip struct {
	ID    uint    `json:"id"`
	Name  string  `json:"name"`
	Note  *string `json:"note"`
	Image *string `json:"image"`
}

// get record list
// @Summary record list
// @Description record list for personal user
// @Tags Record
// @Accept json
// @Produce json
// @Param recordList query recordListRequest true "recordList"
// @Success 200 {object} listResponse
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

	data, count, err := repo.GetRecordList(paginateCondition, middleware.Uid)

	if err != nil {
		c.JSON(422, gin.H{
			"message": "get record list error",
			"error":   err.Error(),
		})

		c.Abort()

		return
	}

	var convertedRecords []record
	for _, r := range *data {
		convertedRecord := record{
			ID:     r.ID,
			Weight: r.Weight,
			Reps:   int(r.Reps),
			Sets:   0,
			Note:   r.Note,
			Equip: equip{
				ID:    r.Equip.ID,
				Name:  r.Equip.Name,
				Note:  r.Equip.Note,
				Image: r.Equip.Image,
			},
		}

		convertedRecords = append(convertedRecords, convertedRecord)
	}

	c.JSON(200, listResponse{
		Page:    paginate.Page,
		PerPage: paginate.PerPage,
		Data:    convertedRecords,
		Total:   *count,
	})
}
