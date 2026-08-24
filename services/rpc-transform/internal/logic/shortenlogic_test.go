package logic

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"iot-zero/services/rpc-transform/internal/svc"
	"iot-zero/services/rpc-transform/pb/transform"

	"iot-zero/services/rpc-transform/model"
)

func TestShortenLogic_Shorten(t *testing.T) {
	ast := assert.New(t)

	// Build mock models and svc context
	ctl := gomock.NewController(t)
	shortModel := model.NewMockExampleModel(ctl)
	svcCtx := &svc.ServiceContext{
		Model: shortModel,
	}
	// build platform logic
	logic := NewShortenLogic(context.Background(), svcCtx)

	// Failed to simulate model insert
	shortModel.EXPECT().Insert(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("insert error")).
		Times(1)
	_, err := logic.Shorten(&transform.ShortenReq{Url: "testUrl"})
	ast.NotNil(err)

	// Simulate model insert success
	shortModel.EXPECT().Insert(gomock.Any(), gomock.Any()).
		Return(nil, nil).
		Times(1)
	resp, err := logic.Shorten(&transform.ShortenReq{Url: "testUrl"})
	ast.Nil(err)
	ast.True(resp.Shorten != "")
}
