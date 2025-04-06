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

	// Green checkmark (ok, nothing to do)
	greenCheckMark := "\033[42m\033[97m \u2713 \033[0m"
	// Yellow exclamation mark (update available)
	yellowExclamation := "\033[43m\033[30m \u0021 \033[0m"
	// Red X (web site not found)
	redX := "\033[91m \u2717 \033[0m"
	// Gray circle (no update check)
	grayCircle := "\033[90m \u25CB \033[0m"
	asciiSymbols := []string{greenCheckMark, yellowExclamation, redX, grayCircle}

	// Define table headers
	list := [][]string{
		{"NAME", "ID", "SIZE", "MODIFIED", "+UPD"},
	}

	// Initialize pflag
	pflag.Parse()

	ctx := context.Background()
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	// Get model list
	modelsPtr, err := client.List(ctx)
	if err != nil {
		log.Fatal(err)
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
		status := asciiSymbols[0] // all ok

		// Calc day diff
		daysDiff := datetools.DaysDifference(model.ModifiedAt, time.Now())
		modified := fmt.Sprintf("%s (%dd)", model.ModifiedAt.Format("02-01-2006"), daysDiff)

		if *chkUpdates {
			// Get details from web page by model name
			ow := scraper.NewOllamaWeb(model.Name)
			err := ow.GetModelInfo()
			if err != nil {
				status = asciiSymbols[2] // error occured
			} else {
				// Compare ID and last modified date
				if (daysDiff > ow.Days) || (digest != ow.Digest) {
					status = asciiSymbols[1] // update found
				}
			}
		} else {
			status = asciiSymbols[3] // nothing done
		}

		// Write table entry
		entry := []string{model.Name, digest, formatbytes.FormatBytes(model.Size), modified, status}
		list = append(list, entry)
	}

	fmt.Println("OllaMMan - the Ollama Model Manager")
	fmt.Println("\n" + markdown.MarkdownTable(list))
}
