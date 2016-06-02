package flanders

import (
	"fmt"
	"time"

	"lab.getweave.com/weave/flanders/db"
	"lab.getweave.com/weave/flanders/hep"
	"lab.getweave.com/weave/flanders/log"
)

func processPacket(packet []byte) error {

	defer func() {
		if r := recover(); r != nil {
			log.Err("Sip parser panicked, but I recovered...")
			log.Err(string(packet))
			log.Err(fmt.Sprintf("%[ 02]X", packet))
		}
	}()

	hepMsg, err := hep.NewHepMsg(packet)
	if err != nil {
		return fmt.Errorf("Unable to parse packet: %s", err)
	}

	if hepMsg.SipMsg == nil || hepMsg.SipMsg.StartLine == nil {
		return nil
	}

	switch hepMsg.SipMsg.StartLine.Method {
	case "OPTIONS":
		return nil
	case "SUBSCRIBE":
		return nil
	case "NOTIFY":
		return nil
	case "REGISTER":
		return nil
	}

	if hepMsg.SipMsg.Cseq != nil {
		switch hepMsg.SipMsg.Cseq.Method {
		case "OPTIONS":
			return nil
		case "SUBSCRIBE":
			return nil
		case "NOTIFY":
			return nil
		case "REGISTER":
			return nil
		}
	}

	var datetime time.Time

	//log.Debug(string(packet))
	if hepMsg.Timestamp != 0 {
		datetime = time.Unix(int64(hepMsg.Timestamp), int64(hepMsg.TimestampMicro)*1000)
	} else {
		datetime = time.Now()
	}

	dbObject := db.NewDbObject()

	dbObject.Datetime = datetime
	dbObject.MicroSeconds = datetime.Nanosecond() / 1000
	dbObject.Method = hepMsg.SipMsg.StartLine.Method + hepMsg.SipMsg.StartLine.Resp
	dbObject.ReplyReason = hepMsg.SipMsg.StartLine.RespText
	dbObject.SourceIp = hepMsg.Ip4SourceAddress
	dbObject.SourcePort = hepMsg.SourcePort
	dbObject.DestinationIp = hepMsg.Ip4DestinationAddress
	dbObject.DestinationPort = hepMsg.DestinationPort
	dbObject.CallId = hepMsg.SipMsg.CallId
	dbObject.FromUser = hepMsg.SipMsg.From.URI.User
	dbObject.FromDomain = hepMsg.SipMsg.From.URI.Host
	dbObject.FromTag = hepMsg.SipMsg.From.Tag
	dbObject.ToUser = hepMsg.SipMsg.To.URI.User
	dbObject.ToDomain = hepMsg.SipMsg.To.URI.Host
	dbObject.ToTag = hepMsg.SipMsg.To.Tag
	dbObject.UserAgent = hepMsg.SipMsg.UserAgent
	dbObject.Cseq = hepMsg.SipMsg.Cseq.Val
	for _, header := range hepMsg.SipMsg.Headers {
		if header.Header == "x-cid" {
			dbObject.CallIdAleg = header.Val
		}
	}
	dbObject.Msg = hepMsg.SipMsg.Msg

	h.broadcast <- *dbObject

	err = dbObject.Save()
	if err != nil {
		return fmt.Errorf("unable to save sip message: %s", err)
	}

	return nil
}