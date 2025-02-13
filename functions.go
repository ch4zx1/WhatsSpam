package main

import (
	"context"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"os"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func sendSpoofedReplyThis(chatID types.JID, spoofedID types.JID, msgID string, text string, msg *waProto.Message) (*waProto.Message, *whatsmeow.SendResponse, error) {
	newmsg := &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text:        proto.String(text),
			PreviewType: waProto.ExtendedTextMessage_IMAGE.Enum(),
			ContextInfo: &waProto.ContextInfo{
				Participant:   proto.String(spoofedID.String()),
				QuotedMessage: msg.GetExtendedTextMessage().GetContextInfo().GetQuotedMessage(),
			},
		},
	}
	resp, err := cli.SendMessage(context.Background(), chatID, newmsg)
	if err != nil {
		log.Errorf("Error sending reply message: %v", err)
		return msg, &resp, err
	} else {
		log.Infof("Message sent (server timestamp: %s)", resp.Timestamp)
		return msg, &resp, nil
	}
}

func sendSpoofedReplyMessage(chatID types.JID, fromID types.JID, msgID string, replyText string, myTtext string) (*waProto.Message, *whatsmeow.SendResponse, error) {
	msg := &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text: proto.String(myTtext),
			ContextInfo: &waProto.ContextInfo{
				Participant: proto.String(fromID.String()),
				QuotedMessage: &waProto.Message{
					Conversation: proto.String(replyText),
				},
			},
		},
	}
	resp, err := cli.SendMessage(context.Background(), chatID, msg)
	if err != nil {
		log.Errorf("Error sending reply message: %v", err)
		return msg, &resp, err
	} else {
		log.Infof("Message sent (server timestamp: %s)", resp.Timestamp)
		return msg, &resp, nil
	}
}

func sendSpoofedReplyImg(chatID types.JID, fromID types.JID, msgID string, file string, replyText string, myTtext string) (*waProto.Message, *whatsmeow.SendResponse, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		log.Errorf("Failed to read %s: %v", file, err)
		return &waProto.Message{}, &whatsmeow.SendResponse{}, err
	}
	uploaded, err := cli.Upload(context.Background(), data, whatsmeow.MediaImage)
	if err != nil {
		log.Errorf("Failed to upload file: %v", err)
		return &waProto.Message{}, &whatsmeow.SendResponse{}, err
	}

	msg := &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text:        proto.String(myTtext),
			PreviewType: waProto.ExtendedTextMessage_IMAGE.Enum(),
			ContextInfo: &waProto.ContextInfo{
				Participant: proto.String(fromID.String()),
				QuotedMessage: &waProto.Message{
					ImageMessage: &waProto.ImageMessage{
						Caption:           proto.String(replyText),
						DirectPath:        proto.String(uploaded.DirectPath),
						MediaKey:          uploaded.MediaKey,
						Mimetype:          proto.String(http.DetectContentType(data)),
						FileEncSHA256:     uploaded.FileEncSHA256,
						FileSHA256:        uploaded.FileSHA256,
						FileLength:        proto.Uint64(uint64(len(data))),
						Height:            proto.Uint32(100),
						Width:             proto.Uint32(100),
						MediaKeyTimestamp: proto.Int64(time.Now().Unix()),
					},
				},
			},
		},
	}
	resp, err := cli.SendMessage(context.Background(), chatID, msg)
	if err != nil {
		log.Errorf("Error sending reply message: %v", err)
		return msg, &resp, err
	} else {
		log.Infof("Message sent (server timestamp: %s)", resp.Timestamp)
		return msg, &resp, nil
	}
}

func sendSpoofedTalkDemo(chatJID types.JID, spoofedJID types.JID, toGender string, language string, spoofedFile string) {
	msgmap := make(map[string]map[string]map[int]string)
	msgmap["msg1"] = make(map[string]map[int]string)
	msgmap["msg1"]["generic"] = make(map[int]string)
	msgmap["msg1"]["generic"][0] = ("Primeira")
	msgmap["msg2"] = make(map[string]map[int]string)
	msgmap["msg2"]["generic"] = make(map[int]string)
	msgmap["msg2"]["generic"][0] = ("Segunda")
	msgmap["msg3"] = make(map[string]map[int]string)
	msgmap["msg3"]["generic"] = make(map[int]string)
	msgmap["msg3"]["generic"][0] = "Terceira"

	_, err := cli.SendMessage(context.Background(), chatJID, &waProto.Message{Conversation: proto.String(msgmap[language]["generic"][0])})

	if err != nil {
		log.Errorf("Error on sending spoofed msg: %v", err)
	} else {
		log.Infof("mensagem disparada para %s, via (%s - num conectado) ", chatJID.String(), spoofedJID.String())
	}
}

func sendConversationMessage(recipient_jid types.JID, text string) (*waProto.Message, *whatsmeow.SendResponse, error) {
	msg := &waProto.Message{Conversation: proto.String(text)}
	resp, err := cli.SendMessage(context.Background(), recipient_jid, msg)
	if err != nil {
		log.Errorf("Error sending message: %v", err)
		return msg, &resp, err
	} else {
		log.Infof("Message sent (server timestamp: %s)", resp.Timestamp)
		return msg, &resp, nil
	}
}

func sendMessage(recipient_jid types.JID, msg *waProto.Message) (*waProto.Message, *whatsmeow.SendResponse, error) {
	resp, err := cli.SendMessage(context.Background(), recipient_jid, msg)
	if err != nil {
		log.Errorf("Error sending message: %v", err)
		return msg, &resp, err
	} else {
		log.Infof("Message sent (server timestamp: %s)", resp.Timestamp)
		return msg, &resp, nil
	}
}

func getMsg(evt *events.Message) string {
	msg := ""

	if evt.Message.Conversation != nil {
		msg = *evt.Message.Conversation
	}

	if evt.Message.ExtendedTextMessage != nil {
		msg = *evt.Message.ExtendedTextMessage.Text
	}

	return msg
}

func download(evt_type string, file interface{}, mimetype string, evt *events.Message, rawEvt interface{}) (err error) {
	if file != nil {
		exts, _ := mime.ExtensionsByType(mimetype)
		file_name := fmt.Sprintf("%s%s", evt.Info.ID, exts[0])
		if mimetype == "text/vcard" {
			data := file.(*waProto.ContactMessage)
			err = postEventFile(evt_type, rawEvt, nil, file_name, []byte(*data.Vcard))
		} else {
			data, err := cli.Download(file.(whatsmeow.DownloadableMessage))
			if err != nil {
				postError(evt_type, fmt.Sprintf("%s Failed to download", evt_type), rawEvt)
				return err
			}
			err = postEventFile(evt_type, rawEvt, nil, file_name, data)
		}
		if err != nil {
			postError(evt_type, fmt.Sprintf("%s Failed to save event", evt_type), rawEvt)
			return err
		}
		return nil
	}
	return errors.New("File is nil")
}

func parseJID(arg string) (types.JID, bool) {
	if arg[0] == '+' {
		arg = arg[1:]
	}
	if !strings.ContainsRune(arg, '@') {
		return types.NewJID(arg, types.DefaultUserServer), true
	} else {
		recipient, err := types.ParseJID(arg)
		if err != nil {
			log.Errorf("Invalid JID %s: %v", arg, err)
			return recipient, false
		} else if recipient.User == "" {
			log.Errorf("Invalid JID %s: no server specified", arg)
			return recipient, false
		}
		return recipient, true
	}
}
