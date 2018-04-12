package main

import "io"
import "bytes"
import "strings"
import "os"
import "fmt"

func resolveArgs() (*os.File, *os.File, int, int) {
    inf, _ := os.Open(os.Args[1])
    outf, isNew := openOrCreate(os.Args[2])
    nif, _ := lineCounter(inf)
    nof, _ := lineCounter(outf)
    
    inf.Seek(0, 0)
    outf.Seek(0, 0)
    
    if isNew {
        outf.Write(getInsert())
    }
    
    return inf, outf, nif, nof
}

func getInsert() ([]byte) {
    format := "INSERT INTO `mdl_logstore_standard_log` (`%s`) VALUES\n"
    keys := Log{}.JsonKeys()
    insert := fmt.Sprintf(format, strings.Join(keys, "`,`"))
    
    return []byte(insert)
}

func lineCounter(r io.Reader) (int, error) {
    buf := make([]byte, 32*1024)
    count := 0
    lineSep := []byte{'\n'}
    
    for {
        c, err := r.Read(buf)
        count += bytes.Count(buf[:c], lineSep)

        switch {
        case err == io.EOF:
            return count, nil

        case err != nil:
            return count, err
        }
    }
}

func fileExists(filename string) bool {
    _, err := os.Stat(filename)
    return os.IsNotExist(err)
}

func openOrCreate(filename string) (*os.File, bool) {
    isNew := fileExists(filename)
    outputFile, _ := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
    
    return outputFile, isNew
}
