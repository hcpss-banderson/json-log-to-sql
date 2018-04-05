package main

import (
    "bytes"
    "fmt"
    "strconv"
    "math"
)

const TYPE_VALUE_SEPARATOR = ':'
const VALUES_SEPARATOR = ';'

func PhpSerialize(value interface{}) (result string, err error) {
    buf := new(bytes.Buffer)
    err = phpEncodeValue(buf, value)
    if err == nil {
        result = buf.String()
    } else {
        panic(err)
    }
    
    return
}

func phpEncodeValue(buf *bytes.Buffer, value interface{}) (err error) {    
    switch t := value.(type) {
    default:
        err = fmt.Errorf("Unexpected type %T", t)
    case []interface{}:
        buf.WriteString("a")
        buf.WriteRune(TYPE_VALUE_SEPARATOR)
        err = phpEncodeSimpleArrayCore(buf, t)
    case *interface{}:
        myInterface := *t
        myMap := myInterface.(map[string]interface{})
        
        err = phpEncodeValue(buf, myMap)
        
        // myMapInterface := make(map[interface{}]interface{})
        // for key, value := range myMap {
        //     myMapInterface[key] = value
        // }
        // err = phpEncodeValue(buf, myMapInterface)
    case map[string]interface{}:
        myMapInterface := make(map[interface{}]interface{})
        for key, value := range t {
            myMapInterface[key] = value
        }
        err = phpEncodeValue(buf, myMapInterface)
    case bool:
        buf.WriteString("b")
        buf.WriteRune(TYPE_VALUE_SEPARATOR)
        if t {
            buf.WriteString("1")
        } else {
            buf.WriteString("0")
        }
        buf.WriteRune(VALUES_SEPARATOR)
    case nil:
        buf.WriteString("N")
        buf.WriteRune(VALUES_SEPARATOR)
    case int, int64, int32, int16, int8:
        buf.WriteString("i")
        buf.WriteRune(TYPE_VALUE_SEPARATOR)
        strValue := fmt.Sprintf("%v", t)
        buf.WriteString(strValue)
        buf.WriteRune(VALUES_SEPARATOR)
    case float32:
        buf.WriteString("d")
        buf.WriteRune(TYPE_VALUE_SEPARATOR)
        strValue := strconv.FormatFloat(float64(t), 'f', -1, 64)
        buf.WriteString(strValue)
        buf.WriteRune(VALUES_SEPARATOR)
    case float64:
        // We seem to get a lot of integers here.
        remainder := math.Mod(t, 1)
        
        if remainder == 0 {
            // This is an intiger.
            err = phpEncodeValue(buf, int(t))  
        } else {
            // OK, this really is a float64.
            buf.WriteString("d")
            buf.WriteRune(TYPE_VALUE_SEPARATOR)
            strValue := strconv.FormatFloat(float64(t), 'f', -1, 64)
            buf.WriteString(strValue)
            buf.WriteRune(VALUES_SEPARATOR)
        }
    case string:
        buf.WriteString("s")
        buf.WriteRune(TYPE_VALUE_SEPARATOR)
        phpEncodeString(buf, t)
        buf.WriteRune(VALUES_SEPARATOR)
    case map[interface{}]interface{}:
        buf.WriteString("a")
        buf.WriteRune(TYPE_VALUE_SEPARATOR)
        err = phpEncodeArrayCore(buf, t)
    }
    return
}

func phpEncodeString(buf *bytes.Buffer, strValue string) {
    valLen := strconv.Itoa(len(strValue))
    buf.WriteString(valLen)
    buf.WriteRune(TYPE_VALUE_SEPARATOR)
    buf.WriteRune('"')
    buf.WriteString(strValue)
    buf.WriteRune('"')
}

func phpEncodeSimpleArrayCore(buf *bytes.Buffer, arrValue []interface{}) (err error) {
    valLen := strconv.Itoa(len(arrValue))
    buf.WriteString(valLen)
    buf.WriteRune(TYPE_VALUE_SEPARATOR)
    
    buf.WriteRune('{')
    for v := range arrValue {
        if err = phpEncodeValue(buf, v); err != nil {
            break
        }
        
        // if intKey, _err := strconv.Atoi(fmt.Sprintf("%v", k)); _err == nil {
        //     if err = phpEncodeValue(buf, intKey); err != nil {
        //         break
        //     }
        // } else {
        //     if err = phpEncodeValue(buf, k); err != nil {
        //         break
        //     }
        // }
        // if err = phpEncodeValue(buf, v); err != nil {
        //     break
        // }
    }
    buf.WriteRune('}')
    return err
}

func phpEncodeArrayCore(buf *bytes.Buffer, arrValue map[interface{}]interface{}) (err error) {
    valLen := strconv.Itoa(len(arrValue))
    buf.WriteString(valLen)
    buf.WriteRune(TYPE_VALUE_SEPARATOR)

    buf.WriteRune('{')
    for k, v := range arrValue {
        if intKey, _err := strconv.Atoi(fmt.Sprintf("%v", k)); _err == nil {
            if err = phpEncodeValue(buf, intKey); err != nil {
                break
            }
        } else {
            if err = phpEncodeValue(buf, k); err != nil {
                break
            }
        }
        if err = phpEncodeValue(buf, v); err != nil {
            break
        }
    }
    buf.WriteRune('}')
    return err
}
