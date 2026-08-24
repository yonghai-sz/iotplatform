package utils_test

import (
	"iot-zero/pkg/utils"
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		salt     string
		want     string
	}{
		{
			name:     "test1",
			password: "123456",
			salt:     "abc456",
			want:     "b7d6c2ee39faf9179ef1ae83d2dbb6842da3915a8b3370f376c2675c953fc67fa27e2d411f620be21fc1f80b338aeb7bd1ea58162b7ceb6c1aac991b6a5ef65a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.HashPassword(tt.password, tt.salt)
			t.Logf("got: %v", got)
			if got != tt.want {
				t.Errorf("HashPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
