package main

import (
	"os"
	"time"

	"gopkg.in/toast.v1"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()

	app.Name = "toast"
	app.Usage = "Windows 10 toasts"
	app.Version = "v1"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Jacob Marshall",
			Email: "go-toast@jacobmarshall.co",
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "app-id, id",
			Usage: "the app identifier (used for grouping multiple toasts)",
		},
		cli.StringFlag{
			Name: "title, t",
			Usage: "the main toast title/heading",
		},
		cli.StringFlag{
			Name: "message, m",
			Usage: "the toast's main message (new lines as separator)",
		},
		cli.StringFlag{
			Name: "icon, i",
			Usage: "the app icon path (displays to the left of the toast)",
		},
		cli.StringSliceFlag{
			Name: "action, a",
			Usage: "optional action button",
		},
		cli.StringSliceFlag{
			Name: "action-type, at",
			Usage: "the type of action button",
		},
		cli.StringSliceFlag{
			Name: "action-arg, aa",
			Usage: "the action button argument",
		},
	}

	app.Action = func(c *cli.Context) error {
		appID := c.String("app-id")
		title := c.String("title")
		message := c.String("message")
		icon := c.String("icon")

		var actions []toast.Action
		actionTexts := c.StringSlice("action")
		actionTypes := c.StringSlice("action-type")
		actionArgs := c.StringSlice("action-arg")

		for index, actionLabel := range actionTexts {
			var actionType string = "protocol"
			var actionArg string
			if len(actionTypes) > index {
				actionType = actionTypes[index]
			}
			if len(actionArgs) > index {
				actionArg = actionArgs[index]
			}
			actions = append(actions, toast.Action{
				Type: actionType,
				Label: actionLabel,
				Arguments: actionArg,
			})
		}

		notification := &toast.Notification{
			AppID:   appID,
			Title:   title,
			Message: message,
			Icon:    icon,
			Actions: actions,
		}

		return notification.Push()
	}

	app.Run(os.Args)
}
