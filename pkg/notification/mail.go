package notification

import (
        "fmt"
        "bytes"
        "encoding/base64"
        "io/ioutil"
        "net/smtp"
        "strings"
        "time"
)

func init(){
  SNS["mail"] = MSend
  conf = configer.ConfigParse()
  DEFAULTMAIL,_ = NewMail(conf.SmtpServer.Username,conf.SmtpServer.Password,conf.SmtpServer.SmtpAddress,conf.Port)
  message    = Message{"from": conf.SmtpServer.Username,
                   "to":  conf.Notifications,
                   "cc":  conf.Notifications_cc,
                   "bcc":  conf.Notifications_bcc,
           }
}

func MSend(state string,msg client.SimpleInfo,b []byte)error{
  conf = configer.ConfigParse()
  mclient,_= NewMail(conf.SmtpServer.Username,conf.SmtpServer.Password,conf.SmtpServer.SmtpAddress,conf.Port)
  if err :=mclient.Send(state,msg,b);err != nil{
     info.Println(err)
     return err
  }
  return nil
}

var (
conf        *(configer.Obj)
DEFAULTMAIL *mail
message     *Message
)

type mail struct {
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

func NewMail(username, password, smtpServer, port string) (*mail, error) {
        auth := smtp.PlainAuth("", username, password, smtpServer)
        return &mail{
                user:     username,
                password: password,
                host:     smtpServer,
                port:     port,
                auth:     auth,
        }, nil
}

func (m *mail) Send(state string,msg client.SimpleInfo,b []byte) error {
        var imgsrc string
        buffer := bytes.NewBuffer(nil)
        boundary := "BamboerBoundary"
        Header := make(map[string]string)
        Header["From"] = message.from
        Header["To"] = strings.Join(message.to, ";")
        Header["Cc"] = strings.Join(message.cc, ";")
        Header["Bcc"] = strings.Join(message.bcc, ";")
        Header["Subject"] = message.subject
        Header["Content-Type"] = "multipart/related;boundary=" + boundary
        Header["Mime-Version"] = "1.0"
        Header["Date"] = time.Now().String().UTC()

        m.writeHeader(buffer, Header)
        body := "\r\n--" + boundary + "\r\n"
        body += "Content-Type: text/html; charset=UTF-8 \r\n"
        body += "\r\n" + message.body + "\r\n"
        buffer.WriteString(body)

        if message.attachment.WithFile {
                attachment := "\r\n--" + boundary + "\r\n"
                attachment += "Content-Transfer-Encoding:base64\r\n"
                attachment += "Content-Disposition:attachment\r\n"
                attachment += "Content-Type:" + message.attachment.ContentType + ";name=\"" + message.attachment.Name + "\"\r\n"
                attachment += "Content-ID: <" + message.attachment.Name + "> \r\n\r\n"
                imgsrc = "<p><img src=\"cid:" + message.attachment.Name + "> \r\n\t\t\t"
                buffer.WriteString(attachment)
                defer func() {
                        if err := recover(); err != nil {
                                fmt.Println("Error: ", err)
                        }
                }()
                m.writeFile(buffer, b)
        }

        var template = `
<html>
	<body>
		<p>text:%s</p><br>
		%s			
	</body>
</html>`
        var content = fmt.Sprintf(template, message.body, imgsrc)
        body := "\r\n--" + boundary + "\r\n"
        body += "Content-Type: text/html; charset=UTF-8 \r\n"
        body += content
        buffer.WriteString(body)

        buffer.WriteString("\r\n--" + boundary + "--")
        smtp.SendMail(m.host + m.port, m.auth, m.user, message.to, buffer.Bytes())
        return nil
}

func (m *mail) writeHeader(buffer *bytes.Buffer, Header map[string]string) string {
        header := ""
        for key, value := range Header {
                header += key + ":" + value + "\r\n"
        }
        header += "\r\n"
        buffer.WriteString(header)
        return header
}

func (m *mail) writeFile(buffer *bytes.Buffer, b []byte) {
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
