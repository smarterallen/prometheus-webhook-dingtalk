package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "strings"
    "time"
    "io/ioutil"
    "net/http"
)

const (
    DingDingUrl = "https://oapi.dingtalk.com/robot/send?access_token=access_token"
)

type Text struct {
    Content string `json:"content"`
}

type Msg struct {
    MsgType string `json:"msgtype"`
    Text    Text   `json:"text"`
}

type Alert struct {
    Status      string            `json:"status"`
    Labels      map[string]string `json:"labels"`
    Annotations map[string]string `json:"annotations"`
    StartsAt    time.Time         `json:"startsAt"`
    EndsAt      time.Time         `json:"endsAt"`
}

type Notification struct {
    Version           string            `json:"version"`
    GroupKey          string            `json:"groupKey"`
    Status            string            `json:"status"`
    Receiver          string            `json:"receiver"`
    GroupLabels       map[string]string `json:"groupLabels"`
    CommonLabels      map[string]string `json:"commonLabels"`
    CommonAnnotations map[string]string `json:"commonAnnotations"`
    ExternalURL       string            `json:"externalURL"`
    Alerts            []Alert           `json:"alerts"`
}

func Dingtalk(w http.ResponseWriter, r *http.Request) {
    b, _ := ioutil.ReadAll(r.Body)
    defer r.Body.Close()
    var notification Notification
    log.Println(string(b))
    json.Unmarshal(b, &notification)
    log.Println(notification)
    contents := []string{}
    headers := fmt.Sprintf("group: %s  status:%s", notification.CommonLabels["group"], notification.Status)
    log.Println(headers)
    contents = append(contents, headers)
    for _, each := range notification.Alerts {
        body := fmt.Sprintf("status:%s %s", each.Status, each.Annotations["summary"])
        contents = append(contents, body)
    }
    strings.Join(contents, "\n")
    msg := Msg{
        MsgType: "text",
        Text: Text{
            Content: strings.Join(contents, "\n"),
        },
    }
    msgJson, _ := json.Marshal(msg)
    req, _ := http.NewRequest("POST", DingDingUrl, bytes.NewBuffer(msgJson))
    req.Header.Add("Content-Type", "application/json")
    client := http.Client{}
    res,_ := client.Do(req)
    defer res.Body.Close()
    body, _ := ioutil.ReadAll(res.Body)
    log.Printf("%s\n", body)
    fmt.Fprint(w, "hello world\n")
}

func main() {
    http.HandleFunc("/send", Dingtalk)
    http.ListenAndServe(":8090", nil)
}
