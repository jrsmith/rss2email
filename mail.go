package main

import (
    "fmt"
    "time"
    "net/smtp"
    "bytes"
)

type Email struct {
    Recipients []string
    Sender string
    Subject string
    MimeType string
    Body string
}

func sendItem(subject string, content string) {

    fmt.Println("Sending mail")

    var buffer bytes.Buffer

    auth := smtp.PlainAuth(
        "",
        config.SMTP.Username,
        config.SMTP.Password,
        config.SMTP.Host,
    )

    buffer.WriteString("To: "+config.ToEmail[0]+"\n")
    buffer.WriteString("Subject: ")
    buffer.WriteString((*subject))
    buffer.WriteString("\nMIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n")
    buffer.WriteString((*content))

    time.Sleep(10 * 1e9)

    err := smtp.SendMail(
        config.SMTP.OutgoingServer,
        auth,
        config.SMTP.From,
        config.ToEmail,
        buffer.Bytes(),
    )

    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println("Mail sent")
    }

}