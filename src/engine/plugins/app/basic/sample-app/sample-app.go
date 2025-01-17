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
	"strings"
)

func main() {
	optDiscover := getopt.BoolLong("discover", 0, "Application Discover")
	optQuiesce := getopt.BoolLong("quiesce", 0, "Application Quiesce")
	optUnquiesce := getopt.BoolLong("unquiesce", 0, "Application Unquiesce")
	optPreRestore := getopt.BoolLong("preRestore", 0, "Application Pre Restore")
	optPostRestore := getopt.BoolLong("postRestore", 0, "Application Post Restore")
	optInfo := getopt.BoolLong("info", 0, "Application Plugin Information")
	optHelp := getopt.BoolLong("help", 0, "Help")
	getopt.Parse()

	if *optHelp {
		getopt.Usage()
		os.Exit(0)
	}

	//load env parameters
	configMap := getEnvParams()

	if *optQuiesce {
		printEnv(configMap)
		quiesce(configMap)
	} else if *optUnquiesce {
		printEnv(configMap)
		unquiesce(configMap)
	} else if *optPreRestore {
		printEnv(configMap)
		preRestore(configMap)
	} else if *optPostRestore {
		printEnv(configMap)
		postRestore(configMap)
	} else if *optInfo {
		info()
	} else if *optDiscover {
		discover()
	} else {
		getopt.Usage()
		os.Exit(0)
	}
}

func discover() {
	var discoverResult util.DiscoverResult = setDiscoverResult()

	//output json
	b, err := json.Marshal(discoverResult)
	if err != nil {
		fmt.Println("ERROR " + err.Error())
	} else {
		fmt.Println(string(b))
	}
}

func quiesce(configMap map[string]string) {
	printEnv(configMap)
	fmt.Println("INFO *** Application quiesce ***")
}

func unquiesce(configMap map[string]string) {
	printEnv(configMap)
	fmt.Println("INFO *** Application unquiesce ***")
}

func preRestore(configMap map[string]string) {
	printEnv(configMap)
	fmt.Println("INFO *** Application Pre Restore ***")
}

func postRestore(configMap map[string]string) {
	printEnv(configMap)
	fmt.Println("INFO *** Application Post Restore ***")
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

func setDiscoverResult() (discoverResult util.DiscoverResult) {
	var data []string
	data = append(data, "/path/to/data/file1")
	data = append(data, "/path/to/data/file2")

	var logs []string
	logs = append(logs, "/path/to/logs/file1")
	logs = append(logs, "/path/to/logs/file2")

	var discoverInst1 util.Discover
	discoverInst1.Instance = "inst1"
	discoverInst1.DataFilePaths = data
	discoverInst1.LogFilePaths = logs

	var discoverInst2 util.Discover
	discoverInst2.Instance = "inst2"
	discoverInst2.DataFilePaths = data
	discoverInst2.LogFilePaths = logs

	var discoverList []util.Discover
	discoverList = append(discoverList, discoverInst1)
	discoverList = append(discoverList, discoverInst2)

	var messages []util.Message
	msg := util.SetMessage("INFO", "*** Application Discovery ***")
	messages = append(messages, msg)

	for _, discover := range discoverList {
		dataFiles := strings.Join(discover.DataFilePaths, " ")
		logFiles := strings.Join(discover.LogFilePaths, " ")
		msg := util.SetMessage("INFO", "Instance ["+discover.Instance+"] data files: ["+dataFiles+"] log files: ["+logFiles+"]")
		messages = append(messages, msg)
	}

	result := util.SetResult(0, messages)
	discoverResult.Result = result
	discoverResult.DiscoverList = discoverList

	return discoverResult
}

func setPlugin() (plugin util.Plugin) {
	plugin.Name = "sample"
	plugin.Description = "A sample plugin"
	plugin.Version = "1.0.0"
	plugin.Type = "app"

	var capabilities []util.Capability
	var discoverCap util.Capability
	discoverCap.Name = "discover"

	var quiesceCap util.Capability
	quiesceCap.Name = "quiesce"

	var unquiesceCap util.Capability
	unquiesceCap.Name = "unquiesce"

	var preRestoreCap util.Capability
	preRestoreCap.Name = "preRestore"

	var postRestoreCap util.Capability
	postRestoreCap.Name = "postRestore"

	var infoCap util.Capability
	infoCap.Name = "info"

	capabilities = append(capabilities, discoverCap, quiesceCap, unquiesceCap, preRestoreCap, postRestoreCap, infoCap)

	plugin.Capabilities = capabilities

	return plugin
}

func printEnv(configMap map[string]string) {
	config, err := util.ConfigMapToJson(configMap)
	if err != nil {
		fmt.Println("ERROR " + err.Error())
	}
	fmt.Println("DEBUG Config Parameters: " + config)
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
	configMap["SampleAppVar1"] = os.Getenv("SampleAppVar1")
	configMap["SampleAppVar2"] = os.Getenv("SampleAppVar2")

	return configMap
}
