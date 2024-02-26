package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	"github.com/Vladimir77715/my-tcp-pow/internal/config"
	"github.com/Vladimir77715/my-tcp-pow/internal/errror"
	"github.com/Vladimir77715/my-tcp-pow/internal/parser"
	"github.com/Vladimir77715/my-tcp-pow/internal/tcp/server/handlers"
)

type SecretKey string

var SecretKeyVal SecretKey = "secret_key"

var ErrSecretKeyNotFound = errors.New("secret key not found")

type Server struct {
	cfg      *config.ServerConfig
	handlers map[int]handlers.Handler
	secret   []byte
}

const (
	Quit int = iota
	RequestChallenge
	SendSolution
	WisdomRequest
)

var (
	ErrServerStop       = errors.New("server has been stopped")
	ErrWrongCommandType = errors.New("wrong command type")
)

func New(cfg *config.ServerConfig, handlers map[int]handlers.Handler) *Server {
	bUid, _ := uuid.New().MarshalBinary()
	return &Server{cfg: cfg, handlers: handlers, secret: bUid}
}

func (s *Server) Start(ctx context.Context) error {
	const methodName = "Start"
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.cfg.Address, s.cfg.Port))
	if err != nil {
		return errror.WrapError(err, methodName)
	}
	defer listener.Close()

	// Создание канала для приема сигналов
	sysChan := make(chan os.Signal, 1)

	// Настройка Notify для корректного завершения сервера
	signal.Notify(sysChan, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGSTOP)

	for {
		select {
		case <-sysChan:
			return ErrServerStop
		default:
			conn, err := listener.Accept()
			if err != nil {
				return errror.WrapError(err, methodName)
			}
			go s.handleConn(ctx, conn)
		}
	}
}

func (s *Server) handleConn(ctx context.Context, conn net.Conn) {
	const methodName = "handleConn"

	defer conn.Close()

	rawData, err := parser.Encode(conn)
	if err != nil {
		_, err = conn.Write([]byte(errror.WrapError(err, methodName).Error()))
		if err != nil {
			log.Print(err.Error())
		}

		return
	}

	if rawData.Command < 0 || rawData.Command > WisdomRequest {
		_, err = conn.Write([]byte(errror.WrapError(ErrWrongCommandType, methodName).Error()))
		if err != nil {
			log.Print(err.Error())
		}

		return
	}

	if rawData.Command == Quit {
		return
	}

	if h, ok := s.handlers[rawData.Command]; ok {
		out, handleErr := h.Handle(context.WithValue(ctx, SecretKeyVal, s.secret), rawData)
		if handleErr != nil {
			conn.Write([]byte(errror.WrapError(err, methodName).Error()))
			return
		}

		_, err = conn.Write(parser.Decode(out))
		if err != nil {
			log.Print(err.Error())
		}

	} else {
		_, err = conn.Write([]byte(errror.WrapError(ErrWrongCommandType, methodName).Error()))
		if err != nil {
			log.Print(err.Error())
		}

		return
	}
}

func (h *Server) generateToken() (string, *time.Time, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	expire := time.Now().Add(time.Minute * 5)
	// Создание токена с Claims (утверждениями)
	claims := token.Claims.(jwt.MapClaims)
	claims["admin"] = true
	claims["exp"] = expire.Unix()

	// Подписание токена нашим секретным ключом
	tokenString, err := token.SignedString(h.secret)
	if err != nil {
		return "", nil, err
	}

	return tokenString, &expire, nil
}
