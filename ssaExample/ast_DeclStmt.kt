package GoToJava

import java.io.BufferedReader
class ast_DeclStmt {

	var Decl: Any? = null
}

fun read_ast_DeclStmt(buffReader: BufferedReader, id: Int): ast_DeclStmt {
	val res = ast_DeclStmt()
    if (id != -1) {
        if (ptrMap.containsKey(id)) {
            return ptrMap[id] as ast_DeclStmt
        }
        ptrMap[id] = res
    }
    var line: String
    var split: List<String>
    var id: Int
    var readType: String

	line = buffReader.readLine()
	if (line == "end") {
        return res
    }
    split = line.split(" ")
    readType = split[1]
    id = -1
    if (split.size > 2) {
        id = split[2].toInt()
    }
    res.Decl = mapDec[readType]?.invoke(buffReader, id) as Any?

	buffReader.readLine()
	return res
}
