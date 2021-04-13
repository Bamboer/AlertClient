package notification

import (
        "fmt"
        "time"
        "bytes"
        "strconv"
        "strings"
        "net/smtp"
        "encoding/base64"
        "grafana/pkg/client"
        "grafana/pkg/configer"
)

var (
        message *Message
)

func init() {
        SNS["mail"] = MSend
}

func MSend(state string, msg client.SimpleInfo, b []byte) error {
        conf := configer.ConfigParse()
        notifications  := strings.Split(conf.Notifications,",")
        notifications_cc := strings.Split(conf.Notifications_cc,",")
        notifications_bcc := strings.Split(conf.Notifications_bcc,",")
        message = &Message{From: "SVoice " + conf.SmtpServer.Username,
                To:   notifications,
                Cc:   notifications_cc,
                Bcc:  notifications_bcc,
                Attachment: Attachment{
                        WithFile:    true,
                        ContentType: "image/png",
                        Name:        "graph.png",
                },
        }
        c := TimeController(conf.SmtpServer.StartTime, conf.SmtpServer.EndTime)
        if control := <-c; !control {
                err := fmt.Errorf("scheduler time mail send closed.")
                return err
        }

        mclient, _ := NewMail(conf.SmtpServer.Username, conf.SmtpServer.Password, conf.SmtpServer.SmtpAddress, conf.SmtpServer.Port)
        if err := mclient.Send(state, msg, b); err != nil {
                info.Println(err)
                return err
        }
        info.Println("mail sent.")
        return nil
}

type mailor struct {
        User     string
        Password string
        Host     string
        Port     string
        Auth     smtp.Auth
}

type Message struct {
        From        string
        To          []string
        Cc          []string
        Bcc         []string
        Subject     string
        Body        string
        ContentType string
        Attachment  Attachment
}

type Attachment struct {
        Name        string
        ContentType string
        WithFile    bool
}

func NewMail(username, password, smtpServer, port string) (*mailor, error) {
        auth := smtp.PlainAuth("", username, password, smtpServer)
        return &mailor{
                User:     username,
                Password: password,
                Host:     smtpServer,
                Port:     port,
                Auth:     auth,
        }, nil
}

func (m *mailor) Send(state string, msg client.SimpleInfo, b []byte) error {
        var imgsrc string
        buffer := bytes.NewBuffer(nil)
        boundary := "BamboerBoundary"
        Header := make(map[string]string)
        Header["From"] = message.From
        Header["To"] = strings.Join(message.To, ";")
        Header["Cc"] = strings.Join(message.Cc, ";")
        Header["Bcc"] = strings.Join(message.Bcc, ";")
        Header["Subject"] = msg.Name
        if state == "ok" {
                Header["Subject"] = msg.Name + " Recovery !"
        }
        Header["Content-Type"] = "multipart/related;boundary=" + boundary
        Header["Mime-Version"] = "1.0"
        Header["Date"] = time.Now().UTC().String()

        m.WriteHeader(buffer, Header)

        if message.Attachment.WithFile {
                attachment := "\r\n--" + boundary + "\r\n"
                attachment += "Content-Transfer-Encoding:base64\r\n"
                //                attachment += "Content-Disposition:attachment\r\n"
                attachment += "Content-Type:" + message.Attachment.ContentType + ";name=\"" + message.Attachment.Name + "\"\r\n"
                attachment += "Content-ID: <" + message.Attachment.Name + "> \r\n\r\n"
                imgsrc = "<p><img src=\"cid:" + message.Attachment.Name + "\"></p>"
                buffer.WriteString(attachment)
                defer func() {
                        if err := recover(); err != nil {
                                info.Println("Error: ", err)
                        }
                }()
                m.WriteFile(buffer, b)
        }

        body := "\r\n--" + boundary + "\r\n"
        body += "Content-Type: text/html; charset=UTF-8 \r\n"
        body += renderMessage(state, imgsrc, msg)
        buffer.WriteString(body)

        buffer.WriteString("\r\n--" + boundary + "--")
        if err := smtp.SendMail(m.Host+m.Port, m.Auth, m.User, message.To, buffer.Bytes());err !=nil{
            return err
        }
        return nil
}

func TimeController(start, end string) <-chan bool {
        c := make(chan bool,1)
        defer close(c)
        var start_time, end_time time.Time
        starthour := strings.Split(start, ":")[0]
        startmin := strings.Split(start, ":")[1]

        endhour:= strings.Split(end, ":")[0]
        endmin := strings.Split(end, ":")[1]

        startHour,_ := strconv.Atoi(starthour)
        startMin,_ := strconv.Atoi(startmin)
        endHour,_ := strconv.Atoi(endhour)
        endMin,_ := strconv.Atoi(endmin)
        if startHour > 24 || startHour < 0 || endHour > 24 || endHour < 0 || start == end || startMin > 60 || startMin < 0 || endMin > 60 || endMin < 0 {
                info.Println("smtp server send time error, that will ignored.")
                c <- true
                return c
        }
        now_time := time.Now()
        start_time = time.Date(now_time.Year(), now_time.Month(), now_time.Day(), startHour, startMin, 0, 0, now_time.Location())
        end_time = time.Date(now_time.Year(), now_time.Month(), now_time.Day(), endHour, endMin, 0, 0, now_time.Location())
        if startHour > endHour {
                end_time = time.Date(now_time.Year(), now_time.Month(), now_time.Day()+1, endHour, endMin, 0, 0, now_time.Location())
        }

        switch {
        case now_time.Before(start_time):
                info.Println("smtp send time not till.")
                c <- false
                return c
        case now_time.After(end_time):
                info.Println("smtp send time not passed.")
                c <- false
                return c
        case now_time.Before(end_time):
                info.Println("smtp send time.")
                c <- true
                return c
        }
       return c
}

func renderMessage(state, imgsrc string, msg client.SimpleInfo) string {
        var template string
        if state == "alerting" {
                template = `
<html>
       <body>
                %s
                <p>Alarm: %s</p>
                <p>Metric: %s</p>
                <p>Value: %v</p>
                <p>Detail: It's need %v second at least to recovery.</p>
                <p>Time: %v </p>
      </body>
</html>
`
        }
        if state == "ok" {
                template = `
<html>
       <body>
                %s
                <p>Alarm: %s Recovery !</p>
                <p>Metric: %s</p>
                <p>Value: %v</p>
                <p>Detail: It's need %v second at least to recovery.</p>
                <p>Time: %v </p>
       </body>
</html>
`
        }
        var content = fmt.Sprintf(template, imgsrc, msg.Name, msg.AlertMetrics, msg.AlertValues, strconv.Itoa(msg.Frequency), time.Now().UTC().String())
        return content
}

func (m *mailor) WriteHeader(buffer *bytes.Buffer, Header map[string]string) string {
        header := ""
        for key, value := range Header {
                header += key + ":" + value + "\r\n"
        }
        header += "\r\n"
        buffer.WriteString(header)
        return header
}

func (m *mailor) WriteFile(buffer *bytes.Buffer, b []byte) {
        payload := make([]byte, base64.StdEncoding.EncodedLen(len(b)))
        base64.StdEncoding.Encode(payload, b)
        buffer.WriteString("\r\n")
        for index, line := 0, len(payload); index < line; index++ {
                buffer.WriteByte(payload[index])
                if (index+1)%76 == 0 {
                        buffer.WriteString("\r\n")
                }
        }
}
