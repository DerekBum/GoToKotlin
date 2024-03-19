package GoToJava

import java.io.BufferedReader
import jacodbInst.*
import jacodbInst.impl.location.GoInstLocationImpl
class ssa_Jump : ssaToJacoInst {

	var anInstruction: generatedInlineStruct_000? = null

	override fun createJacoDBInst(parent: GoMethod): GoJumpInst {
        return GoJumpInst(
            GoInstLocationImpl(
                anInstruction!!.block!!.Index!!.toInt(),
                0,
                parent,
            ),
            parent,
            GoInstRef(
                anInstruction!!.block!!.Index!!.toInt()
            )
        )
    }
}

fun read_ssa_Jump(buffReader: BufferedReader, id: Int): ssa_Jump {
	val res = ssa_Jump()
    if (id != -1) {
        if (ptrMap.containsKey(id)) {
            return ptrMap[id] as ssa_Jump
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
