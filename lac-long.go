package main

import (
	"encoding/csv"
	"fmt"
	"os"
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
type Route struct {
	name string
	url  string
}

func main() {
	contestants := []string{
		"200222284/daniel-rodas-bautista",
		"200039487/louis-thomas-schreiber",
		"112474537/francis-lessard",
	}
	results := make(map[Route]map[string]bool)
	contestantChallenges := make(map[string][]Route)

	for i := 0; i < len(contestants); i++ {
		ticks := GetTicksForContestant(contestants[i])
		lacLongTicks := GetLacLongTicks(ticks)
		fmt.Printf("lac long ticks for %s: %d\n", contestants[i], len(lacLongTicks))

		for j := 0; j < len(lacLongTicks); j++ {
			route := Route{
				name: lacLongTicks[j].routeName,
				url:  lacLongTicks[j].url,
			}
			if results[route] == nil {
				results[route] = make(map[string]bool)
			}
			results[route][contestants[i]] = true
		}
	}

	for route, result := range results {
		for i := 0; i < len(contestants); i++ {
			if !result[contestants[i]] {
				contestantChallenges[contestants[i]] = append(contestantChallenges[contestants[i]], route)
			}
		}
	}

	for person, routes := range contestantChallenges {
		fmt.Println(person)
		for i := 0; i < len(routes); i++ {
			fmt.Println(routes[i].name)
		}
		fmt.Println("-------------------")
	}
	// TODO: make this into an API
	// TODO: return in order?

}

func GetTicksForContestant(contestantId string) []Tick {
	fileName := fmt.Sprintf("./%s_ticks.csv", strings.Replace(contestantId, "/", "_", -1))
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0)

	if err != nil {
		fmt.Println("There was an error")
	}

	csvReader := csv.NewReader(file)
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
		isValid = tick.leadStyle != "Fell/Hung" && tick.leadStyle != "Pinkpoint"
	}
	if tick.routeType == "Boulder" {
		isValid = tick.style == "Send" || tick.style == "Flash"
	}

	return isValid
}
