package GoToJava

import java.io.BufferedReader
class ast_SelectorExpr {

	var X: Any? = null
	var Sel: ast_Ident? = null
}

fun read_ast_SelectorExpr(buffReader: BufferedReader, id: Int): ast_SelectorExpr {
	val res = ast_SelectorExpr()
    if (id != -1) {
        if (ptrMap.containsKey(id)) {
            return ptrMap[id] as ast_SelectorExpr
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
    res.X = mapDec[readType]?.invoke(buffReader, id) as Any?

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
    res.Sel = mapDec[readType]?.invoke(buffReader, id) as ast_Ident?

	buffReader.readLine()
	return res
}
