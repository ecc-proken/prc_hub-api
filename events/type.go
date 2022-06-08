package events

type Event struct {
	Id                  uint64          `json:"id"`
	UserId              uint64          `json:"user_id,omitempty"`
	Title               string          `json:"title"`
	Description         *string         `json:"description,omitempty"`
	Location            string          `json:"location,omitempty"`
	Published           bool            `json:"published"`
	Completed           bool            `json:"completed"`
	AutoNotifyDocuments bool            `json:"auto_notify_documents_enabled"`
	Documents           []EventDocument `json:"documents,omitempty"`
}

type EventDocument struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}
