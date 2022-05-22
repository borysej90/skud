package app

type checkAccessReq struct {
	ReaderID int64  `json:"reader_id"`
	PassCard string `json:"pass_card"`
}

type checkAccessResp struct {
	Access  bool   `json:"access"`
	Message string `json:"message,omitempty"`
}
