---
title: I wrote a sms survey application

date: 2017-07-10T18:26:08+00:00
url: /index.php/i-wrote-a-sms-survey-application/
al2fb_facebook_link_id:
  - 55700847_10103583294316818
al2fb_facebook_link_time:
  - 2017-07-10T18:26:15+00:00
al2fb_facebook_link_picture:
  - custom=https://pbs.twimg.com/profile_images/2789370216/8765b6ba61039a987bdc1b3bc922bdbf_400x400.png
categories:
  - General

---
My pastor was looking for a sms survey service, requirements: A number to give out, a way to gather the responses, easy to use, cheap as possible. Googling turned up a few results with no clear pricing and salesmen, one I found was about $50 a month. Since I work at a telephone company, I decided just to make one using our API.

I am calling the application &#8220;Survey Hail&#8221;, as its a survey service, and a play of words on the pastors name. As with all code I write, this will is open sourced. It works as follows.

If you are the admin running the service for yourself / others. Assuming you are running a Linux server (you really should be, its the only real server OS in the world). Install docker and golang. Download the open-source code <https://github.com/Jmainguy/survey_hail>. Run the make command. Sign up for an account at [The Bandwidth API][1]. Pricing is $0.35 a month for a number, inbound sms is free (people texting your app their responses), outbound MMS to deliver the report is $0.01 a message. So, easily under $1 a month to run. You get your api / credentials, add them to a config.yaml file the service uses, then run run.sh. The server is now running.

The server checks in every 30 seconds for new messages, it stores messages it finds in a sqlite3 database, if it receives a message with the text of &#8220;report&#8221; from a admin phone number, it will send a csv report of all other text messages to that number. Users text their survey responses to the number, the admin texts &#8220;report&#8221; when they want a report.

Ez pz, lemon squeezy.

 [1]: https://catapult.inetwork.com/pages/home.jsf