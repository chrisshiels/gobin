// 'head.go'.
// Chris Shiels.


package main


import (
    "bufio"
    "flag"
    "fmt"
    "io"
    "os"
)


const exitsuccess = 0
const exitfailure = 1


type header struct {
    filename string
    reader *bufio.Reader
}


func newheader(filename string, file *os.File) *header {
    return &header{ filename: filename,
                    reader: bufio.NewReader(file) }
}


func (t *header) head(nlines int) error {
    linen := 0

    var s string
    var err error
    for s, err = t.reader.ReadString('\n');
        err == nil;
        s, err = t.reader.ReadString('\n') {
        fmt.Println(s[0:len(s) - 1])
        linen++
        if linen == nlines {
            break
        }
    }

    if err != io.EOF {
        return err
    }
    return nil
}


func main() {
    flag.Usage = func() {
        fmt.Println("Usage:  head [ -n n ] [ file ... ]")
        flag.PrintDefaults()
    }
    flagn := flag.Int("n", 10, "Output the last n lines")

    // Note flag.Parse() will also handle '-h' and '--help' and will exit with
    // exit status 2.
    flag.Parse()


    if len(flag.Args()) == 0 {
        header := newheader("-", os.Stdin)
        header.head(*flagn)
    } else {
        for i, filename := range flag.Args() {
            file, err := os.Open(filename)
            if err != nil {
                fmt.Fprintf(os.Stderr, "head: %s\n", err)
                continue
            }

            if len(flag.Args()) > 1 {
                if i != 0 {
                    fmt.Println()
                }
                fmt.Printf("==> %s <==\n", filename)
            }

            header := newheader(filename, file)
            err = header.head(*flagn)
            if err != nil {
                fmt.Fprintf(os.Stderr, "head: %s\n", err)
                continue
            }
        }
    }

    os.Exit(exitsuccess)
}
