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
	"fossul/src/engine/util"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// GetStatus godoc
// @Description Status and version information for the service
// @Accept  json
// @Produce  json
// @Success 200 {string} string
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /status [get]
func GetStatus(w http.ResponseWriter, r *http.Request) {
	var status util.Status
	status.Msg = "OK"
	status.Version = version

	json.NewEncoder(w).Encode(status)
}

// PluginList godoc
// @Description List storage or archive plugins
// @Param pluginType path string true "plugin type (storage|archive)"
// @Accept  json
// @Produce  json
// @Success 200 {array} string
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /pluginList/{pluginType} [get]
func PluginList(w http.ResponseWriter, r *http.Request) {
	var plugins []string

	params := mux.Vars(r)
	var pluginType string = params["pluginType"]

	var storagePluginDir string
	if pluginType == "storage" {
		storagePluginDir = pluginDir + "/storage"
	} else if pluginType == "archive" {
		storagePluginDir = pluginDir + "/archive"
	} else {
		log.Println("ERROR plugin type " + pluginType + " must be storage|archive")

		_ = json.NewDecoder(r.Body).Decode(&plugins)
		json.NewEncoder(w).Encode(plugins)
	}

	fileInfo, err := ioutil.ReadDir(storagePluginDir)
	if err != nil {
		log.Println(err)
	}

	for _, file := range fileInfo {
		plugins = append(plugins, file.Name())
	}

	_ = json.NewDecoder(r.Body).Decode(&plugins)
	json.NewEncoder(w).Encode(plugins)
}

// PluginInfo godoc
// @Description Plugin information and version
// @Param config body util.Config true "config struct"
// @Param pluginName path string true "name of plugin"
// @Param pluginType path string true "plugin type (storage|archive)"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.PluginInfoResult
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /pluginInfo/{pluginName}/{pluginType} [post]
func PluginInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var pluginName string = params["pluginName"]
	var pluginType string = params["pluginType"]

	var pluginInfoResult util.PluginInfoResult
	var pluginInfo util.Plugin
	var result util.Result
	var messages []util.Message

	config, err := util.GetConfig(w, r)
	printConfigDebug(config)

	if err != nil {
		message := util.SetMessage("ERROR", "Couldn't read config! "+err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)
		pluginInfoResult.Result = result

		_ = json.NewDecoder(r.Body).Decode(&pluginInfoResult)
		json.NewEncoder(w).Encode(pluginInfoResult)

		return
	}

	pluginPath := util.GetPluginPath(pluginName)

	if pluginPath == "" {
		var plugin string = pluginDir + "/" + pluginType + "/" + pluginName
		if _, err := os.Stat(plugin); os.IsNotExist(err) {
			msg := util.SetMessage("ERROR", "Plugin not found! "+err.Error())
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			pluginInfoResult.Result = result

			_ = json.NewDecoder(r.Body).Decode(&pluginInfoResult)
			json.NewEncoder(w).Encode(pluginInfoResult)
		}

		var resultSimple util.ResultSimple
		resultSimple = util.ExecutePluginSimple(config, pluginType, plugin, "--info")
		if resultSimple.Code != 0 {
			msg := util.SetMessage("ERROR", "Plugin Info failed!")
			messages = append(messages, msg)
			result := util.SetResult(1, messages)
			pluginInfoResult.Result = result

			_ = json.NewDecoder(r.Body).Decode(&pluginInfoResult)
			json.NewEncoder(w).Encode(pluginInfoResult)
		} else {
			pluginInfoString := strings.Join(resultSimple.Messages, " ")
			json.Unmarshal([]byte(pluginInfoString), &pluginInfo)

			pluginInfoResult.Result.Code = resultSimple.Code
			pluginInfoResult.Plugin = pluginInfo

			_ = json.NewDecoder(r.Body).Decode(&pluginInfoResult)
			json.NewEncoder(w).Encode(pluginInfoResult)
		}
	} else {
		if pluginType == "storage" {
			plugin, err := util.GetStorageInterface(pluginPath)

			if err != nil {
				msg := util.SetMessage("ERROR", err.Error())
				messages = append(messages, msg)
				result = util.SetResult(1, messages)
				pluginInfoResult.Result = result

				_ = json.NewDecoder(r.Body).Decode(&pluginInfoResult)
				json.NewEncoder(w).Encode(pluginInfoResult)
			} else {
				pluginInfo := plugin.Info()

				pluginInfoResult.Result.Code = 0
				pluginInfoResult.Plugin = pluginInfo

				_ = json.NewDecoder(r.Body).Decode(&pluginInfoResult)
				json.NewEncoder(w).Encode(pluginInfoResult)
			}
		} else if pluginType == "archive" {
			plugin, err := util.GetArchiveInterface(pluginPath)
			if err != nil {
				msg := util.SetMessage("ERROR", err.Error())
				messages = append(messages, msg)
				result = util.SetResult(1, messages)
				pluginInfoResult.Result = result

				_ = json.NewDecoder(r.Body).Decode(&pluginInfoResult)
				json.NewEncoder(w).Encode(pluginInfoResult)
			} else {
				pluginInfo := plugin.Info()

				pluginInfoResult.Result.Code = 0
				pluginInfoResult.Plugin = pluginInfo

				_ = json.NewDecoder(r.Body).Decode(&pluginInfoResult)
				json.NewEncoder(w).Encode(pluginInfoResult)
			}
		} else {
			msg := util.SetMessage("ERROR", "Invalid plugin type ["+pluginType+"], must be app|archive")
			messages = append(messages, msg)
			result = util.SetResult(1, messages)
			pluginInfoResult.Result = result

			_ = json.NewDecoder(r.Body).Decode(&pluginInfoResult)
			json.NewEncoder(w).Encode(pluginInfoResult)
		}
	}
}
