package writer

import (
	"agent/internal/cluster"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

var (
	Accept  = 0
	Dropped = 0
)

func RunWriter(ct *cluster.Replicas, mh string) {
	time.Sleep(5 * time.Second)
	log.Info().Msg("Run as Writer")
	ct.MasterConn.Exec("CREATE TABLE IF NOT EXISTS test (id integer PRIMARY KEY);")
	FirstTest(ct)
	SecondTest(ct)
}
func FirstTest(ct *cluster.Replicas) {
	log.Info().Msg("Run first test")
	db := ct.MasterConn
	for i := 0; i < 10000; i++ {
		if !Write(db, i) {
			log.Error().Msg("Cannot write line to DB")
		}
		if i == 5000 {
			log.Info().Msg("Shutdown slave")

			http.Get("http://pg-slave:8080/shutdown")
		}
	}
	http.Get("http://pg-slave:8080/accept")

	ShowResults()
}
func SecondTest(ct *cluster.Replicas) {
	log.Info().Msg("Run second test")
	db := ct.MasterConn
	for i := 0; i < 1000000; i++ {
		if !Write(db, i) {
			log.Error().Msg("Cannot write line to DB")
		}
		if i == 500000 {
			log.Info().Msg("Shutdown master")
			db = ct.SlaveConn

			http.Get("http://pg-master:8080/shutdown")
			http.Get("http://pg-slave:8080/promote")

			time.Sleep(10 * time.Second)
		}
	}

	ShowResults()
}

func Write(db *sql.DB, number int) bool {
	log.Info().Int("number", number).Msg("Write num")
	psqlInfo := fmt.Sprintf("INSERT INTO public.test (id) VALUES (%d)", number)
	_, err := db.Exec(psqlInfo)
	if err != nil {
		Dropped += 1
		return false
	}
	Accept += 1
	return true
}

func ShowResults() {
	log.Info().Int("Accepted: ", Accept).Msg("Results")
	log.Info().Int("Dropped: ", Dropped).Msg("Results")
	Accept = 0
	Dropped = 0
}
