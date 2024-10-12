package main

import (
	"os"

	"github.com/Kartik-Kumar12/Rate-Limiter/rate_limiter_system/client/services"
	"github.com/rs/zerolog/log"
)

func main() {

	if err := services.ExecuteSequentially(); err != nil {
		log.Error().Err(err).Msg("Failed to execute requests sequentially")
		os.Exit(1)
	}
	if err := services.ExecuteConcurrently(); err != nil {
		log.Error().Err(err).Msg("Failed to execute requests sequentially")
		os.Exit(1)
	}

}
