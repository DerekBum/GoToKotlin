class ssa_Function {

	var name: String? = null
	var Object: types_Func? = null
	var method: ssa_selection? = null
	var Signature: types_Signature? = null
	var pos: Int? = null
	var Synthetic: String? = null
	var syntax: Any? = null
	var info: types_Info? = null
	var goversion: String? = null
	var parent: ssa_Function? = null
	var Pkg: ssa_Package? = null
	var Prog: ssa_Program? = null
	var Params: Array<ssa_Parameter>? = null
	var FreeVars: Array<ssa_FreeVar>? = null
	var Locals: Array<ssa_Alloc>? = null
	var Blocks: Array<ssa_BasicBlock>? = null
	var Recover: ssa_BasicBlock? = null
	var AnonFuncs: Array<ssa_Function>? = null
	var referrers: Array<Any>? = null
	var anonIdx: Int? = null
	var typeparams: types_TypeParamList? = null
	var typeargs: Array<Any>? = null
	var topLevelOrigin: ssa_Function? = null
	var generic: ssa_generic? = null
	var currentBlock: ssa_BasicBlock? = null
	var vars: Map<types_Var, Any>? = null
	var namedResults: Array<ssa_Alloc>? = null
	var targets: ssa_targets? = null
	var lblocks: Map<types_Label, ssa_lblock>? = null
	var subst: ssa_subster? = null
}
