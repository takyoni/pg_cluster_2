package slave

import (
	"agent/internal/cluster"
	"os/exec"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func RunSlave(ct *cluster.Replicas) {
	log.Info().Msg("Run as Slave")
	server := gin.Default()
	server.GET("/shutdown", Shutdown)
	server.GET("/accept", Accept)
	server.GET("/promote", Promote)
	server.Run(":8080")
}

func Shutdown(c *gin.Context) {
	err := exec.Command("iptables", "-A", "INPUT", "-p", "tcp", "--dport", "5432", "-j", "DROP").Run()
	if err != nil {
		log.Err(err).Msg("Cannot block input d5432 connections to Master")
	}
	err = exec.Command("iptables", "-A", "INPUT", "-p", "tcp", "--sport", "5432", "-j", "DROP").Run()
	if err != nil {
		log.Err(err).Msg("Cannot block input s5432 connections to Master")
	}

}

func Accept(c *gin.Context) {
	cmd := exec.Command("iptables", "-F")
	err := cmd.Run()

	if err == nil {
		log.Info().Msg("Success block connections to Master")
	}
}

func Promote(c *gin.Context) {
	cmd := exec.Command("touch", "/tmp/touch_me_to_promote_to_me_master")
	err := cmd.Run()

	if err == nil {
		log.Info().Msg("Success promote to Master")
	}
}
