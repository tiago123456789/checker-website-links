package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tiago123456789/checker-website-links/checkerLinks"
)

func main() {
	apiKey := flag.String("api-key", "", "API key (required)")
	linkToCheck := flag.String("link", "", "Link to check (required)")
	limitLinks := flag.Int("limit", 100, "Limit links to check (optional)")
	disableCache := flag.Bool("disable-cache", false, "Disable cache (optional)")
	timeout := flag.Int("timeout", 5, "Timeout (optional)")
	filenameOutput := flag.String("output", "output.json", "Output file (optional)")
	maxTimeMsAcceptedPerRequest := flag.Int("max-time-ms-accepted", 5000, "Max time ms accepted (optional)")
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

	checkerLinks := checkerLinks.NewCheckerLink(
		*linkToCheck,
		*apiKey,
		*limitLinks,
		*filenameOutput,
		*maxTimeMsAcceptedPerRequest,
	)

	checkerLinks.Run(*disableCache, *timeout)

	fmt.Println("All links checked successfully")
	os.Exit(0)

}
