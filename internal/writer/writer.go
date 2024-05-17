package writer

import (
	"agent/internal/cluster"
	"context"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

var (
	Accept  = 0
	Dropped = 0
)

func RunWriter(ct *cluster.Replicas) {
	time.Sleep(5 * time.Second)
	log.Info().Msg("Run as Writer")
	FirstTest(ct)
	SecondTest(ct)
}
func FirstTest(ct *cluster.Replicas) {
	log.Info().Msg("Run first test")
	tset := true
	_, err := ct.MasterConn.Exec("CREATE TABLE IF NOT EXISTS test (id integer);")
	if err != nil {
		log.Err(err).Msg("Cannot create table")
	}
	for i := 0; i < 1000000; i++ {
		if tset {
			result := WriteMasterLine(ct, i)
			if !result {
				log.Info().Msg("Cannot write to master")
				tset = false
			}
		} else {
			log.Info().Msg("Write to slave")
			WriteSlaveLine(ct, i)
		}
		if i == 500000 {
			log.Info().Msg("Shutdown master")
			http.Get("http://pg-slave:8080/shutdown")
		}
	}
	http.Get("http://pg-slave:8080/accept")
	WriteResults()
}
func SecondTest(ct *cluster.Replicas) {
	log.Info().Msg("Run second test")
	tset := true
	_, err := ct.MasterConn.Exec("CREATE TABLE IF NOT EXISTS test (id integer);")
	if err != nil {
		log.Err(err).Msg("Cannot create table")
	}
	for i := 0; i < 1000000; i++ {
		if tset {
			result := WriteMasterLine(ct, i)
			if !result {
				log.Info().Msg("Cannot write to master")
				tset = false
			}
		} else {
			log.Info().Msg("Write to slave")
			WriteSlaveLine(ct, i)
		}
		if i == 500000 {
			log.Info().Msg("Shutdown master")
			http.Get("http://pg-master:8080/shutdown")
		}
	}
	WriteResults()
}
func WriteMasterLine(ct *cluster.Replicas, number int) bool {
	ctx := context.Background()
	tx, err := ct.MasterConn.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return false
	}
	_, err = tx.ExecContext(ctx, "INSERT INTO test (id) VALUES ($1)", number)
	if err != nil {
		Dropped += 1
		return false
	}
	Accept += 1
	return true
}
func WriteSlaveLine(ct *cluster.Replicas, number int) bool {
	ctx := context.Background()
	tx, err := ct.SlaveConn.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return false
	}
	_, err = tx.ExecContext(ctx, "INSERT INTO test (id) VALUES ($1)", number)
	if err != nil {
		Dropped += 1
		return false
	}
	Accept += 1
	return true
}
func WriteResults() {
	log.Info().Int("Accepted: ", Accept).Msg("Results")
	log.Info().Int("Dropped: ", Dropped).Msg("Results")
	Accept = 0
	Dropped = 0
}
