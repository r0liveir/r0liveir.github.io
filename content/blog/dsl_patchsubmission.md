---
title: "Patch Submission"
date: "2026-06-17"
author: "Oliver"
tags: ["linux", "kernel", "dsl"]
---

Hello, welcome to the last post of phase 1 of the Free Software Development Course.

## Patch work

For this patch, my colleague and I decided to take a code duplication issue found under `drm/amdgpu`, which, under RiseUp issues list, is **Duplication number 8**, for the following two functions, with 100% similarity:
- `gmc_v10_vm_fault_interrupt_state`
- `gmc_v12_0_vm_fault_interrupt_state`

These functions are related to the Graphic Memory Controller (GMC) part, while v10 and v12_0 refer to two controller versions. Taking a closer look at both functions:

```C 
static int
gmc_v10_0_vm_fault_interrupt_state(struct amdgpu_device *adev,
				   struct amdgpu_irq_src *src, unsigned int type,
				   enum amdgpu_interrupt_state state)
{
	switch (state) {
	case AMDGPU_IRQ_STATE_DISABLE:
		/* MM HUB */
		amdgpu_gmc_set_vm_fault_masks(adev, AMDGPU_MMHUB0(0), false);
		/* GFX HUB */
		/* This works because this interrupt is only
		 * enabled at init/resume and disabled in
		 * fini/suspend, so the overall state doesn't
		 * change over the course of suspend/resume.
		 */
		if (!adev->in_s0ix)
			amdgpu_gmc_set_vm_fault_masks(adev, AMDGPU_GFXHUB(0), false);
		break;
	case AMDGPU_IRQ_STATE_ENABLE:
		/* MM HUB */
		amdgpu_gmc_set_vm_fault_masks(adev, AMDGPU_MMHUB0(0), true);
		/* GFX HUB */
		/* This works because this interrupt is only
		 * enabled at init/resume and disabled in
		 * fini/suspend, so the overall state doesn't
		 * change over the course of suspend/resume.
		 */
		if (!adev->in_s0ix)
			amdgpu_gmc_set_vm_fault_masks(adev, AMDGPU_GFXHUB(0), true);
		break;
	default:
		break;
	}

	return 0;
}

```

as for the gmc v12 version, it is mostly the same, with one little detail (the `type` parameter here is being passed as `unsigned` instead of `unsigned int`):

```C
static int gmc_v12_0_vm_fault_interrupt_state(struct amdgpu_device *adev,
					      struct amdgpu_irq_src *src, unsigned type,
					      enum amdgpu_interrupt_state state)
{
...
}
```

taking a look at where in these files the functions are used, we can see something akin to:

```C
static const struct amdgpu_irq_src_funcs gmc_v12_0_irq_funcs = {
	.set = gmc_v12_0_vm_fault_interrupt_state,
	.process = gmc_v12_0_process_interrupt,
};

```

and similar for gmc v10

both types are essentially the same in C programming, thus, we can refactor this properly. Our proposed steps were then:
- Move `vm_fault_interrupt_state` functions into a helper `amdgpu_gmc_vm_fault_interrupt_state`, found under `amdgpu_gmc.c` file 
- Update both `irq_funcs` from v10 and v12.0 to handle the new helper function
- Update `amdgpu_gmc.h` header file to insert this new function, so that both files can make use of it 

# Viewing the solution

This is how the solution turned out to be. Under `amdgpu_gmc.c`:

```C 
int amdgpu_gmc_vm_fault_interrupt_state(struct amdgpu_device *adev,
				   struct amdgpu_irq_src *src, unsigned int type,
				   enum amdgpu_interrupt_state state)
{
	switch (state) {
	case AMDGPU_IRQ_STATE_DISABLE:
		/* MM HUB */
		amdgpu_gmc_set_vm_fault_masks(adev, AMDGPU_MMHUB0(0), false);
		/* GFX HUB */
		/* This works because this interrupt is only
		 * enabled at init/resume and disabled in
		 * fini/suspend, so the overall state doesn't
		 * change over the course of suspend/resume.
		 */
		if (!adev->in_s0ix)
			amdgpu_gmc_set_vm_fault_masks(adev, AMDGPU_GFXHUB(0), false);
		break;
	case AMDGPU_IRQ_STATE_ENABLE:
		/* MM HUB */
		amdgpu_gmc_set_vm_fault_masks(adev, AMDGPU_MMHUB0(0), true);
		/* GFX HUB */
		/* This works because this interrupt is only
		 * enabled at init/resume and disabled in
		 * fini/suspend, so the overall state doesn't
		 * change over the course of suspend/resume.
		 */
		if (!adev->in_s0ix)
			amdgpu_gmc_set_vm_fault_masks(adev, AMDGPU_GFXHUB(0), true);
		break;
	default:
		break;
	}

	return 0;
}
```

Under `amdgpu_gmc.h`:

```C
int amdgpu_gmc_vm_fault_interrupt_state(struct amdgpu_device *adev,
				   struct amdgpu_irq_src *src, unsigned int type,
				   enum amdgpu_interrupt_state state);
```

And the two imports for the struct looked like:

```C 
 static const struct amdgpu_irq_src_funcs gmc_v12_0_irq_funcs = {
	.set = gmc_v12_0_vm_fault_interrupt_state,
	.set = amdgpu_gmc_vm_fault_interrupt_state,
 	.process = gmc_v12_0_process_interrupt,
 };
```

# Submitting the patch 

With FLUSP Tutorial, environment configuration and patch submission with kw became pretty easy to understand and to implement. Following this, we've sent our patch in May, 20. So far, we heard no responses from the manteiners, thus i've sent out another patch with subject `PATCH RESEND` to ping again and, hopefully, hear feedback on our patch. We also plan on sending more patches in the future. 

# Overall experience

This initial experience with Linux contribution was pretty great! I look forward to submitting more patches in the future and, if an opportunity arises, contribute more with the Kernel and work with it professionally.
