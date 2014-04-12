计算每个词的逆文档频率(Inverse Document Frequency，缩写为IDF):

逆文档频率(IDF) = log10( 语料库的文档总数 / (包含该词的文档数) +1 )


如果一个词越常见，那么分母就越大，逆文档频率就越小越接近0。分母之所以要加1，是为了避免分母为0（即所有文档都不包含该词）。

log10 表示对得到的值取以10为底的对数。


./idf.conf 配置文件： 

{
	"ExtractWebpagePath": "E:/SearchEngine/extract/webpage/",
	"ExtractUrlDataPath": "E:/SearchEngine/extract/extracturldata.db"
}

ExtractWebpagePath 中的文件名以 _body.txt 结尾的文档是抽取出来的网页正文，可以看做成一个语料库。
ExtractUrlDataPath 是将文档抽取了正文生成 _body.txt 的数据库。



然后将得到的逆文档频率信息以如下形式保存到sqlite3数据库中，数据库名为 idf.db ， 只有一个表为 data 。
同时也将这个数据保存到文本文件， idf.txt ，方便随时查看。

word  idf

例如：

中国  0.603
蜜蜂  2.713
养殖  2.410
