package demo

import (
	"context"
	"math/rand"

	"github.com/kiru997/go-ex/pkg/tracingutil"
)

type Service interface {
	Demo(ctx context.Context) int
	Rand(ctx context.Context) int
}

type demoService struct {
}

func NewService() Service {
	return &demoService{}
}

func (*demoService) Demo(ctx context.Context) int {
	_, span := tracingutil.Start(ctx, "demoService.Demo")
	defer span.End()
	return rand.Int()
}

func (d *demoService) Rand(ctx context.Context) int {
	_, span := tracingutil.Start(ctx, "demoService.Rand")
	defer span.End()
	return d.Demo(ctx)
}
