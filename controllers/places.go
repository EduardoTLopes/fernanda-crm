package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/EduardoTLopes/fernanda/database"
	"github.com/EduardoTLopes/fernanda/models"
)

func ShowPlace(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var id = r.PathValue("ID")
		var place models.Place
		database.DB.First(&place, id)

		Templates.ExecuteTemplate(w, "ShowPlace", place)
	}

	if r.Method == http.MethodPatch {
		var id = r.PathValue("ID")
		var place models.Place
		database.DB.First(&place, id)
		place.Contacted = !place.Contacted
		database.DB.Model(&place).Update("Contacted", place.Contacted)
		http.Redirect(w, r, "/places", http.StatusSeeOther)
	}

	if r.Method == http.MethodDelete {
		id := r.PathValue("ID")
		var place models.Place
		if err := database.DB.First(&place, id).Error; err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		database.DB.Delete(&place)
		http.Redirect(w, r, "/places", http.StatusSeeOther)
	}
}

func ShowAllPlaces(w http.ResponseWriter, r *http.Request) {
	var places []models.Place
	database.DB.Order("Name ASC").Find(&places)
	Templates.ExecuteTemplate(w, "ListPlaces", places)
}

func AddPlaces(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		Templates.ExecuteTemplate(w, "AddPlaces", nil)
		return
	}
	if r.Method == http.MethodPost {
		log.Default().Println("POST")
		apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")

		latitude := r.FormValue("latitude")
		longitude := r.FormValue("longitude")
		radius := "3000"
		placeType := "bakery,restaurant"

		url := fmt.Sprintf("https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=%s,%s&radius=%s&type=%s&key=%s", latitude, longitude, radius, placeType, apiKey)

		resp, err := http.Get(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var result struct {
			Results []Place `json:"results"`
		}
		if err := json.Unmarshal(body, &result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		places := result.Results
		for i, place := range places {
			var dbPlace = models.Place{}

			database.DB.Where(&models.Place{PlaceID: place.PlaceID}).Where("deleted_at IS NOT NULL").First(&dbPlace)
			if dbPlace.ID == 0 {
				detailUrl := fmt.Sprintf("https://maps.googleapis.com/maps/api/place/details/json?place_id=%s&fields=formatted_phone_number&key=%s", place.PlaceID, apiKey)
				detailResp, err := http.Get(detailUrl)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer detailResp.Body.Close()

				detailBody, err := io.ReadAll(detailResp.Body)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				var detailResult struct {
					Result struct {
						FormattedPhoneNumber string `json:"formatted_phone_number"`
					} `json:"result"`
				}
				if err := json.Unmarshal(detailBody, &detailResult); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				places[i].FormattedPhoneNumber = detailResult.Result.FormattedPhoneNumber
				var dbPlace = models.Place{
					PlaceID:              place.PlaceID,
					Name:                 place.Name,
					FormattedPhoneNumber: detailResult.Result.FormattedPhoneNumber,
					PriceLevel:           place.PriceLevel,
					Rating:               place.Rating,
					Types:                strings.Join(place.Types, ","),
					Vicinity:             place.Vicinity,
					Contacted:            false,
					ContactName:          "",
					Notes:                "",
				}
				database.DB.Create(&dbPlace)
			}
		}

		Templates.ExecuteTemplate(w, "AddPlaces", places)
	}

}

func UpdateContactNameAndNotes(w http.ResponseWriter, r *http.Request) {
	var id = r.PathValue("ID")
	var place models.Place
	database.DB.First(&place, id)
	if r.Method == http.MethodPost {
		contactName := r.FormValue("contactName")
		notes := r.FormValue("notes")
		place.ContactName = contactName

		// Initialize or update Notes
		place.Notes = notes // Append new notes to existing (or newly initialized) string
		database.DB.Model(&place).Update("ContactName", place.ContactName)
		err := database.DB.Model(&place).Update("Notes", place.Notes)
		if err.Error != nil {
			log.Default().Println(err.Error)
		}
		http.Redirect(w, r, "/place/"+id, http.StatusSeeOther)
	}
}
