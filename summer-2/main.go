package main

import (
	"encoding/json"
	"fmt"
	"summer-2/marshal"
	"summer-2/model"
	"summer-2/unmarshal"
)

// 示例结构体
type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	// 随便搞个结构体
	p := Person{
		Name: "liuxian",
		Age:  18,
	}
	// 测试序列化
	jsonData, err := marshal.JSON(p)
	if err != nil {
		fmt.Println("序列化失败:", err)
		return
	}
	// 打印json数据
	fmt.Println(string(jsonData))
	unData := model.Person{}
	// 调用官方json反序列化包验证
	err = json.Unmarshal(jsonData, &unData)
	if err != nil {
		panic(err)
	}
	// 打印反序列化后数据
	fmt.Println(unData)

	// 测试反序列化
	jsonData, err = json.Marshal(p)
	if err != nil {
		panic(err)
	}
	newP := model.Person{}
	// 调用自己的反序列化
	err = unmarshal.UnJSON(jsonData, &newP)
	if err != nil {
		panic(err)
	}
	fmt.Println(newP)
}
