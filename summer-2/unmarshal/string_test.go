package unmarshal

//
//import (
//	"errors"
//	"fmt"
//	"testing"
//)
//
//func TestString(t *testing.T) {
//	data := []byte(`"name":"liuxian","age":18}`)
//	//data := []byte(`{"name":"liuxian","age":18}`)
//
//	var name string
//	var age string
//
//	// 解析name字段
//	remaining, err := parseString(data, &name)
//	if err != nil {
//		fmt.Println("解析name字段失败:", err)
//		return
//	}
//	fmt.Println("name:", name)
//
//	// 解析age字段
//	remaining, err = parseString(remaining[1:], &age)
//	if err != nil {
//		fmt.Println("解析age字段失败:", err)
//		return
//	}
//	fmt.Println("age:", age)
//}
//func TestOB(t *testing.T) {
//	data := []byte(`{"name":"liuxian","age":18}`)
//	result := make(map[string]interface{}) // Declare result as a map pointer
//	err := parseObject(data, &result)      // Pass the address of the map
//	if err != nil {
//		fmt.Println("解析JSON对象失败:", err)
//		return
//	}
//
//	// Access the parsed values
//	name, ok := result["name"].(string)
//	age, ok := result["age"].(float64)
//	if !ok {
//		fmt.Println(errors.New("无效的JSON对象"))
//	}
//
//	fmt.Println("name:", name)
//	fmt.Println("age:", age)
//}
