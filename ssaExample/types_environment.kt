class types_environment {

	var decl: types_declInfo? = null
	var scope: types_Scope? = null
	var pos: Int? = null
	var iota: Any? = null
	var errpos: Any? = null
	var inTParamList: Boolean? = null
	var sig: types_Signature? = null
	var isPanic: Map<ast_CallExpr, Boolean>? = null
	var hasLabel: Boolean? = null
	var hasCallOrRecv: Boolean? = null
}
