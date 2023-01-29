package services

import (
	"fmt"

	"github.com/imroc/req"
	"github.com/kavenegar/kavenegar-go"
	"github.com/mhd7966/hodhod/configs"
	"github.com/mhd7966/hodhod/log"
)

func KavenegarSendSMS(message string, numbers []string) error {
	c := configs.Cfg.Kaveh
	api := kavenegar.New(c.Token)

	// We Could Use Response But Don't Need it for Now ;)
	_, err := api.Message.Send(c.Number, numbers, message, nil)
	if err != nil {
		log.Log.Debug("SMS. Kaveh SMS have an error", err)
		return err
	}

	return nil
}

func SignalSendSMS(message string, numbers []string) error {
	c := configs.Cfg.Signal
	header := req.Header{
		"Authorization": "Bearer " + c.Token,
	}
	data := req.Param{
		"from":    c.Number,
		"message": message,
		"numbers": numbers,
	}

	r, err := req.Post("https://panel.signalads.com/rest/api/v1/message/send.json", header, req.BodyJSON(&data))
	if err != nil {
		log.Log.Debug("SMS. Signal SMS have an error", err)
		return err
	}

	var response map[string]interface{}
	r.ToJSON(&response)

	if valid, ok := response["success"].(bool); !ok && !valid {
		log.Log.Debug("SMS. error on Signal Response")
		return fmt.Errorf("error on Signal Response")
	}

	return nil
}

// package drivers

// import (
// 	"fmt"
// 	log "github.com/sirupsen/logrus"
// 	"os"

// 	"github.com/imroc/req"
// 	"github.com/kavenegar/kavenegar-go"
// )

// func SendMessage(messages string, numbers []string) {
// 	err := KavenegarSendSMS(messages, numbers)
// 	if err == nil {
// 		return
// 	}
// 	log.Error(err.Error())

// 	err = SignalSendSMS(messages, numbers)
// 	if err != nil {
// 		log.Error(err.Error())
// 		log.Error("[SMS] [All drivers are not available]")
// 	}

// 	log.Infof("[SMS] [To: %v]", numbers)
// }

// func KavenegarSendSMS(message string, numbers []string) error {
// 	api := kavenegar.New(os.Getenv("KAVENEGAR_KEY"))

// 	// We Could Use Response But Don't Need it for Now ;)
// 	_, err := api.Message.Send(os.Getenv("KAVENEGAR_NUMBER"), numbers, message, nil)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func SignalSendSMS(message string, numbers []string) error {
// 	header := req.Header{
// 		"Authorization": "Bearer " + os.Getenv("SIGNAL_KEY"),
// 	}
// 	data := req.Param{
// 		"from":    os.Getenv("SIGNAL_NUMBER"),
// 		"message": message,
// 		"numbers": numbers,
// 	}

// 	r, err := req.Post("https://panel.signalads.com/rest/api/v1/message/send.json", header, req.BodyJSON(&data))
// 	if err != nil {
// 		return err
// 	}

// 	var response map[string]interface{}
// 	r.ToJSON(&response)

// 	if valid, ok := response["success"].(bool); !ok && !valid {
// 		return fmt.Errorf("Error on Signal Response")
// 	}

// 	return nil
// }
