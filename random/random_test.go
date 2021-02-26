package random_test

import (
	"testing"
	"trace/random"
)

func TestRandomBase(t *testing.T) {
	uint64_id := random.RandomUint64()
	t.Log(uint64_id)

	uuid, err := random.RandomUUID()
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(uuid)
}
