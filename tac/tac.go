// 'tac.go'.
// Chris Shiels.


package main


import (
    "bufio"
    "fmt"
    "io"
    "os"
)


const exitsuccess = 0
const exitfailure = 1


func tac(reader io.Reader, writer io.Writer) error {
    var scanner *bufio.Scanner
    var l []string
    var err error

    scanner = bufio.NewScanner(reader)
    scanner.Split(bufio.ScanLines)

    for scanner.Scan() {
        l = append(l, scanner.Text())
        // Note:
        // It looks like append doubles cap when len reaches cap.
        //fmt.Printf("Info:  len: %d, cap: %d\n", len(l), cap(l))
    }

    if err = scanner.Err(); err != nil {
        return err
    }

    for i := len(l) - 1; i >= 0; i-- {
        fmt.Fprintln(writer, l[i])
    }

    return nil
}


func _main(stdin io.Reader,
           stdout io.Writer,
           stderr io.Writer,
           args []string) (exitstatus int) {
    var err error

    if len(args[1:]) == 0 {
        err = tac(stdin, stdout)
        if err != nil {
            fmt.Fprintf(stderr, "tac: %s\n", err)
            return exitfailure
        }
    } else {
        for _, filename := range args[1:] {
            var file *os.File

            if file, err = os.Open(filename); err != nil {
                fmt.Fprintf(stderr, "tac: %s\n", err)
                continue
            }

            err = tac(file, stdout)
            _ = file.Close()
            if err != nil {
                fmt.Fprintf(stderr, "tac: %s\n", err)
                continue
            }
        }
    }

    return exitsuccess
}


func main() {
    os.Exit(_main(os.Stdin, os.Stdout, os.Stderr, os.Args))
}
