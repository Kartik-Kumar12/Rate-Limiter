package client

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/Kartik-Kumar12/Rate-Limiter/rate_limiter_system/client/structs"
	"github.com/Kartik-Kumar12/Rate-Limiter/rate_limiter_system/common/utils"
)

const (
	IPConfigFileName = "config.json"
)

func runSequenctially() error {

	jsonBytes, err := utils.ReadFileContent(IPConfigFileName)
	if err != nil {
		log.Error().Err(err).Msg("Error reading IP configs")
		return err
	}

	var ipConfigs structs.IpAddressConfig
	if err = json.Unmarshal(jsonBytes, &ipConfigs); err != nil {
		log.Error().Err(err).Msg("Error Unmarshalling IP configs")
		return err
	}

	for _, ip := range ipConfigs.IpAddresses {
		url := "http://localhost:8080/ping?ip=" + string(ip)

		resp, err := http.Get(url)
		if err != nil {
			log.Error().Err(err).Msgf("Error making request for ip %v : %v", string(ip), err)
			return err
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error().Err(err).Msgf("Error making request for ip %v : %v", string(ip), err)
			return err
		}
		log.Info().Msgf("Response for IpAddress %v : \nStatus Code : %v\nMessage : %v", string(ip), resp.StatusCode, body)
	}
	return nil
}

func main() {

	if err := runSequenctially(); err != nil {
		log.Error().Err(err).Msg("")
		os.Exit(1)
	}

}
