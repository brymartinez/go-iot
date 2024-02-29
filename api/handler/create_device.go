package handler

import (
	"encoding/json"
	"fmt"
	"go-iot/api/common"
	"go-iot/api/model"
	"go-iot/api/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Construct error message from validation errors
func constructErrorMessage(err error) string {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return "Unknown validation error"
	}

	errMsg := make(map[string]string)
	for _, e := range validationErrors {
		errMsg[e.Field()] = e.Tag()
	}

	jsonErrMsg, err := json.Marshal(errMsg)
	if err != nil {
		return "Error generating error message"
	}

	return string(jsonErrMsg)
}

func setDefaultValues(config *model.DeviceConfig) {
	if config.IsEnabled == nil {
		defaultValue := true
		config.IsEnabled = &defaultValue
	}
	if config.IsInteractive == nil {
		defaultValue := true
		config.IsInteractive = &defaultValue
	}
	if config.SendFrequency == nil {
		defaultValue := "5m"
		config.SendFrequency = &defaultValue
	}
	if config.Connection == nil {
		defaultValue := ""
		config.Connection = &defaultValue
	}
	if config.Version == nil {
		defaultValue := ""
		config.Version = &defaultValue
	}
}

func CreateDevice(c *gin.Context) {
	var requestBody model.CreateDeviceDTO

	if err := c.BindJSON(&requestBody); err != nil {
		errMsg := "Failed to process request"
		common.BadRequestError(c, errMsg)
		return
	}

	// Set default values for fields if they are not provided
	setDefaultValues(&requestBody.Config)

	validate := validator.New()
	if err := validate.Struct(requestBody); err != nil {
		// Construct error message
		errMsg := constructErrorMessage(err)
		common.BadRequestError(c, errMsg)
		return
	}

	db, err := common.ConnectToDB()
	if err != nil {
		fmt.Printf("Error connecting to db, %d", err)
		common.InternalServerError(c)
		return
	}

	var device model.Device

	// Check if serial_no already exists
	err = db.Model(&model.Device{}).Where("serial_no = ? AND status != 'DELETED'", requestBody.SerialNo).Select(&device)
	if err != nil {
		if err.Error() == "pg: no rows in result set" {
			// expected
			publicId := service.GenerateID(requestBody.Class)
			fmt.Printf("Public id is %s\n", publicId)
			device.PublicID = publicId
			device.Status = "PROVISIONED"
			device.Config = requestBody.Config
			device.Class = requestBody.Class
			device.SerialNo = requestBody.SerialNo
			device.Name = requestBody.Name
			device.CreatedAt = time.Now()
			device.UpdatedAt = time.Now()
			_, err = db.Model(&device).Insert(&device)
			if err != nil {
				fmt.Printf("Error saving to db, %d", err)
				common.InternalServerError(c)
				return
			}

			c.JSON(200, device)
		} else {
			fmt.Printf("Error getting device, %d", err)
			common.InternalServerError(c)
			return
		}
	} else {
		common.DeviceExistsError(c)
		return
	}

}
