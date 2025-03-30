package main

import (
	"context"
	"fmt"
	"github.com/spf13/pflag"
	"log"
	"ollaman/datetools"
	"ollaman/formatbytes"
	"ollaman/markdown"
	"ollaman/scraper"
	"sort"
	"time"

	"github.com/ollama/ollama/api"
)

var sortByDate = pflag.BoolP("order-date", "d", false, "Sort by date (oldest first)")
var sortByName = pflag.BoolP("order-name", "n", false, "Sort alphabetically by name")

func main() {

	// 1. Update available:
	//    - 🔶 : `\U0001F536`
	//
	// 2. Model not found:
	//    - ❌ : `\U0000274C`
	//
	// 3. Is already latest version:
	//    - ✅ : `\U00002705`

	symbols := []string{"  \U00002705", "  \U0001F536", "  \U0000274C"}
	list := [][]string{
		{"NAME", "ID", "SIZE", "MODIFIED", "UPDATE"},
	}

	// Initialize flag
	pflag.Parse()

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
		sort.Slice(modelsPtr.Models, func(i, j int) bool {
			return modelsPtr.Models[i].Name < modelsPtr.Models[j].Name
		})
	}

	if *sortByDate {
		sort.Slice(modelsPtr.Models, func(i, j int) bool {
			return modelsPtr.Models[i].ModifiedAt.Before(modelsPtr.Models[j].ModifiedAt)
		})
	}

	// Iterate over models in ListResponse
	for _, model := range modelsPtr.Models {

		status := symbols[0]

		// Get details from web page by model name
		ow := scraper.NewOllamaWeb(model.Name)
		ow.GetModelInfo()

		// Compare ID and last modified date
		days := datetools.DaysDifference(model.ModifiedAt, time.Now())
		digest := model.Digest[:12]
		if (days > ow.Days) || (digest != ow.Digest) {
			status = symbols[1] // update found
		}
		// fmt.Printf("Details: %+v\n\n", model.Details)

		entry := []string{model.Name, digest, formatbytes.FormatBytes(model.Size), model.ModifiedAt.Format(time.RFC3339), status}
		list = append(list, entry)
	}

	fmt.Println("OllaMan - the Ollama Model Manager")
	fmt.Println("\n" + markdown.MarkdownTable(list))
}
