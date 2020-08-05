# liveRedirect
用于提供http重定向访问yy/huya直播间M3U8,FLV直播流的作用，让它可以在potplayerer之类的播放器中直接播放，不需要在网页中忍受各种弹窗浮层。

http://项目ip:端口/platformName/房间id 的形式在potplayerer打开

## 编译&运行
1、安装golang环境

2、go mod tidy（好像是下载依赖的命令，我没试过，我都是用goland来开发自动下载依赖的）

3、go build -o lr(如果是交叉编译，需要设置编译平台，比如linux :SET GOOS=linux、 SET GOARCH=amd64)

4、执行lr (默认监听了5000端口)

### 已支持直播平台
| platformName |  平台| 
| :-----| ----: | 
| huya | 虎牙 | 
| yy | yy | 
| huajiao | 花椒 | 
| 2cp | 棉花糖直播 | 
| zhanqi | 战旗 | 
| 51lm | 羚萌直播 | 


### demo
http://170.178.163.226:5000/huya/774810

### 在potplayer中用播放列表的形式打开

将下面内容保存成playlist.dpl 然后就可以用potplayer 打开
```
DAUMPLAYLIST
playname=http://170.178.163.226:5000/huya/774810
playtime=4812
topindex=0
saveplaypos=0
1*file*http://170.178.163.226:5000/huya/616702
1*title*呆呆
2*file*http://170.178.163.226:5000/huya/774810
2*title*韩涵
2*duration2*10539383
3*file*http://170.178.163.226:5000/huya/599610
3*title*像风
4*file*http://170.178.163.226:5000/huya/912597
4*title*小蒜
4*duration2*19397155
5*file*http://170.178.163.226:5000/huya/kx2020
5*title*开心心
5*duration2*10539383


```


### 二次开发

参考[real_url](https://github.com/wbt5/real-url) 项目增加其他直播的解析就可以了。


## 代码逻辑
基本所有的解析逻辑都是照着下面的代码抄的
- [real_url](https://github.com/wbt5/real-url)
