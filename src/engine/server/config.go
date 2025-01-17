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
	"fossul/src/engine/util"
)

func GetConsolidatedConfig(profileName, configName, policyName string) (util.Config, error) {
	conf := configDir + "/" + profileName + "/" + configName + "/" + configName + ".conf"
	config, err := util.ReadConfig(conf)

	if err != nil {
		return config, err
	}

	config.ProfileName = profileName
	config.ConfigName = configName

	backupRetention := util.GetBackupRetention(policyName, config.BackupRetentions)
	config.SelectedBackupRetention = backupRetention
	archiveRetention := util.GetArchiveRetention(policyName, config.ArchiveRetentions)
	config.SelectedArchiveRetention = archiveRetention
	config.SelectedBackupPolicy = policyName

	if config.AppPlugin != "" {
		appConf := configDir + "/" + profileName + "/" + configName + "/" + config.AppPlugin + ".conf"
		appConfigMap, err := util.ReadConfigToMap(appConf)

		if err != nil {
			return config, err
		}

		config.AppPluginParameters = appConfigMap
	}

	if config.StoragePlugin != "" {
		storageConf := configDir + "/" + profileName + "/" + configName + "/" + config.StoragePlugin + ".conf"
		storageConfigMap, err := util.ReadConfigToMap(storageConf)

		if err != nil {
			return config, err
		}

		config.StoragePluginParameters = storageConfigMap
	}

	if config.ArchivePlugin != "" {
		archiveConf := configDir + "/" + profileName + "/" + configName + "/" + config.ArchivePlugin + ".conf"
		archiveConfConfigMap, err := util.ReadConfigToMap(archiveConf)

		if err != nil {
			return config, err
		}

		config.ArchivePluginParameters = archiveConfConfigMap
	}

	return config, nil
}
