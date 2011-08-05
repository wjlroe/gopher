package main

import (
	"xml"
	"http"
	"log"
	"fmt"
)

type Site struct {
	Name string
	Apikey string
	RoomName string
	Rooms []Room
	ChosenRoom *Room
}

func (site *Site) CampfireUrl(path string) (url string) {
	url = "https://" + site.Apikey + ":X@" + site.Name + ".campfirenow.com" + path
	return url
}

func (site *Site) GetRooms() {
	url := site.CampfireUrl("/rooms.xml")

	parsed_url := ParseURL(url)

	//fmt.Printf("Going to request URL: %s\n", parsedUrl.String())
	response, err := http.Get(parsed_url.String())
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Printf("Status: %s\n", response.Status)
		log.Fatal("Could not list rooms")
	}

	parser := xml.NewParser(response.Body)
	rooms := Rooms{Room: nil}
	err = parser.Unmarshal(&rooms, nil)
	if err != nil {
		log.Fatal("Error unmarshalling xml:", err)
	}
	site.Rooms = rooms.Room
	fmt.Println("Rooms",rooms)
}

func (site *Site) Whoami() (user *User) {
	url := site.CampfireUrl("/users/me.xml")
	parsed_url := ParseURL(url)

	response, err := http.Get(parsed_url.String())
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Printf("Status: %s\n", response.Status)
		log.Fatal("Could not get me")
	}

	parser := xml.NewParser(response.Body)
	user = &User{}
	err = parser.Unmarshal(&user, nil)
	if err != nil {
		log.Fatal("Error unmarshalling xml:", err)
	}
	return user
}

func (site *Site) JoinRoom() {
	for i, _ := range site.Rooms {
		if site.Rooms[i].Name == site.RoomName {
			site.ChosenRoom = &site.Rooms[i]
			break
		}
	}
	if site.ChosenRoom == nil {
		log.Fatal("Could not find chosen room: %s in available rooms: %s", site.RoomName, site.Rooms)
	}
	me := site.Whoami()
	site.ChosenRoom.Site = site
	site.ChosenRoom.Join(me)
}