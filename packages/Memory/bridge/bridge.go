package bridge

import (
	"encoding/json"
	"fmt"
	"main/packages/Memory/instance"
	Rbx "main/packages/Memory/instance"
	"main/packages/Memory/memory"
	"reflect"
	"regexp"
	"strconv"
	"sync"
	"syscall"
	"time"
)

const dataMaxLen = 199998
const payloadMatch = "^[A-Fa-f0-9]{8}"

var PEER_TYPE = map[string]int{
	"Roblox":   0,
	"External": 1,
}

var SENDER_TYPE = map[string]int{
	"R2E": 0, // Roblox to external
	"E2R": 1, // External to Roblox
}

func extractBits(value, field, width int) int {
	return (value >> field) & width
}

type BridgeChannel struct {
	rbx           *Rbx.RobloxInstances
	mem           *memory.Luna
	Handle        syscall.Handle
	Name          string
	States        Rbx.Instance
	Peer0         Rbx.Instance
	Peer1         Rbx.Instance
	InstanceRefs  Rbx.Instance
	BuffersCaches map[int]map[int]Rbx.Instance
}

func NewBridgeChannel(handle syscall.Handle, rbx *instance.RobloxInstances) *BridgeChannel {
	return &BridgeChannel{
		rbx:    rbx,
		mem:    rbx.Mem,
		Handle: handle,
		BuffersCaches: map[int]map[int]Rbx.Instance{
			0: {},
			1: {},
		},
	}
}

func (bc *BridgeChannel) Initialize(channelContainer Rbx.Instance) {
	bc.Name = channelContainer.Name()
	bc.States = channelContainer.WaitForChild("States", 1)
	bc.Peer0 = channelContainer.WaitForChild("Peer0", 1)
	bc.Peer1 = channelContainer.WaitForChild("Peer1", 1)
	bc.InstanceRefs = channelContainer.WaitForChild("InstanceRefs", 1)
}

func (bc *BridgeChannel) GetChannelStates() (bool, bool, bool, *int) {
	if bc.States.Address == 0 || bc.States.Value() == nil {
		return false, false, false, nil
	}

	data := bc.States.Value()
	if data == nil || reflect.TypeOf(data).String() == "string" {
		return false, false, false, nil
	}

	packedValue := int(bc.States.Value().(float64))

	isUsed := extractBits(packedValue, 0, 1) == 1
	responding := extractBits(packedValue, 1, 1) == 1
	responded := extractBits(packedValue, 2, 1) == 1
	sender := extractBits(packedValue, 3, 1)

	var senderPtr *int
	if isUsed {
		senderPtr = &sender
	}

	return isUsed, responding, responded, senderPtr
}

var lock sync.Mutex

func (bc *BridgeChannel) SetChannelStates(isUsed, responding, responded bool, sender int) {
	lock.Lock()
	defer lock.Unlock()
	if bc.States.Address < 1000 {
		return
	}
	result := 0
	if isUsed {
		result |= 0b0001
	}
	if responding {
		result |= 0b0010
	}
	if responded {
		result |= 0b0100
	}
	result |= (sender << 3) & 0b1000
	bc.States.SetValue(result)
}

func (bc *BridgeChannel) GetBufferData(containerType int) string {
	var container Rbx.Instance
	if containerType == 0 {
		container = bc.Peer0
	} else if containerType == 1 {
		container = bc.Peer1
	} else {
		return ""
	}

	buffersCache := bc.BuffersCaches[containerType]

	c := container.GetChildren()
	childrenCount := len(c)
	result := ""

	for bufferIdx := 0; bufferIdx < childrenCount; bufferIdx++ {
		var bufferObj Rbx.Instance
		if val, ok := buffersCache[bufferIdx]; ok {
			bufferObj = val
		} else {
			bufferObj = container.FindFirstChild(strconv.Itoa(bufferIdx), true)
			if bufferObj.Address > 1000 {
				buffersCache[bufferIdx] = bufferObj
			}
		}
		if bufferObj.Address > 1000 {
			if bufferObj.Value() != nil {
				result += bufferObj.Value().(string)
			}
		}
	}

	if len(result) > 0 {
		re := regexp.MustCompile(payloadMatch)
		bufferSizeMatch := re.FindStringSubmatch(result)
		if len(bufferSizeMatch) > 0 {
			matchLen := len(bufferSizeMatch[0])
			bufferSize, err := strconv.ParseInt(bufferSizeMatch[0], 16, 64)
			if err != nil {
				return ""
			}
			startIdx := matchLen + 1
			endIdx := int(bufferSize) + matchLen + ((int(bufferSize) / dataMaxLen) + 1)
			if endIdx > len(result) {
				endIdx = len(result)
			}
			return result[startIdx:endIdx]
		} else {
			return ""
		}
	}
	return ""
}

func (bc *BridgeChannel) SetBufferData(newData string) bool {
	buffersCache := bc.BuffersCaches[PEER_TYPE["External"]]

	for bufferPos := 0; bufferPos < len(newData); bufferPos += dataMaxLen {
		bufferIdx := bufferPos / dataMaxLen

		var bufferObj Rbx.Instance
		if val, ok := buffersCache[bufferIdx]; ok {
			bufferObj = val
		} else {
			bufferObj = bc.Peer1.FindFirstChild(strconv.Itoa(bufferIdx), true)
			if bufferObj.Address > 1000 {
				buffersCache[bufferIdx] = bufferObj
			} else {
				return false
			}
		}

		endPos := bufferPos + dataMaxLen
		if endPos > len(newData) {
			endPos = len(newData)
		}

		bufferObj.SetValue(newData[bufferPos:endPos])

	}

	return true
}

type Bridge struct {
	rbx               *instance.RobloxInstances
	mem               *memory.Luna
	Handle            syscall.Handle
	Channels          []*BridgeChannel
	Sessions          map[string]int
	QueuedDatas       []string
	CallbacksRegistry map[string]func(int, []interface{}) []interface{}
	RobloxTerminated  bool
	MainContainer     Rbx.Instance
	ModuleHolder      Rbx.Instance
	mutex             sync.Mutex
}

func NewBridge(mem *instance.RobloxInstances) *Bridge {
	b := &Bridge{
		rbx:               mem,
		mem:               mem.Mem,
		Handle:            mem.Mem.ProcessHandle,
		Channels:          []*BridgeChannel{},
		Sessions:          make(map[string]int),
		QueuedDatas:       []string{},
		CallbacksRegistry: make(map[string]func(int, []interface{}) []interface{}),
		RobloxTerminated:  false,
	}

	go b.bridgeListener()
	go b.bridgeQueueSched()

	return b
}

func (b *Bridge) Start(newPID int, mainContainer Rbx.Instance) {
	b.MainContainer = mainContainer
	b.ModuleHolder = b.MainContainer.WaitForChild("ModuleHolder", 5)
	channels := b.MainContainer.WaitForChild("Channels", 5)

	if b.ModuleHolder.Address < 1000 || channels.Address < 1000 {
		return
	}

	b.Channels = []*BridgeChannel{}
	b.Sessions = make(map[string]int)
	b.QueuedDatas = []string{}
	b.CallbacksRegistry = make(map[string]func(int, []interface{}) []interface{})

	for channelIdx := 0; channelIdx < 8; channelIdx++ {
		channelContainer := channels.FindFirstChild(strconv.Itoa(channelIdx), false)
		if channelContainer.Address < 1000 {
			continue
		}
		channelObj := NewBridgeChannel(b.Handle, b.rbx)
		channelObj.Initialize(channelContainer)
		b.Channels = append(b.Channels, channelObj)

		time.Sleep(time.Millisecond * 50)
	}

	b.RobloxTerminated = false
	//go b.processHandler(newPID)
}

var last time.Time

func (b *Bridge) Send(action string, args []interface{}) {

	if time.Since(last) <= time.Millisecond*50 {
		time.Sleep(time.Millisecond * 50)
	}

	last = time.Now()

	if b.RobloxTerminated {
		return
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	session := b.Sessions[action]
	payload, err := processData(action, session, args)
	if err != nil {
		return
	}

	b.QueuedDatas = append(b.QueuedDatas, payload)
	b.Sessions[action] = session + 1
}

func (b *Bridge) RegisterCallback(callbackName string, callback func(int, []interface{}) []interface{}) {
	b.CallbacksRegistry[callbackName] = callback
}

func processData(action string, session int, args []interface{}) (string, error) {
	data := []interface{}{action, session, args}
	result, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	dataLenHex := fmt.Sprintf("%08x", len(result))
	return fmt.Sprintf("%s|%s", dataLenHex, string(result)), nil
}

func handleCallback(cbname string, channel *BridgeChannel, callback func(int, []interface{}) []interface{}, session int, args []interface{}) {
	returnedArgs := callback(session, args)
	payload, err := processData(cbname, session, returnedArgs)
	if err != nil {
		return
	}
	setSuccess := channel.SetBufferData(payload)
	channel.SetChannelStates(setSuccess, false, true, SENDER_TYPE["E2R"])
}

func (b *Bridge) bridgeQueueSched() {
	for {
		time.Sleep(1 * time.Millisecond)

		if b.RobloxTerminated || len(b.QueuedDatas) == 0 {
			continue
		}

		channel := b.getAvailableChannel()
		if channel == nil {
			continue
		}

		b.mutex.Lock()
		payload := b.QueuedDatas[0]
		b.QueuedDatas = b.QueuedDatas[1:]
		b.mutex.Unlock()

		setSuccess := channel.SetBufferData(payload)

		channel.SetChannelStates(setSuccess, false, false, SENDER_TYPE["E2R"])
	}
}

func (b *Bridge) bridgeListener() {
	for {
		time.Sleep(1 * time.Millisecond)
		if b.RobloxTerminated {
			break
		}
		for _, channel := range b.Channels {
			isUsed, _, _, sender := channel.GetChannelStates()
			if sender != nil && *sender == SENDER_TYPE["E2R"] || !isUsed {
				continue
			}

			rawData := channel.GetBufferData(PEER_TYPE["Roblox"])

			if rawData == "" {
				channel.SetChannelStates(false, false, false, SENDER_TYPE["E2R"])
				continue
			}

			var receivedData []interface{}
			err := json.Unmarshal([]byte(rawData), &receivedData)
			if err != nil {
				channel.SetChannelStates(false, false, false, SENDER_TYPE["E2R"])
				continue
			}
			if len(receivedData) < 3 {
				channel.SetChannelStates(false, false, false, SENDER_TYPE["E2R"])
				continue
			}
			action, ok := receivedData[0].(string)
			if !ok {
				channel.SetChannelStates(false, false, false, SENDER_TYPE["E2R"])
				continue
			}
			sessionFloat, ok := receivedData[1].(float64)
			if !ok {
				channel.SetChannelStates(false, false, false, SENDER_TYPE["E2R"])
				continue
			}
			session := int(sessionFloat)
			rawArgs, ok := receivedData[2].([]interface{})
			if !ok {
				channel.SetChannelStates(false, false, false, SENDER_TYPE["E2R"])
				continue
			}

			callback, exists := b.CallbacksRegistry[action]

			if exists {
				actionArgs := []interface{}{}
				for _, valueInfo := range rawArgs {
					valueSlice, ok := valueInfo.([]interface{})
					if !ok || len(valueSlice) < 2 {
						continue
					}
					valueType, _ := valueSlice[0].(string)
					value := valueSlice[1]
					if valueType == "Instance" {
						v := channel.InstanceRefs.FindFirstChild(fmt.Sprintf("%v", value), true)
						value = v.Value()
					} else if valueType == "table" {
						valueStr, ok := value.(string)
						if ok {
							var tableValue interface{}
							err := json.Unmarshal([]byte(valueStr), &tableValue)
							if err == nil {
								value = tableValue
							}
						}
					}
					actionArgs = append(actionArgs, value)
				}
				handleCallback(action, channel, callback, session, actionArgs)
				channel.SetChannelStates(true, true, false, SENDER_TYPE["E2R"])
			} else {
				channel.SetChannelStates(false, false, false, SENDER_TYPE["E2R"])
			}
		}
	}
}

func (b *Bridge) getAvailableChannel() *BridgeChannel {
	for _, channel := range b.Channels {
		isUsed, _, _, _ := channel.GetChannelStates()
		if !isUsed {
			return channel
		}
	}
	return nil
}
