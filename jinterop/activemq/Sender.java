//
// Copyright © 2012-2016 Guy M. Allard
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

import javax.jms.JMSException;
import javax.jms.Queue;
import javax.jms.QueueConnection;
import javax.jms.QueueConnectionFactory;
import javax.jms.QueueSender;
import javax.jms.QueueSession;
import javax.jms.TextMessage;
import javax.naming.Context;
import javax.naming.InitialContext;
import javax.naming.NamingException;

/**
 * Sends messages on the queue.
 */
public class Sender
{

  public static void main(String[] args) throws Exception
  {
 		System.out.println("Sends messages on the queue...");
    new Sender().go();
  } // end of main
  
  private void go() throws NamingException, JMSException
  {
    Context ictx = new InitialContext();
    QueueConnectionFactory qcf = (QueueConnectionFactory) ictx.lookup("ConnectionFactory");
    Queue queue = (Queue) ictx.lookup("Guys");
    ictx.close();

    QueueConnection qconn = qcf.createQueueConnection();
    QueueSession qsess = qconn.createQueueSession(true, 0);
    QueueSender qsend = qsess.createSender(queue);
    TextMessage msg = qsess.createTextMessage();

    int i;
    for (i = 0; i < Constants.NUM_MESSAGES; i++) {
    	if (Constants.DO_WORK_WAIT)
    	{
      	try
		    {
        		Thread.sleep(Constants.WORK_WAIT_MS); // Simulate message build
		    }
        	catch(InterruptedException ie)
		    {
        		// Ignore
		    }
    	}
    	String s = "Test Number " + i;
    	msg.setText(s);
    	qsend.send(msg);
   		System.out.println("Sent: " + s);
    	qsess.commit();	//
    }
   	System.out.println(i + " messages sent.");
    qconn.close();

  } // end of go
} // end of class Sender
