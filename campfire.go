package main

import (
	"fmt"
	"http"
	"log"
	"bytes"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"willroe-goconf.googlecode.com/hg"
)

// const SERVER = "ardourcreations.campfirenow.com"

// // ardour:
// const APIKEY = "84a8e2a71d7947c45e8ec92e2a13bdbe009e4cdc"

// // madred:
// //const APIKEY = "5decd9115ca6aebd3fa796a58c51fb509a6a9cdb"
// const ROOM = "Test"
var caught_signals = []string{"SIGINT","SIGTERM"}

func ParseURL(url string) (parsed_url *http.URL) {
	parsed_url, err := http.ParseURL(url)
	if err != nil {
		log.Fatal(err)
	}
	return parsed_url
}

func Post(url string, body string, content_type string) (response *http.Response, err os.Error) {
	parsed_url := ParseURL(url)

	postBody := bytes.NewBufferString(body)
	fmt.Printf("Going to post to: %s\n", url)
	request, err := http.NewRequest("POST", parsed_url.String(), postBody)
	request.Header = http.Header{
  		"Content-Type": {content_type},
		"Content-Length": {strconv.Itoa(postBody.Len())},
  	}
	request.ContentLength = int64(postBody.Len())
	response, err = http.DefaultClient.Do(request)
	return response, err
}

func CatchSignals(room *Room) {
	for {
		signal := <-signal.Incoming

		signal_parts := strings.Split(signal.String(), ":", 2)
		signal_name := signal_parts[0]
		log.Printf("Signal name: %s\n", signal_name)
		for _,name := range(caught_signals) {
			if name == signal_name {
				log.Printf("%s - Leave all rooms!\n", name)
				room.Leave()
				log.Fatal("Received interupt")
			}
		}
		log.Printf("Signal: %s\n", signal.String())
	}
}

func getConfig(config_file string, site_name string) (site *Site) {
	config, err := conf.ReadConfigFile(config_file)
	if err != nil {
		log.Fatal("Error reading config file: ", err)
	}
	name, err := config.GetString(site_name, "name")
	apikey, err := config.GetString(site_name, "apikey")
	room, err := config.GetString(site_name, "room")
	if name == "" {
		log.Fatal("Name is blank for site: ", site_name)
	}
	if apikey == "" {
		log.Fatal("Apikey is blank for site: ", site_name)
	}
	if room == "" {
		log.Fatal("Room is blank for site: ", site_name)
	}
	site = &Site{Name: name, Apikey: apikey, RoomName: room}
	return site
}

func main() {
	if len(os.Args) > 1 {
		cmd := os.Args[1]
		if cmd == "command" {
			if len(os.Args) > 2 {
				output := Command(strings.Join(os.Args[2:]," "))
				fmt.Printf("Command output: %s\n", output)
				os.Exit(0)
			} else {
				log.Fatal("You need to give me a command to run!")
			}
		}
	}

	if len(os.Args) < 3 {
                log.Fatal("Usage: campfire config_file site_name")
        }
        config_file := os.Args[1]
        site_name := os.Args[2]

	site := getConfig(config_file, site_name)
	site.GetRooms()
	site.JoinRoom()
}
