package record

import (
	"jd_workout_golang/app/middleware"
	recordRepo "jd_workout_golang/app/repositories/record"
	"strconv"
	"github.com/gin-gonic/gin"
)

// delete record
// @Summary delete record
// @Description delete record
// @Tags Record
// @Accept json
// @Produce json
// @Param id path string true "record id"
// @Success 200 {string} string "{'message': 'delete success'}"
// @Failure 422 {string} string "{'message': '缺少必要欄位', 'error': 'error message'}"
// @Failure 403 {string} string "{'message': 'jwt token error', 'error': 'error message'}"
// @Router /record/{id} [delete]
// @Security Bearer
func DeleteRecord(c *gin.Context) {
	id := c.Param("id")

	recordId, err := strconv.ParseUint(id, 10, 32)

	if err != nil {
		c.JSON(422, gin.H{
			"message": "uri id error",
		})

		c.Abort()

		return
	}

	record, err := recordRepo.GetRecord(uint64(recordId), middleware.Uid)

	if err != nil {
		c.JSON(422, gin.H{
			"message": "record not found",
			"error":   err.Error(),
		})

		return
	}

	recordRepo.Delete(record)

	if err != nil {
		c.JSON(422, gin.H{
			"message": "delete error",
			"error":   err.Error(),
		})

		c.Abort()

		return
	}

	c.JSON(200, gin.H{
		"message": "delete success",
	})
}
