package domain

import (
	"time"
)

type (
	FileType string
)

const (
	Image FileType = "image"
	Video FileType = "video"
	Other FileType = "other"
)

// размер, дата загрузки, айди пользователя, ссылка на внешнее хранилище
type File struct {
	ID          int64     `json:"-"           db:"id"`
	UserID      int64     `json:"-"           db:"user_id"`
	Name        string    `json:"name"        db:"name"         bson:"name"`
	Size        int64     `json:"size"        db:"size"         bson:"size"`
	Type        FileType  `json:"type"        db:"type"         bson:"type"`
	ContentType string    `json:"contentType" db:"content_type" bson:"contentType"`
	URL         string    `json:"url"         db:"url"          bson:"url,omitempty"`
	UploadAt    time.Time `json:"uploadAt"    db:"upload_at"    bson:"uploadAt"`
}
