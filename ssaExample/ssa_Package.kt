class ssa_Package {

	var Prog: ssa_Program? = null
	var Pkg: types_Package? = null
	var Members: Map<String, Any>? = null
	var objects: Map<Any, Any>? = null
	var init: ssa_Function? = null
	var debug: Boolean? = null
	var syntax: Boolean? = null
	var buildOnce: sync_Once? = null
	var ninit: Int? = null
	var info: types_Info? = null
	var files: Array<ast_File>? = null
	var created: Array<ssa_Function>? = null
	var initVersion: Map<Any, String>? = null
}
