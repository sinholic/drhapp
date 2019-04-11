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

// SkillCreate is for creating a skill for a user in our system
func SkillCreate(c *gin.Context) {
	// Assign skill var to model.Skill so we can bind it with BindJSON
	var skill model.Skill

	// Binding JSON value to skill
	if err := c.BindJSON(&skill); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
			Status: http.StatusBadRequest,
			Message: model.Message{
				Error: err.Error(),
			},
		})
		return
	}

	// Validate with govalidator
	_, err := validator.ValidateStruct(skill)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
			Status: http.StatusBadRequest,
			Message: model.Message{
				Error: err.Error(),
			},
		})
		return
	}

	// Check if the user already have a skill in our system, then set nothing
	// Otherwise create a new row for the user
	if dbc := db.Create(&skill); dbc.Error != nil {
		// Create failed, do something e.g. return, panic etc.
		// db.Create(&skill)
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.JSONResults{
			Status: http.StatusInternalServerError,
			Message: model.Message{
				Error: dbc.Error.Error(),
			},
		})
		return
	}
	helper.Log.Println("Skill for : " + string(skill.UserID) + " got created")

	// Return the skill that login
	c.JSON(http.StatusOK, model.JSONResults{
		Status: http.StatusOK,
		Result: skill,
		Message: model.Message{
			Success: "Skill " + config.SuccessCreate,
		},
	})
}

// SkillUpdate to update user skill in our database, we can also use one function in a place
// But it would be hassle to read and review the code later
func SkillUpdate(c *gin.Context) {
	// Assign skill var to model.Skill so we can bind it with BindJSON
	var skill model.Skill

	// Check user ID is valid int
	if userid, err := strconv.Atoi(c.Param("userID")); err == nil {
		myresult := db.Where(model.Skill{UserID: uint(userid)})
		if myresult.Error != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
				Status: http.StatusBadRequest,
				Message: model.Message{
					Error: err.Error(),
				},
			})
			return
		}

		// Binding JSON value to skill
		if err := c.BindJSON(&skill); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
				Status: http.StatusBadRequest,
				Message: model.Message{
					Error: err.Error(),
				},
			})
			return
		}

		// Validate with govalidator
		_, err := validator.ValidateStruct(skill)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
				Status: http.StatusBadRequest,
				Message: model.Message{
					Error: err.Error(),
				},
			})
			return
		}

		db.Save(&skill)

		// Skill update success
		c.JSON(http.StatusOK, model.JSONResults{
			Status: http.StatusOK,
			Result: myresult.Value,
			Message: model.Message{
				Success: "Skill " + config.SuccessUpdate,
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
