package main

import (
	"xml"
	"log"
	"fmt"
	"strings"
	"http"
	"io/ioutil"
	"bytes"
	"bufio"
	"json"
)

type Rooms struct {
	Room []Room
}

func (rooms *Rooms) String() string {
	ret := make([]string, len(rooms.Room))
	for i, r := range rooms.Room {
		ret[i] = r.String()
	}
	return strings.Join(ret, ", ")
}

type Room struct {
	Site *Site
	Id int
	Name string
	Topic string
	Full bool
	Users []*User
	Me *User
}

func (room *Room) String() string {
	users := "("
	for i,_ := range(room.Users) {
		users += room.Users[i].String() + ", "
	}
	users += ")"
	return room.Name + " Members: (" + users + ")"
}

func (room *Room) StreamingUrl() (url string) {
	url = fmt.Sprintf("http://%s:X@streaming.campfirenow.com/room/%d/live.json", room.Site.Apikey, room.Id)
	return url
}

func (room *Room) Show() {
	room_path := fmt.Sprintf("/room/%d.xml", room.Id)
	url := room.Site.CampfireUrl(room_path)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal("Could not get room")
	}
	defer response.Body.Close()

	parser := xml.NewParser(response.Body)
	err = parser.Unmarshal(&room, nil)
	if err != nil {
		log.Fatal("Error unmarshalling room:", err)
	}
	var body []byte
	body, err = ioutil.ReadAll(response.Body)
	fmt.Printf("Body: %s\n", string(body))
}

func (room *Room) Join(me *User) {
	room.Show()
	fmt.Printf("Room: %s\n", room)
	path := fmt.Sprintf("/room/%d/join.xml", room.Id)
	url := room.Site.CampfireUrl(path)

	parsed_url := ParseURL(url)

	response, err := Post(parsed_url.String(), "", "")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Printf("Status: %s\n", response.Status)
		log.Printf("RequestMethod: %s\n", response.Request.Method)
		var body []byte
		body, err = ioutil.ReadAll(response.Body)
		fmt.Printf("Body: %s\n", string(body))
		log.Fatal("Could not join room")
	}
	room.Me = me
	room.Say("Beep Beep. I have nothing to say yo.")

	// catch all signals...
	go CatchSignals(room)

	room.Stream()
}

func (room *Room) Leave() {
	path := fmt.Sprintf("/room/%d/leave.xml", room.Id)
	url := room.Site.CampfireUrl(path)

	parsed_url := ParseURL(url)

	response, err := Post(parsed_url.String(), "", "")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Printf("Status: %s\n", response.Status)
		log.Printf("RequestMethod: %s\n", response.Request.Method)
		var body []byte
		body, err = ioutil.ReadAll(response.Body)
		fmt.Printf("Body: %s\n", string(body))
		log.Fatal("Could not join room")
	}
}

func (room *Room) Say(say_this string) {
	path := fmt.Sprintf("/room/%d/speak.xml", room.Id)
	url := room.Site.CampfireUrl(path)

	request_xml := `<message>
<body>`+say_this+`</body>
</message>`

	fmt.Printf("Request XML: \n%s\n", request_xml)

	response, err := Post(url, request_xml, "application/xml")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != 201 {
		log.Printf("Status: %s\n", response.Status)
		log.Fatal("Could not say anything!")
	}
}

func (room *Room) FindUser(user_id int) (user *User) {
	for _,user := range(room.Users) {
		if user.Id == user_id {
			return user
		}
	}
	return nil
}

func (room *Room) Stream() {
	url := room.StreamingUrl()
	fmt.Printf("Going to stream from: %s\n", url)
	parsedUrl := ParseURL(url)

	response, err := http.Get(parsedUrl.String())
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Printf("Status: %s\n", response.Status)
		log.Fatal("Could not stream room")
	}

	reader := bufio.NewReader(response.Body)
	for {
		line, err := reader.ReadBytes('}')
		if err != nil {
			log.Fatal(err)
		}
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		fmt.Printf("Received: %s\n\n", line)
		message := &Message{}
		err = json.Unmarshal(line, &message)
		fmt.Printf("Message: %s\n", message.String())

		if strings.Contains(strings.ToLower(message.Body), "gopher") &&
			message.User_Id != room.Me.Id {
			user := room.FindUser(message.User_Id)
			user_name := ""
			if user != nil {
				user_name = user.Name
			}

			command_return := Command(message.Body)
			if len(command_return) > 0 {
				command_return[0] = user_name + ": " + command_return[0]
				//say_back := user_name + ": " + command_return[0]
				//room.Say(say_back)
				for i,_ := range command_return {
					room.Say(command_return[i])
				}
			}
		}
	}
}
