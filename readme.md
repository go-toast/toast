# Toast

A go package for Windows 10 toast notifications.

As seen in [jacobmarshall/pokevision-cli](https://github.com/jacobmarshall/pokevision-cli).

![Toast](./screenshot-toast.png)

![Action centre](./screenshot-action-centre.png)

## Example

```go
package main

import (
    "log"

    "gopkg.in/toast.v1"
)

func main() {
    notification := toast.Notification{
        AppID: "Example App",
        Title: "My notification",
        Message: "Some message about how important something is...",
        Icon: "go.png",
        Actions: []toast.Action{
            {"protocol", "I'm a button", ""},
            {"protocol", "Me too!", ""},
        },
    }
    err := notification.Push()
    if err != nil {
        log.Fatalln(err)
    }
}
```
