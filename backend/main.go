package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Donor struct {
	Name       string `json:"name"`
	BloodGroup string `json:"blood_group"`
	Phone      string `json:"phone"`
	City       string `json:"city"`
}

const csvFile = "donors.csv"

func validatePhone(phone string) error {
	re := regexp.MustCompile(`^[0-9]{10}$`)
	if !re.MatchString(phone) {
		return errors.New("phone number must be exactly 10 digits")
	}
	return nil
}

// CSV mein data save karne ka function
func saveToCSV(d Donor) error {
	if err := validatePhone(d.Phone); err != nil {
		return err
	}
	f, err := os.OpenFile(csvFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	return writer.Write([]string{d.Name, d.BloodGroup, d.Phone, d.City})
}

// CSV se saara data read karne ka function
func readFromCSV() ([]Donor, error) {
	f, err := os.Open(csvFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Donor{}, nil
		}
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var donors []Donor
	for _, rec := range records {
		donors = append(donors, Donor{
			Name:       rec[0],
			BloodGroup: rec[1],
			Phone:      rec[2],
			City:       rec[3],
		})
	}
	return donors, nil
}

func deleteDonor(phone string) error {
	donors, _ := readFromCSV()
	f, err := os.Create(csvFile)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	for _, d := range donors {
		if d.Phone != phone {
			writer.Write([]string{d.Name, d.BloodGroup, d.Phone, d.City})
		}
	}
	return nil
}

func donorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method == "POST" {
		var d Donor

		if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		if err := saveToCSV(d); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(d)
	} else if r.Method == "GET" {
		donors, _ := readFromCSV()
		group := r.URL.Query().Get("blood_group")

		if group != "" {
			filtered := []Donor{}
			for _, d := range donors {
				if strings.EqualFold(d.BloodGroup, group) {
					filtered = append(filtered, d)
				}
			}
			json.NewEncoder(w).Encode(filtered)
		} else {
			json.NewEncoder(w).Encode(donors)
		}
	} else if r.Method == "DELETE" {
		phone := r.URL.Query().Get("phone")
		deleteDonor(phone)
		w.WriteHeader(http.StatusOK)
	} else if r.Method == "PUT" {
		var updatedDonor Donor
		json.NewDecoder(r.Body).Decode(&updatedDonor)
		deleteDonor(updatedDonor.Phone)
		saveToCSV(updatedDonor)
		json.NewEncoder(w).Encode(updatedDonor)
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Disposition", "attachment; filename=donors_list.csv")
	w.Header().Set("Content-Type", "text/csv")

	// File read karke user ko bhej dena
	http.ServeFile(w, r, csvFile)
}

func main() {
	http.HandleFunc("/donors", donorHandler)
	http.HandleFunc("/download", downloadHandler)
	fmt.Println("Server running on :8080. Data saving to donors.csv")
	http.ListenAndServe(":8080", nil)
}
