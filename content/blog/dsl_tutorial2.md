---
title: "Tutorial 2"
date: "2026-04-21"
author: "Oliver"
tags: ["linux", "kernel", "dsl"]
---

Hello, welcome to the second post about DSL26's tutorials!

## Wandering kernel lands

In this tutorial, the premise was to compile the kernel using **cross-compilation** for the ARM architecture, and booting this custom kernel in a VM.

In order to reduce build time, in Tutorial 1 we got the `vm_mod_list` file to make a minimal set of modules for compilation. Also, for the **append to kernel release** option, I went with `softwarelivre26`. 

It went pretty well, with almost no errors. During the proccess I've madesome silly mistakes, such as forgetting to **install modules inside the VM**, using `make ... modules_install`. But other than that, it went smoothly.

