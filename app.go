package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const FacebookMsgAPIToken = "token"

const BotID = "botid"

func webhook(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	switch req.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(req.URL.Query().Get("hub.challenge")))

	case http.MethodPost:

		var receivedMessage ReceivedMessage

		err := JSONUnmarshaler(req.Body, &receivedMessage)

		if err != nil {
			log.Println("Failed to Parse Webhook")
			return
		}

		for _, entry := range receivedMessage.Entry {
			for _, messages := range entry.Messaging {
				if messages.Sender.ID.String() != BotID {
					senderID, _ := messages.Sender.ID.Int64()
					if messages.Message != nil {
						switch {
						case messages.Message.Text != nil:
							sendFacebookMessage(int(senderID), *messages.Message.Text)
						case messages.Message.Attachments != nil && len(messages.Message.Attachments) != 0:
							for _, attachment := range messages.Message.Attachments {
								sendFacebookImage(int(senderID), attachment.Payload.URL)
							}
						default:
						}
					} else if messages.PostBack != nil {
						sendFacebookMessage(int(senderID), messages.PostBack.Payload)
					}
				}
			}
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
	}
}

func JSONUnmarshaler(body io.Reader, i interface{}) error {

	bufbody := new(bytes.Buffer)
	length, err := bufbody.ReadFrom(body)
	if err != nil && err != io.EOF {
		log.Println("EOF Error")
		return err
	}

	log.Println(bufbody.String())
	if err = json.Unmarshal(bufbody.Bytes()[:length], &i); err != nil {
		log.Println("JSONUnmarshaler Failed : " + err.Error())
		log.Println("Body : ")
		log.Println(bufbody.String())
		return err
	}

	return nil
}

func sendFacebookMessage(recipient int, text string) error {

	reqBody, err := json.Marshal(&SendMessage{
		Recipient: &ID{
			ID: recipient,
		},
		Message: &Text{Text: text},
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "https://graph.facebook.com/v2.8/me/messages?access_token="+FacebookMsgAPIToken, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json; charset=UTF-8")

	log.Println(req.URL.String())

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil

}

func sendFacebookImage(recipient int, url string) error {

	reqBody, err := json.Marshal(&SendImage{
		Recipient: &ID{
			ID: recipient,
		},
		Message: &AttachmentMessage{
			Attachment: &Attachment{
				Payload: &Payload{
					URL: url,
				},
				Type: "image",
			},
		},
	})
	if err != nil {
		return err
	}

	fmt.Println("SendImage JSON : \r\n" + string(reqBody))

	req, err := http.NewRequest(http.MethodPost, "https://graph.facebook.com/v2.8/me/messages?access_token="+FacebookMsgAPIToken, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json; charset=UTF-8")

	log.Println(req.URL.String())

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil

}

func sendFacebookButton(recipient string) error {

	reqBody := []byte(strings.Replace(`{  
    "recipient":{  
        "id":"USER_ID"
    },
    "message":{  
        "attachment":{  
            "type":"template",
            "payload":{  
                "template_type":"button",
                "text":"ボタン選択",
                "buttons":[  
                    {  
                        "type":"postback",
                        "title":"ボタン1",
                        "payload":"button1"
                    },
                    {  
                        "type":"postback",
						"title":"ボタン2",
                        "payload":"button2"
                    }
                ]
            }
        }
    }
}`, "USER_ID", recipient, 1))

	fmt.Println("SendButton JSON : \r\n" + string(reqBody))

	req, err := http.NewRequest(http.MethodPost, "https://graph.facebook.com/v2.8/me/messages?access_token="+FacebookMsgAPIToken, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json; charset=UTF-8")

	log.Println(req.URL.String())

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil

}

func main() {

	var port string

	if port = os.Getenv("PORT"); len(port) == 0 {
		port = "8080"
	}

	http.HandleFunc("/webhook", webhook)

	fmt.Println("Starting app on port " + port)
	http.ListenAndServe(":"+port, nil)
}
