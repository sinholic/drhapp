package config

import (
	"agit.com/smartdashboard-backend/helper"
	"agit.com/smartdashboard-backend/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	db    *gorm.DB
	datas [][]string
	// indexRevcode int
)

// LoadDB is function to load our DB
func LoadDB() *gorm.DB {
	// Specify the database connection here, for more examples of MySQL connection you can check here https://github.com/go-sql-driver/mysql#examples
	db, err := gorm.Open("mysql", "root:mekise@tcp/smartdashboard?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		helper.Log.Fatal("error in connecting to database")
	}
	return db
}

// MigrateDB is a function to migrating our database
func MigrateDB() {
	// Load our DB and assign to the db
	db := LoadDB()
	// Migrate our models, make sure if you create a new model you put the model in here so gorm will create the table automatically
	db.AutoMigrate(
		&model.Revcode{},
		&model.User{},
		&model.Profile{},
		&model.Contact{},
		// &model.Socialmedia{},
		&model.Education{},
		&model.Skill{},
		&model.Project{},
		&model.Scientificwork{},
		&model.Languageproficiency{},
		&model.Award{},
	)
}

// PopulateDataInRev is create some reference data in our refcode table
func PopulateDataInRev() {
	db := LoadDB()

	datas = [][]string{
		// Add reference for Gender
		{"GENDER", "Male"},
		{"GENDER", "Female"},
		// Add reference for Religion
		{"RELIGION", "Islam"},
		{"RELIGION", "Protestant"},
		{"RELIGION", "Katolik"},
		{"RELIGION", "Hindu"},
		{"RELIGION", "Buddha"},
		{"RELIGION", "Konghucu"},
		{"RELIGION", "Other"},
		// Add reference for MaritalStatus
		{"MARRIED", "Single"},
		{"MARRIED", "Married"},
		// Add reference for Education
		{"EDUCATION", "Elementary School"},
		{"EDUCATION", "Junior High School"},
		{"EDUCATION", "High School"},
		{"EDUCATION", "Diploma"},
		{"EDUCATION", "Bachelor"},
		{"EDUCATION", "Master"},
		{"EDUCATION", "Doctor"},
		// Add reference for Education Type
		{"EDTYPE", "Formal"},
		{"EDTYPE", "Non Formal"},
		// Add reference for Nationality
		{"NATION", "Indonesian"},
		{"NATION", "Foreigner"},
		// Add reference for IdentityType
		{"IDENTITY", "ID Card"},
		{"IDENTITY", "Passport"},
		// Add reference for Socialmedia Type
		{"SOCMED", "Facebook"},
		{"SOCMED", "Twitter"},
		{"SOCMED", "Instagram"},
		// Add reference for Skill Type
		{"SKILLS", "Operating System"},
		{"SKILLS", "Application Programming"},
		{"SKILLS", "Database"},
		{"SKILLS", "Scripting Language"},
		{"SKILLS", "Web Server"},
		{"SKILLS", "Application Package"},
		{"SKILLS", "Network"},
		{"SKILLS", "Methodologies"},
		// Add reference for Role type
		{"ROLES", "User"},
		{"ROLES", "Admin"},
		{"ROLES", "Super Admin"},
	}

	indexRevcode := 0
	referenceName, tempRefName := "", "GENDER"
	for _, value := range datas {
		referenceName = value[0]
		referenceValue := value[1]

		helper.Log.Println(indexRevcode)
		helper.Log.Println("Reference name : " + referenceName)
		helper.Log.Println("Reference value : " + referenceValue)

		// log.Print(revcode)

		if tempRefName != referenceName {
			tempRefName = value[0]
			indexRevcode = 1
			// log.Println("Not Same")
		} else {
			indexRevcode++
			// log.Println("Same")
		}

		revcode := model.Revcode{
			RevCode: indexRevcode,
			Type:    referenceName,
			Value:   referenceValue,
		}

		if dbc := db.Find(&revcode, model.Revcode{Type: referenceName, Value: referenceValue}); dbc.Error != nil {
			db.Create(&revcode)
			helper.Log.Println("Reference " + referenceName + " and value " + referenceValue + " is created")
		} else {
			db.Find(&revcode, model.Revcode{Type: referenceName, Value: referenceValue}).Updates(model.Revcode{RevCode: indexRevcode, Type: referenceName, Value: referenceValue})
			helper.Log.Println("Reference " + referenceName + " and value " + referenceValue + " is already exists")
		}
		// fmt.Println(value[1])
	}

	db.Close()
}
