import java.util.Map;

public class types_environment {

	public types_declInfo decl;
	public types_Scope scope;
	public int pos;
	public Object iota;
	public Object errpos;
	public Boolean inTParamList;
	public types_Signature sig;
	public Map<ast_CallExpr, Boolean> isPanic;
	public Boolean hasLabel;
	public Boolean hasCallOrRecv;
}
