package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Prajithp/gchat-notifier/config"
	"github.com/gorilla/mux"
	"net/http"
)

type Text struct {
	Text string `json:"text"`
}

type KeyValue struct {
	TopLabel         string `json:"topLabel"`
	Content          string `json:"content"`
	ContentMultiline bool   `json:"contentMultiline"`
}

type Widgets struct {
	KeyValue      *KeyValue `json:"keyValue,omitempty"`
	TextParagraph *Text     `json:"textParagraph,omitempty"`
}

type Sections struct {
	Widgets []Widgets `json:"widgets"`
}

type Cards struct {
	Sections []Sections `json:"sections"`
}

type SliceCard struct {
	Cards []Cards `json:"cards"`
}

type Alerts struct {
	Alerts []Alert `json:"alerts"`
}

type Alert struct {
	Annotations  map[string]interface{} `json:"annotations"`
	EndsAt       string                 `json:"sendsAt"`
	GeneratorURL string                 `json:"generatorURL"`
	Labels       map[string]interface{} `json:"labels"`
	StartsAt     string                 `json:"startsAt"`
}

var colorMap = map[string]string{
	"warning":  "#ffc107",
	"success":  "#2cbe4e",
	"critical": "#ff0000",
}

func Notification(c *config.Config, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var alert Alerts
	err := json.NewDecoder(r.Body).Decode(&alert)
	if err != nil {
		fmt.Println(err)
		respondError(w, http.StatusBadRequest, "Invalie json")
		return
	}
	defer r.Body.Close()

	var channel config.Channel
	for _, x := range c.Channels {
		if vars["channel"] == x.Name {
			channel = x
		}
	}

	if channel.Name == "" {
		respondError(w, http.StatusNotFound, "channel not found")
		return
	}

	for _, a := range alert.Alerts {
		widget := make([]Widgets, 0, len(channel.Labels))
		for _, label := range channel.Labels {
			keyvalue := KeyValue{label, a.Labels[label].(string), true}
			w := Widgets{}
			w.KeyValue = &keyvalue
			widget = append(widget, w)
		}
		Labelsection := Sections{Widgets: widget}

		severity := a.Labels["severity"].(string)
		alertname := a.Labels["alertname"].(string)
		title := fmt.Sprintf("<b>%s - <font color=\"%s\">%s</font></b>", alertname, colorMap[severity], severity)
		paragraph := Text{title}

		ww := Widgets{}
		ww.TextParagraph = &paragraph
		Headerwidget := []Widgets{ww}
		Headersection := Sections{Widgets: Headerwidget}

		card := Cards{
			Sections: []Sections{Headersection, Labelsection},
		}
		s := SliceCard{
			Cards: []Cards{card},
		}
		bodymarshal, _ := json.Marshal(&s)

		client := &http.Client{}
		req, err := http.NewRequest("POST", channel.Url, bytes.NewBuffer(bodymarshal))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			fmt.Println(err)
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
		_, err = http.Post(channel.Url, "application/json", bytes.NewBuffer(bodymarshal))
		if err != nil {
		    fmt.Printf("%s", "unable to post data to google")
		}
	}

	res := make(map[string]string)
	res["message"] = vars["Ok"]
	respondJSON(w, http.StatusOK, res)
}
