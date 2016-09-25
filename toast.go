package toast

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/nu7hatch/gouuid"
)

var toastTemplate *template.Template

var (
	ErrorInvalidAudio    error = errors.New("toast: invalid audio")
	ErrorInvalidDuration       = errors.New("toast: invalid duration")
)

type toastAudio string

const (
	Default        toastAudio = "ms-winsoundevent:Notification.Default"
	IM                        = "ms-winsoundevent:Notification.IM"
	Mail                      = "ms-winsoundevent:Notification.Mail"
	Reminder                  = "ms-winsoundevent:Notification.Reminder"
	SMS                       = "ms-winsoundevent:Notification.SMS"
	LoopingAlarm              = "ms-winsoundevent:Notification.Looping.Alarm"
	LoopingAlarm2             = "ms-winsoundevent:Notification.Looping.Alarm2"
	LoopingAlarm3             = "ms-winsoundevent:Notification.Looping.Alarm3"
	LoopingAlarm4             = "ms-winsoundevent:Notification.Looping.Alarm4"
	LoopingAlarm5             = "ms-winsoundevent:Notification.Looping.Alarm5"
	LoopingAlarm6             = "ms-winsoundevent:Notification.Looping.Alarm6"
	LoopingAlarm7             = "ms-winsoundevent:Notification.Looping.Alarm7"
	LoopingAlarm8             = "ms-winsoundevent:Notification.Looping.Alarm8"
	LoopingAlarm9             = "ms-winsoundevent:Notification.Looping.Alarm9"
	LoopingAlarm10            = "ms-winsoundevent:Notification.Looping.Alarm10"
	LoopingCall               = "ms-winsoundevent:Notification.Looping.Call"
	LoopingCall2              = "ms-winsoundevent:Notification.Looping.Call2"
	LoopingCall3              = "ms-winsoundevent:Notification.Looping.Call3"
	LoopingCall4              = "ms-winsoundevent:Notification.Looping.Call4"
	LoopingCall5              = "ms-winsoundevent:Notification.Looping.Call5"
	LoopingCall6              = "ms-winsoundevent:Notification.Looping.Call6"
	LoopingCall7              = "ms-winsoundevent:Notification.Looping.Call7"
	LoopingCall8              = "ms-winsoundevent:Notification.Looping.Call8"
	LoopingCall9              = "ms-winsoundevent:Notification.Looping.Call9"
	LoopingCall10             = "ms-winsoundevent:Notification.Looping.Call10"
	Silent                    = ""
)

type toastDuration string

const (
	Short toastDuration = "short"
	Long                = "long"
)

func init() {
	toastTemplate = template.New("toast")
	toastTemplate.Parse(`
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
[Windows.UI.Notifications.ToastNotification, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
[Windows.Data.Xml.Dom.XmlDocument, Windows.Data.Xml.Dom.XmlDocument, ContentType = WindowsRuntime] | Out-Null

$APP_ID = '{{if .AppID}}{{.AppID}}{{else}}Windows App{{end}}'

$template = @"
<toast activationType="{{.ActivationType}}" launch="{{.ActivationArguments}}" duration="{{.Duration}}">
    <visual>
        <binding template="ToastGeneric">
            {{if .Icon}}
            <image placement="appLogoOverride" src="{{.Icon}}" />
            {{end}}
            {{if .Title}}
            <text><![CDATA[{{.Title}}]]></text>
            {{end}}
            {{if .Message}}
            <text><![CDATA[{{.Message}}]]></text>
            {{end}}
        </binding>
    </visual>
    {{if .Audio}}
	<audio src="{{.Audio}}" loop="{{.Loop}}" />
	{{else}}
	<audio silent="true" />
	{{end}}
    {{if .Actions}}
    <actions>
        {{range .Actions}}
        <action activationType="{{.Type}}" content="{{.Label}}" arguments="{{.Arguments}}" />
        {{end}}
    </actions>
    {{end}}
</toast>
"@

$xml = New-Object Windows.Data.Xml.Dom.XmlDocument
$xml.LoadXml($template)
$toast = New-Object Windows.UI.Notifications.ToastNotification $xml
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier($APP_ID).Show($toast)
    `)
}

type Notification struct {
	// The name of your app. This value shows up in Windows 10's Action Centre, so make it
	// something readable for your users. It can contain spaces, however special characters
	// (eg. é) are not supported.
	AppID string

	// The main title/heading for the toast notification.
	Title string

	// The single/multi line message to display for the toast notification.
	Message string

	// An optional path to an image on the OS to display to the left of the title & message.
	Icon string

	// The type of notification level action (like toast.Action)
	ActivationType string

	// The activation/action arguments (invoked when the user clicks the notification)
	ActivationArguments string

	// Optional action buttons to display below the notification title & message.
	Actions []Action

	// The audio to play when displaying the toast
	Audio toastAudio

	// Whether to loop the audio (default false)
	Loop bool

	// How long the toast should show up for (short/long)
	Duration toastDuration
}

// Defines an actionable button.
// See https://msdn.microsoft.com/en-us/windows/uwp/controls-and-patterns/tiles-and-notifications-adaptive-interactive-toasts for more info.
//
// Only protocol type action buttons are actually useful, as there's no way of receiving feedback from the
// user's choice. Examples of protocol type action buttons include: "bingmaps:?q=sushi" to open up Windows 10's
// maps app with a pre-populated search field set to "sushi".
//
//     toast.Action{"protocol", "Open Maps", "bingmaps:?q=sushi"}
type Action struct {
	Type      string
	Label     string
	Arguments string
}

func (n *Notification) buildXML() (string, error) {
	var out bytes.Buffer
	err := toastTemplate.Execute(&out, n)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// Builds the Windows PowerShell script & invokes it, causing the toast to display.
//
// Note: Running the PowerShell script is by far the slowest process here, and can take a few
// seconds in some cases.
//
//     notification := toast.Notification{
//         AppID: "Example App",
//         Title: "My notification",
//         Message: "Some message about how important something is...",
//         Icon: "go.png",
//         Actions: []toast.Action{
//             {"protocol", "I'm a button", ""},
//             {"protocol", "Me too!", ""},
//         },
//     }
//     err := notification.Push()
//     if err != nil {
//         log.Fatalln(err)
//     }
func (n *Notification) Push() error {
	xml, _ := n.buildXML()
	return invokeTemporaryScript(xml)
}

func Audio(name string) (toastAudio, error) {
	switch strings.ToLower(name) {
	case "default":
		return Default, nil
	case "im":
		return IM, nil
	case "mail":
		return Mail, nil
	case "reminder":
		return Reminder, nil
	case "sms":
		return SMS, nil
	case "loopingalarm":
		return LoopingAlarm, nil
	case "loopingalarm2":
		return LoopingAlarm2, nil
	case "loopingalarm3":
		return LoopingAlarm3, nil
	case "loopingalarm4":
		return LoopingAlarm4, nil
	case "loopingalarm5":
		return LoopingAlarm5, nil
	case "loopingalarm6":
		return LoopingAlarm6, nil
	case "loopingalarm7":
		return LoopingAlarm7, nil
	case "loopingalarm8":
		return LoopingAlarm8, nil
	case "loopingalarm9":
		return LoopingAlarm9, nil
	case "loopingalarm10":
		return LoopingAlarm10, nil
	case "loopingcall":
		return LoopingCall, nil
	case "loopingcall2":
		return LoopingCall2, nil
	case "loopingcall3":
		return LoopingCall3, nil
	case "loopingcall4":
		return LoopingCall4, nil
	case "loopingcall5":
		return LoopingCall5, nil
	case "loopingcall6":
		return LoopingCall6, nil
	case "loopingcall7":
		return LoopingCall7, nil
	case "loopingcall8":
		return LoopingCall8, nil
	case "loopingcall9":
		return LoopingCall9, nil
	case "loopingcall10":
		return LoopingCall10, nil
	case "silent":
		return Silent, nil
	default:
		return Default, ErrorInvalidAudio
	}
}

func Duration(name string) (toastDuration, error) {
	switch strings.ToLower(name) {
	case "short":
		return Short, nil
	case "long":
		return Long, nil
	default:
		return Short, ErrorInvalidDuration
	}
}

func invokeTemporaryScript(content string) error {
	id, _ := uuid.NewV4()
	file := filepath.Join(os.TempDir(), id.String()+".ps1")
	defer os.Remove(file)
	err := ioutil.WriteFile(file, []byte(content), 0600)
	if err != nil {
		return err
	}
	if err = exec.Command("PowerShell", "-ExecutionPolicy", "Bypass", "-File", file).Run(); err != nil {
		return err
	}
	return nil
}
