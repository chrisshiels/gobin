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


func grep(pattern string, reader io.Reader) error {
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
            fmt.Println(scanner.Text())
        }
    }

    if err = scanner.Err(); err != nil {
        return err
    }

    return nil
}


func main() {
    var err error

    if len(os.Args[1:]) == 0 {
        fmt.Fprintln(os.Stderr, "Usage:  grep pattern [ file ... ]")
        os.Exit(exitsuccess)
    } else if len(os.Args[1:]) == 1 {
        err = grep(os.Args[1], os.Stdin)
        _ = os.Stdin.Close()
        if err != nil {
            fmt.Fprintf(os.Stderr, "grep: %s\n", err)
            os.Exit(exitfailure)
        }
    } else {
        for _, filename := range os.Args[2:] {
            var file *os.File

            if file, err = os.Open(filename); err != nil {
                fmt.Fprintf(os.Stderr, "grep: %s\n", err)
                continue
            }

            err = grep(os.Args[1], file)
            _ = file.Close()
            if err != nil {
                fmt.Fprintf(os.Stderr, "grep: %s\n", err)
                continue
            }
        }
    }

    os.Exit(exitsuccess)
}
