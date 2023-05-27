package model

type Chat struct{
	ID string `json:"id"`
	From string `json:"from"`
	To string `json:"to"`
	Msg string `json:"msg"`
	Timestamp int64 `json:"timestamp"`
}

type ContactList struct{
	Username string `json:"username"`
	LastActivity string `json:"last_activity"`
}