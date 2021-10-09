package models

//Create Struct
type User struct {
	Id       string `json:"Id"`
	Name     string `json:"Name"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
}
type Post struct {
	Id       ID     `json:"Id"`
	Caption  string `json:"Caption"`
	ImageURL string `json:"ImageURL"`
	PostedTS string `json:"PostedTS"`
}
type ID struct {
	UserID string `json:"UserId"`
	PostID string `json:"PostId"`
}
