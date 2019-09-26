package controllertests

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

	err := refreshCandidateTable()
	if err != nil {
		log.Fatal(err)
	}

	candidate := models.Candidate{
			Name: "Steve Smith",
			Phone_No:    "423 456 7890",
			Status: 0,
	}

	err = server.DB.Model(&models.Candidate{}).Create(&candidate).Error
	if err != nil {
		return models.Candidate{}, err
	}
	return candidate, nil
}

func seedCandidates() ([]models.Candidate, error) {

	var err error
	if err != nil {
		return nil, err
	}
	candidates := []models.Candidate{
		models.Candidate{
			Name: "Steve Waugh",
			Phone_No:    "123 456 7890",
			Status: 2,
		},
		models.Candidate{
			Name: "Mark Waugh",
			Phone_No:    "223 456 7890",
			Status: 1,
		},
	}
	for i, _ := range candidates {
		err := server.DB.Model(&models.Candidate{}).Create(&candidates[i]).Error
		if err != nil {
			return []models.Candidate{}, err
		}
	}
	return candidates, nil
}
