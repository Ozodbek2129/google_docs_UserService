package email

import (
        "bytes"
        "fmt"
        "html/template"
        "log"
        "math/rand"
        "net/smtp"
        "strconv"
        "time"
)

func Email(email string) (string, error) {

        // Seed the random number generator with a cryptographically secure value
        source := rand.NewSource(time.Now().UnixNano())
        myRand := rand.New(source)

        // Generate a random 6-digit number (100000 to 999999)
        randomNumber := myRand.Intn(900000) + 100000
        code := strconv.Itoa(randomNumber)

        err := SendCode(email, code)

        if err != nil {
                return "", err
        }

        return "Sizning emailingizga xabar yuborildi", nil
}

func SendCode(email string, code string) error {
        // sender data
        from := "articanconnection@gmail.com"
        password := "colo twdh fabv kcvj"

        // Receiver email address
        to := []string{
                email,
        }

        // smtp server configuration.
        smtpHost := "smtp.gmail.com"
        smtpPort := "587"

        // Authentication.
        auth := smtp.PlainAuth("", from, password, smtpHost)

        t, err := template.ParseFiles("api/email/template.html")
        if err != nil {
                log.Fatalf("Error parsing template: %v", err)
                return err
        }

        var body bytes.Buffer

        mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
        body.Write([]byte(fmt.Sprintf("Subject: Your verification code \n%s\n\n", mimeHeaders)))

        err = t.Execute(&body, struct {
                Passwd string
        }{
                Passwd: code,
        })
        if err != nil {
                log.Fatalf("Error executing template: %v", err)
                return err
        }

        // Sending email.
        err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
        if err != nil {
                log.Fatalf("Error sending email: %v", err)
                return err
        }
        return nil
}