package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
)

var dataFile = fmt.Sprintf("%v/vct-data.json", getWorkingDir())

func LoadData() (*Data, error) {
	dat, err := ioutil.ReadFile(dataFile)
	if err != nil {
		return nil, err
	}
	data := Data{}

	err = json.Unmarshal(dat, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func SaveData(data Data) {
	dat, err := json.MarshalIndent(data, "", "  ")
	check(err)

	file, err := os.Create(dataFile)
	check(err)
	defer file.Close()

	_, err = file.Write(dat)
	check(err)
	err = file.Sync()
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getHomeDir() string {
	usr, err := user.Current()
	check(err)
	return usr.HomeDir
}

func getWorkingDir() string {
	dir, err := os.Getwd()
	check(err)
	return dir
}