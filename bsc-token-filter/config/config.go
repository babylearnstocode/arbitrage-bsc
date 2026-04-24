package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ArchiveUrl       string
	RPCUrl           string
	V2FactoryAddress string
	IpcPath          string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}

	return &Config{
		ArchiveUrl:       os.Getenv("INFURA_URL"),
		RPCUrl:           os.Getenv("LOCAL_RPC_URL"),
		V2FactoryAddress: os.Getenv("V2_FACTORY_ADDRESS"),
		IpcPath:          os.Getenv("LOCAL_IPC_PATH"),
	}
}
