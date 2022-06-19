package io

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//Padding is 1 + 000000 bunch of zeroes
//Add 64 bits of the original message to the result of above ??
//Initialize MD Buffer, the computed of the 512 bits output

//Config type
type Config struct {
	Path string `json:"Hash Path"`
}

//Function that will return an array of 512 bits needed for MD5
//Path to the file is in the config
//If the file is a .txt, its contents will be hashed instead of the file
func GetFile() ([]byte, error) {
	//Get the config
	confF, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	confBody, err := ioutil.ReadAll(confF)
	if err != nil {
		return nil, err
	}
	var conf Config
	err = json.Unmarshal(confBody, &conf)
	if err != nil {
		return nil, err
	}

	//Get the hash file
	hashF, err := os.Open(conf.Path)
	if err != nil {
		return nil, err
	}

	hashData, err := ioutil.ReadAll(hashF)
	if err != nil {
		return nil, err
	}

	return hashData, nil
}
