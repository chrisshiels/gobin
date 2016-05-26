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


func cat(reader io.Reader) error {
    bytes := make([]byte, 65536)

    var n int
    var err error
    for n, err = reader.Read(bytes); err == nil; n, err = reader.Read(bytes) {
        if _, err = os.Stdout.Write(bytes[0:n]); err != nil {
            break
        }
    }

    if err == io.EOF {
        err = nil
    }

    return err
}


func main() {
    var err error

    if len(os.Args[1:]) == 0 {
        err = cat(os.Stdin)
        _ = os.Stdin.Close()
        if err != nil {
            fmt.Fprintf(os.Stderr, "cat: %s\n", err)
            os.Exit(exitfailure)
        }
    } else {
        for _, filename := range os.Args[1:] {
            var file *os.File

            if file, err = os.Open(filename); err != nil {
                fmt.Fprintf(os.Stderr, "cat: %s\n", err)
                continue
            }

            err = cat(file)
            _ = file.Close()
            if err != nil {
                fmt.Fprintf(os.Stderr, "cat: %s\n", err)
                continue
            }
        }
    }

    os.Exit(exitsuccess)
}
