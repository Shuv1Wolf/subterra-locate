package events

type BLEBeaconRawEventV1 struct {
	DeviceMAC string           `json:"dev_mac"` // mac address
	DeviceId  string           `json:"dev_id"`  // device ID
	Count     int              `json:"count"`   // number of beacons
	Beacons   []BLEBeaconRawV1 `json:"e"`       // list of beacons
}

type BLEBeaconRawV1 struct {
	Id      string `json:"id"`  // beacon ID
	Rssi    int    `json:"r"`   // dBm
	Txpower int    `json:"txp"` // transmit power
}

type BLEBeaconHistoryEventV1 struct {
}
