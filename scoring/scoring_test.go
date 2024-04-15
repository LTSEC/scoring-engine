package scoring

import (
	"testing"
	"time"

	"github.com/LTSEC/scoring-engine/config"
)

func TestScoring(t *testing.T) {
	go func() {
		t.Log(ScoringStartup(config.Parse("E:/Projects/scoring-engine/tests/example.yaml")))
	}()
	time.Sleep(1 * time.Second)
	t.Log(ScoringToggle(false))
}
