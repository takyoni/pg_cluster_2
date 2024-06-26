package arbiter

import (
	"agent/internal/cluster"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type Server struct {
	ct *cluster.Replicas
}

func RunArbiter(ct *cluster.Replicas) {
	log.Info().Msg("Run as Arbiter")
	handler := &Server{ct: ct}

	server := gin.New()
	server.Use(gin.Recovery())
	server.GET("/master", handler.MasterStatus)
	server.GET("/ping", handler.Ping)

	server.Run(":8080")
}

func (s *Server) MasterStatus(c *gin.Context) {
	result := s.ct.CheckMaster()
	if !result {
		c.JSON(http.StatusBadGateway, gin.H{"master": result})
	} else {
		c.JSON(http.StatusOK, gin.H{"master": result})
	}
}

func (s *Server) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}
