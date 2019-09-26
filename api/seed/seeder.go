package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/gargprateek248/manage_candidates/api/models"
)

var candidates = []models.Candidate{
	models.Candidate{
		Name: "Prateek Garg",
		Phone_No:    "9873346570",
		Status: 0,
	},
	models.Candidate{
		Name: "Martin Luther",
		Phone_No:    "987 44 5670",
		Status: 1,
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Candidate{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.Candidate{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	for i, _ := range candidates {
		err = db.Debug().Model(&models.Candidate{}).Create(&candidates[i]).Error
		if err != nil {
			log.Fatalf("cannot seed candidates table: %v", err)
		}
	}
}
