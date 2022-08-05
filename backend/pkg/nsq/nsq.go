package nsq

import (
	"encoding/json"
	"fmt"
	"qiu/backend/pkg/e"
	log "qiu/backend/pkg/logging"
	"qiu/backend/pkg/setting"

	"github.com/nsqio/go-nsq"
)

var Producer *nsq.Producer
var ConsumerUploadImage *nsq.Consumer

func Setup() {
	config := nsq.NewConfig()
	var err error
	Producer, err = nsq.NewProducer(setting.NsqSetting.Nsqd, config)

	if err != nil {
		panic(err)
	}
	go runConsumerUploadImage()
}

func runConsumerUploadImage() {
	config := nsq.NewConfig()
	ConsumerUploadImage, _ := nsq.NewConsumer(e.TOPIC_UPLOAD_IMAGE, "ch", config)
	// ConsumerUploadImage.SetLogger(nsq.Lo, nsq.LogLevelInfo)
	defer ConsumerUploadImage.Stop()
	ConsumerUploadImage.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		var msg UploadImageMessage
		err := json.Unmarshal(message.Body, &msg)
		if err != nil {
			log.Logger.Panic(err)
		}
		log.Logger.Info("[Nsq] Got a message: ", msg)
		return nil
	}))
	err := ConsumerUploadImage.ConnectToNSQLookupd(setting.NsqSetting.NsqLookup)
	if err != nil {
		fmt.Println("setting.NsqSetting.NsqLookup", setting.NsqSetting.NsqLookup)
		log.Logger.Panic(err, setting.NsqSetting.NsqLookup)
		return
	}
	<-ConsumerUploadImage.StopChan
	// sigChan := make(chan os.Signal, 1)
	// signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	// <-sigChan
}
