package main

type UpdateStrategy string
type OperationName string

const (
	IncrementalUpdateStrategy UpdateStrategy = "INCREMENTAL"
	BatchUpdateStrategy       UpdateStrategy = "BATCH"
)

const (
	CreateCollectionOperationName OperationName = "CREATE_COLLECTION"
	InsertOneOperationName        OperationName = "INSERT_ONE"
	DeleteOneOperationName        OperationName = "DELETE_ONE"
	UpdateOneOperationName        OperationName = "UPDATE_ONE"
)
