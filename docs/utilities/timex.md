---
sidebar_position: 2
---

# Timex

The `timex` package provides three custom time types — `DateTime`, `Date`, and `Time` — that serve as drop-in replacements for `time.Time` with built-in JSON serialization, database scanning, and rich manipulation methods.

These types are used throughout the framework, including in all audit model fields (`CreatedAt`, `UpdatedAt`).

## Types Overview

| Type | Format | Go Layout | Example |
| --- | --- | --- | --- |
| `timex.DateTime` | `YYYY-MM-DD HH:mm:ss` | `time.DateTime` | `"2024-03-15 14:30:00"` |
| `timex.Date` | `YYYY-MM-DD` | `time.DateOnly` | `"2024-03-15"` |
| `timex.Time` | `HH:mm:ss` | `time.TimeOnly` | `"14:30:00"` |

All three types implement:
- `json.Marshaler` / `json.Unmarshaler` — clean JSON format without timezone
- `sql.Scanner` / `driver.Valuer` — database compatibility
- `encoding.TextMarshaler` / `encoding.TextUnmarshaler`

All three types share the same conversion (`Unwrap`, `Format`, `String`),
wire/database integration (`MarshalJSON`, `UnmarshalJSON`, `MarshalText`,
`UnmarshalText`, `Scan`, `Value`), and comparison (`Between`, an open interval)
method families, plus type-specific timestamp helpers (`UnixMilli`,
`UnixMicro`, `UnixNano` on `DateTime`; `ToDuration` on `Time`). JSON uses
plain layouts with no `T` separator.

## DateTime

### Creating

```go
// Current time
now := timex.Now()

// From time.Time
dt := timex.Of(time.Now())

// From string
dt, err := timex.Parse("2024-03-15 14:30:00")

// With custom format
dt, err := timex.Parse("15/03/2024 14:30", "02/01/2006 15:04")

// From Unix timestamp
dt := timex.FromUnix(1710510600, 0)
dt := timex.FromUnixMilli(1710510600000)
dt := timex.FromUnixMicro(1710510600000000)
```

The `DateTime` constructors are `Now`, `Of`, `Parse`, `FromUnix`, `FromUnixMilli`,
and `FromUnixMicro`.

The lenient `Parse*` entry points try the given layout first, then fall back
to common formats (RFC3339, ISO-8601, …). Both paths interpret zone-less input
in the **local** timezone (v0.38 fix — the fallback used to assume UTC, which
could shift a date-only string by up to a day); inputs carrying an explicit
offset keep it on either path.

### Accessing Components

```go
dt.Year()       // 2024
dt.Month()      // time.March
dt.Day()        // 15
dt.Hour()       // 14
dt.Minute()     // 30
dt.Second()     // 0
dt.Weekday()    // time.Friday
dt.YearDay()    // 75
dt.Location()   // *time.Location
dt.Nanosecond() // nanosecond component
```

### Arithmetic

```go
dt.Add(2 * time.Hour)    // Add duration
dt.AddDate(1, 2, 3)      // Add years, months, days
dt.AddDays(7)             // Add days
dt.AddMonths(3)           // Add months
dt.AddYears(1)            // Add years
dt.AddHours(5)            // Add hours
dt.AddMinutes(30)         // Add minutes
dt.AddSeconds(90)         // Add seconds
```

### Comparison

```go
dt.Equal(other)           // Equality
dt.Before(other)          // Before check
dt.After(other)           // After check
dt.Between(start, end)    // Open-interval range check: start < dt < end
dt.IsZero()               // Zero value check
```

### Time Boundaries

```go
dt.BeginOfMinute()   // 2024-03-15 14:30:00
dt.EndOfMinute()     // 2024-03-15 14:30:59.999...
dt.BeginOfHour()     // 2024-03-15 14:00:00
dt.EndOfHour()       // 2024-03-15 14:59:59.999...
dt.BeginOfDay()      // 2024-03-15 00:00:00
dt.EndOfDay()        // 2024-03-15 23:59:59.999...
dt.BeginOfWeek()     // Sunday of current week
dt.EndOfWeek()       // Saturday of current week
dt.BeginOfMonth()    // 2024-03-01 00:00:00
dt.EndOfMonth()      // 2024-03-31 23:59:59.999...
dt.BeginOfQuarter()  // 2024-01-01 00:00:00
dt.EndOfQuarter()    // 2024-03-31 23:59:59.999...
dt.BeginOfYear()     // 2024-01-01 00:00:00
dt.EndOfYear()       // 2024-12-31 23:59:59.999...
```

### Weekday Navigation

```go
dt.Monday()     // Monday of current week
dt.Tuesday()    // Tuesday of current week
dt.Wednesday()  // ...
dt.Thursday()
dt.Friday()
dt.Saturday()
dt.Sunday()
```

### Conversion

```go
dt.Unwrap()      // → time.Time
dt.String()      // → "2024-03-15 14:30:00"
dt.Format(layout) // Custom format
dt.Unix()        // Unix seconds
dt.UnixMilli()   // Unix milliseconds
dt.UnixMicro()   // Unix microseconds
dt.UnixNano()    // Unix nanoseconds
dt.Since()       // Duration since dt
dt.Until()       // Duration until dt
dt.Sub(other)    // Duration between
```

## Date

### Creating

```go
now := timex.NowDate()
d := timex.DateOf(time.Now())  // Strips time components
d, err := timex.ParseDate("2024-03-15")
```

The `Date` constructors are `NowDate`, `DateOf`, and `ParseDate`.

### Methods

`Date` offers the same boundary and comparison methods as `DateTime`, but operates on date-level granularity:

```go
d.AddDays(7)
d.AddMonths(1)
d.AddYears(1)
d.BeginOfWeek()
d.EndOfMonth()
d.Monday() // ... through Sunday()
d.Between(start, end)
d.Location()
```

## Time

### Creating

```go
now := timex.NowTime()
t := timex.TimeOf(time.Now())  // Strips date components
t, err := timex.ParseTime("14:30:00")
```

The `Time` constructors are `NowTime`, `TimeOf`, and `ParseTime`.

### Methods

```go
t.AddHours(2)
t.AddMinutes(30)
t.AddSeconds(90)
t.AddMilliseconds(500)
t.AddMicroseconds(1000)
t.AddNanoseconds(1000000)
t.Hour()
t.Minute()
t.Second()
t.Nanosecond()
t.ToDuration()
t.BeginOfMinute()
t.EndOfHour()
t.Between(start, end)
```

`Between` uses an open interval for all three types: values equal to `start` or
`end` return `false`.

## JSON Behavior

```json
{
  "createdAt": "2024-03-15 14:30:00",
  "birthday": "1990-05-20",
  "startTime": "09:00:00"
}
```

No timezone suffix, no `T` separator — clean, human-readable formats.

The concrete methods are `MarshalJSON`, `UnmarshalJSON`, `MarshalText`, and
`UnmarshalText`; database integration uses `Scan` and `Value`.

## Error Sentinels

| Error | Meaning |
| --- | --- |
| `ErrInvalidDateTimeFormat` | `DateTime` parsing or JSON/text decoding received an invalid format |
| `ErrInvalidDateFormat` | `Date` parsing or JSON/text decoding received an invalid format |
| `ErrInvalidTimeFormat` | `Time` parsing or JSON/text decoding received an invalid format |
| `ErrFailedScan` | database `Scan` received an invalid value |
| `ErrUnsupportedDestType` | scan destination type is unsupported |

## Database Usage

All three types work seamlessly with Bun ORM:

```go
type Event struct {
    bun.BaseModel `bun:"table:events"`
    orm.Model

    StartDate timex.Date     `json:"startDate" bun:"start_date,type:date"`
    StartTime timex.Time     `json:"startTime" bun:"start_time,type:time"`
    CreatedAt timex.DateTime `json:"createdAt" bun:"created_at,type:timestamp"`
}
```
