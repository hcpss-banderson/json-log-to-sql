package main

import "bufio"
import "fmt"
import "encoding/json"
import "io"
import "github.com/stoicperlman/fls"

func main() {
    inputFile, outputFile, numInputFileLines, numOutputFileLines := resolveArgs()
    defer inputFile.Close()
    defer outputFile.Close()

    lineNumber := numOutputFileLines - 1
    lineFile := fls.LineFile(inputFile)    
    lineFile.SeekLine(int64(lineNumber), io.SeekStart)
    scanner := bufio.NewScanner(lineFile)
    for scanner.Scan() {
        line := scanner.Text()[12:]
        lineNumber++
        
        log := Log{}
        if err := json.Unmarshal([]byte(line), &log); err != nil {
            panic(err)
        }
        
        lineEnder := ","
        if lineNumber == numInputFileLines {
            lineEnder = ";"
        }
        insert := fmt.Sprintf("%s%s\n", log.ToInsert(), lineEnder)
        
        outputFile.Write([]byte(insert))
    }
}
