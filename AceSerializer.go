package AceSerializer

import (
    "errors"
    "fmt"
    "math"
    "reflect"
    "regexp"
    "strings"
)

const serInf, serNegInf = "1.#INF", "-1.#INF"

func serializeStringHelper(str string) (result string, err error) {
    var n = str[0]
    if n == 0x1E {
        result = "\x7E\x7A"
    } else if n <= 0x20 {
        result = "\x7E" + string(n+64)
    } else if n == 0x5E {
        result = "\x7E\x7D"
    } else if n == 0x7E {
        result = "\x7E\x7C"
    } else if n == 0x7F {
        result = "\x7E\x7B"
    } else {
        err = errors.New("error type of str: " + str)
    }
    return
}

func changeNumber(v interface{}) string {
    switch v.(type) {
    case float64:
        if math.IsInf(v.(float64), 1) {
            return serInf
        } else if math.IsInf(v.(float64), -1) {
            return serNegInf
        } else {
            return fmt.Sprintf("%.14g", v)
        }
    }
    return AnyToStr(v)
}

func serializeValue(v interface{}, res map[int]string, nres int) (int, error) {

    var err error

    switch reflect.ValueOf(v).Kind() {
    case reflect.String:
        res[nres+1] = "^S"
        reg := regexp.MustCompile(`[[:cntrl:] \x5E\x7E\x7F]`)
        result := AnyToStr(v)
        for _, x := range reg.FindAllString(result, -1) {
            newX, err := serializeStringHelper(x)
            if err != nil {
                return 0, err
            }
            result = strings.Replace(result, x, newX, 1)
        }
        res[nres+2] = result
        nres = nres + 2
    case reflect.Map:
        nres++
        res[nres] = "^T"

        value := reflect.ValueOf(v)
        for _, k := range value.MapKeys() {
            kType := k.Kind()

            if kType == reflect.Int {
                nres, err = serializeValue(k.Int()+1, res, nres)
                if err != nil {
                    return 0, err
                }
            } else if kType == reflect.Interface {
                k := k.Interface()
                switch k.(type) {
                case int:
                    nres, err = serializeValue(k.(int)+1, res, nres)
                    if err != nil {
                        return 0, err
                    }
                case string:
                    nres, err = serializeValue(k.(string), res, nres)
                    if err != nil {
                        return 0, err
                    }
                default:
                    return 0, errors.New("error map of key.")
                }
            } else {
                nres, err = serializeValue(k.String(), res, nres)
                if err != nil {
                    return 0, err
                }
            }
            nres, err = serializeValue(value.MapIndex(k).Interface(), res, nres)
            if err != nil {
                return 0, err
            }
        }

        nres++
        res[nres] = "^t"
    case reflect.Bool:
        nres++
        if v == true {
            res[nres] = "^B"
        } else {
            res[nres] = "^b"
        }
    case reflect.Float64, reflect.Int:
        str := changeNumber(v)
        if str == serInf || str == serNegInf || AnyToStr(v) == str {
            res[nres+1] = "^N"
            res[nres+2] = str
            nres = nres + 2
        } else {
            m, e := math.Frexp(AnyToFloat(v))
            res[nres+1] = "^F"
            res[nres+2] = fmt.Sprintf("%.0f", m*math.Pow(2, 53))
            res[nres+3] = "^f"
            res[nres+4] = AnyToStr(e - 53)
            nres = nres + 4
        }
    case reflect.Invalid:
        nres++
        res[nres] = "^Z"
    case reflect.Slice, reflect.Array:
        nres++
        res[nres] = "^T"

        value := reflect.ValueOf(v)

        for i := 0; i < value.Len(); i++ {
            item := value.Index(i).Interface()

            if item == nil {
                res[nres+1] = ""
                res[nres+2] = ""
                nres = nres + 2
            } else {
                nres, err = serializeValue(i+1, res, nres)
                if err != nil {
                    return 0, err
                }

                nres, err = serializeValue(item, res, nres)
                if err != nil {
                    return 0, err
                }
            }
        }

        nres++
        res[nres] = "^t"
    default:
        return 0, errors.New("Cannot serialize a value of type")
    }

    return nres, nil
}

func Serialize(args ...interface{}) (string, error) {
    nres := 0
    serializeTbl := make(map[int]string)
    serializeTbl[nres] = "^1"

    for _, arg := range args {
        r, err := serializeValue(arg, serializeTbl, nres)
        if err != nil {
            return "", err
        }
        nres = r
    }

    serializeTbl[nres+1] = "^^"
    num := len(serializeTbl)
    result := make([]string, num)

    for k, v := range serializeTbl {
        result[k] = v
    }

    return strings.Join(result, ""), nil
}

func deserializeStringHelper(escape string) (result string, err error) {
    if escape < "~\x7A" {
        result = string(escape[1] - 64)
    } else if escape == "~\x7A" {
        result = "\x1E"
    } else if escape == "~\x7B" {
        result = "\x7F"
    } else if escape == "~\x7C" {
        result = "\x7E"
    } else if escape == "~\x7D" {
        result = "\x5E"
    } else {
        err = errors.New("DeserializeStringHelper got called for '" + escape + "'?!?")
    }
    return
}

func deserializeNumberHelper(number string) interface{} {
    if number == serNegInf {
        return math.Inf(-1)
    } else if number == serInf {
        return math.Inf(1)
    } else {
        return AnyToFloat(number)
    }
}

func gmatch(str string, reg string) func() map[string]string {

    result := regexp.MustCompile(reg).FindAllString(str, -1)
    index := -1
    num := len(result)

    return func() map[string]string {
        index++
        if index == num {
            return nil
        } else {
            rs := []rune(result[index])
            lth := len(rs)
            return map[string]string{
                "ctl":  string(rs[0:2]),
                "data": string(rs[2:lth]),
            }
        }
    }
}

func deserializeValue(iter func() map[string]string, single bool, ctl string, data string) (interface{}, error) {
    if !single {
        r := iter()
        ctl = r["ctl"]
        data = r["data"]
    }

    if len(ctl) == 0 {
        return nil, errors.New("Supplied data misses AceSerializer terminator ('^^')")
    }

    if ctl == "^^" {
        return nil, nil
    }

    var res interface{}
    if ctl == "^S" {
        reg := regexp.MustCompile("~.")
        res = data
        for _, v := range reg.FindAllString(data, -1) {
            _v, err := deserializeStringHelper(v)
            if err != nil {
                return nil, err
            }
            res = strings.Replace(res.(string), v, _v, 1)
        }
    } else if ctl == "^N" {
        res = deserializeNumberHelper(data)
        if res == nil {
            return nil, errors.New("Invalid serialized number: '" + data + "'")
        }
    } else if ctl == "^F" {
        r := iter()
        if r["ctl"] != "^f" {
            return nil, errors.New("Invalid serialized floating-point number, expected '^f', not '" + r["ctl"] + "'")
        }
        m := AnyToFloat(data)
        e := AnyToFloat(r["data"])
        if m == 0 || e == 0 {
            return nil, errors.New("Invalid serialized floating-point number, expected mantissa and exponent, got 0")
        }
        res = m * math.Pow(2, e)
    } else if ctl == "^B" {
        res = true
    } else if ctl == "^b" {
        res = false
    } else if ctl == "^Z" {
        res = nil
    } else if ctl == "^T" {
        res := make(map[interface{}]interface{})

        for {
            r := iter()

            if r["ctl"] == "^t" {
                break
            }

            k, err := deserializeValue(iter, true, r["ctl"], r["data"])
            if err != nil {
                return nil, err
            }

            if k == nil {
                return nil, errors.New("Invalid AceSerializer table format (no table end marker)")
            }

            if reflect.ValueOf(k).Kind() == reflect.Int {
                k = reflect.ValueOf(k).Int() - 1
            }

            r = iter()
            v, err := deserializeValue(iter, true, r["ctl"], r["data"])

            if err != nil {
                return nil, err
            }

            if v == nil {
                return nil, errors.New("Invalid AceSerializer table format (no table end marker)")
            }
            res[k] = v
        }

        return res, nil
    } else {
        return nil, errors.New("Invalid AceSerializer control code '" + ctl + "'")
    }

    return res, nil
}

func Deserialize(str string) ([]interface{}, error) {

    str = regexp.MustCompile(`[\x01-\x20\x7F]`).ReplaceAllString(str, "")

    iter := gmatch(str, `(\^.)([^^]*)`)
    r := iter()

    if r == nil || r["ctl"] != "^1" {
        return nil, errors.New("Supplied data is not AceSerializer data (rev 1)")
    }

    result := make([]interface{}, 0)

    for {
        r := iter()
        if r == nil {
            break
        }

        res, err := deserializeValue(iter, true, r["ctl"], r["data"])
        if err != nil {
            return nil, err
        }

        if r["ctl"] != "^^" {
            result = append(result, res)
        }
    }

    return result, nil
}
