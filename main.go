package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func main() {
	sourceUrl := "https://updates.jenkins.io/update-center.actual.json"

	// CLI Arguments
	// var buildDateFlag = flag.Bool("b", false, "Print the build date")
	// var delimeterFlag = flag.String("d", ",", "The delimeter used to separate returned data")
	// var nameFlag = flag.Bool("n", false, "Print the plugin name")
	// var pluginsFlag = flag.String("p", "", "A comma-separated list of plugin names to return information for. Use * to list all plugins")
	// var sizeFlag = flag.Bool("s", false, "Print the size in bytes")
	// var sha1Flag = flag.Bool("sha1", false, "Print the sha1")
	// var sha256Flag = flag.Bool("sha256", false, "Print the sha256")
	// var urlFlag = flag.Bool("u", false, "Print the download url")
	// var versionFlag = flag.Bool("v", false, "Print the version")
	flag.Parse()

	defer handleErrors()

	data := getUpdateData(sourceUrl)

	fmt.Println(data["core"].(map[string]interface{})["url"])
}

func handleErrors() {
	if message := recover(); message != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", message)
		os.Exit(1)
	}
}

func getUpdateData(url string) map[string]interface{} {
	response, error := http.Get(url)
	if error != nil {
		panic(error)
	}

	if response.StatusCode != 200 {
		panic(fmt.Sprintf("updates.jenkins.io responded with %s", strconv.Itoa(response.StatusCode)))
	}

	responseBody, readError := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if readError != nil {
		panic(readError)
	}

	var data map[string]interface{}
	parseError := json.Unmarshal(responseBody, &data)
	if parseError != nil {
		panic(parseError)
	}

	return data
}
