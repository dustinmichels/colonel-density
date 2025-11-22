package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// City represents a KFC city/location with its details
type City struct {
	URL       string
	DataCount int
	PlaceName string
	StateCode string
}

// getCitiesOnStatePage fetches and parses a state page to extract city data
func getCitiesOnStatePage(stateURL string, stateCode string) ([]City, error) {
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

	var cities []City

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

		cities = append(cities, City{
			URL:       fullURL,
			DataCount: dataCount,
			PlaceName: placeName,
			StateCode: stateCode,
		})
	})

	return cities, nil
}
