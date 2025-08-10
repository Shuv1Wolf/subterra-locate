#include <Arduino.h>
#include <sstream>
#include <string>

#include <WiFi.h>
#define MQTT_MAX_PACKET_SIZE 1844
#include <PubSubClient.h>

// BLE
#include <BLEDevice.h>
#include <BLEUtils.h>
#include <BLEScan.h>
#include <BLEAdvertisedDevice.h>

#include "credentials.h"

// ---------- НАСТРОЙКИ ----------
static const char* DEVICE_ID = "device$100";

// ---------- ПРОТОКОЛ МАЯКА ----------
static const uint16_t SERVICE_UUID_16 = 0xFFF0;
// Формат: [ver:1][id_len:1][beacon_id:id_len][txp:1]
struct ParsedSD {
  uint8_t  ver;
  String   beacon_id; 
  int8_t   txp;
};

static const int SCAN_WINDOW_SEC = 4; 

// ---------- СЕТЬ / MQTT ----------
WiFiClient espClient;
PubSubClient client(espClient);

// ---------- БУФЕР МАЯКОВ ----------
struct BeaconItem {
  String beacon_id;
  int    rssi;
  int8_t txp;
  bool   valid;
};
static BeaconItem beacons[32];
static uint8_t beaconCount = 0;

// ---------- УТИЛИТЫ ----------
static void resetBuffer() {
  for (auto &b : beacons) b.valid = false;
  beaconCount = 0;
}

static BeaconItem* upsertBeacon(const String& beacon_id) {
  for (auto &b : beacons) {
    if (b.valid && b.beacon_id == beacon_id) return &b;
  }
  for (auto &b : beacons) {
    if (!b.valid) {
      b.valid = true;
      b.beacon_id = beacon_id;
      b.rssi = -127;
      b.txp  = -127;
      beaconCount++;
      return &b;
    }
  }
  return nullptr;
}

static bool parseServiceData(const std::string &serviceData, ParsedSD &out) {
  if (serviceData.size() < 3) return false;
  const uint8_t* p = (const uint8_t*)serviceData.data();

  uint8_t ver    = p[0];
  uint8_t id_len = p[1];

  if (serviceData.size() < (size_t)(2 + id_len + 1)) return false;

  const char* id_ptr = (const char*)(p + 2);
  out.ver       = ver;
  out.beacon_id = String(id_ptr, id_len);
  out.txp       = (int8_t)p[2 + id_len];
  return true;
}

static String jsonEscape(const String& s) {
  String out; out.reserve(s.length() + 4);
  for (size_t i = 0; i < s.length(); ++i) {
    char c = s[i];
    if (c == '\"' || c == '\\') { out += '\\'; out += c; }
    else out += c;
  }
  return out;
}

// ---------- BLE CALLBACKS ----------
class MyAdvertisedDeviceCallbacks : public BLEAdvertisedDeviceCallbacks {
public:
  void onResult(BLEAdvertisedDevice advertisedDevice) override {
    if (!advertisedDevice.haveServiceData()) return;

    for (size_t i = 0; i < advertisedDevice.getServiceDataCount(); ++i) {
      BLEUUID su = advertisedDevice.getServiceDataUUID(i);
      if (!su.equals(BLEUUID((uint16_t)SERVICE_UUID_16))) continue;

      std::string sd = advertisedDevice.getServiceData(i);
      ParsedSD p{};
      if (!parseServiceData(sd, p)) continue;

      BeaconItem* it = upsertBeacon(p.beacon_id);
      if (!it) return;

      int rssi = advertisedDevice.haveRSSI() ? advertisedDevice.getRSSI() : -127;
      if (rssi > it->rssi) it->rssi = rssi;
      it->txp = p.txp;

      Serial.printf("[beacon] id=%s rssi=%d txp=%d ver=%u\n",
                    it->beacon_id.c_str(), it->rssi, it->txp, p.ver);
    }
  }
};

// ---------- WIFI / MQTT ----------
void connectWiFi() {
  if (WiFi.status() == WL_CONNECTED) return;
  Serial.println("Connecting to WiFi...");
  WiFi.mode(WIFI_STA);
  WiFi.begin(ssid, password);
  uint32_t t0 = millis();
  while (WiFi.status() != WL_CONNECTED && millis() - t0 < 15000) {
    delay(300);
    Serial.print(".");
  }
  Serial.println();
  if (WiFi.status() == WL_CONNECTED) {
    Serial.print("WiFi OK, IP: ");
    Serial.println(WiFi.localIP());
  } else {
    Serial.println("WiFi not connected");
  }
}

void connectMQTT() {
  if (client.connected()) return;
  client.setServer(mqttServer, mqttPort);
  Serial.print("Connecting MQTT ... ");
  if (client.connect("ESP32Client", mqttUser, mqttPassword)) {
    Serial.println("OK");
  } else {
    Serial.print("failed, rc=");
    Serial.println(client.state());
  }
}

// ---------- СКАН / ПУБЛИКАЦИЯ ----------
void scanOnce() {
  resetBuffer();
  BLEScan* pBLEScan = BLEDevice::getScan();
  static MyAdvertisedDeviceCallbacks cb;
  pBLEScan->setAdvertisedDeviceCallbacks(&cb);
  pBLEScan->setActiveScan(true);
  pBLEScan->setInterval(100);
  pBLEScan->setWindow(80);

  Serial.printf("Scanning %ds...\n", SCAN_WINDOW_SEC);
  pBLEScan->start(SCAN_WINDOW_SEC, false);
  pBLEScan->stop();

  Serial.printf("Beacons in window: %u\n", beaconCount);
}

void publishNow() {
  if (beaconCount < 3) {
    Serial.printf("Skip publish: only %u beacon(s) detected\n", beaconCount);
    return;
  }

  const String dev_mac = WiFi.macAddress();

  String payload = "{\"dev_mac\":\"";
  payload += dev_mac;
  payload += "\",\"dev_id\":\"";
  payload += DEVICE_ID;
  payload += "\",\"count\":";
  payload += String(beaconCount);
  payload += ",\"e\":[";
  uint8_t added = 0;
  for (auto &b : beacons) {
    if (!b.valid) continue;
    if (added) payload += ",";
    payload += "{\"id\":\"";
    payload += jsonEscape(b.beacon_id);
    payload += "\",\"r\":";
    payload += String(b.rssi);
    payload += ",\"txp\":";
    payload += String((int)b.txp);
    payload += "}";
    added++;
  }
  payload += "]}";

  Serial.print("Payload length: ");
  Serial.println(payload.length());
  Serial.println(payload);

  static uint8_t msgbuf[MQTT_MAX_PACKET_SIZE];
  payload.getBytes(msgbuf, payload.length() + 1);

  if (!client.connected()) connectMQTT();
  if (client.connected()) {
    bool ok = client.publish("/ble/rssi", msgbuf, payload.length(), false);
    Serial.print("MQTT publish: ");
    Serial.println(ok ? "OK" : "FAIL");
  }
}

// ---------- SETUP / LOOP ----------
void setup() {
  Serial.begin(115200);
  delay(100);
  BLEDevice::init("");

  connectWiFi();
  connectMQTT();
}

void loop() {
  scanOnce();
  connectWiFi();
  connectMQTT();
  client.loop();
  publishNow();
  delay(500);
}
