package http

//easyjson:json
type paylaod struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}
