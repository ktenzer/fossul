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
package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fossul/src/engine/util"
	"log"
	"net/http"
)

func ArchivePluginList(auth Auth, pluginType string) ([]string, error) {
	var plugins []string

	req, err := http.NewRequest("GET", "http://"+auth.StorageHostname+":"+auth.StoragePort+"/pluginList/"+pluginType, nil)
	if err != nil {
		return plugins, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return plugins, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&plugins); err != nil {
			return plugins, err
		}
	} else {
		return plugins, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return plugins, nil

}

func ArchivePluginInfo(auth Auth, config util.Config, pluginName, pluginType string) (util.PluginInfoResult, error) {
	var pluginInfoResult util.PluginInfoResult

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://"+auth.StorageHostname+":"+auth.StoragePort+"/pluginInfo/"+pluginName+"/"+pluginType, b)
	if err != nil {
		return pluginInfoResult, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return pluginInfoResult, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&pluginInfoResult); err != nil {
			return pluginInfoResult, err
		}
	} else {
		return pluginInfoResult, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return pluginInfoResult, nil
}

func Archive(auth Auth, config util.Config) (util.Result, error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://"+auth.StorageHostname+":"+auth.StoragePort+"/archive", b)
	if err != nil {
		log.Println("NewRequest: ", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result, err
		}
	} else {
		return result, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result, nil
}

func ArchiveList(auth Auth, profileName, configName, policyName string, config util.Config) (util.Archives, error) {
	var archives util.Archives

	config = SetAdditionalConfigParams(profileName, configName, policyName, config)

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://"+auth.StorageHostname+":"+auth.StoragePort+"/archiveList", b)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return archives, err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return archives, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&archives); err != nil {
			return archives, err
		}
	} else {
		return archives, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return archives, nil
}

func ArchiveDelete(auth Auth, config util.Config) (util.Result, error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://"+auth.StorageHostname+":"+auth.StoragePort+"/archiveDelete", b)
	if err != nil {
		return result, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result, err
		}
	} else {
		return result, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result, nil
}
