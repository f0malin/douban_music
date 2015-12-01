# douban_music
A music download robot for music.douban.com
一个帮你从豆瓣音乐人下载mp3的小工具。
主要是为了给大家学习go语言时做个参考，程序中用到了：http请求、下载、正则表达式、id3修改等技术。

使用方法：
1、在当前目录下新建music文件夹
2、找到豆瓣某音乐人的主页，如：http://site.douban.com/Ceekay/，作为参数传入即可。
$GOPATH/bin/douban_music http://site.douban.com/Ceekay/

目前下载后会自动更新mp3的id3 tag，不过只支持v2的，v1的请自行使用其他工具修改。
