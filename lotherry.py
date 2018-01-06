import logging
import os
import time

from selenium import webdriver
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.by import By
from selenium.common.exceptions import NoAlertPresentException

logging.basicConfig(level=logging.INFO,
        format='%(asctime)s [%(levelname)s] (%(threadName)-10s) %(message)s',
        )


MAX_ATTEMPTS = int(os.getenv('MAX_ATTEMPTS', 50))
OUTPUT_FILE = os.getenv('OUTPUT_FILE', 'winning_numbers')

def acceptAlert(driver):
    try:
        time.sleep(0.02)
        alert = driver.switch_to.alert
        alert.accept()
        logging.info("Alert accepted")
    except NoAlertPresentException:
        a = 1

def main():
    logging.info("Configuring options...")
    options = webdriver.ChromeOptions()
    options.add_argument('headless')
    options.add_argument('window-size=1200x600')
    logging.info("Starting webdriver...")
    driver = webdriver.Chrome(chrome_options=options)

    output = open('file_{}_{}'.format(OUTPUT_FILE, int(time.time())), 'w')
    output.write("      key     |       address     |       balance     \n")
    output.flush()

    while True:
        logging.info("Fetching page...")
        driver.get("http://www.ethersecret.com/random")
        acceptAlert(driver)
        rows = driver.find_elements(By.TAG_NAME, "tr")
        acceptAlert(driver)
        for rowIdx, row in enumerate(rows):
            if (rowIdx == 0):
                continue
            columns = row.find_elements(By.TAG_NAME, "td")
            acceptAlert(driver)
            for attempt in range(1, MAX_ATTEMPTS):
                logging.info("Attempt %d of %d", attempt, MAX_ATTEMPTS)
                rowInfo = {}
                for idx, val in enumerate(columns):
                    acceptAlert(driver)
                    if (idx == 0):
                        rowInfo['key'] = val.text
                        acceptAlert(driver)
                    if (idx == 1):
                        rowInfo['address'] = val.text
                        acceptAlert(driver)
                    if (idx == 2):
                        rowInfo['balance'] = val.text
                        acceptAlert(driver)
                try:
                    balance = rowInfo['balance']
                    if (balance == "" or balance == None):
                        continue
                    logging.info("%s | %s | %s", rowInfo['key'], rowInfo['address'], rowInfo['balance'])
                    if (balance != "0"):
                        logging.info("Bingo: %s | %s | %s", rowInfo['key'], rowInfo['address'], rowInfo['balance'])
                        output.write(rowInfo['key'] + " " + rowInfo['address']
                                + " " + rowInfo['balance'] + "\n")
                        output.flush()
                    break
                except KeyError:
                    time.sleep(1)
                    continue
                logging.info("%s | %s | %s", rowInfo['key'], rowInfo['address'], rowInfo['balance'])

if __name__ == "__main__":
    main()
