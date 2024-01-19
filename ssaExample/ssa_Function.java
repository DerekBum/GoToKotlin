import java.util.Map;

public class ssa_Function {

	public String name;
	public types_Func object;
	public ssa_selection method;
	public types_Signature Signature;
	public Integer pos;
	public String Synthetic;
	public Object syntax;
	public types_Info info;
	public String goversion;
	public ssa_Function parent;
	public ssa_Package Pkg;
	public ssa_Program Prog;
	public ssa_Parameter[] Params;
	public ssa_FreeVar[] FreeVars;
	public ssa_Alloc[] Locals;
	public ssa_BasicBlock[] Blocks;
	public ssa_BasicBlock Recover;
	public ssa_Function[] AnonFuncs;
	public Object[] referrers;
	public Integer anonIdx;
	public types_TypeParamList typeparams;
	public Object[] typeargs;
	public ssa_Function topLevelOrigin;
	public ssa_generic generic;
	public ssa_BasicBlock currentBlock;
	public Map<types_Var, Object> vars;
	public ssa_Alloc[] namedResults;
	public ssa_targets targets;
	public Map<types_Label, ssa_lblock> lblocks;
	public ssa_subster subst;
}
