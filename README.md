# lotherry

Get rich or code trying

## How to try to get rich

1. Download chrome driver for you OS

```bash
wget https://chromedriver.storage.googleapis.com/2.33/chromedriver_mac64.zip
unzip chromedriver_mac64.zip
```
2. Install and configure go on your machine: https://golang.org/dl/
3. Run the `.go` file

```bash
cd lotherry/go/src/app
CHROME_DRIVER_LOCATION=/where/you/downloaded/chrome/driver/file go run lotherry.go
tail -f $(ls | grep winning | tail -n 1)
```

Note: For other config params set [these environment variables](https://github.com/alexandrulaurus/lotherry/blob/master/go/src/app/lotherry.go#L13-L16)

## Docker

```bash
//TODO: Need to install chrome libs
```
