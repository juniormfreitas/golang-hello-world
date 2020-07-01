package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Author string : Author's name
const Author = "J. de Freitas"

// Version int : App's Version
const Version = 1.02

// Delay int : One second delay
const Delay = 1

// HTTPOK int : Http response when the resquest resturns properly
const HTTPOK = 200

func main() {

	introduction()

	for {
		showIntroMenu()

		switch getUserInput() {
		case 1:
			checkStatus()
		case 2:
			retreiveLogs()
		case 3:
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println("Unknown command")
			os.Exit(-1)
		}
	}
}

func introduction() {
	fmt.Println("-----------------------")
	fmt.Println("Http Response Checker")
	fmt.Println("Author: ", Author)
	fmt.Println("Version: ", Version)
	fmt.Println("-----------------------")
}

func getUserInput() int {
	var userInput int

	fmt.Scan(&userInput)
	fmt.Println("Your input is", userInput)

	return userInput
}

func showIntroMenu() {
	fmt.Println("1 - Check status")
	fmt.Println("2 - Show logs")
	fmt.Println("3 - Exit")
}

func retreiveLogs() {
	fmt.Println("Retrieving logs...")

	file, err := ioutil.ReadFile("logs.txt")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(file))
	}
}

func checkStatus() {
	fmt.Println("Checking status...")

	sites := showSites()

	for _, site := range sites {
		testSite(site)
		time.Sleep(Delay * time.Second)
	}
}

func showSites() []string {
	var sites []string

	filename := "sites.txt"

	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		fmt.Println(filename, "does not exist.")
	} else {
		file, err := os.Open(filename)

		if err != nil {
			fmt.Println("Error to open the file: ", err)
		}

		reader := bufio.NewReader(file)

		for {
			row, fileErr := reader.ReadString('\n')
			row = strings.TrimSpace(row)

			sites = append(sites, row)

			if fileErr == io.EOF {
				break
			}
		}

		file.Close()
	}

	return sites
}

func testSite(site string) {
	httpResp, err := http.Get(site)

	fmt.Println("Testing ", site)

	if err != nil {
		fmt.Println(err)
	} else {
		if &httpResp.StatusCode != nil {
			if httpResp.StatusCode == HTTPOK {
				fmt.Println("Status: ", httpResp.StatusCode, " - Successfully loaded")
			} else {
				fmt.Println("Status: ", httpResp.StatusCode, " We are unable to load this website")
			}
			registerLog(site, httpResp.StatusCode, httpResp.Status)
		} else {
			fmt.Println("Something went wrong")
		}
	}
}

func registerLog(site string, statusCode int, status string) {
	file, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println(err)
	} else {
		file.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - ")
		file.WriteString(site + " - status - " + strconv.FormatInt(int64(statusCode), 10) + " - " + status + "\n")
	}

	file.Close()
}
