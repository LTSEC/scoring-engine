package scoring

import (
	"testing"
	"time"

	"github.com/LTSEC/scoring-engine/config"
)

func TestScoring(t *testing.T) {
	go func() {
		t.Log(ScoringStartup(config.Parse("../tests/test_yaml.yaml")))
	}()
	time.Sleep(10 * time.Second)
}
