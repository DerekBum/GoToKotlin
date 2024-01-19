class types_Checker {

	var conf: types_Config? = null
	var ctxt: types_Context? = null
	var fset: token_FileSet? = null
	var pkg: types_Package? = null
	var Info: types_Info? = null
	var version: types_version? = null
	var objMap: Map<Any, types_declInfo>? = null
	var impMap: Map<types_importKey, types_Package>? = null
	var valids: types_instanceLookup? = null
	var pkgPathMap: Map<String, Map<String, Boolean>>? = null
	var seenPkgMap: Map<types_Package, Boolean>? = null
	var files: Array<ast_File>? = null
	var posVers: Map<token_File, types_version>? = null
	var imports: Array<types_PkgName>? = null
	var dotImportMap: Map<types_dotImportKey, types_PkgName>? = null
	var recvTParamMap: Map<ast_Ident, types_TypeParam>? = null
	var brokenAliases: Map<types_TypeName, Boolean>? = null
	var unionTypeSets: Map<types_Union, types__TypeSet>? = null
	var mono: types_monoGraph? = null
	var firstErr: Any? = null
	var methods: Map<types_TypeName, Array<types_Func>>? = null
	var untyped: Map<Any, types_exprInfo>? = null
	var delayed: Array<types_action>? = null
	var objPath: Array<Any>? = null
	var cleaners: Array<Any>? = null
	var environment: types_environment? = null
	var indent: Int? = null
}
