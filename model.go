package main

import "encoding/json"

type ReceivedMessage struct {
	Object string   `json:"object"`
	Entry  []*Entry `json:"entry"`
}

type Entry struct {
	ID        json.Number  `json:"id"`
	Time      json.Number  `json:"time"`
	Messaging []*Messaging `json:"messaging"`
}

type Messaging struct {
	Sender    *Sender     `json:"sender"`
	Recipient *Recipient  `json:"recipient"`
	Timestamp json.Number `json:"timestamp"`
	Message   *Message    `json:"message"`
	Delivery  *Delivery   `json:"delivery"`
	PostBack  *PostBack   `json:"postback"`
}

type PostBack struct {
	Payload string `json:"payload"`
}

type Sender struct {
	ID json.Number `json:"id"`
}

type Recipient struct {
	ID json.Number `json:"id"`
}

type Message struct {
	Attachments []*Attachment `json:"attachments"`
	MID         string        `json:"mid"`
	Seq         json.Number   `json:"seq"`
	Text        *string       `json:"text"`
	StickerID   *json.Number  `json:"sticker_id"`
}

type Delivery struct {
	Mids      []string    `json:"mids"`
	Watermark json.Number `json:"watermark"`
	Seq       json.Number `json:"seq"`
}

type ID struct {
	ID int `json:"id"`
}

type Text struct {
	Text string `json:"text"`
}

type SendMessage struct {
	Recipient *ID   `json:"recipient"`
	Message   *Text `json:"message"`
}

type SendImage struct {
	Recipient *ID                `json:"recipient"`
	Message   *AttachmentMessage `json:"message"`
}

type AttachmentMessage struct {
	Attachment *Attachment `json:"attachment"`
}

type Attachment struct {
	Payload *Payload `json:"payload"`
	Type    string   `json:"type"`
}

type Payload struct {
	URL string `json:"url"`
}

type FBUserProfile struct {
	ID        *json.Number `json:"fid"`
	FirstName string       `json:"first_name"`
	//Gender     string       `json:"gender"`
	LastName   string `json:"last_name"`
	ProfilePic string `json:"profile_pic"`
}

type UserExist struct {
	Result bool `json:"result"`
}

type RegistResponse struct {
	Result string `json:"result"`
}

type DefaultResponse struct {
	Result json.Number `json:"result"`
}

type MatchResponse struct {
	Fid        json.Number `json:"fid"`
	ID         int         `json:"id"`
	Name       string      `json:"name"`
	ProfileImg string      `json:"profile_img"`
}
