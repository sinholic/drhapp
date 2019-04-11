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

// ScientificworkCreate is for creating a scientificwork for a user in our system
func ScientificworkCreate(c *gin.Context) {
	// Assign scientificwork var to model.Scientificwork so we can bind it with BindJSON
	var scientificwork model.Scientificwork

	// Binding JSON value to scientificwork
	if err := c.BindJSON(&scientificwork); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
			Status: http.StatusBadRequest,
			Message: model.Message{
				Error: err.Error(),
			},
		})
		return
	}

	// Validate with govalidator
	_, err := validator.ValidateStruct(scientificwork)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
			Status: http.StatusBadRequest,
			Message: model.Message{
				Error: err.Error(),
			},
		})
		return
	}

	// Check if the user already have a scientificwork in our system, then set nothing
	// Otherwise create a new row for the user
	if dbc := db.Create(&scientificwork); dbc.Error != nil {
		// Create failed, do something e.g. return, panic etc.
		// db.Create(&scientificwork)
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.JSONResults{
			Status: http.StatusInternalServerError,
			Message: model.Message{
				Error: dbc.Error.Error(),
			},
		})
		return
	}
	helper.Log.Println("Scientificwork for : " + string(scientificwork.UserID) + " got created")

	// Return the scientificwork that login
	c.JSON(http.StatusOK, model.JSONResults{
		Status: http.StatusOK,
		Result: scientificwork,
		Message: model.Message{
			Success: "Scientificwork " + config.SuccessCreate,
		},
	})
}

// ScientificworkUpdate to update user scientificwork in our database, we can also use one function in a place
// But it would be hassle to read and review the code later
func ScientificworkUpdate(c *gin.Context) {
	// Assign scientificwork var to model.Scientificwork so we can bind it with BindJSON
	var scientificwork model.Scientificwork

	// Check user ID is valid int
	if userid, err := strconv.Atoi(c.Param("userID")); err == nil {
		myresult := db.Where(model.Scientificwork{UserID: uint(userid)})
		if myresult.Error != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
				Status: http.StatusBadRequest,
				Message: model.Message{
					Error: err.Error(),
				},
			})
			return
		}

		// Binding JSON value to scientificwork
		if err := c.BindJSON(&scientificwork); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
				Status: http.StatusBadRequest,
				Message: model.Message{
					Error: err.Error(),
				},
			})
			return
		}

		// Validate with govalidator
		_, err := validator.ValidateStruct(scientificwork)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
				Status: http.StatusBadRequest,
				Message: model.Message{
					Error: err.Error(),
				},
			})
			return
		}

		db.Save(&scientificwork)

		// Scientificwork update success
		c.JSON(http.StatusOK, model.JSONResults{
			Status: http.StatusOK,
			Result: myresult.Value,
			Message: model.Message{
				Success: "Scientificwork " + config.SuccessUpdate,
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
