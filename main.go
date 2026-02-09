package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
)

type Exoplanet struct {
	PlanetName string  `json:"pl_name"`
	Hostname   string  `json:"hostname"`
	DiscYear   int     `json:"disc_year"`
	Mass       float64 `json:"pl_bmasse"` // planet mass
	Radius     float64 `json:"pl_rade"`   // planet radius
}

func main() {
	apiKey := os.Getenv("NASA_API_KEY")
	if apiKey == "" {
		apiKey = "DEMO_KEY"
	}

	// Updated URL with correct table and format
	url := fmt.Sprintf("https://exoplanetarchive.ipac.caltech.edu/TAP/sync?query=select+pl_name,hostname,disc_year,pl_bmasse,pl_rade+from+ps+where+tran_flag=1&format=json&api_key=%s", apiKey)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}

	var planets []Exoplanet
	err = json.Unmarshal(body, &planets)
	if err != nil {
		fmt.Println("Error parsing:", err)
		return
	}

	sort.Slice(planets, func(i, j int) bool {
		return planets[i].PlanetName < planets[j].PlanetName
	})

	// sort by discovery year (newest first)
	sort.Slice(planets, func(i, j int) bool {
		return planets[i].DiscYear > planets[j].DiscYear
	})

	fmt.Printf("\nFound %d planets:\n", len(planets))
	for i, p := range planets {
		if i >= 5 {
			break
		}
		fmt.Printf("%s (host: %s, discovered: %d)\n", p.PlanetName, p.Hostname, p.DiscYear)
	}
}
