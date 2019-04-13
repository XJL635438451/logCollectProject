该日志收集系统，只是收集到了当前机器（不能跨机器）上的不同目录下的日志：
1. 如果后面支持收集多台机器的日志时用的是ip来区分不同机器（当前程序虽然看似支持多台机器，其实ip列表只有当前机器ip）
2. 一台机器可以收集多个路径下的日志

可以通过客户端动态的修改etcd中的数据(要收集的日志路径)：
1. 增加
2. 删除
3. 修改
修改完成之后后台可以热加载，实现动态的修改。
在etcd中保存的key:value格式如下：
key = /home/work/ 在配置文件中
etcd的key为上面的key与机器ip的拼接：
etcdKey = /home/work/192.168.1.0
value = "[{"path":"D:/project/nginx/logs/access2.log","topic":"nginx_log"},
        {"path":"D:/project/nginx/logs/error2.log","topic":"nginx_log_err"}]"
其中value要是字符串，value中存放的是一台机器的所有日志收集的路径和topic。