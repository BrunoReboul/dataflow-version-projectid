// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the 'License');
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an 'AS IS' BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"log"
	"os"
)

type dfAsset struct {
	Name     string `json:"name"`
	Resource struct {
		Parent string `json:"parent"`
		Data   struct {
			Environment struct {
				UserAgent struct {
					Version string `json:"version"`
				}
			} `json:"environment"`
			ProjectID string `json:projectId`
		} `json:"data"`
	} `json:"resource"`
}

func main() {
	var filename string
	var dataflowJob dfAsset
	scannerBufferSizeKiloBytes := 128
	var records [][]string

	flag.StringVar(&filename, "filename", "./df.json", "dataflow job cai export")
	flag.Parse()

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scannerBuffer := make([]byte, scannerBufferSizeKiloBytes*1024)
	scanner.Buffer(scannerBuffer, scannerBufferSizeKiloBytes*1024)
	for scanner.Scan() {
		// fmt.Println(scanner.Text())
		json.Unmarshal(scanner.Bytes(), &dataflowJob)
		record := []string{dataflowJob.Resource.Data.ProjectID,
			dataflowJob.Resource.Data.Environment.UserAgent.Version,
			dataflowJob.Name}
		records = append(records, record)
		// fmt.Printf("%s %s %s\n", dataflowJob.Resource.Data.ProjectID, dataflowJob.Resource.Data.Environment.UserAgent.Version, dataflowJob.Name)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fileOut, err := os.OpenFile("result.csv",
		os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	csvWriter := csv.NewWriter(fileOut)
	csvWriter.WriteAll(records)
	csvWriter.Flush()
}
