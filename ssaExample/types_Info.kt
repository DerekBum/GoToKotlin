class types_Info {

	var Types: Map<Any, types_TypeAndValue>? = null
	var Instances: Map<ast_Ident, types_Instance>? = null
	var Defs: Map<ast_Ident, Any>? = null
	var Uses: Map<ast_Ident, Any>? = null
	var Implicits: Map<Any, Any>? = null
	var Selections: Map<ast_SelectorExpr, types_Selection>? = null
	var Scopes: Map<Any, types_Scope>? = null
	var InitOrder: Array<types_Initializer>? = null
}
