package models

type RawData struct {
	Command int
	Payload []byte
}

type HashData struct {
	ZeroCount int    `json:"zero_count"`
	Data      string `json:"data"`
	Nonce     int    `json:"nonce"`
}

type AccessToken struct {
	Token string `json:"token"`
}
