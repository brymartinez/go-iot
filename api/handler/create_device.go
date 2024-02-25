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

func CreateDevice(c *gin.Context) {
	var requestBody model.CreateDeviceDTO

	if err := c.BindJSON(&requestBody); err != nil {
		// DO SOMETHING WITH THE ERROR
		fmt.Printf("Error parsing body, %d", err)
		common.InternalServerError(c)
		return
	}

	// Track whether fields were provided in the request body
	var isEnabledProvided, isInteractiveProvided, sendFrequencyProvided, connectionProvided, versionProvided bool

	// Check if fields were provided and set default values if necessary
	if !isEnabledProvided && requestBody.Config.IsEnabled == false {
		requestBody.Config.IsEnabled = true
	} // TODO - WRONG
	if !isInteractiveProvided && requestBody.Config.IsInteractive == false {
		requestBody.Config.IsInteractive = true
	} // TODO - WRONG
	if !sendFrequencyProvided && requestBody.Config.SendFrequency == "" {
		requestBody.Config.SendFrequency = "5m"
	}
	if !connectionProvided && requestBody.Config.Connection == "" {
		requestBody.Config.Connection = ""
	} // TODO - WRONG
	if !versionProvided && requestBody.Config.Version == "" {
		requestBody.Config.Version = ""
	} // TODO - WRONG

	validate := validator.New()
	if err := validate.Struct(requestBody); err != nil {
		// Construct error message
		errMsg := constructErrorMessage(err)
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

	var device model.Device
	publicId := service.IDGenerator(requestBody.Class)
	fmt.Printf("Public id is %s\n", publicId)
	device.PublicID = publicId
	device.Status = "PROVISIONED"
	device.Config = requestBody.Config
	device.Class = requestBody.Class
	device.Name = requestBody.Name
	device.CreatedAt = time.Now()
	device.UpdatedAt = time.Now()
	_, err = db.Model(&device).Insert(&device)
	if err != nil {
		fmt.Printf("Error saving to db, %d", err)
		common.InternalServerError(c)
		return
	}

	c.IndentedJSON(200, device)
}
