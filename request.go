package main

import (
	"crypto/tls"
	"io/ioutil"
	"math"
	"net"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Struct for parsed URL
type request struct {
	profile    bool
	iterations int
	url        string
	path       string
}

// Struct for profile results
type profile struct {
	url              string
	iterations       int
	fastest          int
	slowest          int
	mean             int
	median           int
	percentageError  float64
	errors           []int
	smallestResponse int
	largestResponse  int
}

// Struct for individual responses
type response struct {
	code int
	time int
	size int
}

func parseUrlPath(url string) (string, string) {
	// Strip http or https from url
	regex := regexp.MustCompile("^(http://|https://)")
	url = regex.ReplaceAllString(strings.ToLower(url), "")

	// Get location of '/' from url, and if it exists
	// separate the string at that location
	parsedUrl, path := url, "/"
	separator := strings.Index(url, "/")

	if separator != -1 {
		path = parsedUrl[separator:]
		parsedUrl = parsedUrl[:separator]
	}
	return parsedUrl, path
}

func createRequest(request request) (string, int) {
	// Create a 10 second timeout
	timeout, _ := time.ParseDuration("10s")

	// Create new socket dialer
	dialer := net.Dialer{
		Timeout: timeout,
	}

	// Try to create new connection
	connection, err := tls.DialWithDialer(&dialer, "tcp", request.url+":https", nil)

	// If an error occured, print it to the user and exit
	if err != nil {
		printError(err.Error())
		os.Exit(0)
	}

	// Defer closing the connection
	defer connection.Close()

	// Create the request string
	requestString := []byte("GET " + request.path + " HTTP/1.0\r\nHost: " + request.url + "\r\n\r\n")

	// Write the string out to the connection
	connection.Write(requestString)

	// Read in response
	response, err := ioutil.ReadAll(connection)
	responseString := string(response)
	// Check if any errors occured
	if err != nil {
		printError(err.Error())
		os.Exit(0)
	}

	// Extract code from the response string
	code, _ := strconv.Atoi(responseString[9:12])

	// Return the body and code
	return responseString, code
}

func sendRequest(responses chan response, request request) {
	// Create start time
	start := time.Now()

	// Create request and get the results
	responseBody, code := createRequest(request)

	// If user isn't profiling just print the body
	if !request.profile {
		printBody(responseBody)
	}

	// Create response struct to pass through the channel
	responseStruct := response{
		code: code,
		time: int(time.Since(start).Milliseconds()),
		size: len(responseBody),
	}

	// Pass the channel through the channel
	responses <- responseStruct
}

func handleRequest(request request) {
	// Get sanitized url & path
	request.url, request.path = parseUrlPath(request.url)

	// create go channel for passing responses
	responses := make(chan response, request.iterations)

	// Create variables for calculating profile results
	var timeSum int
	times := make([]int, 0)
	sizes := make([]int, 0)
	codes := make([]int, 0)

	// create temp response to populate with channel
	response := response{}

	// Loop through requests as many times as specified
	for index := 0; index < request.iterations; index++ {
		// Call function with channel
		go sendRequest(responses, request)
		// Get results from channel
		response = <-responses

		// Populate variables
		timeSum += response.time
		times = append(times, response.time)
		sizes = append(sizes, response.size)
		codes = append(codes, response.code)
	}
	// Close channel
	close(responses)

	// Only check errors and construct profile
	// struct if we are profiling
	if request.profile {
		// Check if any errors
		badCodes := make([]int, 0)

		for _, code := range codes {
			if code != 200 {
				badCodes = append(badCodes, code)
			}
		}

		// Sort times & sizes arrays
		sort.Ints(times)
		sort.Ints(sizes)

		// Construct profile struct from data
		profile := profile{
			url:              request.url,
			iterations:       request.iterations,
			fastest:          times[0],
			slowest:          times[request.iterations-1],
			mean:             timeSum / request.iterations,
			median:           times[int(math.Floor(float64(len(times))/2.0))],
			errors:           badCodes,
			percentageError:  float64(len(badCodes) / request.iterations),
			largestResponse:  sizes[0],
			smallestResponse: sizes[len(sizes)-1],
		}

		// Finally print the profile data
		printProfile(profile)
	}
}

func main() {
	// Get arguments from cli
	arguments := os.Args[1:]

	// Empty request struct
	var result request

	// Validate arguments passed correctly (would love to do this differently)
	if len(arguments) > 0 {

		// Check if no arguments, or help flag passed
		if len(arguments) <= 0 || arguments[0] == "--help" {
			helpMessage()
		}

		// Check if we got 2 arguments, and the 1st is the url flag
		if len(arguments) == 2 && arguments[0] == "--url" {
			// Construct request struct
			result.iterations = 1
			result.profile = false
			result.url = arguments[1]
			handleRequest(result)
			return
		}

		// Check if we got 4 arguments, and that --url and --profile are correct
		if len(arguments) == 4 && arguments[0] == "--url" && arguments[2] == "--profile" {
			// Construct request struct
			result.profile = true
			result.url = arguments[1]

			// Ensure argument 4 can be parsed as an integer
			if iterations, err := strconv.Atoi(arguments[3]); err == nil {
				result.iterations = iterations
				handleRequest(result)
				return
			} else {
				invalidMessage()
			}
		}
	}
	// If we got here something went wrong
	invalidMessage()
}
