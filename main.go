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
    PlanetName string `json:"pl_name"`
    Hostname   string `json:"pl_hostname"`
    DiscYear   int    `json:"pl_disc"`
}

func main() {
    apiKey := os.Getenv("NASA_API_KEY")
    if apiKey == "" {
        apiKey = "DEMO_KEY"
    }

    url := fmt.Sprintf("https://exoplanetarchive.ipac.caltech.edu/cgi-bin/nstedAPI/nph-nstedAPI?table=exoplanets&format=json&api_key=%s", apiKey)

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

    for i, p := range planets {
        if i >= 5 {
            break
        }
        fmt.Printf("%s (discovered: %d)\n", p.PlanetName, p.DiscYear)
    }
}