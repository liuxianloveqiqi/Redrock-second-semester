package unmarshal

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func UnJSON(data []byte, result interface{}) error {
	// 利用反射获取结果数据的值和类型

	resultType := reflect.TypeOf(result)

	// 检查结果数据类型是否为指针类型
	if resultType.Kind() != reflect.Ptr {
		return errors.New("目标对象必须是指针类型")
	}

	// 创建一个指向结果数据的指针
	resultPtr := reflect.ValueOf(result).Elem()

	// 解析JSON数据
	err := parseValue(data, resultPtr)
	if err != nil {
		return err
	}

	return nil
}

//// 自己实现JSON反序列化函数
//func UnJSON(data []byte, v interface{}) error {
//	// 反射获取目标对象的类型信息
//	targetType := reflect.TypeOf(v)
//
//	// 如果目标对象不是指针类型，则返回错误
//	if targetType.Kind() != reflect.Ptr {
//		return errors.New("目标对象必须是指针类型")
//	}
//
//	// 将JSON数据为解析为map[string]interface{}
//	jsonMap, err := parseJSON(data)
//	if err != nil {
//		return err
//	}
//
//	// 获取指针指向的实际对象的类型信息
//	targetType = targetType.Elem()
//
//	// 反射创建目标对象的实例
//	targetValue := reflect.New(targetType).Elem()
//
//	// 遍历目标对象的每个字段
//	for i := 0; i < targetType.NumField(); i++ {
//		// 获取字段信息和值
//		field := targetType.Field(i)
//		fieldValue := targetValue.Field(i)
//
//		// 获取字段的JSON标签名
//		tagName := field.Tag.Get("json")
//
//		// 如果JSON标签名为""，则忽略该字段
//		if tagName == "" {
//			continue
//		}
//
//		// 从JSON数据中获取对应的值
//		value, ok := jsonMap[tagName]
//		if !ok {
//			// 如果没有找到该字段，继续continue下一个字段
//			continue
//		}
//
//		// 使用反射设置字段的值
//		err := setFieldValue(value, fieldValue)
//		if err != nil {
//			return fmt.Errorf("设置字段%s的值失败: %s", field.Name, err)
//		}
//	}
//
//	// 将反序列化后的值设置回目标对象
//	reflect.ValueOf(v).Elem().Set(targetValue)
//
//	return nil
//}

// 解析JSON数据为map[string]interface{}
func parseJSON(data []byte) (map[string]interface{}, error) {
	var jsonData interface{}

	// 解析JSON数据
	err := parseValue(data, &jsonData)
	if err != nil {
		return nil, err
	}

	// 将解析后的数据转换为map[string]interface{}
	jsonMap, ok := jsonData.(map[string]interface{})
	if !ok {
		return nil, errors.New("解析结果不是map[string]interface{}类型")
	}

	return jsonMap, nil
}

// 解析JSON值
func parseValue(data []byte, v interface{}) error {
	// 去除空白字符
	data = skipWhitespace(data)

	// 根据值的类型进行解析
	switch data[0] {
	case '{':
		return parseObject(data, v)
	case '[':
		return parseArray(data, v)
	case '"':
		_, err := parseString(data, v)
		return err
	case 't':
		return parseTrue(data, v)
	case 'f':
		return parseFalse(data, v)
	case 'n':
		return parseNull(data, v)
	default:
		return parseNumber(data, v)
	}
}

// 跳过空白字符
func skipWhitespace(data []byte) []byte {
	i := 0
	for i < len(data) && (data[i] == ' ' || data[i] == '\t' || data[i] == '\n' || data[i] == '\r') {
		i++
	}
	return data[i:]
}

// 解析JSON对象

func parseObject(data []byte, v interface{}) error {
	// 创建一个空的 map[string]interface{}
	obj := make(map[string]interface{})

	// 跳过起始的"{"
	data = data[1:]
	fmt.Println(string(data))
	for {
		data = skipWhitespace(data)

		// 判断对象是否为空
		if len(data) == 0 || data[0] == '}' {
			break
		}

		// 解析键
		key, err := parseString(data, nil)
		if err != nil {
			return err
		}
		fmt.Println("xxx", string(key))
		fmt.Println("无效的JSON对象", data)
		// 更新 data 跳过键的内容
		data = data[len(key):]
		data = skipWhitespace(data)
		fmt.Println("you的JSON对象:", data)

		if len(data) == 0 || data[0] != ':' {
			return errors.New("无效的JSON对象")
		}

		// 解析值
		var value interface{}
		err = parseValue(data[1:], &value)
		if err != nil {
			return err
		}

		// 将键值对存入对象
		obj[string(key[1:len(key)-1])] = value

		data = skipWhitespace(data)
		if len(data) == 0 || data[0] == '}' {
			break
		} else if data[0] != ',' {
			return errors.New("无效的JSON对象")
		}

		data = data[1:]
	}

	// 将解析后的对象存入目标变量
	reflectValue := reflect.ValueOf(v).Elem()
	reflectValue.Set(reflect.ValueOf(obj))

	return nil
}

// 解析JSON数组
func parseArray(data []byte, v interface{}) error {
	// 创建一个空的[]interface{}
	arr := make([]interface{}, 0)

	// 跳过起始的"["
	data = data[1:]

	for {
		data = skipWhitespace(data)

		// 判断数组是否为空
		if data[0] == ']' {
			break
		}

		// 解析值
		value := make(map[string]interface{})
		err := parseValue(data, &value)
		if err != nil {
			return err
		}

		// 将元素存入数组
		arr = append(arr, value)

		data = skipWhitespace(data)
		if data[0] == ',' {
			data = data[1:]
		}
	}

	// 将解析后的数组存入目标变量
	reflectValue := reflect.ValueOf(v).Elem()
	reflectValue.Set(reflect.ValueOf(arr))

	return nil
}

func parseString(data []byte, v interface{}) ([]byte, error) {
	if len(data) < 2 || data[0] != '"' {
		return nil, errors.New("无效的JSON字符串")
	}

	// 查找字符串的结束位置
	endIndex := 1
	for endIndex < len(data) {
		if data[endIndex] == '"' && data[endIndex-1] != '\\' {
			break
		}
		endIndex++
	}

	if endIndex >= len(data) {
		return nil, errors.New("无效的JSON字符串")
	}

	// 提取字符串内容
	str := data[0 : endIndex+1]

	// 将解析后的字符串存入目标变量
	if v != nil {
		value, ok := v.(*string)
		if !ok {
			return nil, errors.New("无效的目标变量类型")
		}
		*value = string(str[1 : len(str)-1])
	}

	return str, nil
}

// 解析JSON中的true值
func parseTrue(data []byte, v interface{}) error {
	if len(data) < 4 || string(data[:4]) != "true" {
		return errors.New("无效的JSON true值")
	}

	// 将解析后的true值存入目标变量
	if v != nil {
		reflectValue := reflect.ValueOf(v).Elem()
		reflectValue.SetBool(true)
	}

	return nil
}

// 解析JSON中的false值
func parseFalse(data []byte, v interface{}) error {
	if len(data) < 5 || string(data[:5]) != "false" {
		return errors.New("无效的JSON false值")
	}

	// 将解析后的false值存入目标变量
	if v != nil {
		reflectValue := reflect.ValueOf(v).Elem()
		reflectValue.SetBool(false)
	}

	return nil
}

// 解析JSON中的null值
func parseNull(data []byte, v interface{}) error {
	if len(data) < 4 || string(data[:4]) != "null" {
		return errors.New("无效的JSON null值")
	}

	// 将解析后的null值存入目标变量
	if v != nil {
		reflectValue := reflect.ValueOf(v).Elem()
		reflectValue.Set(reflect.Zero(reflectValue.Type()))
	}

	return nil
}

// 解析JSON中的数字
func parseNumber(data []byte, v interface{}) error {
	// 查找数字的结束位置
	endIndex := 0
	for endIndex < len(data) && isValidNumberChar(data[endIndex]) {
		endIndex++
	}

	if endIndex == 0 {
		return errors.New("无效的JSON数字")
	}

	// 提取数字内容
	number := data[0:endIndex]

	// 解析数字
	parsedNumber, err := parseNumberValue(number)
	if err != nil {
		return err
	}

	// 将解析后的数字存入目标变量
	if v != nil {
		reflectValue := reflect.ValueOf(v).Elem()

		switch reflectValue.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intValue, err := strconv.ParseInt(parsedNumber, 10, 64)
			if err != nil {
				return err
			}
			reflectValue.SetInt(intValue)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			uintValue, err := strconv.ParseUint(parsedNumber, 10, 64)
			if err != nil {
				return err
			}
			reflectValue.SetUint(uintValue)
		case reflect.Float32, reflect.Float64:
			floatValue, err := strconv.ParseFloat(parsedNumber, 64)
			if err != nil {
				return err
			}
			reflectValue.SetFloat(floatValue)
		default:
			return errors.New("无效的目标变量类型")
		}
	}

	return nil
}

// 检查字符是否是有效的数字字符
func isValidNumberChar(char byte) bool {
	return (char >= '0' && char <= '9') || char == '-' || char == '+' || char == '.' || char == 'e' || char == 'E'
}

// 解析数字字符串为对应的值
func parseNumberValue(number []byte) (string, error) {
	// 特殊情况：处理负数
	if len(number) > 1 && number[0] == '-' && number[1] == '0' {
		return string(number), errors.New("无效的JSON数字")
	}

	// 检查是否为十进制数
	if number[0] == '0' && len(number) > 1 {
		return string(number), errors.New("无效的JSON数字")
	}

	return string(number), nil
}

//// 使用反射设置字段的值
//func setFieldValue(value interface{}, fieldValue reflect.Value) error {
//	// 根据字段的类型进行类型断言，并设置对应的值
//	switch fieldValue.Kind() {
//	case reflect.String:
//		fieldValue.SetString(value.(string))
//	case reflect.Int:
//		intValue, err := strconv.Atoi(value.(string))
//		if err != nil {
//			return err
//		}
//		fieldValue.SetInt(int64(intValue))
//	default:
//		return fmt.Errorf("不支持的字段类型: %s", fieldValue.Kind().String())
//	}
//
//	return nil
//}
