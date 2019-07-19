package dns

import (
	"github.com/spf13/viper"
	"golang.org/x/net/dns/dnsmessage"
	"log"
	"net"
	"strings"

	"../burst"
)

const (
	dnsPacketLen uint16 = 512
	tld          string = "burst."
)

type request struct {
	Host string
	Type string
	Data string
	TTL  uint32
}

type BurstDNS struct {
	conn *net.UDPConn
}

func isBurstTLD(message dnsmessage.Message) bool {
	return strings.HasSuffix(message.Questions[0].Name.String(), tld)
}

func sendPacket(conn *net.UDPConn, addr net.UDPAddr, message dnsmessage.Message) {
	packed, err := message.Pack()
	if err != nil {
		log.Println(err)
		return
	}

	_, err = conn.WriteToUDP(packed, &addr)
	if err != nil {
		log.Println(err)
	}
}

func (s *BurstDNS) answer(addr net.UDPAddr, m dnsmessage.Message) {
	// remove .TLD.
	domain := m.Questions[0].Name.String()
	aliasName := domain[:len(domain)-7]

	if records, err := burst.GetRecords(aliasName); err == nil {
		for _, record := range records {
			req := request{
				Host: domain,
				Type: record.Type,
				Data: record.Data,
				TTL:  record.TTL,
			}

			res, err := ToResource(req)
			if err != nil {
				log.Println(err)
				continue
			}

			m.Answers = append(m.Answers, res)
		}
	}

	go sendPacket(s.conn, addr, m)
}

func (s *BurstDNS) Listen() {
	var err error
	udpPort := viper.GetInt("dns.port")
	s.conn, err = net.ListenUDP("udp", &net.UDPAddr{Port: udpPort})
	if err != nil {
		log.Fatal(err)
	}
	defer s.conn.Close()

	log.Printf("Listening udp :%d\n", udpPort)

	for {
		buf := make([]byte, dnsPacketLen)

		_, addr, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("Error ReadFromUDP: %s\n", err)
			continue
		}

		var m dnsmessage.Message

		err = m.Unpack(buf)
		if err != nil {
			log.Printf("Error Unpack message: %s\n", err)
			continue
		}
		if len(m.Questions) != 1 {
			log.Printf("Must be one question: %s", m.Questions)
			continue
		}

		if !isBurstTLD(m) {
			continue
		}

		log.Printf("Got request %v, %v from: %s\n", m, m.Header, addr)

		go s.answer(*addr, m)
	}
}
