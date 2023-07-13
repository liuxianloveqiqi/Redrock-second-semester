package marshal

import (
	"bytes"
	"reflect"
	"strconv"
)

func JSON(data interface{}) ([]byte, error) {
	// 先利用反射获取数据的类型和值
	dataType := reflect.TypeOf(data)
	value := reflect.ValueOf(data)
	// 根据数据类型选择对应的序列化方式
	switch dataType.Kind() {
	case reflect.Slice, reflect.Array:
		// 切片或数组类型
		buffer := bytes.Buffer{}
		buffer.WriteString("[")
		// 递归序列化切片或数组的每个元素
		for i := 0; i < value.Len(); i++ {
			// 递归调用JSON函数对每个元素进行序列化
			elemData, err := JSON(value.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			// 追加到buffer
			buffer.Write(elemData)
			// 添加逗号分隔符
			if i < value.Len()-1 {
				buffer.WriteString(",")
			}
		}
		// 添加最右方括号
		buffer.WriteString("]")

		return buffer.Bytes(), nil

	case reflect.Map:
		// map类型
		buffer := bytes.Buffer{}
		buffer.WriteString("{")
		// 获取map所有的key
		keys := value.MapKeys()

		for i := 0; i < len(keys); i++ {
			// 获取 map 的键和值
			key := keys[i]
			value := value.MapIndex(key)

			// 对键进行序列化
			keyData, err := JSON(key.Interface())
			if err != nil {
				return nil, err
			}

			// 对值进行序列化
			valueData, err := JSON(value.Interface())
			if err != nil {
				return nil, err
			}

			// 将序列化后的键和值添加到缓冲区中，并添加引号
			// 使用 strconv.Quote将键序列化为一个带有引号的字符串
			buffer.WriteString(strconv.Quote(string(keyData)))
			// 用":"隔开
			buffer.WriteString(":")
			// 同理给值加引号
			buffer.WriteString(strconv.Quote(string(valueData)))

			// 如果不是最后一个键值对，就用逗号进行分隔
			if i < len(keys)-1 {
				buffer.WriteString(",")
			}
		}

		buffer.WriteString("}")

		return buffer.Bytes(), nil

	case reflect.String:
		// 字符串类型
		str := data.(string)
		return []byte(strconv.Quote(str)), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// 整数类型
		num := reflect.ValueOf(data).Int()
		return []byte(strconv.FormatInt(num, 10)), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// 无符号整数类型
		num := reflect.ValueOf(data).Uint()
		return []byte(strconv.FormatUint(num, 10)), nil

	case reflect.Float32, reflect.Float64:
		// 浮点数类型
		num := reflect.ValueOf(data).Float()
		return []byte(strconv.FormatFloat(num, 'f', -1, 64)), nil

	case reflect.Bool:
		// 布尔类型
		b := reflect.ValueOf(data).Bool()
		if b {
			return []byte("true"), nil
		}
		return []byte("false"), nil

	default:
		// 对于其他复杂类型，递归处理字段
		return handleComplexType(data)
	}
}

// 处理复杂类型的序列化
func handleComplexType(data interface{}) ([]byte, error) {
	// 获取数据的值和类型
	value := reflect.ValueOf(data)
	dataType := value.Type()

	// 创建一个缓冲区，用于拼接JSON字符串
	buffer := bytes.Buffer{}
	buffer.WriteString("{")

	// 遍历结构体的字段
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := dataType.Field(i)

		// 获取字段的名称和标签
		name := fieldType.Name
		tag := fieldType.Tag.Get("json")

		// 如果标签为空，则默认使用字段名作为JSON的键
		if tag == "" {
			tag = name
		}

		// 这次也是采用递归处理字段的值
		fieldData, err := JSON(field.Interface())
		if err != nil {
			return nil, err
		}

		// 拼接JSON键值对
		buffer.WriteString("\"" + tag + "\":")
		buffer.Write(fieldData)

		// 添加','隔开
		if i < value.NumField()-1 {
			buffer.WriteString(",")
		}
	}
	// 补上}
	buffer.WriteString("}")

	return buffer.Bytes(), nil
}
