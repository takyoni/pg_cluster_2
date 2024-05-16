package writer

import (
	"agent/internal/cluster"
	"context"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func RunWriter(ct *cluster.Replicas) {
	time.Sleep(5 * time.Second)
	log.Info().Msg("Run as Writer")
	tset := true
	ct.MasterConn.Exec("CREATE TABLE IF NOT EXISTS test (id integer);")
	for i := 0; i < 1000000; i++ {
		if tset {
			result := WriteMasterLine(ct, i)
			log.Info().Msg("Cannot write to master")
			if !result {
				tset = false
			}
		} else {
			WriteSlaveLine(ct, i)
		}
		if i == 500000 {
			http.Get("http://pg-master:8080/shutdown")
		}
	}
}
func WriteMasterLine(ct *cluster.Replicas, number int) bool {
	ctx := context.Background()
	tx, err := ct.MasterConn.BeginTx(ctx, nil)
	if err != nil {
		return false
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, "INSERT INTO test (id) VALUES (?)", number)
	return err == nil
}
func WriteSlaveLine(ct *cluster.Replicas, number int) bool {
	ctx := context.Background()
	tx, err := ct.SlaveConn.BeginTx(ctx, nil)
	if err != nil {
		return false
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, "INSERT INTO test (id) VALUES (?)", number)
	return err == nil
}
