package notify

import (
	"encoding/json"
	"fmt"
)

const (
	// %s - organization id
	// %s - flow id
	// %s - access token of the developer
	flowdockApiUrl          = "https://api.flowdock.com/messages?access_token=%s"
	flowdockStartedMessage  = "*Building* %s, commit <%s|%s>, author %s"
	flowdockSuccessMessage  = "*Success* %s, commit <%s|%s>, author %s"
	flowdockFailureMessage  = "*Failed* %s, commit <%s|%s>, author %s"
	flowdockMessageType		= "message"
)

type Flowdock struct {
	FlowToken		string `yaml:"flow_token,omitempty"`
	Started			bool   `yaml:"on_started,omitempty"`
	Success			bool   `yaml:"on_success,omitempty"`
	Failure			bool   `yaml:"on_failure,omitempty"`
}

func (f *Flowdock) Send(context *Context) error {
	switch {
	case context.Commit.Status == "Started" && f.Started:
		return f.sendStarted(context)
	case context.Commit.Status == "Success" && f.Success:
		return f.sendSuccess(context)
	case context.Commit.Status == "Failure" && f.Failure:
		return f.sendFailure(context)
	}

	return nil
}

func getFlowdockMessage(context *Context, message string) string {
	url := getBuildUrl(context)
	return fmt.Sprintf(message, context.Repo.Name, url, context.Commit.HashShort(), context.Commit.Author)
}

func (f *Flowdock) url() string {
	return fmt.Sprintf(flowdockApiUrl, f.FlowToken)
}

func (f *Flowdock) sendStarted(context *Context) error {
	return f.send(getFlowdockMessage(context, flowdockStartedMessage))
}

func (f *Flowdock) sendSuccess(context *Context) error {
	return f.send(getFlowdockMessage(context, flowdockSuccessMessage))
}

func (f *Flowdock) sendFailure(context *Context) error {
	return f.send(getFlowdockMessage(context, flowdockFailureMessage))
}

// helper function to send HTTP requests
func (f *Flowdock) send(message string) error {
	// data will get posted in this format
	data := struct {
		Event    string `json:"event"`
		Content  string `json:"content"`
	}{flowdockMessageType, message}

	// data json encoded
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	go sendJson(f.url(), payload)

	return nil
}
