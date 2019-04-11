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

// ProfileCreate is for creating a profile for a user in our system
func ProfileCreate(c *gin.Context) {
	// Assign profile var to model.Profile so we can bind it with BindJSON
	var profile model.Profile

	// Binding JSON value to profile
	if err := c.BindJSON(&profile); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
			Status: http.StatusBadRequest,
			Message: model.Message{
				Error: err.Error(),
			},
		})
		return
	}

	// Validate with govalidator
	_, err := validator.ValidateStruct(profile)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
			Status: http.StatusBadRequest,
			Message: model.Message{
				Error: err.Error(),
			},
		})
		return
	}

	// Check if the user already have a profile in our system, then set nothing
	// Otherwise create a new row for the user
	if dbc := db.Create(&profile); dbc.Error != nil {
		// Create failed, do something e.g. return, panic etc.
		// db.Create(&profile)
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.JSONResults{
			Status: http.StatusInternalServerError,
			Message: model.Message{
				Error: dbc.Error.Error(),
			},
		})
		return
	}
	helper.Log.Println("Profile for : " + profile.FullName + " got created")

	// Return the profile that login
	c.JSON(http.StatusOK, model.JSONResults{
		Status: http.StatusOK,
		Result: profile,
		Message: model.Message{
			Success: "Profile " + config.SuccessCreate,
		},
	})
}

// ProfileUpdate to update user profile in our database, we can also use one function in a place
// But it would be hassle to read and review the code later
func ProfileUpdate(c *gin.Context) {
	// Assign profile var to model.Profile so we can bind it with BindJSON
	var profile model.Profile

	// Check user ID is valid int
	if userid, err := strconv.Atoi(c.Param("userID")); err == nil {
		myresult := db.Where(model.Profile{UserID: uint(userid)}).First(&profile)
		if myresult.Error != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
				Status: http.StatusBadRequest,
				Message: model.Message{
					Error: err.Error(),
				},
			})
			return
		}

		// Binding JSON value to profile
		if err := c.BindJSON(&profile); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
				Status: http.StatusBadRequest,
				Message: model.Message{
					Error: err.Error(),
				},
			})
			return
		}

		// Validate with govalidator
		_, err := validator.ValidateStruct(profile)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
				Status: http.StatusBadRequest,
				Message: model.Message{
					Error: err.Error(),
				},
			})
			return
		}
		// Save profile with all value
		db.Save(&profile)
		// Profile update success
		c.JSON(http.StatusOK, model.JSONResults{
			Status: http.StatusOK,
			Result: myresult.Value,
			Message: model.Message{
				Success: "Profile " + config.SuccessUpdate,
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
