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

// EducationCreate is for creating a education for a user in our system
func EducationCreate(c *gin.Context) {
	// Assign education var to model.Education so we can bind it with BindJSON
	var education model.Education

	// Binding JSON value to education
	if err := c.BindJSON(&education); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
			Status: http.StatusBadRequest,
			Message: model.Message{
				Error: err.Error(),
			},
		})
		return
	}

	// Validate with govalidator
	_, err := validator.ValidateStruct(education)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
			Status: http.StatusBadRequest,
			Message: model.Message{
				Error: err.Error(),
			},
		})
		return
	}

	// Check if the user already have a education in our system, then set nothing
	// Otherwise create a new row for the user
	if dbc := db.Create(&education); dbc.Error != nil {
		// Create failed, do something e.g. return, panic etc.
		// db.Create(&education)
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.JSONResults{
			Status: http.StatusInternalServerError,
			Message: model.Message{
				Error: dbc.Error.Error(),
			},
		})
		return
	}
	helper.Log.Println("Education for : " + string(education.UserID) + " got created")

	// Return the education that login
	c.JSON(http.StatusOK, model.JSONResults{
		Status: http.StatusOK,
		Result: education,
		Message: model.Message{
			Success: "Education " + config.SuccessCreate,
		},
	})
}

// EducationUpdate to update user education in our database, we can also use one function in a place
// But it would be hassle to read and review the code later
func EducationUpdate(c *gin.Context) {
	// Assign education var to model.Education so we can bind it with BindJSON
	var education model.Education

	// Check user ID is valid int
	if userid, err := strconv.Atoi(c.Param("userID")); err == nil {
		myresult := db.Where(model.Education{UserID: uint(userid)})
		if myresult.Error != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
				Status: http.StatusBadRequest,
				Message: model.Message{
					Error: err.Error(),
				},
			})
			return
		}

		// Binding JSON value to education
		if err := c.BindJSON(&education); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
				Status: http.StatusBadRequest,
				Message: model.Message{
					Error: err.Error(),
				},
			})
			return
		}

		// Validate with govalidator
		_, err := validator.ValidateStruct(education)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
				Status: http.StatusBadRequest,
				Message: model.Message{
					Error: err.Error(),
				},
			})
			return
		}

		db.Save(&education)

		// Education update success
		c.JSON(http.StatusOK, model.JSONResults{
			Status: http.StatusOK,
			Result: myresult.Value,
			Message: model.Message{
				Success: "Education " + config.SuccessUpdate,
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
