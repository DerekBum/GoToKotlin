package GoToJava

import java.io.BufferedReader
import jacodbInst.*
class ssa_Builtin : ssaToJacoExpr, ssaToJacoValue {

	var name: String? = null
	var sig: types_Signature? = null

	override fun createJacoDBExpr(): GoBuiltin {
        return GoBuiltin(
            0,
            name!!,
            sig!!
        )
    }

    override fun createJacoDBValue(): GoValue {
        return createJacoDBExpr()
    }
}

fun read_ssa_Builtin(buffReader: BufferedReader, id: Int): ssa_Builtin {
	val res = ssa_Builtin()
    if (id != -1) {
        if (ptrMap.containsKey(id)) {
            return ptrMap[id] as ssa_Builtin
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
    res.name = mapDec[readType]?.invoke(buffReader, id) as String?

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
    res.sig = mapDec[readType]?.invoke(buffReader, id) as types_Signature?

	buffReader.readLine()
	return res
}
