package expenses

type Expenses struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Amount string   `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

type Err struct {
	Message string `json:"message"`
}
