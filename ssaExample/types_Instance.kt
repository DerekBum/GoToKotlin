package GoToJava

import java.io.BufferedReader
import jacodbInst.*
class types_Instance {

	var TypeArgs: types_TypeList? = null
	var Type: Any? = null
}

fun read_types_Instance(buffReader: BufferedReader, id: Int): types_Instance {
	val res = types_Instance()
    if (id != -1) {
        if (ptrMap.containsKey(id)) {
            return ptrMap[id] as types_Instance
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
    res.TypeArgs = mapDec[readType]?.invoke(buffReader, id) as types_TypeList?

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

	buffReader.readLine()
	return res
}
