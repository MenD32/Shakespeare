package generators

import (
	"github.com/MenD32/Shakespeare/pkg/trace"
)

type Generator interface {
	Generate() (trace.TraceLog, error)
}