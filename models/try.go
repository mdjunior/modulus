package models

// Try is struct for try a attempt on game
type Try struct {
	Name  string `json:"name" form:"name" query:"name"`
	Value int    `json:"value" form:"value" query:"value"`
}
