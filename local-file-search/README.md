
scan.go 编译生成 scan.exe，该程序首先从skipdir.txt中读取过滤掉的目录，每条目录占一行。

得到的文件信息存放在localfile.db(sqlite3数据库)中，localfile.db中只有一张表 info。

每次运行scan.exe程序都将删除info，然后重新构造 info 。

info 表的内容为：

文件名，路径
name path


