package a

type hoge struct {
	FooBar   int
	FizzBazz int `gorm:"primaryKey"`
	fuga     int // want "fuga is NOT exported"
}
