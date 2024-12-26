
```go
type Cast struct {
  Actor, Role string
}

type Serie struct {
  Cast []Cast
}

series := map[string]Serie{
  "A-Team": {Cast: []Cast{
    {Actor: "George Peppard", Role: "Hannibal"},
    {Actor: "Dwight Schultz", Role: "Murdock"},
    {Actor: "Mr. T", Role: "Baracus"},
    {Actor: "Dirk Benedict", Role: "Faceman"},
  }},
}

value, _ := lookup.DotP(series, "A-Team.Cast.0.Actor")
fmt.Println(q, "->", value.Interface())
```

