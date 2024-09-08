package long_memory_task

import "chat_service/repository/storage/tcvectordb"

type LongMemoryTask struct {
	db *tcvectordb.TCVectorDB
}

func Init(db *tcvectordb.TCVectorDB) *LongMemoryTask {
	task := &LongMemoryTask{
		db: db,
	}
	return task
}

func (p *LongMemoryTask) Exec(input string) error {
	return nil
}
