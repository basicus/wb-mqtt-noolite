package mqtt

// Message MQTT Сообщение
type Message struct {
	Topic   string
	Retain  bool
	Payload string
}

// Packet Пакет сообщений
type Packet struct {
	messages []*Message
}

// NewPacket Формирование нового пакета из сообщений
func NewPacket(messages ...*Message) *Packet {
	return &Packet{
		messages: messages,
	}
}

// Add Добавление сообщения в пакет
func (p *Packet) Add(m ...*Message) {
	p.messages = append(p.messages, m...)
}
