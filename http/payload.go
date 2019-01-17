package http

//easyjson:json
type paylaod struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}
