extern const int D1PWM ;
extern const int D1DIR ;


//  from global_vars
extern const int D2PWM ;
extern const int D2DIR;

extern bool getAccData;
extern bool getGyroAngleData;
extern bool getAccAngleData;
//  end from

//  from commands
extern String commandSendAcc;
extern String commandSendGyroAngle;
extern String commandSendAccAngle;
//  end from

#define commandStartChar 62 //>
#define commandDataChar 58  //:
#define commandEndChar 60   //<


void sendData(String  command, String  data) {

  if (gotmessage ) {
    Udp.beginPacket(outIp, outPort);

    Udp.print(">");
    Udp.print(command);
    Udp.print(":");
    Udp.print(data);
    Udp.print("<");
    Udp.endPacket();

    //    Serial.print(command);
    //    Serial.print(": ");
    //    Serial.println(data);

  }
  //  else {
  //    Serial.print(command);
  //    Serial.print(": ");
  //    Serial.println(data);
  //  }
  yield();
}


void commandHandler(String command , String data ) {

  if (command == "qw") {
    //   is used to check command string to int transform
    int dataInt = data.toInt();
    sendData(command, String(dataInt));

  } else if (command == "dr") {
    // set both drives refference
    int separator = data.indexOf("|");
    String d1Data = data.substring(0, separator);
    String d2Data = data.substring(separator);

    //     HANDLE D1
    bool isNegative = false;
    int dataInt;
    if (strstr(d1Data.c_str(), "-")) {
      isNegative = true;
      data = data.substring( 1, data.length());
    }
    dataInt = data.toInt();


    if (dataInt > PWMRANGE) {
      dataInt = PWMRANGE;
    } else if (dataInt < 0) {
      dataInt = 0;
    }
    digitalWrite (D1DIR, isNegative);
    digitalWrite (D1DIR, !isNegative);
    analogWrite(D1PWM, dataInt);

    //    HANDLE D2
     isNegative = false;
     dataInt;
    if (strstr(d2Data.c_str(), "-")) {
      isNegative = true;
      data = data.substring( 1, data.length());
    }
    dataInt = data.toInt();


    if (dataInt > PWMRANGE) {
      dataInt = PWMRANGE;
    } else if (dataInt < 0) {
      dataInt = 0;
    }
    digitalWrite (D2DIR, isNegative);
    digitalWrite (D2DIR, !isNegative);
    analogWrite(D2PWM, dataInt);


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
    digitalWrite (D1DIR, isNegative);
    digitalWrite (D1DIR, !isNegative);
    analogWrite(D1PWM, dataInt);


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
    digitalWrite (D2DIR, isNegative);
    digitalWrite (D2DIR, !isNegative);
    analogWrite(D2PWM, dataInt);
  }
}
