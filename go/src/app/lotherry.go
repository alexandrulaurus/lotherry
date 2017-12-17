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


//CHROME_DRIVER_LOCATION
const DefaultDriverLocation = "chromedriver"
//SCRAPPING_URL
const DefaultScrappingUrl = ""
//OUTPUT_FILE
const DefaultOutputFile = "winning_numbers"
//MAX_RETRIES
const DefaultRetryAttempts = "20"

func main() {

	chromeDriverLocation := getEnvironmentVariable("CHROME_DRIVER_LOCATION", DefaultDriverLocation)
	chromeDriver := webdriver.NewChromeDriver(chromeDriverLocation)
	err := chromeDriver.Start()
	checkError("Cannot start driver", err)

	chromeOptions := make(map[string][]string)
	chromeOptions["args"] = []string{"--headless", "--disable-gpu", "--window-size=1280,800"}
	desired := webdriver.Capabilities{"Platform": "Linux", "chromeOptions" : chromeOptions}

	required := webdriver.Capabilities{}
	session, err := chromeDriver.NewSession(desired, required)
	checkError("Cannot create session", err)
	time.Sleep(10 * time.Second)
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
		log.Println("Fetching next page")
		err = session.Url(url)
		//checkError("Cannot fetch url", err)
		session.AcceptAlert()

		rows, err := session.FindElements(webdriver.FindElementStrategy("tag name"), "tr")
		checkError("Cannot fetch records", err)

		for rowIndex := 1; rowIndex < len(rows); rowIndex++ {
			rowData, err := rows[rowIndex].FindElements(webdriver.FindElementStrategy("tag name"), "td")
			checkError("Cannot fetch row data", err)

			for attempt := 1; attempt < maxRetries; attempt++ {
				rowInfo := parseRow(rowData)
				balance := rowInfo["balance"]

				if len(balance) == 0 {
					log.Printf("Attempt %d of %d. Balance not ready. Retrying", attempt, maxRetries)
					continue
				}
				if len(balance) != 0 {
					log.Printf("%s %s %s", rowInfo["privateKey"], rowInfo["address"], balance)
					if balance != "0" {
						log.Printf("Bingo: %s %s %s %s", rowInfo["privateKey"], rowInfo["address"], rowInfo["balance"], balance)
						csvRow := []string{rowInfo["privateKey"], rowInfo["address"], balance}
						err = writer.Write(csvRow)
						checkError("Cannot write to file", err)
						writer.Flush()
					}
					break
				}
				if attempt == maxRetries - 1 {
					log.Printf("Exhausted all attempts for row %d", rowIndex)
				}
			}
		}
		var seconds = random(5, 10)
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

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max - min) + min
}

