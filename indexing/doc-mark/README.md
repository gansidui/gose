对文档进行标号

文档是指那些已经被分析过了的文档，即已经抽取了正文和标题的文档， 在配置文件中指定了路径

{
	"ExtractUrlDataPath": "D:/SearchEngine/extract/extracturldata.db"
}

然后将标号的结果生成sqlite数据库保存，名为 docmark.db, 表为data
字段为：
md5    id

其中 md5 为主键，这样保证标号唯一，每次标号都先计算数据库中最大的标号值，然后再不断递增标号。

create table data(md5 varchar(32) not null primary key, id integer not null)