package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/spf13/viper"
)

type Configuration struct {
	GithubToken string `mapstructure:"token"`
}

func loadConstants(homeDir string){
	viper.Set("HOME", homeDir)
	agenda := path.Join(homeDir, "/Documents/Agenda/")
	viper.Set("AGENDA", agenda)
	viper.Set("PROJECTS", agenda+"/projects")

}
func LoadConfig() {
	var config Configuration
	viper.SetConfigName(".GoSeq")
	viper.SetConfigType("yaml")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Unable to find home directory: %v", err)
	}

	loadConstants(homeDir)
	viper.AddConfigPath(homeDir + "/.config/")
	viper.AddConfigPath(homeDir)

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {

			fmt.Printf("Configuration file not found.\nOpen https://github.com/settings/tokens\nPlease enter your GitHub token:\n")
			reader := bufio.NewReader(os.Stdin)
			token, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalf("Error reading input: %v", err)
			}
			config.GithubToken = strings.TrimSpace(token)

			viper.Set("token", config.GithubToken)

			configFilePath := homeDir + "/.config/.GoSeq.yaml"

			err = viper.SafeWriteConfigAs(configFilePath)
			if err != nil {
				if os.IsExist(err) {
					err = viper.WriteConfigAs(configFilePath)
				}
				if err != nil {
					log.Fatalf("Error writing config file: %v", err)
				}
			}
		} else {
			log.Fatal(err)
		}
	} else {
		err = viper.Unmarshal(&config)
		if err != nil {
			log.Fatalf("Unable to decode into struct: %v", err)
		}
	}

}
