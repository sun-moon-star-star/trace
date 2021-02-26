package trace

import (
	"errors"
	"testing"
)

func TestSpan(t *testing.T) {
	spaner := NewSpaner()
	spaner.Tags["id"] = 123456789
	spaner.Tags["name"] = "zhaolu"

	spaner.Logs["error"] = errors.New("unknown error")
	spaner.Logs["action"] = "success"

	spaner.Baggages["data-access"] = "(0, ok)"

	spaner.End()
	t.Log(spaner.FormatMapStrategy(spaner.Tags))
	t.Log(spaner.FormatMapStrategy(spaner.Logs))
	t.Log(spaner.FormatMapStrategy(spaner.Baggages))

	t.Log(spaner.FormatSpanerStrategy(spaner))
}
