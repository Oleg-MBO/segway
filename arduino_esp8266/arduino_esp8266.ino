#include <dummy.h>

#include <Wire.h>

#include <ESP8266WiFi.h>
#include <WiFiUdp.h>

//#include <OSCMessage.h>


#define udpBufferLength 255
#define commandStartChar 62 //>
#define commandDataChar 58  //:
#define commandEndChar 60   //<

extern const int D1PWM;
extern const int D1DIR;
extern const int D2PWM;
extern const int D2DIR;

extern const char* ssid;
extern const char* password;

extern bool getAccData;
extern bool getGyroAngleData;
extern bool getAccAngleData;

WiFiUDP Udp;
unsigned int localUdpPort = 4210;  // local port to listen on
char incomingPacket[udpBufferLength];  // buffer for incoming packets


const long interval = 1000;           // interval at which to blink (milliseconds)
unsigned long previousMillis = 0;        // will store last time LED was updated

long counder = 0;

char buf[40];
bool gotmessage = false;
IPAddress outIp;
uint16_t outPort;



unsigned long timeLastCommangGot = 0;;

void setup()
{
  pinMode (D1PWM, OUTPUT);
  pinMode (D1DIR, OUTPUT);
  pinMode (D2PWM, OUTPUT);
  pinMode (D2DIR, OUTPUT);

  digitalWrite (D1DIR, HIGH);
  digitalWrite (D2DIR, HIGH);

  analogWrite(D1PWM, 0);
  analogWrite(D2PWM, 0);


  Serial.begin(115200);
  Serial.println("\n\nStarting.....");

  mpu_setup();
  Serial.println();

  Serial.printf("Connecting to %s ", ssid);
  WiFi.begin(ssid, password);

  int countWait = 0;
  while (WiFi.status() != WL_CONNECTED) {
    countWait++;
    delay(500);
    yield();
    Serial.print(".");
    if (countWait > 20) {
      ESP.restart();
    }
  }
  Serial.println(" connected");

  Udp.begin(localUdpPort);
  Serial.print("Now listening at IP ");
  Serial.print( WiFi.localIP().toString().c_str());

  Serial.print(":");
  Serial.println(localUdpPort);


}





void loop()
{
  if (WiFi.status() != WL_CONNECTED) {
    Serial.println("\n\nReset..");
    ESP.restart();
  }
  mpu_loop();

  //  __________________________________________________

  int packetSize = Udp.parsePacket();
  //  START HANDLE COMMAND
  if (packetSize)  {
    timeLastCommangGot = millis();
    // receive incoming UDP packets
    //    Serial.printf("Received % d bytes from % s, port % d\n", packetSize, Udp.remoteIP().toString().c_str(), Udp.remotePort());
    int len = Udp.read(incomingPacket, 255);
    if (len > 0) {
      incomingPacket[len] = 0;
    }

    outIp = Udp.remoteIP();
    outPort = Udp.remotePort();
    //    Serial.print(outIp);
    //    Serial.print(": ");
    //    Serial.println(outPort);

    gotmessage = true;

    if (incomingPacket[0] != commandStartChar) {
      Serial.println("UDP packet hasn`t start char");
      return;
    }
    int dataCharPos = 0;
    int endCharPos = 0;


    for (int i = 0; i < udpBufferLength; i++) {
      if (dataCharPos == 0 && incomingPacket[i] == commandDataChar) {
        dataCharPos = i;
      }
      if (endCharPos == 0 && incomingPacket[i] == commandEndChar) {
        endCharPos = i;
        break;
      }
    }
    if (endCharPos == 0 ) {
      Serial.println("UDP packet hasn`t end char");
      return;
    }

    String bufferStr = String(incomingPacket);
    String command = bufferStr.substring(1, dataCharPos);
    String data = bufferStr.substring(dataCharPos + 1, endCharPos);
    commandHandler(command, data);
  }
  //  END HANDLE COMMAND

  unsigned long  now = millis();
  //  unsigned long lastCommandDelay = now - timeLastCommangGot;
  if ((now - timeLastCommangGot) > 2000) {
    //    if 2s haven`t got messages
    timeLastCommangGot = now;


    analogWrite(D1PWM, 0);
    analogWrite(D2PWM, 0);
    Serial.print("Curently listening at IP ");
    Serial.print( WiFi.localIP().toString().c_str());
    Serial.print(":");
    Serial.println(localUdpPort);
  }

}


