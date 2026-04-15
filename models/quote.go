package models

type Quote struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	QuoteLine string `json:"quoteline"`
}

