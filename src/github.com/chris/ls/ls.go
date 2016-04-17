// 'ls.go'.
// Chris Shiels.


package main


import (
    "fmt"
    "io/ioutil"
    "os"
)


const exitsuccess = 0
const exitfailure = 1


func ls(dirname string) error {

    fileinfos, err := ioutil.ReadDir(dirname)

    if err != nil {
        return err
    }

    for _, fileinfo := range fileinfos {

        var symbol string = ""
        if fileinfo.Mode().IsDir() {
            symbol = "/"
        }
        fmt.Printf("%s%s\n", fileinfo.Name(), symbol)
    }

    return nil
}


func main() {

    var err error

    if len(os.Args[1:]) == 0 {
        err = ls(".")
    } else if len(os.Args[1:]) == 1 {
        err = ls(os.Args[1])
    } else {
        for n, arg := range os.Args[1:] {
            if n > 0 {
                fmt.Println()
            }
            fmt.Printf("%s:\n", arg)

            err = ls(arg)
            if err != nil {
                break
            }
        }
    }

    if err != nil {
        fmt.Printf("Error:  %s\n", err)
        os.Exit(exitfailure)
    }

    os.Exit(exitsuccess)
}
