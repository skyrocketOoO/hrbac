package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
	"github.com/skyrocketOoO/erx/erx"
	"github.com/skyrocketOoO/serverx/api"
	"github.com/skyrocketOoO/serverx/internal/boot"
	"github.com/skyrocketOoO/serverx/internal/controller"
	"github.com/skyrocketOoO/serverx/internal/controller/middleware"
	"github.com/skyrocketOoO/serverx/internal/global"
	"github.com/skyrocketOoO/serverx/internal/service/exter/db"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "server",
	Short: "The main service command",
	Long:  ``,
	// Args:  cobra.MinimumNArgs(1),
	Run: RunServer,
}

func RunServer(cmd *cobra.Command, args []string) {
	if err := boot.InitAll(); err != nil {
		log.Fatal().Msgf("Initialization failed: %v", err)
	}

	restController := controller.NewHandler()

	router := gin.Default()
	router.Use(middleware.Cors())
	api.Bind(router, restController)

	port, _ := cmd.Flags().GetString("port")
	// router.Run(":" + port)

	// Create a new server instance
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Info().Msgf("Starting server on port %s...", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("Server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error().Msgf("Server forced to shutdown: %v", err)
	} else {
		log.Info().Msg("Server shut down gracefully.")
	}

	if err := db.Close(); err != nil {
		log.Error().Msgf("Error closing database connection: %v", err)
	} else {
		log.Info().Msg("Database connection closed successfully")
	}
}

func init() {
	Cmd.Flags().StringP("port", "p", "8080", "port")
	Cmd.Flags().
		StringVarP(&global.Database, `database`, "d", "postgres", `"postgres", "mysql"`)
	Cmd.Flags().
		StringVarP(&global.Env, `env`, "e", "dev", `"dev", "prod"`)

	Cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		validDatabases := map[string]bool{"postgres": true, "mysql": true}
		if !validDatabases[global.Database] {
			return erx.Errorf("invalid database value: %s. Must be one of: postgres, mysql",
				global.Database)
		}

		validEnvs := map[string]bool{"dev": true, "prod": true}
		if !validEnvs[global.Env] {
			return erx.Errorf("invalid environment value: %s. Must be one of: dev, prod", global.Env)
		}

		port, _ := cmd.Flags().GetString("port")
		if _, err := strconv.Atoi(port); err != nil || port == "" {
			return erx.Errorf("invalid port value: %s. Must be a valid number", port)
		}

		return nil
	}
}
