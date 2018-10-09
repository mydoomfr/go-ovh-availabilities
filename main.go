package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

// TODO : Unitest

var (
	configFile string
	config     = Configuration{}
)

const (
	API_AVAILABLE          = "available"
	API_UNAVAILABLE        = "unavailable"
	API_DEFAULT_DATACENTER = "default"
	DATETIME_FORMAT        = "2006/01/02 15:04:05"
	MAIL_TEMPLATE_HTML     = "template_mail.html"
	CONFIG_FILE            = "config.json"
)

// Availabilities Hello, world!
type Availabilities []struct {
	Hardware    string `json:"hardware"`
	Region      string `json:"region"`
	Datacenters []struct {
		Datacenter       string `json:"datacenter"`
		Availability     string `json:"availability"`
		LastAvailability string
	}
}

func (a *Availabilities) getAvailabilities() {
	myClient := &http.Client{Timeout: 10 * time.Second}
	u, err := url.Parse(config.ApiUrl)
	if err != nil {
		log.Fatal(err)
	}

	response, err := myClient.Get(u.String())
	if err != nil {
		log.Println(err.Error())
		return
	}

	if response.StatusCode != 200 {
		b, _ := ioutil.ReadAll(response.Body)
		log.Println(string(b))
		return
	}

	err = json.NewDecoder(response.Body).Decode(&a)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *Availabilities) checkAvailabilities(c *Configuration) {
	for _, a := range *a {
		for _, c := range c.Wanted {

			if !strings.EqualFold(a.Hardware, c.Hardware) {
				continue
			}

			if len(c.Region) > 0 && !strings.EqualFold(a.Region, c.Region) {
				continue
			}

			for i := 0; i < len(a.Datacenters); i++ {
				d := &a.Datacenters[i]

				// Prevent twice notification
				if strings.EqualFold(d.Datacenter, API_DEFAULT_DATACENTER) {
					continue
				}

				if len(c.Datacenters) > 0 && !Contains(c.Datacenters, d.Datacenter) {
					continue
				}

				// If 1st iteration : Check if server is already available
				// Else print the diff between now/old value
				if d.LastAvailability == "" {
					if !strings.EqualFold(d.Availability, API_UNAVAILABLE) {
						notify(a.Hardware, a.Region, d.Availability, d.Datacenter, c.Name)
					}
				} else {
					if !strings.EqualFold(d.Availability, d.LastAvailability) {
						notify(a.Hardware, a.Region, d.Availability, d.Datacenter, c.Name)
					}
				}

				// Save Availability for next loop
				d.LastAvailability = d.Availability
			}
		}
	}
}

func notify(hardware, region, availability, datacenter, name string) {
	var now = time.Now().Format(DATETIME_FORMAT)
	var status = API_UNAVAILABLE
	if availability != API_UNAVAILABLE {
		status = API_AVAILABLE
	}

	fmt.Printf("%s %s (%s) is now : %s (%s) in %s (%s) !\n",
		now, name, hardware, status, availability, region, datacenter)

	// Send Mail
	r := NewRequest([]string{config.Mail.To}, config.Mail.Object)
	r.Send(MAIL_TEMPLATE_HTML,
		map[string]string{
			"datetime": now, "name": name, "hardware": hardware, "region": region, "status": status,
			"availability": availability, "datacenter": datacenter})
}

func (a *Availabilities) start() {
	a.getAvailabilities()
	if a.isEmpty() {
		log.Printf("Data not found, retry in %v seconds...\n", config.Refresh)
	} else {
		go a.checkAvailabilities(&config)
	}
	time.Sleep(time.Duration(config.Refresh) * time.Second)
}

func (a Availabilities) isEmpty() bool {
	return reflect.DeepEqual(a, Availabilities{})
}

// Contains tells whether a contains x.
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func main() {
	// Args
	flag.StringVar(&configFile, "c", CONFIG_FILE, "config file")
	flag.Parse()

	log.Print("Configuration file : ", configFile)

	// Load Config
	config.loadConfiguration(configFile)
	ovh := &Availabilities{}
	for {
		ovh.start()
	}
}
