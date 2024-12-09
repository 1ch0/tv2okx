package v1

type TrendingViewRequest struct {
	State         string `json:"state" validate:"required" description:"状态"`
	Coin          string `json:"coin" validate:"required" description:"货币类型"`
	IndicatorType string `json:"indicator_type" validate:"required" description:"指标类型"`
	Price         string `json:"price" validate:"required" description:"价格"`
}
