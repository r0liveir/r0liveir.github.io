---
title: "Tutorial 1"
date: "2026-04-21"
author: "Oliver"
tags: ["linux", "kernel", "dsl"]
---

Hello, welcome to the first post about DSL26's tutorials!

## And so the journey begins

The first tutorial describes setting up a Linux kernel test environment using `QEMU`. The first time I tried it (during vacation), it was a mess... but mostly because I forgot to set the working directory to be `/home/lk-dev`, and instead it was in my user's home directory, which (to no suprise) messed up permissions, and before I realized, I had a VM with no connection.

Other than that, it was pretty fine. Sure, the amount of new info thrown at you when first diving into this world is huuge, but FLUSP's team made sure to give a welcoming course, although I had some issues:

- The link for **nocloud debian image** wasn't working properly, so I had to get one by myself in `cloud.debian.org`. Specifically, I used [this one](https://cloud.debian.org/images/cloud/bookworm/20260210-2384/debian-12-nocloud-arm64-20260210-2384.qcow2)
- During step **2.4 (use libvirst to streamline managing VMs)**, i got an error regarding the `create_vm_virsh()` function, in which it said that `/dev/<vdaX>` (in my case, `vda2`) was not found. Searching more through this, i've found that this was due to a mismatch on how QEMU and libvirt handle default hardware. At the line `--disk path="${VM_DIR}/arm64_img.qcow2"`, it wasn't told to use the `virtio` bus with `vda` as device target, which means the guest OS was seeing it under a different name (like `/dev/sda`). So i changed this line to `--disk path="${VM_DIR}/arm64_img.qcow2",bus=virtio,target=vda`.
- When configuring SSH, `sshd` service was not found. This was due to SSH service not being installed by default on this image, which was easily fixed by installing the `openssh-server` with `apt`.

At last, we fetch the list of modules loaded in the guest kernel in order to save into a file and reduce Linux build time for the next tutorial.

Thanks for reading!
