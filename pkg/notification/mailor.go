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
        message = &Message{from: conf.SmtpServer.Username,
                to:   notifications,
                cc:   notifications_cc,
                bcc:  notifications_bcc,
                attachment: Attachment{
                        WithFile:    true,
                        ContentType: "image/png",
                        Name:        "graph.png",
                },
        }
        c := timeController(conf.SmtpServer.StartTime, conf.SmtpServer.EndTime)
        if control := <-c; !control {
                info.Println("scheduler time mail send closed.")
                return nil
        }
        mclient, _ := NewMail(conf.SmtpServer.Username, conf.SmtpServer.Password, conf.SmtpServer.SmtpAddress, conf.SmtpServer.Port)
        if err := mclient.Send(state, msg, b); err != nil {
                info.Println(err)
                return err
        }
        return nil
}

type mailor struct {
        user     string
        password string
        host     string
        port     string
        auth     smtp.Auth
}

type Message struct {
        from        string
        to          []string
        cc          []string
        bcc         []string
        subject     string
        body        string
        contentType string
        attachment  Attachment
}

type Attachment struct {
        Name        string
        ContentType string
        WithFile    bool
}

func NewMail(username, password, smtpServer, port string) (*mailor, error) {
        auth := smtp.PlainAuth("", username, password, smtpServer)
        return &mailor{
                user:     username,
                password: password,
                host:     smtpServer,
                port:     port,
                auth:     auth,
        }, nil
}

func (m *mailor) Send(state string, msg client.SimpleInfo, b []byte) error {
        var imgsrc string
        buffer := bytes.NewBuffer(nil)
        boundary := "BamboerBoundary"
        Header := make(map[string]string)
        Header["From"] = message.from
        Header["To"] = strings.Join(message.to, ";")
        Header["Cc"] = strings.Join(message.cc, ";")
        Header["Bcc"] = strings.Join(message.bcc, ";")
        Header["Subject"] = msg.Name
        if state == "ok" {
                Header["Subject"] = msg.Name + " Recovery !"
        }
        Header["Content-Type"] = "multipart/related;boundary=" + boundary
        Header["Mime-Version"] = "1.0"
        Header["Date"] = time.Now().UTC().String()

        m.writeHeader(buffer, Header)

        if message.attachment.WithFile {
                attachment := "\r\n--" + boundary + "\r\n"
                attachment += "Content-Transfer-Encoding:base64\r\n"
                //                attachment += "Content-Disposition:attachment\r\n"
                attachment += "Content-Type:" + message.attachment.ContentType + ";name=\"" + message.attachment.Name + "\"\r\n"
                attachment += "Content-ID: <" + message.attachment.Name + "> \r\n\r\n"
                imgsrc = "<p><img src=\"cid:" + message.attachment.Name + "></p><br>\r\n\t\t\t"
                buffer.WriteString(attachment)
                defer func() {
                        if err := recover(); err != nil {
                                info.Println("Error: ", err)
                        }
                }()
                m.writeFile(buffer, b)
        }

        body := "\r\n--" + boundary + "\r\n"
        body += "Content-Type: text/html; charset=UTF-8 \r\n"
        body += renderMessage(state, imgsrc, msg)
        buffer.WriteString(body)

        buffer.WriteString("\r\n--" + boundary + "--")
        smtp.SendMail(m.host+m.port, m.auth, m.user, message.to, buffer.Bytes())
        return nil
}

func timeController(start, end string) <-chan bool {
        c := make(chan bool)
        defer close(c)
        var start_time, end_time time.Time
        starthour := strings.Split(start, ":")
         startmin := strings.Split(start, ":")

        endhour, endmin := strings.Split(end, ":")
        endhour, endmin := strings.Split(end, ":")

        startHour := strconv.Atoi(starthour)
        startMin := strconv.Atoi(startmin)
        endHour := strconv.Atoi(endhour)
        endMin := strconv.Atoi(endmin)
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
                c <- false
                return c
        case now_time.After(end_time):
                c <- false
                return c
        case now_time.Before(end_time):
                c <- true
                return c
        }
}

func renderMessage(state, imgsrc string, msg client.SimpleInfo) string {
        var template string
        if state == "alerting" {
                template = `
<html>
        <body>
                %s
                <p>Alarm: %s</p><br>
                <p>Metric: %s</p><br>
                <p>Value: %s</p><br>
                <p>Detail: %s</p><br>
        </body>
</html>
`
        }
        if state == "ok" {
                template = `
<html>
        <body>
                %s
                <p>Alarm: %s Recovery !</p><br>
                <p>Metric: %s</p><br>
                <p>Value: %s</p><br>
                <p>Detail: It\'s least need %s second recovery.</p><br>
                <p>Time: %s </p><br>
        </body>
</html>
`
        }
        var content = fmt.Sprintf(template, imgsrc, msg.Name, msg.AlertMetrics, msg.AlertValues, strconv.Itoa(msg.Frequency), time.Now().UTC().String())
        return content
}

func (m *mailor) writeHeader(buffer *bytes.Buffer, Header map[string]string) string {
        header := ""
        for key, value := range Header {
                header += key + ":" + value + "\r\n"
        }
        header += "\r\n"
        buffer.WriteString(header)
        return header
}

func (m *mailor) writeFile(buffer *bytes.Buffer, b []byte) {
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
