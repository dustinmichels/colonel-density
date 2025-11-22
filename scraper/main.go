package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

func main() {
	// All 48 continental US state abbreviations (excluding Alaska and Hawaii)
	states := []string{
		"al", "az", "ar", "ca", "co", "ct", "de", "fl", "ga", "id",
		"il", "in", "ia", "ks", "ky", "la", "me", "md", "ma", "mi",
		"mn", "ms", "mo", "mt", "ne", "nv", "nh", "nj", "nm", "ny",
		"nc", "nd", "oh", "ok", "or", "pa", "ri", "sc", "sd", "tn",
		"tx", "ut", "vt", "va", "wa", "wv", "wi", "wy",
	}

	// Channel to collect all cities
	citiesChan := make(chan []City, len(states))
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
		}(state)
	}

	// Wait for all goroutines to complete
	go func() {
		wg.Wait()
		close(citiesChan)
	}()

	// Collect all cities
	var allCities []City
	for cities := range citiesChan {
		allCities = append(allCities, cities...)
	}

	// Print summary
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Printf("SUMMARY: Found %d total cities across %d states\n", len(allCities), len(states))
	fmt.Println(strings.Repeat("=", 80))

	// Print first 10 cities as examples
	fmt.Println("\nFirst 10 cities:")
	for i, city := range allCities {
		if i >= 10 {
			break
		}
		fmt.Printf("%d. %s (%s) - Count: %d\n   URL: %s\n",
			i+1, city.PlaceName, strings.ToUpper(city.StateCode), city.DataCount, city.URL)
	}

	// Print statistics by state
	fmt.Println("\nCities by state:")
	stateCounts := make(map[string]int)
	for _, city := range allCities {
		stateCounts[city.StateCode]++
	}
	for _, state := range states {
		if count, exists := stateCounts[state]; exists {
			fmt.Printf("  %s: %d cities\n", strings.ToUpper(state), count)
		}
	}
}
