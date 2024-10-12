package services

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/Kartik-Kumar12/Rate-Limiter/rate_limiter_system/client/structs"
	"github.com/Kartik-Kumar12/Rate-Limiter/rate_limiter_system/common/utils"
	"github.com/rs/zerolog/log"
)

const (
	IPConfigFileName = "config/ip_address.json"
)

func makePingRequest(ip string) {
	url := "http://localhost:8080/ping?ip=" + ip
	resp, err := http.Get(url)
	if err != nil {
		log.Error().Err(err).Msgf("Error making request for IP %v", ip)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msgf("Error reading response for IP %v", ip)
		return
	}
	log.Info().Msgf("Response for IP %v: Status Code: %v, Message: %v", ip, resp.StatusCode, string(body))
}

func ExecuteSequentially() error {
	configBytes, err := utils.ReadFileContent(IPConfigFileName)
	if err != nil {
		log.Error().Err(err).Msg("Error reading IP config file")
		return err
	}

	var ipConfig structs.IpAddressConfig
	if err := json.Unmarshal(configBytes, &ipConfig); err != nil {
		log.Error().Err(err).Msg("Error unmarshalling IP config")
		return err
	}

	for _, ip := range ipConfig.IpAddresses {
		makePingRequest(string(ip))
	}
	return nil
}

func ExecuteConcurrently() error {
	configBytes, err := utils.ReadFileContent(IPConfigFileName)
	if err != nil {
		log.Error().Err(err).Msg("Error reading IP config file")
		return err
	}

	var ipConfig structs.IpAddressConfig
	if err := json.Unmarshal(configBytes, &ipConfig); err != nil {
		log.Error().Err(err).Msg("Error unmarshalling IP config")
		return err
	}

	var wg sync.WaitGroup
	for _, ip := range ipConfig.IpAddresses {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			makePingRequest(ip)
		}(string(ip))
	}
	wg.Wait()
	return nil
}
