package notice

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

var EmailSrv *EmailService

func SendAlarm(processName, command string) {
	if EmailSrv == nil {
		return
	}

	EmailSrv.SendWithHtmlTpl(processName, command)
}

// SmtpConfig smtp config
type SmtpConfig struct {
	Host      string             // 邮件服务器地址
	Port      int                // 端口
	User      string             // 发送邮件用户账号
	Pwd       string             // 授权密码
	Tpl       *template.Template // 发送消息模板
	ToUser    string             // 邮件接收方
	ServerURL string             // 外部访问的URL
}

type EmailService struct {
	Cfg    SmtpConfig
	Client *mail.SMTPClient
}

func NewEmailService(cfg SmtpConfig) *EmailService {
	mailSrv := &EmailService{}
	mailSrv.Cfg = cfg

	client := mail.NewSMTPClient()

	// SMTP Client
	client.Host = cfg.Host
	client.Port = cfg.Port
	client.Username = cfg.User
	client.Password = cfg.Pwd
	client.Encryption = mail.EncryptionSTARTTLS
	client.ConnectTimeout = 10 * time.Second
	client.SendTimeout = 10 * time.Second

	// Connect to client
	smtpClient, err := client.Connect()
	if err != nil {
		fmt.Printf("client.Connect error %s\n", err.Error())
		return nil
	}

	mailSrv.Client = smtpClient
	return mailSrv
}

func (s *EmailService) SendEmail(htmlBody string) error {
	// Create the email message
	email := mail.NewMSG()

	email.SetFrom(fmt.Sprintf("From Supervisord <%s>", s.Cfg.User)).
		AddTo(s.Cfg.ToUser).
		SetSubject("Supervisord Exception Notice")

	// Get from each mail
	email.GetFrom()
	email.SetBody(mail.TextHTML, htmlBody)

	// Send with high priority
	email.SetPriority(mail.PriorityHigh)

	// always check error after send
	if email.Error != nil {
		return email.Error
	}

	// Pass the client to the email message to send it
	return email.Send(s.Client)
}

func (s *EmailService) SendWithHtml(tplPath string) error {
	emailHtml, err := os.ReadFile(tplPath)
	if err != nil {
		fmt.Println("os.ReadFile error :" + err.Error())
		return err
	}

	err = s.SendEmail(string(emailHtml))
	if err != nil {
		fmt.Printf("SendEmail Error : %s\n", err.Error())
		return err
	}
	return nil
}

type ProcessError struct {
	ProgramName string
	Command     string
	ExceptTime  string
	ServerAddr  string
}

func (s *EmailService) SendWithHtmlTpl(processName, command string) error {
	// 利用给定数据渲染模板
	proErr := ProcessError{
		ProgramName: processName,
		Command:     command,
		ExceptTime:  time.Now().Format("2006-01-02 15:04:05"),
		ServerAddr:  s.Cfg.ServerURL,
	}
	var body bytes.Buffer
	err := s.Cfg.Tpl.Execute(&body, proErr)
	if err != nil {
		fmt.Println("Tpl.Execute,err", err)
		return err
	}

	err = s.SendEmail(body.String())
	if err != nil {
		fmt.Printf("SendEmail Error: %s\n", err.Error())
		return err
	}

	return nil
}
