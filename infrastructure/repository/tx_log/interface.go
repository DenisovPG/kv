package tx_log

type TxLogger interface {
	LogPut(key, value string)
	LogGet(key string)
	LogDelete(key string)
	Run()
	Stop()
}
