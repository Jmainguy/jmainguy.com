---
title: "Large Language Model: A Journey into AI and PC Upgrades"
date: 2023-10-30
draft: false
---

I found myself feeling unfulfilled on a Saturday, watching some rather lackluster football games. It was then that I decided I wanted to dive into the world of Large Language Models (LLMs), all thanks to a [fascinating article about Google investing $2 billion into a company I had never even heard of](https://news.ycombinator.com/item?id=38048155).

## Exploring LLMs with Hugging Face

In the past, I'd dabbled with using [Hugging Face](https://huggingface.co/) for generating AI-driven images, back when it was the hottest trend a few months ago. So, I naturally gravitated back to it. At the very top of their leader board was a model called [ShiningValiant](https://huggingface.co/ValiantLabs/ShiningValiant). However, it came with a significant requirement - downloading around 300GB of data. I thought, "Do I have enough disk space for this?"

## A Disk Space Dilemma

To my surprise, I discovered that I'd only partitioned 500GB of my 2TB NVMe disk. Without wasting any time, I used parted to delete the unused 200GB partition and resize the main one to a whopping 1TB. Then, it was time to tackle extending a btrfs LUKS encrypted partition. Here are the commands that worked for me:


```/bin/bash
# Find out the name of my active crypt
dmsetup status

# Resize it to the maximum partition size
cryptsetup resize luks-03bdbe6b-c97c-4541-959e-cf83ea6f7f31

# Extend the file system
btrfs filesystem resize max /
```

After installing llm and experimenting with various models, my PC suddenly froze because it was running low on RAM. 

I like to take a break when running into issues, as it helps resolve most problems when you come back with a fresh mindset. I started playing some guitar and before I knew it my kids had surrounded me with their instruments. Cora Beth was on the Keys, while Deacon was rocking out with his Ukelele he got for his birthday recently. We ran through a good set list.

* Smells Like Teen Spirit
* Dammit
* I'll Fly Away
* Star Spangled Banner
* Amazing Grace
* Brain Stew
* Die, Die My Darling
* Last Caress
  

## The RAM Revelation

When I came back to hacking, my quest to experiment with LLMs unveiled an unexpected revelation: I thought I had 32GB of RAM, but it turns out I only had 16GB. It had been a few years since I built this PC, my second ever, so it was definitely time for an upgrade. I turned to chatGPT for guidance on figuring out my motherboard and current RAM without cracking open the case:

```/bin/bash
# Find RAM details
sudo dmidecode --type 17

# Discover motherboard info
inxi -M
```

## The PC Part Picker Solution

With this newfound information, I headed to [PC Part Picker](https://pcpartpicker.com/) to see what upgrades were possible. It helped me select a [fantastic 64GB RAM kit](https://pcpartpicker.com/product/KdgQzy/corsair-vengeance-lpx-64-gb-4-x-16-gb-ddr4-3200-memory-cmk64gx4m4e3200c16) for approximately $155 after tax.

By the way, did you know that even when Amazon or other retailers do not collect sales tax for you, you're still responsible for it? Failing to pay your owed taxes can lead to a cumbersome monthly chore of filling out a [sales and use tax](https://www.ncdor.gov/taxes-forms/sales-and-use-tax) form. Companies that collect sales tax are actually doing you a favor and saving you from this burden. If you do not like sales tax, take it up with your elected representatives.

## Waiting for the RAM and Future Endeavors

While eagerly awaiting the arrival of my new RAM, I continued to experiment with LLMs for the next couple of days. I attempted to use ROCm/PyTorch with my Radeon card (gfx803), but I soon learned that my card was considered too old and no longer supported. Trying to build the necessary software myself proved challenging and even if I succeeded, it would have been a temporary hack.

I enjoyed playing with [text-generation-webui](https://github.com/oobabooga/text-generation-webui/tree/main) and was able to get it working on cpu, but the speeds were horrendous. 

[chat-ui](https://github.com/huggingface/chat-ui) looks interesting and I might play with that if I ever get good speeds with these models.

## Summary
In 2023, the consensus seems to be that for AI/ML, you should be using an Nvidia GPU. Their 40x series appears to be outstanding, with some people even vouching for the 30x series if you can find a good deal. Personally, I'm getting an M2 Macbook from work in about a month, and it's rumored to be quite capable for AI/ML. So, I'll be revisiting my AI experiments without having to spend more on a GPU for the time being.
