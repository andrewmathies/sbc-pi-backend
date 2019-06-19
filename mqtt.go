package bg

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var client MQTT.Client
var qos int

// MQTT functions

func handleMsg(beanID string, msg Msg) {
	log.Println("Handling MQTT Message\nBeanID: ", beanID, "\nMsg: ", msg)

	switch msg.Header {
	case "Hello":
		// create new unit and stuff it in the dict
		unit := Unit{Version: msg.Version, BeanID: beanID, Name: "", State: Idle}
		dict[ksuid.New().String()] = unit
	case "StartUpdate":
		// the backend published this, so do nothing
	case "Complete":
		// update status of unit and push that to frontend???
		id := msg.ID
		unit := dict[id]
		unit.State = Idle
		dict[id] = unit
	case "Fail":
		// update status of unit and push that to frontend???
		id := msg.ID
		unit := dict[id]
		unit.State = Failed
		dict[id] = unit
	default:
		log.Println("ERROR: unexpected MQTT message ", msg.Header)
	}
}

func publishMsg(beanID string, msg Msg) {
	log.Println("Publishing StartUpdate version on /unit/", beanID, "/")
	json, encodeErr := json.Marshal(msg)
	
	if encodeErr != nil {
		log.Println("ERROR: couldn't marshal msg for mqtt message ", encodeErr)
		return
	}

    token := client.Publish("/unit/" + beanID + "/", byte(qos), false, string(json))
    token.Wait()
}

// default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, mqttMsg MQTT.Message) {
	var msg Msg
	unpackErr := json.Unmarshal(mqttMsg.Payload(), &msg)
	if unpackErr != nil {
		log.Println("ERROR: couldn't unpack MQTT message ", unpackErr)
		return
	}

	topic := mqttMsg.Topic()
	topicParts := strings.Split(topic, "/")
	if len(topicParts) != 4 {
		log.Println("ERROR: badly formed MQTT topic ", topicParts)
		return
	}

	beanID := topicParts[2]
	handleMsg(beanID, msg)
}

func setupMQTT(tlsConfig *tls.Config) {
	qos = 1

	// opts contains broker address and other config info
	opts := MQTT.NewClientOptions().AddBroker("tls://saturten.com:8883")
  	opts.SetClientID("go-simple")
	opts.SetDefaultPublishHandler(f)
	opts.SetTLSConfig(tlsConfig)
	opts.SetUsername("andrew")
	opts.SetPassword("1plus2is3")
	
	// initiate connection with broker
	client = MQTT.NewClient(opts)
  	if token := client.Connect(); token.Wait() && token.Error() != nil {
    	panic(token.Error())
	}
	log.Println("Connected to MQTT broker")
	
	// subscribe to wildcard topic
	if token := client.Subscribe("/unit/+/", byte(qos), nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	log.Println("Subscribed to /unit/+/")
}