class ssa_BasicBlock {

	var Index: Int? = null
	var Comment: String? = null
	var parent: ssa_Function? = null
	var Instrs: Array<Any>? = null
	var Preds: Array<ssa_BasicBlock>? = null
	var Succs: Array<ssa_BasicBlock>? = null
	var dom: ssa_domInfo? = null
	var gaps: Int? = null
	var rundefers: Int? = null
}
