package processor

import (
	"hash/fnv"
	"strconv"
	"strings"

	"iot-zero/services/tcp-server/internal/protocol/model"
)

// Return 32-bit hash result.
func fnv32(src string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(src))
	return h.Sum32()
}

func genAuthCode(d *model.Device) string {
	var splitByte byte = '_'

	codeBuilder := new(strings.Builder)
	codeBuilder.WriteString(d.ID)
	codeBuilder.WriteByte(splitByte)
	codeBuilder.WriteString(d.Plate)
	codeBuilder.WriteByte(splitByte)
	codeBuilder.WriteString(d.Phone)

	return strconv.Itoa(int(fnv32(codeBuilder.String())))
}
