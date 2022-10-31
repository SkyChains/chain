// Copyright (C) 2019-2022, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package snowstorm

import (
	"context"

	"github.com/luxdefi/luxd/ids"
	"github.com/luxdefi/luxd/snow/events"
	"github.com/luxdefi/luxd/utils/wrappers"
)

var _ events.Blockable = (*acceptor)(nil)

type acceptor struct {
	g        *Directed
	errs     *wrappers.Errs
	deps     ids.Set
	rejected bool
	txID     ids.ID
}

func (a *acceptor) Dependencies() ids.Set { return a.deps }

func (a *acceptor) Fulfill(ctx context.Context, id ids.ID) {
	a.deps.Remove(id)
	a.Update(ctx)
}

func (a *acceptor) Abandon(_ context.Context, id ids.ID) { a.rejected = true }

func (a *acceptor) Update(context.Context) {
	// If I was rejected or I am still waiting on dependencies to finish or an
	// error has occurred, I shouldn't do anything.
	if a.rejected || a.deps.Len() != 0 || a.errs.Errored() {
		return
	}
	a.errs.Add(a.g.accept(a.txID))
}
