package checkerLinks

type Payload struct {
	Url     string `json:"url"`
	Sitemap string `json:"sitemap"`
	Limit   int    `json:"limit"`
}

type Link struct {
	Url         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Response struct {
	Success bool   `json:"success"`
	Links   []Link `json:"links"`
}

type CheckResult struct {
	Url          string `json:"url"`
	Status       int    `json:"status"`
	TimeDuration int64  `json:"time_duration_ms"`
}
