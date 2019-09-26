package modeltests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/gargprateek248/manage_candidates/api/controllers"
	"github.com/gargprateek248/manage_candidates/api/models"
)

var server = controllers.Server{}
var candidateInstance = models.Candidate{}

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()

	os.Exit(m.Run())
}

func Database() {

	var err error

	TestDbDriver := os.Getenv("TestDbDriver")

	if TestDbDriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbUser"), os.Getenv("TestDbName"), os.Getenv("TestDbPassword"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
}

func refreshCandidateTable() error {
	err := server.DB.DropTableIfExists(&models.Candidate{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.Candidate{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneCandidate() (models.Candidate, error) {

	refreshCandidateTable()

	candidate := models.Candidate{
		Name: "Pet",
		Phone_No:    "120 650 7689",
		Status: 2,
	}

	err := server.DB.Model(&models.Candidate{}).Create(&candidate).Error
	if err != nil {
		log.Fatalf("cannot seed candidates table: %v", err)
	}
	return candidate, nil
}

func seedCandidates() error {

	candidates := []models.Candidate{
		models.Candidate{
		Name: "Peter",
		Phone_No:    "320 650 7689",
		Status: 1,
		},
		models.Candidate{
		Name: "Peter Moore",
		Phone_No:    "345 650 7689",
		Status: 0,
		},
	}

	for i, _ := range candidates {
		err := server.DB.Model(&models.Candidate{}).Create(&candidates[i]).Error
		if err != nil {
			return err
		}
	}
	return nil
}

//func seedCandidates() ([]models.Candidate, error) {

//	var err error

//	if err != nil {
//		return []models.Candidate{}, err
//	}
//	var candidates = []models.Candidate{
//		models.Candidate{
//		Name: "Tom",
//		Phone_No:    "420 650 7689",
//		Status: 2,
//		},
//		models.Candidate{
//		Name: "Tommy",
//		Phone_No:    "520 650 7689",
//		Status: 1,
//		},
//	}

//	for i, _ := range candidates {
//		err = server.DB.Model(&models.Candidate{}).Create(&candidates[i]).Error
//		if err != nil {
//			log.Fatalf("cannot seed candidates table: %v", err)
//		}
//	}
//	return candidates, nil
//}
