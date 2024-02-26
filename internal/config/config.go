package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Vladimir77715/my-tcp-pow/internal/errror"
)

type ServerConfig struct {
	Address          string
	Port             int64
	MinSolutionRange int64
	MaxSolutionRange int64
}

type ClientConfig struct {
	Address string
}

const (
	defHost             = "0.0.0.0"
	defPort             = "8080"
	defSolutionRangeMin = "1"
	defSolutionRangeMax = "5"
	defAddr             = "0.0.0.0:8080"
)

func ParseServerConfig() (*ServerConfig, error) {
	const methodName = "ParseServerConfig"

	addr := GetenvOrDefVal("SERVER_ADDRESS", defHost)

	port, err := strconv.ParseInt(GetenvOrDefVal("SERVER_PORT", defPort), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error from method %s: %v", methodName, err)
	}

	minSolRange, err := strconv.ParseInt(GetenvOrDefVal("SERVER_MIN_SOLUTION_RANGE", defSolutionRangeMin), 10, 64)
	if err != nil {
		return nil, errror.WrapError(err, methodName)
	}

	maxSolRange, err := strconv.ParseInt(GetenvOrDefVal("SERVER_MAX_SOLUTION_RANGE", defSolutionRangeMax), 10, 64)
	if err != nil {
		return nil, errror.WrapError(err, methodName)
	}

	return &ServerConfig{Address: addr, Port: port, MinSolutionRange: minSolRange, MaxSolutionRange: maxSolRange}, nil
}

func ParseClientConfig() *ClientConfig {
	addr := GetenvOrDefVal("SERVER_ADDRESS", defAddr)

	return &ClientConfig{Address: addr}
}

func GetenvOrDefVal(key, defVal string) string {
	val := os.Getenv(key)
	if val == "" {
		val = defVal
	}

	return val
}
