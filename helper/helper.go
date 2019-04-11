package helper

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"agit.com/smartdashboard-backend/model"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	// Log define here
	Log *log.Logger
)

// InitLogFile is a function to create a log to file, we need this log when we deploy the application to binary file
func InitLogFile() {
	var logpath = flag.String("logpath", "/tmp/agit.com-smartdashboard-backend.log", "Log Path")
	writeLog(*logpath)
}

func writeLog(logpath string) {
	println("LogFile: " + logpath)
	file, err := os.Create(logpath)
	if err != nil {
		panic(err)
	}
	Log = log.New(file, "", log.LstdFlags|log.Llongfile)
}

// HashAndSalt is en encrypting function with bcrypt
// We can also change the hash function with another hash function (more info go to https://github.com/golang/crypto)
func HashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		Log.Println(err)
	}

	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

// Auth is a function to create and verify our token from the JSON
func Auth(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	// Parse token that got from Header to our SigningMethodH
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			Log.Fatalf("Unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(time.Now().Format("2019-03-29")), nil
	})

	// Verify the token
	if token != nil && err == nil {
		Log.Println("token verified")
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.JSONResults{
			Status: http.StatusUnauthorized,
			Message: model.Message{
				Error: err.Error(),
			},
		})
	}
}

// CreateToken - Create a JWT access token.
func CreateToken(user model.User, jwtKey string) (string, error) {

	expireToken := time.Now().Add(time.Hour * 48).Unix()

	// Set-up claims
	claims := model.TokenClaims{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
		Email:    user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "smartdashboard-backend-auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(jwtKey))

	return tokenString, err
}
