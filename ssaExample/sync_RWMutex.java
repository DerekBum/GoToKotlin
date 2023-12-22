import java.util.Map;

public class sync_RWMutex {

	public sync_Mutex w;
	public long writerSem;
	public long readerSem;
	public atomic_Int32 readerCount;
	public atomic_Int32 readerWait;
}
