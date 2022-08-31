<script>
    import Modal from "../index.svelte"
    import {profile} from "../../../store";
    import {variables} from "../../../store"

    let visible = true;

    let code = `
#include "ESP8266WiFi.h"
#include "Adafruit_MQTT.h"
#include "Adafruit_MQTT_Client.h"
#include "DHT.h"
#define DHTPIN 4 // пин к которому подключен датчик

DHT dht(DHTPIN, DHT11);

#define WLAN_SSID       "ivanola"
#define WLAN_PASS       "23031998"

/************************* MQTT Login, Password *********************************/

#define AIO_SERVER      "159.122.98.156"
#define AIO_SERVERPORT  8443                   // use 8883 for SSL
#define AIO_USERNAME    "mosquitto"
#define AIO_KEY         ""

WiFiClient client;
Adafruit_MQTT_Client mqtt(&client, AIO_SERVER, AIO_SERVERPORT, AIO_USERNAME, AIO_KEY);

Adafruit_MQTT_Publish tempP = Adafruit_MQTT_Publish(&mqtt, AIO_USERNAME "/home/temp");

Adafruit_MQTT_Publish humidP = Adafruit_MQTT_Publish(&mqtt, AIO_USERNAME "/home/humid");

void MQTT_connect();

void setup() {
  Serial.begin(115200);
  delay(10);

  WiFi.begin(WLAN_SSID, WLAN_PASS);
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  Serial.println();

  Serial.println("WiFi connected");
  Serial.println("IP address: "); Serial.println(WiFi.localIP());

  // Setup MQTT subscription for onoff feed.
  dht.begin();
}

uint32_t x=0;

void loop() {
  delay(10000); // 2 секунды задержки
  float h = dht.readHumidity(); //Измеряем влажность
  float t = dht.readTemperature(); //Измеряем температуру
  // Ensure the connection to the MQTT server is alive (this will make the first
  // connection and automatically reconnect when disconnected).  See the MQTT_connect
  // function definition further below.
  MQTT_connect();

  // this is our 'wait for incoming subscription packets' busy subloop
  // try to spend your time here


  // Now we can publish stuff!

  if (! tempP.publish(t)) {
    Serial.println(F("Failed"));
  } else {
    Serial.println(F("OK!"));
  }

  if (! humidP.publish(h)) {
    Serial.println(F("Failed"));
  } else {
    Serial.println(F("OK!"));
  }

  // ping the server to keep the mqtt connection alive
  // NOT required if you are publishing once every KEEPALIVE seconds
  /*
  if(! mqtt.ping()) {
    mqtt.disconnect();
  }
  */
}

// Function to connect and reconnect as necessary to the MQTT server.
// Should be called in the loop function and it will take care if connecting.
void MQTT_connect() {
  int8_t ret;

  // Stop if already connected.
  if (mqtt.connected()) {
    return;
  }

  Serial.print("Connecting to MQTT... ");

  uint8_t retries = 3;
  while ((ret = mqtt.connect()) != 0) { // connect will return 0 for connected
       Serial.println(mqtt.connectErrorString(ret));
       Serial.println("Retrying MQTT connection in 5 seconds...");
       mqtt.disconnect();
       delay(5000);  // wait 5 seconds
       retries--;
       if (retries == 0) {
         // basically die and wait for WDT to reset me
         while (1);
       }
  }
  Serial.println("MQTT Connected!");
}`;
</script>

{#if visible}
    <Modal width="800">
        <h1>Гайд по подключению устройств</h1>
        <div style=" overflow: auto; height: 400px; text-align: left; margin: 10px 30px;">
        <div style="text-align: left; width: 700px; margin: auto;">
            Для работы передачи данных между устройствами и платформой используется протокол MQTT.
            Для подключения вам понадобится данные MQTT сервера: <br><br>
            <strong>MQTT Хост: </strong> 159.122.98.156<br>
            <strong>MQTT порт: </strong> 1883<br>
            <strong>MQTT имя пользователя</strong> mosquitto<br>
            <strong>MQTT пароль </strong> (оставьте поле пустым)<br>

            <strong>Чтобы подписаться на получение данных с датчика используйте следующий топик</strong>
            <pre><code>value/{$profile.mqttToken}/<b>mqttLabelOfDevice</b></code></pre>
            <div> - где mqttLabelOfDevice - Label устройства, доступный на странице детального просмотра. <br>
                Этот топик также используется для публикации значений
            </div>
        </div>

        <div style="margin: 10px auto; width: 700px; text-align: left">Пример кода для Esp8266 на ArduinoIDE</div>
        <div style="overflow: auto; height: 400px; text-align: left; width: 700px; margin: auto">
            <pre><code class="language-javascript">{@html code}</code></pre>
        </div>
        </div>
    </Modal>
{/if}