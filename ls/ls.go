// 'ls.go'.
// Chris Shiels.


package main


import (
    "fmt"
    "io"
    "io/ioutil"
    "os"
)


const exitsuccess = 0
const exitfailure = 1


func ls(stdout io.Writer, dirname string) error {
    fileinfos, err := ioutil.ReadDir(dirname)
    if err != nil {
        return err
    }

    for _, fileinfo := range fileinfos {
        var symbol string = ""
        if fileinfo.Mode().IsDir() {
            symbol = "/"
        }
        fmt.Fprintf(stdout, "%s%s\n", fileinfo.Name(), symbol)
    }

    return nil
}


func _main(stdin io.Reader,
           stdout io.Writer,
           stderr io.Writer,
           args []string) (exitstatus int) {
    var err error

    if len(args[1:]) == 0 {
        err = ls(stdout, ".")
    } else if len(args[1:]) == 1 {
        err = ls(stdout, args[1])
    } else {
        for n, arg := range args[1:] {
            if n > 0 {
                fmt.Fprintln(stdout)
            }
            fmt.Fprintf(stdout, "%s:\n", arg)

            err = ls(stdout, arg)
            if err != nil {
                break
            }
        }
    }

    if err != nil {
        fmt.Fprintf(stderr, "ls: %s\n", err)
        return exitfailure
    }

    return exitsuccess
}


func main() {
    os.Exit(_main(os.Stdin, os.Stdout, os.Stderr, os.Args))
}
