import java.util.Map;

public class ssa_BasicBlock {

	public Integer Index;
	public String Comment;
	public ssa_Function parent;
	public Object[] Instrs;
	public ssa_BasicBlock[] Preds;
	public ssa_BasicBlock[] Succs;
	public ssa_domInfo dom;
	public Integer gaps;
	public Integer rundefers;
}
