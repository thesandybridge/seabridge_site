# Arch Linux Installation Guide

The following guide assumes you have installed Arch on a a bootable USB and have already booted into the drive. The official Arch wiki has a great guide, but it doesn't go into specifics and can be vague. I have attempted to dive into the specifics and include some solutions to problems I ran into.

This guide will walk through how to configure Arch using Grub as a bootloader, Gnome display manager and i3wm for the desktop environment.

## Table of Contents

1. [Check Network Connection](#check-network-connection)
2. [Partition the disks](#partition-the-disks)
3. [Format the partitions](#format-the-partitions)
	1. [Mount the file systems](#mount-the-file-systems)
4. [Install Linux](#install-linux)
5. [Configure the system](#configure-the-system)
	1. [Localization](#localization)
	2. [Network configuration](#network-configuration)
	3. [Root](#root)
6. [Install Grub Bootloader](#install-grub-bootloader)
7. [Install a desktop environment](#install-a-desktop-environment)
8. [Create a user](#create-a-normal-user)
9. [Install i3wm](#install-i3wm-alternative-desktop-environment)
10. [Complete](#complete)

## Check Network Connection
- Check to make sure you are connected to the internet.
	`$ ip link`  you can verify you are connected by `$ ping archlinux.org`

- Ethernet should automatically connect, if you are using wireless run [iwctl](https://wiki.archlinux.org/title/Iwctl)

## Partition the disks

The following assumes you are installing Arch on an unallocated drive with Free space.

`$ fdisk -l` will list your drives. Find the drive you want to install arch to and then make note of this in the following steps.
	for example, the drive might say /dev/sdb

`$ cfdisk /dev/the_disk_to_be_partitioned`
- This will open a disk partition tool that will save the headache of setting up drives. You can manually do this with fdisk as well.

You will need to create 3 partitions. Replace sdb* with your partition names.

If you are configuring Arch on a VM use fdisk to create the partitions. You don’t need the EFI partition.

| Partition         | Partition Type     | Size|
|--------------|-----------|------------|
| /dev/sdb1 | EFI system partition      | AT least 300M       |
| /dev/sdb2      | Linux swap | Half of your RAM up to 4GB |
| /dev/sdb3      | Linux Filesystem (root) | Remaining space |


## Format the partitions

Once you have successfully created the partitions and wrote them to the drive continue with the following. (The guide will continue with the partition naming convention as defined in the examples above.)

- Create an Ext4 file system on the root partition
- Create the swap
- Format the EFI system

```bash
mkfs.ext4 /dev/sdb3
mkswap /dev/sdb2
mkfs.fat -F 32 /dev/sdb1
```

### Mount the file systems

- Mount the root volume to `/mnt`.
- Configure the swap.

```bash
mount /dev/sdb3 /mnt
swapon /dev/sdb2
```

## Install Linux

Install Linux using [pacstrap](https://man.archlinux.org/man/pacstrap.8) to the mounted file system.
```bash
pacstrap -K /mnt base linux linux-firmware
```

## Configure the system

- Generate an [fstab](https://wiki.archlinux.org/title/Fstab) file.
```bash
genfstab -U /mnt >> /mnt/etc/fstab
```

- Change into the root of the new system
```bash
arch-chroot /mnt
```

- Set the timezone and configure the clock

```bash
ln -sf /usr/share/zoneinfo/America/New_York /etc/localtime
hwclock --systohc
```

- Install sudo and vim

```bash
pacman -S sudo vim
```

### Localization

```bash
vim /etc/locale.gen
```

This will output a list of locals, scroll down and uncomment the two that precede with `en_US`.

- Generate the locales with `locale-gen`.
- Create the locale.conf

```bash
echo LANG=en_US.UTF-8 > /etc/locale.conf
export LANG=en_US.UTF-8
```

### Network configuration

- Create the hostname file `$ echo myhostname > /etc/hostname`.
- Configure the hosts file `$ vim /etc/hosts`.

```bash
127.0.0.1        localhost
::1              localhost
127.0.1.1        myhostname
```

### Root

- Set a root password with `$ passwd`

## Install Grub bootloader

```bash
pacman -S grub efibootmgr os-prober mtools
```

- Create a mount point for your efi partition and mount it.
```bash
mkdir /boot/efi
mount /dev/sdb1 /boot/efi
```

- Install the bootloader
```bash
grub-install --target=x86_64-efi --bootloader-id=Arch
```
( you can name the bootloader-id to whatever you want, this is the label that shows in grub and your bios.)
- Generate the grub.cfg file.
```bash
grub-mkconfig -o /boot/grub/grub.cfg
```

## Install a desktop environment

- Install the X environment
```bash
pacman -S xorg-server xorg-apps
```

- Install device drivers (you don't need all three. choose the appropriate one for your system).
```bash
pacman -S nvidia nvidia-utils
pacman -S xf86-video-ati
pacman -S xf86-video-intel
```

- Install gnome or whichever desktop display manager and environment you prefer.
```bash
pacman -S gnome gnome-extra networkmanager
```

- enable your display manager and network manager
```bash
systemctl enable gdm
systemctl enable NetworkManager
```

## Create a normal user

```bash
useradd -m -G wheel youruser
passwd youruser
```

After the user has been created you will need to enable the `wheel` permission group.

```bash
visudo
```

Scroll down and uncomment the line `%wheel ALL=(ALL) ALL`.

## Install i3wm (alternative desktop environment)

```bash
pacman -S i3-wm i3status
```

## Complete

You have fully configured the environment and are ready to reboot. Follow the steps below.

```bash
exit
umount -R /mnt
reboot
```
