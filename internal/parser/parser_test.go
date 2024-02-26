package parser_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/Vladimir77715/my-tcp-pow/internal/models"
	"github.com/Vladimir77715/my-tcp-pow/internal/parser"
	"github.com/Vladimir77715/my-tcp-pow/internal/parser/mocks"
	"github.com/Vladimir77715/my-tcp-pow/internal/reader"
	"github.com/Vladimir77715/my-tcp-pow/internal/tcp/server"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	t.Run("Success parse", func(t *testing.T) {
		rd := mocks.NewReader(t)

		rd.EXPECT().Read(mock.MatchedBy(func(b []byte) bool {
			write([]byte(strconv.Itoa(server.RequestChallenge)), b)
			return true
		})).Return(0, nil).Once()

		rd.EXPECT().Read(mock.MatchedBy(func(b []byte) bool {
			write([]byte("{}"), b)
			return true
		})).Return(2, nil).Once()
		rd.EXPECT().Read(mock.MatchedBy(func(b []byte) bool {
			write([]byte(parser.EndLine), b)
			return true
		})).Return(3, nil).Once()

		data, err := parser.Encode(rd)
		require.NoError(t, err)
		require.NotNil(t, data)
		require.Equal(t, *data, models.RawData{Command: 1, Payload: []byte("{}")})
	})

	t.Run("Success parse without payload", func(t *testing.T) {
		rd := mocks.NewReader(t)

		rd.EXPECT().Read(mock.MatchedBy(func(b []byte) bool {
			write([]byte(strconv.Itoa(server.RequestChallenge)), b)
			return true
		})).Return(0, nil).Once()

		rd.EXPECT().Read(mock.MatchedBy(func(b []byte) bool {
			write([]byte(parser.EndLine), b)
			return true
		})).Return(3, nil).Once()

		data, err := parser.Encode(rd)
		require.NoError(t, err)
		require.NotNil(t, data)
		require.Equal(t, *data, models.RawData{Command: 1, Payload: nil})
	})

	t.Run("Client input timeout", func(t *testing.T) {
		rd := mocks.NewReader(t)

		rd.EXPECT().Read(mock.MatchedBy(func(b []byte) bool {
			time.Sleep(11 * time.Second)
			return true
		})).Return(0, nil).Once()

		_, err := parser.Encode(rd)
		require.ErrorIs(t, err, reader.ErrTimeout)
	})

	t.Run("Wrong command format", func(t *testing.T) {
		rd := mocks.NewReader(t)

		rd.EXPECT().Read(mock.MatchedBy(func(b []byte) bool {
			write([]byte("nn"), b)
			return true
		})).Return(0, nil).Once()

		_, err := parser.Encode(rd)
		require.ErrorIs(t, err, parser.ErrWrongCommandFormat)
	})
}

func write(data []byte, b []byte) {
	for i := 0; i < len(data); i++ {
		b[i] = data[i]
	}
}
