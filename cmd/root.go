package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/adrg/xdg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	appName     = "ghostship"
	cfgFile, _  = xdg.ConfigFile(path.Join(appName, "config.toml"))
	cfgDir      = filepath.Dir(cfgFile)
	userLicense string

	buildParams BuildParams

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "ghostship",
		Short: "ghostship manages your shell's prompt",
		Long:  `ghostship manages your shell's prompt`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				_ = cmd.Help()
				os.Exit(0)
			}
		},
	}
)

// InitCobra adds all child commands to the root command and sets flags appropriately.
// This is called by main.init(). It only needs to happen once to the rootCmd.
func InitCobra(bp BuildParams) error {
	buildParams = bp
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	var debug bool
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debug logging")
	_ = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", cfgFile, "default config file")
	_ = viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
}

func initConfig() {
	if viper.GetBool("debug") {
		log.SetLevel(log.DebugLevel)
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(xdg.Home)
		viper.AddConfigPath(cfgDir)
		viper.AddConfigPath(path.Join(xdg.ConfigHome, appName))
		viper.SetConfigType("toml")
		viper.SetConfigName(appName)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Warnf("Warning: config file not found: %v, %v", cfgFile, err)
		} else {
			log.Warnf("Warning: Error in processing config file: %v, %v", cfgFile, err)
		}
	}
	log.Debugf("Using config file: %v", viper.ConfigFileUsed())

	rootCmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Determine the naming convention of the flags when
		// represented in the config file
		configName := f.Name
		log.Debugf("Processing flag: %v", f.Name)

		// Apply the viper config value to the flag
		// when the flag is not set and viper has a value
		if !f.Changed && viper.IsSet(configName) {
			val := viper.Get(configName)
			_ = rootCmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
