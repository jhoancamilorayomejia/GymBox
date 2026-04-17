package models

type Customer struct {
	IDCustomer int    `json:"idcustomer"`
	Cedula     string `json:"cedula"`
	Name       string `json:"name"`
	LastName   string `json:"lastname"`
	Phone      string `json:"phone"`
	Password   string `json:"password"`
}
