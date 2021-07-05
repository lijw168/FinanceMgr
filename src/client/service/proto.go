package service

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

const (
	FinanceMgrMagic = 0xF8
	HeaderLength    = 8
	Int32BytesCount = 4
)

//the message structure definition
type Packet struct {
	//header
	Magic    uint8
	Version  uint8
	OpCode   uint8
	Reserved uint8
	Size     int32
	//body
	Buf      []byte
	RecvTime time.Time
}

func NewPacket() *Packet {
	pkt := new(Packet)
	pkt.Magic = FinanceMgrMagic
	pkt.Version = 0x01
	return pkt
}

func (p *Packet) String() string {
	strRecvTime := p.RecvTime.String()
	strBuf := fmt.Sprintf("{Magic:%d,Version:%d,OpCode:%d,Size:%d,receive time:%s,Buf:%s}\r\n",
		p.Magic, p.Version, p.OpCode, p.Size, strRecvTime, string(p.Buf))
	return strBuf
}

func (p *Packet) marshalHeader() (out []byte) {
	out = make([]byte, HeaderLength)
	out[0] = p.Magic
	out[1] = p.Version
	out[2] = p.OpCode
	out[3] = p.Reserved
	p.Size = (int32)(len(p.Buf))
	binary.LittleEndian.PutUint32(out[4:4+Int32BytesCount], uint32(p.Size))
	return
}

func (p *Packet) unmarshalHeader(in []byte) (err error) {
	p.Magic = in[0]
	if p.Magic != FinanceMgrMagic {
		return errors.New("BadMagic:" + strconv.Itoa(int(p.Magic)))
	}
	p.Version = in[1]
	p.OpCode = in[2]
	p.Reserved = in[3]
	//little endian
	p.Size = int32(binary.LittleEndian.Uint32(in[4 : 4+Int32BytesCount]))
	return
}

func (p *Packet) WriteToConn(c net.Conn) (err error) {
	//firstly the header
	bufHdr := p.marshalHeader()
	if _, err = c.Write(bufHdr); err != nil {
		return err
	}
	if p.Size != 0 {
		//then the body
		if _, err = c.Write(p.Buf); err != nil {
			return err
		}
	}
	return
}

func (p *Packet) ReadFromConn(c net.Conn) (err error) {
	bufHrd := make([]byte, HeaderLength)
	p.RecvTime = time.Now()
	if _, err = io.ReadFull(c, bufHrd); err != nil {
		return err
	}

	if err = p.unmarshalHeader(bufHrd); err != nil {
		return err
	}
	if p.Size != 0 {
		p.Buf = make([]byte, p.Size)
		if _, err = io.ReadFull(c, p.Buf); err != nil {
			return err
		}
	}
	return
}
