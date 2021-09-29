---
title: deleting files with sinatrarb

date: 2015-03-05T14:23:04+00:00
url: /index.php/deleting-files-with-sinatrarb/
al2fb_facebook_link_id:
  - 55700847_10101922362183128
al2fb_facebook_link_time:
  - 2015-03-05T14:23:13+00:00
al2fb_facebook_link_picture:
  - media=https://jmainguy.com/?al2fb_image=1
categories:
  - General

---
With sinatrarb you can create a quick web app in ruby very easily, their documentation is great.

http://www.sinatrarb.com/intro.html

One thing you can do, is send a file, like say a tar file.

send_file &#8220;certs/#{filename}/#{filename}.tar&#8221;, :filename => &#8220;#{filename}.tar&#8221;

I was hosting this on openshift which has a hard drive limitation of 1gb, so I needed to cleanup after each send. Turns out, you cannot simply place FileUtils.rm\_rf(&#8216;directorypath/name&#8217;) after send\_file to do the job. The reason is, send_file does not know its complete, and so FileUtils cannot remove the files its currently sending. The answer everyone on the internet suggests is scheduling a cron job to delete the files regularly. Since I was on openshift, and did not really want to use cron for this, I instead placed a job at the begining of each POST, that removes all directories older than 60 minutes.

In short, this is my code to delete directories older than 60 minutes.

def removedir (minutes)
    
%x[find certs/* -type d -cmin +#{minutes} -exec rm -rf {} \\;]
  
end

and then it gets called at the beginning of each POST

post &#8216;/json&#8217; do
    
removedir(60)