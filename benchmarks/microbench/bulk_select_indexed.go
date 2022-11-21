package main

import (
	"fmt"

	"github.com/Shopify/mybench"
)

type BulkSelectIndexed struct {
	mybench.WorkloadConfig
	table *mybench.Table
}

func NewBulkSelectIndexed(config *mybench.BenchmarkConfig, table *mybench.Table, eventRate float64) mybench.AbstractWorkload {
	eventRate = eventRate * config.Multiplier
	var workloadInterface mybench.WorkloadInterface[MicroBenchContextData] = &BulkSelectIndexed{
		WorkloadConfig: mybench.NewWorkloadConfigWithDefaults(mybench.WorkloadConfig{
			Name:           "BulkSelectIndexed",
			DatabaseConfig: config.DatabaseConfig,
			RateControl: mybench.RateControlConfig{
				EventRate: eventRate,
			},
		}),
		table: table,
	}

	workload, err := mybench.NewWorkload(workloadInterface)
	if err != nil {
		panic(err)
	}

	return workload
}

func (c *BulkSelectIndexed) Config() mybench.WorkloadConfig {
	return c.WorkloadConfig
}

func (c *BulkSelectIndexed) Event(ctx mybench.WorkerContext[MicroBenchContextData]) error {
	args := []interface{}{
		c.table.Generate(ctx.Rand, "idx2"),
	}

	_, err := ctx.Data.Statement.Execute(args...)
	return err
}

func (c *BulkSelectIndexed) NewContextData(conn *mybench.Connection) (MicroBenchContextData, error) {
	var err error
	contextData := MicroBenchContextData{}

	query := fmt.Sprintf("SELECT * FROM %s WHERE idx2 = ?", c.table.Name)
	contextData.Statement, err = conn.Prepare(query)
	return contextData, err
}
