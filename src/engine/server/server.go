package main
 
import (
    "log"
    "net/http"
    "fossil/src/engine/util"
    "os"
)

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

var runningWorkflowMap map[string]string = make(map[string]string)
 
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

    StartCron()
    err = LoadCronSchedules()
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Starting server service on port [" + port + "]")
    log.Fatal(http.ListenAndServe(":" + port, router))
}