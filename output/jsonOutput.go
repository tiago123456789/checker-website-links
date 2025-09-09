package output

import (
	"encoding/json"
	"fmt"
	"os"
)

type JsonOutput struct {
}

func NewJsonOutput() *JsonOutput {
	return &JsonOutput{}
}

func (j *JsonOutput) GenerateOutputFile(filename string, output map[string]interface{}) {
	jsonData, err := json.Marshal(output)
	if err != nil {
		fmt.Errorf("Error marshalling output: %v", err)
		return
	}

	err = os.WriteFile(filename+".json", jsonData, 0644)
	if err != nil {
		fmt.Errorf("Error writing output file: %v", err)
		return
	}
}
