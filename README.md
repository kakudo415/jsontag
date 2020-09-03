# jsontag

カーソル位置に対応する構造体のフィールドにJSON Tagをつける。

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
