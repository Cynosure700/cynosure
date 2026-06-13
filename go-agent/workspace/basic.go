package main

import "fmt"

func main() {
	// 变量声明
	var name string = "Go"
	age := 10 // 简短声明
	var version float64 = 1.21

	// 常量
	const greeting = "Hello"

	// 条件语句
	if age > 5 {
		fmt.Println(greeting, name, "版本", version)
	}

	// 循环
	for i := 0; i < 3; i++ {
		fmt.Println("循环第", i+1, "次")
	}

	// 数组
	var arr [3]int = [3]int{1, 2, 3}

	// 切片
	slice := []string{"a", "b", "c"}
	slice = append(slice, "d")

	// Map
	m := map[string]int{"x": 10, "y": 20}

	// 函数调用
	sum := add(3, 4)
	fmt.Println("3+4 =", sum)

	// 打印其他结构
	fmt.Println("数组:", arr)
	fmt.Println("切片:", slice)
	fmt.Println("Map:", m)
}

func add(a, b int) int {
	return a + b
}
