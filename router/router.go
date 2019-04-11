package router

import (
	"agit.com/smartdashboard-backend/controller"

	"agit.com/smartdashboard-backend/helper"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

// Lists is a list of router in smartdashboard-backend
func Lists() *gin.Engine {
	// Set the router as the default one shipped with Gin
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// Serve frontend static files
	router.Use(static.Serve("/", static.LocalFile("./views", true)))

	// Setup route group for the API
	apiRouter := router.Group("/api")
	{
		// apiRouter.Use(gzip.Gzip(gzip.DefaultCompression))
		// We use an authorization from SigningMethodHS256 and we need that token every request from the API
		// Group for User
		userRouter := apiRouter.Group("/user")
		{
			userRouter.POST("/login", controller.UserLoginLDAP)
			userRouter.POST("/view/:userID", helper.Auth, controller.UserView)
		}
		// Group for Profile
		profileRouter := apiRouter.Group("/profile")
		{
			profileRouter.POST("/create", helper.Auth, controller.ProfileCreate)
			profileRouter.POST("/update/:profileID", helper.Auth, controller.ProfileUpdate)
		}
		// Group for Contact
		contactRouter := apiRouter.Group("/contact")
		{
			contactRouter.POST("/create", helper.Auth, controller.ContactCreate)
			contactRouter.POST("/update/:contactID", helper.Auth, controller.ContactUpdate)
			// contactRouter.POST("/view/:userID", helper.Auth, contact.ViewProfile)
		}
		// Group for Education
		educationRouter := apiRouter.Group("/education")
		{
			educationRouter.POST("/create", helper.Auth, controller.EducationCreate)
			educationRouter.POST("/update/:educationID", helper.Auth, controller.EducationUpdate)
			// contactRouter.POST("/view/:userID", helper.Auth, contact.ViewProfile)
		}
		// Group for Skill
		skillRouter := apiRouter.Group("/skill")
		{
			skillRouter.POST("/create", helper.Auth, controller.SkillCreate)
			skillRouter.POST("/update/:skillID", helper.Auth, controller.SkillUpdate)
			// contactRouter.POST("/view/:userID", helper.Auth, contact.ViewProfile)
		}
		// Group for Project
		projectRouter := apiRouter.Group("/project")
		{
			projectRouter.POST("/create", helper.Auth, controller.ProjectCreate)
			projectRouter.POST("/update/:projectID", helper.Auth, controller.ProjectUpdate)
			// contactRouter.POST("/view/:userID", helper.Auth, contact.ViewProfile)
		}
		// Group for Scientificwork
		scientificworkRouter := apiRouter.Group("/scientificwork")
		{
			scientificworkRouter.POST("/create", helper.Auth, controller.ScientificworkCreate)
			scientificworkRouter.POST("/update/:scientificworkID", helper.Auth, controller.ScientificworkUpdate)
			// contactRouter.POST("/view/:userID", helper.Auth, contact.ViewProfile)
		}
		// Group for Languageproficiency
		languageproficiencyRouter := apiRouter.Group("/languageproficiency")
		{
			languageproficiencyRouter.POST("/create", helper.Auth, controller.LanguageproficiencyCreate)
			languageproficiencyRouter.POST("/update/:languageproficiencyID", helper.Auth, controller.LanguageproficiencyUpdate)
			// contactRouter.POST("/view/:userID", helper.Auth, contact.ViewProfile)
		}
		// Group for Award
		awardRouter := apiRouter.Group("/award")
		{
			awardRouter.POST("/create", helper.Auth, controller.AwardCreate)
			awardRouter.POST("/update/:awardID", helper.Auth, controller.AwardUpdate)
			// contactRouter.POST("/view/:userID", helper.Auth, contact.ViewProfile)
		}
	}

	return router
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Accept-Encoding", "gzip")
		c.Writer.Header().Set("Content-Type", "application/json")
		if c.Request.Method == "OPTIONS" {
			c.Abort()
			return
		}
		c.Next()
	}
}
