package noderunnerai;

import java.awt.Point;
import java.io.IOException;
import java.util.List;
import org.json.*;

/**
 *
 */
public abstract class BasicRunner {

    protected String name;
    protected String room;
    protected int level;

    private TCPClient client;

    public BasicRunner(String room, int level) {
        this.name = room;
        this.room = room;
        this.level = level;

        this.client = new TCPClient();
    }

    public BasicRunner(String room) {
        this(room, 1);
    }

    private void move(int direction) {
        client.send("move", "direction", "" + direction, "room", room);
    }

    private void dig(int direction) {
        client.send("dig", "direction", "" + direction, "room", room);
    }

    public void run() {
        this.client.join(name, room);

        try {
            String next = client.readNext();

            while (next != null) {

                JSONObject obj = new JSONObject(next);
                
                String event = obj.getString("event");
                
                if(event.equals("error")) {
                    System.err.println("Error : " + obj.getString("data"));
                    System.exit(-1);
                }
                
                JSONObject data = obj.getJSONObject("data");
                                
                switch(event) {
                    
                    case "start":
                        JSONArray tiles = data.getJSONArray("tiles");
                        
                        List<Object> list = tiles.toList();
                        
                        String[] arr = list.toArray(new String[list.size()]);
                        
                        start(arr);
                        break;
                        
                    case "next":
                        sendNext();
                        
                        break;
                }

                next = client.readNext();
            }
        } catch (IOException ex) {
            ex.printStackTrace();
        }
    }

    private void sendNext() {
        Move m = this.next();

        if (m.event == Event.MOVE) {
            move(m.direction.getValue());
        } else if (m.event == Event.DIG) {
            dig(m.direction.getValue());
        }
    }

    public abstract void start(String[] grid);

    public abstract Move next();
}
