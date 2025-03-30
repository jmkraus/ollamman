package scraper

import (
	"log"
	"net/http"
	"ollaman/datetools"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type OllamaWeb struct {
	URL    string
	doc    *goquery.Document
	Digest string
	Days   int16
}

func NewOllamaWeb(modelName string) *OllamaWeb {
	return &OllamaWeb{
		URL: "https://ollama.com/library/" + modelName,
	}
}

func (ow *OllamaWeb) fetchWebPage() {
	// Request the HTML page.
	res, err := http.Get(ow.URL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	ow.doc = doc
}

func (ow *OllamaWeb) GetModelInfo() {

	// Hole die Webseite
	if ow.doc == nil {
		ow.fetchWebPage()
	}

	// Wähle das gewünschte div-Element und extrahiere die Paragraphen
	var paragraphs []string
	ow.doc.Find("#file-explorer section div div div p").Each(func(index int, item *goquery.Selection) {
		text := strings.TrimSpace(item.Text())
		paragraphs = append(paragraphs, text)
	})

	ow.Digest = paragraphs[2][:12]
	date, _ := datetools.ParseRelativeDate(paragraphs[0])
	ow.Days = datetools.DaysDifference(date, time.Now())
}
