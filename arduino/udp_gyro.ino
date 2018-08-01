#include <dummy.h>

#include <Wire.h>

#include <ESP8266WiFi.h>
#include <WiFiUdp.h>

//#include <OSCMessage.h>




#define udpBufferLength 255

#define commandStartChar 62 //>
#define commandDataChar 58  //:
#define commandEndChar 60   //<

const char* ssid = "2_4G";
const char* password = "123567785962325";

WiFiUDP Udp;
unsigned int localUdpPort = 4210;  // local port to listen on
char incomingPacket[udpBufferLength];  // buffer for incoming packets
char  replyPacket[] = "Hi there! Got the message :-)";  // a reply string to send back


const long interval = 1000;           // interval at which to blink (milliseconds)
unsigned long previousMillis = 0;        // will store last time LED was updated

long counder = 0;

char buf[40];
bool gotmessage = false;
IPAddress outIp;
uint16_t outPort;

const int D1PWM = 13;
const int D1DIR = 14;

const int D2PWM = 12;
const int D2DIR = 16;

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
  while (WiFi.status() != WL_CONNECTED)
  {
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


void sendData(String  command, String  data) {

  if (gotmessage ) {
    Udp.beginPacket(outIp, outPort);

    Udp.print(">");
    Udp.print(command);
    Udp.print(":");
    Udp.print(data);
    Udp.print("<");
    Udp.endPacket();
    //    Serial.print("send udp: ");
    //    Serial.print(command);
    //    Serial.print(": ");
    //    Serial.println(data);

  }
  //  else {
  //    Serial.print(command);
  //    Serial.print(": ");
  //    Serial.println(data);
  //  }
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
  if (packetSize)  {
    timeLastCommangGot = millis();
    // receive incoming UDP packets
//    Serial.printf("Received % d bytes from % s, port % d\n", packetSize, Udp.remoteIP().toString().c_str(), Udp.remotePort());
    int len = Udp.read(incomingPacket, 255);
    if (len > 0)
    {
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


    if (command == "qw") {
      int dataInt = data.toInt();
      sendData(command, String(dataInt));

    } else if (command == "dr1") {
      bool isNegative = false;
      int dataInt;
      if (strstr(data.c_str(), "-")) {
        isNegative = true;
        data = data.substring( 1, data.length());
      }
      dataInt = data.toInt();


      if (dataInt > PWMRANGE) {
        dataInt = PWMRANGE;
      } else if (dataInt < 0) {
        dataInt = 0;
      }
      if (!isNegative) {
        digitalWrite (D1DIR, HIGH);
      } else {
        digitalWrite (D1DIR, LOW);
      }
      analogWrite(D1PWM, dataInt);
      //      if (isNegative) {
      //        dataInt = dataInt * -1;
      //      }
      //      sendData("dr1", String(dataInt));

    } else if (command == "dr2") {
      bool isNegative = false;
      int dataInt;
      if (strstr(data.c_str(), "-")) {
        isNegative = true;
        data = data.substring( 1, data.length());
      }
      dataInt = data.toInt();


      if (dataInt > PWMRANGE) {
        dataInt = PWMRANGE;
      } else if (dataInt < 0) {
        dataInt = 0;
      }
      if (!isNegative) {
        digitalWrite (D2DIR, HIGH);
      } else {
        digitalWrite (D2DIR, LOW);
      }
      analogWrite(D2PWM, dataInt);
      //      if (isNegative) {
      //        dataInt = dataInt * -1;
      //      }
      //      sendData("dr2", String(dataInt));
    }

  }
  //  END HANDE COMMAND
  unsigned long  now = millis();
  //  unsigned long lastCommandDelay = now - timeLastCommangGot;
  if ((now - timeLastCommangGot) > 2000) {
    timeLastCommangGot = now;
    //    Serial.print("Havent got commands ");
    //    Serial.print(lastCommandDelay / 1000);
    //    Serial.println(" s");

    Serial.print("Curently listening at IP ");
    Serial.print( WiFi.localIP().toString().c_str());
    Serial.print(":");
    Serial.println(localUdpPort);
  }

}


