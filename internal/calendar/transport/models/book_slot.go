package models

type BookSlotRequest struct {
	Urgent bool    `json:"urgent"`
	Title  *string `json:"title,omitempty"`
}
