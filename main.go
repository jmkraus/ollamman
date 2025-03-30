package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"ollaman/datetools"
	"ollaman/formatbytes"
	"ollaman/markdown"
	"ollaman/scraper"
	"sort"
	"time"

	"github.com/ollama/ollama/api"
)

func main() {

	// 1. Update verfügbar:
	//   		- 🔶 Orangefarbene Raute: `\U0001F536`
	//
	// 2. Datei nicht mehr gefunden:
	//   		- ❌ Rotes Kreuz: `\U0000274C`
	//
	// 3. Neueste Version (Ihr bereits gewähltes):
	//   		- ✅ Weißer Haken auf grünem Hintergrund: `\U00002705`

	symbols := []string{"  \U00002705", "  \U0001F536", "  \U0000274C"}
	list := [][]string{
		{"NAME", "ID", "SIZE", "MODIFIED", "UPDATE"},
	}

	// Initialize flag
	sortByName := flag.Bool("on", false, "Sort by name")
	sortByDate := flag.Bool("od", false, "Sort by date")
	flag.Parse()

	ctx := context.Background()
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	// Get model list
	modelsPtr, err := client.List(ctx)
	if err != nil {
		// Error handling
		fmt.Println("Error:", err)
		return
	}

	if *sortByName {
		// Sortieren nach Name
		sort.Slice(modelsPtr.Models, func(i, j int) bool {
			return modelsPtr.Models[i].Name < modelsPtr.Models[j].Name
		})
	}

	// Falls das Flag gesetzt ist, nach Datum sortieren
	if *sortByDate {
		sort.Slice(modelsPtr.Models, func(i, j int) bool {
			return modelsPtr.Models[i].ModifiedAt.Before(modelsPtr.Models[j].ModifiedAt)
		})
	}

	// Iterate over models in ListResponse
	for _, model := range modelsPtr.Models {

		status := symbols[0]
		// fmt.Printf("Name: %s, Model: %s, ModifiedAt: %s, Size: %d, Digest: %.12s\n",
		// model.Name, model.Model, model.ModifiedAt.Format(time.RFC3339), model.Size, model.Digest)

		// Get details from web page by model name
		myOllama := scraper.NewOllamaWeb(model.Name)
		myOllama.GetModelInfo()

		// Compare ID and last modified date
		days := datetools.DaysDifference(model.ModifiedAt, time.Now())
		digest := model.Digest[:12]
		if (days > myOllama.Days) || (digest != myOllama.Digest) {
			status = symbols[1] // update found
		}
		// fmt.Printf("Details: %+v\n\n", model.Details)

		entry := []string{model.Name, digest, formatbytes.FormatBytes(model.Size), model.ModifiedAt.Format(time.RFC3339), status}
		list = append(list, entry)
	}

	fmt.Println("OllaMan - the Ollama Update Manager")
	fmt.Println("\n" + markdown.MarkdownTable(list))
}
