package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
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

type AnswerPackages struct {
	onlyInChoice1 []BinaryPackage
	onlyInChoice2 []BinaryPackage
	HigherVersion []BinaryPackage
}

type structuredData struct {
	Length   int             `json:"length"`
	Packages []BinaryPackage `json:"packages"`
}

func main() {
	start := time.Now()
	var (
		choice1  string
		choice2  string
		wg       sync.WaitGroup
		testChan = make(chan structuredData, 2)
	)
	fmt.Println("Please, enter names of your branches:")
	fmt.Scanf("%s %s", &choice1, &choice2)
	wg.Add(2)
	go getPackages(&wg, choice1, testChan)
	go getPackages(&wg, choice2, testChan)
	wg.Wait()
	set1, set2 := <-testChan, <-testChan
	between, secondStart := time.Since(start), time.Now()
	var (
		uniqMap map[string][]*BinaryPackage
		Ans     AnswerPackages
	)
	uniqMap = make(map[string][]*BinaryPackage, set1.Length)
	for i, gotPackage := range set1.Packages {
		uniqMap[gotPackage.Name] = []*BinaryPackage{&set1.Packages[i], nil}
	}
	for i, gotPackage := range set2.Packages {
		if _, exist := uniqMap[gotPackage.Name]; exist {
			uniqMap[gotPackage.Name][1] = &set2.Packages[i]
		} else {
			uniqMap[gotPackage.Name] = []*BinaryPackage{nil, &set2.Packages[i]}
		}
	}
	for _, val := range uniqMap {
		if val[1] == nil && val[0] != nil {
			Ans.onlyInChoice1 = append(Ans.onlyInChoice1, *val[0])
		}
		if val[0] == nil && val[1] != nil {
			Ans.onlyInChoice2 = append(Ans.onlyInChoice2, *val[1])
		}
		if val[0] != nil && val[1] != nil {
			if compareVersion(val[0].Version, val[1].Version) == 1 {
				Ans.HigherVersion = append(Ans.HigherVersion, *val[0])
			}
		}
	}
	wg.Add(3)
	go createFile(&wg, choice1, "Uniq", Ans.onlyInChoice1)
	go createFile(&wg, choice2, "Uniq", Ans.onlyInChoice2)
	go createFile(&wg, "", "Updated", Ans.HigherVersion)
	wg.Wait()
	calculations, end := time.Since(secondStart), time.Since(start)
	fmt.Printf("%s\nTime used: %fs to get response, %fs to calculations\nTotal Time: %fs",
		"Done! Check your folder.", between.Seconds(), calculations.Seconds(), end.Seconds())
}

func getPackages(wg *sync.WaitGroup, choice string, data chan structuredData) {
	defer wg.Done()
	var SD structuredData
	choice1Response, err := http.Get(apiUrl + choice)
	if err != nil || choice1Response.StatusCode != 200 {
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

func createFile(wg *sync.WaitGroup, choice, mode string, ansSlice []BinaryPackage) {
	defer wg.Done()
	file, _ := os.Create(path.Join(choice + mode + ".json"))
	chMarshaled, _ := json.Marshal(ansSlice)
	_, err := file.Write(chMarshaled)
	if err != nil {
		log.Fatalf("Cant write in file %s: %s", choice, err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("%s", err)
		}
	}(file)
}

func compareVersion(version1 string, version2 string) int {
	length1, length2, ptr1, ptr2 := len(version1), len(version2), 0, 0
	for ptr1 < length1 || ptr2 < length2 {
		n1 := 0
		for ptr1 < length1 && '0' <= version1[ptr1] && version1[ptr1] <= '9' {
			n1 = n1*10 + int(version1[ptr1]-'0')
			ptr1++
		}
		n2 := 0
		for ptr2 < length2 && '0' <= version2[ptr2] && version2[ptr2] <= '9' {
			n2 = n2*10 + int(version2[ptr2]-'0')
			ptr2++
		}
		if n1 > n2 {
			return 1
		}
		if n1 < n2 {
			return -1
		}
		ptr1, ptr2 = ptr1+1, ptr2+1
	}
	return 0
}
