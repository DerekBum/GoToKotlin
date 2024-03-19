package GoToJava

import java.io.BufferedReader
import jacodbInst.*
class ast_ValueSpec {

	var Doc: ast_CommentGroup? = null
	var Names: List<ast_Ident>? = null
	var Type: Any? = null
	var Values: List<Any>? = null
	var Comment: ast_CommentGroup? = null
}

fun read_ast_ValueSpec(buffReader: BufferedReader, id: Int): ast_ValueSpec {
	val res = ast_ValueSpec()
    if (id != -1) {
        if (ptrMap.containsKey(id)) {
            return ptrMap[id] as ast_ValueSpec
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
    res.Doc = mapDec[readType]?.invoke(buffReader, id) as ast_CommentGroup?

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
    res.Names = mapDec[readType]?.invoke(buffReader, id) as List<ast_Ident>?

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
    res.Type = mapDec[readType]?.invoke(buffReader, id) as Any?

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
    res.Values = mapDec[readType]?.invoke(buffReader, id) as List<Any>?

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
    res.Comment = mapDec[readType]?.invoke(buffReader, id) as ast_CommentGroup?

	buffReader.readLine()
	return res
}
