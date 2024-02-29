package handler

import (
	"fmt"
	"go-iot/api/common"
	"go-iot/api/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func mergeStructs(config1 model.DeviceConfig, config2 model.DeviceConfig) model.DeviceConfig {
	// Handle pointer fields individually to properly override
	if config2.IsEnabled != nil {
		config1.IsEnabled = config2.IsEnabled
	}
	if config2.IsInteractive != nil {
		config1.IsInteractive = config2.IsInteractive
	}
	if config2.Connection != nil {
		config1.Connection = config2.Connection
	}
	if config2.SendFrequency != nil {
		config1.SendFrequency = config2.SendFrequency
	}
	if config2.Version != nil {
		config1.Version = config2.Version
	}

	return config1
}

func ConfigureDevice(c *gin.Context) {
	id := c.Param("id")
	var requestBody model.ConfigureDeviceDTO

	if err := c.BindJSON(&requestBody); err != nil {
		errMsg := "Failed to process request"
		common.BadRequestError(c, errMsg)
		return
	}

	validate := validator.New()
	if err := validate.Struct(requestBody); err != nil {
		// Handle validation error
		errMsg := constructErrorMessage(err)
		fmt.Printf("Validation error: %s\n", errMsg)
		common.BadRequestError(c, errMsg)
		return
	}

	db, err := common.ConnectToDB()
	if err != nil {
		fmt.Printf("Error connecting to db, %d", err)
		common.InternalServerError(c)
		return
	}

	var updatedDevice model.Device
	err = db.Model(&model.Device{}).Where("public_id = ? AND status='PROVISIONED'", id).Select(&updatedDevice)
	if err != nil {
		fmt.Printf("Error getting device, %d", err)
		if err.Error() == "pg: no rows in result set" {
			common.NotFoundError(c)
			return
		} else {
			common.InternalServerError(c)
			return
		}
	}

	updatedDevice.Config = mergeStructs(updatedDevice.Config, requestBody.Config)
	updatedDevice.UpdatedAt = time.Now()

	_, err = db.Model(&updatedDevice).Where("public_id = ? AND status='PROVISIONED'", id).Update(&updatedDevice)
	if err != nil {
		fmt.Printf("Error saving to db, %d", err)
		common.InternalServerError(c)
		return
	}

	c.JSON(200, updatedDevice)

}
