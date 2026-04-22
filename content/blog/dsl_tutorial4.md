---
title: "Tutorial 4"
date: "2026-04-21"
author: "Oliver"
tags: ["linux", "kernel", "dsl"]
---

Hello, welcome to the fourth post about DSL26's tutorials!

## Entering dungeons

The 4th tutorial was used to introduce character devices, abstractions that provide an interface for interacting with hardware or software components as a stream of files, and are located at the `/dev` directory. They're essentially files for communicating with hardware drivers, one byte at a time. 

I can say for sure that that this was the hardest tutorial to follow, with exercises being even harder to understand (where to start? what to do?). The follow-along part was pretty straightforward, but still hard to understand what was I doing. 

For the proposed exercises, the first one could be solved by employing the following line in the `simple_char_open` function: 

```C
pr_info("Device opened - Major %u and minor %u\n", imajor(inode), iminor(inode));
```

Which would then print the major and minor device numbers. The second one was very hard, but eventually I found a solution (with a lil help of AI) that proposed the use of structs like the following:

```C
struct device_data {
	struct cdev s_cdev;
	char *s_buf;
	size_t size;
	struct mutex lock;
};

```

Which would then be used to implement more than one device number and make it keep separate buffers for each of them. My solution was something like this:

```C
#include <linux/init.h>
#include <linux/module.h>

#include <linux/kdev_t.h> /* for MAJOR and MINOR */
#include <linux/cdev.h> /* for cdev */
#include <linux/fs.h> /* for chrdev functions */
#include <linux/slab.h> /* for malloc */
#include <linux/string.h> /* for strlen() */
#include <linux/uaccess.h> /* copy_to_user() */
#include <linux/mutex.h> // For mutexes/locks

// Ex2: define a container struct for keeping separate buffers
struct device_data {
	struct cdev s_cdev;
	char *s_buf;
	size_t size;
	struct mutex lock;
};

#define S_BUFF_SIZE 4096
#define MINOR_NUMS 2

static dev_t dev_id;
struct device_data *char_devices;

static int simple_char_open(struct inode *inode, struct file *file)
{
	struct device_data *data;
	data = container_of(inode->i_cdev, struct device_data, s_cdev);
	file->private_data = data;

	// Exercise 1: modify this to print major and minor device numbers on device open
	pr_info("Device opened - Major %u and minor %u\n", imajor(inode), iminor(inode));

	pr_info("%s: %s\n", KBUILD_MODNAME, __func__);
	return 0;
}

static ssize_t simple_char_read(struct file *file, char __user *buffer,
				size_t count, loff_t *ppos)
{
	struct device_data *data = file->private_data;
	int n_bytes;

	if (mutex_lock_interruptible(&data->lock)) return -ERESTARTSYS;
	pr_info("%s: %s about to read %ld bytes from buffer position %lld\n",
		KBUILD_MODNAME, __func__, count, *ppos);

	n_bytes = count - copy_to_user(buffer, data->s_buf + *ppos, count);
	*ppos += n_bytes;

	mutex_unlock(&data->lock);
	return n_bytes;
}

static ssize_t simple_char_write(struct file *file, const char __user *buffer,
				size_t count, loff_t *ppos)
{
	struct device_data *data = file->private_data;
	int n_bytes;

	if (mutex_lock_interruptible(&data->lock)) return -ERESTARTSYS;

	pr_info("%s: %s about to write %ld bytes to buffer position %lld\n",
		KBUILD_MODNAME, __func__, count, *ppos);

	n_bytes = count - copy_from_user(data->s_buf + *ppos, buffer, count);
	mutex_unlock(&data->lock);

	return n_bytes;
}

static int simple_char_release(struct inode *inode, struct file *file)
{
	pr_info("%s: %s\n", KBUILD_MODNAME, __func__);
	return 0;
}

static const struct file_operations simple_char_fops = {
	.owner = THIS_MODULE,
	.open = simple_char_open,
	.release = simple_char_release,
	.read = simple_char_read,
	.write = simple_char_write,
};

static int __init simple_char_init(void)
{
	int ret;

	pr_info("Initialize %s module.\n", KBUILD_MODNAME);

	ret = alloc_chrdev_region(&dev_id, 0, MINOR_NUMS, "simple_char");
	if (ret < 0)
		return ret;
	
	char_devices = kmalloc_array(MINOR_NUMS, sizeof(struct device_data), GFP_KERNEL);
	if (!char_devices) {
		unregister_chrdev_region(dev_id, MINOR_NUMS);
		return -ENOMEM;
	}
	// Exercise 2
	for (int i = 0; i < MINOR_NUMS; i++) {
		dev_t dev_no = MKDEV(MAJOR(dev_id), i);

		char_devices[i].s_buf = kmalloc(S_BUFF_SIZE, GFP_KERNEL);

		strcpy(char_devices[i].s_buf, "This is data from simple_char buffer.");

		mutex_init(&char_devices[i].lock);

		cdev_init(&char_devices[i].s_cdev, &simple_char_fops);
		char_devices[i].s_cdev.owner = THIS_MODULE;
		cdev_add(&char_devices[i].s_cdev, dev_no, 1);
	}	
	
	// For now, return 0
	return 0;
}

static void __exit simple_char_exit(void)
{
	/*
	 * Undoes the device ID mapping and frees cdev struct, removing the
	 * character device from the system.
	 */
	for (int i = 0; i < MINOR_NUMS; i++) {
		cdev_del(&char_devices[i].s_cdev);
		kfree(char_devices[i].s_buf);
	}
	kfree(char_devices);
	/* Unregisters (disassociate) the device numbers allocated. */
	unregister_chrdev_region(dev_id, MINOR_NUMS);
	pr_info("%s exiting.\n", KBUILD_MODNAME);
}

module_init(simple_char_init);
module_exit(simple_char_exit);

MODULE_AUTHOR("A Linux kernel student <name.surname@usp.br>");
MODULE_DESCRIPTION("A simple character device driver example.");
MODULE_LICENSE("GPL");

```

Although it would be nice to say it worked, when I tested it, It gave some errors. But hey, it's a start! In the future I'll probably come back to this exercise.
