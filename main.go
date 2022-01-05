package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

func main() {
	sourceUrl := "https://updates.jenkins.io/update-center.actual.json"

	defer handleErrors()

	// CLI Arguments
	PrintDelimeterFlag := flag.String("d", "", "The delimeter used to separate returned data. Prints table output by default")
	headerFlag := flag.Bool("h", true, "Print the header row")
	// PluginsFlag := flag.String("p", "", "A comma-separated list of plugin names to return information for. Use * to list all plugins")

	type column struct {
		displayName string
		enabled     *bool
		keyName     string
	}

	columns := []column{
		{"NAME", flag.Bool("n", false, "Print the plugin name"), "name"},
		{"VERSION", flag.Bool("v", false, "Print the version"), "version"},
		{"SIZE", flag.Bool("s", false, "Print the size in bytes"), "size"},
		{"SHA1", flag.Bool("sha1", false, "Print the sha1"), "sha1"},
		{"SHA256", flag.Bool("sha256", false, "Print the sha256"), "sha256"},
		{"BUILD DATE", flag.Bool("b", false, "Print the build date"), "buildDate"},
		{"URL", flag.Bool("u", true, "Print the download url"), "url"},
	}
	flag.Parse()

	// We only care about columns that are enabled so we filter them early and save some iterations and conditionals later
	columnsFiltered := []column{}
	for _, col := range columns {
		if *col.enabled {
			columnsFiltered = append(columnsFiltered, col)
		}
	}
	if len(columnsFiltered) == 0 {
		panic("The command requires at least one piece of data should be returned")
	}
	if len(columnsFiltered) == 1 {
		*headerFlag = false
	}

	data := getUpdateData(sourceUrl)

	// setup the stdout writer using the tabwriter package
	var delimeter string
	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 4, ' ', 0)
	delimeter = "\t"
	defer writer.Flush()
	if *PrintDelimeterFlag != "" {
		// We still need a tabwriter object but this one will not add any spacing. Just a simple delimeter
		writer = tabwriter.NewWriter(os.Stdout, 0, 0, 0, []byte(*PrintDelimeterFlag)[0], tabwriter.AlignRight)
		delimeter = *PrintDelimeterFlag
	}

	// Header Row
	var header string
	if *headerFlag {
		for _, col := range columnsFiltered {
			header += col.displayName + delimeter
		}
		fmt.Fprintln(writer, strings.TrimSuffix(header, delimeter))
	}

	// Print data for the main jenkins.war
	var buffer string
	core := data["core"].(map[string]interface{})
	for _, col := range columnsFiltered {
		buffer += fmt.Sprint(core[col.keyName]) + delimeter
	}
	fmt.Fprintln(writer, strings.TrimSuffix(buffer, delimeter))
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
