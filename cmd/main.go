package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

var (
	properties *viper.Viper
)

const envVarDDService = "DD_SERVICE"

func main() {
	fmt.Println("Hello, World!")
	// read all the properties from os.environ or .env
	properties = initConfiguration("./env")
}

func initConfiguration(path string) *viper.Viper {
	viperConfigManager := viper.NewWithOptions(viper.KeyDelimiter("_"))
	viperConfigManager.SetConfigName("application")
	viperConfigManager.SetConfigType("yaml")
	viperConfigManager.AddConfigPath(path)
	err := viperConfigManager.BindEnv(envVarDDService)
	if err != nil {
		log.Println("error in init configuration")
		//log.Warnf("Failed to bind a configuration key to the '%v' environment variable with error %v",
		//envVarDDService, err)
	}

	viperConfigManager.AutomaticEnv()
	viperConfigManager.AllowEmptyEnv(true)

	viperConfigManager.WatchConfig()
	viperConfigManager.OnConfigChange(func(e fsnotify.Event) {
		// TODO - notify observers
		//log.Infof("config file changed: %s", e.Name)
	})

	err = viperConfigManager.ReadInConfig()
	if err != nil {
		log.Fatal()
		//log.Fatal(fmt.Errorf("unable to start %s due to missing applicaiton config %v", serviceName, err))
	}

	//log.Infof("loading application config from %v", viperConfigManager.ConfigFileUsed())
	log.Printf("loading application config from %v", viperConfigManager.ConfigFileUsed())
	return viperConfigManager
}
