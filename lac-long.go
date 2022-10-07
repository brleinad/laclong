package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strings"
)

type TickIndex int8

const (
	DateIndex TickIndex = iota
	RouteIndex
	GradeIndex
	NotesIndex
	UrlIndex
	PitchesIndex
	AreaIndex
	AveStarsIndex
	StarsIndex
	StyleIndex
	LeadStyleIndex
	RouteTypeIndex
)

type Tick struct {
	date      string
	routeName string
	grade     string
	notes     string
	url       string
	pitches   string
	area      string
	aveStars  string
	stars     string
	style     string
	leadStyle string
	routeType string
}

func main() {
	ticks := GetTicks()

	lacLongTicks := GetLacLongTicks(ticks)
	fmt.Print(lacLongTicks)
	// fmt.Println(numberOfTicks)
	// get ticks for each contestant
	// create map of route -> route info and contestants that have done it
	// For each contestant return routes not done
}

func GetTicks() []Tick {
	const URL = "https://www.mountainproject.com/user/200222284/daniel-rodas-bautista/tick-export"
	ticksResponse, err := http.Get(URL)

	if err != nil {
		fmt.Println("There was an error")
	}

	csvReader := csv.NewReader(ticksResponse.Body)
	csvTicks, err := csvReader.ReadAll()

	if err != nil {
		fmt.Println("There was an error")
	}

	return csv2Ticks(csvTicks)
}

func csv2Ticks(csvTicks [][]string) []Tick {
	var ticks []Tick
	for i := 0; i < len(csvTicks); i++ {
		tick := Tick{
			date:      csvTicks[i][DateIndex],
			routeName: csvTicks[i][RouteIndex],
			grade:     csvTicks[i][GradeIndex],
			notes:     csvTicks[i][NotesIndex],
			url:       csvTicks[i][UrlIndex],
			pitches:   csvTicks[i][PitchesIndex],
			area:      csvTicks[i][AreaIndex],
			aveStars:  csvTicks[i][AveStarsIndex],
			stars:     csvTicks[i][StarsIndex],
			style:     csvTicks[i][StyleIndex],
			leadStyle: csvTicks[i][LeadStyleIndex],
			routeType: csvTicks[i][RouteTypeIndex],
		}
		ticks = append(ticks, tick)
	}
	return ticks
}

func GetLacLongTicks(ticks []Tick) []Tick {
	var lacLongTicks []Tick

	for i := 0; i < len(ticks); i++ {
		if IsInLacLong(ticks[i]) && IsValidSend(ticks[i]) {
			lacLongTicks = append(lacLongTicks, ticks[i])
		}
	}
	return lacLongTicks
}

func IsInLacLong(tick Tick) bool {
	const lacLongAreaName = "International > North America > Canada > Quebec > 03. Quebec City, Charlevoix, Portneuf > Lac Long"
	return strings.Contains(tick.area, lacLongAreaName)
}

func IsValidSend(tick Tick) bool {
	isValid := false
	if tick.style == "Lead" {
		isValid = tick.leadStyle == "Redpoint" || tick.leadStyle == "Onsight" || tick.leadStyle == "Flash"
	}
	if tick.routeType == "Boulder" {
		isValid = tick.style == "Send" || tick.style == "Flash"
	}

	return isValid
}
