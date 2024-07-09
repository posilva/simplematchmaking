package cmd

import (
	"os"

	"github.com/posilva/simplematchmaking/cmd/simplematchmaking/app"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "simplematchmaking",
	Short: "Service for Matchmaking",
	Long: ` Service to handle Matchmaking manager
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		app.Run()
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.BindPFlag("local", cmd.Flags().Lookup("local"))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.simplematchmaking.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("local", "l", false, "Run the service locally against using docker compose")
}
