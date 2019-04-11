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

// LanguageproficiencyCreate is for creating a languageproficiency for a user in our system
func LanguageproficiencyCreate(c *gin.Context) {
	// Assign languageproficiency var to model.Languageproficiency so we can bind it with BindJSON
	var languageproficiency model.Languageproficiency

	// Binding JSON value to languageproficiency
	if err := c.BindJSON(&languageproficiency); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
			Status: http.StatusBadRequest,
			Message: model.Message{
				Error: err.Error(),
			},
		})
		return
	}

	// Validate with govalidator
	_, err := validator.ValidateStruct(languageproficiency)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
			Status: http.StatusBadRequest,
			Message: model.Message{
				Error: err.Error(),
			},
		})
		return
	}

	// Check if the user already have a languageproficiency in our system, then set nothing
	// Otherwise create a new row for the user
	if dbc := db.Create(&languageproficiency); dbc.Error != nil {
		// Create failed, do something e.g. return, panic etc.
		// db.Create(&languageproficiency)
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.JSONResults{
			Status: http.StatusInternalServerError,
			Message: model.Message{
				Error: dbc.Error.Error(),
			},
		})
		return
	}
	helper.Log.Println("Languageproficiency for : " + string(languageproficiency.UserID) + " got created")

	// Return the languageproficiency that login
	c.JSON(http.StatusOK, model.JSONResults{
		Status: http.StatusOK,
		Result: languageproficiency,
		Message: model.Message{
			Success: "Languageproficiency " + config.SuccessCreate,
		},
	})
}

// LanguageproficiencyUpdate to update user languageproficiency in our database, we can also use one function in a place
// But it would be hassle to read and review the code later
func LanguageproficiencyUpdate(c *gin.Context) {
	// Assign languageproficiency var to model.Languageproficiency so we can bind it with BindJSON
	var languageproficiency model.Languageproficiency

	// Check user ID is valid int
	if userid, err := strconv.Atoi(c.Param("userID")); err == nil {
		myresult := db.Where(model.Languageproficiency{UserID: uint(userid)})
		if myresult.Error != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
				Status: http.StatusBadRequest,
				Message: model.Message{
					Error: err.Error(),
				},
			})
			return
		}

		// Binding JSON value to languageproficiency
		if err := c.BindJSON(&languageproficiency); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
				Status: http.StatusBadRequest,
				Message: model.Message{
					Error: err.Error(),
				},
			})
			return
		}

		// Validate with govalidator
		_, err := validator.ValidateStruct(languageproficiency)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
				Status: http.StatusBadRequest,
				Message: model.Message{
					Error: err.Error(),
				},
			})
			return
		}

		db.Save(&languageproficiency)

		// Languageproficiency update success
		c.JSON(http.StatusOK, model.JSONResults{
			Status: http.StatusOK,
			Result: myresult.Value,
			Message: model.Message{
				Success: "Languageproficiency " + config.SuccessUpdate,
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
