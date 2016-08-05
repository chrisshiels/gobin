// 'grep.go'.
// Chris Shiels.


package main


import (
    "bufio"
    "fmt"
    "io"
    "os"
    "regexp"
)


const exitsuccess = 0
const exitfailure = 1


func grep(pattern string, reader io.Reader, writer io.Writer) error {
    var scanner *bufio.Scanner
    var err error

    scanner = bufio.NewScanner(reader)
    scanner.Split(bufio.ScanLines)

    for scanner.Scan() {
        var matched bool
        matched, err = regexp.MatchString(pattern, scanner.Text())
        if err != nil {
            return err
        }
        if matched {
            fmt.Fprintln(writer, scanner.Text())
        }
    }

    if err = scanner.Err(); err != nil {
        return err
    }

    return nil
}


func _main(stdin io.Reader,
           stdout io.Writer,
           stderr io.Writer,
           args []string) (exitstatus int) {
    var err error

    if len(args[1:]) == 0 {
        fmt.Fprintln(stderr, "Usage:  grep pattern [ file ... ]")
        return exitsuccess
    } else if len(args[1:]) == 1 {
        err = grep(args[1], stdin, stdout)
        if err != nil {
            fmt.Fprintf(stderr, "grep: %s\n", err)
            return exitfailure
        }
    } else {
        for _, filename := range args[2:] {
            var file *os.File

            if file, err = os.Open(filename); err != nil {
                fmt.Fprintf(stderr, "grep: %s\n", err)
                continue
            }

            err = grep(args[1], file, stdout)
            _ = file.Close()
            if err != nil {
                fmt.Fprintf(stderr, "grep: %s\n", err)
                continue
            }
        }
    }

    return exitsuccess
}


func main() {
    os.Exit(_main(os.Stdin, os.Stdout, os.Stderr, os.Args))
}
