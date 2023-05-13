package main

import (
	//	"../common/const/device"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/sibivishnu/Weather/common/const/device"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"path"
	"strings"
)

type MailDoc struct {
	Params   interface{}
	Subject  string
	Template string
}

type DailyMailDoc struct {
	TotalCount int
	Cat1Count  int
	Cat2Count  int
	Cat3Count  int
}

func sendEmail(mailDoc MailDoc) error {

	// load the actual template text
	body, err := ioutil.ReadFile(path.Join("/templates", mailDoc.Template))
	if err != nil {
		log.Printf(err.Error())
		return err
	}

	from := mail.Address{smtpFromName, smtpFromAddress}
	to := mail.Address{"", smtpRecipients}

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = mailDoc.Subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	tmpl, err := template.New("body").Parse(string(body))
	var doc bytes.Buffer
	err = tmpl.Execute(&doc, mailDoc.Params)
	if err != nil {
		log.Printf("1")
		log.Printf(err.Error())
		return err
	}

	message += "\r\n" + doc.String()

	// Connect to the SMTP Server
	servername := smtpHost + ":" + smtpPort

	log.Printf(servername)

	host, _, _ := net.SplitHostPort(servername)
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		log.Printf(err.Error())
		return err
	}

	cl, err := smtp.NewClient(conn, host)
	defer cl.Quit()
	if err != nil {
		log.Printf(err.Error())
		return err
	}

	// Auth
	if err = cl.Auth(auth); err != nil {
		log.Printf(err.Error())
		return err
	}

	// To && From
	if err = cl.Mail(from.Address); err != nil {
		log.Printf(err.Error())
		return err
	}

	if err = cl.Rcpt(to.Address); err != nil {
		log.Printf(err.Error())
		return err
	}

	// Data
	w, err := cl.Data()
	if err != nil {
		log.Printf(err.Error())
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Printf(err.Error())
		return err
	}

	err = w.Close()
	if err != nil {
		log.Printf(err.Error())
		return err
	}

	return nil
}

func sendDailyMail() {
	dmd := getMailCountInfo()

	mailSendDoc := MailDoc{}
	mailSendDoc.Subject = "Daily device count info"
	mailSendDoc.Params = dmd
	mailSendDoc.Template = "emailDeviceCount"

	err := sendEmail(mailSendDoc)
	if err != nil {
		log.Printf("Error sending the mail : %+v\n", err)
	}
}

func getMailCountInfo() DailyMailDoc {

	deviceListMap, _ := redisInstance.QueryCache("devicerequested:*")
	dmd := DailyMailDoc{}

	for dKey, _ := range deviceListMap {
		keyArr := strings.Split(dKey, ":")
		if len(keyArr) < 2 {
			log.Printf("Wrong device requested key %s", dKey)
			continue
		}

		deviceID := keyArr[1]
		key := "device:" + deviceID

		data, err := redisInstance.GetCachedData(key)

		// No data found from the cache
		if err != nil {
			log.Printf("Device %s not found in cache", deviceID)
			continue
		}

		var d device.Device
		json.Unmarshal(data, &d)

		if d.Category == "1" {
			dmd.Cat1Count = dmd.Cat1Count + 1
		} else if d.Category == "2" {
			dmd.Cat2Count = dmd.Cat2Count + 1
		} else if d.Category == "3" {
			dmd.Cat3Count = dmd.Cat3Count + 1
		}

	}

	dmd.TotalCount = dmd.Cat1Count + dmd.Cat2Count + dmd.Cat3Count
	log.Printf("Total : %d , Category 1 : %d , Category 2 : %d , Category 3 : %d", dmd.TotalCount, dmd.Cat1Count, dmd.Cat2Count, dmd.Cat3Count)

	return dmd
}
