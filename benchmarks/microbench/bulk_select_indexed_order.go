package main

import (
	"fmt"

	"github.com/Shopify/mybench"
)

type BulkSelectIndexedOrder struct {
	mybench.WorkloadConfig
	table      *mybench.Table
	orderField string
}

func NewBulkSelectIndexedOrder(config MicroBenchConfig, table *mybench.Table, eventRate float64, orderField string) mybench.AbstractWorkload {
	eventRate = eventRate * config.Multiplier
	var workloadInterface mybench.WorkloadInterface[MicroBenchContextData] = &BulkSelectIndexedOrder{
		WorkloadConfig: mybench.NewWorkloadConfigWithDefaults(mybench.WorkloadConfig{
			Name:           "BulkSelectIndexedOrdered_" + orderField,
			DatabaseConfig: config.DatabaseConfig,
			RateControl: mybench.RateControlConfig{
				EventRate: eventRate,
			},
		}),
		table:      table,
		orderField: orderField,
	}

	workload, err := mybench.NewWorkload(workloadInterface)
	if err != nil {
		panic(err)
	}

	return workload
}

func (c *BulkSelectIndexedOrder) Config() mybench.WorkloadConfig {
	return c.WorkloadConfig
}

func (c *BulkSelectIndexedOrder) Event(ctx mybench.WorkerContext[MicroBenchContextData]) error {
	args := []interface{}{
		c.table.Generate(ctx.Rand, "idx2"),
	}

	_, err := ctx.Data.Statement.Execute(args...)
	return err
}

func (c *BulkSelectIndexedOrder) NewContextData(conn *mybench.Connection) (MicroBenchContextData, error) {
	var err error
	contextData := MicroBenchContextData{}

	query := fmt.Sprintf("SELECT * FROM %s WHERE idx2 = ? ORDER BY %s", c.table.Name, c.orderField)
	contextData.Statement, err = conn.Prepare(query)
	return contextData, err
}