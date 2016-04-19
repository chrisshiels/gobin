// 'find.go'.
// Chris Shiels.


package main


import (
    "fmt"
    "os"
    "path/filepath"
)


const exitsuccess = 0
const exitfailure = 1


type predicate func(path string, info os.FileInfo) (matched bool, err error)


func makepredicatename(name string) (p predicate, err error) {
    return func(path string, info os.FileInfo) (matched bool, err error) {
        return filepath.Match(name, filepath.Base(path))
    }, nil
}


func makepredicateprint() (p predicate, err error) {
    return func(path string, info os.FileInfo) (matched bool, err error) {
        fmt.Println(path)
        return true, nil
    }, nil
}


func makepredicatetype(filetype string) (p predicate, err error) {
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

    return nil, fmt.Errorf("Unrecognised type %s", filetype)
}


type find struct {
    listpaths []string
    listpredicates []predicate
}


func (f *find) parse(args []string) error {
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
            p, err = makepredicatename(args[i + 1])
            i++
        } else if args[i] == "-print" {
            p, err = makepredicateprint()
        } else if args[i] == "-type" && i + 1 < len(args) {
            p, err = makepredicatetype(args[i + 1])
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


func (f *find) walkfunc() filepath.WalkFunc {
    return func(path string, info os.FileInfo, err error) error {
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error: % s\n", err)
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


func (f *find) walk() error {
    for _, path := range f.listpaths {
        err := filepath.Walk(path, f.walkfunc())
        if err != nil {
            return err
        }
    }

    return nil
}


func main() {

    var f find
    var err error

    if err = f.parse(os.Args[1:]); err != nil {
        fmt.Fprintf(os.Stderr, "Error:  %s\n", err)
        os.Exit(exitfailure)
    }

    if err = f.walk(); err != nil {
        fmt.Fprintf(os.Stderr, "Error:  %s\n", err)
        os.Exit(exitfailure)
    }

    os.Exit(exitsuccess)
}
