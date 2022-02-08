package models

/*
TradeIn describes a car to add to the database
*/
type TradeIn struct {
	Make  string `json:"make"`
	Model string `json:"model"`
}
