package service

import (
	"context"
	"fmt"
	"github.com/nikoksr/notify/service/dingding"

	"github.com/1ch0/tv2okx/pkg/server/config"
	"github.com/1ch0/tv2okx/pkg/server/infrastructure/datastore"
	apis "github.com/1ch0/tv2okx/pkg/server/interfaces/api/dto/v1"
	"github.com/1ch0/tv2okx/pkg/server/utils"
	"github.com/1ch0/tv2okx/pkg/server/utils/log"
)

// NotifyService is service for systemInfoCollection
type NotifyService interface {
	Init(ctx context.Context) error
	Do(ctx context.Context, req *apis.TrendingViewRequest) error
}

type notifyServiceImpl struct {
	Config config.Config       `json:"config"`
	Store  datastore.DataStore `inject:"datastore"`
	dd     *dingding.Service
}

// NewNotifyService return a systemInfoCollectionService
func NewNotifyService() NotifyService {
	return &notifyServiceImpl{}
}

func (n notifyServiceImpl) Init(ctx context.Context) error {
	//cfg := &dingding.Config{
	//	Token:  "dddd",
	//	Secret: "xxx",
	//}
	//n.dd = dingding.New(cfg)

	return nil
}

func (n notifyServiceImpl) Do(ctx context.Context, req *apis.TrendingViewRequest) error {
	if err := utils.Validate.Struct(req); err != nil {
		return err
	}
	//msg := fmt.Sprintf("标记种类："+coin+"，标记方向："+state+"，触发K线时段："+indicatorType+"，当前价格："+price)
	msg := fmt.Sprintf("标记种类[%s] 标记方向[%s] 触发K线时段[%s] 当前价格[%s]",
		req.Coin, req.State, req.IndicatorType, req.Price)

	// TODO: 发送通知
	log.Logger.Infof(msg)
	//err := n.dd.Send(context.Background(), "subject", "content")
	//if err != nil {
	//	return err
	//}
	return nil
}
