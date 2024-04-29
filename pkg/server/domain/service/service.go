package service

import (
	"context"
	"fmt"
)

// needInitData register the service that need to init data
var needInitData []DataInit

// InitServiceBean init all service instance
func InitServiceBean() []interface{} {

	trendingView := NewTrendingViewService()
	notifyService := NewNotifyService()
	cmcService := NewCMCService()

	needInitData = []DataInit{
		trendingView,
		notifyService,
		cmcService,
	}

	return []interface{}{
		trendingView,
		notifyService,
		cmcService,
	}
}

// DataInit the service set that needs init data
type DataInit interface {
	Init(ctx context.Context) error
}

// InitData init data
func InitData(ctx context.Context) error {
	for _, init := range needInitData {
		if err := init.Init(ctx); err != nil {
			return fmt.Errorf("database init failure %w", err)
		}
	}
	return nil
}
