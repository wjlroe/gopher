package main


type User struct {
	Id int
	Name string
	Email_Address string
	Admin bool
	Created_At string
	Type string
	Avatar_Url string
}

func (user *User) String() string {
	if user != nil {
		return user.Name
	}
	return ""
}