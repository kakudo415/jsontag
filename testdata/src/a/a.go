package a

type hoge struct {
	FooBar   int
	FizzBazz int
}

func main() {
	var fuga = hoge{
		FooBar:   3,
		FizzBazz: 2,
	}
	println(fuga)
}
