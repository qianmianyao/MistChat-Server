package dot

import "time"

type MessageType string

const (
	TextMessage   MessageType = "text"
	ImageMessage  MessageType = "image"
	VideoMessage  MessageType = "video"
	FileMessage   MessageType = "file"
	SystemMessage MessageType = "system"
)

type Source struct {
	Uid  string `json:"uid"`
	Name string `json:"name"`
}

type Attachment struct {
	URL    string `json:"url"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	Size   int64  `json:"size"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

type Content struct {
	Text       string      `json:"text,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Attachment *Attachment `json:"attachment,omitempty"`
}

type DataMessage struct {
	Type    MessageType `json:"type"`
	Content Content     `json:"content"`
}

type ReadStatus struct {
	ReadBy   []string `json:"readBy"`
	UnreadBy []string `json:"unreadBy"`
}

type Envelope struct {
	Source      Source      `json:"source"`
	Message     DataMessage `json:"message"`
	ReadStatus  *ReadStatus `json:"readStatus,omitempty"`
	Destination string      `json:"destination"`
	Timestamp   time.Time   `json:"timestamp"`
}
