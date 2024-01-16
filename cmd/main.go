package main

import (
	"GoGinBoilerPlate/server"
	"GoGinBoilerPlate/server/props"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

var (
	properties *viper.Viper
)

const envVarDDService = "DD_SERVICE"

func main() {
	fmt.Println("Hello, World!")
	// read all the properties from os.environ or .env
	properties = initConfiguration("./env")

	//set up server
	serverProperties, err := readServerConfig()
	srv := server.NewServer(serverProperties) //&sc
	defer srv.Shutdown()
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

func readServerConfig() (*props.Properties, error) {

	yamlFile, err := os.ReadFile("./env/application.yaml")

	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	var data props.Properties
	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return &data, err
}
