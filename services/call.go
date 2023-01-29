package services

import (
	"github.com/abr-ooo/hodhod/configs"
	"github.com/kavenegar/kavenegar-go"
	"github.com/abr-ooo/hodhod/log"
)

func Call(messages string, number string) error {
	c := configs.Cfg
	api := kavenegar.New(c.Kaveh.Token)

	_, err := api.Call.MakeTTS(number, messages, &kavenegar.CallParam{})
	if err != nil {
		log.Log.Errorf("[Call] [To: %s] have error : %s", number, err)
		return err
	}
	log.Log.Infof("[Call] [To: %s] succesful", number)

	return nil
}
