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


func cat(reader io.Reader, writer io.Writer) error {
    bytes := make([]byte, 65536)

    var n int
    var err error
    for n, err = reader.Read(bytes); err == nil; n, err = reader.Read(bytes) {
        if _, err = writer.Write(bytes[0:n]); err != nil {
            break
        }
    }

    if err == io.EOF {
        err = nil
    }

    return err
}


func _main(stdin io.Reader,
           stdout io.Writer,
           stderr io.Writer,
           args []string) (exitstatus int) {
    var err error

    if len(args[1:]) == 0 {
        err = cat(stdin, stdout)
        if err != nil {
            fmt.Fprintf(stderr, "cat: %s\n", err)
            return exitfailure
        }
    } else {
        for _, filename := range args[1:] {
            var file *os.File

            if file, err = os.Open(filename); err != nil {
                fmt.Fprintf(stderr, "cat: %s\n", err)
                continue
            }

            err = cat(file, stdout)
            _ = file.Close()
            if err != nil {
                fmt.Fprintf(stderr, "cat: %s\n", err)
                continue
            }
        }
    }

    return exitsuccess
}


func main() {
    os.Exit(_main(os.Stdin, os.Stdout, os.Stderr, os.Args))
}
