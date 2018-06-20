package datastore

import "github.com/sniperkit/watchub/shared/model"

type Execstore interface {
	Executions() ([]model.Execution, error)
}
