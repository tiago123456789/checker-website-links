package output

import (
	"fmt"
	"os"

	"github.com/tiago123456789/checker-website-links/types"
)

type CsvOutput struct {
}

func NewCsvOutput() *CsvOutput {
	return &CsvOutput{}
}

func (j *CsvOutput) GenerateOutputFile(filename string, output map[string]interface{}) {
	csvData := "url,status,time_duration_ms,excedded_max_time\n"

	// output := map[string]interface{}{
	// 	"links_ok":    linksOk,
	// 	"links_error": linksError,
	// }

	linksOk := output["links_ok"].([]types.CheckResult)
	for i := 0; i < len(linksOk); i++ {
		link := linksOk[i]
		csvData += fmt.Sprintf("%s, %d, %d, false\n", link.Url, link.Status, link.TimeDuration)
	}

	linksError := output["links_error"].([]types.CheckResult)
	for i := 0; i < len(linksError); i++ {
		link := linksError[i]
		csvData += fmt.Sprintf("%s, %d, %d, false\n", link.Url, link.Status, link.TimeDuration)
	}

	linksOkSpendMoreThanMaxTime := output["links_ok_spend_more_than_max_time"].([]types.CheckResult)
	for i := 0; i < len(linksOkSpendMoreThanMaxTime); i++ {
		link := linksOkSpendMoreThanMaxTime[i]
		csvData += fmt.Sprintf("%s, %d, %d, true\n", link.Url, link.Status, link.TimeDuration)
	}

	linksErrorSpendMoreThanMaxTime := output["links_error_spend_more_than_max_time"].([]types.CheckResult)
	for i := 0; i < len(linksErrorSpendMoreThanMaxTime); i++ {
		link := linksErrorSpendMoreThanMaxTime[i]
		csvData += fmt.Sprintf("%s, %d, %d, true\n", link.Url, link.Status, link.TimeDuration)
	}

	err := os.WriteFile(filename+".csv", []byte(csvData), 0644)
	if err != nil {
		fmt.Errorf("Error writing output file: %v", err)
		return
	}
}
