package table

type StartSessionRequest struct {
	QRToken string `json:"qr_token" validate:"required,uuid"`
}

type StartSessionResponse struct {
	SessionToken string `json:"session_token"`
	TableID      string `json:"table_id"`
	TableNo      string `json:"table_no"`
	ExpiresAt    string `json:"expires_at"`
}

type CheckoutCompletedEvent struct {
	EventID   string `json:"event_id"`
	Event     string `json:"event"`
	SessionID string `json:"session_id"`
	TableID   string `json:"table_id"`
}
