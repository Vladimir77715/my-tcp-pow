package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/Vladimir77715/my-tcp-pow/internal/config"
	"github.com/Vladimir77715/my-tcp-pow/internal/errror"
	"github.com/Vladimir77715/my-tcp-pow/internal/models"
	"github.com/Vladimir77715/my-tcp-pow/internal/parser"
	"github.com/Vladimir77715/my-tcp-pow/internal/tcp/server"
)

// Функция для вычисления хеша с условием
func computeProofOfWork(originalData string, zeroCount int) int {
	nonce := 0

	bld := strings.Builder{}

	for i := 0; i < zeroCount; i++ {
		bld.WriteRune('0')
	}

	zeroCountStr := bld.String()

	for {
		// Конвертируем nonce из int в строку для конкатенации
		nonceStr := strconv.Itoa(nonce)
		// Вычисляем хеш
		hash := sha256.Sum256([]byte(originalData + nonceStr))
		hashStr := hex.EncodeToString(hash[:])
		// Проверяем условие хеша
		if hashStr[:zeroCount] == zeroCountStr { // Для простоты предполагаем, что ищем 4 нуля
			return nonce
		}
		nonce++
	}
}

func main() {
	clientCfg := config.ParseClientConfig()

	cl := client{clientCfg.Address}

	hashData, err := cl.requestForTask()
	if err != nil {
		log.Fatalln(err)
	}

	hashData.Nonce = computeProofOfWork(hashData.Data, hashData.ZeroCount)

	at, err := cl.sendSolution(hashData)
	if err != nil {
		log.Fatalln(err)
	}

	wisd, err := cl.getWisdom(at)
	if err != nil {
		log.Fatalln(err)
	}

	log.Print(wisd)
}

type client struct {
	addr string
}

func (c *client) requestForTask() (*models.HashData, error) {
	const methodName = "requestForTask"

	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return nil, errror.WrapError(err, methodName)
	}

	_, err = conn.Write(parser.Decode(&models.RawData{Command: 1}))
	if err != nil {
		return nil, errror.WrapError(err, methodName)
	}

	data, err := parser.Encode(conn)
	if err != nil {
		return nil, errror.WrapError(err, methodName)
	}

	hashData := models.HashData{}

	err = json.Unmarshal(data.Payload, &hashData)
	if err != nil {
		return nil, errror.WrapError(err, methodName)
	}

	return &hashData, nil
}

func (c *client) sendSolution(data *models.HashData) (*models.AccessToken, error) {
	const methodName = "sendSolution"

	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return nil, errror.WrapError(err, methodName)
	}

	b, err := json.Marshal(data)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = conn.Write(parser.Decode(&models.RawData{Command: server.SendSolution, Payload: b}))
	if err != nil {
		return nil, errror.WrapError(err, methodName)
	}

	respData, err := parser.Encode(conn)
	if err != nil {
		return nil, errror.WrapError(err, methodName)
	}

	accessToken := models.AccessToken{}

	err = json.Unmarshal(respData.Payload, &accessToken)
	if err != nil {
		return nil, errror.WrapError(err, methodName)
	}

	println(accessToken.Token)

	return &accessToken, nil
}

func (c *client) getWisdom(at *models.AccessToken) (string, error) {
	const methodName = "getWisdom"

	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return "", errror.WrapError(err, methodName)
	}

	b, err := json.Marshal(at)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = conn.Write(parser.Decode(&models.RawData{Command: server.WisdomRequest, Payload: b}))
	if err != nil {
		return "", errror.WrapError(err, methodName)
	}

	respData, err := parser.Encode(conn)
	if err != nil {
		return "", errror.WrapError(err, methodName)
	}

	return string(respData.Payload), nil
}
