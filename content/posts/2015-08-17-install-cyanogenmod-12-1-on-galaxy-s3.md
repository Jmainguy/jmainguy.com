---
title: Install Cyanogenmod 12.1 on Galaxy S3
type: post
date: 2015-08-17T20:27:55+00:00
url: /index.php/install-cyanogenmod-12-1-on-galaxy-s3/
al2fb_facebook_link_id:
  - 55700847_10102199251335048
al2fb_facebook_link_time:
  - 2015-08-17T20:28:00+00:00
al2fb_facebook_link_picture:
  - media=http://jmainguy.com/?al2fb_image=1
categories:
  - General

---
Phone broke over the weekend, now I get to use the wife&#8217;s old S3 to replace my S3, I hope it can do the job =). I insist on using cyanogenmod on my phones as its an awesome OS. This is how I installed it today from my Fedora 22 laptop.

1. Install needed tools on Laptop
    
a. dnf install android-tools heimdall
  
2. Download gapps, cyanogenmod, and CWM
    
a. http://download2.clockworkmod.com/recoveries/recovery-clockwork-6.0.4.5-d2spr.img
    
b. https://download.cyanogenmod.org/get/jenkins/121988/cm-12.1-20150817-NIGHTLY-d2spr.zip
    
c. https://github.com/cgapps/vendor_google/blob/builds/arm/gapps-5.1-arm-2015-07-17-13-29.zip
  
3. Install CWM
    
a. Power off the Galaxy S III (Sprint) and connect the USB adapter to the computer but not to the Galaxy S III (Sprint), yet.
    
b. Boot the Galaxy S III (Sprint) into download mode. Hold Home, Volume Down & Power. Accept the disclaimer on the device. Then, insert the USB cable into the device.
    
c. heimdall flash &#8211;RECOVERY recovery-clockwork-6.0.4.5-d2spr.img &#8211;no-reboot
    
d. Unplug the USB cable from your device.
    
e. You can now manually reboot the phone into recovery mode. Hold Vol Up, Home & Power. Be sure to reboot into recovery immediately after having installed the custom recovery. Otherwise the custom recovery will be overwritten and the device will reboot (appearing as though your custom recovery failed to install).
  
4. Install Cyanogenmod and Gapps
    
a. After booting into CWM, and ensuring it installed correctly, boot back up the stock OS.
    
b. Transfer files (plug usb back in)
      
i. adb push cm-12.1-20150817-NIGHTLY-d2spr.zip /sdcard/
      
ii. adb push gapps-5.1-arm-2015-07-17-13-29.zip /sdcard/
    
c. Boot into CWM again.
      
i. reboot the phone into recovery mode. Hold Vol Up, Home & Power.
    
d. Choose install zips, and install both zips.
    
e. Lolipop Profit.

Sources:
  
http://wiki.cyanogenmod.org/w/Install\_CM\_for_d2spr
  
https://download.cyanogenmod.org/?device=d2spr
  
http://wiki.cyanogenmod.org/w/Doc:\_adb\_intro
  
https://clockworkmod.com/rommanager
  
https://github.com/cgapps/vendor_google/blob/builds/README.md