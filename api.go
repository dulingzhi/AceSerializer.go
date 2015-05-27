package AceSerializer

import (
    "fmt"
    "reflect"
    "strconv"
    "time"
)

func AnyToFloat(v interface{}) float64 {
    switch v.(type) {
    // #1
    // 这个nil必须保持，不然在检索结构的方法时，有可能会陷入死循环
    case nil:
        return 0
    // #2
    case bool:
        if v == true {
            return 1
        }
    // #3
    // 这玩意可真算不上优雅啊，go怎么就没有泛型呢？
    case int:
        if conv, ok := v.(int); ok {
            return float64(conv)
        }
    case int8:
        if conv, ok := v.(int8); ok {
            return float64(conv)
        }
    case int16:
        if conv, ok := v.(int16); ok {
            return float64(conv)
        }
    case int32:
        if conv, ok := v.(int32); ok {
            return float64(conv)
        }
    case int64:
        if conv, ok := v.(int64); ok {
            return float64(conv)
        }
    case uint:
        if conv, ok := v.(uint); ok {
            return float64(conv)
        }
    case uint8:
        if conv, ok := v.(uint8); ok {
            return float64(conv)
        }
    case uint16:
        if conv, ok := v.(uint16); ok {
            return float64(conv)
        }
    case uint32:
        if conv, ok := v.(uint32); ok {
            return float64(conv)
        }
    case uint64:
        if conv, ok := v.(uint64); ok {
            return float64(conv)
        }   // 这里仍然是有问题
    // #4
    case float32:
        if conv, ok := v.(float32); ok {
            return float64(conv)
        }
    case float64:
        if conv, ok := v.(float64); ok {
            return float64(conv)
        }
    // #5
    case string:
        if conv, ok := v.(string); ok {
            return StrToFloat(conv)
        }
    // #6
    case time.Time:
        if conv, ok := v.(time.Time); ok {
            return float64(conv.Unix())
        }
    // #999
    default:
        // 数组、切片、Map转类型是什么类型呢？
        return AnyToFloat(CallAnyStructMethod(v, "Float"))
    }
    return 0
}

// 注意，所有其他的AnyTo转换，都不处理[]byte，因为实际上[]byte的情况会比较复杂，他可能包含了encode/gob的编码格式，也可能是json格式
// 也可能用户自己打包的，所以我们不做任何处理
// 但AnyToStr的话还是要处理，尝试最简单的转换
func AnyToStr(v interface{}) string {
    switch v.(type) {
    // #1
    // 这个nil必须保持，不然在检索结构的方法时，有可能会陷入死循环
    case nil:
        return ""
    // #2
    // 布尔类型，应该返回个啥呢？真头疼，暂时先返回一个1吧，总比返回了true好
    case bool:
        if v == true {
            return "1"
        }
    // #3
    // 这玩意可真算不上优雅啊，go怎么就没有泛型呢？
    case int:
        if conv, ok := v.(int); ok {
            return strconv.Itoa(conv)
        }
    case int8:
        if conv, ok := v.(int8); ok {
            return strconv.Itoa(int(conv))
        }
    case int16:
        if conv, ok := v.(int16); ok {
            return strconv.Itoa(int(conv))
        }
    case int32:
        if conv, ok := v.(int32); ok {
            return strconv.Itoa(int(conv))
        }   // 32bit 64bit系统都能涵盖了这个值
    case int64:
        if conv, ok := v.(int64); ok {
            return fmt.Sprint(conv)
        }
    case uint:
        if conv, ok := v.(uint); ok {
            return fmt.Sprint(conv)
        }
    case uint8:
        if conv, ok := v.(uint8); ok {
            return strconv.Itoa(int(conv))
        }
    case uint16:
        if conv, ok := v.(uint16); ok {
            return strconv.Itoa(int(conv))
        }
    case uint32:
        if conv, ok := v.(uint32); ok {
            return fmt.Sprint(conv)
        }   // 32无负数整型，转int就少了一截了
    case uint64:
        if conv, ok := v.(uint64); ok {
            return fmt.Sprint(conv)
        }   // 64位无负数整型，就更加是少了一截了。
    // #4
    case float32:
        if conv, ok := v.(float32); ok {
            return strconv.FormatFloat(float64(conv), 'f', -1, 64)
        }
    case float64:
        if conv, ok := v.(float64); ok {
            return strconv.FormatFloat(conv, 'f', -1, 64)
        }
    // #5
    case []byte:
        if conv, ok := v.([]byte); ok {
            return string(conv)
        }
    case string:
        if conv, ok := v.(string); ok {
            return conv
        }
    // #6
    case time.Time:
        if conv, ok := v.(time.Time); ok {
            return conv.String()
        }
    // #999
    default:
        // 数组、切片、Map转类型是什么类型呢？
        return AnyToStr(CallAnyStructMethod(v, "String"))
    }
    return ""
}

// 转型最好优先转型到最大的值，然后再往底缩进
// 更精确的做法，应该是根据位长，来做出适当的判断但过度优化，又不如直接用go提供一些方法
// 所以这个方法只是确保值的有效性转换，性能在能考虑的条件下，才考虑
func AnyToInt64(v interface{}) int64 {
    switch v.(type) {
    // #1
    // 这个nil必须保持，不然在检索结构的方法时，有可能会陷入死循环
    case nil:
        return 0
    // #2
    case bool:
        if v == true {
            return 1
        }
    // #3
    // 这玩意可真算不上优雅啊，go怎么就没有泛型呢？
    case int:
        if conv, ok := v.(int); ok {
            return int64(conv)
        }
    case int8:
        if conv, ok := v.(int8); ok {
            return int64(conv)
        }
    case int16:
        if conv, ok := v.(int16); ok {
            return int64(conv)
        }
    case int32:
        if conv, ok := v.(int32); ok {
            return int64(conv)
        }
    case int64:
        if conv, ok := v.(int64); ok {
            return int64(conv)
        }
    case uint:
        if conv, ok := v.(uint); ok {
            return int64(conv)
        }
    case uint8:
        if conv, ok := v.(uint8); ok {
            return int64(conv)
        }
    case uint16:
        if conv, ok := v.(uint16); ok {
            return int64(conv)
        }
    case uint32:
        if conv, ok := v.(uint32); ok {
            return int64(conv)
        }
    case uint64:
        if conv, ok := v.(uint64); ok {
            return int64(conv)
        }   // 这里仍然是有问题
    // #4
    case float32:
        if conv, ok := v.(float32); ok {
            return int64(conv)
        }
    case float64:
        if conv, ok := v.(float64); ok {
            return int64(conv)
        }
    // #5
    case string:
        if conv, ok := v.(string); ok {
            return StrToFInt64(conv)
        }
    // #6
    case time.Time:
        if conv, ok := v.(time.Time); ok {
            return conv.Unix()
        }
    // #999
    default:
        // 数组、切片、Map转类型是什么类型呢？
        return AnyToInt64(CallAnyStructMethod(v, "Int"))
    }
    return 0
}

func AnyToInt(v interface{}) int {
    return int(AnyToInt64(v))
}

func StrToFloat(value string) float64 {
    float, err := strconv.ParseFloat(value, 64)
    if err != nil {
        return 0
    }
    return float
}

func CallAnyStructMethod(v interface{}, method string) interface{} {
    ref := reflect.ValueOf(v)
    refKind := ref.Kind()
    if refKind == reflect.Ptr {
        refKind = ref.Elem().Kind()
    }
    // 如果是结构的话，尝试检索一下他是否有Int、ToInt的函数
    if refKind == reflect.Struct {
        fn := ref.MethodByName(method)
        if fn.IsValid() {
            rs := fn.Call(nil)
            if len(rs) > 0 {
                return rs[0].Interface()
            }
        }
    }
    return nil
}

func StrToFInt64(value string) int64 {
    float := StrToFloat(value)
    return int64(float)
}
