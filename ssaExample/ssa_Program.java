import java.util.Map;

public class ssa_Program {

	public token_FileSet Fset;
	public Map<String, ssa_Package> imported;
	public Map<types_Package, ssa_Package> packages;
	public long mode;
	public typeutil_MethodSetCache MethodSets;
	public ssa_canonizer canon;
	public types_Context ctxt;
	public sync_Mutex methodsMu;
	public typeutil_Map methodSets;
	public ssa_tpWalker parameterized;
	public sync_Mutex runtimeTypesMu;
	public typeutil_Map runtimeTypes;
	public sync_Mutex objectMethodsMu;
	public Map<types_Func, ssa_Function> objectMethods;
}
