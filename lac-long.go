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

func main() {
	ticks := GetTicks()
	// fmt.Print(ticks)
	numberOfTicks := GetLacLongTicks(ticks)
	fmt.Println(numberOfTicks)
}

func GetTicks() [][]string {
	const URL = "https://www.mountainproject.com/user/200222284/daniel-rodas-bautista/tick-export"
	ticksResponse, err := http.Get(URL)

	if err != nil {
		fmt.Println("There was an error")
	}

	csvReader := csv.NewReader(ticksResponse.Body)
	ticks, err := csvReader.ReadAll()

	if err != nil {
		fmt.Println("There was an error")
	}

	return ticks
}

func GetLacLongTicks(ticks [][]string) int {
	numberOfLacLongTicks := 0
	for i := 0; i < len(ticks); i++ {
		if IsInLacLong(ticks[i]) && IsValidSend(ticks[i]) {
			fmt.Println(ticks[i][RouteIndex])
			numberOfLacLongTicks++
		}
	}
	return numberOfLacLongTicks
}

func IsInLacLong(tick []string) bool {
	const lacLongAreaName = "International > North America > Canada > Quebec > 03. Quebec City, Charlevoix, Portneuf > Lac Long"
	return strings.Contains(tick[AreaIndex], lacLongAreaName)
}

func IsValidSend(tick []string) bool {
	isValid := false
	if tick[StyleIndex] == "Lead" && (tick[LeadStyleIndex] == "Redpoint" || tick[LeadStyleIndex] == "Onsight" || tick[LeadStyleIndex] == "Flash") {
		isValid = true
	}
	if tick[RouteTypeIndex] == "Boulder" && tick[StyleIndex] == "Send" {
		isValid = true
	}

	return isValid
}
