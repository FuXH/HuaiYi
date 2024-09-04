package tcvectordb

import (
	"context"
	"fmt"
	"time"

	"chat_service/entity"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/ai_document_set"
)

type TCVectorDB struct {
	client *tcvectordb.RpcClient
}

func NewTCVectorDB(url, user, key string) (*TCVectorDB, error) {
	cli, err := tcvectordb.NewRpcClient(url, user, key, &tcvectordb.ClientOption{
		Timeout:         30 * time.Second,
		ReadConsistency: tcvectordb.EventualConsistency})
	fmt.Println("NewRpcClient err: ", err, cli)
	if err != nil {
		fmt.Println("NewRpcClient err: ", err)
		return nil, err
	}
	return &TCVectorDB{
		client: cli,
	}, nil
}

// ListDatabase
// test
func (p *TCVectorDB) ListDatabase(ctx context.Context) {
	dbList, err := p.client.ListDatabase(ctx)
	if err != nil {
		fmt.Println("ListDatabase err: ", err)
	}
	for _, val := range dbList.AIDatabases {
		fmt.Println("AIDatabases: ", val)
	}
}

// LoadAndSplitText 上传doc(markdown格式)
func (p *TCVectorDB) LoadAndSplitText(ctx context.Context, database, collection string,
	filePath string, metadata map[string]interface{}) (*tcvectordb.LoadAndSplitTextResult, error) {
	coll := p.client.AIDatabase(database).CollectionView(collection)
	res, err := coll.LoadAndSplitText(ctx, tcvectordb.LoadAndSplitTextParams{
		LocalFilePath: filePath,
		MetaData:      metadata,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

// SearchByFileName 通过doc名称查找
func (p *TCVectorDB) SearchByFileName(ctx context.Context, database, collection string,
	searchMessage *entity.TCVectorSearchMessage) ([]tcvectordb.AIDocumentSet, error) {
	coll := p.client.AIDatabase(database).CollectionView(collection)
	res, err := coll.Query(ctx, tcvectordb.QueryAIDocumentSetParams{
		DocumentSetName: []string{searchMessage.FileName},
	})
	if err != nil {
		return nil, err
	}
	return res.Documents, nil
}

// SearchByFilterSql 通过doc的metadata查找，返回匹配的N个doc
func (p *TCVectorDB) SearchByFilterSql(ctx context.Context, database, collection string,
	searchMessage *entity.TCVectorSearchMessage) ([]tcvectordb.AIDocumentSet, error) {
	db := p.client.AIDatabase(database)
	coll := db.CollectionView(collection)

	res, err := coll.Query(ctx, tcvectordb.QueryAIDocumentSetParams{
		Filter: tcvectordb.NewFilter(searchMessage.FilterSql),
		Limit:  searchMessage.Limit,
	})
	if err != nil {
		return nil, err
	}
	return res.Documents, nil
}

// SearchByText 通过filterSql和输入文本查找，返回匹配的N个doc
func (p *TCVectorDB) SearchByText(ctx context.Context, database, collection string,
	searchMessage *entity.TCVectorSearchMessage) ([]tcvectordb.AISearchDocumentSet, error) {
	db := p.client.AIDatabase(database)
	coll := db.CollectionView(collection)

	// 查找与给定查询向量相似的向量。支持输入文本信息检索与输入文本相似的内容，同时，支持搭配标量字段的 Filter 表达式一并检索。
	enableRerank := true
	res, err := coll.Search(ctx, tcvectordb.SearchAIDocumentSetsParams{
		Content:     searchMessage.Content,
		ExpandChunk: []int{1, 0},
		Filter:      tcvectordb.NewFilter(searchMessage.FilterSql),
		Limit:       searchMessage.Limit,
		RerankOption: &ai_document_set.RerankOption{
			Enable:                &enableRerank,
			ExpectRecallMultiples: 2.5,
		},
	})
	if err != nil {
		return nil, err
	}
	return res.Documents, nil
}

// Update 更新doc的metadata
func (p *TCVectorDB) Update(ctx context.Context, database, collection string,
	updateField map[string]interface{}, filterSql string) error {
	db := p.client.AIDatabase(database)
	coll := db.CollectionView(collection)

	_, err := coll.Update(ctx, updateField, tcvectordb.UpdateAIDocumentSetParams{
		Filter: tcvectordb.NewFilter(filterSql),
	})
	if err != nil {
		return err
	}
	return nil
}
