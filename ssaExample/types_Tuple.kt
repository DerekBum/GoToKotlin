package GoToJava

import java.io.BufferedReader
class types_Tuple {

	var vars: List<types_Var>? = null
}

fun read_types_Tuple(buffReader: BufferedReader, id: Int): types_Tuple {
	val res = types_Tuple()
    if (id != -1) {
        if (ptrMap.containsKey(id)) {
            return ptrMap[id] as types_Tuple
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
    res.vars = mapDec[readType]?.invoke(buffReader, id) as List<types_Var>?

	buffReader.readLine()
	return res
}
