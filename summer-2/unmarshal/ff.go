package unmarshal

//
//import (
//	"errors"
//	"reflect"
//	"strconv"
//)
//
// UnJSON 反序列化JSON格式的字节切片为Go数据结构

//
//// 解析JSON值
//func parseValue(data []byte, v reflect.Value) error {
//	// 去除空白字符
//	data = skipWhitespace(data)
//
//	// 根据值的类型进行解析
//	switch data[0] {
//	case '{':
//		return parseObject(data, v)
//	case '[':
//		return parseArray(data, v)
//	case '"':
//		return parseString(data, v)
//	case 't':
//		return parseTrue(data, v)
//	case 'f':
//		return parseFalse(data, v)
//	case 'n':
//		return parseNull(data, v)
//	default:
//		return parseNumber(data, v)
//	}
//}
//
//// 解析JSON对象
//func parseObject(data []byte, v reflect.Value) error {
//	// 创建一个空的 map[string]interface{}
//	obj := make(map[string]interface{})
//
//	// 跳过起始的"{"
//	data = data[1:]
//
//	for {
//		data = skipWhitespace(data)
//
//		// 判断对象是否为空
//		if len(data) == 0 || data[0] == '}' {
//			break
//		}
//
//		// 解析键
//		key, err := parseString(data, nil)
//		if err != nil {
//			return err
//		}
//
//		data = skipWhitespace(key[1:])
//		if len(data) == 0 || data[0] != ':' {
//			return errors.New("无效的JSON对象")
//		}
//
//		// 解析值
//		fieldName := string(key[1 : len(key)-1])
//		fieldValue := reflect.New(v.Type().Elem().FieldByName(fieldName).Type).Elem()
//		err = parseValue(data[1:], fieldValue)
//		if err != nil {
//			return err
//		}
//
//		// 将键值对存入对象
//		obj[fieldName] = fieldValue.Interface()
//
//		data = skipWhitespace(data)
//		if len(data) == 0 || data[0] == '}' {
//			break
//		} else if data[0] != ',' {
//			return errors.New("无效的JSON对象")
//		}
//
//		data = data[1:]
//	}
//
//	// 将解析后的对象存入目标变量
//	v.Set(reflect.ValueOf(obj))
//
//	return nil
//}
//// 跳过空白字符
//func skipWhitespace(data []byte) []byte {
//	i := 0
//	for i < len(data) && (data[i] == ' ' || data[i] == '\t' || data[i] == '\n' || data[i] == '\r') {
//		i++
//	}
//	return data[i:]
//}
//// 解析JSON字符串
//func parseString(data []byte, v reflect.Value) ([]byte, error) {
//	// 查找字符串的结束位置
//	endIndex := 0
//	for endIndex < len(data) && data[endIndex] != '"' {
//		// 转义字符
//		if data[endIndex] == '\\' {
//			endIndex++
//		}
//		endIndex++
//	}
//
//	if endIndex >= len(data) {
//		return nil, errors.New("无效的JSON字符串")
//	}
//
//	// 提取字符串内容
//	str := data[0 : endIndex+1]
//
//	// 将解析后的字符串存入目标变量
//	if v.IsValid() {
//		v.SetString(string(str[1 : len(str)-1]))
//	}
//
//	return str, nil
//}
//
//// 解析JSON中的数字
//func parseNumber(data []byte, v interface{}) error {
//	// 查找数字的结束位置
//	endIndex := 0
//	for endIndex < len(data) && isValidNumberChar(data[endIndex]) {
//		endIndex++
//	}
//
//	if endIndex == 0 {
//		return errors.New("无效的JSON数字")
//	}
//
//	// 提取数字内容
//	number := data[0:endIndex]
//
//	// 解析数字
//	parsedNumber, err := parseNumberValue(number)
//	if err != nil {
//		return err
//	}
//
//	// 将解析后的数字存入目标变量
//	if v != nil {
//		reflectValue := reflect.ValueOf(v).Elem()
//
//		switch reflectValue.Kind() {
//		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
//			intValue, err := strconv.ParseInt(parsedNumber, 10, 64)
//			if err != nil {
//				return err
//			}
//			reflectValue.SetInt(intValue)
//		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
//			uintValue, err := strconv.ParseUint(parsedNumber, 10, 64)
//			if err != nil {
//				return err
//			}
//			reflectValue.SetUint(uintValue)
//		case reflect.Float32, reflect.Float64:
//			floatValue, err := strconv.ParseFloat(parsedNumber, 64)
//			if err != nil {
//				return err
//			}
//			reflectValue.SetFloat(floatValue)
//		default:
//			return errors.New("无效的目标变量类型")
//		}
//	}
//
//	return nil
//}
//
//// 检查字符是否是有效的数字字符
//func isValidNumberChar(char byte) bool {
//	return (char >= '0' && char <= '9') || char == '-' || char == '+' || char == '.' || char == 'e' || char == 'E'
//}
//
//// 解析数字字符串为对应的值
//func parseNumberValue(number []byte) (string, error) {
//	// 特殊情况：处理负数
//	if len(number) > 1 && number[0] == '-' && number[1] == '0' {
//		return string(number), errors.New("无效的JSON数字")
//	}
//
//	// 检查是否为十进制数
//	if number[0] == '0' && len(number) > 1 {
//		return string(number), errors.New("无效的JSON数字")
//	}
//
//	return string(number), nil
//}
//
