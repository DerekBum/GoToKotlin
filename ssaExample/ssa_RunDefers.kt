package GoToJava

import java.io.BufferedReader
class ssa_RunDefers {

	var anInstruction: generatedInlineStruct_000? = null
}

fun read_ssa_RunDefers(buffReader: BufferedReader, id: Int): ssa_RunDefers {
	val res = ssa_RunDefers()
    if (id != -1) {
        if (ptrMap.containsKey(id)) {
            return ptrMap[id] as ssa_RunDefers
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
    res.anInstruction = mapDec[readType]?.invoke(buffReader, id) as generatedInlineStruct_000?

	buffReader.readLine()
	return res
}
