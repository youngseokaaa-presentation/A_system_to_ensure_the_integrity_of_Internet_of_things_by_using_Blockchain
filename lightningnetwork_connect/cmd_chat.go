package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/binary"
	
	"log"
	"strings"
	"time"
	"os/exec"

	"github.com/lightningnetwork/lnd/keychain"
	"github.com/lightningnetwork/lnd/lntypes"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/routing/route"
	"google.golang.org/grpc"

	"github.com/jroimartin/gocui"
	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
	"github.com/lightningnetwork/lnd/lnrpc/signrpc"
	"github.com/urfave/cli"
)

var chatCommand = cli.Command{
	Name:      "chat",
	Category:  "Chat",
	ArgsUsage: "recipient_pubkey",
	Usage:     "Use lnd as a p2p messenger application.",
	Action:    actionDecorator(chat),
	Flags: []cli.Flag{
		cli.Uint64Flag{
			Name:  "amt_msat",
			Usage: "payment amount per chat message",
			Value: 1000,
		},
	},
}

var byteOrder = binary.BigEndian

const (
	tlvMsgRecord    = 34349334
	tlvSigRecord    = 34349337
	tlvSenderRecord = 34349339
	tlvTimeRecord   = 34349343

	// TODO: Reference lnd master constant when available.
	tlvKeySendRecord = 5482373484
)

type messageState uint8

const (							//채팅 보낼때 성공하면 v뜨고 실패하면 x뜨게 하는 변수
	statePending messageState = iota

	stateDelivered

	stateFailed
)

type chatLine struct {					//내가 보내는 채팅에 들어가는 내용들 채팅, 보내는 사람, 받는 사람, 수수료, 시간 등등 채팅 gui에 들어가는 값들
	text      string
	sender    route.Vertex
	recipient *route.Vertex
	state     messageState
	fee       uint64
	timestamp time.Time
}

func(this *chatLine) Format() string {
		return fmt.Sprintf("%s",this.text)
		}

var (
	msgLines       []chatLine
	destination    *route.Vertex
	runningBalance map[route.Vertex]int64 = make(map[route.Vertex]int64)

	keyToAlias = make(map[route.Vertex]string)
	aliasToKey = make(map[string]route.Vertex)

	self route.Vertex
)

func initAliasMaps(conn *grpc.ClientConn) error {
	client := lnrpc.NewLightningClient(conn)

	graph, err := client.DescribeGraph(
		context.Background(),
		&lnrpc.ChannelGraphRequest{},
	)
	if err != nil {
		return err
	}

	aliasCount := make(map[string]int)
	for _, node := range graph.Nodes {
		alias := node.Alias
		aliasCount[alias]++
	}

	for _, node := range graph.Nodes {
		alias := node.Alias

		key, err := route.NewVertexFromStr(node.PubKey)
		if err != nil {
			return err
		}

		if aliasCount[alias] > 1 {
			alias += "-" + node.PubKey[:6]
		}

		aliasToKey[alias] = key
		aliasToKey[key.String()] = key

		keyToAlias[key] = alias
	}

	info, err := client.GetInfo(context.Background(), &lnrpc.GetInfoRequest{})
	if err != nil {
		return err
	}

	self, err = route.NewVertexFromStr(info.IdentityPubkey)
	if err != nil {
		return err
	}

	return nil
}

func setDest(destStr string) {				//커맨드창에서 send 입력할때 주소값 받아서 여기로 가져옴
	dest, err := route.NewVertexFromStr(destStr)
	if err == nil {
		destination = &dest
	}

	if dest, ok := aliasToKey[destStr]; ok {
		destination = &dest
	}
}

func chat(ctx *cli.Context) error {
	chatMsgAmt := int64(ctx.Uint64("amt_msat"))

	conn := getClientConn(ctx, false)
	defer conn.Close()

	err := initAliasMaps(conn)
	if err != nil {
		return err
	}

	if ctx.NArg() != 0 {
		destStr := ctx.Args().First()
		setDest(destStr)
	}

	mainRpc := lnrpc.NewLightningClient(conn)
	client := routerrpc.NewRouterClient(conn)
	signClient := signrpc.NewSignerClient(conn)

	req := &lnrpc.InvoiceSubscription{}
	rpcCtx := context.Background()
	stream, err := mainRpc.SubscribeInvoices(rpcCtx, req)
	if err != nil {
		return err
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	
	addMsg_py :=func(line chatLine) int {
		msgLines = append(msgLines,line)
		return len(msgLines)-1
		}
	sendreply_func :=func() error{				//받은 명령이 temp일 경우 파이썬을 실행해 온도를 측정하고 상대방에게 전송한다.
		var g *gocui.Gui
		var test string
		cmd := exec.Command("python3","AdafruitDHT.py","11","4")
                out, err :=cmd.CombinedOutput()
                if err != nil{
                fmt.Println(fmt.Sprint(err)+string(out))
                }
                test = string(out)	
		newMsg := test
	
		d := *destination
		msgIdx := addMsg_py(chatLine{
			sender:    self,
			text:      newMsg,
			recipient: &d,
		})	

		payAmt := runningBalance[*destination]
		if payAmt < chatMsgAmt {
			payAmt = chatMsgAmt
		}
		if payAmt > 10*chatMsgAmt {
			payAmt = 10 * chatMsgAmt
		}

		var preimage lntypes.Preimage
		if _, err := rand.Read(preimage[:]); err != nil {
			return err
		}
		hash := preimage.Hash()

		// Message sending time stamp
		timestamp := time.Now().UnixNano()
		var timeBuffer [8]byte
		byteOrder.PutUint64(timeBuffer[:], uint64(timestamp))

		// Sign all data.
		signData, err := getSignData(
			self, *destination, timeBuffer[:], []byte(newMsg),
		)
		if err != nil {
			return err
		}

		signResp, err := signClient.SignMessage(context.Background(), &signrpc.SignMessageReq{
			Msg: signData,
			KeyLoc: &signrpc.KeyLocator{
				KeyFamily: int32(keychain.KeyFamilyNodeKey),
				KeyIndex:  0,
			},
		})
		if err != nil {
			return err
		}
		signature := signResp.Signature

		customRecords := map[uint64][]byte{
			tlvMsgRecord:     []byte(newMsg),
			tlvSenderRecord:  self[:],
			tlvTimeRecord:    timeBuffer[:],
			tlvSigRecord:     signature,
			tlvKeySendRecord: preimage[:],
		}

		req := routerrpc.SendPaymentRequest{
			PaymentHash:       hash[:],
			AmtMsat:           payAmt,
			FinalCltvDelta:    40,
			Dest:              destination[:],
			FeeLimitMsat:      chatMsgAmt * 10,
			TimeoutSeconds:    30,
			DestCustomRecords: customRecords,
		}

		go func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			stream, err := client.SendPayment(ctx, &req)
			if err != nil {
				g.Update(func(g *gocui.Gui) error {
					return err
				})
				return
			}

			for {
				status, err := stream.Recv()
				if err != nil {
					break
				}

				switch status.State {
				case routerrpc.PaymentState_SUCCEEDED:
					msgLines[msgIdx].fee = uint64(status.Route.TotalFeesMsat)
					runningBalance[*destination] -= payAmt
					msgLines[msgIdx].state = stateDelivered
					g.Update(updateView)
					break

				case routerrpc.PaymentState_IN_FLIGHT:

				default:
					msgLines[msgIdx].state = stateFailed
					g.Update(updateView)
					break
				}
			}
		}()

		return nil
	}


	addMsg := func(line chatLine) int {			//받은 메세지가 temp일 경우 위의 온도측정을 하는 함수로 이동해 메세지를 보내게 한다.

		var temp string ="temp"
		if(line.Format()==temp){
		sendreply_func()
		}else{
		}
		msgLines = append(msgLines, line)
		return len(msgLines) - 1
	}

	sendMessage := func(g *gocui.Gui, v *gocui.View) error {	//whatsat에 있는 채팅을 보내는 부분이다.
		if len(v.BufferLines()) == 0 {
			return nil
		}
		newMsg := v.BufferLines()[0]

		v.Clear()
		if err := v.SetCursor(0, 0); err != nil {
			return err
		}
		if err := v.SetOrigin(0, 0); err != nil {
			return err
		}

		if newMsg[0] == '/' {
			destHex := newMsg[1:]
			setDest(destHex)

			updateView(g)

			return nil
		}

		if destination == nil {
			return nil
		}

		d := *destination
		msgIdx := addMsg(chatLine{
			sender:    self,
			text:      newMsg,
			recipient: &d,
		})

		err := updateView(g)
		if err != nil {
			return err
		}

		payAmt := runningBalance[*destination]
		if payAmt < chatMsgAmt {
			payAmt = chatMsgAmt
		}
		if payAmt > 10*chatMsgAmt {
			payAmt = 10 * chatMsgAmt
		}

		var preimage lntypes.Preimage
		if _, err := rand.Read(preimage[:]); err != nil {
			return err
		}
		hash := preimage.Hash()

		// Message sending time stamp
		timestamp := time.Now().UnixNano()
		var timeBuffer [8]byte
		byteOrder.PutUint64(timeBuffer[:], uint64(timestamp))

		// Sign all data.
		signData, err := getSignData(
			self, *destination, timeBuffer[:], []byte(newMsg),
		)
		if err != nil {
			return err
		}

		signResp, err := signClient.SignMessage(context.Background(), &signrpc.SignMessageReq{
			Msg: signData,
			KeyLoc: &signrpc.KeyLocator{
				KeyFamily: int32(keychain.KeyFamilyNodeKey),
				KeyIndex:  0,
			},
		})
		if err != nil {
			return err
		}
		signature := signResp.Signature

		customRecords := map[uint64][]byte{
			tlvMsgRecord:     []byte(newMsg),
			tlvSenderRecord:  self[:],
			tlvTimeRecord:    timeBuffer[:],
			tlvSigRecord:     signature,
			tlvKeySendRecord: preimage[:],
		}

		req := routerrpc.SendPaymentRequest{
			PaymentHash:       hash[:],
			AmtMsat:           payAmt,
			FinalCltvDelta:    40,
			Dest:              destination[:],
			FeeLimitMsat:      chatMsgAmt * 10,
			TimeoutSeconds:    30,
			DestCustomRecords: customRecords,
		}

		go func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			stream, err := client.SendPayment(ctx, &req)
			if err != nil {
				g.Update(func(g *gocui.Gui) error {
					return err
				})
				return
			}

			for {
				status, err := stream.Recv()
				if err != nil {
					break
				}

				switch status.State {
				case routerrpc.PaymentState_SUCCEEDED:
					msgLines[msgIdx].fee = uint64(status.Route.TotalFeesMsat)
					runningBalance[*destination] -= payAmt
					msgLines[msgIdx].state = stateDelivered
					g.Update(updateView)
					break

				case routerrpc.PaymentState_IN_FLIGHT:

				default:
					msgLines[msgIdx].state = stateFailed
					g.Update(updateView)
					break
				}
			}
		}()

		return nil
	}

	err = g.SetKeybinding("send", gocui.KeyEnter, gocui.ModNone, sendMessage)
	if err != nil {
		return err
	}

	go func() {
		returnErr := func(err error) {
			g.Update(func(g *gocui.Gui) error {
				return err
			})
		}

		for {
			invoice, err := stream.Recv()
			if err != nil {
				returnErr(err)
				return
			}

			if invoice.State != lnrpc.Invoice_SETTLED {
				continue
			}

			var customRecords map[uint64][]byte
			for _, htlc := range invoice.Htlcs {
				if htlc.State == lnrpc.InvoiceHTLCState_SETTLED {
					customRecords = htlc.CustomRecords
					break
				}
			}
			if customRecords == nil {
				continue
			}

			msg, ok := customRecords[tlvMsgRecord]
			if !ok {
				continue
			}

			signature, ok := customRecords[tlvSigRecord]
			if !ok {
				continue
			}

			timestampBytes, ok := customRecords[tlvTimeRecord]
			if !ok {
				continue
			}
			timestamp := time.Unix(
				0,
				int64(byteOrder.Uint64(timestampBytes)),
			)

			senderBytes, ok := customRecords[tlvSenderRecord]
			if !ok {
				continue
			}
			sender, err := route.NewVertexFromBytes(senderBytes)
			if err != nil {
				// Invalid sender pubkey
				continue
			}

			signData, err := getSignData(sender, self, timestampBytes, msg)
			if err != nil {
				returnErr(err)
				return
			}

			verifyResp, err := signClient.VerifyMessage(
				context.Background(),
				&signrpc.VerifyMessageReq{
					Msg:       signData,
					Signature: signature,
					Pubkey:    sender[:],
				})
			if err != nil {
				returnErr(err)
				return
			}

			if !verifyResp.Valid {
				continue
			}

			if destination == nil {
				destination = &sender
			}

			addMsg(chatLine{
				sender:    sender,
				text:      string(msg),
				timestamp: timestamp,
			})
			g.Update(updateView)

			amt := invoice.AmtPaid
			runningBalance[*destination] += amt
		}
	}()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
}

func layout(g *gocui.Gui) error {				//채팅창의 ui에 관련된 부분이다. 
	g.Cursor = true

	maxX, maxY := g.Size()
	if v, err := g.SetView("messages", 0, 0, maxX-1, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Messages "
	}

	if v, err := g.SetView("send", 0, maxY-4, maxX-1, maxY-1); err != nil {
		if _, err := g.SetCurrentView("send"); err != nil {
			return err
		}

		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
	}

	updateView(g)

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func updateView(g *gocui.Gui) error {					
	const (
		maxSenderLen = 16
	)

	sendView, _ := g.View("send")
	if destination == nil {
		sendView.Title = " Set a destination by typing /pubkey "
	} else {
		alias := keyToAlias[*destination]
		sendView.Title = fmt.Sprintf(" Send to %v [balance: %v msat]",
			alias, runningBalance[*destination])
	}

	messagesView, _ := g.View("messages")

	messagesView.Clear()
	cols, rows := messagesView.Size()

	startLine := len(msgLines) - rows
	if startLine < 0 {
		startLine = 0
	}

	for _, line := range msgLines[startLine:] {
		text := line.text
		var r string
		if line.recipient != nil {
			r = keyToAlias[*line.recipient]
		} else {
			r = fmt.Sprintf("sent: %v",
				line.timestamp.Format(time.ANSIC))
		}




		text += fmt.Sprintf(" \x1b[34m(%v)\x1b[0m", r)

		var amtDisplay string
		if line.state == stateDelivered {
			amtDisplay = formatMsat(line.fee)
		}

		maxTextFieldLen := cols - len(amtDisplay) - maxSenderLen + 5
		maxTextLen := maxTextFieldLen
		if line.state != statePending {
			maxTextLen -= 2
		}
		if len(text) > maxTextLen {
			text = text[:maxTextLen-3] + "..."
		}
		paddingLen := maxTextFieldLen - len(text)
		switch line.state {
		case stateDelivered:
			text += " \x1b[34m✔️\x1b[0m"
			paddingLen -= 2
		case stateFailed:
			text += " \x1b[31m✘\x1b[0m"
			paddingLen -= 2
		}

		text += strings.Repeat(" ", paddingLen)

		senderAlias := keyToAlias[line.sender]
		if len(senderAlias) > maxSenderLen {
			senderAlias = senderAlias[:maxSenderLen]
		}
		fmt.Fprintf(messagesView, "%16v: %v \x1b[34m%v\x1b[0m",
			senderAlias,
			text, amtDisplay,
		)

		fmt.Fprintln(messagesView)
	}
	return nil
}

func formatMsat(msat uint64) string {
	wholeSats := msat / 1000
	msats := msat % 1000
	var msatsStr string
	if msats > 0 {
		msatsStr = fmt.Sprintf(".%03d", msats)
		msatsStr = strings.TrimRight(msatsStr, "0")
	}
	return fmt.Sprintf("[%d%-4s sat]",
		wholeSats, msatsStr,
	)
}

func getSignData(sender, recipient route.Vertex, timestamp []byte, msg []byte) ([]byte, error) {
	var signData bytes.Buffer

	// Write sender.
	if _, err := signData.Write(sender[:]); err != nil {
		return nil, err
	}

	// Write recipient.
	if _, err := signData.Write(recipient[:]); err != nil {
		return nil, err
	}

	// Write time.
	if _, err := signData.Write(timestamp); err != nil {
		return nil, err
	}

	// Write message.
	if _, err := signData.Write(msg); err != nil {
		return nil, err
	}

	return signData.Bytes(), nil
}
