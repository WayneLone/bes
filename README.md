# By Example Scaffold

代码示例项目脚手架

## 由来

[go by example](https://github.com/mmcgrana/gobyexample) 这个项目把示例源码、模板、工具都放到一个项目里维护。对于模板和工具的变动实际上会很小，但是示例源码的变动就会很大。另外，工具都是基于 shell 脚本编写的，对于 windows 用户不是那么友好，甚至需要安装 go 语言 sdk 才能生成示例项目静态文件。

所以，索性就把模板和工具独立出来，将工具做成 cli，这样就可以不安装 go 语言运行时和使用 shell 脚本了。再者，可以自由的生成其它语言的 by example 了。

模板和工具是基于 [7958694](https://github.com/mmcgrana/gobyexample/commit/7958694c0ea91d3bef545cc4857a53e8c5eab48d) 这次提交构建的。如果发现工具和模板有新的变动，可以参考变动调整该脚手架的源码。

## 额外的东西

安利一下 gobyexample 这个项目，里面放的都是 go 的各种示例，可以作为速查手册

官网：https://gobyexample.com/
源码：https://github.com/mmcgrana/gobyexample

中文官网（同步有落后）：https://gobyexample-cn.github.io/
中文源码（同步有落后）：https://github.com/gobyexample-cn/gobyexample
