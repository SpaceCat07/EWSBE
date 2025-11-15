#include <Wire.h>
#include <Adafruit_GFX.h>
#include <Adafruit_SH110X.h>
#include "Adafruit_SHT4x.h"
#include <Adafruit_BMP3XX.h>
#include <Adafruit_INA219.h>
#include <SPI.h>
#include <SD.h>

#define CS_PIN 5  // pin CS SD Card
File dataFile;

// MH-Z19 UART
#define MHZ_RX 17
#define MHZ_TX 16
HardwareSerial mhz19_uart(1);
uint8_t getppm_cmd[9] = {0xFF, 0x01, 0x86, 0x00, 0x00, 0x00, 0x00, 0x00, 0x79};


// Tambahan untuk curah hujan
#define RAIN_PIN 13       // Ganti sesuai pin digital yang digunakan
volatile unsigned int rainTicks = 0;
float rainAmount = 0.0;

void IRAM_ATTR countRain() {
  rainTicks++;
}

// OLED
#define SCREEN_WIDTH 128
#define SCREEN_HEIGHT 64
#define OLED_RESET -1
#define I2C_ADDRESS 0x3C
Adafruit_SH1106G display(SCREEN_WIDTH, SCREEN_HEIGHT, &Wire, OLED_RESET);

// Sensor
Adafruit_SHT4x sht4 = Adafruit_SHT4x();
Adafruit_BMP3XX bmp;
Adafruit_INA219 ina219;

// Ultrasonik
#define TRIG_PIN 0
#define ECHO_PIN 2


// Anemometer
volatile byte rpmcount = 0;
volatile unsigned long last_micros;
unsigned long timeold = 0;
float kecepatan_meter_per_detik = 0.0;
int GPIO_pulse = 14;
volatile boolean flag = false;

// Analog Pins
#define VOLTAGE_SENSOR_PIN 34  // GPIO36, A0
#define WIND_VANE_PIN      35  // GPIO39, A1

// ISR Anemometer
void IRAM_ATTR rpm_anemometer() {
  flag = true;
}

void setup() {

  Serial.begin(115200);
  //mhzSerial.begin(9600, SERIAL_8N1, 16, 17);
  Wire.begin();
  mhz19_uart.begin(9600, SERIAL_8N1, MHZ_RX, MHZ_TX);
  delay(1000); // Warming up (kalau perlu, bisa delay(60000))
  // Interrupt untuk tipping bucket rain gauge
  pinMode(RAIN_PIN, INPUT_PULLUP);
  attachInterrupt(digitalPinToInterrupt(RAIN_PIN), countRain, FALLING);


  // Inisialisasi SD Card
  if (!SD.begin(CS_PIN)) {
    Serial.println("Gagal inisialisasi kartu SD");
    display.setCursor(0, 0);
    display.println("SD Card gagal!");
    display.display();
    while (true);
  }
  Serial.println("Kartu SD siap dipakai");

  // Buka logg.csv, cek apakah kosong → tulis header
// Buka catat.csv, cek apakah kosong → tulis header
  dataFile = SD.open("/catat.csv", FILE_WRITE);

  if (dataFile) {
    if (dataFile.size() == 0) {
      dataFile.println("waktu_ms,suhu,lembap,tekanan,ketinggian,co2,jarak,angin,arahAngin,teganganINA,arusINA,teganganSensor,curahHujan");
      Serial.println("Header ditulis ke catat.csv");
    }
    dataFile.close();
  } else {
    Serial.println("Gagal membuka catat.csv");
  }

  pinMode(TRIG_PIN, OUTPUT);
  pinMode(ECHO_PIN, INPUT);
  pinMode(GPIO_pulse, INPUT_PULLUP);
  attachInterrupt(digitalPinToInterrupt(GPIO_pulse), rpm_anemometer, RISING);

  display.begin(I2C_ADDRESS, true);
  display.clearDisplay();
  display.setTextSize(1);
  display.setTextColor(SH110X_WHITE);

  // SHT40
  if (!sht4.begin()) {
    display.setCursor(0, 0);
    display.println("SHT4x gagal!");
    display.display();
    while (1) delay(10);
  }
  sht4.setPrecision(SHT4X_HIGH_PRECISION);
  sht4.setHeater(SHT4X_NO_HEATER);

  //pinMode(MHZ_PWM_PIN, INPUT);

  // BMP390
  if (!bmp.begin_I2C()) {
    display.setCursor(0, 0);
    display.println("BMP390 gagal!");
    display.display();
    while (1) delay(10);
  }
  bmp.setTemperatureOversampling(BMP3_OVERSAMPLING_8X);
  bmp.setPressureOversampling(BMP3_OVERSAMPLING_4X);
  bmp.setIIRFilterCoeff(BMP3_IIR_FILTER_COEFF_3);
  bmp.setOutputDataRate(BMP3_ODR_50_HZ);

  // INA219
  if (!ina219.begin()) {
    Serial.println("INA219 gagal!");
    display.setCursor(0, 0);
    display.println("INA219 gagal!");
    display.display();
    while (1) delay(10);
  }
}


float readJsnDistance() {
  digitalWrite(TRIG_PIN, LOW);
  delayMicroseconds(2);
  digitalWrite(TRIG_PIN, HIGH);
  delayMicroseconds(10);
  digitalWrite(TRIG_PIN, LOW);
  long duration = pulseIn(ECHO_PIN, HIGH, 30000);
  if (duration == 0) return -1;
  return duration * 0.0343 / 2.0;
}
void loggToSD(unsigned long waktu, float suhu, float lembap, float tekanan, float ketinggian,
             int co2, float jarak, float angin, float arahAngin, float busVoltage, float current_mA,
             float voltSensor, float rain) {
  dataFile = SD.open("/catat.csv", FILE_APPEND);
  if (dataFile) {
    dataFile.printf("%lu,%.2f,%.2f,%.2f,%.2f,%d,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f\n",
      waktu, suhu, lembap, tekanan, ketinggian, co2, jarak, angin, arahAngin,
      busVoltage, current_mA, voltSensor, rain);
    dataFile.close();
    Serial.println("Data tersimpan ke log.csv");
  } else {
    Serial.println("Gagal menulis ke log.csv");
  }
}

int readCO2_UART() {
  while (mhz19_uart.available()) mhz19_uart.read(); // flush

  mhz19_uart.write(getppm_cmd, 9);
  mhz19_uart.flush();

  uint32_t start = millis();
  while (mhz19_uart.available() < 9) {
    if (millis() - start > 1000) return -1; // timeout
  }

  uint8_t response[9];
  mhz19_uart.readBytes(response, 9);

  if (response[0] == 0xFF && response[1] == 0x86 && checksum(response) == response[8]) {
    return response[2] * 256 + response[3];
  } else {
    return -1;
  }
}

uint8_t checksum(uint8_t *packet) {
  uint8_t sum = 0;
  for (int i = 1; i < 8; i++) sum += packet[i];
  return 0xFF - sum + 1;
}


float readRainAmount() {
  static unsigned int lastTicks = 0;
  unsigned int currentTicks;

  noInterrupts();
  currentTicks = rainTicks;
  interrupts();

  float deltaTicks = currentTicks - lastTicks;
  lastTicks = currentTicks;

  // Misalnya: 1 tipping = 0.2 mm
  float rain_mm = deltaTicks * 0.2;
  rainAmount += rain_mm;

  return rain_mm;  // return curah hujan selama loop terakhir
}


float hitungKecepatanAngin() {
  static const float waktuUkur = 10.0;
  static unsigned long lastCheck = millis();

  if ((millis() - lastCheck) >= waktuUkur * 1000) {
    detachInterrupt(digitalPinToInterrupt(GPIO_pulse));
    float rps = float(rpmcount) / waktuUkur;
    float v = (-0.0181 * rps * rps) + (1.3859 * rps) + 1.4055;
    if (v <= 1.5) v = 0.0;
    kecepatan_meter_per_detik = v;
    rpmcount = 0;
    lastCheck = millis();
    attachInterrupt(digitalPinToInterrupt(GPIO_pulse), rpm_anemometer, RISING);
  }

  if (flag && (micros() - last_micros >= 5000)) {
    rpmcount++;
    last_micros = micros();
    flag = false;
  }

  return kecepatan_meter_per_detik;
}

float readVoltageSensor() {
  int adcVal = analogRead(VOLTAGE_SENSOR_PIN);
  float voltage = (adcVal / 4095.0) * 3.3 * 5.0; // misalnya pakai voltage divider (1:5)
  return voltage;
}

float readWindVane() {
  int adcVal = analogRead(WIND_VANE_PIN);
  float voltage = (adcVal / 4095.0) * 3.3;
  float angle = (voltage / 3.3) * 360.0;
  // Clamp nilai agar tidak keluar batas
  if (angle < 0) angle = 0;
  if (angle > 360) angle = 360;
  return angle;
}

void loop() {
  // Sensor readings
  float rain = readRainAmount();
  sensors_event_t humEvent, tempEvent;
  sht4.getEvent(&humEvent, &tempEvent);
  float suhu = tempEvent.temperature;
  float lembap = humEvent.relative_humidity;

  if (!bmp.performReading()) return;
  float tekanan = bmp.pressure / 100.0;
  float ketinggian = bmp.readAltitude(1013.25);

  int co2 = readCO2_UART();


  float jarak = readJsnDistance();
  float angin = hitungKecepatanAngin();

  float busVoltage = ina219.getBusVoltage_V();
  float current_mA = ina219.getCurrent_mA();

  float voltSensor = readVoltageSensor();
  float arahAngin = readWindVane();

  // Serial print
  Serial.printf("T: %.1fC, H: %.1f%%, P: %.1fhPa, Alt: %.1fm, CO2: %dppm, Jarak: %.1fcm, Angin: %.1fm/s, INA219: %.2fV %.2fmA, V: %.2fV, WindDir: %.1f°, Hujan: %.1fmm\n",
    suhu, lembap, tekanan, ketinggian, co2, jarak, angin, busVoltage, current_mA, voltSensor, arahAngin, rain);


  // OLED display
  display.clearDisplay();
  display.setCursor(0, 0);
  display.printf("T:%.0fC H:%.0f%%\n", suhu, lembap);
  display.printf("P:%.0fhPa CO2:%s\n", tekanan, co2 == -1 ? "Err" : String(co2).c_str());
  display.printf("Jrk:%.0fcm A:%.1fm/s\n", jarak, angin);
  // display.printf("INA:%.1fV %.0fmA\n", busVoltage, current_mA);
  display.printf("INA:%.0fmA\n", current_mA);
  display.printf("V:%.1fV Dir:%.0fdeg\n", voltSensor, arahAngin);
  display.printf("Hujan:%.1fmm", rain);
  display.display();

  unsigned long waktu = millis();
  loggToSD(waktu, suhu, lembap, tekanan, ketinggian, co2, jarak, angin, arahAngin,
          busVoltage, current_mA, voltSensor, rain);


  delay(1000);
}