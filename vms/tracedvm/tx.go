// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package tracedvm

import (
	"context"

	"go.opentelemetry.io/otel/attribute"

	"github.com/skychains/chain/snow/consensus/snowstorm"
	"github.com/skychains/chain/trace"

	oteltrace "go.opentelemetry.io/otel/trace"
)

var _ snowstorm.Tx = (*tracedTx)(nil)

type tracedTx struct {
	snowstorm.Tx

	tracer trace.Tracer
}

func (t *tracedTx) Verify(ctx context.Context) error {
	ctx, span := t.tracer.Start(ctx, "tracedTx.Verify", oteltrace.WithAttributes(
		attribute.Stringer("txID", t.ID()),
	))
	defer span.End()

	return t.Tx.Verify(ctx)
}

func (t *tracedTx) Accept(ctx context.Context) error {
	ctx, span := t.tracer.Start(ctx, "tracedTx.Accept", oteltrace.WithAttributes(
		attribute.Stringer("txID", t.ID()),
	))
	defer span.End()

	return t.Tx.Accept(ctx)
}

func (t *tracedTx) Reject(ctx context.Context) error {
	ctx, span := t.tracer.Start(ctx, "tracedTx.Reject", oteltrace.WithAttributes(
		attribute.Stringer("txID", t.ID()),
	))
	defer span.End()

	return t.Tx.Reject(ctx)
}
