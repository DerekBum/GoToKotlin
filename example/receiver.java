import java.io.BufferedReader;
import java.io.FileNotFoundException;
import java.io.FileReader;
import java.io.IOException;
import java.util.*;

public class receiver {
    private int lineId = 0;

    public Object returnFilled(List<String> all) {
        String[] split = all.get(lineId).split(" ");
        
    }

    public static void main(String[] args) throws IOException {
        BufferedReader br = new BufferedReader(new FileReader("filled.txt"));
        List<String> sb = new ArrayList<String>();

        try {
            String line = br.readLine();

            while (line != null) {
                sb.add(line);
                line = br.readLine();
            }
        } finally {
            br.close();
        }

        Object res = returnFilled(sb);
    }
}
