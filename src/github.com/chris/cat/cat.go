// 'cat.go'.
// Chris Shiels.


package main


import (
    "fmt"
    "io"
    "os"
)


const exitsuccess = 0
const exitfailure = 1


func cat(file *os.File) error {

    var n int
    var err error
    bytes := make([]byte, 65536)

    for n, err = file.Read(bytes); err == nil; n, err = file.Read(bytes) {
        if _, err = os.Stdout.Write(bytes[0:n]); err != nil {
            break
        }
    }

    if (err != io.EOF) {
        return err
    }

    return nil
}


func main() {

    var err error

    if len(os.Args[1:]) == 0 {
        err = cat(os.Stdin)
        _ = os.Stdin.Close()
    } else {
        for _, filename := range os.Args[1:] {
            var file *os.File
            if file, err = os.Open(filename); err == nil {
                err = cat(file)
                _ = file.Close()
                if err != nil {
                    break
                }
            }
        }
    }

    if err != nil {
        fmt.Printf("Error:  %s\n", err)
        os.Exit(exitfailure)
    }

    os.Exit(exitsuccess)
}
