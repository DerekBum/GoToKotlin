class sync_RWMutex {

	var w: sync_Mutex? = null
	var writerSem: Long? = null
	var readerSem: Long? = null
	var readerCount: atomic_Int32? = null
	var readerWait: atomic_Int32? = null
}
