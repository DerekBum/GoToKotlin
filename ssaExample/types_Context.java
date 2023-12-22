import java.util.Map;

public class types_Context {

	public sync_Mutex mu;
	public Map<String, types_ctxtEntry[]> typeMap;
	public int nextID;
	public Map<Object, int> originIDs;
}
