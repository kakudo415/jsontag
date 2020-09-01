package a

type hoge struct {
	FooBar   int // want "fooBar"
	FizzBazz int // want "fizzBazz"
}

func main() {
	var fuga = hoge{
		FooBar:   3,
		FizzBazz: 2,
	}
	println(fuga)
}
