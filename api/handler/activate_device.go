package handler

import (
	"fmt"
	"go-iot/api/common"
	"go-iot/api/model"

	"github.com/gin-gonic/gin"
)

func ActivateDevice(c *gin.Context) {
	id := c.Param("id")

	db, err := common.ConnectToDB()
	if err != nil {
		fmt.Printf("Error connecting to db, %d", err)
		common.InternalServerError(c)
		return
	}

	var device model.Device
	err = db.Model(&model.Device{}).Where("public_id = ? AND status='PROVISIONED'", id).Select(&device)
	if err != nil {
		fmt.Printf("Error getting device, %d", err)
		if err.Error() == "pg: no rows in result set" {
			common.NotFoundError(c)
		} else {
			common.InternalServerError(c)
		}
	}

	device.Status = "PENDING"
	// TODO - send to SNS
	_, err = db.Model(&device).Where("public_id = ? AND status='PROVISIONED'", id).Update(&device)
	if err != nil {
		fmt.Printf("Error saving to db, %d", err)
		common.InternalServerError(c)
	}

	c.JSON(200, device)
}
