package main

type Input struct {
	A int32
	B int32
}

type Output struct {
	Sum int32
}

func Fx(input *Input) (output *Output) {
	output = &Output{
		Sum: input.A + input.B,
	}
	return
}
