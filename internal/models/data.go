package models

type Data struct {
	Login    string `json:"login"`
	Password string `json:"password,omitempty"`
	IP       string `json:"ip"`
}
