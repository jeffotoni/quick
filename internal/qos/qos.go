package qos

import (
    "errors"
    "os"
)

func FileExist(path string) error {
    _, err := os.Stat(path)
    if os.IsNotExist(err) {
        return errors.New("error file not exist")
    } else if err != nil {
        return err
    } else {
        return nil
    }
}
