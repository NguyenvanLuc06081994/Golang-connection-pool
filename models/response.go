package models

type Response struct {
	Elapsed  int64      `json:"elapsed"`
	Average  float64    `json:"average"`
	Products []*Product `json:"products"`
}
