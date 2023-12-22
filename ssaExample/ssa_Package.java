import java.util.Map;

public class ssa_Package {

	public ssa_Program Prog;
	public types_Package Pkg;
	public Map<String, Object> Members;
	public Map<Object, Object> objects;
	public ssa_Function init;
	public Boolean debug;
	public Boolean syntax;
	public sync_Once buildOnce;
	public int ninit;
	public types_Info info;
	public ast_File[] files;
	public ssa_Function[] created;
	public Map<Object, String> initVersion;
}
