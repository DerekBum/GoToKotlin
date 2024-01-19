class types_Scope {

	var parent: types_Scope? = null
	var children: Array<types_Scope>? = null
	var number: Int? = null
	var elems: Map<String, Any>? = null
	var pos: Int? = null
	var end: Int? = null
	var comment: String? = null
	var isFunc: Boolean? = null
}
