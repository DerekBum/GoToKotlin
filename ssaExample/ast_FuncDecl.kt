class ast_FuncDecl {

	var Doc: ast_CommentGroup? = null
	var Recv: ast_FieldList? = null
	var Name: ast_Ident? = null
	var Type: ast_FuncType? = null
	var Body: ast_BlockStmt? = null
}
