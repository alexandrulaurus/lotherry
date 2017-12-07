package main

import (
    "github.com/fedesog/webdriver"
	"log"
	"time"
	"math/rand"
	"os"
	"encoding/csv"
)

func main() {

	chromeDriver := webdriver.NewChromeDriver("/Users/alexandru/Documents/chromedriver")
	err := chromeDriver.Start()
	if err != nil {
		log.Fatal(err)
	}
	desired := webdriver.Capabilities{"Platform": "Linux"}
	required := webdriver.Capabilities{}
	session, err := chromeDriver.NewSession(desired, required)
	if err != nil {
		log.Fatal(err)
	}
	err = session.Url("http://www.ethersecret.com/random")
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create("/Users/alexandru/Documents/lotherry/winning_numbers")
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	csvRow := []string{"privateKey", "address", "balance"}
	err = writer.Write(csvRow)
	checkError("Cannot write to file", err)

	count := 0
	for count < 10 {
		rows, err := session.FindElements(webdriver.FindElementStrategy("tag name"), "tr")
		if err != nil {
			log.Fatal(err)
		}

		for rowIndex := range rows {
			rowData, rowErr := rows[rowIndex].FindElements(webdriver.FindElementStrategy("tag name"), "td")
			if rowErr != nil {
				log.Fatal(rowErr)
			}

			rowInfo := fetchBalance(0, rowData)
			log.Println(rowInfo["privateKey"] + " " + rowInfo["address"] + " " + rowInfo["balance"])

			if rowInfo["balance"] != "0" {
				csvRow := []string{rowInfo["privateKey"], rowInfo["address"], rowInfo["balance"]}
				err := writer.Write(csvRow)
				checkError("Cannot write to file", err)
			}
		}
		var seconds= int(rand.Intn(30))
		log.Printf("Sleeping for %d", seconds)
		time.Sleep(time.Duration(seconds) * time.Second)
		count++
	}

	session.Delete()
	chromeDriver.Stop()
}

func fetchBalance(attempt int, rowData []webdriver.WebElement) (map[string]string) {
	rowInfo := parseRow(rowData)
	if attempt < 20 && len(rowInfo["balance"]) == 0 {
		log.Println("Balance not ready. Sleeping...")
		time.Sleep(1 * time.Second)
		attempt++
		return fetchBalance(attempt, rowData)
	} else if attempt >= 10 && len(rowInfo["balance"]) == 0 {
		log.Println("Exhausted all attempts. Returning")
		return rowInfo
	} else {
		return rowInfo
	}
}

func parseRow(rowData []webdriver.WebElement)(map[string]string){
	var rowInfo = make(map[string]string)
	for column := range rowData {
		text, textFetchErr := rowData[column].Text()
		if textFetchErr != nil {
			log.Println(textFetchErr)
		}
		if column == 2 {
			rowInfo["balance"] = text
		} else if column == 0 {
			rowInfo["privateKey"] = text
		} else if column == 1 {
			rowInfo["address"] = text
		}
	}
	return rowInfo
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}