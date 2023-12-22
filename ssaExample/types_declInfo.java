import java.util.Map;

public class types_declInfo {

	public types_Scope file;
	public types_Var[] lhs;
	public Object vtyp;
	public Object init;
	public Boolean inherited;
	public ast_TypeSpec tdecl;
	public ast_FuncDecl fdecl;
	public Map<Object, Boolean> deps;
}
