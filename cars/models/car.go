package models

/*
Car describes a car to retrieve from the database
*/
type Car struct {
	TradeIn
	Id int64 `json:"id"`
}
