package models

// ShortenURLRequest описывает запрос пользователя.
type ShortenURLRequest struct {
	URL string `json:"url"`
}

// ShortenURLResponse описывает ответ сервера.
type ShortenURLResponse struct {
	Result string `json:"result"`
}
