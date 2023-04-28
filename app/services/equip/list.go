package equip

import (
	"jd_workout_golang/app/middleware"
	"jd_workout_golang/app/models"
	repo "jd_workout_golang/app/repositories/equip"
	recordRepo "jd_workout_golang/app/repositories/record"
	"time"

	"github.com/gin-gonic/gin"
)

type equipListRequest struct {
	Page    int ` json:"currentPage" form:"currentPage"`
	PerPage int ` json:"perPage" form:"perPage"`
}

type equipListResponse struct {
	Page    int            `json:"currentPage" form:"currentPage"`
	PerPage int            `json:"perPage" form:"perPage"`
	Data    []equipExpand `json:"data"`
	Total   int64          `json:"total"`
}

type equipExpand struct {
	models.Equip `json:"equip"`
	MaxWeightRecord maxWeightRecord `json:"maxWeightRecord"`
	MaxVolumeRecord lastMaxWeightRecord `json:"maxVolumeRecord"`
	// lastRecords []models.Record
}

type maxWeightRecord struct {
	ID uint `json:"id"`
	Weight float32 `json:"weight"`
	Reps uint `json:"reps"`
	DayVolumn float32 `json:"dayVolumn"`
}

type lastMaxWeightRecord struct {
	ID uint `json:"id"`
	MaxWeight float32 `json:"maxWeight"`
	MaxWeightReps uint `json:"maxWeightReps"`
	DayVolumn float32 `json:"dayVolumn"`
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

	ids := []uint{}
	for _, v := range *data {
		ids = append(ids, v.ID)
	}

	maxWeight := recordRepo.GetMaxRecord(ids, time.Now().Format("2006-01-02") + " 23:59:59")
	hash := map[uint]recordRepo.RecordWithVolumn{}
	for _, v := range *maxWeight {
		hash[v.EquipId] = v
	}

	lastMaxWeight := recordRepo.GetMaxRecord(ids, time.Now().Format("2006-01-02") + " 00:00:00")
	lastHash := map[uint]recordRepo.RecordWithVolumn{}
	for _, v := range *lastMaxWeight {
		lastHash[v.EquipId] = v
	}

	equipData := []equipExpand{}
	for _, v := range *data {
	
		equipData = append(equipData, equipExpand{
			Equip: v,
			MaxWeightRecord: maxWeightRecord{
				ID: v.ID,
				Weight: hash[v.ID].Weight,
				Reps: hash[v.ID].Reps,
				DayVolumn: float32(hash[v.ID].Volumn),
			},
			MaxVolumeRecord: lastMaxWeightRecord{
				ID: v.ID,
				MaxWeight: lastHash[v.ID].Weight,
				MaxWeightReps: lastHash[v.ID].Reps,
				DayVolumn: float32(lastHash[v.ID].Volumn),
			},
			// lastRecords: v.Record,
		})
	}
	

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
		Data:    equipData,
		Total:   *count,
	})
}
