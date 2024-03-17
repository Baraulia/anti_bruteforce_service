package models

type Data struct {
	Login    string `json:"login"`
	Password string `json:"password,omitempty"`
	Ip       string `json:"ip"`
}
