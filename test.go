package main

import "irss-gateway/dispatcher"

func main() {
	content := "第一次写软件测评，在网络上收集了很多关于此类软件的使用，总结了一下：\n收趣官网：收趣 - 我的云端收藏夹，稍后阅读神器界面截图：PC网页截图Android界面截图介绍：    PC端只有WEB端，且收趣的插件在Google插件中心已下架（2021.6.9），这意味着你只能通过压缩包的形式运行收趣插件（但是不用其他方式每次启动会提示插件的安全性，会很麻烦）。插件支持FireFox/Chrome/Safari等等，由于插件在Google下架故没有评测。\n   收趣个人感觉无人维护了，官方客服从没有回复过。导出书签需要高级会员。\n优点：支持软件内打开网页，对部分网站/APP给予了适配（小猿APP/微信公众号等等），支持修改软件内访问的浏览器UA，支持免费缓存（仅适配的网页）/无限标签/无限书签等等，支持收藏时顺便进行分类（这个很重要，避免麻烦）\n缺点：无人维护，(未完待续)"
	hash, err := dispatcher.UploadPassage(content)
	if err != nil {
		panic(err)
	}
	print(hash)
}
