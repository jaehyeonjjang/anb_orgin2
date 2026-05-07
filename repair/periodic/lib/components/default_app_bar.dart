import 'package:common_control/common_control.dart';
import 'package:periodic/controllers/auth_controller.dart';

class DefaultAppBar extends CWidget {
  DefaultAppBar({super.key, this.leading, this.actions, this.centerActions});

  final c = Get.find<AuthController>();

  final Widget? leading;
  final List<Widget>? actions;
  final List<Widget>? centerActions;

  @override
  build(BuildContext context) {
    const titleBackgroundColor = Color.fromARGB(255, 101, 192, 240);

    var title = c.title;

    if (title == '') {
      if (c.periodic.id > 0) {
        title = c.periodic.apt.name;

        if (title == '') {
          title = c.apt.name;
        }
      } else if (c.apt.id > 0) {
        title = c.apt.name;
      }
    }

    final titleWidget = Obx(() => Text(c.title == '' ? title : c.title,
        style: const TextStyle(fontSize: 20, fontWeight: FontWeight.w500)));

    final hasCenter = centerActions != null && centerActions!.isNotEmpty;

    return AppBar(
      iconTheme: IconThemeData(
        color: Theme.of(context).primaryColor,
      ),
      title: titleWidget,
      backgroundColor: titleBackgroundColor,
      titleTextStyle: const TextStyle(color: Colors.black87),
      bottomOpacity: 0.0,
      elevation: 0.0,
      centerTitle: true,
      leading: leading,
      actions: actions,
      flexibleSpace: hasCenter
          ? SafeArea(
              child: Align(
                // 가로축 0.5 = 화면 가운데와 오른쪽 끝의 중간 지점
                alignment: const Alignment(0.5, 0),
                child: Row(
                  mainAxisSize: MainAxisSize.min,
                  children: centerActions!,
                ),
              ),
            )
          : null,
    );
  }
}
