package main

import (
	"os"
	"fmt"
	"github.com/pborman/getopt/v2"
	"engine/util"
	"encoding/json"
	"engine/util/k8s"
)

func main() {
	optAction := getopt.StringLong("action",'a',"","backup|backupList|backupDelete|info")
	optHelp := getopt.BoolLong("help", 0, "Help")
	getopt.Parse()

	if *optHelp {
		getopt.Usage()
		os.Exit(0)
	}

	if getopt.IsSet("action") != true {
		fmt.Println("ERROR: incorrect parameter")
		getopt.Usage()
		os.Exit(1)
	}

	//load env parameters
	var namespace string = os.Getenv("Namespace")
	var serviceName string = os.Getenv("ServiceName")
	var accessWithinCluster string = os.Getenv("AccessWithinCluster")

	if *optAction == "backup" {
		backup(namespace,serviceName,accessWithinCluster)
	} else if *optAction == "backupList" {
		backupList()
	} else if *optAction == "backupDelete" {
		backupDelete()		
	} else if *optAction == "info" {
		info()			
	} else {
		fmt.Println("ERROR: incorrect parameter", *optAction)
		getopt.Usage()
		os.Exit(1)
	}
}	

func backup (namespace,serviceName,accessWithinCluster string) {
	fmt.Println("Performing container backup")
	podName := k8s.GetPod(namespace,serviceName,accessWithinCluster)
	fmt.Println("found pod", podName)


}

func backupList () {
	fmt.Println("Performing backup list")
}

func backupDelete () {
	fmt.Println("Performing backup delete")
}

func info () {
	var plugin util.Plugin = setPlugin()

	//output json
	b, err := json.Marshal(plugin)
    if err != nil {
        fmt.Println(err)
        return
	}
	
	fmt.Println(string(b))
}

func setPlugin() (plugin util.Plugin) {
	plugin.Name = "sample"
	plugin.Description = "A sample plugin"
	plugin.Type = "app"

	var capabilities []util.Capability
	var backupCap util.Capability
	backupCap.Name = "backup"

	var backupListCap util.Capability
	backupListCap.Name = "backupList"

	var backupDeleteCap util.Capability
	backupDeleteCap.Name = "backupDelete"

	var infoCap util.Capability
	infoCap.Name = "info"

	capabilities = append(capabilities,backupCap,backupListCap,backupDeleteCap,infoCap)

	plugin.Capabilities = capabilities

	return plugin
}