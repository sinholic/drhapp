package controller

import (
	"log"
	"net/http"
	"strconv"

	"agit.com/smartdashboard-backend/config"
	"agit.com/smartdashboard-backend/helper"
	"agit.com/smartdashboard-backend/model"
	validator "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// ContactCreate is for creating a contact for a user in our system
func ContactCreate(c *gin.Context) {
	// Assign contact var to model.Contact so we can bind it with BindJSON
	var contact model.Contact

	// Binding JSON value to contact
	if err := c.BindJSON(&contact); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
			Status: http.StatusBadRequest,
			Message: model.Message{
				Error: err.Error(),
			},
		})
		return
	}

	// Validate with govalidator
	_, err := validator.ValidateStruct(contact)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
			Status: http.StatusBadRequest,
			Message: model.Message{
				Error: err.Error(),
			},
		})
		return
	}

	// Check if the user already have a contact in our system, then it would error
	// Otherwise create a new row for the user
	if dbc := db.Create(&contact); dbc.Error != nil {
		// Create failed, do something e.g. return, panic etc.
		// db.Create(&contact)
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.JSONResults{
			Status: http.StatusInternalServerError,
			Message: model.Message{
				Error: dbc.Error.Error(),
			},
		})
		return
	}
	helper.Log.Println("Contact for : " + string(contact.UserID) + " got created")

	// Return the contact that login
	c.JSON(http.StatusOK, model.JSONResults{
		Status: http.StatusOK,
		Result: contact,
		Message: model.Message{
			Success: "Contact " + config.SuccessCreate,
		},
	})
}

// ContactUpdate to update user contact in our database, we can also use one function in a place
// But it would be hassle to read and review the code later
func ContactUpdate(c *gin.Context) {
	// Assign contact var to model.Contact so we can bind it with BindJSON
	var contact model.Contact

	// Check user ID is valid int
	if userid, err := strconv.Atoi(c.Param("userID")); err == nil {
		myresult := db.Where(model.Contact{UserID: uint(userid)}).First(&contact)
		if myresult.Error != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
				Status: http.StatusBadRequest,
				Message: model.Message{
					Error: err.Error(),
				},
			})
			return
		}
		log.Println(myresult)

		// Binding JSON value to contact
		if err := c.BindJSON(&contact); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
				Status: http.StatusBadRequest,
				Message: model.Message{
					Error: err.Error(),
				},
			})
			return
		}

		log.Println(contact.Address)

		// Validate with govalidator
		_, err := validator.ValidateStruct(contact)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
				Status: http.StatusBadRequest,
				Message: model.Message{
					Error: err.Error(),
				},
			})
			return
		}
		// Save contact param with all value
		db.Save(&contact)

		c.JSON(http.StatusOK, model.JSONResults{
			Status: http.StatusOK,
			Result: myresult,
			Message: model.Message{
				Success: "Contact " + config.SuccessUpdate,
			},
		})
	} else {
		// The user ID is invalid
		c.AbortWithStatusJSON(http.StatusNotFound, model.JSONResults{
			Status: http.StatusNotFound,
			Message: model.Message{
				Error: "Invalid request",
			},
		})
	}
}
