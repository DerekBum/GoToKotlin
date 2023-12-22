import java.util.Map;

public class types_Info {

	public Map<Object, types_TypeAndValue> Types;
	public Map<ast_Ident, types_Instance> Instances;
	public Map<ast_Ident, Object> Defs;
	public Map<ast_Ident, Object> Uses;
	public Map<Object, Object> Implicits;
	public Map<ast_SelectorExpr, types_Selection> Selections;
	public Map<Object, types_Scope> Scopes;
	public types_Initializer[] InitOrder;
}
