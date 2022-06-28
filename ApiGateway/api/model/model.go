package model

type User struct {
	Id          string  `json:"id"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Login       string  `json:"login"`
	Password    string  `json:"password"`
	Email       string  `json:"email"`
	Bio         string  `json:"bio"`
	PhoneNumber string  `json:"phone_number"`
	TypeId      int64   `json:"type_id"`
	Status      bool    `json:"status"`
	Address     Address `json:"address"`
	Posts       []*Post
}
type Address struct {
	Id         string `json:"id"`
	UserId     string `json:"user_id"`
	Country    string `json:"country"`
	City       string `json:"city"`
	District   string `json:"district"`
	PostalCode int64  `json:"postal_code"`
}
type Post struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	UserId      string   `json:"user_id"`
	Medias      []*Media `json:"medias"`
}
type Media struct {
	Id     string `json:"id"`
	PostId string `json:"post_id"`
	Type   string `json:"type"`
	Link   string `json:"link"`
}
type Id struct {
	Id string `json:"id"`
}
type Check struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type Code struct {
	Codd string `json:"code"`
}
type Login struct {
	EmailAddress string `json:"email_address"`
	Password     string `json:"password"`
}
type LoginRes struct {
	UserData *User
	Token    string
	Refresh  string
}
type JwtReqMod struct {
	Token string `json:"token"`
}
