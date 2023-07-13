package marshal

import (
	"bytes"
	"reflect"
	"strconv"
)

// JSON 将Go数据结构转换为JSON格式的字节切片
func JSON(data interface{}) ([]byte, error) {
	// 利用反射获取数据的类型
	dataType := reflect.TypeOf(data)

	// 根据数据类型选择对应的序列化方式
	switch dataType.Kind() {
	case reflect.String:
		// 如果是字符串类型，则直接转义并返回JSON格式的字符串
		str := data.(string)
		return []byte(strconv.Quote(str)), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// 如果是整数类型，则将其转换为字符串并返回JSON格式的字符串
		num := reflect.ValueOf(data).Int()
		return []byte(strconv.FormatInt(num, 10)), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// 如果是无符号整数类型，则将其转换为字符串并返回JSON格式的字符串
		num := reflect.ValueOf(data).Uint()
		return []byte(strconv.FormatUint(num, 10)), nil

	case reflect.Float32, reflect.Float64:
		// 如果是浮点数类型，则将其转换为字符串并返回JSON格式的字符串
		num := reflect.ValueOf(data).Float()
		return []byte(strconv.FormatFloat(num, 'f', -1, 64)), nil

	case reflect.Bool:
		// 如果是布尔类型，则返回JSON格式的布尔值
		b := reflect.ValueOf(data).Bool()
		if b {
			return []byte("true"), nil
		}
		return []byte("false"), nil

	default:
		// 对于其他复杂类型，递归处理字段，并构建JSON格式的字节切片
		return handleComplexType(data)
	}
}

// handleComplexType 处理复杂类型的序列化
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

// 示例结构体
type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}