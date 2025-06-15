package cmd

import (
	"dailyalu-server/internal/container"
	"dailyalu-server/internal/router"
	"dailyalu-server/pkg/app_log/zap_log"
	"dailyalu-server/pkg/db/postgres"
	"dailyalu-server/pkg/mailer/smtp"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the API server",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize logger
		if err := zap_log.Init(
			viper.GetString("logging.level"),
			viper.GetString("logging.format"),
			viper.GetString("logging.file_directory"),
		); err != nil {
			return fmt.Errorf("failed to initialize logger: %w", err)
		}

		// Initialize database
		db, err := postgres.NewConnection(postgres.Config{
			Host:     viper.GetString("database.host"),
			Port:     viper.GetInt("database.port"),
			User:     viper.GetString("database.user"),
			Password: viper.GetString("database.password"),
			DBName:   viper.GetString("database.name"),
			SSLMode:  viper.GetString("database.sslmode"),
		})
		if err != nil {
			return fmt.Errorf("failed to connect to database: %w", err)
		}

		// newSes, err := ses.InitSes(context.TODO())
		// if err != nil {
		// 	return fmt.Errorf("failed to get mailer service: %w", err)
		// }

		newSmtp := smtp.InitSmtp()

		// Initialize dependency container
		cont := container.NewContainer(
			db,
			newSmtp,
			viper.GetString("jwt.secret"),
			viper.GetString("jwt.refresh-secret-key"),
			viper.GetDuration("jwt.expiry")*time.Hour,
			viper.GetDuration("jwt.refresh-expiry")*time.Hour,
		)
		defer cont.Close()

		// Initialize Fiber app
		app := fiber.New(fiber.Config{
			AppName: "DailyAlu API Server",
		})

		// Add global middleware
		app.Use(cont.GetErrorMiddleware().Handle())

		// Setup routes
		router.SetupUserRoutes(
			app,
			cont.GetUserHandler(),
			cont.GetSecurityMiddleware(),
		)

		router.SetupActivityRoutes(
			app,
			cont.GetActivityHandler(),
			cont.GetSecurityMiddleware(),
		)

		router.SetupToolsRoutes(
			app,
			cont.GetSecurityMiddleware(),
		)

		// router.SetupChildrenRoutes(
		// 	app,
		// 	cont.GetChildrenHandler(),
		// 	cont.GetSecurityMiddleware(),
		// )

		// Start server
		port := viper.GetInt("server.port")
		// Bind to all network interfaces (0.0.0.0) instead of just localhost
		//return app.Listen(fmt.Sprintf("0.0.0.0:%d", port))

		return app.Listen(fmt.Sprintf(":%d", port))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
