## 通过应用包名备份恢复安卓应用数据
## Backup and restore Android app data by app package name

```shell
➜  app-data-backup-or-restore ./app_data_manager
请选择操作类型:
1: 备份应用数据
2: 恢复应用数据
请输入选项序号: 1
请输入应用包名: com.whatsapp
请输入本地备份路径 (默认: ./backup_data):
* daemon not running; starting now at tcp:5037
* daemon started successfully
removing leading '/' from member names
removing leading '/' from member names
/sdcard/com.whatsapp_backup_20241218_160419.tar.gz: 1 file pulled, 0 skipped. 27.1 MB/s (59285379 bytes in 2.086s)
./backup_data 备份完成!
➜  app-data-backup-or-restore ./app_data_manager
请选择操作类型:
1: 备份应用数据
2: 恢复应用数据
请输入选项序号: 2
请输入应用包名: com.whatsapp
请输入备份文件(绝对路径)路径: /Users/h1code2/GolandProjects/app-data-backup-or-restore/backup_data/com.whatsapp_backup_20241218_160419.tar.gz
/Users/h1code2/GolandProjects/app-data-backup-or-restore/backup_data/com.whatsapp_backup_20241218_160419.tar.gz: 1 file pushed, 0 skipped. 17.1 MB/s (59285379 bytes in 3.300s)
Success
SELinux: Loaded file_contexts
恢复完成!
```

MIT License

Copyright (c) [2024] [h1code2]

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
