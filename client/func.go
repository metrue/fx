package main

type Input struct {
	Age int32
}

type Output struct {
	Name string
	Age int32
}

func Fx(input *Input) (output *Output) {
	output = &Output{
		Name: "my name is fx",
		Age: input.Age,
	}
	return
}
