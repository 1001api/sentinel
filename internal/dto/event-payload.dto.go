package dto

type KeyPayload struct {
	PublicKey string `json:"PublicKey"`
}

type PrivateKeyPayload struct {
	PrivateKey string `json:"PrivateKey"`
}

type UserIDPayload struct {
	UserID string `json:"UserID"`
}
