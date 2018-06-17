//
// Copyright © 2012-2018 Guy M. Allard
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

import java.util.Random;

import javax.jms.JMSException;
import javax.jms.Message;
import javax.jms.ObjectMessage;
import javax.jms.BytesMessage;
import javax.jms.MapMessage;
import javax.jms.StreamMessage;
import javax.jms.Queue;
import javax.jms.QueueConnection;
import javax.jms.QueueConnectionFactory;
import javax.jms.QueueReceiver;
import javax.jms.QueueSession;
import javax.jms.TextMessage;
import javax.naming.Context;
import javax.naming.InitialContext;
import javax.naming.NamingException;


/**
 * Requests messages on the queue.
 */
public class Receiver
{

  public static void main(String[] args) throws Exception
  {
 		System.out.println("Requests to receive messages...");
    new Receiver().go();
  } // end of main
  
  
  private void go() throws NamingException, JMSException
  {
  	Random gener = new Random(System.currentTimeMillis());
    Context ictx = new InitialContext();
    Queue queue = (Queue) ictx.lookup("Guys");
    QueueConnectionFactory qcf = (QueueConnectionFactory) ictx.lookup("ConnectionFactory");
    ictx.close();

    QueueConnection qconn = qcf.createQueueConnection();
    QueueSession qsess = qconn.createQueueSession(true, 0);
    QueueReceiver qrec = qsess.createReceiver(queue);
    Message msg;

    qconn.start();
   	System.out.println("Waits are for " + Constants.RECEIVE_WAIT + "(ms)");
    int i;
    // for (i = 0; i < 10; i++) {
    for (i = 0;; i++) {    	
      // msg = qrec.receive();
      msg = qrec.receive(Constants.RECEIVE_WAIT);
      if (msg == null) break;
      if (msg instanceof TextMessage)
      {
     		System.out.println("Text Msg received: " + ((TextMessage) msg).getText());
      }
      else if (msg instanceof ObjectMessage)
      {
    		System.out.println("Object Msg received: "
           + ((ObjectMessage) msg).getObject());
      }
      else if (msg instanceof BytesMessage)
      {
            BytesMessage bmsg = (BytesMessage)msg;
            long len = bmsg.getBodyLength();
            byte[] buff = new byte[(int)len];
            int rlen = bmsg.readBytes(buff);
            String body = new String(buff);
    		System.out.println("Bytes Msg received: "
                + "len:" + rlen + "\n" + "body:" + body);
      }
      else if (msg instanceof StreamMessage)
      {
    		System.out.println("Stream Msg received: ");
            System.out.println(msg);
      }
      else if (msg instanceof MapMessage)
      {
    		System.out.print("Map Msg received: ");
            System.out.println(msg);
      }
      else
      {
     		System.out.println("Some Other Msg received: " + msg);
      }
      int nextWait = Constants.MIN_WORK + gener.nextInt(Constants.WORK_SPAN);
     	System.out.println("Next Input Process Wait Length: " + nextWait);
      try
	    {
       	Thread.sleep(nextWait);	// process work
	    }
        catch(InterruptedException ie)
	    {
        	//
	    }
    }
    qsess.commit();
   	System.out.println(i + " messages received.");
    qconn.close();
  	
  }
} // end of class Receiver
