participleutil是分词包

calculate-idf 计算每个词的逆文档频率 【可以根据某个标准语料库计算出来的，这样就无须频繁更新】

doc-mark 对文档进行标号

calculate-tf-idf 计算TF-IDF 并建立倒排索引



执行顺序：

（go run calculate_idf_main.go 平时不需要执行）

go run doc_mark_main.go

go run calculate_tf_idf_main.go

