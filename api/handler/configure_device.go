package handler

import (
	"fmt"
	"go-iot/api/common"
	"go-iot/api/model"

	"dario.cat/mergo"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func mergeStructs(config1 model.DeviceConfig, config2 model.DeviceConfig) (model.DeviceConfig, error) {
	// Merge config2 into config1, overriding fields from config1
	if err := mergo.Merge(&config1, config2, mergo.WithOverride); err != nil {
		fmt.Println("Error merging configs:", err)
		return config1, err
	}

	return config1, nil
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

	fmt.Printf("%+v\n", requestBody)

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

	newConfig, err := mergeStructs(updatedDevice.Config, requestBody.Config)
	if err != nil {
		fmt.Printf("Error merging structs, %d", err)
		common.InternalServerError(c)
	}
	updatedDevice.Config = newConfig

	_, err = db.Model(&updatedDevice).Where("public_id = ? AND status='PROVISIONED'", id).Update(&updatedDevice)
	if err != nil {
		fmt.Printf("Error saving to db, %d", err)
		common.InternalServerError(c)
	}

	c.IndentedJSON(200, updatedDevice)

}
