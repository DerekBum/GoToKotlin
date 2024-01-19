class ssa_FreeVar {

	var name: String? = null
	var typ: Any? = null
	var pos: Int? = null
	var parent: ssa_Function? = null
	var referrers: Array<Any>? = null
	var outer: Any? = null
}
