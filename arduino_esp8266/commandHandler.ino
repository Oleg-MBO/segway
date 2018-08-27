
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
}


void commandHandler(String command , String data ) {

  if (command == "qw") {
    //   is used to check command string to int transform
    int dataInt = data.toInt();
    sendData(command, String(dataInt));

  } else if (command == "dr1") {
    //    //      set PWM to drive 1
    //    bool isNegative = false;
    //    int dataInt;
    //    if (strstr(data.c_str(), "-")) {
    //      isNegative = true;
    //      data = data.substring( 1, data.length());
    //    }
    //    dataInt = data.toInt();
    //
    //
    //    if (dataInt > PWMRANGE) {
    //      dataInt = PWMRANGE;
    //    } else if (dataInt < 0) {
    //      dataInt = 0;
    //    }
    //    if (!isNegative) {
    //      digitalWrite (D1DIR, HIGH);
    //    } else {
    //      digitalWrite (D1DIR, LOW);
    //    }
    //    analogWrite(D1PWM, dataInt);

    //    DriveConf driveConf = parseDriveConf(data);
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
    //      set PWM to drive 2
    //    bool isNegative = false;
    //    int dataInt;
    //    if (strstr(data.c_str(), "-")) {
    //      isNegative = true;
    //      data = data.substring( 1, data.length());
    //    }
    //    dataInt = data.toInt();
    //
    //
    //    if (dataInt > PWMRANGE) {
    //      dataInt = PWMRANGE;
    //    } else if (dataInt < 0) {
    //      dataInt = 0;
    //    }
    //    if (!isNegative) {
    //      digitalWrite (D2DIR, HIGH);
    //    } else {
    //      digitalWrite (D2DIR, LOW);
    //    }
    //    analogWrite(D2PWM, dataInt);
    //
    //    DriveConf driveConf = parseDriveConf(data);
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


  } else if (command == commandSendAcc) {
    if (data == "1") {
      getAccData = true;
    } else if (data == "0") {
      getAccData = false;
    }
    sendData(commandSendAcc, getAccData ? "1" : "0");


  } else if (command == commandSendGyroAngle) {
    if (data == "1") {
      getGyroAngleData = true;
    } else if (data == "0") {
      getGyroAngleData = false;
    }
    sendData(commandSendGyroAngle, getGyroAngleData ? "1" : "0");


  } else if (command == commandSendAccAngle) {
    if (data == "1") {
      getAccAngleData = true;
    } else if (data == "0") {
      getAccAngleData = false;
    }
    sendData(commandSendAccAngle, getAccAngleData ? "1" : "0");
  }
}
