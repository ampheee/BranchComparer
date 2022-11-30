package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

const apiUrl = "https://rdb.altlinux.org/api/export/branch_binary_packages/"

type BinaryPackage struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Realise string `json:"realise"`
	Arch    string `json:"arch"`
	Disttag string `json:"disttag"`
	Source  string `json:"source"`
	Epoch   int    `json:"epoch"`
}

type structuredData struct {
	Length   int             `json:"length"`
	Packages []BinaryPackage `json:"packages"`
}

type check struct {
	seen int
	BinaryPackage
}

func main() {
	start := time.Now()
	var (
		choice1  string = "p10"
		choice2  string = "sisyphus"
		wg       sync.WaitGroup
		testChan = make(chan structuredData, 2)
	)
	wg.Add(2)
	go getPackages(&wg, choice1, testChan)
	go getPackages(&wg, choice2, testChan)
	wg.Wait()
	set1, set2 := <-testChan, <-testChan
	between := time.Since(start)
	var uniqMap map[string]check
	if set1.Length > set2.Length {
		uniqMap = make(map[string]check, set1.Length)
		for _, gotPackage := range set1.Packages {
			uniqMap[gotPackage.Name] = check{seen: -1, BinaryPackage: gotPackage}
		}
		//} else {
		//	uniqMap = make(map[string]map[string]interface{}, set2.Length)
		//	for _, gotPackage := range set2.Packages {
		//		uniqMap[gotPackage.Name] = map[string]interface{}{
		//			"seen":    -1,
		//			"version": gotPackage.Version,
		//		}
		//	}
	}
	end := time.Since(start)
	fmt.Println(between, end)
}

func getPackages(wg *sync.WaitGroup, choice string, data chan structuredData) {
	defer wg.Done()
	var SD structuredData
	choice1Response, err := http.Get(apiUrl + choice)
	if err != nil && choice1Response.StatusCode != 200 {
		log.Fatalf("cant get api of %s: %s", choice, err)
	}
	choice1Data, err := io.ReadAll(choice1Response.Body)
	if err != nil {
		log.Fatalf("cant read data from body %s response: %s", choice, err)
	}
	err = json.Unmarshal(choice1Data, &SD)
	if err != nil {
		log.Fatalf("cant unmarshal %s: %s", choice, err)
	}
	data <- SD
}
