package handler

import (
	"fmt"
	"go-iot/api/common"
	"go-iot/api/model"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

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
	err = db.Model(&model.Device{}).Where("public_id = ?", id).Select(&updatedDevice)
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

	c.IndentedJSON(200, updatedDevice)

}
