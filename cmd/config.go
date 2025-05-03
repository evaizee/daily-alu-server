package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
}

var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get configuration value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := viper.Get(key)
		if value == nil {
			return fmt.Errorf("configuration key not found: %s", key)
		}
		fmt.Printf("%v\n", value)
		return nil
	},
}

func init() {
	configCmd.AddCommand(configGetCmd)
	rootCmd.AddCommand(configCmd)

	// Add default configuration values
	viper.SetDefault("server.port", 3000)
	viper.SetDefault("server.env", "development")
	viper.SetDefault("server.apikey", "") // Empty string means no master key
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("jwt.expiry", 24)
	viper.SetDefault("logging.level", "debug")
	viper.SetDefault("logging.format", "json")

	//add ses values
	viper.SetDefault("aws.ses.sender", "no-reply@dailyalu.com")
	viper.SetDefault("aws.ses.access_key", "")
	viper.SetDefault("aws.ses.access_secret_key", "")
	viper.SetDefault("aws.region", "ap-southeast-1")

	// Rate limiter configuration
	viper.SetDefault("ratelimit.enabled", true)
	viper.SetDefault("ratelimit.default.max", 60)         // 60 requests
	viper.SetDefault("ratelimit.default.expiration", 60)  // 60 seconds (1 minute)

	// Endpoint-specific rate limits
	viper.SetDefault("ratelimit.endpoints.post.api_v1_auth_login.max", 100)
	viper.SetDefault("ratelimit.endpoints.post.api_v1_auth_login.expiration", 60)
	viper.SetDefault("ratelimit.endpoints.post.api_v1_auth_register.max", 20)
	viper.SetDefault("ratelimit.endpoints.post.api_v1_auth_register.expiration", 60)

	// Enable environment variable overrides
	viper.SetEnvPrefix("DAILYALU")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}
