import 'dart:io';

import 'package:flutter/material.dart';

//const serverUrl = 'http://10.0.2.2:9107'; // 안드로이드 에뮬레이터용
const serverUrl = 'http://192.168.0.168:9107'; // 실제 디바이스용
//const serverUrl = 'http://localhost:9107'; // 웹용
//const serverUrl = 'https://service.anbweb.kr'; // 운영 서버

const primaryColor = Color.fromARGB(255, 33, 56, 190);
const selectColor = Color.fromARGB(255, 232, 31, 99);
const backgroundColor = Color.fromARGB(255, 101, 192, 240);
const titleBackgroundColor = Color.fromARGB(255, 101, 192, 240);

String platform() {
  try {
    if (Platform.isAndroid) {
      return 'android';
    } else if (Platform.isIOS) {
      return 'ios';
    }
  } catch (e) {
    return 'web';
  }

  return '';
}

bool isExistFile(String filename) {
  if (platform() == 'web') {
    return true;
  }

  return File(filename).existsSync();
}
