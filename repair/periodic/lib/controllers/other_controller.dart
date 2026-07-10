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
  final _isLoading = true.obs; // 로딩 상태 추가

  bool get modified => _modified.value;
  set modified(value) => _modified.value = value;

  int get tab => _tab.value;
  set tab(value) => _tab.value = value;

  List<Periodicother> get periodicothers => _periodicothers;
  set periodicothers(value) => _periodicothers.value = value;

  bool get isLoading => _isLoading.value; // 로딩 상태 getter
  set isLoading(value) => _isLoading.value = value;

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

  // tab16 상태 체크박스/라디오 저장소 (행 인덱스 0-7)
  final _tab16States = <int, List<String>>{}.obs;
  Map<int, List<String>> get tab16States => _tab16States;

  List<String> getTab16State(int rowIndex) {
    return _tab16States[rowIndex] ?? [];
  }

  bool isTab16Checked(int rowIndex, String option) {
    return getTab16State(rowIndex).contains(option);
  }

  void toggleTab16Checkbox(int rowIndex, String option) {
    final list = List<String>.from(_tab16States[rowIndex] ?? []);
    if (list.contains(option)) {
      list.remove(option);
    } else {
      list.add(option);
    }
    _tab16States[rowIndex] = list;
    _tab16States.refresh();
    _onTab16StateChanged();
  }

  void setTab16Radio(int rowIndex, String option) {
    _tab16States[rowIndex] = [option];
    _tab16States.refresh();
    _onTab16StateChanged();
  }

  void _onTab16StateChanged() {
    modified = true;
    final BlueprintController blueprintController =
        Get.find<BlueprintController>();
    blueprintController.setTab16Modified(true);

    final authController = Get.find<AuthController>();
    if (authController.autosave == true) {
      save();
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
    try {
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

      // 부착물 상태(체크박스/라디오) 로드
      final states16 = await storage16.getItem('other_16_states');
      if (states16 != null && states16 != '') {
        final Map<String, dynamic> decoded = json.decode(states16);
        decoded.forEach((key, value) {
          _tab16States[int.parse(key)] = List<String>.from(value);
        });
        hasTab16Data = true;
      }

      // LocalStorage에 state가 없으면 periodicothers(category=16)에서 로드
      // (다른 기기/서버 재로드 대비)
      if (_tab16States.isEmpty) {
        for (var item in periodicothers) {
          if (item.category != 16) continue;
          final rowIndex = item.order - 160;
          if (rowIndex < 0 || rowIndex > 7) continue;
          if (item.status.isEmpty) continue;
          _tab16States[rowIndex] = item.status.split(',');
        }
      }

      // tab16 데이터가 있으면 modified 설정
      if (hasTab16Data) {
        blueprintController.setTab16Modified(true);
      }

      // 데이터 로드 완료 후 UI 업데이트 명시적으로 트리거
      _tab16Images.refresh();
      _tab16States.refresh();
    } catch (e) {
      print('Error loading tab16 data: $e');
    } finally {
      // 로딩 완료
      isLoading = false;
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
    // _tab16States 값을 category=16 periodicothers 로 동기화
    // (Row 0 → Order 160, Row N → Order 160+N)
    _syncTab16StatesToPeriodicothers();

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

    // tab16 상태(체크박스/라디오)도 저장
    final Map<String, dynamic> statesToSave = {};
    _tab16States.forEach((key, value) {
      statesToSave[key.toString()] = value;
    });
    await storage.setItem('other_16_states', json.encode(statesToSave));

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

  // _tab16States (row 0-7) → periodicothers (Category=16, Order 160+N).status 로 동기화
  // Row 0: 부착물 등 (radio a-e), Row 1-7: 각 부재 체크박스
  void _syncTab16StatesToPeriodicothers() {
    for (var i = 0; i < periodicothers.length; i++) {
      final item = periodicothers[i];
      if (item.category != 16) continue;

      final rowIndex = item.order - 160;
      if (rowIndex < 0 || rowIndex > 7) continue;

      final states = _tab16States[rowIndex] ?? [];
      final newStatus = states.join(',');

      if (item.status != newStatus) {
        periodicothers[i].status = newStatus;
        periodicothers[i].change = 1;
      }
    }
  }
}
