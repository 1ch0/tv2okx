package service

import (
	"context"

	"github.com/1ch0/tv2okx/pkg/server/config"
	"github.com/1ch0/tv2okx/pkg/server/infrastructure/datastore"
	apis "github.com/1ch0/tv2okx/pkg/server/interfaces/api/dto/v1"
)

// CMCService is service for systemInfoCollection
type CMCService interface {
	Init(ctx context.Context) error
	Transfer(ctx context.Context, req *apis.CMCTransferRequest) error
}

type CMCServiceImpl struct {
	Config config.Config       `json:"config"`
	Store  datastore.DataStore `inject:"datastore"`
	//okxSwap *goexv2swap.PrvApi
}

// NewCMCService return a systemInfoCollectionService
func NewCMCService() CMCService {
	return &CMCServiceImpl{}
}

func (cmc CMCServiceImpl) Init(ctx context.Context) error {
	//okxCfg := cmc.Config.OKX
	//okxSwap := goexv2.OKx.Swap.NewPrvApi(
	//	options.WithApiKey(okxCfg.APIKey),
	//	options.WithApiSecretKey(okxCfg.APISecret),
	//	options.WithPassphrase(okxCfg.PassPhrase))
	//cmc.okxSwap = okxSwap

	return nil
}

func (cmc CMCServiceImpl) Transfer(ctx context.Context, req *apis.CMCTransferRequest) error {

	return nil
}

func (cmc CMCServiceImpl) createOrder(ctx context.Context, req *apis.CMCTransferRequest) error {
	return nil
}
