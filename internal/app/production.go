package app

import (
	"context"
	"time"
	"unicorn/factory"
)

type productionLine struct {
	factory  factory.Factory
	logistic *logisticsCenter
}

// NewProductionLine creates a new production line.
func NewProductionLine(factory factory.Factory, logistics *logisticsCenter) *productionLine {
	return &productionLine{
		factory:  factory,
		logistic: logistics,
	}
}

// StartProduction starts producing unicorns at rate.
func (pl *productionLine) StartProduction(ctx context.Context, rate time.Duration) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(rate):
			pl.logistic.HandleUnicorn(pl.factory.NewUnicorn())
		}
	}
}
