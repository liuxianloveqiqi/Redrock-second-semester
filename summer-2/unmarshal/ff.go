package unmarshal

//
//import (
//	"errors"
//	"strconv"
//	"unicode"
//)
//
//// Unmarshal 函数用于将 JSON 字符串解析为结构体
//func Unmarshal(data []byte, v interface{}) error {
//	// 将传入的 v 参数转换为结构体指针
//	ptr, ok := v.(*Person)
//	if !ok {
//		return errors.New("只支持解析 *Person 类型")
//	}
//
//	// 解析 JSON 字符串的索引位置
//	index := 0
//
//	// 跳过空白字符
//	skipWhitespace(data, &index)
//
//	// 检查 JSON 字符串是否以左花括号开始
//	if data[index] != '{' {
//		return errors.New("JSON 字符串必须以左花括号开始")
//	}
//	index++
//
//	// 解析 JSON 键值对
//	for {
//		skipWhitespace(data, &index)
//
//		// 检查是否到达了 JSON 字符串的末尾
//		if index >= len(data) {
//			return errors.New("JSON 字符串不完整")
//		}
//
//		// 检查是否到达了 JSON 对象的结束符号
//		if data[index] == '}' {
//			break
//		}
//
//		// 解析键
//		key, err := parseString(data, &index)
//		if err != nil {
//			return err
//		}
//
//		// 跳过冒号
//		if data[index] != ':' {
//			return errors.New("JSON 键值对缺少冒号")
//		}
//		index++
//
//		// 解析值
//		if key == "name" {
//			name, err := parseString(data, &index)
//			if err != nil {
//				return err
//			}
//			ptr.Name = name
//		} else if key == "age" {
//			age, err := parseInt(data, &index)
//			if err != nil {
//				return err
//			}
//			ptr.Age = age
//		} else {
//			// 忽略未知键
//			skipValue(data, &index)
//		}
//
//		// 跳过逗号
//		if data[index] == ',' {
//			index++
//		}
//	}
//
//	// 检查 JSON 字符串是否以右花括号结束
//	if data[index] != '}' {
//		return errors.New("JSON 字符串必须以右花括号结束")
//	}
//
//	return nil
//}
//
//// 辅助函数：跳过空白字符
//func skipWhitespace(data []byte, index *int) {
//	for *index < len(data) && unicode.IsSpace(rune(data[*index])) {
//		*index++
//	}
//}
//
//// 辅助函数：解析字符串
//func parseString(data []byte, index *int) (string, error) {
//	// 跳过初始引号
//	*index++
//
//	start := *index
//	for *index < len(data) {
//		// 检查是否到达字符串结尾
//		if data[*index] == '"' && data[*index-1] != '\\' {
//			str := string(data[start:*index])
//			*index++
//			return str, nil
//		}
//		*index++
//	}
//
//	return "", errors.New("JSON 字符串不完整")
//}
//
//// 辅助函数：解析整数
//func parseInt(data []byte, index *int) (int, error) {
//	start := *index
//	for *index < len(data) {
//		// 检查是否到达整数结尾
//		if !unicode.IsDigit(rune(data[*index])) {
//			num, err := strconv.Atoi(string(data[start:*index]))
//			if err != nil {
//				return 0, errors.New("无效的整数")
//			}
//			return num, nil
//		}
//		*index++
//	}
//
//	return 0, errors.New("JSON 字符串不完整")
//}
