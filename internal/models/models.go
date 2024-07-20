package models

// ShortenURLRequest описывает запрос пользователя.
type ShortenURLRequest struct {
	URL string `json:"url"`
}

// ShortenURLResponse описывает ответ сервера.
type ShortenURLResponse struct {
	Result string `json:"result"`
}

type URLRecord struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
