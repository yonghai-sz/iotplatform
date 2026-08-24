package logic

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"iot-zero/services/platform-api/internal/svc"
	"iot-zero/services/platform-api/internal/types"

	"iot-zero/services/rpc-transform/transformer"
)

func TestShortenLogic_Shorten(t *testing.T) {
	ast := assert.New(t)

	// Build mock and svc context
	ctl := gomock.NewController(t)
	mockTransformer := transformer.NewMockTransformer(ctl)

	l := NewShortenLogic(context.Background(),
		&svc.ServiceContext{
			Transformer: mockTransformer,
		})

	// Failed to simulate transformer Shorten
	mockTransformer.EXPECT().Shorten(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("call err")).
		Times(1)
	_, err := l.Shorten(&types.ShortenReq{})
	ast.NotNil(err)

	// Simulate transformer Shorten success
	mockTransformer.EXPECT().Shorten(gomock.Any(), gomock.Any()).
		Return(&transformer.ShortenResp{Shorten: "123"}, nil).
		Times(1)
	resp, err := l.Shorten(&types.ShortenReq{})
	ast.Nil(err)
	ast.Equal(resp, &types.ShortenResp{Shorten: "123"})
}
