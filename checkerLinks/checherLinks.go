package checkerLinks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/tiago123456789/checker-website-links/types"
)

type CheckerLink struct {
	baseLink                    string
	apiKey                      string
	limit                       int
	links                       []types.Link
	filename                    string
	maxTimeMsAcceptedPerRequest int
	location                    types.Location
	output                      types.Output
}

func NewCheckerLink(
	baseLink string,
	apiKey string,
	limit int,
	filename string,
	maxTimeMsAcceptedPerRequest int,
	location types.Location,
	output types.Output,
) *CheckerLink {
	return &CheckerLink{
		baseLink:                    baseLink,
		apiKey:                      apiKey,
		limit:                       limit,
		filename:                    filename,
		maxTimeMsAcceptedPerRequest: maxTimeMsAcceptedPerRequest,
		location:                    location,
		output:                      output,
	}
}

func (c *CheckerLink) getLinks() []types.Link {
	payload := types.Payload{
		Url:     c.baseLink,
		Sitemap: "include",
		Limit:   c.limit,
	}

	if c.location.Country != "" || c.location.Language != "" {
		payload.Location = c.location
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.firecrawl.dev/v2/map", bytes.NewReader(jsonData))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}

	var response types.Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error unmarshalling response body:", err)
		return nil
	}

	return response.Links
}

func (c *CheckerLink) checkLinks(link types.Link, disableCache bool, timeout int) types.CheckResult {
	if strings.Contains(link.Url, ":443") {
		return types.CheckResult{
			Url:    link.Url,
			Status: http.StatusOK,
		}
	}

	if disableCache {
		link.Url = fmt.Sprintf("%s?v=%d", link.Url, time.Now().Unix())
	}

	startAt := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link.Url, nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	endAt := time.Now()

	timeDuration := endAt.Sub(startAt).Milliseconds()
	if err := ctx.Err(); err != nil {
		return types.CheckResult{
			Url:          link.Url,
			Status:       http.StatusRequestTimeout,
			TimeDuration: timeDuration,
		}
	}

	if err != nil {
		return types.CheckResult{
			Url:          link.Url,
			Status:       http.StatusInternalServerError,
			TimeDuration: timeDuration,
		}
	}

	err = resp.Body.Close()
	if err != nil {
		return types.CheckResult{
			Url:          link.Url,
			Status:       resp.StatusCode,
			TimeDuration: timeDuration,
		}
	}

	return types.CheckResult{
		Url:          link.Url,
		Status:       resp.StatusCode,
		TimeDuration: timeDuration,
	}
}

func (c *CheckerLink) checkTimeSpendToReceiveResponse(
	links []types.CheckResult,
	maxTimeMsAcceptedPerRequest int,
) []types.CheckResult {
	linksSpendMoreThanMaxTime := []types.CheckResult{}
	for _, link := range links {
		if link.TimeDuration > int64(maxTimeMsAcceptedPerRequest) {
			linksSpendMoreThanMaxTime = append(linksSpendMoreThanMaxTime, link)
		}
	}
	return linksSpendMoreThanMaxTime
}

func (c *CheckerLink) Run(disableCache bool, timeout int) ([]types.CheckResult, []types.CheckResult) {
	links := c.getLinks()
	concurrency := make(chan struct{}, 10)

	wg := sync.WaitGroup{}
	var linksOk []types.CheckResult
	var linksError []types.CheckResult

	var lockCounter sync.Mutex

	totalLinks := len(links)
	totalCheckedLinks := 0

	for _, link := range links {
		wg.Add(1)
		concurrency <- struct{}{}
		go func(link types.Link) {
			defer wg.Done()
			result := c.checkLinks(link, disableCache, timeout)
			if result.Status == http.StatusOK {
				linksOk = append(linksOk, result)
			} else {
				linksError = append(linksError, result)
			}
			lockCounter.Lock()
			totalCheckedLinks++
			lockCounter.Unlock()
			fmt.Printf("Checked the link: %s total checked links: %d/%d\n", link.Url, totalCheckedLinks, totalLinks)
			<-concurrency
		}(link)
	}

	wg.Wait()

	output := map[string]interface{}{
		"links_ok":    linksOk,
		"links_error": linksError,
	}

	linksOkSpendMoreThanMaxTime := c.checkTimeSpendToReceiveResponse(linksOk, c.maxTimeMsAcceptedPerRequest)
	linksErrorSpendMoreThanMaxTime := c.checkTimeSpendToReceiveResponse(linksError, c.maxTimeMsAcceptedPerRequest)

	output["links_ok_spend_more_than_max_time"] = linksOkSpendMoreThanMaxTime
	output["links_error_spend_more_than_max_time"] = linksErrorSpendMoreThanMaxTime

	c.output.GenerateOutputFile(c.filename, output)

	if len(linksOkSpendMoreThanMaxTime) > 0 || len(linksErrorSpendMoreThanMaxTime) > 0 {
		linksSpendMoreThanMaxTimeString := ""
		for _, link := range linksOkSpendMoreThanMaxTime {
			linksSpendMoreThanMaxTimeString += fmt.Sprintf("Url: %s | Status: %d | Time duration: %d ms\n", link.Url, link.Status, link.TimeDuration)
		}

		for _, link := range linksErrorSpendMoreThanMaxTime {
			linksSpendMoreThanMaxTimeString += fmt.Sprintf("Url: %s | Status: %d | Time duration: %d ms\n", link.Url, link.Status, link.TimeDuration)
		}

		fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
		fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
		fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")

		fmt.Println("Links spend more than max time accepted per request:\n" + linksSpendMoreThanMaxTimeString)
		os.Exit(1)
	}

	return linksOk, linksError
}
