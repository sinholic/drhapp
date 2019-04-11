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

// ProjectCreate is for creating a project for a user in our system
func ProjectCreate(c *gin.Context) {
	// Assign project var to model.Project so we can bind it with BindJSON
	var project model.Project

	// Binding JSON value to project
	if err := c.BindJSON(&project); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
			Status: http.StatusBadRequest,
			Message: model.Message{
				Error: err.Error(),
			},
		})
		return
	}

	// Validate with govalidator
	_, err := validator.ValidateStruct(project)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
			Status: http.StatusBadRequest,
			Message: model.Message{
				Error: err.Error(),
			},
		})
		return
	}

	// Check if the user already have a project in our system, then set nothing
	// Otherwise create a new row for the user
	if dbc := db.Create(&project); dbc.Error != nil {
		// Create failed, do something e.g. return, panic etc.
		// db.Create(&project)
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.JSONResults{
			Status: http.StatusInternalServerError,
			Message: model.Message{
				Error: dbc.Error.Error(),
			},
		})
		return
	}
	helper.Log.Println("Project for : " + string(project.UserID) + " got created")

	// Return the project that login
	c.JSON(http.StatusOK, model.JSONResults{
		Status: http.StatusOK,
		Result: project,
		Message: model.Message{
			Success: "Project " + config.SuccessCreate,
		},
	})
}

// ProjectUpdate to update user project in our database, we can also use one function in a place
// But it would be hassle to read and review the code later
func ProjectUpdate(c *gin.Context) {
	// Assign project var to model.Project so we can bind it with BindJSON
	var project model.Project

	// Check user ID is valid int
	if userid, err := strconv.Atoi(c.Param("userID")); err == nil {
		myresult := db.Where(model.Project{UserID: uint(userid)})
		if myresult.Error != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
				Status: http.StatusBadRequest,
				Message: model.Message{
					Error: err.Error(),
				},
			})
			return
		}

		// Binding JSON value to project
		if err := c.BindJSON(&project); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
				Status: http.StatusBadRequest,
				Message: model.Message{
					Error: err.Error(),
				},
			})
			return
		}

		// Validate with govalidator
		_, err := validator.ValidateStruct(project)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
				Status: http.StatusBadRequest,
				Message: model.Message{
					Error: err.Error(),
				},
			})
			return
		}

		db.Save(&project)

		// Project update success
		c.JSON(http.StatusOK, model.JSONResults{
			Status: http.StatusOK,
			Result: myresult.Value,
			Message: model.Message{
				Success: "Project " + config.SuccessUpdate,
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
