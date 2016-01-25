# applog -- The Application Logger

The applog library takes a UI-focused approach to logging for an application. It's primarily intended to be used for informing the application user of
things like in-progress tasks, errors, and important messages.

The library is also intended to make supporting an application easier, by
logging everything, even if it only displays a subset of information to
the user.

## Examples

``` go
package main

import (
    "time"

    "github.intel.com/hpdd/applog"
)

func main() {
    applog.StartTask("Doing some long-running process")
    time.Sleep(10 * time.Second)
    applog.EndTask()    

    name := "Fred"
    if 2 + 2 == 5 {
        applog.Fail("Reality is broken, %s!", name)
    } else {
        applog.User("Things are fine, %s!", name)
    }
}
```

Will display something like (with a spinner until complete):

    Doing some long-running process... Done.

    Things are fine, Fred!

All but the final EndTask() call may be omitted if there is a sequential
series of tasks started with StartTask().

A call to Fail() with an error object will prepend the text "ERROR: " to
the Error() string before exiting, otherwise it just prints the arguments
before exiting.
