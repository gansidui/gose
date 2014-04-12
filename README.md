A simple search engine in the golang

-----------------------------------------------------------------------

使用说明：

首先在 D 盘下建立文件夹 SearchEngine, 当然在其他路径下也可以，但必须得更改各个配置文件。



爬虫(spider):
	根据说明配置好爬虫， go run spider_main.go ， 下载的网页数据保存在 D:/SearchEngine/down/ 下。


分析(analyze):
	go run analyze_main.go ， 将爬取下来的网页提取出正文和标题， 提取出来的数据保存在 D:/SearchEngine/extract/ 下。

	
索引(indexing):
	根据 indexing 的说明建立倒排索引
	
查找(search):
	根据搜索串查找文档，用于web中的models模块。
	
显示(web):
	搜索引擎的界面，展示搜索结果。 go run main.go 启动服务器。


	
	
local-file-search 是本地磁盘文件搜索	
	
sqliteadmin 是管理sqlite数据库的可视化工具