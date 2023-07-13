package marshal

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMap(t *testing.T) {
	m := make(map[int]string, 0)
	m[3] = "sd"
	m[2] = "2323"
	jsonData, err := JSON(m)
	if err != nil {
		fmt.Println("序列化失败:", err)
		return
	}
	// 打印json数据
	fmt.Println(string(jsonData))
	unData := make(map[int]string, 0)
	// 调用官方json反序列化包验证
	err = json.Unmarshal(jsonData, &unData)
	if err != nil {
		panic(err)
	}
	// 打印反序列化后数据
	fmt.Println(unData)
}
func TestSlice(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	jsonData, err := JSON(s)
	if err != nil {
		fmt.Println("序列化失败:", err)
		return
	}
	// 打印json数据
	fmt.Println(string(jsonData))
	unData := make([]int, 0)
	// 调用官方json反序列化包验证
	err = json.Unmarshal(jsonData, &unData)
	if err != nil {
		panic(err)
	}
	// 打印反序列化后数据
	fmt.Println(unData)
}
