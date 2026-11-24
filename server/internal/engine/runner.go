package engine

import (
	"context"
	"time"

	"github.com/voidsgnl/bastion/server/internal/store"
)

type Runner interface {
	Run(context.Context, store.Scan, store.Target) ([]store.Finding, error)
}

type FixtureRunner struct{}

func (FixtureRunner) Run(ctx context.Context, scan store.Scan, target store.Target) ([]store.Finding, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return []store.Finding{{
		ID:           "f_" + scan.ID,
		ScanID:       scan.ID,
		AssetID:      "asset_" + target.ID,
		SourceEngine: "fixture",
		Title:        "DEMO: Example weak SSH configuration finding",
		Category:     "network",
		Severity:     "high",
		Confidence:   "unverified",
		Status:       "open",
		Evidence: []store.Evidence{{
			Type:    "other",
			Summary: "Synthetic fixture only; no target was contacted.",
		}},
		Risk: store.Risk{
			Score:       82,
			Explanation: "Demonstration score only; no target assessment was performed.",
		},
		FirstSeen: now,
		LastSeen:  now,
	}}, nil
}
