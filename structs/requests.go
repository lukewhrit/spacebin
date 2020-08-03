package structs

// CreateRequest is the structure a request to POST / should follow
type CreateRequest struct {
	Content   string `json:"content" xml:"content" form:"content"`
	Extension string `json:"extension" xml:"extension" form:"extension"`
}
