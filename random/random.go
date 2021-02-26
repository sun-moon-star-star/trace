package random

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix() ^ int64(os.Getpid()))
}

// for uuid sequenceId
func RandomUint12() uint16 {
	return uint16(rand.Int31n(4096))
}

func RandomUint64() uint64 {
	return uint64(rand.Uint32())<<32 + uint64(rand.Uint32())
}

// 格式: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
func RandomUUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid, nil
}
