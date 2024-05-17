package slave

import (
	"agent/internal/cluster"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func RunSlave(ct *cluster.Replicas) {
	log.Info().Msg("Run as Slave")
	go RunServer(ct)
	for {
		time.Sleep(1 * time.Second)
		arbiter, err := ct.CheckAM()
		if err == nil && !arbiter && !ct.CheckMaster() {
			log.Info().Msg("Promote to Master")

			cmd := exec.Command("touch", "/tmp/touch_me_to_promote_to_me_master")
			err := cmd.Run()

			if err == nil {
				log.Info().Msg("Success promote to Master")
				break
			}

			log.Info().Err(err).Msg("Error promote to Master")
		}
	}
}
func RunServer(ct *cluster.Replicas) {
	server := gin.Default()
	server.GET("/shurdown", Shutdown)
	server.GET("/accept", Accept)
	server.Run(":8080")
}
func Shutdown(c *gin.Context) {
	err := exec.Command("iptables", "-P", "INPUT", "-p", "tcp", "--dport", "5432", "-j", "DROP").Run()
	//cmd := exec.Command("iptables", "-A", "INPUT", "DROP")
	if err != nil {
		log.Err(err).Msg("Cannot block input d5432 connections to Master")
	}
	err = exec.Command("iptables", "-P", "INPUT", "-p", "tcp", "--sport", "5432", "-j", "DROP").Run()
	//cmd := exec.Command("iptables", "-A", "INPUT", "DROP")
	if err != nil {
		log.Err(err).Msg("Cannot block input s5432 connections to Master")
	}

}
func Accept(c *gin.Context) {
	cmd := exec.Command("iptables", "-F")
	//cmd := exec.Command("iptables", "-A", "INPUT", "DROP")
	err := cmd.Run()

	if err == nil {
		log.Info().Msg("Success block connections to Master")
	}
}
