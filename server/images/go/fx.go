package main

import "fmt"

type Input struct {
	name string
}

type Output struct {
	name string
}

func Fx(input Input) (output Output) {
	fmt.Println(input)
	output = Output{"my name is fx"}
	return output
}
