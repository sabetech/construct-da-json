package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type CDCBundleInfo struct {
	DedicatedAccount int    `json:"dedicated_account"`
	Description      string `json:"description"`
	UnitType         int    `json:"unit_type"`
	IsActive         bool   `json:"is_active"`
	UnitTypeValue    string `json:"unit_type_value"`
	IsBonus          bool   `json:"is_bonus"`
	IsPostPaid       bool   `json:"is_postpaid"`
}

func main() {
	fmt.Print("Reading from CSV ...")
	f, err := os.Open("da_ids_from_mtn.csv")

	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	var data []CDCBundleInfo
	var count int32
	postPaidDas := [26]int{11, 14, 15, 150, 151, 153, 154, 155, 157, 158, 159, 16, 160, 161, 166, 167, 169, 17, 178, 180, 182, 184, 185, 186, 187, 192}
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		count++
		if count == 1 {
			continue
		}

		if rec[1] == "" {
			continue
		}

		var daId, unit_type int
		var isBonus bool = false
		var isPostpaid bool = false
		var isActive bool = true

		if n, err := strconv.Atoi(rec[0]); err == nil {
			daId = n
		}

		if n, err := strconv.Atoi(rec[2]); err == nil {
			unit_type = n
		}

		if rec[2] == "1" {
			isActive = false
		}

		for _, x := range postPaidDas {
			if x == daId {
				isPostpaid = true
				break
			}
		}

		if strings.Contains(strings.ToLower(rec[1]), "bonus") {
			isBonus = true
		}

		if strings.Contains(strings.ToLower("MyMTN_50MB"), strings.ToLower(strings.Trim(rec[1], " "))) {
			isBonus = true
		}

		var unitValueType string
		switch unit_type {
		case 0:
			unitValueType = "Time sec"
			break
		case 1:
			unitValueType = "Monetary"
		case 6:
			unitValueType = "Volume Bytes"
			break
		case 5:
			unitValueType = "SMS"
			break
		}

		cdcbundleInfo := CDCBundleInfo{
			DedicatedAccount: daId,
			Description:      strings.Trim(rec[1], " "),
			UnitType:         unit_type,
			IsActive:         isActive,
			UnitTypeValue:    unitValueType,
			IsBonus:          isBonus,
			IsPostPaid:       isPostpaid,
		}

		data = append(data, cdcbundleInfo)

	}

	file, _ := json.MarshalIndent(data, "", " ")

	_ = ioutil.WriteFile("test.json", file, 0644)

}
