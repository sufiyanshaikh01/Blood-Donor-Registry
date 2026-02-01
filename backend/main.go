package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type Donor struct {
	Name       string `json:"name"`
	BloodGroup string `json:"blood_group"`
	Phone      string `json:"phone"`
	City       string `json:"city"`
}

const csvFile = "donors.csv"

// CSV mein data save karne ka function
func saveToCSV(d Donor) error {
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

func donorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method == "POST" {
		var d Donor
		json.NewDecoder(r.Body).Decode(&d)
		if err := saveToCSV(d); err != nil {
			http.Error(w, "Error saving data", 500)
			return
		}
		json.NewEncoder(w).Encode(d)

	} else if r.Method == "GET" { kil
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
	}
}

func main() {
	http.HandleFunc("/donors", donorHandler)
	fmt.Println("Server running on :8080. Data saving to donors.csv")
	http.ListenAndServe(":8080", nil)
}