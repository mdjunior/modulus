package models

// Game is struct for joinned clients on game
type Game struct {
	ID       string           `json:"id" form:"id" query:"id"`
	Name     string           `json:"name" form:"name" query:"name"`
	Status   string           `json:"status" form:"status" query:"status"`
	Users    []map[string]int `json:"users" form:"users" query:"users"`
	MaxUsers int              `json:"max_users" form:"max_users" query:"max_users"`
	Winner   map[string]int   `json:"winner,omitempty" form:"winner" query:"winner"`
}
