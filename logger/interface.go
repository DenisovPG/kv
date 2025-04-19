package logger

type TxLogger interface {
	LogPut(key, value string)
	LogGet(key string)
	LogDelete(key string)
	Run()
	Stop()
}
