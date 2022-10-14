package laclong

import (
	"encoding/csv"
	"log"
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
	Name string
	Url  string
}

func GetLacLongChallenges(contestants []string) map[string][]Route {
	results := make(map[Route]map[string]bool)
	contestantChallenges := make(map[string][]Route)

	for i := 0; i < len(contestants); i++ {
		ticks := GetTicksForContestant(contestants[i])
		lacLongTicks := GetLacLongTicks(ticks)
		// fmt.Printf("lac long ticks for %s: %d\n", contestants[i], len(lacLongTicks))

		for j := 0; j < len(lacLongTicks); j++ {
			route := Route{
				Name: lacLongTicks[j].routeName,
				Url:  lacLongTicks[j].url,
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
				contestantChallenges[GetContestantName(contestants[i])] = append(contestantChallenges[GetContestantName(contestants[i])], route)
			}
		}
	}

	return contestantChallenges
}

func GetContestantName(contestantId string) string {
	name := strings.Split(contestantId, "/")[1]
	name = strings.Replace(name, "-", " ", -1)
	name = strings.ToUpper(name)
	return name
}

func GetTicksForContestant(contestantId string) []Tick {
	// fileName := fmt.Sprintf("/data/%s_ticks.csv", strings.Replace(contestantId, "/", "_", -1))
	file, err := os.OpenFile(GetFileName(contestantId), os.O_RDONLY, 0)

	if err != nil {
		log.Fatal("There was an error opening file ", contestantId, " ERROR:", err)
	}

	csvReader := csv.NewReader(file)
	csvTicks, err := csvReader.ReadAll()

	if err != nil {
		log.Fatal("There was an error getting ticks from CSV file ", contestantId, err)
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
