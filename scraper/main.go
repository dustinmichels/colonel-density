package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

// printGreen prints text in green color
func printGreen(format string, args ...interface{}) {
	fmt.Printf("\033[32m"+format+"\033[0m", args...)
}

// printRed prints text in red color
func printRed(format string, args ...interface{}) {
	fmt.Printf("\033[31m"+format+"\033[0m", args...)
}

func main() {
	// All 48 continental US state abbreviations (excluding Alaska and Hawaii)
	states := []string{
		"al", "az", "ar", "ca", "co", "ct", "de", "fl", "ga", "id",
		"il", "in", "ia", "ks", "ky", "la", "me", "md", "ma", "mi",
		"mn", "ms", "mo", "mt", "ne", "nv", "nh", "nj", "nm", "ny",
		"nc", "nd", "oh", "ok", "or", "pa", "ri", "sc", "sd", "tn",
		"tx", "ut", "vt", "va", "wa", "wv", "wi", "wy",
	}

	// Channels to collect all cities and locations
	citiesChan := make(chan []City, len(states))
	locationsChan := make(chan []Location, len(states)*100) // Larger buffer for locations
	var wg sync.WaitGroup

	// Process each state in parallel
	for _, state := range states {
		wg.Add(1)
		go func(stateCode string) {
			defer wg.Done()

			url := fmt.Sprintf("https://locations.kfc.com/%s", stateCode)
			fmt.Printf("Fetching cities for %s...\n", strings.ToUpper(stateCode))

			cities, err := getCitiesOnStatePage(url, stateCode)
			if err != nil {
				log.Printf("Error fetching %s: %v", stateCode, err)
				return
			}

			citiesChan <- cities
			fmt.Printf("✓ Found %d cities in %s\n", len(cities), strings.ToUpper(stateCode))

			// Now get locations from each city in this state
			for _, city := range cities {
				fmt.Printf("  [%s] Fetching locations from %s...\n",
					strings.ToUpper(stateCode), city.PlaceName)

				locations, err := getLocationsFromCity(city.URL)
				if err != nil {
					log.Printf("  Error fetching locations from %s: %v", city.URL, err)
					printRed("    ✗ Failed to get locations\n")
					continue
				}

				// Report success with expected vs actual count
				printGreen("    ✓ Successfully got %d/%d locations\n", len(locations), city.DataCount)

				if len(locations) != city.DataCount {
					log.Printf("  Warning: Expected %d locations but got %d for %s",
						city.DataCount, len(locations), city.URL)
				}

				if len(locations) > 0 {
					locationsChan <- locations
				}
			}
		}(state)
	}

	// Wait for all goroutines to complete
	go func() {
		wg.Wait()
		close(citiesChan)
		close(locationsChan)
	}()

	// Collect all cities
	var allCities []City
	for cities := range citiesChan {
		allCities = append(allCities, cities...)
	}

	// Collect all locations
	var allLocations []Location
	for locations := range locationsChan {
		allLocations = append(allLocations, locations...)
	}

	// Print summary
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Printf("SUMMARY: Found %d total cities and %d total locations across %d states\n",
		len(allCities), len(allLocations), len(states))
	fmt.Println(strings.Repeat("=", 80))

	// Write Locations CSV file
	os.MkdirAll("out", 0755)

	locationFile, err := os.Create("out/locations.csv")
	if err != nil {
		log.Fatalf("Failed to create locations CSV file: %v", err)
	}
	defer locationFile.Close()

	locationWriter := csv.NewWriter(locationFile)
	defer locationWriter.Flush()

	// CSV header
	locationWriter.Write([]string{"name", "address", "city", "state", "zip_code", "country", "latitude", "longitude"})

	// CSV rows
	for _, loc := range allLocations {
		// Handle nullable latitude and longitude
		latStr := ""
		if loc.Latitude != nil {
			latStr = fmt.Sprintf("%.8f", *loc.Latitude)
		}

		lonStr := ""
		if loc.Longitude != nil {
			lonStr = fmt.Sprintf("%.8f", *loc.Longitude)
		}

		record := []string{
			loc.Name,
			loc.Address,
			loc.City,
			loc.State,
			loc.ZipCode,
			loc.Country,
			latStr,
			lonStr,
		}
		locationWriter.Write(record)
	}

	fmt.Println("Saved locations CSV to out/locations.csv")
}
