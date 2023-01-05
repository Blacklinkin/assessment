package expenses

type Handler struct {
	Database database
}

func (h *Handler) InitialDB() {
	h.Database.InitDatabase()
}

func (h *Handler) CloseDB() {
	h.Database.CloseDatabase()
}
