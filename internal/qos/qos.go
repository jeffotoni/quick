package qos

import "os"

func FileExist(path string) bool {
    _, err := os.Stat(path)
    if os.IsNotExist(err) {
        return false
    } else if err != nil {
        return false
    } else {
        return true
    }
}
