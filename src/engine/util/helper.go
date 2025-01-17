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
package util

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func GetTimestamp() int64 {
	time := time.Now().Unix()

	return time
}

func GetBackupDirFromMap(configMap map[string]string) string {
	backupPath := configMap["BackupDestPath"] + "/" + configMap["ProfileName"] + "/" + configMap["ConfigName"]

	return backupPath
}

func GetBackupDirFromConfig(config Config) string {
	backupPath := config.StoragePluginParameters["BackupDestPath"] + "/" + config.ProfileName + "/" + config.ConfigName

	return backupPath
}

func GetBackupPathFromMap(configMap map[string]string) string {
	backupName := GetBackupName(configMap["BackupName"], configMap["BackupPolicy"], configMap["WorkflowId"], configMap["WorkflowTimestamp"])
	backupPath := configMap["BackupDestPath"] + "/" + configMap["ProfileName"] + "/" + configMap["ConfigName"] + "/" + backupName

	return backupPath
}

func GetBackupPathFromConfig(config Config) string {
	timestampToString := fmt.Sprintf("%d", config.WorkflowTimestamp)
	backupName := GetBackupName(config.StoragePluginParameters["BackupName"], config.SelectedBackupPolicy, config.WorkflowId, timestampToString)
	backupPath := config.StoragePluginParameters["BackupDestPath"] + "/" + config.ProfileName + "/" + config.ConfigName + "/" + backupName

	return backupPath
}

func GetBackupName(name, policy, workflowId, timestamp string) string {
	backupName := fmt.Sprintf(name + "_" + policy + "_" + workflowId + "_" + timestamp)

	return backupName
}

func GetRestoreSrcPath(config Config) (string, error) {
	backupPath := config.StoragePluginParameters["BackupDestPath"] + "/" + config.ProfileName + "/" + config.ConfigName
	backupNameSubString := config.StoragePluginParameters["BackupName"] + "_" + config.SelectedBackupPolicy + "_" + IntToString(config.SelectedWorkflowId)

	fmt.Println("[DEBUG] restore path [" + backupPath + "] search string [" + backupNameSubString + "]")
	files, err := ioutil.ReadDir(backupPath)
	if err != nil {
		return "", err
	}
	for _, f := range files {
		if strings.Contains(f.Name(), backupNameSubString) {
			return backupPath + "/" + f.Name(), nil
		}
	}
	return "", nil
}

func GetRestoreSrcPathFromMap(configMap map[string]string) (string, error) {
	backupPath := configMap["BackupDestPath"] + "/" + configMap["ProfileName"] + "/" + configMap["ConfigName"]
	backupNameSubString := configMap["BackupName"] + "_" + configMap["BackupPolicy"] + "_" + configMap["SelectedWorkflowId"]

	fmt.Println("[DEBUG] restore path [" + backupPath + "] search string [" + backupNameSubString + "]")
	files, err := ioutil.ReadDir(backupPath)
	if err != nil {
		return "", err
	}
	for _, f := range files {
		if strings.Contains(f.Name(), backupNameSubString) {
			return backupPath + "/" + f.Name(), nil
		}
	}
	return "", nil
}

func ConvertEpoch(epoch string) string {
	i := StringToInt64(epoch)
	time := time.Unix(i, 0)

	return time.String()
}

func ConvertEpochToTime(epoch string) time.Time {
	i := StringToInt64(epoch)
	time := time.Unix(i, 0)

	return time
}

func JoinArray(array, combinedArray []string) []string {
	for _, item := range array {
		combinedArray = append(combinedArray, item)
	}

	return combinedArray
}

func ExistsInArray(array []string, str string) bool {
	for _, item := range array {
		if item == str {
			return true
		}
	}
	return false
}

func WriteGob(filePath string, object interface{}) error {
	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

func ReadGob(filePath string, object interface{}) error {

	file, err := os.Open(filePath)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}

func IsDirectory(path string) (bool, error) {
	fileOrDir, err := os.Stat(path)

	if err != nil {
		return false, err
	}
	switch mode := fileOrDir.Mode(); {
	case mode.IsDir():
		return true, err
	case mode.IsRegular():
		return false, err
	}

	return false, err
}

func ExistsPath(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func CreateDir(path string, mode os.FileMode) error {

	if ExistsPath(path) == false {
		if err := os.MkdirAll(path, mode); err != nil {
			return err
		} else {
			return nil
		}
	}
	return nil
}

func RecursiveDirDelete(dir string) error {
	if ExistsPath(dir) == true {
		d, err := os.Open(dir)

		if err != nil {
			return err
		}
		defer d.Close()

		names, err := d.Readdirnames(-1)
		if err != nil {
			return err
		}

		for _, name := range names {
			err = os.RemoveAll(filepath.Join(dir, name))
			if err != nil {
				return err
			}
		}

		err = os.Remove(dir)
		if err != nil {
			return err
		}
	}
	return nil
}

func DirectoryList(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	var dirList []string
	if err != nil {
		return dirList, err
	}

	for _, f := range files {
		if f.IsDir() {
			dirList = append(dirList, f.Name())
		}
	}

	return dirList, nil
}

func DirectoryTreeList(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		files = append(files, path)
		return nil
	})

	if err != nil {
		return files, err
	}

	return files, nil
}

func FileList(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	var fileList []string
	if err != nil {
		return fileList, err
	}

	for _, f := range files {
		if !f.IsDir() {
			fileList = append(fileList, f.Name())
		}
	}

	return fileList, nil
}

func PluginList(path, configName string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	var pluginList []string
	if err != nil {
		return pluginList, err
	}

	for _, f := range files {
		if configName+".conf" == f.Name() {
			continue
		}
		if !f.IsDir() {
			pluginList = append(pluginList, f.Name())
		}
	}

	return pluginList, nil
}

func BoolToString(b bool) string {
	s := strconv.FormatBool(b)

	return s
}

func IntToString(i int) string {
	s := strconv.Itoa(i)

	return s
}

func Int64ToString(i int64) string {
	str := strconv.FormatInt(i, 10)
	return str
}

func StringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Println("ERROR " + err.Error())
	}
	return i
}

func StringToInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Println(err.Error())
	}

	return i
}

func IntInSlice(i int, list []int) bool {
	for _, v := range list {
		if v == i {
			return true
		}
	}
	return false
}
