package model

type Lottery struct {
	ID               int64  `json:"id"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	Prize            string `json:"prize"`
	Status           string `json:"status"`
	AvailableTickets int64  `json:"available_tickets"`
}
