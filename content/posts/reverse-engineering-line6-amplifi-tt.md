---
title: "Reverse Engineering the Line 6 AMPLIFi TT for Linux"
date: 2026-09-09T12:55:00-04:00
draft: false
---

Reverse engineering and disassembly have intrigued me for a while. However, tools like Ghidra intimidated me, and I had no idea how to get started with them. That was about where I stayed until recently.

I have been having so much fun using LLM tools that I figured I would see if one could disassemble a simple program. I picked *Lords of the Realm*, a DOS game from the '90s that I enjoyed playing as a kid. The binary was small, so I figured it was a good place to start.

It turns out GPT-5.6 Sol is pretty great at using Ghidra and disassembling things. It made significant progress figuring out how the game worked and found all sorts of cool information. I had burned a week's worth of tokens letting it work on that when I saw [Everything I own, owned](https://news.ycombinator.com/item?id=49413320) on Hacker News. It inspired me to dig the AMPLIFi TT out of my closet.

![The front of my Line 6 AMPLIFi TT](https://immich.soh.re/api/assets/0e1c687a-3e6e-4aca-889a-f321137daa04/original?key=te_XqV_KVlnJ167prVYS86fw22xZkdwf_ny-FFpOPmhbHEdY8AEgLb2BAkLI3hZdMak)

I originally bought the device using my company karma points at Red Hat. You could give away so many points per month, and the person receiving them could cash them in for things like this or gift cards. You got taxed as if the points converted 1:1 to USD, but you could not actually redeem them at that rate. Flawed, but still free money.

The TT advertised the ability to match its guitar tone to songs, like some kind of magical AI. I could not figure out how that would work and was skeptical, but I had the karma points to spend and it looked interesting. I wanted a way to plug my guitar into the computer and get cool metal sounds without buying a whole pedalboard.

I was sorely disappointed when I got it. The hardware looked and felt great. The software was hot garbage, even back in 2014 or 2015.

Despite having USB, it did not let me operate the device from my computer. USB was basically there to use it as a glorified headphone amp and install firmware updates. To use the tones and presets, I had to connect a phone or tablet over Bluetooth, which worked poorly. The device would fail to pair, disconnect halfway through a session, and make browsing community tones a chore.

![The rear of the AMPLIFi TT, with its amp, RCA, main, optical, USB, FBV, and power connections](https://immich.soh.re/api/assets/29f8ab0e-b849-4356-927a-7c941d2f4952/original?key=te_XqV_KVlnJ167prVYS86fw22xZkdwf_ny-FFpOPmhbHEdY8AEgLb2BAkLI3hZdMak)

The magical AI turned out to be song metadata. It grabbed the artist and song title, then searched the Line 6 database for those tags. The tones sounded great, but the mobile app was horrendous. I got so fed up with it that I put the TT in my closet, where it stayed for the next decade.

Now, a decade later, AI is actually a thing, and it can do disassembly and reverse engineering. With this in mind, and inspired by the Hacker News post, I plugged the TT in and began prompting GPT-5.6 Sol through Codex.

It burned tons of credits, but OpenAI kept resetting them before it had to. I would burn through my week's tokens, wait a few hours, and OpenAI would give me more credits to keep burning. It is addicting, like gambling. [codex-reset.com](https://codex-reset.com/) is a fun site that tracks the resets and predicts when it thinks another one is coming.

![My Codex usage down to 3% remaining](https://immich.soh.re/api/assets/ba1a6b29-9531-46ed-9e20-a2ffbed9d7dd/original?key=te_XqV_KVlnJ167prVYS86fw22xZkdwf_ny-FFpOPmhbHEdY8AEgLb2BAkLI3hZdMak)

Over ten days, Codex recorded **59.9 active hours and 35,549,156 tokens**. I prompted it, and then it began prompting me: reboot the device, hold different buttons while it starts, move cables around, play guitar, listen for distortion, take pictures, and finally take the whole thing apart.

![The AMPLIFi TT main board after I opened the case](https://immich.soh.re/api/assets/6fc8cbec-50b4-4567-a98c-167e775a7663/original?key=te_XqV_KVlnJ167prVYS86fw22xZkdwf_ny-FFpOPmhbHEdY8AEgLb2BAkLI3hZdMak)

I found it quite amusing when Codex prompted me to get a hobby knife and begin scraping glue off the microSD card that Line 6 surely did not want me unplugging. I was scared at first, but after watching it get stuck trying to figure out how to flash the firmware it was writing, I was convinced.

![The AMPLIFi TT microSD card still glued into its socket](https://immich.soh.re/api/assets/53041902-284a-4fdb-b809-e5d2d6298d88/original?key=te_XqV_KVlnJ167prVYS86fw22xZkdwf_ny-FFpOPmhbHEdY8AEgLb2BAkLI3hZdMak)

It was actually quite easy to take the device apart and begin scraping. Once the glue was gone, the microSD card popped out with extremely minimal force, just like it was intended to.

Companies ship the hardware, the firmware to run it, and the computer and mobile drivers used to talk to it. They have given you literally everything you need to understand the device. Now, with AI, you have the skills to use those resources and figure it out.

![A close look at the AMPLIFi TT's SHARC DSP and flash](https://immich.soh.re/api/assets/72dd5f54-78fa-4f2d-9238-9328d1d81456/original?key=te_XqV_KVlnJ167prVYS86fw22xZkdwf_ny-FFpOPmhbHEdY8AEgLb2BAkLI3hZdMak)

You cannot violate their licenses or publish their software, but you can brick and customize your own devices to work the way you want them to, limited only by the physical hardware and your LLM token budget.

## And it worked

I can now plug the TT into my Linux computer over USB and control the thing without touching the old mobile app. I can browse all 100 tones stored on it, switch between them, and edit the amps, cabinets, microphones, effects, and their settings. I can make a change in the GUI and hear it on the TT while I am playing.

I can back up the entire tone bank, import and export Line 6 presets, replace any of the 100 slots, and restore everything if I screw it up. The computer can hold as many tones as I want, so I am no longer limited to the 100 that fit on the box. I can also search Line 6's community tones and download them without fighting the mobile app.

I can switch between the processed guitar sound and a properly clean signal, mute the guitar without muting audio from the computer, and control the outputs the hardware supports. I can also identify the installed firmware and recover or reflash the device over USB. I had to pull the microSD card to get there, but I do not have to keep pulling it every time I want to change something.

The source, documentation, tests, and evidence I can redistribute are in [Jmainguy/line6-tt](https://github.com/Jmainguy/line6-tt). I am looking forward to improving the CLI and GUI and coming up with some cool tones now that the TT works the way I want it to. No longer do I need to throw away a device because the software is horrible. I can just write my own, suited to my needs.

![The Rust Line 6 TT control application showing the tone designer](https://immich.soh.re/api/assets/e0bd6064-1c56-42fe-b224-0accf3d595f6/original?key=te_XqV_KVlnJ167prVYS86fw22xZkdwf_ny-FFpOPmhbHEdY8AEgLb2BAkLI3hZdMak)
