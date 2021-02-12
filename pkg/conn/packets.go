package conn

import (
	"fmt"
	"github.com/viamAhmadi/graphBroker/pkg/util"
	"strconv"
	"strings"
)

type Message struct {
	Type        string
	Forward     byte
	Destination string
	Id          int    // 4 bytes
	Sign        string // 4 bytes
	Content     string // 8 kiloBytes
}

type Done struct {
	Type        string
	Forward     byte
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
	Successful  byte
	List        []string
}

type Error struct {
	Msg         string
	Destination []byte
}

func ConvertToMessage(b *[]byte) (*Message, error) {
	if cap(*b) < 10 {
		return nil, ErrConvertToModel
	}
	i, err := strconv.Atoi(util.RemoveAdditionalCharacters((*b)[1:5]))
	if err != nil {
		return nil, ErrConvertToModel
	}
	return &Message{
		Type:        string((*b)[0]),
		Id:          i,
		Sign:        util.RemoveAdditionalCharacters((*b)[5:9]),
		Destination: util.RemoveAdditionalCharacters((*b)[9:31]),
		Forward:     (*b)[31],
		Content:     string((*b)[32:]),
	}, nil
}

func SerializeMessage(id int, forward byte, sign, destination string, content *string) *[]byte {
	s, _ := strconv.Atoi(sign)
	v := []byte(fmt.Sprintf("m%s%s%s%s%s", util.ConvertIntToBytes(id), util.ConvertIntToBytes(s), util.ConvertDesToBytes(destination), string(forward), *content))
	return &v
}

func (m *Message) GetConnId() string {
	if m.Forward == YES {
		return m.Destination + m.Sign + "s"
	}
	return m.Destination + m.Sign
}

func (m *Message) GetId() string {
	return strconv.Itoa(m.Id) + m.Sign
}

func ConvertToDone(b []byte) (Done, error) {
	if len(b) < 21 {
		return Done{}, ErrConvertToModel
	}

	c, err := strconv.Atoi(util.RemoveAdditionalCharacters(b[27:31]))
	if err != nil {
		return Done{}, err
	}
	return Done{
		Type:        string(b[0]),
		Destination: util.RemoveAdditionalCharacters(b[1:23]),
		Sign:        util.RemoveAdditionalCharacters(b[23:27]),
		Count:       c,
		Forward:     b[31],
	}, nil
}

func SerializeDone(destination, sign string, forward byte, count int) []byte {
	s, _ := strconv.Atoi(sign)
	return []byte(fmt.Sprintf("d%s%s%s%s", util.ConvertDesToBytes(destination), util.ConvertIntToBytes(s), util.ConvertIntToBytes(count), string(forward)))
}

func (d *Done) GetConnId() string {
	if d.Forward == YES {
		return d.Destination + d.Sign + "s"
	}
	return d.Destination + d.Sign
}

func ConvertToSend(b []byte) (Send, error) {
	if len(b) < 25 {
		return Send{}, ErrConvertToModel
	}
	return Send{
		Type:        string(b[0]),
		Destination: util.RemoveAdditionalCharacters(b[1:23]),
		Sign:        util.RemoveAdditionalCharacters(b[23:27]),
	}, nil
}

func SerializeSend(destination, sign string) []byte {
	s, _ := strconv.Atoi(sign)
	return []byte(fmt.Sprintf("s%s%s", util.ConvertDesToBytes(destination), util.ConvertIntToBytes(s)))
}

func ConvertToFactor(b *[]byte) (*Factor, error) {
	if len(*b) < 16 {
		return nil, ErrConvertToModel
	}
	//status := string((*b)[27:28])
	successful := (*b)[27]
	var list []string
	if successful != YES {
		nums := strings.Split(string((*b)[28:]), ".")
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
		Sign:        util.RemoveAdditionalCharacters((*b)[23:27]),
		//Status:      status,
		Successful: successful,
		List:       list,
	}, nil
}

func SerializeFactor(destination, sign string, successful byte, list *[]string) *[]byte {
	tmp := ""
	if successful != YES {
		if list != nil {
			tmp = strings.Join(*list, ".")
		}
	}
	s, _ := strconv.Atoi(sign)
	b := []byte(fmt.Sprintf("f%s%s%s%s", util.ConvertDesToBytes(destination), util.ConvertIntToBytes(s), string(successful), tmp))
	return &b
}
