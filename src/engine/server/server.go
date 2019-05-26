package main
 
import (
    "log"
    "net/http"
    "fossil/src/engine/util"
    "os"
	"github.com/swaggo/http-swagger"
    _ "fossil/src/engine/server/docs"
)

const version = "1.0.0"

var port string = os.Getenv("FOSSIL_SERVER_SERVICE_PORT")
var configDir string = os.Getenv("FOSSIL_SERVER_CONFIG_DIR")
var dataDir string = os.Getenv("FOSSIL_SERVER_DATA_DIR")
var myUser string = os.Getenv("FOSSIL_USERNAME")
var myPass string = os.Getenv("FOSSIL_PASSWORD")
var serverHostname string = os.Getenv("FOSSIL_SERVER_CLIENT_HOSTNAME")
var serverPort string = os.Getenv("FOSSIL_SERVER_CLIENT_PORT")
var appHostname string = os.Getenv("FOSSIL_APP_CLIENT_HOSTNAME")
var appPort string = os.Getenv("FOSSIL_APP_CLIENT_PORT")
var storageHostname string = os.Getenv("FOSSIL_STORAGE_CLIENT_HOSTNAME")
var storagePort string = os.Getenv("FOSSIL_STORAGE_CLIENT_PORT")
var debug string = os.Getenv("FOSSIL_SERVER_DEBUG")

var runningWorkflowMap map[string]string = make(map[string]string)

// @title Fossil Framework Server API
// @version 1.0
// @description APIs for managing Fossil workflows, jobs, profile, and configs
// @description JSON API definition can be retrieved at <a href="/api/v1/swagger/doc.json">/api/v1/swagger/doc.json</a>
// @termsOfService http://swagger.io/terms/

// @contact.name Keith Tenzer
// @contact.url http://www.keithtenzer.com
// @contact.email keith.tenzer@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

func main() {
    
    log.Println("Configs directory [" + configDir + "] Data directory [" + dataDir + "]")
    err := util.CreateDir(configDir,0755)
    if err != nil {
        log.Fatal(err)
    }

    err = util.CreateDir(dataDir,0755)
    if err != nil {
        log.Fatal(err)
    }

    router := NewRouter()
    router.PathPrefix("/api/v1").Handler(httpSwagger.WrapHandler)

    StartCron()
    err = LoadCronSchedules()
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Starting server service on port [" + port + "]")
    log.Fatal(http.ListenAndServe(":" + port, router))
}

func printConfigDebug(config util.Config) {
    if debug == "true" {
		log.Println("[DEBUG]",config)
	}
}

func printConfigMapDebug(configMap map[string]string) {
    if debug == "true" {
		log.Println("[DEBUG]",configMap)
	}
}