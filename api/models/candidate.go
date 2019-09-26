package models

import (
	"errors"
	"html"
	"strings"

	"github.com/jinzhu/gorm"
)

type Status int

const (
	progress  Status = 0
	selected  Status = 1
	open	  Status = 2
)

type Candidate struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"size:255;not null;" json:"name"`
	Phone_No  string    `gorm:"size:100;not null;unique" json:"phone_no"`
	Status    Status    `gorm:"" json:"status"`
}

func (c *Candidate) Prepare() {
	c.ID = 0
	c.Name = html.EscapeString(strings.TrimSpace(c.Name))
	c.Phone_No = html.EscapeString(strings.TrimSpace(c.Phone_No))
	c.Status = 0
}

func (c *Candidate) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if c.Name == "" {
			return errors.New("Required Name")
		}
		if c.Phone_No == "" {
			return errors.New("Required Phone Number")
		}

		return nil

	default:
		if c.Name == "" {
			return errors.New("Required Name")
		}
		if c.Phone_No == "" {
			return errors.New("Required Phone Number")
		}
		return nil
	}
}

func (c *Candidate) SaveCandidate(db *gorm.DB) (*Candidate, error) {

	var err error
	err = db.Debug().Create(&c).Error
	if err != nil {
		return &Candidate{}, err
	}
	return c, nil
}

func (c *Candidate) FindAllCandidates(db *gorm.DB) (*[]Candidate, error) {
	var err error
	candidates := []Candidate{}
	err = db.Debug().Model(&Candidate{}).Limit(100).Find(&candidates).Error
	if err != nil {
		return &[]Candidate{}, err
	}
	return &candidates, err
}

func (c *Candidate) FindCandidateByID(db *gorm.DB, uid uint32) (*Candidate, error) {
	var err error
	err = db.Debug().Model(Candidate{}).Where("id = ?", uid).Take(&c).Error
	if err != nil {
		return &Candidate{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Candidate{}, errors.New("Candidate Not Found")
	}
	return c, err
}

func (c *Candidate) UpdateACandidate(db *gorm.DB, uid uint32) (*Candidate, error) {
	var err error
	db = db.Debug().Model(&Candidate{}).Where("id = ?", uid).Take(&Candidate{}).UpdateColumns(
		map[string]interface{}{
			"name":  c.Name,
			"phone_no":     c.Phone_No,
			"status": c.Status,
		},
	)
	if db.Error != nil {
		return &Candidate{}, db.Error
	}
	// This is the display the updated user
	err = db.Debug().Model(&Candidate{}).Where("id = ?", uid).Take(&c).Error
	if err != nil {
		return &Candidate{}, err
	}
	return c, nil
}

func (c *Candidate) DeleteACandidate(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&Candidate{}).Where("id = ?", uid).Take(&Candidate{}).Delete(&Candidate{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
