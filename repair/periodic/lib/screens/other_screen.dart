import 'dart:io';

import 'package:common_control/common_control.dart';
import 'package:flutter/cupertino.dart';
import 'package:fluttertoast/fluttertoast.dart';
import 'package:image_picker/image_picker.dart';
import 'package:periodic/components/default_app_bar.dart';
import 'package:periodic/controllers/auth_controller.dart';
import 'package:periodic/controllers/other_controller.dart';
import 'package:periodic/models/periodicother.dart';

class OtherScreen extends CWidget {
  OtherScreen({super.key});

  final AuthController authController = Get.find<AuthController>();
  final c = Get.find<OtherController>();

  clickBack(context) {
    if (c.modified == false) {
      Get.back();
      return true;
    }

    final ret = showDialog<void>(
      context: context,
      builder: (context2) {
        return AlertDialog(
          title: const Text('데이터 저장'),
          backgroundColor: Colors.white,
          content: const Text(
              '작업내역이 저장되지 않았습니다.\n저장없이 종료하시겠습니까?\n저장없이 종료 선택시 작업한 내역이 모두 삭제됩니다'),
          actions: <Widget>[
            ElevatedButton(
              child: const Text('닫기'),
              onPressed: () {
                Navigator.pop(context2, false);
              },
            ),
            ElevatedButton(
              child: const Text('저장없이 종료'),
              onPressed: () {
                Navigator.pop(context2, true);
                Get.back();
              },
            )
          ],
        );
      },
    );

    return ret;
  }

  @override
  Widget build(BuildContext context) {
    // arguments로 tab이 전달되면 설정
    final arguments = Get.arguments;
    if (arguments != null && arguments['tab'] != null) {
      c.tab = arguments['tab'];
    }

    return WillPopScope(
        onWillPop: () {
          return clickBack(context);
        },
        child: CScaffold(
          appBar: DefaultAppBar(
            leading: IconButton(
                icon: const Icon(Icons.arrow_back_ios),
                onPressed: () {
                  clickBack(context);
                }),
            actions: [
              Obx(() => authController.autosave == false
                  ? IconButton(
                      icon: const Icon(Icons.save),
                      color: Colors.black,
                      onPressed: () {
                        c.save();

                        Fluttertoast.showToast(
                            msg: '저장되었습니다',
                            toastLength: Toast.LENGTH_SHORT,
                            gravity: ToastGravity.CENTER,
                            timeInSecForIosWeb: 1,
                            backgroundColor: Colors.grey[700],
                            textColor: Colors.white,
                            fontSize: 16.0);
                      })
                  : Container()),
            ],
          ),
          backgroundColor: Colors.white,
          body: Obx(() => body(context)),
        ));
  }

  body(context) {
    // final textStyle = Style(
    //   margin: const EdgeInsets.all(10),
    //   padding: const EdgeInsets.all(10),
    //   textStyle: const TextStyle(fontSize: 20),
    // );

    // final textStyleSelected = Style(
    //   margin: const EdgeInsets.all(10),
    //   padding: const EdgeInsets.all(10),
    //   textStyle: const TextStyle(fontSize: 20, fontWeight: FontWeight.w700),
    //   decoration: const BoxDecoration(
    //     border: Border(
    //       bottom: BorderSide(width: 1.5, color: Colors.blue),
    //     ),
    //   ),
    // );

    return CColumn(children: [
      // commented out tab navigation omitted
      if (c.tab == 16)
        Expanded(child: tab16Body(context))
      else
        CScroll(expanded: true, children: [checklist(context), list(context)])
    ]);
  }

  final picker = ImagePicker();

  Future getImage(int index, ImageSource imageSource) async {
    final image = await picker.pickImage(source: imageSource);

    if (image == null) {
      return;
    }

    var path = image.path;

    var offlinefilename = c.periodicothers[index].offlinefilename;
    if (offlinefilename == '') {
      offlinefilename = path;
    } else {
      var filenames = offlinefilename.split(',');
      filenames.add(path);

      offlinefilename = filenames.join(',');
    }

    c.periodicothers[index].offlinefilename = offlinefilename;
    c.periodicothers[index].change = 1;
    c.modified = true;
    c.updatePeriodicothers();
  }

  imageWidget(String filename, int index, int pos, context) {
    return CContainer(
        width: 20,
        height: 20,
        child: Image.file(File(filename)),
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
                        child: Image.file(File(filename))),
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
                                          var offlinefilename = c
                                              .periodicothers[index]
                                              .offlinefilename;
                                          if (offlinefilename != '') {
                                            var filenames =
                                                offlinefilename.split(',');

                                            if (pos < filenames.length) {
                                              filenames.removeAt(pos);
                                            }

                                            offlinefilename =
                                                filenames.join(',');

                                            c.periodicothers[index]
                                                    .offlinefilename =
                                                offlinefilename;
                                            c.periodicothers[index].change = 1;
                                            c.modified = true;
                                            c.updatePeriodicothers();
                                          }

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
        });
  }

  makeRow(index, Periodicother item, context) {
    List<Widget> images = [];

    if (item.offlinefilename != '') {
      final filenames = item.offlinefilename.split(',');
      for (var i = 0; i < filenames.length; i++) {
        var image = filenames[i];

        if (File(image).existsSync() != true) {
          continue;
        }

        images.add(imageWidget(image, index, i, context));
      }
    }

    var txt1 = '양호';
    var txt2 = '보통';
    if (item.type == 2) {
      txt1 = '없음';
      txt2 = '있음';
    }
    return TableRow(children: [
      Container(
        alignment: Alignment.centerLeft,
        color: Colors.white,
        padding: const EdgeInsets.all(10),
        margin: const EdgeInsets.only(top: 5),
        child: Text(item.name, style: const TextStyle(fontSize: 16)),
      ),
      Container(
        alignment: Alignment.center,
        color: Colors.white,
        child: CRow(children: [
          Radio(
            value: 1,
            groupValue: item.result,
            onChanged: (Object? value) {
              c.periodicothers[index].result = 1;
              c.periodicothers[index].change = 1;
              c.modified = true;
              c.updatePeriodicothers();
            },
          ),
          Text(txt1, style: const TextStyle(fontSize: 16)),
          Radio(
            value: 2,
            groupValue: item.result,
            onChanged: (Object? value) {
              c.periodicothers[index].result = 2;
              c.periodicothers[index].change = 1;
              c.modified = true;
              c.updatePeriodicothers();
            },
          ),
          Text(txt2, style: const TextStyle(fontSize: 16)),
        ]),
      ),
      Container(
        alignment: Alignment.center,
        color: Colors.white,
        child: TextField(
            controller: c.statusController[index],
            onChanged: (value) {
              c.periodicothers[index].status = value;
              c.periodicothers[index].change = 1;
              c.modified = true;
              c.updatePeriodicothers();
            },
            decoration: const InputDecoration(
                contentPadding: EdgeInsets.only(left: 10, right: 10),
                border: InputBorder.none)),
      ),
      Container(
        alignment: Alignment.center,
        color: Colors.white,
        child: TextField(
            controller: c.positionController[index],
            onChanged: (value) {
              c.periodicothers[index].position = value;
              c.periodicothers[index].change = 1;
              c.modified = true;
              c.updatePeriodicothers();
            },
            decoration: const InputDecoration(
                contentPadding: EdgeInsets.only(left: 10, right: 10),
                border: InputBorder.none)),
      ),
      Container(
        alignment: Alignment.center,
        padding: const EdgeInsets.all(10),
        child: CRow(children: [
          CContainer(
            child: const Icon(Icons.add_a_photo),
            onTap: () => getImage(index, ImageSource.camera),
          ),
          const SizedBox(width: 10),
          CContainer(
            child: const Icon(Icons.add_photo_alternate),
            onTap: () => getImage(index, ImageSource.gallery),
          ),
          CRow(
            margin: const EdgeInsets.only(left: 20),
            gap: 5,
            children: images,
          )
        ]),
      ),
    ]);
  }

  Widget list(context) {
    if (c.tab == 10 || c.tab == 11 || c.tab == 12 || c.tab == 15) {
      return Container();
    }

    var category = c.tab;

    if (c.tab == 13) {
      category = 2;
    } else if (c.tab == 14) {
      category = 1;
    }
    List<TableRow> items = [];

    TableRow title = TableRow(children: [
      Container(
          alignment: Alignment.center,
          padding: const EdgeInsets.all(10),
          child: const Text('점검내용')),
      Container(
          alignment: Alignment.center,
          padding: const EdgeInsets.all(10),
          child: const Text('점검결과')),
      Container(
          alignment: Alignment.center,
          padding: const EdgeInsets.all(10),
          child: const Text('상태')),
      Container(
          alignment: Alignment.center,
          padding: const EdgeInsets.all(10),
          child: const Text('해당위치')),
      Container(
          alignment: Alignment.center,
          padding: const EdgeInsets.all(10),
          child: const Text('')),
    ]);

    items.add(title);

    for (var i = 0; i < c.periodicothers.length; i++) {
      Periodicother item = c.periodicothers[i];

      if (item.category != category) {
        continue;
      }

      items.add(makeRow(i, item, context));
    }

    return Container(
        padding: const EdgeInsets.all(10),
        child: SingleChildScrollView(
            child: Table(
          columnWidths: const {
            0: FixedColumnWidth(340),
            1: FixedColumnWidth(180),
            2: FixedColumnWidth(200),
            3: FixedColumnWidth(200),
          },
          border: TableBorder.all(color: Colors.black),
          children: items,
        )));
  }

  makeCheckRow(index, Periodicother item, context) {
    final tab = c.tab;
    List<Widget> images = [];

    if (item.offlinefilename != '') {
      final filenames = item.offlinefilename.split(',');
      for (var i = 0; i < filenames.length; i++) {
        var image = filenames[i];

        if (File(image).existsSync() != true) {
          continue;
        }

        images.add(imageWidget(image, index, i, context));
      }
    }

    // var txt1 = '양호';
    // var txt2 = '보통';
    // if (item.type == 2) {
    //   txt1 = '없음';
    //   txt2 = '있음';
    // }

    final names = item.name.split(',');

    List<Widget> widgets = [];

    if (item.type == 1 || item.type == 3) {
      for (var i = 0; i < names.length; i++) {
        widgets.add(Radio(
          value: names[i],
          groupValue: item.status,
          onChanged: (Object? value) {
            c.periodicothers[index].status = names[i];
            c.periodicothers[index].change = 1;
            c.modified = true;
            c.updatePeriodicothers();
          },
        ));

        widgets.add(Text(names[i], style: const TextStyle(fontSize: 16)));
      }
    } else {
      for (var i = 0; i < names.length; i++) {
        widgets.add(Checkbox(
          value: getCheckboxValue(tab, index, i, names),
          onChanged: (Object? value) {
            setCheckboxValue(names[i], tab, index, i, names);
            c.modified = true;
            c.updatePeriodicothers();
          },
        ));

        widgets.add(Text(names[i], style: const TextStyle(fontSize: 16)));
      }
    }

    return TableRow(children: [
      Container(
        alignment: Alignment.centerLeft,
        color: Colors.white,
        padding: const EdgeInsets.all(10),
        margin: const EdgeInsets.only(top: 5),
        child: Text(item.position, style: const TextStyle(fontSize: 16)),
      ),
      Container(
        alignment: Alignment.centerLeft,
        color: Colors.white,
        margin: const EdgeInsets.only(top: 5),
        child: CRow(children: widgets),
      ),
      item.name == "a,b,c,d,e"
          ? Container()
          : Container(
              alignment: Alignment.center,
              padding: const EdgeInsets.all(10),
              child: CRow(children: [
                CContainer(
                  child: const Icon(Icons.add_a_photo),
                  onTap: () => getImage(index, ImageSource.camera),
                ),
                const SizedBox(width: 10),
                CContainer(
                  child: const Icon(Icons.add_photo_alternate),
                  onTap: () => getImage(index, ImageSource.gallery),
                ),
                CRow(
                  margin: const EdgeInsets.only(left: 20),
                  gap: 5,
                  children: images,
                )
              ]),
            ),
    ]);
  }

  Widget checklist(context) {
    if (c.tab < 10) {
      return Container();
    }

    List<TableRow> items = [];

    TableRow title = TableRow(children: [
      Container(
          alignment: Alignment.center,
          padding: const EdgeInsets.all(10),
          child: const Text('점검내용')),
      Container(
          alignment: Alignment.center,
          padding: const EdgeInsets.all(10),
          child: const Text('상태')),
      Container(
          alignment: Alignment.center,
          padding: const EdgeInsets.all(10),
          child: const Text('')),
    ]);

    items.add(title);

    for (var i = 0; i < c.periodicothers.length; i++) {
      Periodicother item = c.periodicothers[i];

      if (item.category != c.tab) {
        continue;
      }

      items.add(makeCheckRow(i, item, context));
    }

    return Container(
        padding: const EdgeInsets.all(10),
        child: SingleChildScrollView(
            child: Table(
          columnWidths: const {
            0: FixedColumnWidth(220),
            1: FixedColumnWidth(750),
          },
          border: TableBorder.all(color: Colors.black),
          children: items,
        )));
  }

  getCheckboxValue(category, index, i, names) {
    final str = names[i];

    final values = c.periodicothers[index].status.split(',');

    for (var i = 0; i < values.length; i++) {
      if (values[i] == str) {
        return true;
      }
    }

    return false;
  }

  setCheckboxValue(name, category, index, i, names) {
    final check = getCheckboxValue(category, index, i, names);

    final values = c.periodicothers[index].status.split(',');

    List<String> newValue = [];

    if (check == true) {
      for (var i = 0; i < values.length; i++) {
        if (values[i] == name) {
          continue;
        }

        newValue.add(values[i]);
      }
    } else {
      values.add(name);

      for (var i = 0; i < names.length; i++) {
        for (var j = 0; j < values.length; j++) {
          if (names[i] == values[j]) {
            newValue.add(values[j]);
          }
        }
      }
    }
    c.periodicothers[index].status = newValue.join(',');
    c.periodicothers[index].change = 1;
  }

  Future getTab16Image(int rowIndex, ImageSource imageSource) async {
    final image = await picker.pickImage(source: imageSource);
    if (image == null) return;

    c.addTab16Image(rowIndex, image.path);
  }

  void showTab16ImagePreview(
      int rowIndex, int imageIndex, String path, BuildContext context) {
    showGeneralDialog(
      barrierDismissible: false,
      context: context,
      pageBuilder: (popContext, __, ___) {
        return Scaffold(
          body: InkWell(
            onTap: () => Navigator.of(popContext).pop(),
            child: Stack(
              children: [
                SizedBox(
                  width: double.infinity,
                  height: double.infinity,
                  child: Image.file(File(path)),
                ),
                Positioned(
                  top: 30.0,
                  right: 10.0,
                  child: Row(
                    children: [
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
                                    onPressed: () =>
                                        Navigator.of(dialogContext).pop(),
                                  ),
                                  ElevatedButton(
                                    child: const Text('삭제'),
                                    onPressed: () {
                                      c.removeTab16Image(rowIndex, imageIndex);
                                      Navigator.of(dialogContext).pop();
                                      Navigator.of(popContext).pop();
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
                    ],
                  ),
                ),
              ],
            ),
          ),
        );
      },
    );
  }

  Widget tab16Body(BuildContext context) {
    const double rowH = 44.0;
    const double col1W = 90.0;
    const double col2W = 90.0;
    const double col3W = 160.0;
    const double col4W = 280.0; // 이미지 열 폭
    const borderSide = BorderSide(color: Colors.black, width: 0.5);
    const deco = BoxDecoration(border: Border.fromBorderSide(borderSide));

    Widget fixedCell(String text, double w, double h) {
      return Container(
        width: w,
        height: h,
        alignment: Alignment.center,
        padding: const EdgeInsets.symmetric(horizontal: 8),
        decoration: deco,
        child: Text(text,
            textAlign: TextAlign.center, style: const TextStyle(fontSize: 16)),
      );
    }

    Widget flexCell(String text, double h) {
      return Container(
        height: h,
        width: double.infinity,
        alignment: Alignment.centerLeft,
        padding: const EdgeInsets.symmetric(horizontal: 8),
        decoration: deco,
        child: Text(text, style: const TextStyle(fontSize: 16)),
      );
    }

    Widget imageCell(int rowIndex, double h) {
      return Obx(() {
        final images = c.getTab16Images(rowIndex);
        return Container(
          width: col4W,
          height: h,
          padding: const EdgeInsets.symmetric(horizontal: 4, vertical: 4),
          decoration: deco,
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.start,
                children: [
                  IconButton(
                    icon: const Icon(Icons.add_a_photo, size: 24),
                    padding: EdgeInsets.zero,
                    constraints: const BoxConstraints(),
                    onPressed: () =>
                        getTab16Image(rowIndex, ImageSource.camera),
                  ),
                  IconButton(
                    icon: const Icon(Icons.add_photo_alternate, size: 24),
                    padding: EdgeInsets.zero,
                    constraints: const BoxConstraints(),
                    onPressed: () =>
                        getTab16Image(rowIndex, ImageSource.gallery),
                  ),
                ],
              ),
              if (images.isNotEmpty)
                Expanded(
                  child: Wrap(
                    spacing: 2,
                    runSpacing: 2,
                    children: images.asMap().entries.map((entry) {
                      final idx = entry.key;
                      final path = entry.value;
                      return InkWell(
                        onTap: () =>
                            showTab16ImagePreview(rowIndex, idx, path, context),
                        child: Container(
                          width: 24,
                          height: 24,
                          decoration: BoxDecoration(
                            border: Border.all(color: Colors.grey),
                          ),
                          child: Image.file(File(path), fit: BoxFit.cover),
                        ),
                      );
                    }).toList(),
                  ),
                ),
            ],
          ),
        );
      });
    }

    return SingleChildScrollView(
      padding: const EdgeInsets.all(10),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // 헤더 행
          Row(children: [
            Container(
              width: 2 * col1W,
              height: rowH,
              alignment: Alignment.center,
              padding: const EdgeInsets.symmetric(horizontal: 8),
              decoration: deco,
              child: const Text('구분',
                  textAlign: TextAlign.center, style: TextStyle(fontSize: 14)),
            ),
            Container(
              width: col3W,
              height: rowH,
              alignment: Alignment.center,
              padding: const EdgeInsets.symmetric(horizontal: 8),
              decoration: deco,
              child: const Text('조사항목',
                  textAlign: TextAlign.center, style: TextStyle(fontSize: 14)),
            ),
            Expanded(
              child: Container(
                height: rowH,
                alignment: Alignment.center,
                padding: const EdgeInsets.symmetric(horizontal: 8),
                decoration: deco,
                child: const Text('내용',
                    textAlign: TextAlign.center,
                    style: TextStyle(fontSize: 14)),
              ),
            ),
            Container(
              width: col4W,
              height: rowH,
              alignment: Alignment.center,
              padding: const EdgeInsets.symmetric(horizontal: 8),
              decoration: deco,
              child: const Text('이미지',
                  textAlign: TextAlign.center, style: TextStyle(fontSize: 14)),
            ),
          ]),
          // 부착물 등 행
          Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              fixedCell('부착물 등', col1W, 7 * rowH),
              // 소분류 열
              Column(children: [
                fixedCell('정착부', col2W, 3.5 * rowH),
                fixedCell('연결부', col2W, 2.5 * rowH),
                fixedCell('보강부', col2W, rowH),
              ]),
              // 조사항목 열
              Column(children: [
                Container(
                  width: col3W,
                  height: 1.5 * rowH,
                  padding:
                      const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: deco,
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: const [
                      Text('앵커 정착부',
                          textAlign: TextAlign.center,
                          style: TextStyle(fontSize: 16)),
                      SizedBox(height: 4),
                      Text('브라켓 정착부',
                          textAlign: TextAlign.center,
                          style: TextStyle(fontSize: 16)),
                    ],
                  ),
                ),
                fixedCell('용접 정착부', col3W, rowH),
                fixedCell('매립 정착부', col3W, rowH),
                Container(
                  width: col3W,
                  height: 1.5 * rowH,
                  padding:
                      const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: deco,
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: const [
                      Text('볼트 연결부',
                          textAlign: TextAlign.center,
                          style: TextStyle(fontSize: 16)),
                    ],
                  ),
                ),
                fixedCell('용접 연결부', col3W, rowH),
                fixedCell('와이어 로프', col3W, rowH),
              ]),
              // 내용 열 (Expanded)
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    Container(
                      height: 1.5 * rowH,
                      width: double.infinity,
                      padding: const EdgeInsets.symmetric(
                          horizontal: 8, vertical: 4),
                      decoration: deco,
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: const [
                          Text('정착부 콘크리트의 균열, 박락 여부 확인',
                              style: TextStyle(fontSize: 16)),
                          SizedBox(height: 4),
                          Text('앵커 시공상태, 풀림 및 빠짐, 부식 여부 확인',
                              style: TextStyle(fontSize: 16)),
                        ],
                      ),
                    ),
                    flexCell('용접 면적 적정성, 균열, 부식 등 손상 발생 여부 등 확인', rowH),
                    flexCell('정착 철물 매립 길이, 철물 여장(노출) 길이 등', rowH),
                    Container(
                      height: 1.5 * rowH,
                      width: double.infinity,
                      padding: const EdgeInsets.symmetric(
                          horizontal: 8, vertical: 4),
                      decoration: deco,
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: const [
                          Text('볼트 풀림 및 빠짐, 부재 변형, 부식 등 확인',
                              style: TextStyle(fontSize: 16)),
                          SizedBox(height: 4),
                          Text('볼트의 시공상태(볼트 규격, 설치 간격 등) 확인',
                              style: TextStyle(fontSize: 16)),
                        ],
                      ),
                    ),
                    flexCell('용접 면적 적정성, 균열 발생 여부 등 확인', rowH),
                    flexCell('손상, 부식, 변형, 고정클립 수량 및 상태 등 확인', rowH),
                  ],
                ),
              ),
              // 이미지 열
              Column(children: [
                imageCell(0, 1.5 * rowH),
                imageCell(2, rowH),
                imageCell(3, rowH),
                imageCell(4, 1.5 * rowH),
                imageCell(5, rowH),
                imageCell(6, rowH),
              ]),
            ],
          ),
          // 변위 변형 행
          Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              fixedCell('변위 변형', col1W + col2W, rowH),
              fixedCell('기울기 및 배부름', col3W, rowH),
              Expanded(
                child: flexCell('면외방향 기울기 및 배부름 발생 유무', rowH),
              ),
              imageCell(7, rowH),
            ],
          ),
          const SizedBox(height: 20),
          // 점검내용 입력
          Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Container(
                width: 80,
                padding: const EdgeInsets.only(top: 14, right: 8, left: 8),
                child: const Text('점검내용', style: TextStyle(fontSize: 16)),
              ),
              Expanded(
                child: TextField(
                  controller: c.tab16ContentController,
                  maxLines: 5,
                  decoration: InputDecoration(
                    hintText: '점검내용을 입력하세요',
                    hintStyle:
                        const TextStyle(color: Color(0xFF757575), fontSize: 14),
                    border: const OutlineInputBorder(),
                    filled: true,
                    fillColor: Colors.white,
                    contentPadding: const EdgeInsets.all(12),
                  ),
                  onChanged: (_) => c.saveTab16Content(),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}
