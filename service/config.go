package service

import (
    "encoding/json"
    "errors"
    "github.com/racecarparts/dashster/model"
    "io/ioutil"
    "os"
    "os/user"
)

var configFilename = ".dashster_config.json"

func ReadOrCreateConfig() error {
    userFolder, err := user.Current()
    if err != nil {
        return err
    }
    confFilePath := userFolder.HomeDir + string(os.PathSeparator) + configFilename

    confFile, err := os.OpenFile(confFilePath, os.O_CREATE, 0644)
    if err != nil {
        return err
    }
    defer confFile.Close()

    confData, err := ioutil.ReadAll(confFile)

    if len(confData) == 0 {
        model.AppConfig = &model.Config{}
        confData, err := json.MarshalIndent(model.AppConfig, "", "  ")
        if err != nil {
            return err
        }
        err = os.WriteFile(confFilePath, confData, 0644)
        if err != nil {
            return err
        }
        return errors.New("config file was not present, but has been created, finish filling it out at " + confFilePath + ".")
    }

    err = json.Unmarshal(confData, &model.AppConfig)
    if err != nil {
        return err
    }

    return nil
}