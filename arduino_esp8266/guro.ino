extern bool getAccData;
extern bool getGyroAngleData;
extern bool getAccAngleData;

/* ============================================
  I2Cdev device library code is placed under the MIT license
  Copyright (c) 2012 Jeff Rowberg
  Permission is hereby granted, free of charge, to any person obtaining a copy
  of this software and associated documentation files (the "Software"), to deal
  in the Software without restriction, including without limitation the rights
  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
  copies of the Software, and to permit persons to whom the Software is
  furnished to do so, subject to the following conditions:
  The above copyright notice and this permission notice shall be included in
  all copies or substantial portions of the Software.
  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
  THE SOFTWARE.
  ===============================================
*/

/* This driver reads quaternion data from the MPU6060 and sends
   Open Sound Control messages.
  GY-521  NodeMCU
  MPU6050 devkit 1.0
  board   Lolin         Description
  ======= ==========    ====================================================
  VCC     VU (5V USB)   Not available on all boards so use 3.3V if needed.
  GND     G             Ground
  SCL     D1 (GPIO05)   I2C clock
  SDA     D2 (GPIO04)   I2C data
  XDA     not connected
  XCL     not connected
  AD0     not connected
  INT     D8 (GPIO15)   Interrupt pin
*/
#define INTERRUPT_PIN 15 // use pin 15 on ESP8266


//#include "I2Cdev.h"
#include <Wire.h>

#include <MPU6050_tockn.h>

MPU6050 mpu6050(Wire);



// ================================================================
// ===               INTERRUPT DETECTION ROUTINE                ===
// ================================================================

volatile bool mpuInterrupt = false;     // indicates whether MPU interrupt pin has gone high
void dmpDataReady() {
  mpuInterrupt = true;
}

void mpu_setup()
{
#if I2CDEV_IMPLEMENTATION == I2CDEV_ARDUINO_WIRE
  Wire.begin();
  Wire.setClock(400000); // 400kHz I2C clock. Comment this line if having compilation difficulties
#elif I2CDEV_IMPLEMENTATION == I2CDEV_BUILTIN_FASTWIRE
  Fastwire::setup(400, true);
#endif
  mpu6050.begin();
  pinMode(INTERRUPT_PIN, INPUT);
  attachInterrupt(digitalPinToInterrupt(INTERRUPT_PIN), dmpDataReady, RISING);


  mpu6050.calcGyroOffsets(true);
}


void mpu_loop()
{
  // if programming failed, don't try to do anything
  if (!mpuInterrupt) return;
  mpuInterrupt = false;

  mpu6050.update();

  //  Serial.print("angleX : ");
  //  Serial.print(mpu6050.getAngleX());
  //  Serial.print("\tangleY : ");
  //  Serial.print(mpu6050.getAngleY());
  //  Serial.print("\tangleZ : ");
  //  Serial.println(mpu6050.getAngleZ());



  //  SEND commandSendAngles command
  String data = String("X ") + String(mpu6050.getAngleX());
  data = data + String(";Y ");
  data = data + String(mpu6050.getAngleY());
  data = data + String(";Z ");
  data = data + String(mpu6050.getAngleZ());

  //  sendData(commandSendAngles, data );
  sendData(commandSendAngles, XYZData(mpu6050.getAngleX(), mpu6050.getAngleY(), mpu6050.getAngleZ() ));

  //  END SEND commandSendAngles comma


  if (getAccData) {
    //    String data = String("X ") + String(mpu6050.getAccX());
    //    data = data + String(";Y ");
    //    data = data + String(mpu6050.getAccY());
    //    data = data + String(";Z ");
    //    data = data + String(mpu6050.getAccZ());
    sendData(commandSendAcc, XYZData(mpu6050.getAccX(), mpu6050.getAccY(), mpu6050.getAccZ() ));
  }
  if (getGyroAngleData) {
    //    String data = String("X ") + String(mpu6050.getAccX());
    //    data = data + String(";Y ");
    //    data = data + String(mpu6050.getAccY());
    //    data = data + String(";Z ");
    //    data = data + String(mpu6050.getAccZ());
    sendData(commandSendGyroAngle, XYZData(mpu6050.getGyroAngleX(), mpu6050.getGyroAngleY(), mpu6050.getGyroAngleZ()));
  }

  if ( getAccAngleData) {
    sendData(commandSendAccAngle, XYZData(mpu6050.getAccAngleX(), mpu6050.getAccAngleY(), 0));
  }

  sendData(commandSendDataDone, "1" );
}

String XYZData(float x , float y , float z ) {
  String data = String("X ") + String(x) + String(";Y ") + String(y) + String(";Z ") + String(z);
  return data;
}


