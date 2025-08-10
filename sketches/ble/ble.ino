#include <Arduino.h>
#include <BLEDevice.h>
#include <BLEAdvertising.h>
#include <BLEBeacon.h>

#define BEACON_NAME        "beacon$102"
#define BEACON_UUID        "2D7A9F0C-E0E8-4CC9-A71B-A21DB2D034A1"
#define BEACON_MAJOR       1
#define BEACON_MINOR       102
#define BEACON_TXPOWER     0xC5
#define APPLE_MANUFACTURER_ID 0x4C00
#define SERVICE_UUID_16    0xFFF0

// Строковый ID
const char* BEACON_ID_STR = "beacon$102";

#define ADV_INTERVAL_MS    250
#define REFRESH_MS         3000

BLEAdvertising* adv = nullptr;

static std::string buildServiceData() {
  uint8_t ver = 1;
  uint8_t id_len = strlen(BEACON_ID_STR);
  uint8_t txp = (uint8_t)BEACON_TXPOWER;

  std::string sd;
  sd.reserve(1 + 1 + id_len + 1);

  sd.push_back(ver);
  sd.push_back(id_len);
  sd.append(BEACON_ID_STR, id_len);
  sd.push_back(txp);

  return sd;
}

static void setupAdvertising() {
  BLEDevice::init(BEACON_NAME);
  BLEDevice::setPower(ESP_PWR_LVL_N0);

  BLEServer* server = BLEDevice::createServer(); (void)server;

  BLEBeacon ib;
  ib.setManufacturerId(APPLE_MANUFACTURER_ID);
  ib.setProximityUUID(BLEUUID(BEACON_UUID));
  ib.setMajor(BEACON_MAJOR);
  ib.setMinor(BEACON_MINOR);
  ib.setSignalPower((int8_t)BEACON_TXPOWER);

  BLEAdvertisementData advData;
  BLEAdvertisementData scanData;

  advData.setFlags(ESP_BLE_ADV_FLAG_GEN_DISC | ESP_BLE_ADV_FLAG_BREDR_NOT_SPT);
  advData.setManufacturerData(ib.getData());

  scanData.setName(BEACON_NAME);
  scanData.setServiceData(BLEUUID((uint16_t)SERVICE_UUID_16), buildServiceData());

  adv = BLEDevice::getAdvertising();
  adv->setAdvertisementData(advData);
  adv->setScanResponseData(scanData);
  adv->setAdvertisementType(ADV_TYPE_SCAN_IND);

  uint16_t intervalUnits = (uint16_t)(ADV_INTERVAL_MS / 0.625);
  adv->setMinInterval(intervalUnits);
  adv->setMaxInterval(intervalUnits);

  adv->start();
  Serial.println("Beacon with string ID started");
}

void setup() {
  Serial.begin(115200);
  delay(100);
  setupAdvertising();
}

void loop() {
  static uint32_t last = 0;
  if (millis() - last > REFRESH_MS) {
    last = millis();

    BLEAdvertisementData advData;
    BLEAdvertisementData scanData;

    BLEBeacon ib;
    ib.setManufacturerId(APPLE_MANUFACTURER_ID);
    ib.setProximityUUID(BLEUUID(BEACON_UUID));
    ib.setMajor(BEACON_MAJOR);
    ib.setMinor(BEACON_MINOR);
    ib.setSignalPower((int8_t)BEACON_TXPOWER);

    advData.setFlags(ESP_BLE_ADV_FLAG_GEN_DISC | ESP_BLE_ADV_FLAG_BREDR_NOT_SPT);
    advData.setManufacturerData(ib.getData());

    scanData.setName(BEACON_NAME);
    scanData.setServiceData(BLEUUID((uint16_t)SERVICE_UUID_16), buildServiceData());

    adv->stop();
    adv->setAdvertisementData(advData);
    adv->setScanResponseData(scanData);
    adv->start();
  }
  delay(10);
}
