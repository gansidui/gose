
这是一个独立的爬虫程序。

./spider.conf 是爬虫的配置文件，内容如下：

{
	"EntryPath": "D:/golib/src/github.com/gansidui/gose/spider/config/entry.txt",
	"UrlQueuePath": "D:/golib/src/github.com/gansidui/gose/spider/config/urlqueue.db",
	"FilterExtPath": "D:/golib/src/github.com/gansidui/gose/spider/config/filterext.txt",
	"LimitedToThisPath": "D:/golib/src/github.com/gansidui/gose/spider/config/limitedtothis.txt",
	"DownWebpagePath": "D:/SearchEngine/down/webpage/",
	"DownUrlDataPath": "D:/SearchEngine/down/downurldata.db",
	"MaxNum": 20000,
	"IntervalTime": "50ms"
}


这是一个json数据结构，需要严格按照json的格式填写。

带有 / 后缀的是文件夹，必须保证.txt的格式为UTF-8无BOM，windows下默认是有BOM头的，用notepad转换下无BOM即可。

.db文件是sqlite数据库。



"EntryPath": 是爬虫的入口地址，存放在一个文本文件中，每行为一个url， 也就是种子地址，url需要带有scheme,例如 http:// 和 https:// 。

"UrlQueuePath": 是爬虫上次爬行过程中的url队列，这些url等待爬虫的分析，分析url后提取网页中的url，再将提取出来的url插入到队列尾部，如此循环。 
在程序中断后，可以根据UrlQueue中保存的url数据继续爬取互联网，不需要再从入口地址开始爬一遍。这个UrlQueue也可以看做是入口地址。
数据库中只有一个表，名称为data，字段有 md5, url 。

"FilterExtPath": 过滤掉指定扩展名的文件，每行一个扩展名。

"LimitedToThisPath": 指明爬虫仅仅抓取包含了这些子串的url，url只需包含其中一条即可。其他所有的url都将被过滤。每行一个字符串。这条规则允许你可以抓取指定的网站数据。

"DownWebpagePath": 爬虫抓取下来的网页保存在该文件夹中，文件名为url的md5值。

"DownUrlDataPath": 是一个sqlite数据库，以前爬过的网页数据，这样爬虫此次爬的时候就会过滤掉这些ulr。数据库中只有一个表，名称为data，字段有 md5, url, path 。 (path为网页的本地存储路径)

"MaxNum": 预期爬下的网页数

"IntervalTime": 间歇时间，每次爬取一个网页后休息一会，为0的话也意味着放弃本时间片，抓取频率太高有被封掉的风险。建议设置为 500ms



爬虫的爬取策略：

每次运行爬虫，首先从 DownUrlData 中读取历史数据，然后从 UrlQueue 中读入起始url，再读入 Entry 中的入口地址，当然读入的时候得先判断是否以前爬取了这个网页。

多个线程负责爬取网页，一个线程负责将网页数据写入磁盘。为了均衡网站的负荷，每次抓取都暂停 IntervalTime 。












