package main

type Message struct {
	Room_Id int
	Created_At string
	Body string
	Id int
	User_Id int
	Type string
}

func (msg *Message) String() string {
	return msg.Created_At + ": " + msg.Body
}

//{"room_id":1,"created_at":"2009-12-01 23:44:40","body":"hello","id":1, "user_id":1,"type":"TextMessage"}