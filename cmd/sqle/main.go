package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lazada/sqle"
	"github.com/lazada/sqle/internal"
)

func main() {
	flag.Usage = usage
	flag.Parse()

	var (
		i     int
		f     string
		fd    *os.File
		fi    *internal.FileInfo
		err   error
		ok    bool
		files = make(map[string]struct{})
		nam   = sqle.NewCachedConvention(new(sqle.SnakeConvention))
	)
	for _, f = range flag.Args() {
		if _, ok = files[f]; !ok {
			err = filepath.Walk(f, func(path string, info os.FileInfo, err error) error {
				if err == nil && !info.IsDir() && strings.HasSuffix(path, `.go`) && !strings.Contains(path, suffix) {
					files[f] = struct{}{}
				}
				return err
			})
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	if f = os.Getenv(`GOFILE`); f != `` {
		files[f] = struct{}{}
	}
	for f, _ = range files {
		if fi, err = internal.ParseFile(f, *tagFlag, nam); err != nil {
			log.Fatal(err)
		}
		for i = len(f) - 1; i >= 0 && !os.IsPathSeparator(f[i]); i-- {
			if f[i] == '.' {
				f = f[:i] + suffix + f[i:]
				break
			}
		}
		if fd, err = os.OpenFile(f, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644); err != nil {
			log.Fatal(err)
		}
		if err = internal.FileTemplate.Execute(fd, fi); err == nil {
			for _, si := range fi.Struct {
				if err = internal.StructTemplate.Execute(fd, si); err != nil {
					log.Fatal(err)
				}
			}
		} else {
			log.Fatal(err)
		}
		if err = fd.Close(); err != nil {
			log.Fatal(err)
		}
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `Usage:`)
	fmt.Fprintln(os.Stderr, `    sqle [Options] PATH ...`)
	fmt.Fprintln(os.Stderr, `    go generate [Options] PATH ...`)
	fmt.Fprintln(os.Stderr, `Options:`)
	flag.PrintDefaults()
}

var (
	tagFlag = flag.String(`tag`, `sql`, ``)
)

const (
	suffix = `.sqle`
)
