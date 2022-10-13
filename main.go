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

	"github.com/robfig/cron/v3"
)

func main() {
	downloadCsvTicks()
	c := cron.New()
	c.AddFunc("@daily", downloadCsvTicks)
	c.Start()

	http.HandleFunc("/", serveTemplate)
	log.Print("Listening on :3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}

}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	fp := filepath.Join("static", "index.html")
	tmpl, _ := template.ParseFiles(fp)
	data := struct {
		Items []string
	}{
		Items: []string{
			"bob",
			"bab",
		},
	}
	tmpl.Execute(w, data)
}

func downloadCsvTicks() {
	contestants := []string{
		"200222284/daniel-rodas-bautista",
		"200039487/louis-thomas-schreiber",
		"112474537/francis-lessard",
	}
	for _, contestant := range contestants {
		downloadCsvTicksForContestant(contestant)
	}
}

func downloadCsvTicksForContestant(contestantId string) {
	URL := fmt.Sprintf("https://www.mountainproject.com/user/%s/tick-export", contestantId)
	ticksResponse, err := http.Get(URL)
	if err != nil {
		fmt.Println("There was an error")
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
	fmt.Printf("Downloaded a file %s with size %d\n", fileName, size)
}
