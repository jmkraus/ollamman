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
var chkUpdates = pflag.BoolP("check-updates", "c", false, "Check models for updates")

func main() {

	// 1. Update available:
	//    - 🔶 : `\U0001F536`
	//
	// 2. Model not found:
	//    - ❌ : `\U0000274C`
	//
	// 3. Is already latest version:
	//    - ✅ : `\U00002705`
	//
	// 4. No update check done
	//    - ⚪ : `\U000026AA`
	// symbols := []string{"  \U00002705", "  \U0001F536", "  \U0000274C", "  \U000026AA"}

	// Weißes Häkchen auf grünem Grund (Neueste Version)
	checkMarkGreen := "\033[42m\033[97m \u2713 \033[0m"
	// Schwarzes Ausrufezeichen auf gelbem Grund (Update verfügbar)
	exclamationYellow := "\033[43m\033[30m \u0021 \033[0m"
	// Rotes X ohne Hintergrundfarbe (Datei nicht gefunden)
	redX := "\033[91m \u2717 \033[0m"
	// Neutrales Symbol (keine Update-Prüfung) - grauer Kreis
	grayCircle := "\033[90m \u25CB \033[0m"
	asciiSymbols := []string{checkMarkGreen, exclamationYellow, redX, grayCircle}

	list := [][]string{
		{"NAME", "ID", "SIZE", "MODIFIED", "+UPD"},
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

		if 1 == 0 {
			info, err := func() (*api.ShowResponse, error) {
				showReq := &api.ShowRequest{Name: model.Name}
				info, err := client.Show(ctx, showReq)
				return info, err
			}()
			if err != nil {
				fmt.Println(err)
			}

			fmt.Printf("%v\n", info.Capabilities)
		}

		digest := model.Digest[:12]
		status := asciiSymbols[0]

		// Calc day diff
		daysDiff := datetools.DaysDifference(model.ModifiedAt, time.Now())

		if *chkUpdates {
			// Get details from web page by model name
			ow := scraper.NewOllamaWeb(model.Name)
			err := ow.GetModelInfo()
			if err != nil {
				status = asciiSymbols[2]
			} else {
				// Compare ID and last modified date
				if (daysDiff > ow.Days) || (digest != ow.Digest) {
					status = asciiSymbols[1] // update found
				}
				// fmt.Printf("Details: %+v\n\n", model.Details)
			}
		} else {
			status = asciiSymbols[3]
		}

		// Write table entry
		modified := fmt.Sprintf("%s (%dd)", model.ModifiedAt.Format("02-01-2006"), daysDiff)
		entry := []string{model.Name, digest, formatbytes.FormatBytes(model.Size), modified, status}
		list = append(list, entry)
	}

	fmt.Println("OllaMan - the Ollama Model Manager")
	fmt.Println("\n" + markdown.MarkdownTable(list))
}
