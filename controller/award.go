package controller

import (
	"net/http"
	"strconv"

	"agit.com/smartdashboard-backend/config"
	"agit.com/smartdashboard-backend/helper"
	"agit.com/smartdashboard-backend/model"
	validator "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// AwardCreate is for creating a award for a user in our system
func AwardCreate(c *gin.Context) {
	// Assign award var to model.Award so we can bind it with BindJSON
	var award model.Award

	// Binding JSON value to award
	if err := c.BindJSON(&award); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
			Status: http.StatusBadRequest,
			Message: model.Message{
				Error: err.Error(),
			},
		})
		return
	}

	// Validate with govalidator
	_, err := validator.ValidateStruct(award)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
			Status: http.StatusBadRequest,
			Message: model.Message{
				Error: err.Error(),
			},
		})
		return
	}

	// Check if the user already have a award in our system, then set nothing
	// Otherwise create a new row for the user
	if dbc := db.Create(&award); dbc.Error != nil {
		// Create failed, do something e.g. return, panic etc.
		// db.Create(&award)
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.JSONResults{
			Status: http.StatusInternalServerError,
			Message: model.Message{
				Error: dbc.Error.Error(),
			},
		})
		return
	}
	helper.Log.Println("Award for : " + award.Name + " got created")

	// Return the award that login
	c.JSON(http.StatusOK, model.JSONResults{
		Status: http.StatusOK,
		Result: award,
		Message: model.Message{
			Success: "Award " + config.SuccessCreate,
		},
	})
}

// AwardUpdate to update user award in our database, we can also use one function in a place
// But it would be hassle to read and review the code later
func AwardUpdate(c *gin.Context) {
	// Assign award var to model.Award so we can bind it with BindJSON
	var award model.Award

	// Check user ID is valid int
	if userid, err := strconv.Atoi(c.Param("userID")); err == nil {
		myresult := db.Where(model.Award{UserID: uint(userid)}).First(&award)
		if myresult.Error != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
				Status: http.StatusBadRequest,
				Message: model.Message{
					Error: err.Error(),
				},
			})
			return
		}

		// Binding JSON value to award
		if err := c.BindJSON(&award); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
				Status: http.StatusBadRequest,
				Message: model.Message{
					Error: err.Error(),
				},
			})
			return
		}

		// Validate with govalidator
		_, err := validator.ValidateStruct(award)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
				Status: http.StatusBadRequest,
				Message: model.Message{
					Error: err.Error(),
				},
			})
			return
		}
		// Save award with all value
		db.Save(&award)
		// Award update success
		c.JSON(http.StatusOK, model.JSONResults{
			Status: http.StatusOK,
			Result: myresult.Value,
			Message: model.Message{
				Success: "Award " + config.SuccessUpdate,
			},
		})
	} else {
		// The user ID is invalid
		c.AbortWithStatusJSON(http.StatusNotFound, model.JSONResults{
			Status: http.StatusNotFound,
			Message: model.Message{
				Error: http.ErrNotSupported.ErrorString,
			},
		})
	}
}
