package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/yaml"
)

type Config struct {
	AppConfig string `yaml:"appname"`
}

type AppConfig struct {
	path string
}

func init() {
	// log as JSON
	log.SetFormatter(&log.JSONFormatter{})

	// Output everything including stderr to stdout
	log.SetOutput(os.Stdout)

	// set level
	log.SetLevel(log.InfoLevel)
}

func (s *AppConfig) configHandler(w http.ResponseWriter, req *http.Request) {
	
	log.info("Config Path is:" + (*s).path)

	yamlFile, err := ioutil.ReadFile((*s).path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Error(err.Error())
		return
	}
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Error(err.Error())
		return
	}
        log.info("The config is: " + config.AppConfig)

	w.Write([]byte("The config is: " + config.AppConfig))
	log.Info(req.Method + req.RequestURI + " " + req.Proto)
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("hello"))
	log.Info(req.Method + req.RequestURI + " " + req.Proto)
}

func headersHandler(w http.ResponseWriter, req *http.Request) {
	jsonHeaders, err := json.Marshal((*req).Header)
	if err != nil {
		log.Error(err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonHeaders)
	log.Info(req.Method + req.RequestURI + " " + req.Proto)
}

func main() {

	port := ":8080"

	s := AppConfig{
		//path: "resources/config.yaml",
		//path: "/var/run/secrets/vaultproject.io/config.yaml",
		path: os.Getenv("APP_CONFIG_PATH"),
	}
	http.HandleFunc("/config", s.configHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/headers", headersHandler)

	log.Info("Listening on port ", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
