package controller

import (
	"GeoApi/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateZoneInput struct {
	Vertex1 models.Position `json:"vertex1" gorm:"embedded"`
	Vertex2 models.Position `json:"vertex2" gorm:"embedded"`
	Vertex3 models.Position `json:"vertex3" gorm:"embedded"`
	Vertex4 models.Position `json:"vertex4" gorm:"embedded"`
}

func CreateZone(db *gorm.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		var input CreateZoneInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		zone := models.Zone{Vertex1: input.Vertex1, Vertex2: input.Vertex2,
			Vertex3: input.Vertex3, Vertex4: input.Vertex4}
		db.Create(&zone)

		c.JSON(http.StatusOK, gin.H{"data": zone})
	}
	return gin.HandlerFunc(fn)
}

func GetAllZone(db *gorm.DB) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		var zones []models.Zone
		db.Find(&zones)
		ctx.JSON(http.StatusOK, gin.H{"data": zones})
	}
	return gin.HandlerFunc(fn)
}

func FindZoneById(db *gorm.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var zone models.Zone

		if err := db.Where("id = ?", c.Param("id")).First(&zone).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": zone})
	}
	return gin.HandlerFunc(fn)
}
func DeleteZoneById(db *gorm.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var zone models.Zone
		if err := db.Where("id = ?", c.Param("id")).First(&zone).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
			return
		}

		db.Delete(&zone)

		c.JSON(http.StatusOK, gin.H{"data": true})
	}
	return gin.HandlerFunc(fn)
}

type DeviceInput struct {
	Position  models.Position `json:"position" gorm:"embedded"`
	TimeStamp string          `json:"timestamp"`
}

func CreateDevice(db *gorm.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		// fmt.Println(c.Request)

		var input DeviceInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		device := models.Device{Position: input.Position}
		device.TimeStamp = time.Now().Format(time.RFC850)
		db.Create(&device)

		c.JSON(http.StatusOK, gin.H{"data": device, "id": device.ID})
	}
	return gin.HandlerFunc(fn)
}

func GetAllDevice(db *gorm.DB) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		var devices []models.Device
		db.Find(&devices)
		ctx.JSON(http.StatusOK, gin.H{"data": devices})
	}
	return gin.HandlerFunc(fn)
}

func FindDeviceById(db *gorm.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var device models.Device

		if err := db.Where("id = ?", c.Param("id")).First(&device).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": device})
	}
	return gin.HandlerFunc(fn)
}

func CrossProduct(EV models.Position, PV models.Position) float64 {
	return EV.Latitude*PV.Longitude - EV.Longitude*PV.Latitude
}

func MakeVector(V1 models.Position, V2 models.Position) models.Position {
	v := models.Position{Latitude: V2.Latitude - V1.Latitude, Longitude: V2.Longitude - V1.Longitude}
	return v
}

// EV[i] = V[i+1] - V[i], where V[] - vertices in order
// PV[i] = P - V[i]
// Cross[i] = CrossProduct(EV[i], PV[i]) = EV[i].X * PV[i].Y - EV[i].Y * PV[i].X
func IsInZone(devicePos models.Position, verts []models.Position) bool {
	size := 4

	var ev []models.Position
	for i := 0; i < size-1; i++ {
		ev = append(ev, MakeVector(verts[i], verts[i+1]))
	}

	var pv []models.Position
	for i := 0; i < size; i++ {
		pv = append(pv, MakeVector(verts[i], devicePos))
	}

	var cross []float64
	for i := 0; i < size-1; i++ {
		cross = append(cross, CrossProduct(ev[i], pv[i]))
	}

	for i := 0; i < size-1; i++ {
		// fmt.Println(cross[i])
		if cross[i] > 0 {
			return false
		}
	}

	return true
}

func ExpandPositionNumber(pos models.Position) models.Position {
	scale := 1000000.0
	return models.Position{Latitude: pos.Latitude * scale, Longitude: pos.Longitude * scale}
}

func CheckDeviceInZone(device models.Device, db *gorm.DB) bool {
	var zones []models.Zone
	db.Find(&zones)
	var is_in_zone bool = false

	for _, z := range zones {
		vertexes := []models.Position{
			ExpandPositionNumber(z.Vertex1),
			ExpandPositionNumber(z.Vertex2),
			ExpandPositionNumber(z.Vertex3),
			ExpandPositionNumber(z.Vertex4)}
		is_in_zone = IsInZone(ExpandPositionNumber(device.Position), vertexes)
		if is_in_zone {
			break
		}
		// fmt.Println(is_in_zone)
	}
	return is_in_zone
}

func UpdateDevice(db *gorm.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var device models.Device
		if err := db.Where("id = ?", c.Param("id")).First(&device).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
			return
		}

		var input DeviceInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		input.TimeStamp = time.Now().Format(time.RFC850)

		db.Model(&device).Updates(input)

		is_in_zone := CheckDeviceInZone(device, db)

		device.IsInZone = is_in_zone
		db.Save(&device)

		c.JSON(http.StatusOK, gin.H{"data": device, "is_in_zone": is_in_zone})
	}
	return gin.HandlerFunc(fn)
}

func DeleteDeviceById(db *gorm.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var device models.Device
		if err := db.Where("id = ?", c.Param("id")).First(&device).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
			return
		}

		db.Delete(&device)

		c.JSON(http.StatusOK, gin.H{"data": true})
	}
	return gin.HandlerFunc(fn)
}
