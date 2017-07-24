#include <linux/module.h>
#include <linux/init.h>
#include <linux/kernel.h>
#include <linux/timer.h> /* */
#include <linux/fs.h>
#include <linux/platform_device.h>
#include <linux/uio_driver.h>
#include <linux/slab.h> /* kmalloc, kfree */

struct uio_info uio_virtual_device_info = {
    .name = "virtual_uio",
    .version = "1.0",
    .irq = UIO_IRQ_NONE,
    // .mem = ?,
    // .port = ?,
};

static struct timer_list int_timer;

void emulated_intterupt( unsigned long data)
{
    struct uio_info * uinfo = (struct uio_info *)(data);
    printk("making intterupt to userspace\n");
    uio_event_notify(uinfo);
    int_timer.expires = jiffies + HZ;
    add_timer(&int_timer);
}

static int uio_virtual_device_drv_probe(struct platform_device *pdev)
{
    printk("initiating timer\n");
    init_timer(&int_timer);
    int_timer.function = emulated_intterupt;
    int_timer.data = (unsigned long)(&uio_virtual_device_info);
    int_timer.expires = jiffies + HZ;
    add_timer(&int_timer);
    printk("uio_virtual_device_probe( %p)\n", pdev);
    uio_virtual_device_info.mem[0].addr = (unsigned long)kmalloc(1024, GFP_KERNEL);

    if(uio_virtual_device_info.mem[0].addr == 0)
        return -ENOMEM;
    uio_virtual_device_info.mem[0].memtype = UIO_MEM_LOGICAL;
    uio_virtual_device_info.mem[0].size = 1024;

    printk("[%s,%d] uio_virtual_device_info.mem[0].addr:0x%llx, .size :%llu\n",\
            __func__,__LINE__,uio_virtual_device_info.mem[0].addr,\
            uio_virtual_device_info.mem[0].size);

    if(uio_register_device(&pdev->dev, &uio_virtual_device_info))
        return -ENODEV;

    return 0;
}

static int uio_virtual_device_drv_remove(struct platform_device *pdev)
{
    uio_unregister_device(&uio_virtual_device_info);

    return 0;
}

static struct platform_driver virtual_device_drv = {
    .probe  = uio_virtual_device_drv_probe,
    .remove = uio_virtual_device_drv_remove,
    .driver = {
        .name = "VIRTUAL_DEVICE",
        .owner = THIS_MODULE,
    }
};

static void virtual_device_remove(struct device *dev)
{
    kfree((const void *)(uio_virtual_device_info.mem[0].addr));
}

static struct platform_device virtual_device = {
    .name           = "VIRTUAL_DEVICE",
    .id             = -1,
    .dev            = {
        .release  = virtual_device_remove,
    },
};

static int __init uio_virtual_device_init(void)
{
    printk("virtual_device init ok!\n");
    platform_device_register(&virtual_device);

    printk("virtual_device_drv init ok!\n");
    return platform_driver_register(&virtual_device_drv);
}

static void __exit uio_virtual_device_exit(void)
{
    printk("del timer\n");
    del_timer(&int_timer);

    printk("virtual_device remove ok!\n");
    platform_device_unregister(&virtual_device);

    printk("virtual_device_drv remove ok!\n");
    platform_driver_unregister(&virtual_device_drv);
}

module_init(uio_virtual_device_init);
module_exit(uio_virtual_device_exit);

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Peng Gao");
MODULE_DESCRIPTION("Daemon of UIO");
