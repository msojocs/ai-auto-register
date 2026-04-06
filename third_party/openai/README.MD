# OpenAI Token Server

## Proof

这部分没啥，复制粘贴就能用。

## Turnstile

计算 Sentinel Token 的 `t` 值（`turnstileToken`）。

本质是一个JS字节码虚拟机，执行字节码获取浏览器环境的各种信息，然后组成字符串计算出Token。

使用`createFake`创建虚假的浏览器对象，关键时用 `throw new ReferenceError(name + "." + prop + " is not defined")` 不断跑抛常来知道缺什么字段，最后把这些字段都加上就能成功计算出Token了。

有时会用Nodejs debug看情况。

## 其他

使用[Fiddler Everywhere](https://github.com/msojocs/fiddler-everywhere-enhance)来把修改的JS映射到页面上，方便调试。