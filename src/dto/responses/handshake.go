package responses

type Handshake struct {
	Code string `json:"code"`
}

func (h *Handshake) SetHandshake(code string) {
	h.Code = code
}
