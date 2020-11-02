package mail

import (
	"encoding/json"
	"fmt"
	"github.com/polundrra/PriceTracker/internal/tracker/workers"
	"net/smtp"
)

func SendEmail(body []byte) error {
	info := workers.MessageInfo{}

	if err := json.Unmarshal(body, &info); err != nil {
		return err
	}

	from := "from@gmail.com"
	password := "<Email Password>"

	to := info.Emails

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte(fmt.Sprintf("Hey, check out this ad: %s. Now the price is %d", info.Ad, info.NewPrice))

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost + ":" + smtpPort, auth, from, to, message)
	if err != nil {
		return err
	}

	fmt.Println("Emails sent successfully")
	return nil
}
