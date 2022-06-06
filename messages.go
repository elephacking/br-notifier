package main

import (
	"encoding/json"
	"html/template"
	"strings"
	"strconv"
	"encoding/base64"

	"github.com/ashwanthkumar/slack-go-webhook"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Can't be const because need reference to variable for Slack webhook title
var (
	CONFIG  string = "New Badger"
	MESSAGE string = "Command Output"
)

type Sender interface {
	SendSlack() error
	SendEmail() error
}

func senderDispatch(status string, webhookResponse WebhookResponse, response []byte) (Sender, error) {
	if status == "CONFIG" {
		return NewConfigDetails(webhookResponse, response)
	}
	if status == "MESSAGE" {
		return NewMessageDetails(webhookResponse, response)
	}
	log.Warn("unknown status:", status)
	return nil, nil
}

// More information about events can be found here:
type WebhookResponse struct {
	Badger    string   `json:"badger"`
	Config    *ConfigDetails   `json:"badger_config"`
	Message   string `json:"badger_msg"`
}

func NewWebhookResponse(body []byte) (WebhookResponse, error) {
	var response WebhookResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return WebhookResponse{}, err
	}
	return response, nil
}

type ConfigDetails struct {
	Arch string `json:"b_arch"`
	Bld string `json:"b_bld"`
	C2 string `json:"b_c2"`
	C2_id string `json:"b_c2_id"`
	Cookie string `json:"b_cookie"`
	Hostname string `json:"b_h_name"`
	Localip string `json:"b_l_ip"`
	Process_name string `json:"b_p_name"`
	Process_id string `json:"b_pid"`
	Last_seen string `json:"b_seen"`
	User_id string `json:"b_uid"`
	Windows_version string `json:"b_wver"`
	Dead bool `json:"dead"`
	Is_pvt bool `json:"is_pvt"`
	Pipeline string `json:"pipeline"`
	Pvt_master string `json:"pvt_master"`
}

type MessageConfigDetails struct {
	Badger string
	Config ConfigDetails
}

func NewConfigDetails(response WebhookResponse, detailsRaw []byte) (MessageConfigDetails, error) {
	var details ConfigDetails
	msgDetails := MessageConfigDetails{
		Badger: response.Badger,
		Config: ConfigDetails{},
	} 
	if err := json.Unmarshal(detailsRaw, &details); err != nil {
		return msgDetails, err
	}
	msgDetails.Config = details;
	return msgDetails, nil
}

func (w MessageConfigDetails) SendSlack() error {
	orange := "#ffa500"
	red := "#f05b4f"
	attachment := slack.Attachment{}
	if w.Config.User_id[0] == '*' {
		attachment = slack.Attachment{Title: &CONFIG, Color: &red}
	} else {
		attachment = slack.Attachment{Title: &CONFIG, Color: &orange}
	}
	attachment.AddField(slack.Field{Title: "Badger ID", Value: w.Badger})
	attachment.AddField(slack.Field{Title: "Badger Hostname", Value: w.Config.Hostname})
	attachment.AddField(slack.Field{Title: "Badger User ID", Value: w.Config.User_id})
	attachment.AddField(slack.Field{Title: "Badger Windows Version", Value: w.Config.Windows_version})
	attachment.AddField(slack.Field{Title: "Badger OS Build", Value: w.Config.Bld})
	attachment.AddField(slack.Field{Title: "Badger C2", Value: w.Config.C2})
	attachment.AddField(slack.Field{Title: "Badger C2 ID", Value: w.Config.C2_id})
	attachment.AddField(slack.Field{Title: "Badger Cookie", Value: w.Config.Cookie})
	attachment.AddField(slack.Field{Title: "Badger Localip", Value: w.Config.Localip})
	attachment.AddField(slack.Field{Title: "Badger Process Name", Value: w.Config.Process_name})
	attachment.AddField(slack.Field{Title: "Badger Process ID", Value: w.Config.Process_id})
	attachment.AddField(slack.Field{Title: "Badger Last Seen", Value: w.Config.Last_seen})
	attachment.AddField(slack.Field{Title: "Badger is Dead?", Value: strconv.FormatBool(w.Config.Dead)})
	attachment.AddField(slack.Field{Title: "Badger Is Pvt?", Value: strconv.FormatBool(w.Config.Is_pvt)})
	attachment.AddField(slack.Field{Title: "Badger Pipeline", Value: w.Config.Pipeline})
	attachment.AddField(slack.Field{Title: "Badger Pvt Master", Value: w.Config.Pvt_master})
	return sendSlackAttachment(attachment)
}

func (w MessageConfigDetails) SendEmail() error {
	templateString := viper.GetString("email_config_template")
	body, err := getEmailBody(templateString, w)
	if err != nil {
		return err
	}
	return sendEmail("BRBot - Initial Connection", body)
}

type MessageDetails struct {
	Message    string
}

type MessageMessageDetails struct {
	Badger string
	Message MessageDetails
}

func NewMessageDetails(response WebhookResponse, detailsRaw []byte) (MessageMessageDetails, error) {
	var details MessageDetails
	msgDetails := MessageMessageDetails{
		Badger: response.Badger,
		Message: MessageDetails{},
	} 
	data := make([]byte, base64.StdEncoding.DecodedLen(len(detailsRaw)))
	if n, err := base64.StdEncoding.Decode(data, detailsRaw); err != nil {
		return msgDetails, err
	} else {
		details.Message = string(data[:n])
		msgDetails.Message = details
	}
	return msgDetails, nil
}

func (w MessageMessageDetails) SendSlack() error {
	white := "#ffffff"
	attachment := slack.Attachment{Title: &MESSAGE, Color: &white}
	attachment.AddField(slack.Field{Title: "Badger ID", Value: w.Badger})
	attachment.AddField(slack.Field{Title: "Command Output", Value: w.Message.Message})
	return sendSlackAttachment(attachment)
}

func (w MessageMessageDetails) SendEmail() error {
	templateString := viper.GetString("email_message_template")
	body, err := getEmailBody(templateString, w)
	if err != nil {
		return err
	}
	return sendEmail("BRBot - Command Output", body)
}

func getEmailBody(templateValue string, obj interface{}) (string, error) {
	out := new(strings.Builder)
	tpl, err := template.New("email").Parse(templateValue)
	if err != nil {
		return "", err
	}
	if err := tpl.Execute(out, obj); err != nil {
		return "", err
	}
	return out.String(), nil
}