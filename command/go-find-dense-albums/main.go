package main

import (
    "fmt"
    "os"
    "path"
    "sort"

    "encoding/csv"
    "path/filepath"

    "github.com/dsoprea/go-logging"
    "github.com/gonum/stat"
    "github.com/jessevdk/go-flags"
)

type parameters struct {
    Path      string   `long:"path" description:"Path to scan" required:"true"`
    Filespecs []string `long:"filespec" description:"Filespec to match for (can be provided zero or more times)"`
    Sigma     int      `long:"sigma" description:"By default, we look for counts larger than one standard deviation in the positive direction. This uses a different sigma value." default:"1"`
    JustSmall bool     `long:"small" description:"Just show small folders."`
    JustBig   bool     `long:"big" description:"Just show big folders."`
    Verbose   bool     `long:"verbose" description:"Print extra information."`
    PrintAll  bool     `long:"print-all" description:"Print all directories found first."`
}

var (
    arguments = new(parameters)
    wf        filepath.WalkFunc
)

func main() {
    p := flags.NewParser(arguments, flags.Default)

    _, err := p.Parse()
    if err != nil {
        os.Exit(1)
    }

    bins := make(map[string]int)

    wf = func(fullPath string, info os.FileInfo, err error) error {
        defer func() {
            if state := recover(); state != nil {
                err = log.Wrap(state.(error))
            }
        }()

        if info.IsDir() == true {
            return nil
        }

        parentPath := path.Dir(fullPath)
        filename := path.Base(fullPath)

        hit := false
        if len(arguments.Filespecs) > 0 {
            for _, filespec := range arguments.Filespecs {
                matched, err := filepath.Match(filespec, filename)
                log.PanicIf(err)

                if matched == true {
                    hit = true
                    break
                }
            }
        } else {
            hit = true
        }

        if hit == true {
            if _, found := bins[parentPath]; found == true {
                bins[parentPath]++
            } else {
                bins[parentPath] = 1
            }
        }

        return nil
    }

    err = filepath.Walk(arguments.Path, wf)
    log.PanicIf(err)

    // Calculate the standard deviation.

    values := make([]float64, len(bins))
    i := 0
    for _, count := range bins {
        values[i] = float64(count)
        i++
    }

    paths := make(sort.StringSlice, len(bins))
    j := 0
    for path, _ := range bins {
        paths[j] = path
        j++
    }

    paths.Sort()

    if arguments.PrintAll == true {
        for _, binPath := range paths {
            count := bins[binPath]
            fmt.Printf("%s: (%d)\n", binPath, count)
        }

        fmt.Printf("\n")
    }

    mu := stat.Mean(values, nil)
    sigma := stat.StdDev(values, nil)

    if arguments.Verbose == true {
        fmt.Printf("Mean: %.2f\n", mu)
        fmt.Printf("Standard deviation: %.2f\n", sigma)
    }

    threshold := int(mu + sigma*float64(arguments.Sigma))

    if arguments.Verbose == true {
        fmt.Printf("Threshold: %d\n", threshold)
    }

    fmt.Printf("\n")

    w := csv.NewWriter(os.Stdout)
    for _, binPath := range paths {
        count := bins[binPath]
        var record []string
        if arguments.JustSmall == true && count < threshold {
            record = make([]string, 2)
            record[0] = binPath
            record[1] = fmt.Sprintf("%d", count)
        } else if arguments.JustBig == true && count >= threshold {
            record = make([]string, 2)
            record[0] = binPath
            record[1] = fmt.Sprintf("%d", count)
        } else if arguments.JustSmall == false && arguments.JustBig == false {
            var sizeBracket int
            if count < threshold {
                sizeBracket = -1
            } else if count > threshold {
                sizeBracket = 1
            }

            record = make([]string, 2)
            record[0] = fmt.Sprintf("%d", sizeBracket)
            record[1] = binPath
            record[2] = fmt.Sprintf("%d", count)
        }

        if record != nil {
            err = w.Write(record)
            log.PanicIf(err)
        }
    }

    w.Flush()

    err = w.Error()
    log.PanicIf(err)
}
