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

type SiteStatus struct {
	Name   string
	URL    string
	Status string
}

func main() {
	// Load configuration from YAML file
	config, err := loadConfig(getEnv("SITES_YAML_PATH", "./sites.yaml"))
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}

	// Check status of each site
	statuses := make([]SiteStatus, 0, len(config.Sites))
	for _, site := range config.Sites {
		status := SiteStatus{
			Name:   site.Name,
			URL:    site.URL,
			Status: ":red_square:",
		}
		if isSiteUp(site.URL) {
			status.Status = ":green_square:"
		}
		statuses = append(statuses, status)
	}

	// Write result to markdown file
	templatePath := getEnv("README_TEMPLATE_PATH", "./README.md.tmpl")
	err = writeMarkdown("README.md", templatePath, statuses)
	if err != nil {
		log.Fatalf("Failed to write output file: %s", err)
	}

	// Print result to console if show_output is true
	show_output := true // Change this to false if you don't want to print to console
	if show_output {
		for _, status := range statuses {
			fmt.Printf("%s: %s\n", status.Name, status.Status)
		}
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
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

func writeMarkdown(filename string, templateFilename string, statuses []SiteStatus) error {
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
