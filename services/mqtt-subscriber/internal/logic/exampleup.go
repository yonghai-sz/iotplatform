package logic

import (
	"context"
	"encoding/json"
	"runtime/debug"
	"unicode/utf8"

	"github.com/zeromicro/go-zero/core/logx"
)

const maxPayloadLogRunes = 512

func truncateForLog(b []byte, maxRunes int) string {
	s := string(b)
	if maxRunes <= 0 || utf8.RuneCountInString(s) <= maxRunes {
		return s
	}
	runes := []rune(s)
	return string(runes[:maxRunes]) + "…(truncated)"
}

func HandleExampleUp(ctx context.Context, topic string, payload []byte) (err error) {
	logger := logx.WithContext(ctx)
	defer func() {
		if p := recover(); p != nil {
			logger.Errorf("HandleExampleUp panic: %v\n%s", p, string(debug.Stack()))
		}
	}()

	var doc map[string]any
	if err := json.Unmarshal(payload, &doc); err != nil {
		logger.Errorf("invalid JSON topic=%s len=%d err=%v", topic, len(payload), err)
		return err
	}

	logger.Infof("topic=%s field_count=%d", topic, len(doc))
	logger.Debugf("topic=%s payload=%s", topic, truncateForLog(payload, maxPayloadLogRunes))
	return nil
}
