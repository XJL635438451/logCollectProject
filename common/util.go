package common

import (
    "errors"
)

//Standardized error output
func ErrMsg(errList ...interface{}) (err error) {
    var errRes string
    for k, val := range errList {
        if k == 1 {
            errRes += " Error: "
        }
        switch v := val.(type) {
        case string:
            errRes += v
        case error:
            errRes += v.Error()
        }
    }
    err = errors.New(errRes)
    return
}
