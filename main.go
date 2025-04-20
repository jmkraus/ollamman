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
	//const greenCheckMark = "\033[42m\033[97m \u2713 \033[0m"
	const greenCheckMark = "\033[42;97m \u2714 \033[0m"
	// Yellow exclamation mark (update available)
	//const yellowExclamation = "\033[43m\033[30m \u0021 \033[0m"
	const yellowExclamation = "\033[43;30m \u26A0 \033[0m"
	// Red X (web site not found)
	const redX = "\033[91m \u2717 \033[0m"
	// Gray circle (no update check)
	const grayCircle = "\033[90m \u25CB \033[0m"

	// Define table headers
	list := [][]string{
		{"NAME", "ID", "SIZE", "MODIFIED", "CAPABILITIES", "+UPD"},
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

		capabilities := ""
		info, err := func() (*api.ShowResponse, error) {
			showReq := &api.ShowRequest{Name: model.Name}
			info, err := client.Show(ctx, showReq)
			return info, err
		}()
		if err != nil {
			fmt.Println(err)
		} else {
			capabilities = fmt.Sprintf("%v", info.Capabilities)
		}

		digest := model.Digest[:12]
		status := greenCheckMark // all ok

		// Calc day diff
		daysDiff := datetools.DaysDifference(model.ModifiedAt, time.Now())
		modified := fmt.Sprintf("%s (%dd)", model.ModifiedAt.Format("02-01-2006"), daysDiff)

		if *chkUpdates {
			// Get details from web page by model name
			ow := scraper.NewOllamaWeb(model.Name)
			err := ow.GetModelInfo()
			if err != nil {
				status = redX // error occured
			} else {
				// Compare ID and last modified date
				if digest != ow.Digest {
					status = yellowExclamation // update found
				}
			}
		} else {
			status = grayCircle // nothing done
		}

		// Write table entry
		entry := []string{model.Name, digest, formatbytes.FormatBytes(model.Size), modified, capabilities, status}
		list = append(list, entry)
	}

	fmt.Println("OllaMMan - the Ollama Model Manager")
	fmt.Println("\n" + markdown.MarkdownTable(list))
}
