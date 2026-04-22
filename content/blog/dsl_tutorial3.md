---
title: "Tutorial 3"
date: "2026-04-21"
author: "Oliver"
tags: ["linux", "kernel", "dsl"]
---

Hello, welcome to the third post about DSL26's tutorials!

## Fighting modules

In this tutorial, I've followed along on how to create a simple linux kernel module and how to create build configs for new kernel features. It may be useful to say that I've been following these tutorials without using **kworkflow**. Not because I have a problem with it, it's just because at the time I followed the second tutorial without using `kw`.

This was possibly the hardest tutorial to follow, yet. Mostly because now we switch to creating new things and putting them to test inside of the VM, a process that can get quite messy easily. As for the proposed exercises, I answered as follows:

1) When compiling with `allmodconfig`, It went for **3 hours and 50 minutes** of compilation time, and 34M for the Image.gz file, while running with `localmodconfig` with the custom list of needed modules, it went for 17 minutes of compilation time, and 15M size for Image.gz. Both kernels booted up correctly. 

The other two exercises were pretty simple / follow along, so no need to put them here.
