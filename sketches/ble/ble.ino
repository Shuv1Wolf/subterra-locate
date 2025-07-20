#include <sys/time.h>
#include "BLEDevice.h"
#include "BLEUtils.h"
#include "BLEAdvertising.h"
#include "BLEBeacon.h"
#include "esp_sleep.h"

#define GPIO_DEEP_SLEEP_DURATION 3  // sleep duration in seconds
#define BEACON_UUID "0fe366e2-80a7-40fd-adee-89387f69a43d"

RTC_DATA_ATTR static time_t lastBootTime = 0;         // Time of last boot (stored in RTC memory)
RTC_DATA_ATTR static uint32_t bootCount = 0;          // Boot counter (stored in RTC memory)

BLEAdvertising* pAdvertising;
struct timeval now;

// ────── Setup BLE Beacon ──────
void setBeacon() {
  BLEBeacon oBeacon;
  oBeacon.setManufacturerId(0x4C00); // Apple iBeacon ID
  oBeacon.setProximityUUID(BLEUUID(BEACON_UUID));
  oBeacon.setMajor((bootCount & 0xFFFF0000) >> 16);
  oBeacon.setMinor(bootCount & 0xFFFF);

  BLEAdvertisementData advData;
  BLEAdvertisementData scanRespData;
  advData.setFlags(0x04); // BR/EDR not supported

  String serviceData = "";
  serviceData += (char)0x1A; // Length
  serviceData += (char)0xFF; // Type
  serviceData += oBeacon.getData().c_str();  

  advData.addData(std::string(serviceData.c_str()));  
  pAdvertising = BLEDevice::getAdvertising();
  pAdvertising->setAdvertisementData(advData);
  pAdvertising->setScanResponseData(scanRespData);
}

// ────── Print boot/log info ──────
void printBootInfo() {
  gettimeofday(&now, NULL);
  Serial.printf("Boot #%d\n", bootCount++);
  Serial.printf("Time since last boot: %ld s\n", now.tv_sec - lastBootTime);
  lastBootTime = now.tv_sec;
}

// ────── Enter deep sleep ──────
void goToDeepSleep() {
  Serial.println("Entering deep sleep...");
  esp_sleep_enable_timer_wakeup(GPIO_DEEP_SLEEP_DURATION * 1000000ULL);
  esp_deep_sleep_start();
}

void setup() {
  Serial.begin(115200);
  delay(100);

  printBootInfo();

  // Init BLE device
  BLEDevice::init("ESP32 iBeacon");

  // Set up Beacon advertisement
  setBeacon();

  if (pAdvertising) {
    pAdvertising->start();
    Serial.println("Advertising started...");
    delay(100);
    pAdvertising->stop();
  }

  goToDeepSleep(); // Will never return from here
}

void loop() {
  // Empty. All logic is in setup(), and device sleeps after advertisement.
}
