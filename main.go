package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/brleinad/laclong/laclong"
	"github.com/robfig/cron/v3"
)

func main() {
	downloadCsvTicks()
	c := cron.New()
	c.AddFunc("@daily", downloadCsvTicks)
	c.Start()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", serveTemplate)
	log.Print("Listening on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}

func getContestants() []string {
	return []string{
		"200222284/daniel-rodas-bautista",
		"200039487/louis-thomas-schreiber",
		"112474537/francis-lessard",
	}
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	fp := filepath.Join("static", "index.html")
	tmpl, _ := template.ParseFiles(fp)
	data := laclong.GetLacLongChallenges(getContestants())
	tmpl.Execute(w, data)
}

func downloadCsvTicks() {
	for _, contestant := range getContestants() {
		downloadCsvTicksForContestant(contestant)
	}
}

func downloadCsvTicksForContestant(contestantId string) {
	URL := fmt.Sprintf("https://www.mountainproject.com/user/%s/tick-export", contestantId)
	ticksResponse, err := http.Get(URL)
	if err != nil {
		log.fatal("There was an error")
	}
	fileName := fmt.Sprintf("./%s_ticks.csv", strings.Replace(contestantId, "/", "_", -1))
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer ticksResponse.Body.Close()
	size, err := io.Copy(file, ticksResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.Print("Downloaded a file %s with size %d\n", fileName, size)
}
