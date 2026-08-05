import 'dart:convert';
import 'dart:io';

import 'package:common_control/common_control.dart';
import 'package:flutter/cupertino.dart';
import 'package:fluttertoast/fluttertoast.dart';
import 'package:image_picker/image_picker.dart';
import 'package:localstorage/localstorage.dart';
import 'package:percent_indicator/percent_indicator.dart';
import 'package:periodic/components/default_app_bar.dart';
import 'package:periodic/components/painter/painter_controller.dart';
import 'package:periodic/controllers/auth_controller.dart';
import 'package:periodic/controllers/blueprint_controller.dart';
import 'package:periodic/controllers/image_controller.dart';
import 'package:periodic/models/blueprint.dart';
import 'package:periodic/models/periodic.dart';
import 'package:periodic/models/periodicimage.dart';
import 'package:periodic/models/periodicother.dart';
import 'package:periodic/models/upload.dart';

class BlueprintScreen extends CWidget {
  BlueprintScreen({super.key});

  final AuthController authController = Get.find<AuthController>();
  final c = Get.find<BlueprintController>();
  final picker = ImagePicker();
  final TextEditingController buwibyeolNameController = TextEditingController();

  endProcess() async {
    final LocalStorage storageLogin = LocalStorage('login.json');
    await storageLogin.ready;
    await storageLogin.deleteItem('periodic');

    authController.periodic = Periodic();
    authController.title = '';
  }

  clickBack(context) async {
    // 데이터를 받는 중이거나 수정사항이 없으면 바로 종료
    if (c.loading == false || c.modified == false) {
      await endProcess();
      Get.back();
      return true;
    }

    final ret = await showDialog<bool>(
      context: context,
      builder: (context2) {
        return AlertDialog(
          title: const Text('데이터 전송'),
          backgroundColor: Colors.white,
          content: const Text(
              '작업내역이 서버로 전송되지 않았습니다.\n전송 없이 종료하시겠습니까.\n저장없이 종료 선택시 작업한 내역이 모두 삭제됩니다'),
          actions: <Widget>[
            ElevatedButton(
              child: const Text('닫기'),
              onPressed: () {
                Navigator.pop(context2, false);
              },
            ),
            ElevatedButton(
              child: const Text('저장없이 종료'),
              onPressed: () async {
                Navigator.pop(context2, false);
                // 임시로 코드 입력 비활성화
                //showExitCodeDialog(context);
                await endProcess();
                Get.back();
              },
            )
          ],
        );
      },
    );

    if (ret == true) {
      await endProcess();
      Get.back();
    }

    return ret ?? false;
  }

  showExitCodeDialog(context) {
    final TextEditingController codeController = TextEditingController();

    showDialog<void>(
      context: context,
      builder: (context3) {
        return AlertDialog(
          title: const Text('종료 코드 입력'),
          backgroundColor: Colors.white,
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              const Text('저장없이 종료하려면 \'1234\' 코드를 입력하세요'),
              const SizedBox(height: 20),
              TextField(
                controller: codeController,
                autofocus: true,
                decoration: const InputDecoration(
                  border: OutlineInputBorder(),
                  labelText: '코드',
                  filled: true,
                  fillColor: Colors.white,
                ),
              ),
            ],
          ),
          actions: <Widget>[
            ElevatedButton(
              child: const Text('취소'),
              onPressed: () {
                Navigator.pop(context3);
              },
            ),
            ElevatedButton(
              child: const Text('확인'),
              onPressed: () {
                if (codeController.text == '1234') {
                  Navigator.pop(context3);
                  endProcess();
                  Get.back();
                } else {
                  Fluttertoast.showToast(
                      msg: '코드가 올바르지 않습니다',
                      toastLength: Toast.LENGTH_SHORT,
                      gravity: ToastGravity.CENTER,
                      timeInSecForIosWeb: 1,
                      backgroundColor: Colors.red,
                      textColor: Colors.white,
                      fontSize: 16.0);
                }
              },
            )
          ],
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return WillPopScope(
        onWillPop: () async {
          return await clickBack(context);
        },
        child: CScaffold(
          appBar: DefaultAppBar(
            leading: IconButton(
                icon: const Icon(Icons.arrow_back_ios),
                onPressed: () async {
                  await clickBack(context);
                }),
            actions: [
              IconButton(
                icon: const Icon(Icons.apartment),
                color: Colors.black,
                onPressed: () {
                  if (c.loading == false) {
                    return;
                  }

                  Get.toNamed('/other');
                },
              ),
              IconButton(
                icon: const Icon(Icons.photo),
                color: Colors.black,
                onPressed: () {
                  if (c.loading == false) {
                    return;
                  }

                  Get.toNamed('/image');
                },
              ),
            ],
          ),
          backgroundColor: Colors.white,
          body: Obx(() => body(context)),
        ));
  }

  body(context) {
    if (c.loading == false) {
      return Column(crossAxisAlignment: CrossAxisAlignment.center, children: [
        Expanded(child: Container()),
        CColumn(
          height: 200,
          padding: const EdgeInsets.only(left: 20, right: 20),
          crossAxisAlignment: CrossAxisAlignment.center,
          children: [
            const SizedBox(height: 10),
            LinearPercentIndicator(
              animation: true,
              lineHeight: 20.0,
              animationDuration: 10,
              percent: c.percent,
              center: Text('${(c.percent * 100.0).toInt()}%',
                  style: const TextStyle(color: Colors.white)),
              barRadius: const Radius.circular(10),
              progressColor: Colors.green,
            ),
            const SizedBox(height: 50),
            const Text('도면 데이터를 전송받는 중입니다', style: TextStyle(fontSize: 16)),
            const SizedBox(height: 10),
            const Text('시간이 소요되니 잠시 기다려 주세요', style: TextStyle(fontSize: 16)),
          ],
        ),
        Expanded(child: Container()),
      ]);
    }

    // collapsed 상태에 따라 항목 필터링
    var visibleItems = <Blueprint>[];
    for (var i = 0; i < c.items.length; i++) {
      var item = c.items[i];
      visibleItems.add(item);

      // level 1이고 collapsed 상태면 하위 항목 건너뛰기
      if (item.level == 1 && item.collapsed) {
        while (i + 1 < c.items.length && c.items[i + 1].level > 1) {
          i++;
        }
      }
    }

    var items = visibleItems.map((item) {
      // 하위 항목 존재 여부 확인 (원본 c.items에서 확인)
      bool hasChildren = false;
      if (item.level == 1) {
        int originalIndex = c.items.indexOf(item);
        if (originalIndex >= 0 && originalIndex < c.items.length - 1) {
          hasChildren = c.items[originalIndex + 1].level > 1;
        }
      }

      return renderItem(item, hasChildren);
    }).toList();

    return GestureDetector(
      onTap: () => FocusScope.of(context).unfocus(),
      behavior: HitTestBehavior.translucent,
      child: CFixedBottom(
          padding: const EdgeInsets.all(10),
          bottom: ElevatedButton(
            style: ElevatedButton.styleFrom(
              minimumSize: const Size.fromHeight(60), // NEW
            ),
            onPressed: () => clickSend(context),
            child: const Text(
              '전송',
              style: TextStyle(fontSize: 20),
            ),
          ),
          children: items),
    );
  }

  Widget renderItem(Blueprint item, bool hasChildren) {
    // 동적으로 수정 여부 확인
    bool isModified = item.extra["modified"] != null;

    // '공중이 이용하는 부위' 관련 항목의 수정 여부 확인
    if (item.extra['isOther'] == true) {
      if (item.level == 1) {
        // 부모 항목: 수정 여부 표시 안함
        isModified = false;
      } else if (item.level == 2 && item.extra['tab'] != null) {
        // 하위 항목: 해당 탭만 확인
        isModified = c.isTabModified(item.extra['tab'] as int);
      }
    }

    // '사진 자료' 관련 항목의 수정 여부 확인
    if (item.extra['isImage'] == true) {
      if (item.level == 1 && hasChildren) {
        // 부모 항목(하위 항목 있음): 수정 여부 표시 안함
        isModified = false;
      } else if (item.extra['imageType'] != null) {
        // 하위 항목 또는 단일 항목(동입구): 해당 타입 확인
        isModified = c.isImageTypeModified(item.extra['imageType'] as int);
      }
    }

    return CContainer(
      border: Border.all(color: Colors.black),
      backgroundColor:
          isModified ? const Color.fromRGBO(255, 200, 200, 1.0) : Colors.white,
      width: double.infinity,
      margin: EdgeInsets.only(bottom: 10, left: (item.level - 1) * 50),
      child: Row(
        children: [
          if (item.level == 1 && hasChildren) ...[
            const SizedBox(width: 10),
            Icon(
              item.collapsed ? Icons.chevron_right : Icons.expand_more,
              size: 30,
            ),
          ],
          Expanded(
            child: CText(
              item.name,
              margin: item.level == 1 && hasChildren
                  ? const EdgeInsets.fromLTRB(5, 20, 20, 20)
                  : const EdgeInsets.all(20),
            ),
          ),
          // '사진 자료' 하위 항목 및 동입구(level 1) 항목에 카메라/갤러리 아이콘 추가
          if (item.extra['isImage'] == true &&
              item.extra['imageType'] != null &&
              (item.level == 2 || (item.level == 1 && !hasChildren))) ...[
            if (item.extra['imageType'] == 3) ...[
              SizedBox(
                width: 180,
                height: 36,
                child: TextField(
                  controller: buwibyeolNameController,
                  autofocus: false,
                  decoration: const InputDecoration(
                    hintText: '명칭',
                    hintStyle: TextStyle(color: Color(0xFF757575)),
                    isDense: true,
                    contentPadding:
                        EdgeInsets.symmetric(horizontal: 8, vertical: 8),
                    border: OutlineInputBorder(),
                  ),
                ),
              ),
              const SizedBox(width: 8),
            ],
            IconButton(
              icon: const Icon(Icons.add_a_photo),
              onPressed: () => getImageForType(
                  item.extra['imageType'] as int, ImageSource.camera),
            ),
            IconButton(
              icon: const Icon(Icons.add_photo_alternate),
              onPressed: () => getImageForType(
                  item.extra['imageType'] as int, ImageSource.gallery),
            ),
            const SizedBox(width: 40),
            // 썸네일 영역: 고정 폭. 왼쪽부터 채워짐
            SizedBox(
              width: 220, // 5 * (36 + 4 right margin) = 200, 여유 포함
              height: 40,
              child: Obx(() {
                final paths =
                    c.getImagePathsByType(item.extra['imageType'] as int);
                if (paths.isEmpty) {
                  return const SizedBox.shrink();
                }
                // 최근 5장 (뒤에서부터)
                final recent =
                    paths.length > 5 ? paths.sublist(paths.length - 5) : paths;
                final thumbs = <Widget>[];
                for (var i = 0; i < recent.length; i++) {
                  final path = recent[i];
                  Widget img;
                  if (path.startsWith('http')) {
                    img = Image.network(path,
                        fit: BoxFit.cover,
                        errorBuilder: (_, __, ___) =>
                            const Icon(Icons.broken_image, size: 20));
                  } else {
                    final file = File(path);
                    if (!file.existsSync()) {
                      continue;
                    }
                    img = Image.file(file, fit: BoxFit.cover);
                  }
                  thumbs.add(Padding(
                    padding: const EdgeInsets.only(right: 4),
                    child: InkWell(
                      onTap: () => showPreview(path),
                      child: Container(
                        width: 36,
                        height: 36,
                        decoration: BoxDecoration(
                          border: Border.all(color: Colors.grey[400]!),
                          borderRadius: BorderRadius.circular(4),
                        ),
                        clipBehavior: Clip.antiAlias,
                        child: img,
                      ),
                    ),
                  ));
                }
                return Row(
                  mainAxisSize: MainAxisSize.min,
                  mainAxisAlignment: MainAxisAlignment.start,
                  children: thumbs,
                );
              }),
            ),
            const SizedBox(width: 10),
          ],
        ],
      ),
      onTap: () {
        // level 1이고 하위 항목이 있으면 토글만 수행
        if (item.level == 1 && hasChildren) {
          final index = c.items.indexOf(item);
          c.toggleCollapse(index);
          return;
        }

        // '공중이 이용하는 부위' 하위 항목 처리
        if (item.extra['isOther'] == true && item.level == 2) {
          final tab = item.extra['tab'];
          Get.toNamed('/other', arguments: {'tab': tab});
          return;
        }

        // '사진 자료' 하위 항목 및 동입구(level 1) 항목 처리
        if (item.extra['isImage'] == true &&
            item.extra['imageType'] != null &&
            (item.level == 2 || (item.level == 1 && !hasChildren))) {
          final imageType = item.extra['imageType'];
          Get.toNamed('/image', arguments: {'type': imageType});
          return;
        }

        // 그 외는 기존 동작
        if (item.upload != 1) {
          return;
        }

        if (item.filename == '') {
          return;
        }

        authController.setTitle(item);
        Get.toNamed('/write', arguments: {'item': item});
      },
    );
  }

  Widget renderItemSend(int index) {
    final item = c.items[index];

    // 하위 항목 존재 여부 확인
    bool hasChildren = false;
    if (item.level == 1) {
      if (index < c.items.length - 1) {
        hasChildren = c.items[index + 1].level > 1;
      }
    }

    // 동적으로 수정 여부 확인 (메인 화면과 동일한 로직)
    bool isModified = item.extra["modified"] != null;

    if (item.extra['isOther'] == true) {
      if (item.level == 1) {
        isModified = false;
      } else if (item.level == 2 && item.extra['tab'] != null) {
        isModified = c.isTabModified(item.extra['tab'] as int);
      }
    }

    if (item.extra['isImage'] == true) {
      if (item.level == 1 && hasChildren) {
        isModified = false;
      } else if (item.extra['imageType'] != null) {
        isModified = c.isImageTypeModified(item.extra['imageType'] as int);
      }
    }

    Color color = Colors.white;
    if (isModified) {
      color = const Color.fromRGBO(255, 200, 200, 1.0);
      if (item.checked == false) {
        color = const Color.fromRGBO(220, 220, 220, 1.0);
      }
    }

    return CContainer(
      border: Border.all(color: Colors.black),
      backgroundColor: color,
      width: double.infinity,
      margin: EdgeInsets.only(bottom: 10, left: (item.level - 1) * 50),
      child: Row(
        children: [
          if (item.level == 1 && hasChildren) ...[
            const SizedBox(width: 10),
            Icon(
              item.collapsed ? Icons.chevron_right : Icons.expand_more,
              size: 30,
            ),
          ],
          Expanded(
            child: CText(
              item.name,
              margin: item.level == 1 && hasChildren
                  ? const EdgeInsets.fromLTRB(5, 20, 20, 20)
                  : const EdgeInsets.all(20),
            ),
          ),
          isModified
              ? Checkbox(
                  onChanged: (value) => c.setCheck(index, value),
                  value: item.checked)
              : const SizedBox(width: 10),
        ],
      ),
      onTap: () {
        if (item.level == 1 && hasChildren) {
          c.toggleCollapse(index);
        }
      },
    );
  }

  clickSend(context) async {
    var find = false;
    for (var i = 0; i < c.items.length; i++) {
      final item = c.items[i];

      if (item.extra['modified'] != null) {
        find = true;
        break;
      }
    }

    if (c.modified == true) {
      find = true;
    }

    if (find == false) {
      showDialog<void>(
        context: context,
        builder: (context2) {
          return AlertDialog(
            title: const Text('데이터 전송'),
            backgroundColor: Colors.white,
            content: const Text('작업 내역이 없습니다',
                style: TextStyle(color: Colors.red, fontSize: 20)),
            actions: <Widget>[
              ElevatedButton(
                child: const Text('닫기'),
                onPressed: () {
                  navigator!.pop(context2);
                },
              ),
            ],
          );
        },
      );

      return;
    }

    c.setCheckAll();

    final h = Get.height;
    showDialog<void>(
      context: context,
      builder: (context2) {
        return AlertDialog(
          title: const Text('데이터 전송'),
          backgroundColor: Colors.white,
          content: CColumn(width: 800, height: h - 320, children: [
            CContainer(
                height: h - 400,
                child: SingleChildScrollView(child: Obx(() {
                  // collapsed 상태에 따라 항목 필터링 (메인 화면과 동일)
                  final visibleIndexes = <int>[];
                  for (var i = 0; i < c.items.length; i++) {
                    final item = c.items[i];
                    visibleIndexes.add(i);
                    if (item.level == 1 && item.collapsed) {
                      while (
                          i + 1 < c.items.length && c.items[i + 1].level > 1) {
                        i++;
                      }
                    }
                  }
                  return CColumn(
                      children: visibleIndexes
                          .map((idx) => renderItemSend(idx))
                          .toList());
                }))),
            const SizedBox(height: 40),
            const Text('온라인 상태에서만 전송이 가능합니다. 데이터를 전송하시겠습니까.',
                style: TextStyle(color: Colors.red, fontSize: 20)),
          ]),
          actions: <Widget>[
            ElevatedButton(
              child: const Text('취소'),
              onPressed: () {
                navigator!.pop(context2);
              },
            ),
            ElevatedButton(
              child: const Text('데이터 전송'),
              onPressed: () {
                c.cancel = true;
                navigator!.pop(context2);

                sendData(context);
              },
            )
          ],
        );
      },
    );
  }

  sendData(context) {
    Future.microtask(() => sendDataProcess());

    showGeneralDialog(
        barrierDismissible: false,
        context: context,
        pageBuilder: (popContext, __, ___) {
          return Obx(() => AlertDialog(
                title: const Text('데이터 전송'),
                backgroundColor: Colors.white,
                content: SizedBox(
                  width: 950,
                  height: 70,
                  child: CColumn(
                      crossAxisAlignment: CrossAxisAlignment.center,
                      children: [
                        const SizedBox(height: 10),
                        LinearPercentIndicator(
                          animation: true,
                          lineHeight: 20.0,
                          animationDuration: 10,
                          percent: c.percent,
                          center: Text('${(c.percent * 100.0).toInt()}%',
                              style: const TextStyle(color: Colors.white)),
                          barRadius: const Radius.circular(10),
                          progressColor: Colors.green,
                        ),
                        const SizedBox(height: 10),
                        c.sendError == true
                            ? const Text('전송중 오류가 발생했습니다',
                                style: TextStyle(color: Colors.red))
                            : c.percent == 1.0
                                ? const Text('전송이 완료되었습니다')
                                : Container()
                      ]),
                ),
                actions: [
                  c.percent != 1.0
                      ? ElevatedButton(
                          onPressed: () {
                            if (c.sendError == true) {
                              navigator!.pop(popContext);
                              return;
                            }

                            showDialog<void>(
                              context: context,
                              builder: (context2) {
                                return AlertDialog(
                                  title: const Text('작업 취소'),
                                  backgroundColor: Colors.white,
                                  content: const Text('작업을 취소하시겠습니까'),
                                  actions: <Widget>[
                                    ElevatedButton(
                                      child: const Text('닫기'),
                                      onPressed: () {
                                        navigator!.pop(context2);
                                      },
                                    ),
                                    ElevatedButton(
                                      child: const Text('작업 취소'),
                                      onPressed: () {
                                        c.cancel = true;
                                        navigator!.pop(context2);
                                        navigator!.pop(popContext);

                                        Fluttertoast.showToast(
                                            msg: '취소되었습니다',
                                            toastLength: Toast.LENGTH_SHORT,
                                            gravity: ToastGravity.CENTER,
                                            timeInSecForIosWeb: 1,
                                            backgroundColor: Colors.red,
                                            textColor: Colors.white,
                                            fontSize: 16.0);
                                      },
                                    )
                                  ],
                                );
                              },
                            );
                          },
                          child: const Text('취소'))
                      : ElevatedButton(
                          onPressed: () => clickSendFinish(),
                          child: const Text('닫기'))
                ],
              ));
        });
  }

  sendDataProcess() async {
    c.cancel = false;
    c.percent = 0.0;
    c.sendError = false;

    final LocalStorage storage = LocalStorage('periodic.json');

    var total = 1;
    var current = 0;

    List<dynamic> datas = [];

    for (var i = 0; i < c.items.length; i++) {
      if (c.items[i].checked == false) {
        continue;
      }

      final blueprint = c.items[i].id;

      await storage.ready;
      final data = await storage.getItem('data_$blueprint');
      final save = await storage.getItem('save_$blueprint');

      if (save == null || save == '') {
        continue;
      }

      if (data == null || data == '') {
        continue;
      }

      final item = json.decode(data);
      datas.add(item);

      final pointsPtr = item['points'];
      if (pointsPtr == null) {
        continue;
      }

      List<dynamic> points = pointsPtr as List<dynamic>;
      for (var j = 0; j < points.length; j++) {
        Point point = Point.fromJson(points[j]);
        total += point.images.length;
      }
    }

    for (var i = 0; i < c.periodicothers.length; i++) {
      Periodicother item = c.periodicothers[i];

      if (item.offlinefilename == '') {
        continue;
      }

      var filenames = item.offlinefilename.split(',');
      total += filenames.length;
    }

    storage.ready;
    final str = await storage.getItem('periodicimages');

    List<Periodicimage> images = [];
    if (str != null && str != '') {
      images = json
          .decode(str)
          .map<Periodicimage>((json) => Periodicimage.fromJson(json))
          .toList();

      total += images.length;

      for (var i = 0; i < images.length; i++) {
        Periodicimage image = images[i];
        if (c.cancel == true) {
          return;
        }

        var ret = '';
        for (var k = 0; k < 5; k++) {
          if (File(image.offlinefilename).existsSync() != true) {
            break;
          }

          ret = await UploadManager.image(image.offlinefilename);

          if (ret != '') {
            break;
          }

          if (k == 4) {
            c.sendError = true;
            break;
          }

          for (var j = 0; j < k + 1; j++) {
            if (c.cancel == true) {
              return;
            }
            sleep(const Duration(seconds: 1));
          }
        }

        images[i].filename = ret;

        current++;
        c.percent = current / total;
      }
    }

    for (var i = 0; i < datas.length; i++) {
      List<dynamic> points = datas[i]['points'];
      for (var j = 0; j < points.length; j++) {
        Point point = Point.fromJson(points[j]);

        datas[i]['points'][j]['onlineimages'] = [];
        for (var k = 0; k < point.images.length; k++) {
          final image = point.images[k];

          var ret = '';
          for (var l = 0; l < 5; l++) {
            if (File(image).existsSync() != true) {
              break;
            }

            ret = await UploadManager.image(image);

            if (ret != '') {
              break;
            }

            if (l == 4) {
              c.sendError = true;
              break;
            }

            for (var m = 0; m < l + 1; m++) {
              if (c.cancel == true) {
                return;
              }
              sleep(const Duration(seconds: 1));
            }
          }

          datas[i]['points'][j]['onlineimages'].add(ret);

          current++;
          c.percent = current / total;
        }
      }
    }

    for (var i = 0; i < c.periodicothers.length; i++) {
      Periodicother item = c.periodicothers[i];

      if (item.offlinefilename == '') {
        c.periodicothers[i].filename = '';
        continue;
      }

      var filenames = item.offlinefilename.split(',');

      List<String> onlinefilenames = [];

      for (var j = 0; j < filenames.length; j++) {
        var image = filenames[j];
        if (File(image).existsSync() != true) {
          continue;
        }

        final ret = await UploadManager.image(image);
        if (ret != '') {
          onlinefilenames.add(ret);
        }

        current++;
        c.percent = current / total;
      }

      c.periodicothers[i].filename = onlinefilenames.join(',');
    }

    Map<String, dynamic> ret = {
      'user': authController.user.id,
      'id': c.id,
      'datas': datas,
      'images': images,
      'periodicothers': c.periodicothers,
    };

    if (c.cancel == true) {
      return;
    }

    for (var i = 0; i < 5; i++) {
      if (c.cancel == true) {
        return;
      }

      final res = await Http.post('/api/periodic/upload', ret);
      if (res['code'] == 'ok') {
        c.percent = 1.0;
        if (authController.autoclose == true) {
          clickSendFinish();
        }
        return;
      }

      for (var j = 0; j < i + 1; j++) {
        if (c.cancel == true) {
          return;
        }
        sleep(const Duration(seconds: 1));
      }
    }

    c.sendError = true;
  }

  clickSendFinish() async {
    if (c.sendError == true) {
      Get.back();
      c.sendError = false;
      return;
    }

    // 전송 완료 후 change 플래그 초기화 → 재진입 시 붉은 잔상 방지
    for (var i = 0; i < c.periodicothers.length; i++) {
      c.periodicothers[i].change = 0;
    }
    final LocalStorage storageBlueprint = LocalStorage('blueprints.json');
    await storageBlueprint.ready;
    await storageBlueprint.setItem(
        'periodicothers', json.encode(c.periodicothers));

    endProcess();

    c.modified = false;
    Get.back();
    Get.back();
  }

  void showPreview(String path) {
    final ctx = Get.context;
    if (ctx == null) return;
    showGeneralDialog(
      barrierDismissible: false,
      context: ctx,
      pageBuilder: (popContext, __, ___) {
        Widget image;
        if (path.startsWith('http')) {
          image = Image.network(path);
        } else {
          image = Image.file(File(path));
        }
        return Scaffold(
          body: InkWell(
            onTap: () => Navigator.of(popContext).pop(),
            child: Stack(children: [
              SizedBox(
                width: double.infinity,
                height: double.infinity,
                child: image,
              ),
              Positioned(
                top: 30.0,
                right: 10.0,
                child: Row(children: [
                  IconButton(
                    icon: const Icon(CupertinoIcons.trash, size: 30.0),
                    onPressed: () {
                      showDialog<void>(
                        context: popContext,
                        builder: (dialogContext) {
                          return AlertDialog(
                            title: const Text('이미지 삭제'),
                            backgroundColor: Colors.white,
                            content: const Text('이미지를 삭제하시겠습니까'),
                            actions: <Widget>[
                              ElevatedButton(
                                child: const Text('취소'),
                                onPressed: () {
                                  Navigator.of(dialogContext).pop();
                                },
                              ),
                              ElevatedButton(
                                child: const Text('삭제'),
                                onPressed: () async {
                                  Navigator.of(dialogContext).pop();
                                  Navigator.of(popContext).pop();

                                  final removed =
                                      await _deleteImageByPath(path);

                                  Fluttertoast.showToast(
                                      msg: removed
                                          ? '사진이 삭제되었습니다'
                                          : '사진을 찾을 수 없습니다',
                                      gravity: ToastGravity.CENTER,
                                      backgroundColor: Colors.grey[700],
                                      textColor: Colors.white);
                                },
                              ),
                            ],
                          );
                        },
                      );
                    },
                  ),
                  IconButton(
                    icon: const Icon(Icons.close, size: 30.0),
                    onPressed: () => Navigator.of(popContext).pop(),
                  ),
                ]),
              ),
            ]),
          ),
        );
      },
    );
  }

  Future<bool> _deleteImageByPath(String path) async {
    try {
      final LocalStorage storage = LocalStorage('periodic.json');
      await storage.ready;
      final str = await storage.getItem('periodicimages');
      if (str == null || str == '') return false;

      final List<Periodicimage> images = (json.decode(str) as List)
          .map<Periodicimage>((j) => Periodicimage.fromJson(j))
          .toList();

      final base = path.split('/').last.split('\\').last;
      final idx = images.indexWhere((e) {
        if (e.offlinefilename == path || e.filename == path) return true;
        final offBase = e.offlinefilename.split('/').last.split('\\').last;
        final fnBase = e.filename.split('/').last.split('\\').last;
        return base.isNotEmpty && (offBase == base || fnBase == base);
      });

      if (idx < 0) return false;

      // 로컬 파일도 실제 삭제
      final offlinePath = images[idx].offlinefilename;
      if (offlinePath.isNotEmpty && !offlinePath.startsWith('http')) {
        try {
          final f = File(offlinePath);
          if (f.existsSync()) f.deleteSync();
        } catch (_) {}
      }

      images.removeAt(idx);
      await storage.setItem('periodicimages', json.encode(images));

      // ImageController 동기화
      try {
        final imageController = Get.find<ImageController>();
        imageController.images = images;
      } catch (_) {}

      // 썸네일 캐시 갱신
      await c.refreshLastImagePaths();
      c.modified = true;
      c.modifiedImage = true;

      return true;
    } catch (_) {
      return false;
    }
  }

  Future getImageForType(int imageType, ImageSource imageSource) async {
    final image = await picker.pickImage(source: imageSource);

    if (image == null) {
      return;
    }

    var path = image.path;
    var item = Periodicimage();
    item.type = imageType;
    if (imageType == 3) {
      item.name = buwibyeolNameController.text;
      buwibyeolNameController.clear();
    }
    item.offlinefilename = path;

    final LocalStorage storage = LocalStorage('periodic.json');
    await storage.ready;
    final str = await storage.getItem('periodicimages');

    List<Periodicimage> images = [];
    if (str != null && str != '') {
      images = json
          .decode(str)
          .map<Periodicimage>((json) => Periodicimage.fromJson(json))
          .toList();
    }

    images.add(item);
    final newStr = json.encode(images);
    await storage.setItem('periodicimages', newStr);

    // ImageController 업데이트
    try {
      final imageController = Get.find<ImageController>();
      imageController.onInit(); // 이미지 목록 다시 로드
    } catch (e) {
      // ImageController가 없으면 무시
    }

    c.modified = true;
    c.modifiedImage = true;
    c.addModifiedImageType(imageType);
    c.setLastImagePath(imageType, path);

    Fluttertoast.showToast(
        msg: '사진이 추가되었습니다',
        toastLength: Toast.LENGTH_SHORT,
        gravity: ToastGravity.CENTER,
        timeInSecForIosWeb: 1,
        backgroundColor: Colors.grey[700],
        textColor: Colors.white,
        fontSize: 16.0);
  }
}
