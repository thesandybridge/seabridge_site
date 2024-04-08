# Recover Grub
This guide is to get you out of the grub recovery view. This happened to me when trying to configure dual boot with Windows and Arch. Installing Windows first before Arch can prevent this issue from potentially occurring but this is good to know just in case.

This can also happen if the boot drive gets changed somehow, such as changing partition names or moving them around.

1. [List Drives](#list-drives)
2. [Set Root and Prefix](#set-root-and-prefix)
3. [Insmod](#insmod)
4. [Reconfigure Grub](#reconfigure-grub)

## List drives

```bash
ls # you will see a list of drives, such as (hd1, gpt1)(hd2, gpt2) ...
```

Run that command on each drive to see which one contains the Linux file system.

## Set root and prefix

```bash
set root=(hd2,gpt3) # Enter the drive you discovered in the first step
set prefix=(hd2,gpt3)/boot/grub # Enter the path to grub
```

## Insmod

Running the following commands should launch Linux.

```bash
insmod normal
normal
```

## Reconfigure grub

Install the bootloader

```bash
grub-install --target=x86_64-efi --bootloader-id=Arch
```

( you can name the bootloader-id to whatever you want, this is the label that shows in grub and your bios.)

Generate the grub.cfg file.

```bash
grub-mkconfig -o /boot/grub/grub.cfg
```

