package unmarshal

import (
	"errors"
	"strconv"
	"summer-2/model"

	"unicode"
)

func UnJSON(data []byte, v interface{}) error {
	// 将传入的v参数转换为结构体指针
	ptr, ok := v.(*model.Person)
	if !ok {
		return errors.New("只支持解析*Person类型")
	}

	// 解析JSON字符串的索引位置
	index := 0

	// 跳过空白字符
	skipWhitespace(data, &index)

	// 检查是否以左花括号开始
	if data[index] != '{' {
		return errors.New("JSON字符串必须以左花括号开始")
	}
	index++

	// 解析JSON键值对
	for {
		skipWhitespace(data, &index)

		// 检查是否到达了JSON字符串的末尾
		if index >= len(data) {
			return errors.New("JSON字符串不完整")
		}

		// 检查是否到达了的结束符号
		if data[index] == '}' {
			break
		}

		// 解析键
		key, err := parseString(data, &index)
		if err != nil {
			return err
		}

		// 跳过冒号
		if data[index] != ':' {
			return errors.New("JSON键值对缺少冒号")
		}
		index++

		// 解析值
		if key == "name" {
			name, err := parseString(data, &index)
			if err != nil {
				return err
			}
			ptr.Name = name
		} else if key == "age" {
			age, err := parseInt(data, &index)
			if err != nil {
				return err
			}
			ptr.Age = age
		} else {
			// 忽略未知键
			skipValue(data, &index)
		}

		// 跳过逗号
		if data[index] == ',' {
			index++
		}
	}

	// 检查是否以右花括号结束
	if data[index] != '}' {
		return errors.New("JSON字符串必须以右花括号结束")
	}

	return nil
}

// 辅助函数：跳过值
func skipValue(data []byte, index *int) {
	// 记录嵌套层数
	depth := 0
	// 标记是否在字符串内部
	inString := false

	for *index < len(data) {
		switch data[*index] {
		case '{':
			if !inString {
				depth++
			}
		case '}':
			if !inString {
				depth--
				if depth == 0 {
					return
				}
			}
		case '[':
			if !inString {
				depth++
			}
		case ']':
			if !inString {
				depth--
				if depth == 0 {
					return
				}
			}
		case '"':
			if *index > 0 && data[*index-1] != '\\' {
				inString = !inString
			}
		}

		*index++
	}
}

// 跳过空白字符
func skipWhitespace(data []byte, index *int) {
	for *index < len(data) && unicode.IsSpace(rune(data[*index])) {
		*index++
	}
}

// 解析字符串
func parseString(data []byte, index *int) (string, error) {
	// 跳过初始引号
	*index++
	start := *index
	for *index < len(data) {
		// 检查是否到达字符串结尾
		if data[*index] == '"' && data[*index-1] != '\\' {
			str := string(data[start:*index])
			*index++
			return str, nil
		}
		*index++
	}

	return "", errors.New("JSON字符串不完整")
}

// 解析整数
func parseInt(data []byte, index *int) (int, error) {
	start := *index
	for *index < len(data) {
		// 检查是否到达整数结尾
		if !unicode.IsDigit(rune(data[*index])) {
			num, err := strconv.Atoi(string(data[start:*index]))
			if err != nil {
				return 0, errors.New("无效的整数")
			}
			return num, nil
		}
		*index++
	}

	return 0, errors.New("JSON字符串不完整")
}
