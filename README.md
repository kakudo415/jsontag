# jsontag

カーソル位置に対応する構造体のフィールドにJSON Tagをつける。

## 使い方

### フィールドにタグを付ける例

```go
type hoge struct {
	FooBar   int
	FizzBazz int `gorm:"primaryKey"`
	fuga     int
}
```


```bash
go vet -vettool="main.exe" --jsontag.where=24 --jsontag.option=omitempty <source.go>
```

```go
type hoge struct {
	FooBar   int
	FizzBazz int `gorm:"primaryKey" json:"fizzBazz,omitempty"`
	fuga     int
}
```

## オプション

```
--jsontag.where=${cursor_offset}
```

```
--jsontag.option=${option}
```

| option    | 動作                         |
| --------- | ---------------------------- |
|           | \`json:"jsonTag"\`           |
| omitempty | \`json:"jsonTag,omitempty"\` |
| ignore    | \`json:"-"\`                 |
