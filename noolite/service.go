package noolite

import (
	"github.com/jacobsa/go-serial/serial"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"sync"
	"time"
	"wb-noolite-mtrf/config"
)

// Service Сервис для общения с MTRF устройствами. Прием запросов и отправка их в адаптер
type Service struct {
	log                   *logrus.Logger
	config                *config.Config
	port                  io.ReadWriteCloser
	serialOptions         serial.OpenOptions
	sndQueue              chan *Request
	rcvQueue              chan *Response
	closeOnce             sync.Once
	exit                  chan struct{}
	blockSend             chan struct{}
	closeWait             sync.WaitGroup
	sendRequestsOnConnect []*Request
	isWaitingResponse     bool
}

// NewNooliteService Создает сервис для работы адаптером MTRF
func NewNooliteService(log *logrus.Logger, config *config.Config, initialRequests []*Request) (*Service, error) {

	s := &Service{
		log:    log,
		config: config,
		serialOptions: serial.OpenOptions{
			PortName:        config.SerialPort,
			BaudRate:        9600,
			StopBits:        1,
			DataBits:        8,
			MinimumReadSize: PacketLen,
		},
		sndQueue:              make(chan *Request, config.QueueLen),
		rcvQueue:              make(chan *Response, config.QueueLen),
		blockSend:             make(chan struct{}, 1),
		sendRequestsOnConnect: initialRequests,
	}
	go s.worker()
	go s.connectManager()
	s.log.Info("[MTRF] Noolite service successful initialization")
	return s, nil

}

// Close Закрывает порт, закрывает каналы
func (s *Service) Close() {
	s.closeOnce.Do(func() {

		err := s.port.Close()
		if err != nil {
			s.log.Errorf("[MTRF] Error on close service: %s", err)
			return
		}
		close(s.exit)
	})
}

func (s *Service) connectManager() {
	for {
		select {
		case <-s.exit:
			return
		default:

		}
		s.closeWait.Wait()
		s.log.Info("[MTRF] Starting connection")
		go s.worker()
		time.Sleep(time.Millisecond * 500)
	}
}

func (s *Service) worker() {
	defer s.closeWait.Done()
	fClose := func() {
		s.log.Error("[MTRF] Closing worker")
	}
	s.closeWait.Add(1)
	for {
		select {
		case <-s.exit:
			fClose()
			return
		default:
		}
		var err error
		if s.port != nil {
			_ = s.port.Close()
		}

		s.port, err = serial.Open(s.serialOptions)
		if err != nil {
			s.log.Errorf("[MTRF] Serial open error: %s", err)
			select {
			case <-s.exit:
				fClose()
				return
			default:
			}
			time.Sleep(time.Second * 5)
			continue
		}

		s.log.Info("[MTRF] Successfully connected.")
		go s.sendRequests()  // Read send requests queue and send it to adapter
		go s.readResponses() // Read responses from adapter and send it to queue
		time.Sleep(time.Millisecond * 100)
		for _, request := range s.sendRequestsOnConnect {
			s.sndQueue <- request
		}
		return
	}
}

func (s *Service) sendRequests() {
	s.log.Info("[MTRF] Starting Requests Queue reader")
	defer s.closeWait.Done()
	s.closeWait.Add(1)
	for {
		select {
		case <-s.exit:
			s.log.Info("[MTRF] Closing requests queue reader")
			return
		case request := <-s.sndQueue:
			{

				if request == nil {
					return
				}

				s.log.Debugf("[MTRF] Send request %s", request.String())
				_, err := s.port.Write(request.BuildBytes())
				s.isWaitingResponse = request.WaitResponse
				if err != nil {
					s.log.Errorf("Error write serial port: %s", err)
					return
				}
				s.log.Trace("[MTRF] Request successfully sent")

				select {
				case <-s.blockSend:
					s.log.Trace("[MTRF] Response received. Ready to send next command")
				case <-time.After(30 * time.Second):
					s.log.Error("[MTRF] No response. Timeout")
				}
				s.isWaitingResponse = false

			}
		}

	}
}

func (s *Service) readResponses() {
	var buf [17]byte
	defer s.closeWait.Done()
	s.closeWait.Add(1)
	s.log.Info("[MTRF] Starting Response Queue producer")
	for {
		_, err := io.ReadAtLeast(s.port, buf[:], PacketLen)

		if s.isWaitingResponse {
			s.log.Trace("[MTRF] Unblock sending new commands.")
			s.blockSend <- struct{}{}
		}

		if err != nil {
			s.log.Errorf("[MTRF] Error while read from serial port: %s", err)
			s.sndQueue <- &Request{}
			s.log.Info("[MTRF] Closing response queue producer")
			return
		}
		response := Response{}
		err = response.Parse(buf[:])
		if err != nil {
			s.log.Errorf("[MTRF] Receive error for %x:%s", buf, err) // Only logging
		} else {
			s.log.Tracef("[MTRF] Receive response %s: ", response.String())
			select {
			case s.rcvQueue <- &response:
			default:
				log.Println("[MTRF] Receive queue is full, skipped")
			}
		}
	}
}
func (s *Service) Receive() <-chan *Response {
	return s.rcvQueue
}

func (s *Service) Send() chan<- *Request {
	return s.sndQueue
}
