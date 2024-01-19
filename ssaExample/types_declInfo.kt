class types_declInfo {

	var file: types_Scope? = null
	var lhs: Array<types_Var>? = null
	var vtyp: Any? = null
	var init: Any? = null
	var inherited: Boolean? = null
	var tdecl: ast_TypeSpec? = null
	var fdecl: ast_FuncDecl? = null
	var deps: Map<Any, Boolean>? = null
}
