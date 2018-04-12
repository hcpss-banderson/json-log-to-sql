package main

import "fmt"
import "strings"
import "reflect"
import "strconv"

type Log struct {
    EventName         *string      `json:"eventname"`
    Component         *string      `json:"component"`
    Action            *string      `json:"action"`
    Target            *string      `json:"target"`
    ObjectTable       *string      `json:"objecttable"`
    ObjectId          *int         `json:"objectid"`
    Crud              *string      `json:"crud"`
    EduLevel          *int         `json:"edulevel"`
    ContextId         *int         `json:"contextid"`
    ContextLevel      *int         `json:"contextlevel"`
    ContextInstanceId *int         `json:"contextinstanceid"`
    UserId            *int         `json:"userid"`
    CourseId          *int         `json:"courseid"`
    RelatedUserId     *int         `json:"relateduserid"`
    Anonymous         *int         `json:"anonymous"`
    Other             *interface{} `json:"other"`
    TimeCreated       *int         `json:"timecreated"`
    Origin            *string      `json:"origin"`
    Ip                *string      `json:"ip"`
    RealUserId        *int         `json:"realuserid"`
}

func (l Log) JsonKeys() []string {
    log := reflect.TypeOf(l)
    value := reflect.ValueOf(l)
    labels := make([]string, value.NumField())
    
    for i := 0; i < value.NumField(); i++ {
        label, _ := log.Field(i).Tag.Lookup("json")
        labels[i] = label
    }
    
    return labels
}

func (l Log) ToInsert() string {
    v := reflect.ValueOf(l)
    values := make([]string, v.NumField())

    for i := 0; i < v.NumField(); i++ {
        value := v.Field(i).Interface()
        
        switch value.(type) {
        case *int:
            intVal := value.(*int)
            if intVal == nil {
                values[i] = "NULL"
            } else {
                values[i] = strconv.Itoa(*intVal)
            }
        case *string:
            stringVal := value.(*string)
            if stringVal == nil {
                values[i] = "NULL"
            } else {
                values[i] = fmt.Sprintf("'%s'", *stringVal)
            }
        default:
            interfaceValue := value.(*interface{})            
            if interfaceValue == nil {
                values[i] = "NULL"
            } else {
                other, _ := PhpSerialize(value)
                values[i] = fmt.Sprintf("'%s'", strings.Replace(other, "'", "''", -1))
            }
        }
    }
    
    return fmt.Sprintf("%s%s%s", "(", strings.Join(values, ","), ")")
}
