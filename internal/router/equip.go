package router

import (
	"github.com/gin-gonic/gin"
	auth "jd_workout_golang/app/middleware"
	equipAction "jd_workout_golang/app/services/equip"
)

func RegisterEquip(r *gin.RouterGroup) {
	equipGroup := r.Group("/equip").Use(auth.ValidateToken)

	equipGroup.GET("/", equipAction.List)
	equipGroup.POST("/", equipAction.CreateEquip)
	equipGroup.PUT("/:id/weight", equipAction.UpdateWeight)
	equipGroup.PATCH("/:id", equipAction.UpdateEquip)
	equipGroup.DELETE("/:id", equipAction.DeleteEquip)
}
