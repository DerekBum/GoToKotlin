class types_Context {

	var mu: sync_Mutex? = null
	var typeMap: Map<String, Array<types_ctxtEntry>>? = null
	var nextID: Int? = null
	var originIDs: Map<Any, Int>? = null
}
