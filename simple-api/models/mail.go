package models

type Content struct {
	Subject string `json:"subject" bson:"subject"`
	Content string `json:"content" bson:"content"`
}

type Mail struct {
	ID       int     `json:"id" bson:"id"`
	Sender   string  `json:"sender" bson:"sender"`
	Receiver string  `json:"receiver" bson:"receiver"`
	Body     Content `json:"body" bson:"body"`
}
