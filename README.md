# :alarm_clock: Cadence

This package is a simple interface in to calculating the next occurence of an
event based on cron syntax. It also supports a more human-readable syntax to
make, all behind a single interface.

```golang
import "github.com/bradhe/cadence"
```

```golang
package main

import (
    "fmt"

    "github.com/bradhe/cadence"
)

func main() {
    // Supports cron syntax with seconds
    t, _ := cadence.Next("*/5 * * * * *", time.Now()))
    fmt.Printf("Next run: %v\n", t)

    // Special cases some human readable syntax, too!
    t, _ = cadence.Next("every hour", time.Now()))
    fmt.Printf("Next run: %v\n", t)
}
```

