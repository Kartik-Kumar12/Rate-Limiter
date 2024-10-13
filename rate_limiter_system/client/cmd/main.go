package main

import (
	"os"

	"github.com/Kartik-Kumar12/Rate-Limiter/rate_limiter_system/client/services"
	"github.com/Kartik-Kumar12/Rate-Limiter/rate_limiter_system/common/cli"
	"github.com/rs/zerolog/log"
)

func main() {

	cli.SetLogger()

	// MAKE SEQUENTIAL REQUESTS
	// for i := 0; i < 10; i++ {
	// 	fmt.Printf("%vth iteration\n", i+1)
	if err := services.ExecuteSequentially(); err != nil {
		log.Error().Err(err).Msg("Failed to execute requests sequentially")
		os.Exit(1)
	}
	// time.Sleep(200 * time.Millisecond)
	// }

	// MAKE PARALLEL REQUESTS
	// for i := 0; i < 10; i++ {
	// 	if err := services.ExecuteConcurrently(); err != nil {
	// 		log.Error().Err(err).Msg("Failed to execute requests sequentially")
	// 		os.Exit(1)
	// 	}
	// 	time.Sleep(200 * time.Millisecond)
	// }
}
