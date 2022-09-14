package model

type Disease struct {
	ID                int    `json:"id"`
	Name		      string `json:"name"`
	Spec_id           int    `json:"spec_id"`
}