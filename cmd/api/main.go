package main

import (
	"GeoApi/controller"
	"GeoApi/database"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	DB := database.ConnectDatabase()

	r.POST("/zones", controller.CreateZone(DB))
	r.GET("/zones", controller.GetAllZone(DB))
	r.GET("/zones/:id", controller.FindZoneById(DB))
	r.GET("/zones/delete/:id", controller.DeleteZoneById(DB))

	r.POST("/devices", controller.CreateDevice(DB))
	r.GET("/devices", controller.GetAllDevice(DB))
	r.GET("/devices/:id", controller.FindDeviceById(DB))
	r.POST("/devices/:id", controller.UpdateDevice(DB))
	r.GET("/devices/delete/:id", controller.DeleteDeviceById(DB))

	r.Run()
}
