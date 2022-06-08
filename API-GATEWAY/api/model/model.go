package model

type Useri struct {
	Id          string   `json:"id"`
	FirstName   string   `json:"first_name"`
	LastName    string   `json:"last_name"`
	Email       []string `json:"email"`
	Bio         string   `json:"bio"`
	PhoneNumber []string `json:"phone_number"`
	TypeId      int64    `json:"type_id"`
	Status      bool     `json:"status"`
	Address     Address  `json:"address"`
}
type Address struct {
	Id         string `json:"id"`
	UserId     string `json:"user_id"`
	Country    string `json:"country"`
	City       string `json:"city"`
	District   string `json:"district"`
	PostalCode int64  `json:"postal_code"`
}
