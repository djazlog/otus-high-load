package model

// DialogMessage представляет сообщение в диалоге для Redis репозитория
type DialogMessage struct {
	ID          string `redis:"id"`
	FromUserID  string `redis:"from_user_id"`
	ToUserID    string `redis:"to_user_id"`
	Text        string `redis:"text"`
	CreatedAtNs int64  `redis:"created_at"`
	DialogKey   string `redis:"dialog_key"`
}
