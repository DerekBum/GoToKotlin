package GoToJava

import java.io.BufferedReader
class ssa_Const {

	var typ: Any? = null
	var Value: Any? = null
}

fun read_ssa_Const(buffReader: BufferedReader, id: Int): ssa_Const {
	val res = ssa_Const()
    if (id != -1) {
        if (ptrMap.containsKey(id)) {
            return ptrMap[id] as ssa_Const
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
    res.typ = mapDec[readType]?.invoke(buffReader, id) as Any?

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
    res.Value = mapDec[readType]?.invoke(buffReader, id) as Any?

	buffReader.readLine()
	return res
}
