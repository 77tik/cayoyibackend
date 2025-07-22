package chain

import (
	"context"
	"testing"
)

func TestChain(t *testing.T) {
	branch := NewBranch(WithBranchHandlers(
		PrintHandler(),
		RangePrintHandler(10),
	))

	driver, _ := NewDriver(WithDefaultBranch(branch))
	ctx := context.Background()
	driver.Chain(ctx)
}
