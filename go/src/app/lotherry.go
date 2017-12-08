package main

import (
    "github.com/fedesog/webdriver"
	"log"
	"time"
	"math/rand"
	"os"
	"encoding/csv"
	"strconv"
)

const DefaultDriverLocation = "chromedriver"
const DefaultScrappingUrl = ""
const DefaultOutputFile = "winning_numbers"
const DefaultRetryAttempts = "20"

func main() {

	chromeDriverLocation := getEnvironmentVariable("CHROME_DRIVER_LOCATION", DefaultDriverLocation)
	chromeDriver := webdriver.NewChromeDriver(chromeDriverLocation)
	err := chromeDriver.Start()
	checkError("Cannot start driver", err)

	desired := webdriver.Capabilities{"Platform": "Linux"}
	required := webdriver.Capabilities{}
	session, err := chromeDriver.NewSession(desired, required)
	checkError("Cannot create session", err)
	url := getEnvironmentVariable("SCRAPPING_URL", DefaultScrappingUrl)
	err = session.Url(url)
	checkError("Cannot fetch url", err)

	timestamp := time.Now().Format("2006-01-02T15:04:05")
	outputFile := getEnvironmentVariable("OUTPUT_FILE", DefaultOutputFile) + "_" + timestamp
	file, err := os.Create(outputFile)
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	csvRow := []string{"privateKey", "address", "balance"}
	err = writer.Write(csvRow)
	checkError("Cannot write to file", err)

	maxRetries, err := strconv.Atoi(getEnvironmentVariable("MAX_RETRIES", DefaultRetryAttempts))
	checkError("Cannot parse max retries", err)

	for 1 == 1 {
		err = session.Url(url)
		checkError("Cannot fetch url", err)

		rows, err := session.FindElements(webdriver.FindElementStrategy("tag name"), "tr")
		checkError("Cannot fetch records", err)

		for rowIndex := range rows {
			rowData, err := rows[rowIndex].FindElements(webdriver.FindElementStrategy("tag name"), "td")
			checkError("Cannot fetch row data", err)

			rowInfo := fetchBalance(0, maxRetries, rowData)
			log.Println(rowInfo["privateKey"] + " " + rowInfo["address"] + " " + rowInfo["balance"])

			if rowInfo["balance"] != "0" {
				csvRow := []string{rowInfo["privateKey"], rowInfo["address"], rowInfo["balance"]}
				err = writer.Write(csvRow)
				checkError("Cannot write to file", err)
				writer.Flush()
			}
		}
		var seconds= int(rand.Intn(30))
		log.Printf("Sleeping for %d", seconds)
		time.Sleep(time.Duration(seconds) * time.Second)
	}

	session.Delete()
	chromeDriver.Stop()
}

func getEnvironmentVariable(key string, defaultValue string) (string) {
	variableValue := os.Getenv(key)
	if len(variableValue) == 0 {
		log.Printf("No env variable specified for %s. Using default: %s", key, defaultValue)
		return defaultValue
	} else {
		log.Printf("Using env variable %s=%s", key, variableValue)
		return variableValue
	}
}

func fetchBalance(attempt int, maxRetries int, rowData []webdriver.WebElement) (map[string]string) {
	rowInfo := parseRow(rowData)
	if attempt < maxRetries && len(rowInfo["balance"]) == 0 {
		log.Printf("Attempt %d: Balance not ready. Sleeping...", attempt)
		time.Sleep(1 * time.Second)
		attempt++
		return fetchBalance(attempt, maxRetries, rowData)
	} else if attempt >= maxRetries && len(rowInfo["balance"]) == 0 {
		log.Printf("Exhausted all %d attempts. Returning", maxRetries)
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