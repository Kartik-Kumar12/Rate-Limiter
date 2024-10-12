package structs

type IpAddress string

type IpAddressConfig struct {
	IpAddresses []IpAddress `json:"ipAddresses"`
}
