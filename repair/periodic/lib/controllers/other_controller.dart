import 'dart:convert';

import 'package:common_control/common_control.dart';
import 'package:localstorage/localstorage.dart';
import 'package:periodic/controllers/blueprint_controller.dart';
import 'package:periodic/models/periodicother.dart';

import 'auth_controller.dart';

class OtherController extends GetxController {
  final _modified = false.obs;
  final _tab = 10.obs;
  final _periodicothers = <Periodicother>[].obs;

  bool get modified => _modified.value;
  set modified(value) => _modified.value = value;

  int get tab => _tab.value;
  set tab(value) => _tab.value = value;

  List<Periodicother> get periodicothers => _periodicothers;
  set periodicothers(value) => _periodicothers.value = value;

  List<TextEditingController> statusController = [];
  List<TextEditingController> positionController = [];
  final TextEditingController tab16ContentController = TextEditingController();

  // tab16 행별 이미지 저장소 (행 인덱스 0-7)
  final _tab16Images = <int, List<String>>{}.obs;
  Map<int, List<String>> get tab16Images => _tab16Images;

  List<String> getTab16Images(int rowIndex) {
    return _tab16Images[rowIndex] ?? [];
  }

  void addTab16Image(int rowIndex, String path) {
    if (!_tab16Images.containsKey(rowIndex)) {
      _tab16Images[rowIndex] = [];
    }
    _tab16Images[rowIndex]!.add(path);
    _tab16Images.refresh();

    modified = true;
    final BlueprintController blueprintController =
        Get.find<BlueprintController>();
    blueprintController.setTab16Modified(true);

    final authController = Get.find<AuthController>();
    if (authController.autosave == true) {
      save();
    }
  }

  void removeTab16Image(int rowIndex, int imageIndex) {
    if (_tab16Images.containsKey(rowIndex)) {
      _tab16Images[rowIndex]!.removeAt(imageIndex);
      _tab16Images.refresh();

      modified = true;
      final BlueprintController blueprintController =
          Get.find<BlueprintController>();
      blueprintController.setTab16Modified(true);

      final authController = Get.find<AuthController>();
      if (authController.autosave == true) {
        save();
      }
    }
  }

  updatePeriodicothers() {
    final authController = Get.find<AuthController>();
    if (authController.autosave == true) {
      save();
    }

    _periodicothers.refresh();
  }

  @override
  void onInit() async {
    super.onInit();

    final BlueprintController blueprintController =
        Get.find<BlueprintController>();

    for (var i = 0; i < blueprintController.periodicothers.length; i++) {
      Periodicother item = blueprintController.periodicothers[i];

      periodicothers.add(Periodicother(
          id: item.id,
          name: item.name,
          type: item.type,
          result: item.result,
          status: item.status,
          position: item.position,
          filename: item.filename,
          offlinefilename: item.offlinefilename,
          change: item.change,
          category: item.category,
          order: item.order,
          periodic: item.periodic,
          date: item.date));
      statusController.add(TextEditingController(text: item.status));
      positionController.add(TextEditingController(text: item.position));
    }

    // 부착물 점검내용 로드
    final LocalStorage storage16 = LocalStorage('blueprints.json');
    await storage16.ready;
    final content16 = await storage16.getItem('other_16_content');
    bool hasTab16Data = false;
    if (content16 != null && content16.toString().isNotEmpty) {
      tab16ContentController.text = content16;
      hasTab16Data = true;
    }

    // 부착물 이미지 로드
    final images16 = await storage16.getItem('other_16_images');
    if (images16 != null && images16 != '') {
      final Map<String, dynamic> decoded = json.decode(images16);
      decoded.forEach((key, value) {
        _tab16Images[int.parse(key)] = List<String>.from(value);
      });
      hasTab16Data = true;
    }

    // tab16 데이터가 있으면 modified 설정
    if (hasTab16Data) {
      blueprintController.setTab16Modified(true);
    }
  }

  saveTab16Content() async {
    modified = true;
    final BlueprintController blueprintController =
        Get.find<BlueprintController>();
    blueprintController.setTab16Modified(true);

    final authController = Get.find<AuthController>();
    if (authController.autosave == true) {
      save();
    }
  }

  save() async {
    final periodicotherStr = json.encode(periodicothers);

    final LocalStorage storage = LocalStorage('blueprints.json');
    await storage.ready;
    await storage.setItem('periodicothers', periodicotherStr);
    await storage.setItem('other_16_content', tab16ContentController.text);

    // tab16 이미지도 저장
    final Map<String, dynamic> imagesToSave = {};
    _tab16Images.forEach((key, value) {
      imagesToSave[key.toString()] = value;
    });
    await storage.setItem('other_16_images', json.encode(imagesToSave));

    final BlueprintController blueprintController =
        Get.find<BlueprintController>();

    blueprintController.periodicothers.clear();

    for (var i = 0; i < periodicothers.length; i++) {
      Periodicother item = periodicothers[i];

      blueprintController.periodicothers.add(Periodicother(
          id: item.id,
          name: item.name,
          type: item.type,
          result: item.result,
          status: item.status,
          position: item.position,
          filename: item.filename,
          offlinefilename: item.offlinefilename,
          change: item.change,
          category: item.category,
          order: item.order,
          periodic: item.periodic,
          date: item.date));
    }

    blueprintController.modified = true;
    blueprintController.modifiedOther = true;
    modified = false;
  }
}
