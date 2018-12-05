package http

//easyjson:json
type paylaod struct {
	To      string `json:"To"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}
