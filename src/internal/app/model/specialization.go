package model

type Specialization struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Salary int    `json:"salary"`
}

type Float struct {
	Num float64
}

type SpecStatView struct {
	Specialization *Specialization
	Disease        *Disease
	Percent        *Float
}