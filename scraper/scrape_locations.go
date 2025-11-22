package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Location represents a KFC location with its details
type Location struct {
	Name      string
	Address   string
	City      string
	State     string
	ZipCode   string
	Country   string
	Latitude  *float64
	Longitude *float64
}

// getLocationsFromCity fetches and parses a city page to extract location data
func getLocationsFromCity(cityURL string) ([]Location, error) {
	// Fetch the page
	resp, err := http.Get(cityURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", cityURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error for %s: %d %s", cityURL, resp.StatusCode, resp.Status)
	}

	// Parse the HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML for %s: %w", cityURL, err)
	}

	var locations []Location

	// Try to find locations in Format 1: individual location pages with Core-address
	doc.Find(".Core-address").Each(func(i int, s *goquery.Selection) {
		loc := parseLocationFormat1(s, doc)
		locations = append(locations, loc)
	})

	// If no locations found, try Format 2: directory listing with Directory-listTeaser
	if len(locations) == 0 {
		doc.Find("li.Directory-listTeaser").Each(func(i int, s *goquery.Selection) {
			loc := parseLocationFormat2(s)
			locations = append(locations, loc)
		})
	}

	return locations, nil
}

// parseLocationFormat1 parses the detailed format with Core-address and coordinates
func parseLocationFormat1(s *goquery.Selection, doc *goquery.Document) Location {
	loc := Location{}

	// Extract coordinates from meta tags
	if lat, exists := s.Find("meta[itemprop='latitude']").Attr("content"); exists {
		if latFloat, err := strconv.ParseFloat(lat, 64); err == nil {
			loc.Latitude = &latFloat
		}
	}
	if lon, exists := s.Find("meta[itemprop='longitude']").Attr("content"); exists {
		if lonFloat, err := strconv.ParseFloat(lon, 64); err == nil {
			loc.Longitude = &lonFloat
		}
	}

	// Extract address components
	addressBlock := s.Find("address.c-address")

	// Street address
	if street, exists := addressBlock.Find("meta[itemprop='streetAddress']").Attr("content"); exists {
		loc.Address = street
	} else {
		loc.Address = strings.TrimSpace(addressBlock.Find(".c-address-street-1").Text())
	}

	// City
	if city, exists := addressBlock.Find("meta[itemprop='addressLocality']").Attr("content"); exists {
		loc.City = city
	} else {
		loc.City = strings.TrimSpace(addressBlock.Find(".c-address-city").Text())
	}

	// State
	loc.State = strings.TrimSpace(addressBlock.Find(".c-address-state").Text())

	// Zip code
	loc.ZipCode = strings.TrimSpace(addressBlock.Find(".c-address-postal-code").Text())

	// Country
	loc.Country = strings.TrimSpace(addressBlock.Find(".c-address-country-name").Text())

	// Try to get location name from the page title or header
	loc.Name = strings.TrimSpace(doc.Find("h1.Core-title").Text())
	if loc.Name == "" {
		// Fallback: construct name from brand and address
		brand := strings.TrimSpace(doc.Find(".LocationName-brand").Text())
		if brand != "" {
			loc.Name = fmt.Sprintf("%s %s", brand, loc.Address)
		} else {
			loc.Name = loc.Address
		}
	}

	return loc
}

// parseLocationFormat2 parses the directory teaser format without coordinates
func parseLocationFormat2(s *goquery.Selection) Location {
	loc := Location{}

	// Extract name from the title link
	titleLink := s.Find("a.Teaser-titleLink")
	brand := strings.TrimSpace(titleLink.Find(".LocationName-brand").Text())
	geo := strings.TrimSpace(titleLink.Find(".LocationName-geo").Text())

	if brand != "" && geo != "" {
		loc.Name = fmt.Sprintf("%s %s", brand, geo)
	} else if geo != "" {
		loc.Name = geo
	} else {
		loc.Name = strings.TrimSpace(titleLink.Text())
	}

	// Extract address components
	addressBlock := s.Find("address.c-address")

	// Street address
	loc.Address = strings.TrimSpace(addressBlock.Find(".c-address-street-1").Text())

	// City
	loc.City = strings.TrimSpace(addressBlock.Find(".c-address-city").Text())

	// State
	loc.State = strings.TrimSpace(addressBlock.Find(".c-address-state").Text())

	// Zip code
	loc.ZipCode = strings.TrimSpace(addressBlock.Find(".c-address-postal-code").Text())

	// Country
	loc.Country = strings.TrimSpace(addressBlock.Find(".c-address-country-name").Text())

	// Coordinates will be nil if not available

	return loc
}
