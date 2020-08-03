package structs

// Payload is a document object
type Payload struct {
	ContentHash string  `json:"content_hash,omitempty"` // A base64 representation form of the document's content.
	ID          *string `json:"id,omitempty"`           // The document ID.
	Content     *string `json:"content,omitempty"`      // The document content.
	Extension   *string `json:"extension,omitempty"`    // The extension of the document.
	CreatedAt   *int    `json:"created_at,omitempty"`   // The Unix timestamp of when the document was inserted.
	UpdatedAt   *int    `json:"updated_at,omitempty"`   // The Unix timestmap of when the document was last modified.
	Exists      *bool   `json:"exists,omitempty"`       // Whether the document does or does not exist.
}

// Response is a Spacebin API response
type Response struct {
	Error   string  `json:"error"` // .Error() should already be called
	Payload Payload `json:"payload"`
	Status  int     `json:"status"`
}
