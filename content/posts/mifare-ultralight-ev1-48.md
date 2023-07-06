---
title: "Mifare Ultralight Ev1 48"
date: 2023-07-06T08:39:20-04:00
draft: false
---

We went on vacation to our condo in Myrtle Beach this week. The water was warm, and the fireworks went all night every day. The kids seemed to have a good time, my sister and parents also came with us which was a nice treat. 

It appears one of our previous renters must have placed their outdoor rusty beach chair on our carpet and setup for a bit, so there is now a rust stain on the carpet. The toilet flapper was completley worn thru and no longer functioned. Half the light bulbs were burnt out and needing to be replaced. The mattress covers, extra pillows, bed skirts, and oven mitts all seem to have dissapeared. Someone left us a waffle maker though, so that is cool. None of this is really unexpected, the renters tend to treat the place rather roughly and the maid's do not report any damage. Only time stuff gets fixed is when the next guest checks in and reports something. If you are a renter of an AirBnB or VRBO type rental in the future, do the owners a favor and report everything you notice broken so they can make it right for the next guests.

The buildings down here also use Mifare Ultralight EV1 (48 bytes) cards for entry (my sample of 2 buildings I have seen cards for). This is a convenient solution for allowing people inside. However there is a fine balance between convenience and security, and contactless cards are definitely not very secure. The way contactless cards work is as follows. There is a reader (on the door), and there is the card. The reader asks the card who they are, and then allows a whitelist of people inside. Anybody can stand there and overhear (sniff) the conversation and then tell the door they are that person (UID), and they will be let in. You can also just ask the card who they are, and then go tell the door the response.

The [proxmark3](https://hackerwarehouse.com/product/proxmark3-rdv4-kit/) is a device that makes sniffing conversations, reading cards, and writing cards quite easy. I have had a good experience using [Magic ntag 21x Ultralight](https://cyborg.ksecsolutions.com/product/magic-ntag-21x-ntag-i2c-mifare-ultralight-ev1-compatible-tag/) cards from ksec for writing to. I recently also purchased the same kind of tag in a [key fob](https://www.aliexpress.us/item/3256802532239642.html) format that works identically. These cards can mimic a number of mifare cards.

 - UL EV1 48b (The one I always use)
 - UL EV1 128b
 - NTAG 210
 - NTAG 212
 - NTAG 213
 - NTAG 215
 - NTAG 216
 - NTAG I2C 1K
 - NTAG I2C 2K
 - NTAG I2C 1K PLUS
 - NTAG I2C 2K PLUS
 - NTAG 213F
 - NTAG 216F

You can use [this script](https://github.com/RfidResearchGroup/proxmark3/blob/master/client/luascripts/hf_mfu_magicwrite.lua) with the commands referenced in this very nice [blog tutorial](https://lab401.com/blogs/academy/magic-ntag-21x-getting-started) for configuring the card. Just change the name of the script being run on that blog to hf_mfu_magicwrite instead of mfu_magic

To read the original card you want to clone, place the card ontop of your proxmark3 reader, go into the pm3 shell (by running pm3), and then run `hf mfu dump` to dump to a file `hf mfu info` is also a fun command to get some basic information. The dump command will write a bin, eml, and json file, for my purposes I only really care about the json file for now.

Since I do not know the key what was used to write blocks 16-19 (counting from 0) those are not dumped. (The readers for the doors do not use these blocks, they can be used for in house credit / tokens / fake money at hotel bar's and such if your place utilizes that)

I then wrote (well prompted chatGPT to write for me) a golang program to read the signature, and blocks from the json file, and write them to the magic ntag clone. You can read the [source code here](https://github.com/Jmainguy/ev1write). pm3 takes -s as an argument for a list of commands. This is likely a better route for performance, however I wanted to send each command one at a time with -c, and look at the response. This allows me to retry the command if it failed with a timeout (as happens when writing to the card quickly). There is also a hf mfu restore command that is supposed to do basically what this golang program is doing, however it was not working consistent enough for my liking. I think a third alternative would be to write this as a lua or python script for pm3 to natively use as well, however I like golang so I went this route.

These scripts and golang program make it much easier for me to dump, clone, and restore cards going forward (I hope).
