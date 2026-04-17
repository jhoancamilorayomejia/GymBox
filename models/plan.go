package models

type Plan struct {
	IDPlan     int     `json:"idplan"`
	IDCustomer int     `json:"idcustomer"`
	TypePlan   string  `json:"typeplan"`
	Price      float64 `json:"price"`
	DatePay    string  `json:"datepay"`
	DateStart  string  `json:"datestart"`
	DateFinish string  `json:"datefinish"`
	State      string  `json:"state"`
}
