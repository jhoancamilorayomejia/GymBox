package models

type PricePlan struct {
	IDPricePlan int     `json:"idpriceplan"`
	TypePlan    string  `json:"typeplan"`
	Price       float64 `json:"price"`
}
