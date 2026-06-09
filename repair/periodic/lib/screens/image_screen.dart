import 'dart:async';
import 'dart:io';

import 'package:common_control/common_control.dart';
import 'package:flutter/cupertino.dart';
import 'package:image_picker/image_picker.dart';
import 'package:periodic/components/default_app_bar.dart';
import 'package:periodic/controllers/auth_controller.dart';
import 'package:periodic/controllers/blueprint_controller.dart';
import 'package:periodic/controllers/image_controller.dart';
import 'package:periodic/models/periodicimage.dart';

class ImageScreen extends CWidget {
  ImageScreen({super.key});

  final AuthController authController = Get.find<AuthController>();
  final BlueprintController blueprintController =
      Get.find<BlueprintController>();
  final c = Get.find<ImageController>();

  final TextEditingController textEditingController = TextEditingController();
  @override
  Widget build(BuildContext context) {
    // arguments로 type이 전달되면 설정
    final arguments = Get.arguments;
    if (arguments != null && arguments['type'] != null) {
      c.type = arguments['type'];
    }

    return CScaffold(
      appBar: DefaultAppBar(
        leading: IconButton(
            icon: const Icon(Icons.arrow_back_ios),
            onPressed: () {
              Get.back();
            }),
      ),
      backgroundColor: Colors.white,
      body: body(context),
    );
  }

  body(context) {
    return CColumn(padding: const EdgeInsets.all(10), children: [
      Expanded(
          child: SingleChildScrollView(
              scrollDirection: Axis.vertical,
              child: Obx(() => Wrap(children: _makeTabs(context))))),
      bottom(),
    ]);
  }

  bottom() {
    return CColumn(children: [
      CRow(children: [
        Obx(() =>
            Checkbox(value: c.type == 1, onChanged: (value) => c.type = 1)),
        const Text('위치도'),
        Obx(() =>
            Checkbox(value: c.type == 2, onChanged: (value) => c.type = 2)),
        const Text('전경'),
        Obx(() =>
            Checkbox(value: c.type == 3, onChanged: (value) => c.type = 3)),
        const Text('부위별'),
        Obx(() =>
            Checkbox(value: c.type == 10, onChanged: (value) => c.type = 10)),
        const Text('주변공사'),
        Obx(() =>
            Checkbox(value: c.type == 11, onChanged: (value) => c.type = 11)),
        const Text('동입구'),
      ]),
      CRow(children: [
        Obx(() =>
            Checkbox(value: c.type == 20, onChanged: (value) => c.type = 20)),
        const Text('지상층 벽체 해머'),
        Obx(() =>
            Checkbox(value: c.type == 21, onChanged: (value) => c.type = 21)),
        const Text('지하층 벽체 해머'),
        Obx(() =>
            Checkbox(value: c.type == 22, onChanged: (value) => c.type = 22)),
        const Text('지상층 슬래브 해머'),
        Obx(() =>
            Checkbox(value: c.type == 23, onChanged: (value) => c.type = 23)),
        const Text('지하층 슬래브 해머'),
        Obx(() =>
            Checkbox(value: c.type == 24, onChanged: (value) => c.type = 24)),
        const Text('지상층 벽체 탄산화'),
        Obx(() =>
            Checkbox(value: c.type == 25, onChanged: (value) => c.type = 25)),
        const Text('균열 팁 측정'),
        const SizedBox(width: 20),
        Expanded(child: TextField(controller: textEditingController)),
        const SizedBox(width: 20),
        ElevatedButton(
            onPressed: () => getImage(ImageSource.camera),
            child: const Text('카메라')),
        const SizedBox(width: 10),
        ElevatedButton(
            onPressed: () => getImage(ImageSource.gallery),
            child: const Text('갤러리')),
      ]),
    ]);
  }

  clickSave(context) async {}

  final picker = ImagePicker();

  Future getImage(ImageSource imageSource) async {
    final image = await picker.pickImage(source: imageSource);

    if (image == null) {
      return;
    }

    var path = image.path;
    var item = Periodicimage();
    item.type = c.type;
    if (c.type == 3 || c.type == 10) {
      item.name = textEditingController.text;
    }
    item.offlinefilename = path;
    c.images.add(item);
    c.saveImage();

    blueprintController.modified = true;
    blueprintController.addModifiedImageType(c.type);
    blueprintController.setLastImagePath(c.type, path);
  }

  Widget _tab(Periodicimage item, int index, context) {
    final typeNames = <int, String>{
      1: '위치도',
      2: '전경',
      3: '부위별',
      10: '주변공사',
      11: '동입구',
      20: '지상층 벽체 해머',
      21: '지하층 벽체 해머',
      22: '지상층 슬래브 해머',
      23: '지하층 슬래브 해머',
      24: '지상층 벽체 탄산화',
      25: '균열 팁 측정',
    };
    final type = typeNames[item.type] ?? '기타';

    return InkWell(
        onTap: () {
          showGeneralDialog(
              barrierDismissible: false,
              context: context,
              pageBuilder: (popContext, __, ___) {
                return Scaffold(
                    body: InkWell(
                  onTap: () {
                    navigator!.pop(popContext);
                  },
                  child: Stack(children: [
                    SizedBox(
                        width: double.infinity,
                        height: double.infinity,
                        child: Image.file(File(item.offlinefilename))),
                    Positioned(
                      top: 30.0,
                      right: 10.0,
                      child: CRow(children: [
                        IconButton(
                            icon: const Icon(CupertinoIcons.trash, size: 30.0),
                            onPressed: () {
                              showDialog<void>(
                                context: context,
                                builder: (context2) {
                                  return AlertDialog(
                                    title: const Text('이미지 삭제'),
                                    backgroundColor: Colors.white,
                                    content: const Text('이미지를 삭제하시겠습니까'),
                                    actions: <Widget>[
                                      ElevatedButton(
                                        child: const Text('취소'),
                                        onPressed: () {
                                          navigator!.pop(context2);
                                        },
                                      ),
                                      ElevatedButton(
                                        child: const Text('삭제'),
                                        onPressed: () {
                                          blueprintController.modified = true;

                                          c.removeImage(index);

                                          navigator!.pop(context2);
                                          navigator!.pop(popContext);
                                        },
                                      )
                                    ],
                                  );
                                },
                              );
                            }),
                        IconButton(
                            icon: const Icon(Icons.close, size: 30.0),
                            onPressed: () {
                              navigator!.pop(popContext);
                            }),
                      ]),
                    ),
                  ]),
                ));
              });
        },
        child: CColumn(children: [
          Container(
            decoration: BoxDecoration(
                border: Border.all(width: 1, color: Colors.grey[700]!)),
            padding: const EdgeInsets.all(10.0),
            margin: const EdgeInsets.only(right: 10.0),
            width: 300,
            child: Column(children: [
              SizedBox(
                  height: 200, child: Image.file(File(item.offlinefilename))),
              Row(children: [
                Container(
                    decoration: BoxDecoration(
                      borderRadius: BorderRadius.circular(10),
                      color: Colors.grey[700],
                    ),
                    margin: const EdgeInsets.only(top: 5),
                    padding: const EdgeInsets.all(5),
                    child: Text(type,
                        style: const TextStyle(color: Colors.white))),
                const SizedBox(width: 10),
                Expanded(
                    child: Text(
                  item.name,
                  overflow: TextOverflow.ellipsis,
                  maxLines: 1,
                  softWrap: false,
                )),
              ]),
            ]),
          ),
        ]));
  }

  List<Widget> _makeTabs(context) {
    List<Widget> tabs = <Widget>[];

    for (var i = 0; i < c.images.length; i++) {
      var item = c.images[i];
      tabs.add(_tab(item, i, context));
    }

    return tabs;
  }
}
