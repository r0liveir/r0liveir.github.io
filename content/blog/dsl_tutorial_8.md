---
title: "Tutorial 5 to 8"
date: "2026-05-10"
author: "Oliver"
tags: ["linux", "kernel", "dsl"]
---

Hello, welcome to the last post about DSL26 tutorials!

## Some notes

Since tutorials 5 and 6 were focused on the "Sending patches via email" pipeline, I've decided not to make dedicated posts to them. Furthermore, tutorial 7 is, essentially, a showcase of the `iio_dummy` structure which we play with in tutorial 8, so naturally I didn't make a post for it as well.

## Slashing dummies

For this tutorial, I've followed along on installing the `iio_dummy` module in order to test changes to it. During this tutorial, we've added 3 additional channels to represent three axis.

Since the post's date, there has been some changes made to `iio_simple_dummy.c` file. When updating `*read_raw()` for handling new channels, i've found that the logic for reading raw channel data has been split into into a helper function called `__iio_dummy_read_raw()`. And this is where i've put the code for making data accessible to userspace, like this:

```C 
static int __iio_dummy_read_raw(struct iio_dev *indio_dev,
				struct iio_chan_spec const *chan,
				int *val)
{
	struct iio_dummy_state *st = iio_priv(indio_dev);

	guard(mutex)(&st->lock);
	switch (chan->type) {
	case IIO_VOLTAGE:
        ...
	case IIO_MAGN:
		switch(chan->scan_index) {
		case DUMMY_MAGN_X:
			*val = st->buffer_compass[0];
			break;
		case DUMMY_MAGN_Y:
			*val = st->buffer_compass[1];
			break;
		case DUMMY_MAGN_Z:
			*val = st->buffer_compass[2];
			break;
		default:
			return -EINVAL;
		}
		return IIO_VAL_INT;
	default:
		return -EINVAL;
	}
}

```

Notice that since this function doesn't use `int ret` like the original, we return `IIO_VAL_INT` directly inside the switch statement. `iio_dummy_read_raw` now looks like this:

```C 
static int iio_dummy_read_raw(struct iio_dev *indio_dev,
			      struct iio_chan_spec const *chan,
			      int *val,
			      int *val2,
			      long mask)
{
	struct iio_dummy_state *st = iio_priv(indio_dev);
	int ret;

	switch (mask) {

        ...

	case IIO_CHAN_INFO_SCALE:
		switch (chan->type) {
		case IIO_VOLTAGE:

            ...

		case IIO_MAGN:
			*val = 0;
			*val2 = 2;
			ret = IIO_VAL_INT_PLUS_MICRO;
			break;
		default:
			return -EINVAL;
		}

```

and with that, these new additions can be succesfully compiled and tested!
