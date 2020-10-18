package main

import (
	"fmt"
)

// Function to print the help message
func helpMessage() {
	fmt.Printf("-------- Go HTTP Client --------\n")
	fmt.Printf("| Usage:\n|\n")
	fmt.Printf("| ./request --url <url>\n|\n")
	fmt.Printf("| Or:\n|\n")
	fmt.Printf("| ./request --url <url> --profile <number of requests>\n|\n")
	fmt.Printf("| If you would like to profile a website\n|\n")
	fmt.Printf("| Examples:\n")
	fmt.Printf("| ./request --url www.google.com\n")
	fmt.Printf("| ./request --url www.google.com --profile 10\n")
	fmt.Printf("--------------------------------\n")
}

// Function to a message when invalid arguments are passed
func invalidMessage() {
	fmt.Printf("-------- Go HTTP Client --------\n")
	fmt.Printf("| Invalid arguments\n")
	fmt.Printf("| Please use --help for more information\n")
	fmt.Printf("--------------------------------\n")
}

// Function to print when an error occurs
func printError(err string) {
	fmt.Printf("-------- Go HTTP Client --------\n")
	fmt.Printf("| An error has occured:\n")
	fmt.Printf("| %s\n", err)
	fmt.Printf("--------------------------------\n")
}

// Function to print the body of a response when not profiling
func printBody(body string) {
	fmt.Printf("-------- Go HTTP Client --------\n")
	fmt.Printf("| Body (May be many lines): \n")
	fmt.Printf("%s", body)
}

// Function to print all the data collected with the profile struct
func printProfile(profile profile) {
	fmt.Printf("------------ Profile Results ------------\n")
	fmt.Printf("| URL: %s\n", profile.url)
	fmt.Printf("| Number of Requests: %d\n", profile.iterations)
	fmt.Printf("| Fastest Response: %d\n", profile.fastest)
	fmt.Printf("| Slowest Response: %d\n", profile.slowest)
	fmt.Printf("| Mean Response Time: %d\n", profile.mean)
	fmt.Printf("| Median Response Time: %d\n", profile.median)
	fmt.Printf("| Percentage of Requests That Errored: %f\n", profile.percentageError)
	fmt.Printf("| Errors (If Any): \n")

	fmt.Printf("| Largest Response: %d\n", profile.largestResponse)
	fmt.Printf("| Smallest Response: %d\n", profile.smallestResponse)
	fmt.Printf("----------------------------------------\n")
}
