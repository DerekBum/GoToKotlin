class ast_File {

	var Doc: ast_CommentGroup? = null
	var Package: Int? = null
	var Name: ast_Ident? = null
	var Decls: Array<Any>? = null
	var FileStart: Int? = null
	var FileEnd: Int? = null
	var Scope: ast_Scope? = null
	var Imports: Array<ast_ImportSpec>? = null
	var Unresolved: Array<ast_Ident>? = null
	var Comments: Array<ast_CommentGroup>? = null
	var GoVersion: String? = null
}
