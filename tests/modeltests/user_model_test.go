package modeltests

import (
	"log"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres driver
	"github.com/gargprateek248/manage_candidates/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllCandidates(t *testing.T) {

	err := refreshCandidateTable()
	if err != nil {
		log.Fatal(err)
	}

	err = seedCandidates()
	if err != nil {
		log.Fatal(err)
	}

	candidates, err := candidateInstance.FindAllCandidates(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the candidates: %v\n", err)
		return
	}
	assert.Equal(t, len(*candidates), 2)
}

func TestSaveCandidate(t *testing.T) {

	err := refreshCandidateTable()
	if err != nil {
		log.Fatal(err)
	}
	newCandidate := models.Candidate{
		ID:       1,
		Name: "test",
		Phone_No:    "120 987 6543",
		Status: 2,
	}
	savedCandidate, err := newCandidate.SaveCandidate(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the candidates: %v\n", err)
		return
	}
	assert.Equal(t, newCandidate.ID, savedCandidate.ID)
	assert.Equal(t, newCandidate.Phone_No, savedCandidate.Phone_No)
	assert.Equal(t, newCandidate.Name, savedCandidate.Name)
}

func TestGetCandidateByID(t *testing.T) {

	err := refreshCandidateTable()
	if err != nil {
		log.Fatal(err)
	}

	candidate, err := seedOneCandidate()
	if err != nil {
		log.Fatalf("cannot seed candidates table: %v", err)
	}
	foundCandidate, err := candidateInstance.FindCandidateByID(server.DB, candidate.ID)
	if err != nil {
		t.Errorf("this is the error getting one candidate: %v\n", err)
		return
	}
	assert.Equal(t, foundCandidate.ID, candidate.ID)
	assert.Equal(t, foundCandidate.Phone_No, candidate.Phone_No)
	assert.Equal(t, foundCandidate.Name, candidate.Name)
}

func TestUpdateACandidate(t *testing.T) {

	err := refreshCandidateTable()
	if err != nil {
		log.Fatal(err)
	}

	candidate, err := seedOneCandidate()
	if err != nil {
		log.Fatalf("Cannot seed candidate: %v\n", err)
	}

	candidateUpdate := models.Candidate{
		ID:       1,
		Name: "testupdate",
		Phone_No:    "220 987 6543",
		Status: 1,
	}
	updatedCandidate, err := candidateUpdate.UpdateACandidate(server.DB, candidate.ID)
	if err != nil {
		t.Errorf("this is the error updating the candidate: %v\n", err)
		return
	}
	assert.Equal(t, updatedCandidate.ID, candidateUpdate.ID)
	assert.Equal(t, updatedCandidate.Phone_No, candidateUpdate.Phone_No)
	assert.Equal(t, updatedCandidate.Name, candidateUpdate.Name)
}

func TestDeleteACandidate(t *testing.T) {

	err := refreshCandidateTable()
	if err != nil {
		log.Fatal(err)
	}

	candidate, err := seedOneCandidate()

	if err != nil {
		log.Fatalf("Cannot seed candidate: %v\n", err)
	}

	isDeleted, err := candidateInstance.DeleteACandidate(server.DB, candidate.ID)
	if err != nil {
		t.Errorf("this is the error updating the candidate: %v\n", err)
		return
	}
	//one shows that the record has been deleted or:
	// assert.Equal(t, int(isDeleted), 1)

	//Can be done this way too
	assert.Equal(t, isDeleted, int64(1))
}
