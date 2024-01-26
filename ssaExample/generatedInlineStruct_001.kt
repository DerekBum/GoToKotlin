package GoToJava

import java.io.BufferedReader
class generatedInlineStruct_001 {

	var block: ssa_BasicBlock? = null
}

fun read_generatedInlineStruct_001(buffReader: BufferedReader, id: Int): generatedInlineStruct_001 {
	val res = generatedInlineStruct_001()
    if (id != -1) {
        if (ptrMap.containsKey(id)) {
            return ptrMap[id] as generatedInlineStruct_001
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
    res.block = mapDec[readType]?.invoke(buffReader, id) as ssa_BasicBlock?

	buffReader.readLine()
	return res
}
