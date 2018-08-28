package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type yamlFile struct {
	path  string
	isnew bool
	data  map[interface{}]interface{}
}

func NewYamlFile(filepath string) *yamlFile {
	yf := &yamlFile{
		path: filepath,
		data: make(map[interface{}]interface{}),
	}
	return yf
}

func (yf *yamlFile) ReadYamlFile() error {
	// Check existance of config file specified by the path.
	// - not exist: create new map data
	// - exist    : set map data from file
	if _, err := os.Stat(yf.path); os.IsNotExist(err) {
		yf.isnew = true
		yf.data["apiVersion"] = "v1"
		yf.data["kind"] = "Config"
		yf.data["preferences"] = map[interface{}]interface{}{}
		yf.data["clusters"] = make([]interface{}, 0)
		yf.data["users"] = make([]interface{}, 0)
		yf.data["contexts"] = make([]interface{}, 0)
		yf.data["current-context"] = ""
	} else {
		yf.isnew = false
		// Read from yaml file to buffer ([]byte).
		buf, err := ioutil.ReadFile(yf.path)
		if err != nil {
			return errors.New(fmt.Sprintf("cannot read yaml file\n(detail) %s", err.Error()))
		}
		// Unmarshal yaml data. (buffer -> map-struucture)
		if err := yaml.Unmarshal([]byte(buf), &yf.data); err != nil {
			return errors.New(fmt.Sprintf("cannot unmarshal yaml data\n(detail) %s", err.Error()))
		}
	}

	return nil
}

func (yf *yamlFile) WriteYamlFile() error {
	// Marshal yaml data. (map-struucture -> buffer)
	buf, err := yaml.Marshal(&yf.data)
	if err != nil {
		return errors.New(fmt.Sprintf("cannot marshal yaml data\n(detail) %s", err.Error()))
	}

	// Write to yaml file from buffer ([]byte).
	if err = ioutil.WriteFile(yf.path, buf, 0600); err != nil {
		return errors.New(fmt.Sprintf("cannot write yaml file\n(detail) %s", err.Error()))
	}

	return nil
}

func (yf *yamlFile) SetParamCluster(keyname string, endpoint string, cert string) error {
	// Create the new map and set values.
	subdata := make(map[interface{}]interface{})
	subdata["name"] = keyname
	subdata["cluster"] = make(map[interface{}]interface{})
	subdata["cluster"].(map[interface{}]interface{})["server"] = endpoint
	subdata["cluster"].(map[interface{}]interface{})["certificate-authority-data"] = cert

	// Check existance and format of the key "clusters".
	if _, ok := yf.data["clusters"]; !ok {
		return errors.New("illegal config format: key 'clusters' is not exist")
	}
	if _, ok := yf.data["clusters"].([]interface{}); !ok {
		return errors.New("illegal config format: key 'clusters' must have array data")
	}

	// Check existance of data that have same "name" in array(slice).
	// - exist    : overwrite it by new data
	// - not exist: append new data into array
	pos := searchInMapArray(yf.data["clusters"].([]interface{}), "name", keyname)
	if pos >= 0 {
		yf.data["clusters"].([]interface{})[pos] = subdata
	} else {
		yf.data["clusters"] = append(yf.data["clusters"].([]interface{}), subdata)
	}

	return nil
}

func (yf *yamlFile) SetParamUser(keyname string, clustername string) error {
	// Create the new map and set values.
	subdata := make(map[interface{}]interface{})
	subdata["name"] = keyname
	subdata["user"] = make(map[interface{}]interface{})
	subdata["user"].(map[interface{}]interface{})["exec"] = make(map[interface{}]interface{})
	subdata["user"].(map[interface{}]interface{})["exec"].(map[interface{}]interface{})["apiVersion"] = "client.authentication.k8s.io/v1alpha1"
	subdata["user"].(map[interface{}]interface{})["exec"].(map[interface{}]interface{})["command"] = "heptio-authenticator-aws"
	subdata["user"].(map[interface{}]interface{})["exec"].(map[interface{}]interface{})["args"] = make([]interface{}, 3)
	subdata["user"].(map[interface{}]interface{})["exec"].(map[interface{}]interface{})["args"].([]interface{})[0] = "token"
	subdata["user"].(map[interface{}]interface{})["exec"].(map[interface{}]interface{})["args"].([]interface{})[1] = "-i"
	subdata["user"].(map[interface{}]interface{})["exec"].(map[interface{}]interface{})["args"].([]interface{})[2] = clustername

	// Check existance and format of the key "users".
	if _, ok := yf.data["users"]; !ok {
		return errors.New("illegal config format: key 'users' is not exist")
	}
	if _, ok := yf.data["users"].([]interface{}); !ok {
		return errors.New("illegal config format: key 'users' must have array data")
	}

	// Check existance of data that have same "name" in array(slice).
	// - exist    : overwrite it by new data
	// - not exist: append new data into array
	pos := searchInMapArray(yf.data["users"].([]interface{}), "name", keyname)
	if pos >= 0 {
		yf.data["users"].([]interface{})[pos] = subdata
	} else {
		yf.data["users"] = append(yf.data["users"].([]interface{}), subdata)
	}

	return nil
}

func (yf *yamlFile) SetParamContext(keyname string, cluster string, user string) error {
	// Create the new map and set values.
	subdata := make(map[interface{}]interface{})
	subdata["name"] = keyname
	subdata["context"] = make(map[interface{}]interface{})
	subdata["context"].(map[interface{}]interface{})["cluster"] = cluster
	subdata["context"].(map[interface{}]interface{})["user"] = user

	// Check existance and format of the key "contexts".
	if _, ok := yf.data["contexts"]; !ok {
		return errors.New("illegal config format: key 'contexts' is not exist")
	}
	if _, ok := yf.data["contexts"].([]interface{}); !ok {
		return errors.New("illegal config format: key 'contexts' must have array data")
	}

	// Check existance of data that have same "name" in array(slice).
	// - exist    : overwrite it by new data
	// - not exist: append new data into array
	pos := searchInMapArray(yf.data["contexts"].([]interface{}), "name", keyname)
	if pos >= 0 {
		yf.data["contexts"].([]interface{})[pos] = subdata
	} else {
		yf.data["contexts"] = append(yf.data["contexts"].([]interface{}), subdata)
	}

	return nil
}

func (yf *yamlFile) SetParamCurrentContext(context string) error {
	// Overwrite the value of key "current-context".
	yf.data["current-context"] = context

	return nil
}

func searchInMapArray(array []interface{}, key string, value string) int {
	for i := 0; i < len(array); i++ {
		if array[i].(map[interface{}]interface{})[key] == value {
			return i
		}
	}
	return -1
}
