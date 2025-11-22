package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// Location represents a KFC location with its details
type Location struct {
	URL       string
	DataCount int
	PlaceName string
	StateCode string
}

// getLocationsOnStatePage fetches and parses a state page to extract location data
func getLocationsOnStatePage(stateURL string, stateCode string) ([]Location, error) {
	// Fetch the page
	resp, err := http.Get(stateURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", stateURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error for %s: %d %s", stateURL, resp.StatusCode, resp.Status)
	}

	// Parse the HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML for %s: %w", stateURL, err)
	}

	var locations []Location

	// Find all list items in the Directory-content
	doc.Find(".Directory-content li.Directory-listItem").Each(func(i int, s *goquery.Selection) {
		// Find the link within each list item
		link := s.Find("a.Directory-listLink")

		// Extract the href attribute
		href, exists := link.Attr("href")
		if !exists {
			return
		}

		// Extract the data-count attribute
		dataCountStr, exists := link.Attr("data-count")
		if !exists {
			return
		}

		// Parse data-count: remove parentheses and convert to int
		dataCountStr = strings.Trim(dataCountStr, "()")
		dataCount, err := strconv.Atoi(dataCountStr)
		if err != nil {
			log.Printf("Warning: failed to parse data-count '%s' for %s: %v", dataCountStr, href, err)
			dataCount = 0
		}

		// Extract the place name
		placeName := strings.TrimSpace(link.Find(".Directory-listLinkText").Text())

		// Create full URL (relative to base)
		fullURL := fmt.Sprintf("https://locations.kfc.com/%s", href)

		locations = append(locations, Location{
			URL:       fullURL,
			DataCount: dataCount,
			PlaceName: placeName,
			StateCode: stateCode,
		})
	})

	return locations, nil
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

	// Channel to collect all locations
	locationsChan := make(chan []Location, len(states))
	var wg sync.WaitGroup

	// Process each state in parallel
	for _, state := range states {
		wg.Add(1)
		go func(stateCode string) {
			defer wg.Done()

			url := fmt.Sprintf("https://locations.kfc.com/%s", stateCode)
			fmt.Printf("Fetching locations for %s...\n", strings.ToUpper(stateCode))

			locations, err := getLocationsOnStatePage(url, stateCode)
			if err != nil {
				log.Printf("Error fetching %s: %v", stateCode, err)
				return
			}

			locationsChan <- locations
			fmt.Printf("✓ Found %d locations in %s\n", len(locations), strings.ToUpper(stateCode))
		}(state)
	}

	// Wait for all goroutines to complete
	go func() {
		wg.Wait()
		close(locationsChan)
	}()

	// Collect all locations
	var allLocations []Location
	for locations := range locationsChan {
		allLocations = append(allLocations, locations...)
	}

	// Print summary
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Printf("SUMMARY: Found %d total locations across %d states\n", len(allLocations), len(states))
	fmt.Println(strings.Repeat("=", 80))

	// Print first 10 locations as examples
	fmt.Println("\nFirst 10 locations:")
	for i, loc := range allLocations {
		if i >= 10 {
			break
		}
		fmt.Printf("%d. %s (%s) - Count: %d\n   URL: %s\n",
			i+1, loc.PlaceName, strings.ToUpper(loc.StateCode), loc.DataCount, loc.URL)
	}

	// Print statistics by state
	fmt.Println("\nLocations by state:")
	stateCounts := make(map[string]int)
	for _, loc := range allLocations {
		stateCounts[loc.StateCode]++
	}
	for _, state := range states {
		if count, exists := stateCounts[state]; exists {
			fmt.Printf("  %s: %d locations\n", strings.ToUpper(state), count)
		}
	}
}
