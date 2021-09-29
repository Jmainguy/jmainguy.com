---
title: Delete email template duplicates from WHMCS

date: 2013-10-20T03:16:26+00:00
url: /index.php/delete-email-template-duplicates-from-whmcs/
al2fb_facebook_link_id:
  - 55700847_10101100818433388
al2fb_facebook_link_time:
  - 2013-10-20T03:16:28+00:00
al2fb_facebook_link_picture:
  - media=https://jmainguy.com/?al2fb_image=1
categories:
  - General

---
After upgrading WHMCS this weekend I noticed I had duplicate email templates. I used the following sql to delete them.

CREATE TABLE new_table AS
  
SELECT * FROM tblemailtemplates WHERE 1 GROUP BY name;

drop table tblemailtemplates;

RENAME TABLE new_table TO tblemailtemplates;