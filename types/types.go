package types

type Payload struct {
	Url      string   `json:"url"`
	Sitemap  string   `json:"sitemap"`
	Limit    int      `json:"limit"`
	Location Location `json:"location"`
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

type Location struct {
	Country  string `json:"country"`
	Language string `json:"language"`
}

type Output interface {
	GenerateOutputFile(filename string, output map[string]interface{})
}
