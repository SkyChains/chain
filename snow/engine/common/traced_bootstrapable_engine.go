// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package common

import (
	"context"

	"github.com/skychains/chain/trace"
)

var _ BootstrapableEngine = (*tracedBootstrapableEngine)(nil)

type tracedBootstrapableEngine struct {
	Engine
	bootstrapableEngine BootstrapableEngine
	tracer              trace.Tracer
}

func TraceBootstrapableEngine(bootstrapableEngine BootstrapableEngine, tracer trace.Tracer) BootstrapableEngine {
	return &tracedBootstrapableEngine{
		Engine:              TraceEngine(bootstrapableEngine, tracer),
		bootstrapableEngine: bootstrapableEngine,
	}
}

func (e *tracedBootstrapableEngine) Clear(ctx context.Context) error {
	ctx, span := e.tracer.Start(ctx, "tracedBootstrapableEngine.Clear")
	defer span.End()

	return e.bootstrapableEngine.Clear(ctx)
}
