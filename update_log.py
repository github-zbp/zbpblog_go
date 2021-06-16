# coding=utf-8

import time, os
# 该脚本用于定时更改日志文件的日期已达到一天一个日志文件的目的
date = time.strftime("%Y%m%d", time.localtime())

# 从配置文件获取日志文件路径
log_path = ""
with open("./conf/app.conf", "r", encoding="utf-8") as conf:
    while True:
        line = conf.readline()
        if line.find("log_path") != -1:
            line = line.replace(" ", "").strip()
            log_path = line.split("=")[-1]
            break

if log_path:
    log_file = os.path.basename(log_path)
    name, ext = log_file.split(".")
    new_name = name + "-" + date
    new_log_path = os.path.dirname(log_path).strip("/") + "/" + new_name + "." + ext
    os.rename(log_path, new_log_path)