package AceSerializer

import (
    "fmt"
    "reflect"
    "regexp"
    "testing"
)

// Array
var array_str = []interface{}{"a", "b", "c", 1, 2, 3, 1.1234567890123321, 2.456, true, nil, 2}
var map_list = map[interface{}]interface{}{
    0:   "a",
    1:   "b",
    2:   "c",
    "a": 1,
    "b": 2,
    "c": 3,
    "e": [3]string{"1", "2", "3"},
}

func TestSerialize(t *testing.T) {
    r, err := Serialize(array_str)
    if err != nil {
        t.Fatal(err)
    }
    if r != "^1^T^N1^Sa^N2^Sb^N3^Sc^N4^N1^N5^N2^N6^N3^N7^F5059599576362793^f-52^N8^N2.456^N9^B^N11^N2^t^^" {
        t.Fatal("Expecting array.")
    }

    r, err = Serialize(array_str...)
    if err != nil {
        t.Fatal(err)
    }
    if r != "^1^Sa^Sb^Sc^N1^N2^N3^F5059599576362793^f-52^N2.456^B^Z^N2^^" {
        t.Fatal("Expecting multi parameter.")
    }

    r, err = Serialize(map_list)
    if err != nil {
        t.Fatal(err)
    }
    for k, v := range map_list {
        switch k.(type) {
        case int:
            f := fmt.Sprintf("\\^N%d\\^S%s", k.(int)+1, v)
            m, _ := regexp.MatchString(f, r)
            if !m {
                t.Fatal("Expecting map.")
                break
            }
        case string:
            var f string
            if reflect.ValueOf(v).Kind() == reflect.Array {
                f = fmt.Sprintf("\\^S%s\\^T\\^N1\\^S1\\^N2\\^S2\\^N3\\^S3\\^t", k)
            } else {
                f = fmt.Sprintf("\\^S%s\\^N%d", k, v)
            }
            m, _ := regexp.MatchString(f, r)
            if !m {
                t.Fatal("Expecting map.")
                break
            }
        default:
            t.Fatal("Expecting map.")
        }
    }
}

func TestDeserialize(t *testing.T) {

    x, err := Serialize(array_str)
    if err != nil {
        t.Fatal(err)
    }
    r, err := Deserialize(x)
    if err != nil {
        t.Fatal(err)
    }
    if reflect.ValueOf(r[0]).Len() != len(array_str)-1 {
        t.Fatal("Expecting array.")
    }

    x, err = Serialize(array_str...)
    if err != nil {
        t.Fatal(err)
    }
    r, err = Deserialize(x)
    if err != nil {
        t.Fatal(err)
    }
    if reflect.ValueOf(r).Len() != len(array_str) {
        t.Fatal("Expecting multi parameter.")
    }

    x, err = Serialize(map_list)
    if err != nil {
        t.Fatal(err)
    }
    r, err = Deserialize(x)
    if err != nil {
        t.Fatal(err)
    }
    if reflect.ValueOf(r[0]).Len() != len(map_list) {
        t.Fatal("Expecting map.")
    }
}
