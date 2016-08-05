// 'find.go'.
// Chris Shiels.


package main


import (
    "fmt"
    "io"
    "os"
    "path/filepath"
)


const exitsuccess = 0
const exitfailure = 1


type predicate func(path string, info os.FileInfo) (matched bool, err error)


func makepredicatename(stdout io.Writer,
                       stderr io.Writer,
                       name string) (p predicate, err error) {
    return func(path string, info os.FileInfo) (matched bool, err error) {
        return filepath.Match(name, filepath.Base(path))
    }, nil
}


func makepredicateprint(stdout io.Writer,
                        stderr io.Writer) (p predicate, err error) {
    return func(path string, info os.FileInfo) (matched bool, err error) {
        fmt.Fprintln(stdout, path)
        return true, nil
    }, nil
}


func makepredicatetype(stdout io.Writer,
                       stderr io.Writer,
                       filetype string) (p predicate, err error) {
    if filetype == "d" {
        return func(path string, info os.FileInfo) (matched bool, err error) {
            return info.Mode().IsDir(), nil
        }, nil
    }

    if filetype == "f" {
        return func(path string, info os.FileInfo) (matched bool, err error) {
            return info.Mode().IsRegular(), nil
        }, nil
    }

    if filetype == "l" {
        return func(path string, info os.FileInfo) (matched bool, err error) {
            return info.Mode() & os.ModeSymlink == os.ModeSymlink, nil
        }, nil
    }

    return nil, fmt.Errorf("Unrecognised type %s", filetype)
}


type find struct {
    listpaths []string
    listpredicates []predicate
}


func (f *find) parse(stdout io.Writer,
                     stderr io.Writer,
                     args []string) error {
    var i int = 0

    for i < len(args) {
        if args[i][0:1] == "-" {
            break
        }
        f.listpaths = append(f.listpaths, args[i])
        i++
    }

    for i < len(args) {
        var p predicate
        var err error

        if args[i] == "-name" && i + 1 < len(args) {
            p, err = makepredicatename(stdout, stderr, args[i + 1])
            i++
        } else if args[i] == "-print" {
            p, err = makepredicateprint(stdout, stderr)
        } else if args[i] == "-type" && i + 1 < len(args) {
            p, err = makepredicatetype(stdout, stderr, args[i + 1])
            i++
        }

        if err != nil {
            return err
        }

        if p == nil {
            return fmt.Errorf("Unrecognised predicate %s", args[i])
        }

        f.listpredicates = append(f.listpredicates, p)
        i++
    }

    return nil
}


func (f *find) walkfunc(stdout io.Writer,
                        stderr io.Writer) filepath.WalkFunc {
    return func(path string, info os.FileInfo, err error) error {
        if err != nil {
            fmt.Fprintf(stderr, "find: % s\n", err)
        } else {
            for _, predicate := range f.listpredicates {
                var matched bool
                matched, err = predicate(path, info)

                if err != nil {
                    return err
                }

                if !matched {
                    break
                }
            }
        }

        return nil
    }
}


func (f *find) walk(stdout io.Writer,
                    stderr io.Writer) error {
    for _, path := range f.listpaths {
        err := filepath.Walk(path, f.walkfunc(stdout, stderr))
        if err != nil {
            return err
        }
    }

    return nil
}


func _main(stdin io.Reader,
           stdout io.Writer,
           stderr io.Writer,
           args []string) (exitstatus int) {
    var f find
    var err error

    if err = f.parse(stdout, stderr, args[1:]); err != nil {
        fmt.Fprintf(stderr, "find: %s\n", err)
        return exitfailure
    }

    if err = f.walk(stdout, stderr); err != nil {
        fmt.Fprintf(stderr, "find: %s\n", err)
        return exitfailure
    }

    return exitsuccess
}


func main() {
    os.Exit(_main(os.Stdin, os.Stdout, os.Stderr, os.Args))
}
