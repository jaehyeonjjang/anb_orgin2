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

    // 탭 이름 매핑
    String tabName = '';
    if (c.tab == 3) {
      tabName = '부대 점검사항';
    }
    // tab 13, 14는 checklist에서 이름 표시

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
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // 탭 이름 표시
            if (tabName.isNotEmpty)
              Padding(
                padding: const EdgeInsets.only(bottom: 10),
                child: Text(tabName,
                    style: const TextStyle(
                        fontSize: 18, fontWeight: FontWeight.bold)),
              ),
            // 표
            Table(
              columnWidths: const {
                0: FixedColumnWidth(340),
                1: FixedColumnWidth(180),
                2: FixedColumnWidth(200),
                3: FixedColumnWidth(200),
              },
              border: TableBorder.all(color: Colors.black),
              children: items,
            ),
          ],
        ));
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

    var names = item.name.split(',');

    // 체크박스 타입(type == 2)이고 '상태양호'가 없으면 맨 앞에 추가
    if (item.type == 2 && !names.contains('상태양호')) {
      names = ['상태양호', ...names];
    }

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

    // 탭 이름 매핑
    String tabName = '';
    if (c.tab == 10) {
      tabName = '추락방지시설';
    } else if (c.tab == 11) {
      tabName = '도로포장';
    } else if (c.tab == 12) {
      tabName = '도로부 신축 이음부';
    } else if (c.tab == 13) {
      tabName = '환기구 등의 덮개';
    } else if (c.tab == 14) {
      tabName = '외벽 마감재';
    } else if (c.tab == 15) {
      tabName = '강재구조 노후';
    }

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
            child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // 탭 이름 표시
            if (tabName.isNotEmpty)
              Padding(
                padding: const EdgeInsets.only(bottom: 10),
                child: Text(tabName,
                    style: const TextStyle(
                        fontSize: 18, fontWeight: FontWeight.bold)),
              ),
            // 표
            Table(
              columnWidths: const {
                0: FixedColumnWidth(220),
                1: FixedColumnWidth(750),
              },
              border: TableBorder.all(color: Colors.black),
              children: items,
            ),
          ],
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
    // 로딩 중이면 로딩 인디케이터 표시
    return Obx(() {
      if (c.isLoading) {
        return const Center(
          child: CircularProgressIndicator(),
        );
      }

      return _buildTab16Content(context);
    });
  }

  Widget _buildTab16Content(BuildContext context) {
    const borderSide = BorderSide(color: Colors.black, width: 0.5);
    const deco = BoxDecoration(border: Border.fromBorderSide(borderSide));
    const double col1W = 90.0; // 대분류 (정착부/연결부 등)
    const double col2W = 130.0; // 소분류 (앵커 및 브라켓 등)
    const double photoW = 130.0; // 사진 열

    // 라벨 셀 (테두리 + 가운데 정렬)
    Widget labelCell(String text,
        {double? width,
        Alignment alignment = Alignment.center,
        TextAlign align = TextAlign.center,
        FontWeight? weight}) {
      return Container(
        width: width,
        alignment: alignment,
        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 10),
        decoration: deco,
        child: Text(text,
            textAlign: align,
            style: TextStyle(fontSize: 15, fontWeight: weight)),
      );
    }

    // 상태 셀 (체크박스 / 라디오)
    Widget stateCell(int rowIndex, List<String> options,
        {bool isRadio = false}) {
      return Expanded(
        child: Obx(() => Container(
              alignment: Alignment.centerLeft,
              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 8),
              decoration: deco,
              child: Wrap(
                spacing: 12,
                runSpacing: 4,
                crossAxisAlignment: WrapCrossAlignment.center,
                children: options.map((opt) {
                  final selected = c.isTab16Checked(rowIndex, opt);
                  return InkWell(
                    onTap: () => isRadio
                        ? c.setTab16Radio(rowIndex, opt)
                        : c.toggleTab16Checkbox(rowIndex, opt),
                    child: Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Icon(
                          isRadio
                              ? (selected
                                  ? Icons.radio_button_checked
                                  : Icons.radio_button_unchecked)
                              : (selected
                                  ? Icons.check_box
                                  : Icons.check_box_outline_blank),
                          size: 20,
                          color: selected ? Colors.blue : Colors.black54,
                        ),
                        const SizedBox(width: 3),
                        Text(opt, style: const TextStyle(fontSize: 15)),
                      ],
                    ),
                  );
                }).toList(),
              ),
            )),
      );
    }

    // 사진 셀
    Widget photoCell(int rowIndex) {
      return Container(
        width: photoW,
        padding: const EdgeInsets.symmetric(horizontal: 4, vertical: 4),
        decoration: deco,
        child: Obx(() {
          final images = c.getTab16Images(rowIndex);
          return Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  IconButton(
                    icon: const Icon(Icons.add_a_photo, size: 22),
                    padding: EdgeInsets.zero,
                    constraints: const BoxConstraints(),
                    onPressed: () =>
                        getTab16Image(rowIndex, ImageSource.camera),
                  ),
                  const SizedBox(width: 8),
                  IconButton(
                    icon: const Icon(Icons.add_photo_alternate, size: 22),
                    padding: EdgeInsets.zero,
                    constraints: const BoxConstraints(),
                    onPressed: () =>
                        getTab16Image(rowIndex, ImageSource.gallery),
                  ),
                ],
              ),
              if (images.isNotEmpty)
                Padding(
                  padding: const EdgeInsets.only(top: 4),
                  child: Wrap(
                    spacing: 2,
                    runSpacing: 2,
                    alignment: WrapAlignment.center,
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
          );
        }),
      );
    }

    // 데이터 행 (소분류 | 상태 | 사진)
    Widget dataRow(int rowIndex, String subLabel, List<String> options,
        {bool isRadio = false}) {
      return IntrinsicHeight(
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            labelCell(subLabel, width: col2W),
            stateCell(rowIndex, options, isRadio: isRadio),
            photoCell(rowIndex),
          ],
        ),
      );
    }

    // 그룹 행 (대분류 셀이 여러 소분류 행을 세로로 병합)
    Widget group(String majorLabel, List<Widget> subRows) {
      return IntrinsicHeight(
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            labelCell(majorLabel, width: col1W),
            Expanded(child: Column(children: subRows)),
          ],
        ),
      );
    }

    return SingleChildScrollView(
      padding: const EdgeInsets.all(10),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // 탭 이름
          const Padding(
            padding: EdgeInsets.only(bottom: 10),
            child: Text('부착물',
                style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
          ),
          // 헤더 행
          IntrinsicHeight(
            child: Row(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                labelCell('점검 내용',
                    width: col1W + col2W, weight: FontWeight.bold),
                Expanded(child: labelCell('상태', weight: FontWeight.bold)),
                labelCell('사진', width: photoW, weight: FontWeight.bold),
              ],
            ),
          ),
          // 부착물 등 (라디오)
          IntrinsicHeight(
            child: Row(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                labelCell('부착물 등', width: col1W + col2W),
                stateCell(0, const ['a', 'b', 'c', 'd', 'e'], isRadio: true),
                photoCell(0),
              ],
            ),
          ),
          // 정착부
          group('정착부', [
            dataRow(1, '앵커 및 브라켓', const [
              '상태양호',
              '콘크리트 균열',
              '콘크리트 박락',
              '앵커 시공',
              '앵커 풀림',
              '앵커 탈락',
              '앵커 부식'
            ]),
            dataRow(2, '용접', const ['상태양호', '면적 적정성', '균열', '부식', '손상']),
            dataRow(3, '매립', const ['상태양호', '철물매립 길이', '철물 여장(노출) 길이']),
          ]),
          // 연결부
          group('연결부', [
            dataRow(4, '볼트', const ['상태양호', '풀림', '탈락', '부재 변형', '부식']),
            dataRow(5, '용접', const ['상태양호', '면적 적정성', '균열', '부식', '손상']),
          ]),
          // 보강부
          group('보강부', [
            dataRow(
                6, '와이어 로프', const ['상태양호', '손상', '꼬임 및 뒤틀림', '변형', '고정클립 손상']),
          ]),
          // 변위 변형
          group('변위 변형', [
            dataRow(7, '기울기 및 배부름', const ['상태양호', '면외 방향 기울기', '배부름 발생']),
          ]),
          const SizedBox(height: 20),
          // 점검내용 입력
          /*Row(
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
        */
        ],
      ),
    );
  }
}
