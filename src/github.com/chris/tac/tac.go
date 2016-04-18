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


func tac(reader io.Reader) error {
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
        fmt.Println(l[i])
    }

    return nil
}


func main() {
    var err error

    if len(os.Args[1:]) == 0 {
        err = tac(os.Stdin)
        _ = os.Stdin.Close()
    } else {
        for _, filename := range os.Args[1:] {
            var file *os.File
            if file, err = os.Open(filename); err == nil {
                err = tac(file)
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
