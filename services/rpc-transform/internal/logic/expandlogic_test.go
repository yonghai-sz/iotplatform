package logic

import (
	"context"
	"errors"
	"testing"

	"iotplatform/services/rpc-transform/internal/svc"
	"iotplatform/services/rpc-transform/pb/transform"

	"iotplatform/services/rpc-transform/model"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestExpandLogic_Expand(t *testing.T) {
	ast := assert.New(t)

	// Build mock models and svc context
	ctl := gomock.NewController(t)
	shortModel := model.NewMockExampleModel(ctl)
	svcCtx := &svc.ServiceContext{
		Model: shortModel,
	}
	// build expand logic
	logic := NewExpandLogic(context.Background(), svcCtx)

	// Failed to simulate model FindOne
	shortModel.EXPECT().FindOne(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("call find one error")).
		Times(1)
	_, err := logic.Expand(&transform.ExpandReq{})
	ast.NotNil(err)

	// Simulate model FindOne success
	shortModel.EXPECT().FindOne(gomock.Any(), gomock.Any()).
		Return(
			&model.Example{
				Shorten: "testShorten",
				Url:     "testUrl",
			},
			nil,
		).
		Times(1)
	resp, err := logic.Expand(&transform.ExpandReq{})
	ast.Nil(err)
	ast.Equal(resp.Url, "testUrl")
}
