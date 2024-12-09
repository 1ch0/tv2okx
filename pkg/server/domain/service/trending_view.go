package service

import (
	"context"
	"github.com/1ch0/tv2okx/pkg/server/infrastructure/datastore"
	apis "github.com/1ch0/tv2okx/pkg/server/interfaces/api/dto/v1"
)

// TrendingViewService is service for systemInfoCollection
type TrendingViewService interface {
	Init(ctx context.Context) error
	Webhook(ctx context.Context, req *apis.TrendingViewRequest) error
}

type trendingViewServiceImpl struct {
	Store  datastore.DataStore `inject:"datastore"`
	CMC    CMCService          `inject:""`
	Notify NotifyService       `inject:""`
}

// NewTrendingViewService return a systemInfoCollectionService
func NewTrendingViewService() TrendingViewService {
	return &trendingViewServiceImpl{}
}

func (tv trendingViewServiceImpl) Init(ctx context.Context) error {

	return nil
}

func (tv trendingViewServiceImpl) Webhook(ctx context.Context, req *apis.TrendingViewRequest) error {
	// 收到信号

	// 通知
	if err := tv.Notify.Do(ctx, req); err != nil {
		return err
	}

	// 下单
	if err := tv.CMC.Transfer(ctx, &apis.CMCTransferRequest{}); err != nil {
		return err
	}

	return nil
}
