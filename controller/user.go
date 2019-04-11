package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"agit.com/smartdashboard-backend/helper"
	"agit.com/smartdashboard-backend/model"
	validator "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	ldap "gopkg.in/ldap.v3"
)

// UserLoginLDAP is for login user with LDAP user and password
func UserLoginLDAP(c *gin.Context) {
	// Assign user var to model.User so we can bind it with BindJSON
	var u model.User

	// Bind our JSON with user model
	if err := c.BindJSON(&u); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
			Status: http.StatusBadRequest,
			Message: model.Message{
				Error: err.Error(),
			},
		})
		return
	}

	// Validate with govalidator
	_, err := validator.ValidateStruct(u)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
			Status: http.StatusBadRequest,
			Message: model.Message{
				Error: err.Error(),
			},
		})
		return
	}

	// It's a non changeable password user to bind connection to the LDAP
	bindusername := "hr_receiver"
	bindpassword := "Tsel2008"

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", "10.65.181.233", 389))
	if err != nil {
		helper.Log.Println(err)
	}
	// defer l.Close()

	// Reconnect with TLS, since we don't need TLS we can skip it
	// err = l.StartTLS(&tls.Config{InsecureSkipVerify: false})
	// if err != nil {
	// 	helper.Log.Println(err)
	// }

	// First bind with a read only user
	err = l.Bind(bindusername, bindpassword)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
			Status: http.StatusBadRequest,
			Message: model.Message{
				Error: err.Error(),
			},
		})
		return
	}

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		"dc=Telkomsel,dc=co,dc=id",
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		// We are looking the user by their samaccountname ex: (16008853 in Telkomsel)
		fmt.Sprintf("(samaccountname=%s)", u.Username),
		[]string{
			"givenname",
			"sn",
			"mail",
			// Right now we only need these data to be store in our db, when you uncomment the attributes below, the order will be changed
			// Make sure you test it before you add another attribute

			// "cn",
			// "dn",
			// "userprincipalname",
			// "samaccountname",
			// "uid",
			// "telephonenumber",
			// "mailnickname",
			// "physicaldeliveryofficename",
			// "initials",
			// "department",
			// "description",
			// "objectguid",
			// "l",
		},
		nil,
	)

	// Get entries from the searchRequest
	sr, err := l.Search(searchRequest)
	if err != nil {
		helper.Log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
			Status: http.StatusBadRequest,
			Message: model.Message{
				Error: err.Error(),
			},
		})
		return
	}

	// Check if the entries more than 1, because LDAP is an unique login domain
	if len(sr.Entries) != 1 {
		helper.Log.Println("User does not exist or too many entries returned")
		c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
			Status: http.StatusBadRequest,
			Message: model.Message{
				Error: "User does not exist or too many entries returned",
			},
		})
		return
	}

	userdn := sr.Entries[0].DN
	userAttributes := sr.Entries[0].Attributes

	// Bind as the user to verify their password
	err = l.Bind(userdn, u.Password)
	if err != nil {
		helper.Log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
			Status: http.StatusBadRequest,
			Message: model.Message{
				Error: err.Error(),
			},
		})
		return
	}

	// Rebind as the read only user for any further queries
	// err = l.Bind(bindusername, bindpassword)
	// if err != nil {
	// 	helper.Log.Println(err)
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
	// 		Status: http.StatusBadRequest,
	// 		Message: model.Message{
	// 			Error: err.Error(),
	// 		},
	// 	})
	// 	return
	// }

	// Map the user that login from the LDAP to our model
	user := model.User{
		Username:  u.Username,
		Name:      userAttributes[1].Values[0] + " " + userAttributes[0].Values[0],
		Email:     userAttributes[2].Values[0],
		Password:  helper.HashAndSalt([]byte(u.Password)),
		LastLogin: time.Now(),
		Role:      1, // We use default for new user role is User, we can change it on the fly with the api next time
	}

	// Create a token based on user login
	token, err := helper.CreateToken(user, time.Now().Format("2019-03-29"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
			Status: http.StatusBadRequest,
			Message: model.Message{
				Error: err.Error(),
			},
		})
		return
	}

	// Check if the user login already registered to our system then we set the login time
	// Otherwise create a new row for the user
	if dbc := db.Find(&user); dbc.Error != nil {
		// Find user failed then create a new user
		user.Token = token
		db.Create(&user)
		helper.Log.Println("User: " + u.Username + " got created")
	} else {
		db.Find(&user).Updates(model.User{LastLogin: time.Now(), Token: token})
		helper.Log.Println("User: " + u.Username + " login")
	}

	helper.Log.Println(u.Username)
	helper.Log.Println(user.Name)
	helper.Log.Println(u.Password)

	helper.Log.Println(time.Now().Format("2019-03-29"))

	// Return the user that login
	c.JSON(http.StatusOK, model.JSONResults{
		Status: http.StatusOK,
		Result: user,
		Message: model.Message{
			Success: "Login success",
		},
	})
}

// UserCreateAllDetail is a function to save user and all the detail with GORM
func UserCreateAllDetail(c *gin.Context) {
	var u model.User

	// Bind our JSON with user model
	if err := c.BindJSON(&u); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
			Status: http.StatusBadRequest,
			Message: model.Message{
				Error: err.Error(),
			},
		})
		return
	}

	// Validate with govalidator
	_, err := validator.ValidateStruct(u)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
			Status: http.StatusBadRequest,
			Message: model.Message{
				Error: err.Error(),
			},
		})
		return
	}

	user := model.User{
		LastLogin: time.Now(),
		Profile: model.Profile{
			FullName: u.Profile.FullName,
		},
	}

	helper.Log.Print(user)

}

// UserView to view user profile in our database
func UserView(c *gin.Context) {
	// Assign our variable for getting the profile
	var user model.User
	// var profile model.Profile

	// Check user ID is valid int
	if userid, err := strconv.Atoi(c.Param("userID")); err == nil {
		user.ID = uint(userid)
		myresult := db.Preload("Profile").Preload("Contact").Preload("Educations").Preload("Skills").Preload("Projects").Preload("Scientificworks").Preload("Languageproficiencies").Preload("Awards").Find(&user)
		if myresult.Error != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, model.JSONResults{
				Status: http.StatusBadRequest,
				Message: model.Message{
					Error: err.Error(),
				},
			})
			return
		}
		// Return the result of user profile
		c.JSON(http.StatusOK, model.JSONResults{
			Status:  http.StatusOK,
			Result:  myresult.Value,
			Message: model.Message{
				// Success: "Profile " + config.SuccessUpdate,
			},
		})
	} else {
		// the jokes ID is invalid
		c.AbortWithStatusJSON(http.StatusNotFound, model.JSONResults{
			Status: http.StatusNotFound,
			Message: model.Message{
				Error: http.ErrNotSupported.ErrorString,
			},
		})
	}
}
