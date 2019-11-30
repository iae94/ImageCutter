package models

type Image struct {
	Name string
	MimeType string
	Url string
	Size int64
	Headers map[string]string
	FetchCount int
}