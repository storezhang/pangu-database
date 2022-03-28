package database

import (
	`github.com/pangum/pangu`
)

func init() {
	pangu.New().Dependencies(
		newEngine,
		newTx,
	)
}
