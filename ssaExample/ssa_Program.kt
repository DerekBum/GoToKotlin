class ssa_Program {

	var Fset: token_FileSet? = null
	var imported: Map<String, ssa_Package>? = null
	var packages: Map<types_Package, ssa_Package>? = null
	var mode: Long? = null
	var MethodSets: typeutil_MethodSetCache? = null
	var canon: ssa_canonizer? = null
	var ctxt: types_Context? = null
	var methodsMu: sync_Mutex? = null
	var methodSets: typeutil_Map? = null
	var parameterized: ssa_tpWalker? = null
	var runtimeTypesMu: sync_Mutex? = null
	var runtimeTypes: typeutil_Map? = null
	var objectMethodsMu: sync_Mutex? = null
	var objectMethods: Map<types_Func, ssa_Function>? = null
}
