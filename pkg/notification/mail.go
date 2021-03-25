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
        name        string
        contentType string
        withFile    bool
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

func (m *mail) Send(message Message) error {
        buffer := bytes.NewBuffer(nil)
        boundary := "GoBoundary"
        Header := make(map[string]string)
        Header["From"] = message.from
        Header["To"] = strings.Join(message.to, ";")
        Header["Cc"] = strings.Join(message.cc, ";")
        Header["Bcc"] = strings.Join(message.bcc, ";")
        Header["Subject"] = message.subject
        Header["Content-Type"] = "multipart/mixed;boundary=" + boundary
        Header["Mime-Version"] = "1.0"
        Header["Date"] = time.Now().String()

        m.writeHeader(buffer, Header)
        body := "\r\n--" + boundary + "\r\n"
        body += "Content-Type:" + message.contentType + "\r\n"
        body += "\r\n" + message.body + "\r\n"
        buffer.WriteString(body)

        if message.attachment.withFile {
                attachment := "\r\n--" + boundary + "\r\n"
                attachment += "Content-Transfer-Encoding:base64\r\n"
                attachment += "Content-Disposition:attachment\r\n"
                attachment += "Content-Type:" + message.attachment.contentType + ";name=\"" + message.attachment.name + "\"\r\n"
                buffer.WriteString(attachment)
                defer func() {
                        if err := recover(); err != nil {
                                fmt.Println("Error: ", err)
                        }
                }()
                m.writeFile(buffer, message.attachment.name)
        }
        buffer.WriteString("\r\n--" + boundary + "--")
        smtp.SendMail(m.host+m.port, m.auth, m.user, m.to, buffer.Bytes())
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

func (m *mail) writeFile(buffer *bytes.Buffer, fileName string) {
        file, err := ioutil.ReadFile(fileName)
        if err != nil {
                fmt.Println(err)
        }
        payload := make([]byte, base64.StdEncoding.EncodeLen(len(file)))
        base64.StdEncoding.Encode(payload, file)
        buffer.WriteString("\r\n")
        for index, line := 0, len(payload); index < line; index++ {
                buffer.WriteByte(payload[index])
                if (index+1)%76 == 0 {
                        buffer.WriteString("\r\n")
                }
        }
}
