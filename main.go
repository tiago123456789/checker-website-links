package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tiago123456789/checker-website-links/checkerLinks"
	"github.com/tiago123456789/checker-website-links/output"
	"github.com/tiago123456789/checker-website-links/types"
)

func main() {
	apiKey := flag.String("api-key", "", "API key (required)")
	linkToCheck := flag.String("link", "", "Link to check (required)")
	limitLinks := flag.Int("limit", 100, "Limit links to check (optional)")
	disableCache := flag.Bool("disable-cache", false, "Disable cache (optional)")
	timeout := flag.Int("timeout", 5, "Timeout (optional)")
	filenameOutput := flag.String("output", "output", "Output file (optional)")
	maxTimeMsAcceptedPerRequest := flag.Int("max-time-ms-accepted", 5000, "Max time ms accepted (optional)")
	country := flag.String("country", "", "Country (optional)")
	language := flag.String("language", "", "Language (optional)")
	outputFormat := flag.String("output-format", "json", "Output format (optional)")

	flag.Parse()

	if *apiKey == "" {
		fmt.Println("Error: -api-key is required")
		flag.Usage()
		os.Exit(1)
	}

	if *linkToCheck == "" {
		fmt.Println("Error: -link is required")
		flag.Usage()
		os.Exit(1)
	}

	if (*country != "" && *language == "") || (*country == "" && *language != "") {
		fmt.Println("Error: if provide -country or -language, both are required")
		flag.Usage()
		os.Exit(1)
	}

	if *outputFormat != "json" && *outputFormat != "csv" {
		fmt.Println("Error: -output-format must be 'json' or 'csv'")
		flag.Usage()
		os.Exit(1)
	}

	outputFormatMap := map[string]types.Output{
		"json": output.NewJsonOutput(),
		"csv":  output.NewCsvOutput(),
	}

	output := outputFormatMap[*outputFormat]

	checkerLinks := checkerLinks.NewCheckerLink(
		*linkToCheck,
		*apiKey,
		*limitLinks,
		*filenameOutput,
		*maxTimeMsAcceptedPerRequest,
		types.Location{
			Country:  *country,
			Language: *language,
		},
		output,
	)

	checkerLinks.Run(*disableCache, *timeout)

	fmt.Println("All links checked successfully")
	os.Exit(0)

}
