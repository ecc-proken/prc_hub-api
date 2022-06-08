package events

import "prc_hub-api/mysql"

type PostBody struct {
	Title               string  `json:"title" validate:"required,gte=1"`
	Description         *string `json:"description" validate:"omitempty"`
	Location            string  `json:"location" validate:"omitempty,gte=1"`
	Published           bool    `json:"published"`
	Completed           bool    `json:"completed"`
	AutoNotifyDocuments bool    `json:"auto_notify_documents_enabled"`
}

func Post(userId uint64, post PostBody) (e Event, err error) {
	// 書込
	result, err := mysql.Write(
		`INSERT INTO events (user_id, title, description, location, published, completed, auto_notify_documents_enabled)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
		userId, post.Title, post.Description, post.Location, post.Published, post.Completed, post.AutoNotifyDocuments,
	)
	if err != nil {
		return
	}

	// Insertした行のIdを取得
	id, err := result.LastInsertId()
	if err != nil {
		return
	}

	e.Id = uint64(id)
	e.UserId = userId
	e.Title = post.Title
	e.Description = post.Description
	e.Location = post.Location
	e.Published = post.Published
	e.Completed = post.Completed
	e.AutoNotifyDocuments = post.AutoNotifyDocuments
	return
}
