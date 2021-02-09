package conn

import (
	"fmt"
	"github.com/viamAhmadi/graphBroker/pkg/util"
	"strconv"
	"strings"
)

type Connection struct {
	Type                 string // Type of request, c=connection or m=message ..
	Destination          string
	Sign                 string
	Count                int
	FirstMsgId, EndMsgId int
}

type Message struct {
	Type    string
	Id      int    // 4 bytes
	Sign    string // 2 bytes
	Content string // 8 kiloBytes
}

type Done struct {
	Type        string
	Destination string
	Sign        string
	Count       int
}

type Send struct {
	Type        string
	Destination string
	Sign        string
}

type Factor struct {
	Type        string
	Destination string
	Sign        string
	Status      string
	List        []string
}

func ConvertToConnection(b []byte) (Connection, error) {
	if cap(b) < 39 {
		return Connection{}, ErrConvertToModel
	}
	count, err := strconv.Atoi(string(b[25:29]))
	if err != nil {
		return Connection{}, err
	}
	firstMsgId, err := strconv.Atoi(string(b[29:33]))
	if err != nil {
		return Connection{}, err
	}
	endMsgId, err := strconv.Atoi(string(b[33:37]))
	if err != nil {
		return Connection{}, err
	}
	return Connection{Type: string(b[0]),
		Destination: util.RemoveAdditionalCharacters(b[1:23]),
		Sign:        string(b[23:25]),
		Count:       count,
		FirstMsgId:  firstMsgId,
		EndMsgId:    endMsgId,
	}, nil
}

func SerializeConnection(destination, sign string, count, firstMsgId, endMsgId int) []byte {
	return []byte(fmt.Sprintf("c%s%s%d%d%d", util.ConvertDesToBytes(destination), sign, count, firstMsgId, endMsgId))
}

func ConvertToMessage(b *[]byte) (*Message, error) {
	if cap(*b) < 8 {
		return nil, ErrConvertToModel
	}
	i, err := strconv.Atoi(util.RemoveAdditionalCharacters((*b)[1:5]))
	if err != nil {
		return nil, ErrConvertToModel
	}
	return &Message{
		Type:    string((*b)[0]),
		Id:      i,
		Sign:    string((*b)[5:7]),
		Content: string((*b)[7:]),
	}, nil
}

func SerializeMessage(id int, sign, content string) *[]byte {
	v := []byte(fmt.Sprintf("m%s%s%s", util.ConvertIdToBytes(id), sign, content))
	return &v
}

func ConvertToDone(b []byte) (Done, error) {
	if len(b) < 19 {
		return Done{}, ErrConvertToModel
	}

	c, err := strconv.Atoi(util.RemoveAdditionalCharacters(b[25:29]))
	if err != nil {
		return Done{}, err
	}
	return Done{
		Type:        string(b[0]),
		Destination: util.RemoveAdditionalCharacters(b[1:23]),
		Sign:        string(b[23:25]),
		Count:       c,
	}, nil
}

func SerializeDone(destination, sign string, count int) []byte {
	return []byte(fmt.Sprintf("d%s%s%s", util.ConvertDesToBytes(destination), sign, util.ConvertIdToBytes(count)))
}

func ConvertToSend(b []byte) (Send, error) {
	if len(b) != 25 {
		return Send{}, ErrConvertToModel
	}
	return Send{
		Type:        string(b[0]),
		Destination: util.RemoveAdditionalCharacters(b[1:23]),
		Sign:        string(b[23:25]),
	}, nil
}

func SerializeSend(destination, sign string) []byte {
	return []byte(fmt.Sprintf("s%s%s", util.ConvertDesToBytes(destination), sign))
}

func ConvertToFactor(b *[]byte) (*Factor, error) {
	if len(*b) < 16 {
		return nil, ErrConvertToModel
	}
	status := string((*b)[25:26])
	var list []string
	if status != ok {
		nums := strings.Split(string((*b)[26:]), ".")
		for i := 0; i < len(nums); i++ {
			val := nums[i]
			if val == "" {
				continue
			}
			list = append(list, val)
		}
	}
	return &Factor{
		Type:        string((*b)[0]),
		Destination: util.RemoveAdditionalCharacters((*b)[1:23]),
		Sign:        string((*b)[23:25]),
		Status:      status,
		List:        list,
	}, nil
}

func SerializeFactor(destination, sign, status string, list *[]string) *[]byte {
	tmp := ""
	if status != ok {
		if list != nil {
			tmp = strings.Join(*list, ".")
		}
	}
	b := []byte(fmt.Sprintf("f%s%s%s%s", util.ConvertDesToBytes(destination), sign, status, tmp))
	return &b
}
