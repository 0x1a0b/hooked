package discordSender

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/0x1a0b/hooked/config"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func New(secret string) (s *Sender) {

	s = &Sender{
		webhookSecret: secret,
	}

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	s.client = &http.Client{Timeout: time.Second * 20, Transport: tr}

    s.logger = logrus.New()
    s.logger.SetLevel(config.GetLogLevel())

	return
}

type Sender struct {
	webhookSecret string
	client *http.Client
	logger *logrus.Logger
}

func (s *Sender) Send(h Hook) (err error) {

	s.logger.Debugf("Start sending hook")

	if s.webhookSecret == "" {
		s.logger.Errorf("webhook secret is not configured")
		return errors.New("webhook secret is not configured")
	}

	var object []byte
    object, err = json.Marshal(h)
    if err != nil {
    	s.logger.Errorf("error creating json from object: %v", err)
    	return
	}

	var result *http.Response
	result, err = s.client.Post(s.webhookSecret, "application/json", bytes.NewBuffer(object))
	if err != nil {
		s.logger.Errorf("http client termianted abnormal: %v", err)
		return err
	}
	defer result.Body.Close()

	if result.StatusCode != 204 {
		text := "Abnormal Statuscode after wending webhook: " + strconv.Itoa(result.StatusCode)
		s.logger.Errorf("discord response error: %v", text)
		if s.logger.IsLevelEnabled(logrus.DebugLevel) == true {
			var body []byte
			body, err = ioutil.ReadAll(result.Body)
			if err != nil {
				s.logger.Errorf("error serializing body: %v", err)
			} else {
				s.logger.Debugf("response body was: %v", string(body))
			}
		}
		return errors.New(text)
	}

	s.logger.Debugf("Ended sending hook")

	return
}

type Hook struct {
	Content string `json:"content,omitempty"`
	Username string `json:"username,omitempty"`
	AvatarUrl string `json:"avatar_url,omitempty"`
	Tts bool `json:"tts,omitempty"`
	Embeds []Embed `json:"embeds,omitempty"`
}
type Embed struct {
	Title string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Color int `json:"color,omitempty"`
	Url string `json:"url,omitempty"`
	Fields []Field `json:"fields,omitempty"`
	Thumbnail Thumbnail `json:"thumbnail,omitempty"`
	Author Author `json:"author,omitempty"`
	Footer Footer `json:"footer,omitempty"`
}
type Field struct {
	Name string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
	Inline bool `json:"inline,omitempty"`
}
type Thumbnail struct {
    Url string `json:"url,omitempty"`
    Height int `json:"height,omitempty"`
    Width int `json:"width,omitempty"`
}
type Footer struct {
	Text string `json:"text,omitempty"`
	IconUrl string `json:"icon_url,omitempty"`
}
type Author struct {
	Name string `json:"name,omitempty"`
	Url string `json:"url,omitempty"`
	IconUrl string `json:"icon_url,omitempty"`
}
type Provider struct {
	Name string `json:"name,omitempty" example:"ACME Webhooks"`
	Url string `json:"url,omitempty" example:"https://en.wikipedia.org/wiki/Acme_Corporation"`
}
