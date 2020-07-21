package main

import (
	"os"
	"fmt"
	"log"
	"io/ioutil"
	"encoding/json"
	"net/http"

	flags "github.com/jessevdk/go-flags"
)

var options struct {
	InputFile string `long:"in" default:"/dev/stdin"`
	OutputFile string `long:"out" default:"/dev/stdout"`

	NamespacesUrl string `long:"namespaces-url" required:"true" env:"NAMESPACES_URL"`
	ItemsUrl string `long:"items-url" required:"true" env:"ITEMS_URL"`

	EnrichNamespace bool `long:"enrich-namespace" env:"ENRICH_NAMESPACE"`
	EnrichNamespaceSettings bool `long:"enrich-namespace-settings" env:"ENRICH_NAMESPACE_SETTINGS"`
	EnrichMetadata bool `long:"enrich-metadata" env:"ENRICH_METADATA"`
	EnrichSettings bool `long:"enrich-settings" env:"ENRICH_SETTINGS"`
}

func GetItems(url string) ([]map[string]interface{}, error) {
	var items []map[string]interface{}

	resp, err := http.Get(url)
	if err != nil {
		return items, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal(body, &items)
	if err != nil {
		return items, err
	}

	return items, nil
}

func GetItem(url string) (map[string]interface{}, error) {
	var item map[string]interface{}

	resp, err := http.Get(url)
	if err != nil {
		return item, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return item, err
	}
	err = json.Unmarshal(body, &item)
	if err != nil {
		return item, err
	}

	return item, nil
}

func main() {
	var parser = flags.NewParser(&options, flags.Default)
	if _, err := parser.Parse(); err != nil {
		os.Exit(0)
	}

	inputFile, err := os.Open(options.InputFile)
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}
	outputFile, err := os.Open(options.InputFile)
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	jsonData, err := ioutil.ReadAll(inputFile)
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	var items []map[string]interface{}
	json.Unmarshal(jsonData, &items)

	for _,item := range items {
		if (options.EnrichNamespace) {
			namespace, err := GetItem(fmt.Sprintf("%s/%s", options.NamespacesUrl, item["namespace"]))
			if err == nil && options.EnrichNamespaceSettings {
				namespace_settings, err := GetItems(fmt.Sprintf("%s/%s/settings", options.NamespacesUrl, item["namespace"]))
				if err == nil {
					namespace["settings"] = namespace_settings
				}
			}
			item["namespace"] = namespace
		}

		if (options.EnrichMetadata) {
			metadata, err := GetItems(fmt.Sprintf("%s/%s/metadata", options.ItemsUrl, item["id"]))
			if err == nil {
				item["metadata"] = metadata
			}
		}
		if (options.EnrichSettings) {
			settings, err := GetItems(fmt.Sprintf("%s/%s/settings", options.ItemsUrl, item["id"]))
			if err == nil {
				item["settings"] = settings
			}
		}
	}

	b, err := json.Marshal(items)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(b))
	fmt.Fprintln(outputFile, string(b))
}
