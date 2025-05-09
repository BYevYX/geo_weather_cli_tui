package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type GeoapifyResponse struct {
	Results []struct {
		Datasource struct {
			Datasource string            `json:"sourcename"`
			Attribution string           `json:"attribution"`
			License     string           `json:"license"`
			Raw         map[string]interface{} `json:"raw"` // Для нестандартных полей
		} `json:"datasource"`
		Name        string  `json:"name"`
		Country     string  `json:"country"`
		CountryCode string  `json:"country_code"`
		State       string  `json:"state"`
		City        string  `json:"city"`
		Postcode    string  `json:"postcode"`
		Suburb      string  `json:"suburb"`
		Street      string  `json:"street"`
		Housenumber string  `json:"housenumber"`
		Lon         float64 `json:"lon"`
		Lat         float64 `json:"lat"`
		Distance    float64 `json:"distance"`
		Importance  float64 `json:"importance"`
		Timezone struct{
            Name        string `json:"name"`
         } `json:"timezone"`
		Address     struct {
			Name        string `json:"name"`
			Country     string `json:"country"`
			CountryCode string `json:"country_code"`
			State       string `json:"state"`
			City        string `json:"city"`
			Postcode    string `json:"postcode"`
			District    string `json:"district"`
			Suburb      string `json:"suburb"`
			Street      string `json:"street"`
			Housenumber string `json:"housenumber"`
			FullAddress string `json:"formatted"`
		} `json:"address"`
		Bbox struct {
			Lon1 float64 `json:"lon1"`
			Lat1 float64 `json:"lat1"`
			Lon2 float64 `json:"lon2"`
			Lat2 float64 `json:"lat2"`
		} `json:"bbox"`
	} `json:"results"`
	Query struct {
		Text    string `json:"text"`
		Parsed  struct {
			City     string `json:"city"`
			Country  string `json:"country"`
			District string `json:"district"`
		} `json:"parsed"`
	} `json:"query"`
	Status struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"status"`
}

func GetCoordinatesFromAddres(searchText string) GeoapifyResponse {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Unnable to load .env file")
	}

	geoApiKey := os.Getenv("GEO_API_KEY")

	url := fmt.Sprintf("https://api.geoapify.com/v1/geocode/search?text=%s&format=json&apiKey=%s", searchText, geoApiKey)
	resp, err := http.Get(url)

	var coords GeoapifyResponse

	if err != nil {
		log.Fatal("Ошибка запроса:", err)
	}
	defer resp.Body.Close()


	if err != nil {
		log.Fatal("Ошибка чтения ответа:", err)
	}

	
	if err := json.NewDecoder(resp.Body).Decode(&coords); err != nil {
		fmt.Printf("Failed to decode JSON: %v\n", err)
	}

	return coords
}
