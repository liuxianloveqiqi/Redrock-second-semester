package unmarshal

import (
	"fmt"
	"testing"
)

func TestString(t *testing.T) {
	data := []byte(`"name":"liuxian","age":18}`)
	//data := []byte(`{"name":"liuxian","age":18}`)

	var name string
	var age string

	// 解析name字段
	remaining, err := parseString(data, &name)
	if err != nil {
		fmt.Println("解析name字段失败:", err)
		return
	}
	fmt.Println("name:", name)

	// 解析age字段
	remaining, err = parseString(remaining[1:], &age)
	if err != nil {
		fmt.Println("解析age字段失败:", err)
		return
	}
	fmt.Println("age:", age)
}
