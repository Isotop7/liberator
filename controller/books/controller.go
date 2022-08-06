package books

// Controller for books
type Controller struct {
}

// Book struct
type Book struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Language  string `json:"language"`
	Category  string `json:"category"`
	ISBN10    string `json:"isbn10" binding:"len=10"`
	ISBN13    string `json:"isbn13" binding:"len=13"`
	PageCount int    `json:"page_count"`
	Rating    int    `json:"rating"`
}

func (c *Controller) Index() ([]*Book, error) {
	return []*Book{}, nil
}
