package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"gopkg.in/yaml.v2"
)

type Site struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type Config struct {
	Sites []Site `yaml:"sites"`
}

func main() {
	// Load configuration from YAML file
	config, err := loadConfig("sites.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}

	// Check status of each site
	statuses := make(map[string]string)
	for _, site := range config.Sites {
		if isSiteUp(site.URL) {
			statuses[site.Name] = ":green_square:"
		} else {
			statuses[site.Name] = ":red_square:"
		}
	}

	// Write result to markdown file
	err = writeMarkdown("README.md", "README.md.tmpl", statuses)
	if err != nil {
		log.Fatalf("Failed to write output file: %s", err)
	}

	// Print result to console if show_output is true
	show_output := true // Change this to false if you don't want to print to console
	if show_output {
		for site, status := range statuses {
			fmt.Printf("%s: %s\n", site, status)
		}
	}
}

func loadConfig(filename string) (*Config, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func isSiteUp(url string) bool {
	client := http.Client{
		Timeout: time.Second * 5,
	}

	resp, err := client.Get(url)
	if err != nil {
		return false
	}

	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func writeMarkdown(filename string, templateFilename string, statuses map[string]string) error {
	templateData, err := ioutil.ReadFile(templateFilename)
	if err != nil {
		return err
	}

	tmpl, err := template.New("readme").Parse(string(templateData))
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	err = tmpl.Execute(file, statuses)
	if err != nil {
		return err
	}

	return nil
}
