package main

import (
	"flag"
	"log"
	"strings"

	"github.com/go-toast/toast"
)

func main() {
	var appId, title, message, icon, activationType, activationArg, action, actionType, actionArg, audio, duration string
	var loop bool

	flag.StringVar(&appId, "app-id", "", "the app identifier (used for grouping multiple toasts)")
	flag.StringVar(&title, "title", "", "the main toast title/heading")
	flag.StringVar(&message, "message", "", "the toast's main message (new lines as separator)")
	flag.StringVar(&icon, "icon", "", "the app icon path (displays to the left of the toast)")
	flag.StringVar(&activationType, "activation-type", toast.ActionProtocol, "the type of action to invoke when the user clicks the toast")
	flag.StringVar(&activationArg, "activation-arg", "", "the activation argument")
	flag.StringVar(&action, "action", "", "optional action button")
	flag.StringVar(&actionType, "action-type", "", "the type of action button")
	flag.StringVar(&actionArg, "action-arg", "", "the action button argument")
	flag.StringVar(&audio, "audio", string(toast.Silent), "which kind of audio should be played")
	flag.BoolVar(&loop, "loop", false, "whether to loop the audio")
	flag.StringVar(&duration, "duration", string(toast.Short), "how long the toast should display for")
	flag.Parse()

	a, _ := toast.ToAudio(audio)
	d, _ := toast.ToDuration(duration)
	var actions []toast.Action

	xs := strings.Split(action, ",")
	ys := strings.Split(actionType, ",")
	zs := strings.Split(actionArg, ",")
	for i, label := range xs {
		var atype, aarg string
		if len(ys) > i {
			atype = ys[i]
		} else {
			atype = toast.ActionProtocol
		}
		if len(zs) > i {
			aarg = zs[i]
		}
		actions = append(actions, toast.Action{
			Type:      atype,
			Label:     label,
			Arguments: aarg,
		})
	}

	notification := &toast.Notification{
		AppID:               appId,
		Title:               title,
		Message:             message,
		Icon:                icon,
		Actions:             actions,
		ActivationType:      activationType,
		ActivationArguments: activationArg,
		Audio:               a,
		Loop:                loop,
		Duration:            d,
	}

	if err := notification.Push(); err != nil {
		log.Fatalln(err)
	}
}
