package models

type CommonHeaders struct {
	UserAgent      string `header:"User-Agent" binding:"required"`
	AcceptLanguage string `header:"Accept-Language" binding:"required"`
}
