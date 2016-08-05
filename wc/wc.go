// 'wc.go'.
// Chris Shiels.


package main


import (
    "bufio"
    "fmt"
    "io"
    "os"
    "unicode"
)


const exitsuccess = 0
const exitfailure = 1


func wc(reader io.Reader) (nl int, nw int, nm int, nc int, err error) {
    type wordstate int

    const (
        outsideword wordstate = iota
        insideword
    )

    var state wordstate = outsideword
    var bufioreader *bufio.Reader = bufio.NewReader(reader)
    var r rune
    var size int

    for r, size, err = bufioreader.ReadRune()
        err == nil
        r, size, err = bufioreader.ReadRune() {

        nc += size
        nm += 1

        if r == '\n' {
            nl += 1
        }

        if unicode.IsSpace(r) {
            state = outsideword
        } else if state == outsideword {
            state = insideword
            nw += 1
        }
    }

    if err == io.EOF {
        err = nil
    }

    return nl, nw, nm, nc, err
}


func _main(stdin io.Reader,
           stdout io.Writer,
           stderr io.Writer,
           args []string) (exitstatus int) {
    var nl, nw, nm, nc int
    var tl, tw, tm, tc int
    var err error

    if len(args[1:]) == 0 {
        if nl, nw, nm, nc, err = wc(stdin); err == nil {
            fmt.Fprintf(stdout, "%8d%8d%8d%8d\n", nl, nw, nm, nc)
        } else {
            fmt.Fprintf(stderr, "wc: %s\n", err)
            return exitfailure
        }
    } else {
        for _, filename := range args[1:] {
            var file *os.File

            if file, err = os.Open(filename); err != nil {
                fmt.Fprintf(stderr, "wc: %s\n", err)
                continue
            }

            nl, nw, nm, nc, err = wc(file)
            _ = file.Close()
            if err != nil {
                fmt.Fprintf(stderr, "wc: %s\n", err)
                continue
            }

            fmt.Fprintf(stdout, "%8d%8d%8d%8d\t%s\n", nl, nw, nm, nc, filename)
            tl += nl
            tw += nw
            tm += nm
            tc += nc
        }
        fmt.Fprintf(stdout, "%8d%8d%8d%8d\ttotal\n", tl, tw, tm, tc)
    }

    return exitsuccess
}


func main() {
    os.Exit(_main(os.Stdin, os.Stdout, os.Stderr, os.Args))
}
