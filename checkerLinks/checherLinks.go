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
)

type CheckerLink struct {
	baseLink                    string
	apiKey                      string
	limit                       int
	links                       []Link
	filename                    string
	maxTimeMsAcceptedPerRequest int
}

func NewCheckerLink(
	baseLink string,
	apiKey string,
	limit int,
	filename string,
	maxTimeMsAcceptedPerRequest int,
) *CheckerLink {
	return &CheckerLink{
		baseLink:                    baseLink,
		apiKey:                      apiKey,
		limit:                       limit,
		filename:                    filename,
		maxTimeMsAcceptedPerRequest: maxTimeMsAcceptedPerRequest,
	}
}

func (c *CheckerLink) getLinks() []Link {
	payload := Payload{
		Url:     c.baseLink,
		Sitemap: "include",
		Limit:   c.limit,
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

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error unmarshalling response body:", err)
		return nil
	}

	return response.Links
}

func (c *CheckerLink) checkLinks(link Link, disableCache bool, timeout int) CheckResult {
	if strings.Contains(link.Url, ":443") {
		return CheckResult{
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
		return CheckResult{
			Url:          link.Url,
			Status:       http.StatusRequestTimeout,
			TimeDuration: timeDuration,
		}
	}

	if err != nil {
		return CheckResult{
			Url:          link.Url,
			Status:       http.StatusInternalServerError,
			TimeDuration: timeDuration,
		}
	}

	err = resp.Body.Close()
	if err != nil {
		return CheckResult{
			Url:          link.Url,
			Status:       resp.StatusCode,
			TimeDuration: timeDuration,
		}
	}

	return CheckResult{
		Url:          link.Url,
		Status:       resp.StatusCode,
		TimeDuration: timeDuration,
	}
}

func (c *CheckerLink) generateOutputFile(output map[string]interface{}) {

	jsonData, err := json.Marshal(output)
	if err != nil {
		fmt.Errorf("Error marshalling output: %v", err)
		return
	}

	os.WriteFile(c.filename, jsonData, 0644)

}

func (c *CheckerLink) checkTimeSpendToReceiveResponse(
	links []CheckResult,
	maxTimeMsAcceptedPerRequest int,
) []CheckResult {
	linksSpendMoreThanMaxTime := []CheckResult{}
	for _, link := range links {
		if link.TimeDuration > int64(maxTimeMsAcceptedPerRequest) {
			linksSpendMoreThanMaxTime = append(linksSpendMoreThanMaxTime, link)
		}
	}
	return linksSpendMoreThanMaxTime
}

func (c *CheckerLink) Run(disableCache bool, timeout int) ([]CheckResult, []CheckResult) {
	links := c.getLinks()
	concurrency := make(chan struct{}, 5)

	wg := sync.WaitGroup{}
	var linksOk []CheckResult
	var linksError []CheckResult

	for _, link := range links {
		wg.Add(1)
		concurrency <- struct{}{}
		go func(link Link) {
			defer wg.Done()
			fmt.Println("Checking link: ", link.Url)
			result := c.checkLinks(link, disableCache, timeout)
			if result.Status == http.StatusOK {
				linksOk = append(linksOk, result)
			} else {
				linksError = append(linksError, result)
			}
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

	c.generateOutputFile(output)

	return linksOk, linksError
}
