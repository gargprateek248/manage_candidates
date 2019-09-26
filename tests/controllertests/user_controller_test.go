package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/gargprateek248/manage_candidates/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateCandidate(t *testing.T) {

	err := refreshCandidateTable()
	if err != nil {
		log.Fatal(err)
	}
	samples := []struct {
		inputJSON    string
		statusCode   int
		name     string
		phone_no        string
		errorMessage string
	}{
		{
			inputJSON:    `{"name":"Pet", "phone_no": "9873345789", "status": 0}`,
			statusCode:   201,
			name:     "Pet",
			phone_no:        "9873345789",
			errorMessage: "",
		},
		{
			inputJSON:    `{"name":"Frank", "phone_no": "9873345789", "status": 1}`,
			statusCode:   500,
			errorMessage: "Phone No Already Taken",
		},
		{
			inputJSON:    `{"name": "", "phone_no": "9873345780", "status": 2}`,
			statusCode:   422,
			errorMessage: "Required name",
		},
		{
			inputJSON:    `{"name": "Kan", "phone_no": "", "status": 1}`,
			statusCode:   422,
			errorMessage: "Required Phone No",
		},
		{
			inputJSON:    `{"name": "Kan", "phone_no": "9873345723", "status": ""}`,
			statusCode:   422,
			errorMessage: "Required Status",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/candidates", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateCandidate)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["name"], v.name)
			assert.Equal(t, responseMap["phone_no"], v.phone_no)
		}
		if v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetCandidates(t *testing.T) {

	err := refreshCandidateTable()
	if err != nil {
		log.Fatal(err)
	}
	_, err = seedCandidates()
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("GET", "/candidates", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetCandidates)
	handler.ServeHTTP(rr, req)

	var candidates []models.Candidate
	err = json.Unmarshal([]byte(rr.Body.String()), &candidates)
	if err != nil {
		log.Fatalf("Cannot convert to json: %v\n", err)
	}
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(candidates), 2)
}

func TestGetCandidateByID(t *testing.T) {

	err := refreshCandidateTable()
	if err != nil {
		log.Fatal(err)
	}
	candidate, err := seedOneCandidate()
	if err != nil {
		log.Fatal(err)
	}
	candidateSample := []struct {
		id           string
		statusCode   int
		name     string
		phone_no        string
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(candidate.ID)),
			statusCode: 200,
			name:   candidate.Name,
			phone_no:      candidate.Phone_No,
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
	}
	for _, v := range candidateSample {

		req, err := http.NewRequest("GET", "/candidates", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetCandidate)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, candidate.Name, responseMap["name"])
			assert.Equal(t, candidate.Phone_No, responseMap["phone_no"])
		}
	}
}

func TestUpdateCandidate(t *testing.T) {

	var AuthID uint32

	err := refreshCandidateTable()
	if err != nil {
		log.Fatal(err)
	}
	candidates, err := seedCandidates() //we need atleast two candidates to properly check the update
	if err != nil {
		log.Fatalf("Error seeding candidate: %v\n", err)
	}
	// Get only the first candidate
	for _, candidate := range candidates {
		if candidate.ID == 2 {
			continue
		}
		AuthID = candidate.ID
	}

	samples := []struct {
		id             string
		updateJSON     string
		statusCode     int
		updateName string
		updatePhone_No    string
		errorMessage   string
	}{
		{
			// Convert int32 to int first before converting to string
			id:             strconv.Itoa(int(AuthID)),
			updateJSON:     `{"name":"Grand", "phone_no": "9876543210", "status": 1}`,
			statusCode:     200,
			updateName: "Grand",
			updatePhone_No:    "9876543210",
			errorMessage:   "",
		},
		{
			// Remember "kenny@gmail.com" belongs to candidate 2
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"name":"Frank", "phone_no": "9873345789", "status": 2}`,
			statusCode:   500,
			errorMessage: "Phone_No Already Taken",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"name": "", "phone_no": "9873145789", "status": 2}`,
			statusCode:   422,
			errorMessage: "Required Name",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			updateJSON:   `{"name": "Kan", "phone_no": "", "status": 2}`,
			statusCode:   422,
			errorMessage: "Required Phone_No",
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/candidates", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateCandidate)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["name"], v.updateName)
			assert.Equal(t, responseMap["phone_no"], v.updatePhone_No)
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeleteCandidate(t *testing.T) {

	var AuthID uint32

	err := refreshCandidateTable()
	if err != nil {
		log.Fatal(err)
	}

	candidates, err := seedCandidates() //we need atleast two candidates to properly check the update
	if err != nil {
		log.Fatalf("Error seeding candidate: %v\n", err)
	}
	// Get only the first and log him in
	for _, candidate := range candidates {
		if candidate.ID == 2 {
			continue
		}
		AuthID = candidate.ID
	}

	candidateSample := []struct {
		id           string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int32 to int first before converting to string
			id:           strconv.Itoa(int(AuthID)),
			statusCode:   204,
			errorMessage: "",
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
	}
	for _, v := range candidateSample {

		req, err := http.NewRequest("GET", "/candidates", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteCandidate)


		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 401 && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
