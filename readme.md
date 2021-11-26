### APC比對工具

一般使用者操作步驟:

Step 1:
到[release]( http://tcaigitlab.corpnet.auo.com/mfg/l3d/01/apc_compare/-/releases "release")下載，並修改conf/config.ini參數

```commandline
#重要參數說明
[outRecord] 輸出檔案的名稱
CsvPath = compare.csv

[mail] alarm相關設定，可以自行修改收件者跟cc的人
Subject = APC資料錯誤
Content = APC資料與SPC不符，相關文件請參考附加檔案。
Send = Rico.Huang@auo.com
CC = Rico.Huang@auo.com;Rico.Huang@auo.com

..
```
Step 2:根據電腦版本64 or 32位元去執行不同的版本，只需要雙擊解壓縮後的apc_compare.exe便會進行比對與寄出警告信。



### 開發者操作步驟:

需安裝golang(可以請IT協助安裝)，並在windows安裝make指令(若不安裝可以複製Makefile中的build指令來做編譯)。

Step 1:安裝相關套件
``cmd
$ go mod install
``
Step 2:編譯
```cmd
$ make
#測試編譯後的檔案是否正確
$ make run
```

Step 3:release，記得修改makefile上方的版本號
```cmd
$ make release 
```