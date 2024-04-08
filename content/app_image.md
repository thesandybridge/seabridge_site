# Setup Up AppImage
This guide will walk through the steps to install and configure an AppImage for use on a linux machine. This guide assumes you have already downloaded the app image.

## Configure a desktop file

Navigate to `~/.local/share/applications`.

From here you can create a new directory with the same name as the application. You can add a `.desktop` file here with the following.

```bash
[Desktop Entry]
Exec=/opt/app_name/file_name.AppImage
Type=Application
Categories=Something
Name=AppName
```

Now you will need to move AppImage to the opt directory we just showed above.

```bash
cd /opt
mkdir app_name
sudo ~/Downloads/file_name.AppImage /app_name
cd app_name
chmod +x file_name.AppImage
```

Now you you can use the application, this will show up in dmenu and rofi, or if you are not using i3wm it will show up on your desktop. You can add an icon to that desktop file we made earlier.
