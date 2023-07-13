package main

import (
	"encoding/json"
	"fmt"
	"summer-2/marshal"
	"summer-2/unmarshal"
)

func main() {
	// 随便搞个结构体
	p := marshal.Person{
		Name: "liuxian",
		Age:  18,
	}
	//// 序列化为JSON格式的字节切片
	//jsonData, err := marshal.JSON(p)
	//if err != nil {
	//	fmt.Println("序列化失败:", err)
	//	return
	//}
	//// 打印json数据
	//fmt.Println(string(jsonData))
	//unData := marshal.Person{}
	//// 调用官方json反序列化包验证
	//err = json.Unmarshal(jsonData, &unData)
	//if err != nil {
	//	panic(err)
	//}
	//// 打印反序列化后数据
	//fmt.Println(unData)

	jsondata, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err)
	}
	newunData := marshal.Person{}
	fmt.Println(string(jsondata))
	err = unmarshal.UnJSON(jsondata, &newunData)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(newunData)
}
