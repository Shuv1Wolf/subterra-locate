package nats

const (
	NATS_EVENTS_BLE_RSSI_TOPIC  = "events.ble.raw.rssi"
	NATS_EVENTS_DEVICE_POSITION = "events.device.position"
)

const (
	NATS_BEACONS_EVENTS_TOPIC        = "beacons.events"
	NATS_BEACONS_EVENTS_CREATED_TYPE = "beacons.events.created"
	NATS_BEACONS_EVENTS_DELETED_TYPE = "beacons.events.deleted"
	NATS_BEACONS_EVENTS_CHANGED_TYPE = "beacons.events.changed"
)

const (
	NATS_DEVICE_EVENTS_TOPIC        = "device.events"
	NATS_DEVICE_EVENTS_CREATED_TYPE = "device.events.created"
	NATS_DEVICE_EVENTS_DELETED_TYPE = "device.events.deleted"
	NATS_DEVICE_EVENTS_CHANGED_TYPE = "device.events.changed"
)
