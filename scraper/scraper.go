package scraper

import (
	"fmt"
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

func (ow *OllamaWeb) fetchWebPage() error {

	// Request the HTML page
	res, err := http.Get(ow.URL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("%v", res.StatusCode)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		ow.doc = nil
		return err
	} else {
		ow.doc = doc
		return nil
	}
}

func (ow *OllamaWeb) GetModelInfo() error {

	err := ow.fetchWebPage()
	if err != nil {
		return err
	}

	// Extract data from DOM tree
	var paragraphs []string
	ow.doc.Find("#file-explorer section div div div p").Each(func(index int, item *goquery.Selection) {
		text := strings.TrimSpace(item.Text())
		paragraphs = append(paragraphs, text)
	})

	// Format and store extracted values
	ow.Digest = paragraphs[2][:12]
	date, _ := datetools.ParseRelativeDate(paragraphs[0])
	ow.Days = datetools.DaysDifference(date, time.Now())
	return nil
}
