package main

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/brleinad/laclong/laclong"
	"github.com/robfig/cron/v3"
)

//go:embed static/*
var static embed.FS

func main() {
	downloadCsvTicks()
	c := cron.New()
	c.AddFunc("@daily", downloadCsvTicks)
	c.Start()

	fSys, err := fs.Sub(static, "static")
	if err != nil {
		log.Fatal("Failed to serve static ", err)
	}
	fs := http.FileServer(http.FS(fSys))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", serveTemplate)
	log.Print("Listening on :8080...")
	err = http.ListenAndServe(":8080", nil)
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
	// fp := filepath.Join(".", "static", "index.html")
	// tmpl, err := template.ParseFiles(fp)

	var tmpl = template.Must(template.ParseFS(static, "static/*"))
	data := laclong.GetLacLongChallenges(getContestants())
	log.Println("YOYO: ", data)
	err := tmpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Fatal("Failed to execute template")
	}
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
		log.Fatal("There was an error")
	}
	fileName := laclong.GetFileName(contestantId)
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
	log.Printf("Downloaded a file %s with size %d\n", fileName, size)
}
