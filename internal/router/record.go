package router

import (
	"github.com/gin-gonic/gin"
	auth "jd_workout_golang/app/middleware"
	recordAction "jd_workout_golang/app/services/record"
)

func RegisterRecord(r *gin.RouterGroup) {
	equipGroup := r.Group("/record").Use(auth.ValidateToken)

	equipGroup.GET("/", recordAction.List)
	equipGroup.GET("/day-summary", recordAction.DaySummaryList)
	equipGroup.POST("/", recordAction.CreateRecord)
	equipGroup.PATCH("/:id", recordAction.UpdateRecord)
	equipGroup.DELETE("/:id", recordAction.DeleteRecord)
}
