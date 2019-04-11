package model

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

// All the struct models are defined here
// Utility models are below here
type (
	// JSONResults contains our struct to return with HTTP Code
	JSONResults struct {
		Status  int         `json:"status"`
		Result  interface{} `json:"result"`
		Message Message     `json:"message"`
	}
	// Message is a struct for JSONResults
	Message struct {
		Error   string `json:"error"`
		Success string `json:"success"`
	}
	// TokenClaims - Claims for a JWT access token.
	TokenClaims struct {
		ID       uint
		Username string
		Name     string
		Email    string
		jwt.StandardClaims
	}
)

// Automigration models are below here
// Make sure you put the validator with govalidator, for more validator type check out at https://github.com/asaskevich/govalidator
// Form gorm definition, you can check at http://doc.gorm.io/models.html#model-definition
type (
	// Revcode contains our reference for other table
	Revcode struct {
		gorm.Model
		RevCode int    `json:"revcode"`
		Type    string `json:"type"`
		Value   string `json:"value"`
	}

	// User contains information about a single User
	User struct {
		gorm.Model
		Username  string    `json:"username" valid:"required" gorm:"unique"`
		Name      string    `json:"name" valid:"-"`
		Email     string    `json:"email" gorm:"unique"`
		Password  string    `json:"password" valid:"required"`
		Token     string    `json:"token" valid:"-" gorm:"type:text"`
		LastLogin time.Time `json:"lastlogin" valid:"-"`
		Role      int       `json:"role" valid:"-"`
		// Relationship should be optional validation
		// Otherwise you would see the error field from the Relationship field that required
		Profile               Profile               `json:"profile" valid:"-"`
		Contact               Contact               `json:"contact" valid:"-"`
		Educations            []Education           `json:"educations" valid:"-"`
		Skills                []Skill               `json:"skills" valid:"-"`
		Projects              []Project             `json:"projects" valid:"-"`
		Scientificworks       []Scientificwork      `json:"scientificworks" valid:"-"`
		Languageproficiencies []Languageproficiency `json:"languageproficiencies" valid:"-"`
		Awards                []Award               `json:"awards" valid:"-"`
	}
	// Profile contains information about profile in a user
	Profile struct {
		gorm.Model
		UserID         uint      `json:"userid" valid:"required" gorm:"unique"`
		FullName       string    `json:"fullname" valid:"required"`
		NickName       string    `json:"nickname" valid:"required"`
		Gender         int       `json:"gender" valid:"required"`
		PlaceOfBirth   string    `json:"placeofbirth" valid:"required"`
		DateOfBirth    time.Time `json:"dateofbirth" valid:"required"`
		Religion       int       `json:"religion" valid:"required"`
		MaritalStatus  int       `json:"maritalstatus" valid:"required"`
		Children       int       `json:"children" valid:"required"`
		Education      int       `json:"education" valid:"required"`
		Height         int       `json:"height" valid:"required"`
		Weight         int       `json:"weight" valid:"required"`
		Nationality    int       `json:"nationality" valid:"required"`
		IdentityType   int       `json:"identitytype" valid:"required"`
		IdentityNumber string    `json:"identitynumber" valid:"required" gorm:"unique"`
		IdentityValid  time.Time `json:"identityvalid" valid:"required"`
	}
	// Contact contains information about user contact
	Contact struct {
		gorm.Model
		UserID      uint   `json:"userid" valid:"required" gorm:"unique"`
		Address     string `json:"address" valid:"required" gorm:"type:text" `
		MobilePhone string `json:"mobilephone" valid:"required"`
		Email       string `json:"email" valid:"required"`
		Facebook    string `json:"facebook" valid:"-"`
		Twitter     string `json:"twitter" valid:"-"`
		Instagram   string `json:"instagram" valid:"-"`
		// Socialmedias []Socialmedia `json:"socialmedias" valid:"-"`
	}

	// Socialmedia contains
	// Socialmedia struct {
	// 	gorm.Model
	// 	UserID uint   `json:"userid" valid:"-"`
	// 	Type   int    `json:"type" valid:"-"`
	// 	Name   string `json:"name" valid:"-"`
	// }

	// Education contains
	Education struct {
		gorm.Model
		UserID    uint    `json:"userid" valid:"required"`
		Type      int     `json:"type" valid:"required"`
		Level     string  `json:"level" valid:"required"`
		Position  string  `json:"position" valid:"-"`
		Name      string  `json:"name" valid:"required"`
		Organizer string  `json:"organizer" valid:"-"`
		Grade     float64 `json:"grade" valid:"required"`
		Major     string  `json:"major" valid:"required"`
		StartDate string  `json:"startdate" valid:"required"`
		EndDate   string  `json:"enddate" valid:"-"`
	}
	// Skill contains
	Skill struct {
		gorm.Model
		UserID      uint   `json:"userid" valid:"required"`
		Type        int    `json:"type" valid:"required"`
		Name        string `json:"name" valid:"required"`
		Proficiency int    `json:"proficiency" valid:"required"`
		Remark      string `json:"remark" valid:"-" gorm:"type:text"`
		Experience  int    `json:"experience" valid:"required"`
	}
	// Project contains
	Project struct {
		gorm.Model
		UserID           uint      `json:"userid" valid:"required"`
		Name             string    `json:"name" valid:"required"`
		StartDate        time.Time `json:"startdate" valid:"required"`
		EndDate          time.Time `json:"enddate" valid:"required"`
		Description      string    `json:"description" valid:"required" gorm:"type:text"`
		CompanyName      string    `json:"companyname" valid:"required"`
		ClientName       string    `json:"clientname" valid:"required"`
		Platform         string    `json:"platform" valid:"-"`
		Database         string    `json:"database" valid:"-"`
		DevelopmentTools string    `json:"developmenttools" valid:"-" gorm:"type:text"`
		Role             string    `json:"role" valid:"required"`
		JobDescription   string    `json:"jobdescription" valid:"required" gorm:"type:text"`
	}
	// Scientificwork contains
	Scientificwork struct {
		gorm.Model
		UserID      uint   `json:"userid" valid:"-"`
		Name        string `json:"name" valid:"-" gorm:"type:text"`
		Description string `json:"description" valid:"-"`
	}
	// Languageproficiency contains
	Languageproficiency struct {
		gorm.Model
		UserID      uint   `json:"userid" valid:"-"`
		Name        string `json:"name" valid:"-"`
		Proficiency int    `json:"proficiency" valid:"-"`
	}
	// Award contains
	Award struct {
		gorm.Model
		UserID uint   `json:"userid" valid:"-"`
		Name   string `json:"name" valid:"-" gorm:"type:text"`
	}
)
