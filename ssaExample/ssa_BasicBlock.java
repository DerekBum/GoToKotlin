import java.util.Map;

public class ssa_BasicBlock {

	public int Index;
	public String Comment;
	public ssa_Function parent;
	public Object[] Instrs;
	public ssa_BasicBlock[] Preds;
	public ssa_BasicBlock[] Succs;
	public ssa_domInfo dom;
	public int gaps;
	public int rundefers;
}
