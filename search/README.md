{
	"DownUrlDataPath": "D:/SearchEngine/down/downurldata.db",
	"ExtractWebpagePath": "D:/SearchEngine/extract/webpage/",
	"ExtractUrlDataPath": "D:/SearchEngine/extract/extracturldata.db",
	"DicPath": "D:/golib/src/github.com/gansidui/gose/indexing/participleutil/mydic.txt",
	"TfIdfPath": "D:/golib/src/github.com/gansidui/gose/indexing/calculate-tf-idf/tfidf.db",
	"DocMarkPath": "D:/golib/src/github.com/gansidui/gose/indexing/doc-mark/docmark.db"
}

这个包是给web中的models模块使用的，配置文件用的全部是绝对路径。

搜索的原理：

初始化时读取 
word --> id, tfidf
id --> md5
md5 --> url, path

输入一个搜索串，将其按空格拆分成多个句子，然后分别对这些句子进行搜索，

对每个句子进行搜索：将句子进行分词得到多个关键词，通过关键词的索引得到文档id集合，对这些id集合取并集。

然后将各个句子搜索的结果取交集。（这个符合用户的搜索习惯）