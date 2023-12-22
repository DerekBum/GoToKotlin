import java.util.Map;

public class ast_File {

	public ast_CommentGroup Doc;
	public int Package;
	public ast_Ident Name;
	public Object[] Decls;
	public int FileStart;
	public int FileEnd;
	public ast_Scope Scope;
	public ast_ImportSpec[] Imports;
	public ast_Ident[] Unresolved;
	public ast_CommentGroup[] Comments;
	public String GoVersion;
}
