# Daytime

Golang package for working with time of day

## Get package

```bash
go get github.com/SPTolkachev/daytime
```

## Create

Create a new daytime

```go
daytime := New(hour int, minute int, second int)
```

Parse a daytime

```go
daytime := Parse("15:04:05")
```

## Convert to string

```go
str := daytime.String()
fmt.Println(str) // 15:04:05
```

## Convert to time

Bringing to the current day's time.

```go
daytime.Time()
```

Bringing to the near future.

```go
daytime.InTheNearFuture()
```

Bringing to the recent past.

```go
daytime.InTheRecentPast()
```
