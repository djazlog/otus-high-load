package model

type Post struct {
	ID           string `redis:"id"`
	Text         string `redis:"text"`
	AuthorUserId string `redis:"author_user_id"`
	CreatedAtNs  int64  `redis:"created_at"`
	UpdatedAtNs  *int64 `redis:"updated_at"`
}
