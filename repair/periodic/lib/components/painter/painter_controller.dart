import 'dart:async';
import 'dart:convert';
import 'dart:io';
import 'dart:ui' as ui;

import 'package:flutter/material.dart';
import 'package:flutter/services.dart' show ByteData, rootBundle;
import 'package:get/get.dart';
import 'package:image_size_getter/file_input.dart';
import 'package:image_size_getter/image_size_getter.dart';
import 'package:localstorage/localstorage.dart';
import 'package:periodic/controllers/auth_controller.dart';
import 'package:periodic/controllers/blueprint_controller.dart';
import 'package:periodic/models/blueprint.dart';

const basicVertical = 1;
const basicHorizontal = 2;

const basicVerticalLine = 3;
const basicHorizontalLine = 4;

const basicVerticalBreak = 5;
const basicHorizontalBreak = 6;

const basicVerticalLineV = 7;
const basicVerticalLineFree = 11;

const basicVerticalBreakV = 9;
const basicVerticalBreakFree = 13;

const curveBlue = 31;
const curveRed = 32;
const curveGreen = 33;
const curveViolet = 34;
const lineBlue = 41;
const lineRed = 42;
const lineGreen = 43;
const lineViolet = 44;

// 크기 조절 가능한 점선 사각형 (drag로 대각선 끝점 지정)
const dashedRectRed = 51;
const dashedRectBlue = 52;

const inclinationLine = 201;
const inclinationHorizontal = 202;
const inclinationVertical = 203;

const fiberVertical = 301;
const fiberHorizontal = 302;

const materialVertical = 401;
const materialHorizontal = 402;

const crackLineRed = 121;
const crackLineBlue = 122;
const crackLineViolet = 123;
const crackLineGreen = 124;

const crackCurveRed = 126;
const crackCurveBlue = 127;
const crackCurveViolet = 128;
const crackCurveGreen = 129;

// 철근노출 X: red=130, blue=101, green=131, violet=135
const rebarXViolet = 135;
// 부식 ◇: red=132, blue=102, green=136, violet=137
const corrosionGreen = 136;
const corrosionViolet = 137;
// 보 ▲(fill): red=133, blue=103, green=138, violet=139
const beamGreen = 138;
const beamViolet = 139;
// 기타 ■(fill): red=104, blue=134, green=140, violet=141
const otherGreen = 140;
const otherViolet = 141;
// 배관누수 ●(fill mid): red=111, blue=112, green=142, violet=143
const pipeLeakGreen = 142;
const pipeLeakViolet = 143;
// 누수 ○(stroke mid): violet=115, blue=105, red=144, green=145
const leakRed = 144;
const leakGreen = 145;

class Point {
  List<Offset> items = <Offset>[];

  int width;
  LineColor color;
  DrawType type;
  int icon;
  int number;
  bool selected = false;
  bool grouped = false;

  String part = '';
  String member = '';
  String shape = '';
  String weight = '';
  String length = '';
  String count = '';
  String progress = '';
  String remark = '';
  int order = 0;
  List<String> images = <String>[];
  List<String> onlineimages = <String>[];

  Point(
      {required this.items,
      required this.color,
      required this.width,
      required this.type,
      required this.icon,
      required this.number,
      required this.part,
      required this.member,
      required this.shape,
      required this.weight,
      required this.length,
      required this.count,
      required this.progress,
      required this.remark,
      required this.order,
      required this.images,
      required this.onlineimages});

  factory Point.fromJson(Map<String, dynamic> json) {
    final items = (json['items'] as List<dynamic>)
        .map((item) => Offset(item['dx'] as double, item['dy'] as double))
        .toList();
    return Point(
        items: items,
        width: json['width'] as int,
        color: LineColor.values[(json['color'] as int) - 1],
        type: DrawType.values[(json['type'] as int) - 1],
        icon: json['icon'] as int,
        number: json['number'] as int,
        part: json['part'] as String,
        member: json['member'] as String,
        shape: json['shape'] as String,
        weight: json['weight'] as String,
        length: json['length'] as String,
        count: json['count'] as String,
        progress: json['progress'] as String,
        remark: json['remark'] as String,
        order: json['order'] as int,
        images: (json['images'] as List<dynamic>)
            .map((item) => item as String)
            .toList(),
        onlineimages: (json['onlineimages'] as List<dynamic>)
            .map((item) => item as String)
            .toList());
  }

  Map<String, dynamic> toJson() {
    return {
      'items': items.map((item) => {'dx': item.dx, 'dy': item.dy}).toList(),
      'width': width,
      'color': color.code,
      'type': type.code,
      'icon': icon,
      'number': number,
      'part': part,
      'member': member,
      'shape': shape,
      'weight': weight,
      'length': length,
      'count': count,
      'progress': progress,
      'remark': remark,
      'order': order,
      'images': images,
      'onlineimages': onlineimages
    };
  }

  add(Offset offset) {
    items.add(offset);
  }

  isInclination() {
    if (icon >= 200 && icon < 300) {
      return true;
    } else {
      return false;
    }
  }

  isBasicLine() {
    if (icon == basicVerticalLine ||
        icon == basicHorizontalLine ||
        icon == basicVerticalBreak ||
        icon == basicHorizontalBreak ||
        icon == basicVerticalLineV ||
        icon == basicVerticalLineFree ||
        icon == basicVerticalBreakV ||
        icon == basicVerticalBreakFree) {
      return true;
    } else {
      return false;
    }
  }

  isCrack() {
    if (icon == crackLineRed ||
        icon == crackLineBlue ||
        icon == crackLineViolet ||
        icon == crackLineGreen ||
        icon == crackCurveRed ||
        icon == crackCurveBlue ||
        icon == crackCurveViolet ||
        icon == crackCurveGreen) {
      return true;
    } else {
      return false;
    }
  }

  isCrackCurve() {
    if (icon == crackCurveRed ||
        icon == crackCurveBlue ||
        icon == crackCurveViolet ||
        icon == crackCurveGreen) {
      return true;
    } else {
      return false;
    }
  }

  Point clone() {
    final p = Point(
        items: List<Offset>.from(items),
        color: color,
        width: width,
        type: type,
        icon: icon,
        number: number,
        part: part,
        member: member,
        shape: shape,
        weight: weight,
        length: length,
        count: count,
        progress: progress,
        remark: remark,
        order: order,
        images: List<String>.from(images),
        onlineimages: List<String>.from(onlineimages));
    p.selected = selected;
    p.grouped = grouped;
    return p;
  }
}

enum LineColor {
  black(1, '검정'),
  red(2, '빨강'),
  blue(3, '파랑'),
  green(4, '초록'),
  lightblue(5, '하늘'),
  violet(6, '보라');

  final int code;
  final String name;
  const LineColor(this.code, this.name);
}

enum DrawType {
  curve(1, '자유곡선'),
  line(2, '직선'),
  multiline(3, '연속선'),
  icon(4, '아이콘'),
  number(5, '순번'),
  material(6, '부재'),
  txt(7, 'TEXT'),
  inclination(8, '기울기'),
  fiber(9, '강도/탄산화'),
  numberLine(10, '순번/직선'),
  ;

  final int code;
  final String name;
  const DrawType(this.code, this.name);
}

enum Mode {
  draw(1, '그리기'),
  zoom(2, '줌'),
  select(3, '선택'),
  selectEnd(3, '선택'),
  move(4, '이동'),
  moveNumber(5, '숫자이동');

  final int code;
  final String name;
  const Mode(this.code, this.name);
}

enum UndoType {
  draw(1, '그리기'),
  delete(2, '지우기');

  final int code;
  final String name;
  const UndoType(this.code, this.name);
}

class Undo {
  List<Point> points = <Point>[];
  UndoType type = UndoType.draw;

  Undo({required this.points, required this.type});
}

class PainterController extends GetxController {
  Blueprint blueprint = Blueprint();

  bool _autosave = false;
  final _toolbox = true.obs;
  final _groupbox = false.obs;
  final _databox = false.obs;
  final _inclinationbox = false.obs;
  final _fiberbox = false.obs;
  final _materialbox = false.obs;
  final _points = [].obs;
  final _undos = [].obs;
  final _works = [].obs;
  final _isLoaded = false.obs;

  final _current = (-1).obs;

  final _number = 1.obs;

  final _mode = Rx<Mode>(Mode.draw);
  final _zoom = 1.0.obs;

  final _sx = 0.0.obs;
  final _sy = 0.0.obs;
  final _currentSx = 0.0.obs;
  final _currentSy = 0.0.obs;
  final _startSx = 0.0.obs;
  final _startSy = 0.0.obs;
  final _endSx = 0.0.obs;
  final _endSy = 0.0.obs;

  final _width = 0.0.obs;
  final _height = 0.0.obs;
  final _imageWidth = 0.0.obs;
  final _imageHeight = 0.0.obs;

  final _perX = 0.0.obs;
  final _perY = 0.0.obs;

  final _color = Rx<LineColor>(LineColor.red);
  final _line = 1.obs;
  final _type = Rx<DrawType>(DrawType.numberLine);

  final _index = 3.obs;

  final _lineStart = false.obs;
  final _iconZoom = 0.0.obs;
  final _numberZoom = 0.0.obs;
  final _crackZoom = 0.0.obs;

  final _visibleGroup = false.obs;

  final _modified = false.obs;

  final _toolboxPosition = 0.0.obs;

  final _iconset = 1.obs;

  final _lastUpdateTime = 0.obs;
  int get lastUpdateTime => _lastUpdateTime.value;
  set lastUpdateTime(int value) => _lastUpdateTime.value = value;

  // 부재(DrawType.material)/강도·탄산화(icon 301~399) 번호를 같은 동(parent)의
  // 모든 층에 걸쳐 빈 번호나 중복 없이 전체 재정렬한다.
  // (한 층에서 점을 추가/삭제하거나 층에 새로 들어올 때마다 호출)
  Future<void> _renumberAcrossFloors(
      bool Function(Point point) isTarget,
      bool Function(int icon, int type) isJsonTarget,
      {void Function(Point point)? syncPoint,
      void Function(Map<String, dynamic> point, int number)? syncJson}) async {
    if (blueprint.parent == 0) {
      var localNum = 1;
      for (var i = 0; i < points.length; i++) {
        if (!isTarget(points[i])) {
          continue;
        }
        points[i].number = localNum;
        syncPoint?.call(points[i]);
        localNum++;
      }
      updatePoints();
      return;
    }

    final blueprintController = Get.find<BlueprintController>();
    final LocalStorage storage = LocalStorage('periodic.json');
    await storage.ready;

    // 같은 동(parent)에 속한 층들을 목록에 표시되는 순서(위층 -> 아래층) 그대로 사용
    final siblings = blueprintController.items
        .where((item) => item.parent == blueprint.parent)
        .toList();

    var num = 1;

    for (var s = 0; s < siblings.length; s++) {
      final item = siblings[s];

      if (item.id == blueprint.id) {
        // 현재 층: 기존 번호 순서를 유지한 채 메모리 상의 points를 재번호
        final targetIndexes = <int>[];
        for (var i = 0; i < points.length; i++) {
          if (isTarget(points[i])) {
            targetIndexes.add(i);
          }
        }

        targetIndexes
            .sort((a, b) => points[a].number.compareTo(points[b].number));

        for (var i = 0; i < targetIndexes.length; i++) {
          points[targetIndexes[i]].number = num;
          syncPoint?.call(points[targetIndexes[i]]);
          num++;
        }
        continue;
      }

      final str = await storage.getItem('data_${item.id}');
      if (str == null || str == '') {
        continue;
      }

      try {
        final j = json.decode(str) as Map<String, dynamic>;
        final list = j['points'] as List<dynamic>;

        final targetItems = <Map<String, dynamic>>[];
        for (var k = 0; k < list.length; k++) {
          final p = list[k] as Map<String, dynamic>;
          if (!isJsonTarget(p['icon'] as int, p['type'] as int)) {
            continue;
          }
          targetItems.add(p);
        }

        if (targetItems.isEmpty) {
          continue;
        }

        targetItems.sort(
            (a, b) => (a['number'] as int).compareTo(b['number'] as int));

        var changed = false;
        for (var k = 0; k < targetItems.length; k++) {
          if (targetItems[k]['number'] != num) {
            targetItems[k]['number'] = num;
            changed = true;
          }
          if (syncJson != null) {
            final before = json.encode(targetItems[k]);
            syncJson(targetItems[k], num);
            if (json.encode(targetItems[k]) != before) {
              changed = true;
            }
          }
          num++;
        }

        if (changed == true) {
          await storage.setItem('data_${item.id}', json.encode(j));
          await storage.setItem('save_${item.id}', 'y');
          blueprintController.setModified(item.id);
        }
      } catch (e) {
        continue;
      }
    }

    updatePoints();
  }

  Future<void> renumberMaterialAcrossFloors() async {
    await _renumberAcrossFloors(
        (point) => point.type == DrawType.material,
        (icon, type) => type == DrawType.material.code);
  }

  // 강도/탄산화(섬유, icon 301~399, icon=300은 데이터 저장용 더미라 제외)
  // 다른 층 추가/삭제로 번호가 밀리면 SH/N 값도 새 번호에 맞게 다시 계산해야 한다.
  Future<void> renumberFiberAcrossFloors() async {
    await _renumberAcrossFloors(
        (point) => point.icon > 300 && point.icon < 400,
        (icon, type) => icon > 300 && icon < 400,
        syncPoint: (point) {
          point.shape = '${point.number * 2 - 1}';
          point.length = '${point.number * 2}';
          point.weight = '${point.number}';
        },
        syncJson: (point, number) {
          point['shape'] = '${number * 2 - 1}';
          point['length'] = '${number * 2}';
          point['weight'] = '$number';
        });
  }

  // 마지막으로 설정한 폭 값 저장
  final _lastWeight = '0.2'.obs;
  String get lastWeight => _lastWeight.value;
  set lastWeight(String value) => _lastWeight.value = value;

  // 데이터박스 위치 (true: 상단, false: 하단)
  final _databoxTop = false.obs;
  bool get databoxTop => _databoxTop.value;
  set databoxTop(bool value) => _databoxTop.value = value;

  int get iconset => _iconset.value;
  set iconset(int value) => _iconset.value = value;

  final _tempText = ''.obs;

  String get tempText => _tempText.value;
  set tempText(String value) => _tempText.value = value;

  setIconset(value) {
    databox = false;
    inclinationbox = false;
    fiberbox = false;
    materialbox = false;

    iconset = value;

    setMode(Mode.draw);

    var index = basicVerticalLine;

    if (value == 2) {
      index = inclinationLine;
    } else if (value == 3) {
      index = fiberVertical;
    } else if (value == 4) {
      index = materialVertical;
    }

    setIndex(index);
  }

  toolboxPositionToggle() {
    if (toolboxPosition == 0) {
      toolboxPositionInit();
    } else {
      toolboxPosition = 0;
    }
  }

  toolboxPositionInit() {
    // toolboxPosition = Get.width - 235;
    toolboxPosition = Get.width - 260;
  }

  reset() {
    toolboxPositionInit();

    toolbox = true;
    groupbox = false;
    databox = false;
    inclinationbox = false;
    fiberbox = false;
    materialbox = false;
    points = [];
    _undos.clear();
    _undos.refresh();
    _works.clear();
    _works.refresh();

    isLoaded = false;

    iconset = 1;
    current = -1;

    number = 1;

    mode = Mode.draw;
    zoom = 1.0;

    sx = 0.0;
    sy = 0.0;
    currentSx = 0.0;
    currentSy = 0.0;
    startSx = 0.0;
    startSy = 0.0;
    endSx = 0.0;
    endSy = 0.0;

    perX = 0.0;
    perY = 0.0;

    color = LineColor.red;
    line = 1;
    type = DrawType.numberLine;

    setIndex(3);

    lineStart = false;

    visibleGroup = false;

    modified = false;

    final authController = Get.find<AuthController>();
    _iconZoom.value = authController.iconZoom;
    _numberZoom.value = authController.numberZoom;
    _crackZoom.value = authController.crackZoom;
  }

  bool get modified => _modified.value;
  set modified(value) {
    _modified.value = value;

    if (value == true) {
      if (_autosave == true) {
        runSave++;

        //save();
      }
    }
  }

  double get toolboxPosition => _toolboxPosition.value;
  set toolboxPosition(double value) => _toolboxPosition.value = value;

  get points => _points;
  set points(value) => _points.value = value;

  int get current => _current.value;
  set current(value) => _current.value = value;

  int get number => _number.value;
  set number(value) => _number.value = value;

  double get iconZoom {
    return _iconZoom.value;
  }

  set iconZoom(value) {
    _iconZoom.value = value;
    modified = true;

    final authController = Get.find<AuthController>();
    authController.iconZoom = value;
  }

  double get numberZoom {
    return _numberZoom.value;
  }

  set numberZoom(value) {
    _numberZoom.value = value;
    modified = true;

    final authController = Get.find<AuthController>();
    authController.numberZoom = value;
  }

  double get crackZoom {
    return _crackZoom.value;
  }

  set crackZoom(value) {
    _crackZoom.value = value;
    modified = true;

    final authController = Get.find<AuthController>();
    authController.crackZoom = value;
  }

  int get index => _index.value;
  set index(value) => _index.value = value;

  bool get lineStart => _lineStart.value;
  set lineStart(value) => _lineStart.value = value;
  LineColor get color => _color.value;
  set color(value) => _color.value = value;
  int get line => _line.value;
  set line(value) => _line.value = value;
  DrawType get type => _type.value;
  set type(value) => _type.value = value;
  bool get toolbox => _toolbox.value;
  set toolbox(value) => _toolbox.value = value;
  bool get groupbox => _groupbox.value;
  set groupbox(value) => _groupbox.value = value;
  bool get databox => _databox.value;
  set databox(value) => _databox.value = value;
  bool get inclinationbox => _inclinationbox.value;
  set inclinationbox(value) => _inclinationbox.value = value;
  bool get fiberbox => _fiberbox.value;
  set fiberbox(value) => _fiberbox.value = value;
  bool get materialbox => _materialbox.value;
  set materialbox(value) => _materialbox.value = value;
  double get perX => _perX.value;
  set perX(value) => _perX.value = value;
  double get perY => _perY.value;
  set perY(value) => _perY.value = value;
  bool get isLoaded => _isLoaded.value;
  set isLoaded(value) => _isLoaded.value = value;
  Mode get mode => _mode.value;
  set mode(value) => _mode.value = value;
  double get sx => _sx.value;
  set sx(value) => _sx.value = value;
  double get sy => _sy.value;
  set sy(value) => _sy.value = value;
  double get currentSx => _currentSx.value;
  set currentSx(value) => _currentSx.value = value;
  double get currentSy => _currentSy.value;
  set currentSy(value) => _currentSy.value = value;
  double get startSx => _startSx.value;
  set startSx(value) => _startSx.value = value;
  double get startSy => _startSy.value;
  set startSy(value) => _startSy.value = value;
  double get endSx => _endSx.value;
  set endSx(value) => _endSx.value = value;
  double get endSy => _endSy.value;
  set endSy(value) => _endSy.value = value;
  double get zoom => _zoom.value;
  set zoom(value) => _zoom.value = value;
  double get width => _width.value;
  set width(value) => _width.value = value;
  double get height => _height.value;
  set height(value) => _height.value = value;
  double get imageWidth => _imageWidth.value;
  set imageWidth(value) => _imageWidth.value = value;
  double get imageHeight => _imageHeight.value;
  set imageHeight(value) => _imageHeight.value = value;

  bool get visibleGroup => _visibleGroup.value;
  set visibleGroup(value) => _visibleGroup.value = value;

  late ui.Image image;
  var events = StreamController<int>();

  startPoint(Offset offset) async {
    var num = 0;
    // 아이콘 범위 우선 체크 (결함도/기울기/섬유/재료 구분)
    if (index >= 200 && index < 300) {
      // 기울기 (200-299): 기울기 아이콘만 체크
      for (var i = 0; i < points.length; i++) {
        if (points[i].icon < 200 || points[i].icon >= 300) {
          continue;
        }

        if (points[i].number > num) {
          num = points[i].number;
        }
      }

      num++;
    } else if (index >= 300 && index < 400) {
      // 섬유 (300-399): 섬유 아이콘만 체크, icon=300은 데이터 저장용 더미 제외
      for (var i = 0; i < points.length; i++) {
        if (points[i].icon <= 300 || points[i].icon >= 400) {
          continue;
        }

        if (points[i].number > num) {
          num = points[i].number;
        }
      }

      num++;
    } else if (index >= 400 && index < 500) {
      // 재료 (400-499): 재료 아이콘만 체크
      // (같은 동의 다른 층과의 전체 번호 재정렬은 추가 직후 renumberMaterialAcrossFloors()에서 처리)
      for (var i = 0; i < points.length; i++) {
        if (points[i].icon < 400 || points[i].icon >= 500) {
          continue;
        }

        if (points[i].number > num) {
          num = points[i].number;
        }
      }

      num++;
    } else if (type == DrawType.number || type == DrawType.numberLine) {
      // 결함도 (100-199): number/numberLine 타입이면서 100-199 범위만 체크
      for (var i = 0; i < points.length; i++) {
        if (points[i].type != DrawType.number &&
            points[i].type != DrawType.numberLine) {
          continue;
        }
        // 결함도 범위만 체크 (기울기 200-299 제외)
        if (points[i].icon >= 200 && points[i].icon < 300) {
          continue;
        }

        if (points[i].number > num) {
          num = points[i].number;
        }
      }

      num++;
    } else if (type == DrawType.material) {
      for (var i = 0; i < points.length; i++) {
        if (points[i].type != DrawType.material) {
          continue;
        }

        if (points[i].number > num) {
          num = points[i].number;
        }
      }

      num++;
    }

    var point = Point(
        items: [],
        color: color,
        width: 1,
        type: type,
        icon: index,
        number: num,
        part: '',
        member: '',
        shape: '',
        weight: '',
        length: '',
        count: '',
        progress: '',
        remark: '',
        order: 0,
        images: [],
        onlineimages: []);

    if (type == DrawType.number || type == DrawType.numberLine) {
      point.weight = lastWeight; // 마지막 설정값 사용

      if (index == 1 ||
          index == 3 ||
          index == 5 ||
          index == 7 ||
          index == 9 ||
          index == 11 ||
          index == 13) {
        point.member = '벽체';
      } else {
        point.member = '슬래브';
      }

      point.count = '1';
      if (points.length > 0) {
        for (var i = points.length - 1; i >= 0; i--) {
          if (points[i].type == DrawType.number ||
              points[i].type == DrawType.numberLine) {
            point.part = points[i].part;
            break;
          }
        }
      }
    }
    point.add(offset);
    points.add(point);

    if (index > 300 && index < 400) {
      // 강도/탄산화도 부재처럼 같은 동의 모든 층이 측정위치 번호를 공유하도록
      // 전체 재정렬하고, 그 안에서 최종 번호를 기준으로 SH/N 값도 함께 갱신한다.
      // (undo 스냅샷을 찍기 전에 최종 값을 확정해야 undo 시 값이 비지 않는다)
      await renumberFiberAcrossFloors();
    }

    if (type == DrawType.material) {
      await renumberMaterialAcrossFloors();
    }

    clearUndo();
    _works.add(_snapshotPoints());
    _works.refresh();

    if (type == DrawType.number || type == DrawType.numberLine) {
      viewDatabox(points.length - 1);
    } else if (index >= 200 && index < 300) {
      viewDataboxInclination(points.length - 1);
    } else if (index >= 300 && index < 400) {
      viewDataboxFiber(points.length - 1);
    } else if (index >= 400 && index < 500) {
      viewDataboxMaterial(points.length - 1);
    } else {
      databox = false;
      inclinationbox = false;
      fiberbox = false;
      materialbox = false;
    }

    if (type != DrawType.curve &&
        type != DrawType.line &&
        type != DrawType.multiline) {
      modified = true;
    }

    updatePoints();
  }

  addPoint(Offset offset) {
    points[points.length - 1].add(offset);

    clearUndo();

    _works[_works.length - 1] = _snapshotPoints();
    _works.refresh();
  }

  movePoint(Offset offset) {
    points[points.length - 1]
        .items[points[points.length - 1].items.length - 1] = offset;

    clearUndo();

    _works[_works.length - 1] = _snapshotPoints();
    _works.refresh();
  }

  // points 리스트를 깊은 복사하여 스냅샷으로 반환
  // (undo/redo 히스토리가 같은 Point 인스턴스를 공유해 번호 등이 변형되는 것을 방지)
  List<Point> _snapshotPoints() {
    return points.map<Point>((p) => (p as Point).clone()).toList();
  }

  pointClear() {
    points.clear();

    _works.add(_snapshotPoints());
    _works.refresh();

    modified = true;
  }

  initEvent() {
    events.close();
    events = StreamController<int>();
  }

  updateCanvas() {
    events.add(points.length);
  }

  clear() {
    pointClear();
    _works.clear();
    _works.refresh();
    _undos.clear();
    _undos.refresh();

    sx = 0.0;
    sy = 0.0;
    initZoom();
    setIndex(3);
    final authController = Get.find<AuthController>();
    _iconZoom.value = authController.iconZoom;
    _numberZoom.value = authController.numberZoom;
    _crackZoom.value = authController.crackZoom;
    setMode(Mode.draw);
  }

  clearUndo() {
    _undos.clear();
    _undos.refresh();
  }

  initZoom() {
    var zoom = 1.0;

    var rate = width / height;
    var imageRate = imageWidth / imageHeight;

    if (rate > imageRate) {
      zoom = height / imageHeight;
    } else {
      zoom = width / imageWidth;
    }

    setZoom(zoom);
  }

  setZoom(zoom) {
    if (imageWidth * zoom < width && imageHeight * zoom < height) {
      var rate = width / height;
      var imageRate = imageWidth / imageHeight;

      if (rate > imageRate) {
        zoom = height / imageHeight;
      } else {
        zoom = width / imageWidth;
      }
    }

    if (zoom > 5.0) {
      zoom = 5.0;
    }

    sx = ((imageWidth * zoom * perX) - startSx) / zoom;
    sy = ((imageHeight * zoom * perY) - startSy) / zoom;

    this.zoom = zoom;

    if (imageWidth * zoom < width) {
      sx = 0.0;
    } else {
      if (sx * zoom >= imageWidth * zoom - width) {
        sx = (imageWidth * zoom - width) / zoom;
      }
    }

    if (imageHeight * zoom < height) {
      sy = 0.0;
    } else {
      if (sy * zoom >= imageHeight * zoom - height) {
        sy = (imageHeight * zoom - height) / zoom;
      }
    }

    if (sx < 0.0) {
      sx = 0.0;
    }

    if (sy < 0.0) {
      sy = 0.0;
    }

    updateCanvas();
  }

  toolboxToggle() {
    toolbox = !toolbox;
  }

  setMode(value) {
    mode = value;

    lineStart = false;

    if (mode == Mode.draw || mode == Mode.zoom || mode == Mode.move) {
      startSx = 0.0;
      startSy = 0.0;
      endSx = 0.0;
      endSy = 0.0;
    }

    if (mode == Mode.selectEnd) {
      _points.refresh();
    }

    if (mode == Mode.draw || mode == Mode.zoom) {
      clearSelection();
      databox = false;
      inclinationbox = false;
      fiberbox = false;
      materialbox = false;
    } else {
      var icon = 0;
      var count = 0;
      var curr = 0;
      for (var i = 0; i < points.length; i++) {
        var item = points[i];

        if (item.type != DrawType.number && item.type != DrawType.numberLine) {
          continue;
        }

        if (item.selected == true) {
          count++;
          icon = item.icon;
          curr = i;
        }
      }

      if (count == 1) {
        if (icon >= 200 && icon < 300) {
          inclinationbox = true;
        } else if (icon >= 300 && icon < 400) {
          fiberbox = true;
        } else if (icon >= 400 && icon < 500) {
          materialbox = true;
        } else {
          databox = true;
        }
        current = curr;
      } else {
        databox = false;
        inclinationbox = false;
        fiberbox = false;
        materialbox = false;
      }
    }

    updateCanvas();
  }

  modeToggle() {
    if (mode == Mode.draw) {
      mode = Mode.zoom;
    } else {
      mode = Mode.draw;
    }

    lineStart = false;
  }

  setColor(value) {
    color = value;

    lineStart = false;
  }

  setType(value) {
    type = value;

    lineStart = false;
  }

  undo() {
    if (_works.length <= 1) {
      return;
    }

    databox = false;
    inclinationbox = false;
    fiberbox = false;
    materialbox = false;

    var last = _works[_works.length - 1];
    _undos.add(last);
    _works.removeLast();

    _undos.refresh();
    _works.refresh();

    if (_works.isEmpty) {
      points.clear();
    } else {
      points = _works[_works.length - 1]
          .map<Point>((p) => (p as Point).clone())
          .toList();
      modified = true;
    }

    updateCanvas();
  }

  redo() {
    if (_undos.isEmpty) {
      return;
    }

    databox = false;
    inclinationbox = false;
    fiberbox = false;
    materialbox = false;

    var last = _undos[_undos.length - 1];
    // points를 먼저 복원한 후 스냅샷을 만들어 _works에 추가
    // (참조 공유를 방지하여 undo/redo 시 데이터 일관성 유지)
    points = (last as List).map<Point>((p) => (p as Point).clone()).toList();
    _works.add(_snapshotPoints());
    modified = true;
    _undos.removeLast();

    _works.refresh();
    _undos.refresh();

    updateCanvas();
  }

  setIndex(value) {
    index = value;

    if (value == basicVertical) {
      setColor(LineColor.red);
      setType(DrawType.number);
    } else if (value == basicHorizontal) {
      setColor(LineColor.blue);
      setType(DrawType.number);
    } else if (value == curveBlue) {
      setColor(LineColor.lightblue);
      setType(DrawType.curve);
    } else if (value == curveRed) {
      setColor(LineColor.red);
      setType(DrawType.curve);
    } else if (value == curveGreen) {
      setColor(LineColor.green);
      setType(DrawType.curve);
    } else if (value == curveViolet) {
      setColor(LineColor.violet);
      setType(DrawType.curve);
    } else if (value == lineBlue) {
      setColor(LineColor.lightblue);
      setType(DrawType.line);
    } else if (value == lineRed) {
      setColor(LineColor.red);
      setType(DrawType.line);
    } else if (value == lineGreen) {
      setColor(LineColor.green);
      setType(DrawType.line);
    } else if (value == lineViolet) {
      setColor(LineColor.violet);
      setType(DrawType.line);
    } else if (value == dashedRectRed) {
      setColor(LineColor.red);
      setType(DrawType.line);
    } else if (value == dashedRectBlue) {
      setColor(LineColor.blue);
      setType(DrawType.line);
    } else if (value == inclinationLine) {
      setColor(LineColor.red);
      setType(DrawType.line);
    } else if (value == inclinationHorizontal) {
      setColor(LineColor.red);
      setType(DrawType.line);
    } else if (value == inclinationVertical) {
      setColor(LineColor.red);
      setType(DrawType.line);
    } else if (value == fiberVertical) {
      setColor(LineColor.red);
      setType(DrawType.icon);
    } else if (value == fiberHorizontal) {
      setColor(LineColor.blue);
      setType(DrawType.icon);
    } else if (value == materialVertical) {
      setColor(LineColor.red);
      setType(DrawType.material);
    } else if (value == materialHorizontal) {
      setColor(LineColor.blue);
      setType(DrawType.material);
    } else if (value == basicHorizontalLine) {
      setColor(LineColor.blue);
      setType(DrawType.numberLine);
    } else if (value == basicVerticalLine) {
      setColor(LineColor.red);
      setType(DrawType.numberLine);
    } else if (value == basicHorizontalBreak) {
      setColor(LineColor.blue);
      setType(DrawType.numberLine);
    } else if (value == basicVerticalBreak) {
      setColor(LineColor.red);
      setType(DrawType.numberLine);
    } else if (value == basicVerticalLineV) {
      setColor(LineColor.red);
      setType(DrawType.numberLine);
    } else if (value == basicVerticalBreakV) {
      setColor(LineColor.red);
      setType(DrawType.numberLine);
    } else if (value == basicVerticalLineFree) {
      setColor(LineColor.red);
      setType(DrawType.numberLine);
    } else if (value == basicVerticalBreakFree) {
      setColor(LineColor.red);
      setType(DrawType.numberLine);
    } else if (value == crackLineRed) {
      setColor(LineColor.red);
      setType(DrawType.line);
    } else if (value == crackLineBlue) {
      setColor(LineColor.red);
      setType(DrawType.line);
    } else if (value == crackLineViolet) {
      setColor(LineColor.violet);
      setType(DrawType.line);
    } else if (value == crackLineGreen) {
      setColor(LineColor.green);
      setType(DrawType.line);
    } else if (value == crackCurveRed) {
      setColor(LineColor.red);
      setType(DrawType.curve);
    } else if (value == crackCurveBlue) {
      setColor(LineColor.red);
      setType(DrawType.curve);
    } else if (value == crackCurveViolet) {
      setColor(LineColor.violet);
      setType(DrawType.curve);
    } else if (value == crackCurveGreen) {
      setColor(LineColor.green);
      setType(DrawType.curve);
    } else {
      setType(DrawType.icon);
    }
  }

  setIconZoom(zoom) {
    if (zoom <= 0.0) {
      return;
    }

    iconZoom = zoom;

    updateCanvas();
  }

  setNumberZoom(zoom) {
    if (zoom <= 0.0) {
      return;
    }

    numberZoom = zoom;

    updateCanvas();
  }

  setCrackZoom(zoom) {
    if (zoom <= 0.0) {
      return;
    }

    crackZoom = zoom;

    updateCanvas();
  }

  clearSelection() {
    for (var i = 0; i < points.length; i++) {
      var point = points[i];

      point.selected = false;
      point.grouped = false;
    }

    updateCanvas();
  }

  deleteSelection() async {
    var flag = true;
    var find = false;

    while (flag == true) {
      flag = false;

      for (var i = 0; i < points.length; i++) {
        var point = points[i];

        if (point.selected != true) {
          continue;
        }

        points.removeAt(i);
        flag = true;
        find = true;
        break;
      }

      databox = false;
      inclinationbox = false;
      fiberbox = false;
      materialbox = false;
    }

    renumberMaterial();
    await renumberMaterialAcrossFloors();
    await renumberFiberAcrossFloors();
    setMode(Mode.selectEnd);
    groupSort();
    _works.add(_snapshotPoints());
    _works.refresh();
    updateCanvas();

    if (find == true) {
      modified = true;
    }
  }

  // 도면(캔버스)에 그려질 번호를 계산한다.
  // 부재(DrawType.material)는 표(측정위치)의 연속 번호와 별개로
  // 층마다 해당 층 안에서만 1부터 다시 매긴 번호를 보여준다.
  int getDisplayNumber(Point point) {
    if (point.type == DrawType.material) {
      var num = 0;

      for (var i = 0; i < points.length; i++) {
        if (points[i].type != DrawType.material) {
          continue;
        }

        num++;

        if (identical(points[i], point)) {
          return num;
        }
      }
    }

    return point.number;
  }

  renumberMaterial() {
    // 기울기 재번호 (200-299)
    var num = 1;
    for (var i = 0; i < points.length; i++) {
      Point point = points[i];
      if (point.icon < 200 || point.icon >= 300) {
        continue;
      }

      points[i].number = num;
      num++;
    }

  }

  deleteSelectionWithoutNumber() async {
    var flag = true;
    var find = false;

    while (flag == true) {
      flag = false;

      for (var i = 0; i < points.length; i++) {
        var point = points[i];

        if (point.type == DrawType.number ||
            point.type == DrawType.numberLine) {
          continue;
        }

        if (point.selected != true) {
          continue;
        }

        points.removeAt(i);
        flag = true;
        find = true;

        break;
      }
    }

    databox = false;
    inclinationbox = false;
    fiberbox = false;
    materialbox = false;

    renumberMaterial();
    await renumberMaterialAcrossFloors();
    await renumberFiberAcrossFloors();
    setMode(Mode.selectEnd);
    _works.add(_snapshotPoints());
    _works.refresh();
    updateCanvas();

    if (find == true) {
      modified = true;
    }
  }

  groupSort() {
    var flag = true;

    while (flag == true) {
      flag = false;

      for (var i = 0; i < points.length; i++) {
        var find = false;

        for (var j = 0; j < points.length; j++) {
          var item2 = points[j];

          // 부재(DrawType.material) 번호는 동(parent) 전체에 걸쳐 통합 관리되므로
          // 이 층 안의 번호 압축 로직 대상에서 제외한다.
          if (item2.type == DrawType.material) {
            continue;
          }

          // 강도/탄산화(icon 301~399) 번호도 동(parent) 전체에 걸쳐 통합 관리되므로 제외한다.
          if (item2.icon > 300 && item2.icon < 400) {
            continue;
          }

          if (item2.number == i + 1) {
            find = true;
            break;
          }
        }

        if (find == true) {
          continue;
        }

        for (var j = 0; j < points.length; j++) {
          var item2 = points[j];

          if (item2.type == DrawType.material) {
            continue;
          }

          if (item2.icon > 300 && item2.icon < 400) {
            continue;
          }

          if (item2.number > i + 1) {
            flag = true;
            points[j].number--;
          }
        }

        break;
      }
    }
  }

  groupApply() {
    var max = 0;

    for (var i = 0; i < points.length; i++) {
      var item = points[i];

      if (item.grouped == false) {
        if (item.number > max) {
          max = item.number;
        }

        continue;
      }

      points[i].number = 0;
    }

    var number = max + 1;

    for (var i = 0; i < points.length; i++) {
      var flag = false;

      for (var j = 0; j < points.length; j++) {
        var item2 = points[j];

        if (item2.number == i + 1) {
          flag = true;
          break;
        }
      }

      if (flag == true) {
        continue;
      }

      number = i + 1;
      break;
    }

    for (var i = 0; i < points.length; i++) {
      var item = points[i];

      if (item.grouped == false) {
        continue;
      }

      points[i].number = number;
    }

    groupSort();

    updatePoints();

    modified = true;
  }

  groupCancel() {
    var max = 0;

    for (var i = 0; i < points.length; i++) {
      var item = points[i];

      if (item.grouped == false) {
        if (item.number > max) {
          max = item.number;
        }

        continue;
      }

      points[i].number = 0;
    }

    for (var i = 0; i < points.length; i++) {
      var item = points[i];

      if (item.grouped == false) {
        continue;
      }

      max++;
      points[i].number = max;
    }

    groupSort();

    updatePoints();
    modified = true;
  }

  updatePoints() {
    points.refresh();
    updateCanvas();
  }

  // 이미지 추가/삭제 등 데이터 변경 시 현재 스냅샷 업데이트
  updateDataSnapshot() {
    if (_works.isNotEmpty) {
      _works[_works.length - 1] = _snapshotPoints();
      _works.refresh();
    }
    clearUndo();
  }

  viewDatabox(int pos) {
    if (iconset != 1) {
      return;
    }

    databox = true;
    current = pos;
  }

  viewDataboxInclination(int pos) {
    inclinationbox = true;
    current = pos;
  }

  viewDataboxFiber(int pos) {
    fiberbox = true;
    current = pos;
  }

  viewDataboxMaterial(int pos) {
    materialbox = true;
    current = pos;
  }

  setData(int pos, String value) {
    if (pos == 1) {
      points[current].part = value;
    } else if (pos == 2) {
      points[current].member = value;
    } else if (pos == 3) {
      points[current].shape = value;
    } else if (pos == 4) {
      points[current].weight = value;
      lastWeight = value; // 마지막 폭 값 저장
    } else if (pos == 5) {
      points[current].length = value;
    } else if (pos == 6) {
      points[current].count = value;
    } else if (pos == 7) {
      points[current].progress = value;
    } else if (pos == 8) {
      points[current].remark = value;
    }

    // 데이터 변경은 현재 지시선과 같은 undo 단계로 처리
    // (새 스냅샷 추가 대신 마지막 스냅샷을 덮어써서 undo 한 번에 지시선+데이터 함께 되돌림)
    if (_works.isNotEmpty) {
      _works[_works.length - 1] = _snapshotPoints();
      _works.refresh();
    }
    clearUndo();

    updatePoints();
    modified = true;
  }

  setCurrent(Point value) {
    points[current] = value;

    // 데이터 변경은 현재 지시선과 같은 undo 단계로 처리
    if (_works.isNotEmpty) {
      _works[_works.length - 1] = _snapshotPoints();
      _works.refresh();
    }
    clearUndo();
    modified = true;

    updatePoints();
  }

  Point getCurrent() {
    if (current == -1) {
      return Point(
          items: [],
          color: LineColor.black,
          icon: 0,
          number: 0,
          type: DrawType.line,
          width: 0,
          part: '',
          member: '',
          shape: '',
          weight: '',
          length: '',
          count: '',
          progress: '',
          remark: '',
          order: 0,
          images: [],
          onlineimages: []);
    }

    return points[current];
  }

  save() async {
    Map<String, dynamic> items = {
      'id': blueprint.id,
      'points': points.map((item) => item.toJson()).toList(),
      'zoom': zoom,
      'iconzoom': iconZoom,
      'numberzoom': numberZoom,
      'crackzoom': crackZoom,
    };

    final LocalStorage storage = LocalStorage('periodic.json');
    final str = json.encode(items);

    final ready = await storage.ready;
    if (ready) {
      try {
        await storage.setItem('data_${blueprint.id}', str);
        await storage.setItem('save_${blueprint.id}', 'y');
      } catch (e) {
        //print(e);
      }
    } else {
      //print('not ready');
    }

    modified = false;
    final blueprintController = Get.find<BlueprintController>();
    blueprintController.setModified(blueprint.id);
    blueprintController.modified = true;

    if (_autosave == false) {
      databox = false;
      inclinationbox = false;
      fiberbox = false;
      materialbox = false;
    }
  }

  void load() async {
    final authController = Get.find<AuthController>();
    _autosave = authController.autosave;

    final LocalStorage storage = LocalStorage('periodic.json');

    await storage.ready;
    final str = await storage.getItem('data_${blueprint.id}');

    if (str == null || str == '') {
      modified = false;
      _works.add(_snapshotPoints());
      _works.refresh();
      await renumberMaterialAcrossFloors();
      await renumberFiberAcrossFloors();
      return;
    }

    final j = json.decode(str);

    zoom = j['zoom'];
    _iconZoom.value = j['iconzoom'];

    if (_iconZoom.value == 0) {
      _iconZoom.value = authController.iconZoom;
      if (_iconZoom.value == 0) {
        _iconZoom.value = 100;
      }
    }

    try {
      _numberZoom.value = j['numberzoom'];
    } catch (e) {
      _numberZoom.value = 100;
    }

    if (_numberZoom.value == 0) {
      _numberZoom.value = authController.numberZoom;
      if (_numberZoom.value == 0) {
        _numberZoom.value = 100;
      }
    }

    try {
      _crackZoom.value = j['crackzoom'];
    } catch (e) {
      _crackZoom.value = 100;
    }

    if (_crackZoom.value == 0) {
      _crackZoom.value = authController.crackZoom;
      if (_crackZoom.value == 0) {
        _crackZoom.value = 100;
      }
    }

    var items = j['points'].map<Point>((json) => Point.fromJson(json)).toList();

    for (var i = 0; i < items.length; i++) {
      List<String> newImages = [];

      for (var j = 0; j < items[i].images.length; j++) {
        var image = items[i].images[j];
        if (File(image).existsSync() == true) {
          newImages.add(image);
        }
      }

      items[i].images = newImages;
    }

    _points.value = items;
    _works.add(_snapshotPoints());
    _works.refresh();
    _points.refresh();
    updateCanvas();
    modified = false;

    await renumberMaterialAcrossFloors();
    await renumberFiberAcrossFloors();
  }

  removeDataimage(pos) {
    removeDataimageIndex(current, pos);
  }

  removeDataimageIndex(index, pos) {
    Point point = points[index];
    point.images.removeAt(pos);

    // 이미지 삭제도 데이터 변경과 같은 undo 단계로 처리
    updateDataSnapshot();

    updatePoints();
    modified = true;
  }

  removePoint(index) async {
    points.removeAt(index);
    renumberMaterial();
    await renumberMaterialAcrossFloors();
    await renumberFiberAcrossFloors();
    groupSort();
    updatePoints();
    modified = true;
  }

  isUndo() {
    if (_works.length > 1) {
      return true;
    } else {
      return false;
    }
  }

  isRedo() {
    if (_undos.isEmpty) {
      return false;
    } else {
      return true;
    }
  }

  final _runSave = 1.obs;
  int get runSave => _runSave.value;
  set runSave(value) => _runSave.value = value;
  @override
  void onInit() async {
    super.onInit();

    final authController = Get.find<AuthController>();
    _iconZoom.value = authController.iconZoom;
    _numberZoom.value = authController.numberZoom;
    _crackZoom.value = authController.crackZoom;

    debounce(_runSave, (_) {
      save();
    }, time: const Duration(milliseconds: 300));
  }

  // 이전/다음 층 이동 (mode: 1=이전, 2=다음)
  // 이동 후 즉시 새 도면 이미지를 로드하여 캔버스를 갱신한다.
  // 반환값: 실제로 이동했으면 true
  Future<bool> moveFloor(int mode) async {
    databox = false;
    inclinationbox = false;
    fiberbox = false;
    materialbox = false;

    final blueprintController = Get.find<BlueprintController>();
    final authController = Get.find<AuthController>();

    Blueprint prev = Blueprint();
    Blueprint current = Blueprint();

    var find = false;
    for (var i = 0; i < blueprintController.items.length; i++) {
      final item = blueprintController.items[i];
      if (find == true) {
        if (mode == 2) {
          if (item.upload == 1 && item.offlinefilename != '') {
            current = item;
            break;
          }
        } else {
          current = prev;
          break;
        }
      } else {
        if (item.id == blueprint.id) {
          find = true;
          if (mode == 1) {
            current = prev;
            break;
          }
        } else {
          if (item.upload == 1 && item.offlinefilename != '') {
            prev = item;
          }
        }
      }
    }

    if (current.id == 0) {
      return false;
    }

    authController.setTitle(current);
    blueprint = current;
    reset();
    load();

    await loadImage(current.offlinefilename);

    return true;
  }

  // 도면 이미지를 디스크/번들에서 로드하여 c.image 에 반영한다.
  Future<void> loadImage(String imageAssetPath) async {
    if (imageAssetPath.isEmpty) return;

    // 1. 이미지 사이즈 측정
    ui.Size size = const ui.Size(0, 0);
    if (Platform.isAndroid) {
      try {
        final file = File(imageAssetPath);
        final s = ImageSizeGetter.getSize(FileInput(file));
        size = ui.Size(s.width.toDouble(), s.height.toDouble());
      } catch (_) {
        size = const ui.Size(0, 0);
      }
    } else {
      // iOS는 Image 위젯의 resolve 콜백을 사용 (Painter.getImageSizeIOS와 동일)
      final completer = Completer<ui.Size>();
      final image = Image.asset(imageAssetPath);
      image.image.resolve(const ImageConfiguration()).addListener(
            ImageStreamListener((info, _) {
              completer.complete(ui.Size(
                  info.image.width.toDouble(), info.image.height.toDouble()));
            }, onError: (_, __) {
              if (!completer.isCompleted) {
                completer.complete(const ui.Size(0, 0));
              }
            }),
          );
      size = await completer.future;
    }

    imageWidth = size.width;
    imageHeight = size.height;
    initZoom();

    // 2. 실제 이미지 디코드
    ByteData data;
    if (Platform.isAndroid) {
      final file = File(imageAssetPath);
      final bytes = file.readAsBytesSync();
      data = bytes.buffer.asByteData();
    } else {
      data = await rootBundle.load(imageAssetPath);
    }
    final codec = await ui.instantiateImageCodec(
      data.buffer.asUint8List(),
      targetWidth: imageWidth.toInt(),
      targetHeight: imageHeight.toInt(),
    );
    final frame = await codec.getNextFrame();
    image = frame.image;
    isLoaded = true;
    updateCanvas();
  }
}
