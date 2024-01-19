import java.util.Map;

public class types_Checker {

	public types_Config conf;
	public types_Context ctxt;
	public token_FileSet fset;
	public types_Package pkg;
	public types_Info Info;
	public types_version version;
	public Map<Object, types_declInfo> objMap;
	public Map<types_importKey, types_Package> impMap;
	public types_instanceLookup valids;
	public Map<String, Map<String, Boolean>> pkgPathMap;
	public Map<types_Package, Boolean> seenPkgMap;
	public ast_File[] files;
	public Map<token_File, types_version> posVers;
	public types_PkgName[] imports;
	public Map<types_dotImportKey, types_PkgName> dotImportMap;
	public Map<ast_Ident, types_TypeParam> recvTParamMap;
	public Map<types_TypeName, Boolean> brokenAliases;
	public Map<types_Union, types__TypeSet> unionTypeSets;
	public types_monoGraph mono;
	public Object firstErr;
	public Map<types_TypeName, types_Func[]> methods;
	public Map<Object, types_exprInfo> untyped;
	public types_action[] delayed;
	public Object[] objPath;
	public Object[] cleaners;
	public types_environment environment;
	public Integer indent;
}
