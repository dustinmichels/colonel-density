package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Location struct {
	Name    string
	URL     string
	Address string
	City    string
	State   string
	Zip     string
}

func main() {
	// URL to scrape
	baseURL := "https://locations.kfc.com"
	url := "https://locations.kfc.com/ma"

	// Make HTTP request
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	// Check if request was successful
	if resp.StatusCode != 200 {
		log.Fatalf("Failed to fetch page: status code %d", resp.StatusCode)
	}

	// Parse the HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}

	fmt.Println("Extracting location links from KFC Massachusetts page...")
	fmt.Println(strings.Repeat("=", 80))

	// Collect all location URLs
	var locationURLs []string
	doc.Find("a.Directory-listLink").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			// Build full URL
			fullURL := baseURL + href
			locationURLs = append(locationURLs, fullURL)
		}
	})

	fmt.Printf("Found %d locations. Fetching addresses...\n\n", len(locationURLs))

	// Visit each location and extract address
	locations := make([]Location, 0, len(locationURLs))
	for i, locURL := range locationURLs {
		fmt.Printf("[%d/%d] Fetching: %s\n", i+1, len(locationURLs), locURL)

		location := fetchLocationDetails(locURL)
		if location != nil {
			locations = append(locations, *location)
			fmt.Printf("✓ Address: %s\n", location.Address)
			if location.City != "" {
				fmt.Printf("  City: %s, %s %s\n", location.City, location.State, location.Zip)
			}
		} else {
			fmt.Printf("✗ Failed to fetch address\n")
		}
		fmt.Println()

		// Be polite - add small delay between requests
		time.Sleep(500 * time.Millisecond)
	}

	// Print summary
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("SUMMARY")
	fmt.Println(strings.Repeat("=", 80))
	for i, loc := range locations {
		fmt.Printf("%d. %s\n", i+1, loc.Address)
		if loc.City != "" {
			fmt.Printf("   %s, %s %s\n", loc.City, loc.State, loc.Zip)
		}
		fmt.Printf("   URL: %s\n\n", loc.URL)
	}
	fmt.Printf("Total addresses collected: %d/%d\n", len(locations), len(locationURLs))
}

func fetchLocationDetails(url string) *Location {
	// Make HTTP request
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching %s: %v", url, err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Failed to fetch %s: status code %d", url, resp.StatusCode)
		return nil
	}

	// Parse the HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("Error parsing HTML from %s: %v", url, err)
		return nil
	}

	location := &Location{URL: url}

	// Extract street address
	doc.Find(".c-AddressRow .c-address-street-1").Each(func(i int, s *goquery.Selection) {
		location.Address = strings.TrimSpace(s.Text())
	})

	// Extract city, state, zip
	doc.Find(".c-address-city").Each(func(i int, s *goquery.Selection) {
		location.City = strings.TrimSpace(s.Text())
	})

	doc.Find(".c-address-state").Each(func(i int, s *goquery.Selection) {
		location.State = strings.TrimSpace(s.Text())
	})

	doc.Find(".c-address-postal-code").Each(func(i int, s *goquery.Selection) {
		location.Zip = strings.TrimSpace(s.Text())
	})

	// If no address found, return nil
	if location.Address == "" {
		return nil
	}

	return location
}
