package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Vladimir77715/my-tcp-pow/internal/errror"
	"github.com/Vladimir77715/my-tcp-pow/internal/models"
	"github.com/Vladimir77715/my-tcp-pow/internal/reader"
)

//go:generate mockery --name Reader --with-expecter

type Reader interface {
	Read(p []byte) (n int, err error)
}

var ErrWrongCommandFormat = errors.New("wrong command format")

const (
	EndLine = "END"
)

const (
	CommandBatchSize = 2
	PayloadBatchSize = 256
)

func Decode(out *models.RawData) []byte {
	if out == nil {
		return nil
	}

	if out.Payload == nil {
		return []byte(fmt.Sprintf("%d\n%s", out.Command, EndLine))
	}

	return []byte(fmt.Sprintf("%d\n%v\n%s", out.Command, string(out.Payload), EndLine))
}

func Encode(r Reader) (*models.RawData, error) {
	const methodName = "Encode"

	var protocol models.RawData

	cmd, err := parseCommand(r)
	if err != nil {
		return nil, errror.WrapError(err, methodName)
	}
	protocol.Command = cmd

	payload, err := parsePayload(r)
	if err != nil {
		return nil, errror.WrapError(err, methodName)
	}
	protocol.Payload = payload

	return &protocol, nil
}

func parseCommand(r Reader) (int, error) {
	const methodName = "parseCommand"

	batch := make([]byte, CommandBatchSize)

	_, err := reader.New(r, reader.DefTimerDuration, CommandBatchSize).Read(batch)
	if err != nil {
		return 0, errror.WrapError(err, methodName)
	}

	i, err := strconv.ParseInt(string(batch[:1]), 10, 16)
	if err != nil {
		return 0, errror.WrapError(errors.Join(ErrWrongCommandFormat, err), methodName)
	}

	return int(i), nil
}

func parsePayload(r Reader) ([]byte, error) {
	const methodName = "parsePayload"

	var payload []byte

	for {
		batch := make([]byte, PayloadBatchSize)

		readByte, err := reader.New(r, reader.DefTimerDuration, PayloadBatchSize).Read(batch)
		if err != nil {
			return nil, errror.WrapError(err, methodName)
		}

		for _, line := range strings.Split(string(batch[:readByte]), "\n") {
			if line == EndLine {
				return payload, nil
			}
			payload = append(payload, []byte(line)...)
		}

	}
}
