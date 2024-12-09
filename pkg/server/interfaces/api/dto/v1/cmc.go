package v1

type CMCTransferRequest struct {
	State string `json:"state" validate:"required" description:"状态"`
	Coin  string `json:"coin" validate:"required" description:"货币类型"`
}
