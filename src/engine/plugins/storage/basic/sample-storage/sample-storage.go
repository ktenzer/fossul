/*
Copyright 2019 The Fossul Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"encoding/json"
	"fmt"
	"fossul/src/engine/util"
	"github.com/pborman/getopt/v2"
	"os"
)

func main() {
	optBackup := getopt.BoolLong("backup", 0, "Backup")
	optRestore := getopt.BoolLong("restore", 0, "Restore")
	optBackupList := getopt.BoolLong("backupList", 0, "Backup List")
	optBackupDelete := getopt.BoolLong("backupDelete", 0, "Backup Delete")
	optInfo := getopt.BoolLong("info", 0, "Storage Plugin Information")
	optHelp := getopt.BoolLong("help", 0, "Help")
	getopt.Parse()

	if *optHelp {
		getopt.Usage()
		os.Exit(0)
	}

	//load env parameters
	configMap := getEnvParams()

	if *optBackup {
		backup(configMap)
	} else if *optRestore {
		restore(configMap)
	} else if *optBackupList {
		backupList(configMap)
	} else if *optBackupDelete {
		backupDelete(configMap)
	} else if *optInfo {
		info()
	} else {
		getopt.Usage()
		os.Exit(0)
	}
}

func backup(configMap map[string]string) {
	printEnv(configMap)
	fmt.Println("INFO *** Backup ***")
}

func restore(configMap map[string]string) {
	printEnv(configMap)
	fmt.Println("INFO *** Restore ***")
}

func backupList(configMap map[string]string) {
	printEnv(configMap)
	fmt.Println("INFO *** Backup list ***")
}

func backupDelete(configMap map[string]string) {
	printEnv(configMap)
	fmt.Println("INFO *** Backup delete ***")
}

func info() {
	var plugin util.Plugin = setPlugin()

	//output json
	b, err := json.Marshal(plugin)
	if err != nil {
		fmt.Println("ERROR " + err.Error())
	} else {
		fmt.Println(string(b))
	}
}

func setPlugin() (plugin util.Plugin) {
	plugin.Name = "sample-storage"
	plugin.Description = "A sample storage plugin"
	plugin.Version = "1.0.0"
	plugin.Type = "storage"

	var capabilities []util.Capability
	var backupCap util.Capability
	backupCap.Name = "backup"

	var backupListCap util.Capability
	backupListCap.Name = "backupList"

	var backupDeleteCap util.Capability
	backupDeleteCap.Name = "backupDelete"

	var infoCap util.Capability
	infoCap.Name = "info"

	capabilities = append(capabilities, backupCap, backupListCap, backupDeleteCap, infoCap)

	plugin.Capabilities = capabilities

	return plugin
}

func printEnv(configMap map[string]string) {
	config, err := util.ConfigMapToJson(configMap)
	if err != nil {
		fmt.Println("ERROR " + err.Error())
	}
	fmt.Println("DEBUG Config Parameters: " + config + "\n")
}

func getEnvParams() map[string]string {
	configMap := map[string]string{}

	configMap["ProfileName"] = os.Getenv("ProfileName")
	configMap["ConfigName"] = os.Getenv("ConfigName")
	configMap["BackupName"] = os.Getenv("BackupName")
	configMap["SelectedWorkflowId"] = os.Getenv("SelectedWorkflowId")
	configMap["AutoDiscovery"] = os.Getenv("AutoDiscovery")
	configMap["DataFilePaths"] = os.Getenv("DataFilePaths")
	configMap["LogFilePaths"] = os.Getenv("LogFilePaths")
	configMap["BackupPolicy"] = os.Getenv("BackupPolicy")
	configMap["SampleStorageVar1"] = os.Getenv("SampleStorageVar1")
	configMap["SampleStorageVar2"] = os.Getenv("SampleStorageVar2")

	return configMap
}
