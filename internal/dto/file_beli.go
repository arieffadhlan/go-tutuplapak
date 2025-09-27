package dto

type UploadBeliData struct {
	ImageUrl string `json:"imageUrl"`
}

type UploadBeliResponse struct {
	Message string         `json:"message"`
	Data    UploadBeliData `json:"data"`
}
