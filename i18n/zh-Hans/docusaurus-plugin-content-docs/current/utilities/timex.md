---
sidebar_position: 2
---

# Timex

`timex` 包提供三种自定义时间类型 — `DateTime`、`Date` 和 `Time` — 作为 `time.Time` 的替代品，内置 JSON 序列化、数据库扫描和丰富的操作方法。

这些类型在整个框架中广泛使用，包括所有审计模型字段（`CreatedAt`、`UpdatedAt`）。

## 类型概览

| 类型 | 格式 | Go Layout | 示例 |
| --- | --- | --- | --- |
| `timex.DateTime` | `YYYY-MM-DD HH:mm:ss` | `time.DateTime` | `"2024-03-15 14:30:00"` |
| `timex.Date` | `YYYY-MM-DD` | `time.DateOnly` | `"2024-03-15"` |
| `timex.Time` | `HH:mm:ss` | `time.TimeOnly` | `"14:30:00"` |

三种类型都实现了：
- `json.Marshaler` / `json.Unmarshaler` — 无时区的简洁 JSON 格式
- `sql.Scanner` / `driver.Valuer` — 数据库兼容
- `encoding.TextMarshaler` / `encoding.TextUnmarshaler`

## DateTime

### 创建

```go
now := timex.Now()                                    // 当前时间
dt := timex.Of(time.Now())                            // 从 time.Time 转换
dt, err := timex.Parse("2024-03-15 14:30:00")         // 从字符串解析
dt, err := timex.Parse("15/03/2024 14:30", "02/01/2006 15:04") // 自定义格式
dt := timex.FromUnix(1710510600, 0)                   // 从 Unix 时间戳
dt := timex.FromUnixMilli(1710510600000)               // 从毫秒时间戳
```

### 访问组件

```go
dt.Year()       // 2024
dt.Month()      // time.March
dt.Day()        // 15
dt.Hour()       // 14
dt.Minute()     // 30
dt.Second()     // 0
dt.Weekday()    // time.Friday
dt.YearDay()    // 75
```

### 算术运算

```go
dt.Add(2 * time.Hour)    // 加时间段
dt.AddDate(1, 2, 3)      // 加年、月、日
dt.AddDays(7)             // 加天
dt.AddMonths(3)           // 加月
dt.AddYears(1)            // 加年
dt.AddHours(5)            // 加小时
dt.AddMinutes(30)         // 加分钟
dt.AddSeconds(90)         // 加秒
```

### 比较

```go
dt.Equal(other)           // 相等判断
dt.Before(other)          // 早于
dt.After(other)           // 晚于
dt.Between(start, end)    // 范围判断
dt.IsZero()               // 零值判断
```

### 时间边界

```go
dt.BeginOfMinute()   // 2024-03-15 14:30:00
dt.EndOfMinute()     // 2024-03-15 14:30:59.999...
dt.BeginOfHour()     // 2024-03-15 14:00:00
dt.EndOfHour()       // 2024-03-15 14:59:59.999...
dt.BeginOfDay()      // 2024-03-15 00:00:00
dt.EndOfDay()        // 2024-03-15 23:59:59.999...
dt.BeginOfWeek()     // 当前周的周日
dt.EndOfWeek()       // 当前周的周六
dt.BeginOfMonth()    // 2024-03-01 00:00:00
dt.EndOfMonth()      // 2024-03-31 23:59:59.999...
dt.BeginOfQuarter()  // 2024-01-01 00:00:00
dt.EndOfQuarter()    // 2024-03-31 23:59:59.999...
dt.BeginOfYear()     // 2024-01-01 00:00:00
dt.EndOfYear()       // 2024-12-31 23:59:59.999...
```

### 星期导航

```go
dt.Monday()     // 当前周的周一
dt.Tuesday()    // 当前周的周二
dt.Wednesday()  // ...
dt.Thursday()
dt.Friday()
dt.Saturday()
dt.Sunday()
```

### 转换

```go
dt.Unwrap()      // → time.Time
dt.String()      // → "2024-03-15 14:30:00"
dt.Format(layout) // 自定义格式
dt.Unix()        // Unix 秒
dt.UnixMilli()   // Unix 毫秒
dt.Sub(other)    // 两个时间的间隔
```

## Date

### 创建

```go
now := timex.NowDate()
d := timex.DateOf(time.Now())  // 去除时间部分
d, err := timex.ParseDate("2024-03-15")
```

### 方法

`Date` 提供与 `DateTime` 相同的边界和比较方法，但在日期级别操作：

```go
d.AddDays(7)
d.AddMonths(1)
d.AddYears(1)
d.BeginOfWeek()
d.EndOfMonth()
d.Monday() // ... 到 Sunday()
d.Between(start, end)
```

## Time

### 创建

```go
now := timex.NowTime()
t := timex.TimeOf(time.Now())  // 去除日期部分
t, err := timex.ParseTime("14:30:00")
```

### 方法

```go
t.AddHours(2)
t.AddMinutes(30)
t.AddSeconds(90)
t.AddMilliseconds(500)
t.Hour()
t.Minute()
t.Second()
t.BeginOfMinute()
t.EndOfHour()
t.Between(start, end)
```

## JSON 行为

```json
{
  "createdAt": "2024-03-15 14:30:00",
  "birthday": "1990-05-20",
  "startTime": "09:00:00"
}
```

没有时区后缀，没有 `T` 分隔符——干净、人类可读的格式。

## 数据库使用

三种类型可以无缝配合 Bun ORM 使用：

```go
type Event struct {
    bun.BaseModel `bun:"table:events"`
    orm.Model

    StartDate timex.Date     `json:"startDate" bun:"start_date,type:date"`
    StartTime timex.Time     `json:"startTime" bun:"start_time,type:time"`
    CreatedAt timex.DateTime `json:"createdAt" bun:"created_at,type:timestamp"`
}
```
