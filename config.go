package main

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
)

// Configuration
type Configuration struct {
	ApiUrl  string `json:"api"`
	Refresh int    `json:"refresh"`
	Wanted  []struct {
		Name        string   `json:"name"`
		Hardware    string   `json:"hardware"`
		Region      string   `json:"region"`
		Datacenters []string `json:"datacenters"`
	} `json:"wanted"`
	Mail struct {
		From     string `json:"from"`
		To       string `json:"to"`
		Object   string `json:"object"`
		SMTP     struct {
			Active   bool   `json:"active"`
			Server   string `json:"server"`
			Port     uint16 `json:"port"`
			Username string `json:"username"`
			Password string `json:"password"`
		} `json:"smtp"`
		Sendmail struct {
			Active bool   `json:"active"`
			Bin    string `json:"bin"`
		} `json:"sendmail"`
	} `json:"mail"`
}

func (c *Configuration) loadConfiguration(configFile string) {
	f, err := os.Open(configFile)
	defer f.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
	jsonParser := json.NewDecoder(f)
	jsonParser.Decode(&c)
	if c.isEmpty() {
		log.Fatalf("The configuration file \"%s\" is invalid or empty", configFile)
	}
}

func (c Configuration) isEmpty() bool {
	return reflect.DeepEqual(c, Configuration{})
}
