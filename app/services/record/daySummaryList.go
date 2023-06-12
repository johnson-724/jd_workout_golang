package record

import (
	"fmt"
	"jd_workout_golang/app/middleware"
	fsRepo "jd_workout_golang/app/repositories/file"
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

type dateSummaryListResponse struct {
	Page    int         `json:"currentPage" form:"currentPage"`
	PerPage int         `json:"perPage" form:"perPage"`
	Data    []dateGroup `json:"data"`
	Total   int64       `json:"total"`
}

// get personal day summary list
// @Summary day summary list
// @Description record list for personal user
// @Tags Record
// @Accept x-www-form-urlencoded
// @Produce json
// @Param equipList query recordListRequest true "equipList"
// @Success 200 {object} dateSummaryListResponse
// @Failure 422 {string} string "{'message': '缺少必要欄位', 'error': 'error message'}"
// @Failure 403 {string} string "{'message': 'jwt token error', 'error': 'error message'}"
// @Router /record/day-summary [get]
// @Security Bearer
func DaySummaryList(c *gin.Context) {
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

	data, count, err := repo.GetDateSummaryRecords(paginateCondition, middleware.Uid)

	if err != nil {
		c.JSON(422, gin.H{
			"message": "get record list error",
			"error":   err.Error(),
		})

		c.Abort()

		return
	}

	groupByData := groupByRecord(*data)

	c.JSON(200, dateSummaryListResponse{
		Page:    paginate.Page,
		PerPage: paginate.PerPage,
		Data:    groupByData,
		Total:   *count,
	})
}

type dateGroup struct {
	Date   string       `json:"date"`
	Start  time.Time    `json:"start"`
	End    time.Time    `json:"end"`
	Equips []equipGroup `json:"equips"`
}

type equipGroup struct {
	ID      uint           `json:"id"`
	Name    string         `json:"name"`
	Note    *string        `json:"note"`
	Records []recordDetail `json:"records"`
	Image   *string        `json:"image"`
}

type recordDetail struct {
	IDS    []uint   `json:"ids"`
	Weight float32  `json:"weight"`
	Reps   int      `json:"reps"`
	Sets   int      `json:"sets"`
	Note   []string `json:"note"`
}

func groupByRecord(data []repo.RecordByDate) []dateGroup {
	group := make([]dateGroup, 0)

	currentDate := ""
	currentEquipId := 0
	currentRecordKey := ""
	for _, v := range data {
		if currentDate == "" || currentDate != v.Date {
			currentDate = v.Date
			currentEquipId = 0
			currentRecordKey = ""
			group = append(group, dateGroup{
				Date:  currentDate,
				Start: v.CreatedAt,
				End:   v.CreatedAt,
			})
		}

		if currentEquipId == 0 || currentEquipId != int(v.Record.EquipId) {
			currentEquipId = int(v.Record.EquipId)
			file := fsRepo.GinFileStore{
				Path: v.Equip.Image,
			}
			v.Equip.Image = file.GetPath()
			group[len(group)-1].Equips = append(group[len(group)-1].Equips, equipGroup{
				ID:      v.Record.EquipId,
				Name:    v.Record.Name,
				Note:    v.Equip.Note,
				Records: make([]recordDetail, 0),
				Image:   v.Equip.Image,
			})
		}

		key := fmt.Sprintf("%s_%d_%s_%s", currentDate, v.EquipId, strconv.FormatFloat(float64(v.Weight), 'f', 2, 64), strconv.Itoa(int(v.Reps)))
		if currentRecordKey == "" || currentRecordKey != key {
			currentRecordKey = key
			recordList := group[len(group)-1].Equips[len(group[len(group)-1].Equips)-1].Records
			record := recordDetail{
				IDS:    make([]uint, 0),
				Weight: v.Weight,
				Reps:   int(v.Reps),
				Sets:   0,
				Note:   make([]string, 0),
			}
			recordList = append(recordList, record)
			group[len(group)-1].Equips[len(group[len(group)-1].Equips)-1].Records = recordList
		}

		record := &group[len(group)-1].
			Equips[len(group[len(group)-1].Equips)-1].
			Records[len(group[len(group)-1].Equips[len(group[len(group)-1].Equips)-1].Records)-1]

		record.IDS = append(record.IDS, v.ID)
		record.Sets += 1
		if v.Note != "" {
			record.Note = append(record.Note, v.Note)
		}

	}

	return group
}
