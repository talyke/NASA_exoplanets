package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
)

type Exoplanet struct {
	PlanetName string  `json:"pl_name"`
	Hostname   string  `json:"hostname"`
	DiscYear   int     `json:"disc_year"`
	Mass       float64 `json:"pl_bmasse"`
	Radius     float64 `json:"pl_rade"`
	Distance   float64 `json:"sy_dist"`
}

func main() {
	apiKey := os.Getenv("NASA_API_KEY")
	if apiKey == "" {
		apiKey = "DEMO_KEY"
	}

	fmt.Println("🌌 NASA Exoplanet Explorer 🪐")
	fmt.Println(strings.Repeat("=", 50))

	url := fmt.Sprintf("https://exoplanetarchive.ipac.caltech.edu/TAP/sync?query=select+pl_name,hostname,disc_year,pl_bmasse,pl_rade,sy_dist+from+ps+where+tran_flag=1&format=json&api_key=%s", apiKey)

	fmt.Println("🔄 Fetching data from NASA...")
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("❌ Error fetching:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("❌ Error reading:", err)
		return
	}

	var planets []Exoplanet
	err = json.Unmarshal(body, &planets)
	if err != nil {
		fmt.Println("❌ Error parsing:", err)
		return
	}

	fmt.Printf("✅ Loaded %d exoplanets!\n\n", len(planets))

	// Search functionality
	fmt.Print("🔍 Search for a planet/star (or press Enter to skip): ")
	var search string
	fmt.Scanln(&search)

	if search != "" {
		filtered := []Exoplanet{}
		for _, p := range planets {
			if strings.Contains(strings.ToLower(p.PlanetName), strings.ToLower(search)) ||
				strings.Contains(strings.ToLower(p.Hostname), strings.ToLower(search)) {
				filtered = append(filtered, p)
			}
		}
		planets = filtered
		if len(planets) == 0 {
			fmt.Println("❌ No matches found!")
			return
		}
		fmt.Printf("✅ Found %d matches!\n\n", len(planets))
	}

	// sorting menu
	fmt.Println("📊 Choose sorting:")
	fmt.Println("1. By name (A-Z)")
	fmt.Println("2. By discovery year (newest first)")
	fmt.Println("3. By size (largest first)")
	fmt.Println("4. By distance (closest first)")

	var choice int
	fmt.Print("\n➤ Enter choice (1-4): ")
	fmt.Scan(&choice)

	switch choice {
	case 1:
		sort.Slice(planets, func(i, j int) bool {
			return planets[i].PlanetName < planets[j].PlanetName
		})
	case 2:
		sort.Slice(planets, func(i, j int) bool {
			return planets[i].DiscYear > planets[j].DiscYear
		})
	case 3:
		sort.Slice(planets, func(i, j int) bool {
			return planets[i].Radius > planets[j].Radius
		})
	case 4:
		sort.Slice(planets, func(i, j int) bool {
			return planets[i].Distance < planets[j].Distance
		})
	default:
		fmt.Println("Invalid choice, using default (name)")
	}

	// display results
	fmt.Printf("\n🪐 Top 10 Results:\n")
	fmt.Println(strings.Repeat("-", 90))

	for i, p := range planets {
		if i >= 10 {
			break
		}

		sizeEmoji := "🔵"
		if p.Radius < 1.5 {
			sizeEmoji = "🌍"
		} else if p.Radius > 10 {
			sizeEmoji = "🪐"
		}

		fmt.Printf("%s %-25s | Host: %-20s | Year: %4d | Radius: %5.2f⊕ | Dist: %6.1f pc\n",
			sizeEmoji, p.PlanetName, p.Hostname, p.DiscYear, p.Radius, p.Distance)
	}

	fmt.Println(strings.Repeat("-", 90))

	// Statistics
	showStats(planets)

	// Save option
	fmt.Print("\n💾 Save results to file? (y/n): ")
	var save string
	fmt.Scan(&save)
	if save == "y" || save == "Y" {
		saveToFile(planets)
	}
}

func showStats(planets []Exoplanet) {
	var earthLike, superEarth, gasGiant int
	var totalRadius, totalMass float64

	for _, p := range planets {
		totalRadius += p.Radius
		totalMass += p.Mass

		if p.Radius <= 1.5 {
			earthLike++
		} else if p.Radius <= 4 {
			superEarth++
		} else {
			gasGiant++
		}
	}

	total := float64(len(planets))
	fmt.Printf("\n📈 Statistics:\n")
	fmt.Printf("   🌍 Earth-like: %d (%.1f%%)\n", earthLike, float64(earthLike)/total*100)
	fmt.Printf("   🌏 Super-Earths: %d (%.1f%%)\n", superEarth, float64(superEarth)/total*100)
	fmt.Printf("   🪐 Gas Giants: %d (%.1f%%)\n", gasGiant, float64(gasGiant)/total*100)
	fmt.Printf("   📏 Avg radius: %.2f Earth radii\n", totalRadius/total)
}

func saveToFile(planets []Exoplanet) {
	f, err := os.Create("exoplanets_results.txt")
	if err != nil {
		fmt.Println("❌ Error saving:", err)
		return
	}
	defer f.Close()

	fmt.Fprintf(f, "NASA Exoplanet Search Results\n")
	fmt.Fprintf(f, strings.Repeat("=", 80)+"\n\n")

	for _, p := range planets {
		fmt.Fprintf(f, "Planet: %s\n", p.PlanetName)
		fmt.Fprintf(f, "  Host Star: %s\n", p.Hostname)
		fmt.Fprintf(f, "  Discovery Year: %d\n", p.DiscYear)
		fmt.Fprintf(f, "  Radius: %.2f Earth radii\n", p.Radius)
		fmt.Fprintf(f, "  Distance: %.1f parsecs\n\n", p.Distance)
	}

	fmt.Println("✅ Saved to exoplanets_results.txt")
}
